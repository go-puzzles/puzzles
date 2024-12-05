// File:		client.go
// Created by:	Hoven
// Created on:	2024-12-02
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package goredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	redisDialer "github.com/go-puzzles/puzzles/dialer/redis"
	"github.com/redis/go-redis/v9"
)

// Define error types
var (
	// ErrLockAcquireFailed indicates failure to acquire the lock
	ErrLockAcquireFailed = errors.New("failed to acquire lock")
	// ErrLockNotFound indicates the lock does not exist
	ErrLockNotFound = errors.New("lock not found")
	// ErrLockReleaseFailed indicates failure to release the lock
	ErrLockReleaseFailed = errors.New("failed to release lock")
	// ErrLockTimeout indicates lock acquisition timeout
	ErrLockTimeout = errors.New("lock timeout")
)

// Lua script for atomic unlock operation
const unlockScript = `
	if redis.call('get', KEYS[1]) == ARGV[1] then
		return redis.call('del', KEYS[1])
	end
	return 0`

type PuzzleRedisClient struct {
	*redis.Client
	locks sync.Map // stores lock values for validation
}

func NewRedisClient(addr string, db int) *PuzzleRedisClient {
	c := redisDialer.DialGoRedisClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	return &PuzzleRedisClient{Client: c}
}

func NewRedisClientWithAuth(addr string, db int, user, pwd string) *PuzzleRedisClient {
	c := redisDialer.DialGoRedisClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Username: user,
		Password: pwd,
	})
	return &PuzzleRedisClient{Client: c}
}

// TryLock attempts to acquire a distributed lock
func (c *PuzzleRedisClient) TryLock(ctx context.Context, key string, expiration time.Duration) error {
	value := fmt.Sprintf("%s:%d", c.getInstanceID(), time.Now().UnixNano())

	success, err := c.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("%w: %s", ErrLockAcquireFailed, key)
	}

	c.locks.Store(key, value)
	return nil
}

// TryLockWithTimeout attempts to acquire a lock with a timeout
func (c *PuzzleRedisClient) TryLockWithTimeout(ctx context.Context, key string, expiration, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		err := c.TryLock(ctx, key, expiration)
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("%w: %s", ErrLockTimeout, key)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Unlock releases a distributed lock
func (c *PuzzleRedisClient) Unlock(ctx context.Context, key string) error {
	valueI, exists := c.locks.Load(key)
	if !exists {
		return fmt.Errorf("%w: %s", ErrLockNotFound, key)
	}
	value := valueI.(string)

	result, err := c.Eval(ctx, unlockScript, []string{key}, value).Result()
	if err != nil {
		return err
	}

	c.locks.Delete(key)

	if v, ok := result.(int64); !ok || v != 1 {
		return fmt.Errorf("%w: %s", ErrLockReleaseFailed, key)
	}

	return nil
}

// getInstanceID returns the unique identifier for current instance
func (c *PuzzleRedisClient) getInstanceID() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s:%d", hostname, os.Getpid())
}

// SetValue stores any type of value in Redis with automatic type handling
func (c *PuzzleRedisClient) SetValue(ctx context.Context, key string, value any, expiration time.Duration) error {
	redisValue, err := c.convertValueToRedisArg(value)
	if err != nil {
		return fmt.Errorf("failed to convert value to redis arg: %w", err)
	}

	return c.Client.Set(ctx, key, redisValue, expiration).Err()
}

// GetValue retrieves a value from Redis with automatic type conversion
// The result parameter must be a pointer
// Supported types:
// - string
// - []byte
// - int
// - int64
// - float32
// - float64
// - bool
// - time.Time
// For other types, JSON deserialization will be performed
func (c *PuzzleRedisClient) GetValue(ctx context.Context, key string, result any) error {
	rt := reflect.TypeOf(result)
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}

	cmd := c.Client.Get(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return c.convertRedisValueToType(cmd, result)
}

func (c *PuzzleRedisClient) LPushValue(ctx context.Context, key string, values ...any) error {
	args := make([]any, len(values))
	for i, value := range values {
		arg, err := c.convertValueToRedisArg(value)
		if err != nil {
			return fmt.Errorf("failed to convert item %d: %w", i, err)
		}
		args[i] = arg
	}
	return c.LPush(ctx, key, args...).Err()
}

func (c *PuzzleRedisClient) LPopValue(ctx context.Context, key string, result any) error {
	rt := reflect.TypeOf(result)
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}

	cmd := c.LPop(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return c.convertRedisValueToType(cmd, result)
}

func (c *PuzzleRedisClient) RPushValue(ctx context.Context, key string, values ...any) error {
	args := make([]any, len(values))
	for i, value := range values {
		arg, err := c.convertValueToRedisArg(value)
		if err != nil {
			return fmt.Errorf("failed to convert item %d: %w", i, err)
		}
		args[i] = arg
	}
	return c.RPush(ctx, key, args...).Err()
}

func (c *PuzzleRedisClient) RPopValue(ctx context.Context, key string, result any) error {
	rt := reflect.TypeOf(result)
	if rt.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}

	cmd := c.RPop(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return c.convertRedisValueToType(cmd, result)
}

// RangeValue retrieves a range of values from the list and converts them to the specified slice type
// start and stop are inclusive indices
// For example: 0, 10 means get the first 11 elements
// -1 represents the last element, -2 represents the second to last element, and so on
func (c *PuzzleRedisClient) RangeValue(ctx context.Context, key string, start, stop int64, resultPtr any) error {
	resultValue := reflect.ValueOf(resultPtr)
	if !resultValue.IsValid() {
		return fmt.Errorf("result must not be nil")
	}

	if resultValue.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}

	resultValue = resultValue.Elem()
	if resultValue.Kind() != reflect.Slice {
		return fmt.Errorf("result must be a slice")
	}

	cmd := c.LRange(ctx, key, start, stop)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return scanRedisSlice(cmd.Val(), resultValue)
}

// RRangeValue retrieves a range of values from the list starting from the right and converts them to the specified slice type
// start and stop are inclusive indices counted from the right
// For example: 0, 10 means get the first 11 elements from the right
// Note: start and stop are counted from the right, where 0 represents the rightmost element
func (c *PuzzleRedisClient) RRangeValue(ctx context.Context, key string, start, stop int64, resultPtr any) error {
	length, err := c.LLen(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get list length: %w", err)
	}

	leftStart := length - stop - 1
	leftStop := length - start - 1

	return c.RangeValue(ctx, key, leftStart, leftStop, resultPtr)
}

// convertRedisValueToType converts a Redis string value to the specified type
func (c *PuzzleRedisClient) convertRedisValueToType(cmd *redis.StringCmd, result any) (err error) {
	return scan([]byte(cmd.Val()), result)
}

// convertValueToRedisArg converts a value to a format suitable for Redis storage
func (c *PuzzleRedisClient) convertValueToRedisArg(value any) (any, error) {
	switch v := value.(type) {
	case string, int, int64, float32, float64, bool, []byte:
		return v, nil
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("json marshal failed: %w", err)
		}
		return string(jsonBytes), nil
	}
}
