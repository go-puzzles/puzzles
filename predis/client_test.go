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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-puzzles/puzzles/dialer"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestClientCommandDo(t *testing.T) {
	client := NewRedisClient(dialer.DialRedisPool("localhost:6379", 12, 100))

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
	client := NewRedisClient(dialer.DialRedisPool("localhost:6379", 13, 100))

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
	client := NewRedisClient(dialer.DialRedisPool("localhost:6379", 12, 100))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(time.Second * 2)
		if err := client.Lock("testLock"); err != nil {
			t.Errorf("lock error: %v", err)
			return
		}

		fmt.Println("after lockA")
		time.Sleep(time.Second * 4)
		client.UnLock("testLock")
	}()

	// time.Sleep(time.Second)

	if err := client.LockWithBlock("testLock", 10); err != nil {
		t.Errorf("lock error: %v", err)
		return
	}

	fmt.Println("after lockB")

	wg.Wait()
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
