// Deprecated: This package is deprecated and should not be used anymore. Please use the github.com/go-puzzles/puzzles/goredis instead.
package predis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	redisDialer "github.com/go-puzzles/puzzles/dialer/redis"
)

var (
	defaultTTL = time.Second * 10
)

type RedisClient struct {
	pool       *redis.Pool
	lockValues sync.Map
}

func NewRedisClient(pool *redis.Pool) *RedisClient {
	return &RedisClient{
		pool: pool,
	}
}

func NewRedisClientWithAddr(addr string, db int, maxIdle int, password ...string) *RedisClient {
	pool := redisDialer.DialRedisPool(addr, db, maxIdle, password...)
	return NewRedisClient(pool)
}

func (rc *RedisClient) GetPool() *redis.Pool {
	return rc.pool
}

func (rc *RedisClient) GetConn() redis.Conn {
	conn, _ := rc.GetConnWithContext(context.TODO())
	return conn
}

func (rc *RedisClient) GetConnWithContext(ctx context.Context) (redis.Conn, error) {
	return rc.pool.GetContext(ctx)
}

func (rc *RedisClient) Do(command string, args ...any) (reply any, err error) {
	conn := rc.GetConn()
	defer conn.Close()

	return conn.Do(command, args...)
}

func (rc *RedisClient) stringToAny(datas []string) []any {
	resp := make([]any, 0, len(datas))
	for _, data := range datas {
		resp = append(resp, data)
	}
	return resp
}

func (rc *RedisClient) DoWithTransactionPipeline(watchKey []string, commands ...[]any) error {
	conn := rc.GetConn()
	defer conn.Close()

	return rc.TransactionPipeline(conn, watchKey, commands...)
}

func (rc *RedisClient) TransactionPipeline(conn redis.Conn, watchKey []string, commands ...[]any) error {
	if len(watchKey) != 0 {
		_, err := conn.Do("WATCH", rc.stringToAny(watchKey)...)
		if err != nil {
			return errors.Wrap(err, "watchKey")
		}
	}

	if err := conn.Send("MULTI"); err != nil {
		log.Fatalf("Failed to send MULTI: %v", err)
	}

	for _, command := range commands {
		commandName := command[0].(string)
		args := command[1:]
		if err := conn.Send(commandName, args...); err != nil {
			return errors.Errorf("send command: %v args: %v error: %v", commandName, args, err)
		}
	}

	if err := conn.Flush(); err != nil {
		return errors.Wrap(err, "flush")
	}

	if _, err := conn.Do("EXEC"); err != nil {
		return errors.Wrap(err, "exec")
	}

	return nil
}

func (rc *RedisClient) SetWithTTL(key string, value any, ttl time.Duration) error {
	var err error
	switch value.(type) {
	case string, int, int64, float32, float64, bool, []byte:
	default:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("json marshal failed: %w", err)
		}
		value = jsonBytes
	}

	if ttl > 0 {
		_, err = rc.Do("SET", key, value, "EX", int(ttl.Seconds()))
	} else {
		_, err = rc.Do("SET", key, value)
	}

	return errors.Wrap(err, "redis.Set")
}

func (rc *RedisClient) Set(key string, value any) error {
	return rc.SetWithTTL(key, value, 0)
}

func (rc *RedisClient) Get(key string, out any) error {
	reply, err := rc.Do("GET", key)
	if err != nil {
		return errors.Wrap(err, "redis.GET")
	}
	switch ptr := out.(type) {
	case *string:
		*ptr, err = redis.String(reply, nil)
	case *[]byte:
		*ptr, err = redis.Bytes(reply, nil)
	case *int:
		*ptr, err = redis.Int(reply, nil)
	case *int64:
		*ptr, err = redis.Int64(reply, nil)
	case *float32:
		f64, err := redis.Float64(reply, nil)
		if err == nil {
			*ptr = float32(f64)
		}
	case *float64:
		*ptr, err = redis.Float64(reply, nil)
	case *bool:
		*ptr, err = redis.Bool(reply, nil)
	case *time.Time:
		return errors.New("unsupported type: time.Time")
	default:
		var b []byte
		b, err = redis.Bytes(reply, nil)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, out)
	}

	return err
}

func (rc *RedisClient) Delete(key string) error {
	_, err := rc.Do("DEL", key)
	return err
}

