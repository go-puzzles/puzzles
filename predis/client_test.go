// client_test.go
// Created by: Hoven
// Created on: 2024-07-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package predis

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	redisDialer "github.com/go-puzzles/puzzles/dialer/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestClientCommandDo(t *testing.T) {
	client := NewRedisClient(redisDialer.DialRedisPool("localhost:6379", 12, 100))

	res, err := redis.String(client.Do("set", "test-key", 10))
	assert.NoError(t, err)

	t.Logf("command resp: %v", res)

	res, err = redis.String(client.Do("get", "test-key"))
	assert.NoError(t, err)

	t.Logf("command resp: %v", res)
	if res != "10" {
		t.Error("resp no equal")
		return
	}
}

func TestClientSet(t *testing.T) {
	client := NewRedisClient(redisDialer.DialRedisPool("localhost:6379", 13, 100))

	err := client.Set("test-key-set", "test-value")
	assert.NoError(t, err)
	var dataStr string
	err = client.Get("test-key-set", &dataStr)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", dataStr)
	t.Log(dataStr)

	// ================================
	err = client.Set("test-key-set-byte", []byte("test-value"))
	assert.NoError(t, err)
	var dataByte []byte
	err = client.Get("test-key-set-byte", &dataByte)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test-value"), dataByte)
	t.Log(dataByte)

	// ================================
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(map[string]string{"test-key": "test-value"})
	assert.NoError(t, err)

	b := buf.Bytes()
	t.Logf("set b: %v", b)
	err = client.Set("test-gob-set-byte", b)
	assert.NoError(t, err)

	dataByte = []byte{}
	err = client.Get("test-gob-set-byte", &dataByte)
	assert.NoError(t, err)

	t.Log(dataByte)

	dec := gob.NewDecoder(bytes.NewBuffer(dataByte))
	resp := make(map[string]string)
	err = dec.Decode(&resp)

	assert.Equal(t, map[string]string{"test-key": "test-value"}, resp)
}

func TestClientLock(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	const lockKey = "test_lock"

	// 清理可能存在的锁
	client.UnLock(lockKey)

	// 测试获取锁
	err := client.Lock(lockKey)
	assert.NoError(t, err, "首次获取锁应该成功")

	// 启动一个 goroutine 尝试获取已被占用的锁
	done := make(chan bool)
	go func() {
		defer close(done)

		// 尝试获取已被占用的锁，应该失败
		err := client.Lock(lockKey)
		assert.Error(t, err, "获取已被占用的锁应该失败")
		assert.Equal(t, ErrLockFailed, err, "应该返回 ErrLockFailed")

		// 使用 LockWithBlock 等待锁释放
		err = client.LockWithBlock(lockKey, 3)
		assert.NoError(t, err, "等待锁释放后应该成功获取锁")

		// 释放锁
		client.UnLock(lockKey)
	}()

	// 主 goroutine 等待一段时间后释放锁
	time.Sleep(time.Second)
	err = client.UnLock(lockKey)
	assert.NoError(t, err, "释放锁应该成功")

	// 等待另一个 goroutine 完成
	<-done
}

func TestLockTimeout(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	const lockKey = "test_lock_timeout"

	// 清理可能存在的锁
	client.UnLock(lockKey)

	// 设置一个短暂的过期时间
	err := client.Lock(lockKey, time.Second)
	assert.NoError(t, err, "获取锁应该成功")

	// 立即尝试获取锁，应该失败
	err = client.Lock(lockKey)
	assert.Error(t, err, "锁被占用时应该获取失败")

	// 等待锁过期
	time.Sleep(time.Second * 2)

	// 现在应该能够获取锁
	err = client.Lock(lockKey)
	assert.NoError(t, err, "锁过期后应该能够获取")

	// 清理
	client.UnLock(lockKey)
}

