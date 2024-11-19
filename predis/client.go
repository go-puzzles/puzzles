package predis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	redisDialer "github.com/go-puzzles/puzzles/dialer/redis"
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