// LockWithBlock attempts to acquire a lock with the given key in Redis,
// retrying up to a specified maximum number of attempts (maxRetry).
// If the lock cannot be acquired, it will wait for 500 milliseconds
// before retrying, as long as the error returned is ErrLockFailed.
//
// The lock expiration time can be specified as a variadic argument;
// if not provided, a default expiration time will be used.
//
// If the lock is successfully acquired, the method returns nil.
// If the maximum number of retries is reached without acquiring the lock,
// it returns ErrLockFailed.
//
// Parameters:
// - key: The key under which the lock is to be stored.
// - maxRetry: The maximum number of retry attempts to acquire the lock.
// - expires: Optional duration(s) for which the lock should be valid.
//
// Returns:
//   - An error if the lock could not be acquired after maxRetry attempts,
//     or if another error occurred during the lock acquisition process.
func (rc *RedisClient) LockWithBlock(key string, maxRetry int, expires ...time.Duration) error {
	var lastErr error
	for i := 0; i < maxRetry; i++ {
		err := rc.Lock(key, expires...)
		if err == nil {
			return nil
		}

		lastErr = err
		if errors.Is(err, ErrLockFailed) {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		return err
	}

	return lastErr
}

// Lock attempts to acquire a lock with the given key in Redis.
// If the lock is successfully acquired, it sets an expiration time.
// The expiration time can be specified as a variadic argument; if not provided,
// a default expiration time (defaultTTL) will be used.
//
// The lock is acquired using the Redis SET command with the "NX" option,
// which ensures that the lock is only set if the key does not already exist.
//
// If the lock is already held (i.e., the key exists), the method returns
// ErrLockFailed. If any other error occurs during the operation, it is returned.
//
// Parameters:
// - key: The key under which the lock is to be stored.
// - expires: Optional duration(s) for which the lock should be valid.
//
// Returns:
// - An error, if the lock could not be acquired or another error occurred.
func (rc *RedisClient) Lock(key string, expires ...time.Duration) error {
	expire := defaultTTL
	if len(expires) != 0 {
		expire = expires[0]
	}

	value := randomValue()

	reply, err := rc.Do("SET", key, value, "NX", "PX", int(expire.Milliseconds()))
	if err != nil {
		return err
	}
	if reply == nil {
		return ErrLockFailed
	}

	rc.lockValues.Store(key, value)
	return nil
}

func (rc *RedisClient) UnLock(key string) error {
	value, ok := rc.lockValues.Load(key)
	if !ok {
		return ErrLockNotHeld
	}

	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	conn := rc.GetConn()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("EVAL", script, 1, key, value))
	if err != nil {
		return err
	}

	rc.lockValues.Delete(key)

	if reply == 0 {
		return ErrLockNotHeld
	}
	return nil
}

// LPush pushes one or more values to the head of the list stored at key
// If the list does not exist, it is created as an empty list before performing the push operations
//
// Parameters:
// - key: The key of the list
// - values: One or more values to push to the list
//
// Returns:
// - int: The length of the list after the push operations
// - error: An error if the operation failed
func (rc *RedisClient) LPush(key string, values ...any) (int, error) {
	args := make([]any, 0, len(values)+1)
	args = append(args, key)

	for _, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return 0, errors.Wrap(err, "encode value")
		}
		args = append(args, data)
	}

	reply, err := redis.Int(rc.Do("LPUSH", args...))
	if err != nil {
		return 0, errors.Wrap(err, "redis.LPUSH")
	}
	return reply, nil
}

// RPush appends one or more values to the tail of the list stored at key
// If the list does not exist, it is created as an empty list before performing the push operations
//
// Parameters:
// - key: The key of the list
// - values: One or more values to push to the list
//
// Returns:
// - int: The length of the list after the push operations
// - error: An error if the operation failed
func (rc *RedisClient) RPush(key string, values ...any) (int, error) {
	args := make([]any, 0, len(values)+1)
	args = append(args, key)

	for _, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return 0, errors.Wrap(err, "encode value")
		}
		args = append(args, data)
	}

	reply, err := redis.Int(rc.Do("RPUSH", args...))
	if err != nil {
		return 0, errors.Wrap(err, "redis.RPUSH")
	}
	return reply, nil
}

// LPop removes and returns the first element of the list stored at key
// If the list does not exist or is empty, it returns an error
//
// Parameters:
// - key: The key of the list
// - out: A pointer to the variable where the popped value will be stored
//
// Returns:
// - error: An error if the operation failed or if the list is empty
func (rc *RedisClient) LPop(key string, out any) error {
	data, err := redis.Bytes(rc.Do("LPOP", key))
	if err == redis.ErrNil {
		return errors.New("list is empty")
	}
	if err != nil {
		return errors.Wrap(err, "redis.LPOP")
	}

	if err := json.Unmarshal(data, out); err != nil {
		return errors.Wrap(err, "decode value")
	}
	return nil
}

