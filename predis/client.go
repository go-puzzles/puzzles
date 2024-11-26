package predis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

var (
	defaultTTL = time.Second * 10
)

type RedisClient struct {
	pool *redis.Pool
}

func NewRedisClient(pool *redis.Pool) *RedisClient {
	return &RedisClient{
		pool: pool,
	}
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
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "encode")
	}

	if ttl > 0 {
		_, err = rc.Do("SET", key, data, "EX", int(ttl.Seconds()))
	} else {
		_, err = rc.Do("SET", key, data)
	}

	return errors.Wrap(err, "redis.Set")
}

func (rc *RedisClient) Set(key string, value any) error {
	return rc.SetWithTTL(key, value, 0)
}

func (rc *RedisClient) Get(key string, out any) error {
	data, err := redis.Bytes(rc.Do("GET", key))
	if err != nil {
		return errors.Wrap(err, "redis.GET")
	}

	if err := json.Unmarshal(data, &out); err != nil {
		return errors.Wrap(err, "decode")
	}

	return nil
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
func (rc *RedisClient) LockWithBlock(key string, maxRetry int, expires ...time.Duration) (err error) {
	for i := 0; i < maxRetry; i++ {
		err = rc.Lock(key, expires...)
		if err == nil {
			return nil
		}

		if errors.Is(err, ErrLockFailed) {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		return err
	}

	return ErrLockFailed
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
func (rc *RedisClient) Lock(key string, expires ...time.Duration) (err error) {
	expire := defaultTTL
	if len(expires) != 0 {
		expire = expires[0]
	}

	_, err = redis.String(rc.Do("SET", key, time.Now().Unix(), "EX", int(expire.Seconds()), "NX"))
	if err == redis.ErrNil {
		return ErrLockFailed
	}

	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisClient) UnLock(key string) (err error) {
	return rc.Delete(key)
}

func (rc *RedisClient) Close() error {
	return rc.pool.Close()
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