func TestLockWithBlockTimeout(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	const lockKey = "test_lock_block"

	// 清理可能存在的锁
	client.UnLock(lockKey)

	// 首先获取锁
	err := client.Lock(lockKey)
	assert.NoError(t, err, "首次获取锁应该成功")

	// 启动一个 goroutine 来测试 LockWithBlock
	start := time.Now()
	go func() {
		time.Sleep(time.Second * 2)
		client.UnLock(lockKey)
	}()

	// 尝试获取锁，设置较短的重试次数
	err = client.LockWithBlock(lockKey, 5)
	duration := time.Since(start)
	assert.NoError(t, err, "应该成功获取锁")
	assert.True(t, duration >= time.Second*2, "应该等待直到锁被释放")

	// 清理
	client.UnLock(lockKey)
}

func TestRedisClient_ListOperations(t *testing.T) {
	pool := newTestPool()
	client := NewRedisClient(pool)
	defer client.Close()

	// Test key for all list operations
	const testKey = "test_list"

	// Clean up before and after test
	defer client.Delete(testKey)
	client.Delete(testKey)

	t.Run("LPush and RPush operations", func(t *testing.T) {
		// Test LPush
		length, err := client.LPush(testKey, "value1", "value2")
		assert.NoError(t, err)
		assert.Equal(t, 2, length)

		// Test RPush
		length, err = client.RPush(testKey, "value3", "value4")
		assert.NoError(t, err)
		assert.Equal(t, 4, length)

		// Verify length
		length, err = client.LLen(testKey)
		assert.NoError(t, err)
		assert.Equal(t, 4, length)
	})

	t.Run("LPop and RPop operations", func(t *testing.T) {
		// Test LPop
		var value string
		err := client.LPop(testKey, &value)
		assert.NoError(t, err)
		assert.Equal(t, "value2", value)

		// Test RPop
		err = client.RPop(testKey, &value)
		assert.NoError(t, err)
		assert.Equal(t, "value4", value)

		// Verify length after pops
		length, err := client.LLen(testKey)
		assert.NoError(t, err)
		assert.Equal(t, 2, length)
	})

	t.Run("LRange operations", func(t *testing.T) {
		// Clear existing list and add new test data
		client.Delete(testKey)
		_, err := client.RPush(testKey, "item1", "item2", "item3", "item4", "item5")
		assert.NoError(t, err)

		// Test getting all elements
		var allItems []string
		err = client.LRange(testKey, 0, -1, &allItems)
		assert.NoError(t, err)
		assert.Equal(t, []string{"item1", "item2", "item3", "item4", "item5"}, allItems)

		// Test getting a subset of elements
		var subsetItems []string
		err = client.LRange(testKey, 1, 3, &subsetItems)
		assert.NoError(t, err)
		assert.Equal(t, []string{"item2", "item3", "item4"}, subsetItems)

		// Test getting elements with negative indices
		var lastItems []string
		err = client.LRange(testKey, -3, -1, &lastItems)
		assert.NoError(t, err)
		assert.Equal(t, []string{"item3", "item4", "item5"}, lastItems)
	})

	t.Run("Empty list operations", func(t *testing.T) {
		emptyKey := "empty_list"
		defer client.Delete(emptyKey)

		// Test LPop on empty list
		var value string
		err := client.LPop(emptyKey, &value)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "list is empty")

		// Test RPop on empty list
		err = client.RPop(emptyKey, &value)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "list is empty")

		// Test LRange on empty list
		var items []string
		err = client.LRange(emptyKey, 0, -1, &items)
		assert.NoError(t, err)
		assert.Empty(t, items)

		// Test LLen on empty list
		length, err := client.LLen(emptyKey)
		assert.NoError(t, err)
		assert.Equal(t, 0, length)
	})

	t.Run("Complex data types", func(t *testing.T) {
		complexKey := "complex_list"
		defer client.Delete(complexKey)

		// Test struct
		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		// Push complex data
		testData := TestStruct{ID: 1, Name: "test"}
		length, err := client.LPush(complexKey, testData)
		assert.NoError(t, err)
		assert.Equal(t, 1, length)

		// Pop and verify complex data
		var result TestStruct
		err = client.LPop(complexKey, &result)
		assert.NoError(t, err)
		assert.Equal(t, testData, result)
	})
}

// Helper function to create a test Redis pool
func newTestPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func setupTestRedis(t *testing.T) *RedisClient {
	// 使用测试环境的Redis配置
	client := NewRedisClientWithAddr("localhost:6379", 0, 10)
	// 清空当前数据库
	_, err := client.Do("FLUSHDB")
	assert.NoError(t, err)
	return client
}