// RPop removes and returns the last element of the list stored at key
// If the list does not exist or is empty, it returns an error
//
// Parameters:
// - key: The key of the list
// - out: A pointer to the variable where the popped value will be stored
//
// Returns:
// - error: An error if the operation failed or if the list is empty
func (rc *RedisClient) RPop(key string, out any) error {
	data, err := redis.Bytes(rc.Do("RPOP", key))
	if err == redis.ErrNil {
		return errors.New("list is empty")
	}
	if err != nil {
		return errors.Wrap(err, "redis.RPOP")
	}

	if err := json.Unmarshal(data, out); err != nil {
		return errors.Wrap(err, "decode value")
	}
	return nil
}

// LRange returns the specified elements of the list stored at key
// The offsets start and stop are zero-based indexes
// These offsets can also be negative numbers indicating offsets starting at the end of the list
// For example, -1 is the last element of the list, -2 the penultimate, and so on
//
// Parameters:
// - key: The key of the list
// - start: The starting position (inclusive)
// - stop: The ending position (inclusive)
// - out: A pointer to the slice where the range of elements will be stored
//
// Returns:
// - error: An error if the operation failed
func (rc *RedisClient) LRange(key string, start, stop int, out any) error {
	reply, err := redis.ByteSlices(rc.Do("LRANGE", key, start, stop))
	if err != nil {
		return errors.Wrap(err, "redis.LRANGE")
	}

	// Create a temporary slice to store decoded data
	var temp []any
	for _, item := range reply {
		var value any
		if err := json.Unmarshal(item, &value); err != nil {
			return errors.Wrap(err, "decode value")
		}
		temp = append(temp, value)
	}

	// Encode the temporary slice to JSON and then decode it into the output slice
	// This ensures proper type conversion
	encoded, err := json.Marshal(temp)
	if err != nil {
		return errors.Wrap(err, "encode temporary slice")
	}

	if err := json.Unmarshal(encoded, out); err != nil {
		return errors.Wrap(err, "decode to output slice")
	}

	return nil
}

// LLen returns the length of the list stored at key
// If the key does not exist, it is interpreted as an empty list and 0 is returned
//
// Parameters:
// - key: The key of the list
//
// Returns:
// - int: The length of the list at key
// - error: An error if the operation failed
func (rc *RedisClient) LLen(key string) (int, error) {
	reply, err := redis.Int(rc.Do("LLEN", key))
	if err != nil {
		return 0, errors.Wrap(err, "redis.LLEN")
	}
	return reply, nil
}

// SetEX 设置带过期时间的键值对
func (rc *RedisClient) SetEX(key string, value any, seconds int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "encode")
	}

	_, err = rc.Do("SETEX", key, seconds, data)
	return errors.Wrap(err, "redis.SETEX")
}

// SetNX 仅当key不存在时设置值
func (rc *RedisClient) SetNX(key string, value any) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, errors.Wrap(err, "encode")
	}

	reply, err := redis.Int(rc.Do("SETNX", key, data))
	if err != nil {
		return false, errors.Wrap(err, "redis.SETNX")
	}
	return reply == 1, nil
}

// MGet 批量获取多个key的值
func (rc *RedisClient) MGet(keys []string, out map[string]any) error {
	args := make([]any, len(keys))
	for i, key := range keys {
		args[i] = key
	}

	values, err := redis.ByteSlices(rc.Do("MGET", args...))
	if err != nil {
		return errors.Wrap(err, "redis.MGET")
	}

	for i, value := range values {
		if value != nil {
			var v any
			if err := json.Unmarshal(value, &v); err != nil {
				return errors.Wrap(err, "decode")
			}
			out[keys[i]] = v
		}
	}
	return nil
}

// MSet 批量设置多个key-value
func (rc *RedisClient) MSet(values map[string]any) error {
	args := make([]any, 0, len(values)*2)
	for k, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			return errors.Wrap(err, "encode")
		}
		args = append(args, k, data)
	}

	_, err := rc.Do("MSET", args...)
	return errors.Wrap(err, "redis.MSET")
}

// Rename 重命名key
func (rc *RedisClient) Rename(oldKey, newKey string) error {
	_, err := rc.Do("RENAME", oldKey, newKey)
	return errors.Wrap(err, "redis.RENAME")
}

