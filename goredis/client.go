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
func (c *PuzzleRedisClient) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	switch v := value.(type) {
	case string, int, int64, float32, float64, bool, []byte:
		return c.Client.Set(ctx, key, v, expiration).Err()
	default:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("json marshal failed: %w", err)
		}
		return c.Client.Set(ctx, key, jsonBytes, expiration).Err()
	}
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
func (c *PuzzleRedisClient) GetValue(ctx context.Context, key string, result interface{}) error {
	cmd := c.Client.Get(ctx, key)

	var err error
	switch ptr := result.(type) {
	case *string:
		*ptr, err = cmd.Result()
	case *[]byte:
		*ptr, err = cmd.Bytes()
	case *int:
		*ptr, err = cmd.Int()
	case *int64:
		*ptr, err = cmd.Int64()
	case *float32:
		*ptr, err = cmd.Float32()
	case *float64:
		*ptr, err = cmd.Float64()
	case *bool:
		*ptr, err = cmd.Bool()
	case *time.Time:
		*ptr, err = cmd.Time()
	default:
		var b []byte
		b, err = cmd.Bytes()
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, result)
	}

	return err
}