func TestSetEXAndGet(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	key := "test_setex"
	value := map[string]interface{}{
		"name": "test",
		"age":  25,
	}

	// 测试SetEX
	err := client.SetEX(key, value, 1)
	assert.NoError(t, err)

	// 验证值
	var result map[string]interface{}
	err = client.Get(key, &result)
	assert.NoError(t, err)

	// 分别验证每个字段
	assert.Equal(t, value["name"], result["name"])
	assert.InDelta(t, value["age"].(int), result["age"].(float64), 0.0001)

	// 等待过期
	time.Sleep(time.Second * 2)
	err = client.Get(key, &result)
	assert.Error(t, err)
}

func TestSetNX(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	key := "test_setnx"
	value := "test_value"

	// 第一次设置应该成功
	success, err := client.SetNX(key, value)
	assert.NoError(t, err)
	assert.True(t, success)

	// 第二次设置应该失败
	success, err = client.SetNX(key, "new_value")
	assert.NoError(t, err)
	assert.False(t, success)
}

func TestMSetAndMGet(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	values := map[string]interface{}{
		"key1": "value1",
		"key2": float64(123),
		"key3": map[string]interface{}{"nested": "value"},
	}

	// 测试MSet
	err := client.MSet(values)
	assert.NoError(t, err)

	// 测试MGet
	result := make(map[string]interface{})
	keys := []string{"key1", "key2", "key3", "nonexistent"}
	err = client.MGet(keys, result)
	assert.NoError(t, err)

	// 验证结果
	assert.Equal(t, values["key1"], result["key1"])
	assert.Equal(t, values["key2"], result["key2"])
	assert.Equal(t, values["key3"], result["key3"])
	assert.Nil(t, result["nonexistent"])
}

func TestHashOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	hashKey := "test_hash"

	// 测试HSet和HGet
	err := client.HSet(hashKey, "field1", "value1")
	assert.NoError(t, err)

	var result string
	err = client.HGet(hashKey, "field1", &result)
	assert.NoError(t, err)
	assert.Equal(t, "value1", result)

	// 测试HMSet和HMGet
	fields := map[string]interface{}{
		"field2": float64(123),
		"field3": map[string]interface{}{"nested": "value"},
	}
	err = client.HMSet(hashKey, fields)
	assert.NoError(t, err)

	getFields := []string{"field2", "field3", "nonexistent"}
	resultMap := make(map[string]interface{})
	err = client.HMGet(hashKey, getFields, resultMap)
	assert.NoError(t, err)
	assert.Equal(t, fields["field2"], resultMap["field2"])
	assert.Equal(t, fields["field3"], resultMap["field3"])
	assert.Nil(t, resultMap["nonexistent"])

	// 测试HExists
	exists, err := client.HExists(hashKey, "field1")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = client.HExists(hashKey, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)

	// 测试HDel
	err = client.HDel(hashKey, "field1", "field2")
	assert.NoError(t, err)

	exists, err = client.HExists(hashKey, "field1")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestSetOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	setKey := "test_set"

	// 测试SAdd
	err := client.SAdd(setKey, "member1", 123, map[string]string{"nested": "value"})
	assert.NoError(t, err)

	// 测试SCard
	count, err := client.SCard(setKey)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	// 测试SIsMember
	exists, err := client.SIsMember(setKey, "member1")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = client.SIsMember(setKey, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)

	// 测试SMembers
	var members []interface{}
	err = client.SMembers(setKey, &members)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(members))

	// 测试SRem
	err = client.SRem(setKey, "member1")
	assert.NoError(t, err)

	count, err = client.SCard(setKey)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestRename(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	// 设置初始键值
	oldKey := "old_key"
	newKey := "new_key"
	value := "test_value"

	err := client.Set(oldKey, value)
	assert.NoError(t, err)

	// 测试Rename
	err = client.Rename(oldKey, newKey)
	assert.NoError(t, err)

	// 验证旧键不存在
	var result string
	err = client.Get(oldKey, &result)
	assert.Error(t, err)

	// 验证新键存在且值正确
	err = client.Get(newKey, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}