// HSet 设置hash字段的值
func (rc *RedisClient) HSet(key, field string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "encode")
	}

	_, err = rc.Do("HSET", key, field, data)
	return errors.Wrap(err, "redis.HSET")
}

// HGet 获取hash字段的值
func (rc *RedisClient) HGet(key, field string, out any) error {
	data, err := redis.Bytes(rc.Do("HGET", key, field))
	if err == redis.ErrNil {
		return errors.New("field not found")
	}
	if err != nil {
		return errors.Wrap(err, "redis.HGET")
	}

	if err := json.Unmarshal(data, out); err != nil {
		return errors.Wrap(err, "decode")
	}
	return nil
}

// HMSet 批量设置hash字段
func (rc *RedisClient) HMSet(key string, fields map[string]any) error {
	args := make([]any, 0, 1+len(fields)*2)
	args = append(args, key)

	for field, value := range fields {
		data, err := json.Marshal(value)
		if err != nil {
			return errors.Wrap(err, "encode")
		}
		args = append(args, field, data)
	}

	_, err := rc.Do("HMSET", args...)
	return errors.Wrap(err, "redis.HMSET")
}

// HMGet 批量获取hash字段
func (rc *RedisClient) HMGet(key string, fields []string, out map[string]any) error {
	args := make([]any, 0, 1+len(fields))
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}

	values, err := redis.ByteSlices(rc.Do("HMGET", args...))
	if err != nil {
		return errors.Wrap(err, "redis.HMGET")
	}

	for i, value := range values {
		if value != nil {
			var v any
			if err := json.Unmarshal(value, &v); err != nil {
				return errors.Wrap(err, "decode")
			}
			out[fields[i]] = v
		}
	}
	return nil
}

// HDel 删除hash字段
func (rc *RedisClient) HDel(key string, fields ...string) error {
	args := make([]any, 0, 1+len(fields))
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}

	_, err := rc.Do("HDEL", args...)
	return errors.Wrap(err, "redis.HDEL")
}

// HExists 检查hash字段是否存在
func (rc *RedisClient) HExists(key, field string) (bool, error) {
	exists, err := redis.Bool(rc.Do("HEXISTS", key, field))
	if err != nil {
		return false, errors.Wrap(err, "redis.HEXISTS")
	}
	return exists, nil
}

// SAdd 向集合添加成员
func (rc *RedisClient) SAdd(key string, members ...any) error {
	args := make([]any, 0, 1+len(members))
	args = append(args, key)

	for _, member := range members {
		data, err := json.Marshal(member)
		if err != nil {
			return errors.Wrap(err, "encode")
		}
		args = append(args, data)
	}

	_, err := rc.Do("SADD", args...)
	return errors.Wrap(err, "redis.SADD")
}

// SMembers 获取集合所有成员
func (rc *RedisClient) SMembers(key string, out any) error {
	reply, err := redis.ByteSlices(rc.Do("SMEMBERS", key))
	if err != nil {
		return errors.Wrap(err, "redis.SMEMBERS")
	}

	var temp []any
	for _, item := range reply {
		var value any
		if err := json.Unmarshal(item, &value); err != nil {
			return errors.Wrap(err, "decode")
		}
		temp = append(temp, value)
	}

	encoded, err := json.Marshal(temp)
	if err != nil {
		return errors.Wrap(err, "encode temporary slice")
	}

	if err := json.Unmarshal(encoded, out); err != nil {
		return errors.Wrap(err, "decode to output slice")
	}

	return nil
}

// SIsMember 判断元素是否是集合成员
func (rc *RedisClient) SIsMember(key string, member any) (bool, error) {
	data, err := json.Marshal(member)
	if err != nil {
		return false, errors.Wrap(err, "encode")
	}

	return redis.Bool(rc.Do("SISMEMBER", key, data))
}

// SRem 移除集合成员
func (rc *RedisClient) SRem(key string, members ...any) error {
	args := make([]any, 0, 1+len(members))
	args = append(args, key)

	for _, member := range members {
		data, err := json.Marshal(member)
		if err != nil {
			return errors.Wrap(err, "encode")
		}
		args = append(args, data)
	}

	_, err := rc.Do("SREM", args...)
	return errors.Wrap(err, "redis.SREM")
}

// SCard 获取集合成员数量
func (rc *RedisClient) SCard(key string) (int, error) {
	return redis.Int(rc.Do("SCARD", key))
}

func randomValue() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (rc *RedisClient) Close() error {
	return rc.pool.Close()
}
