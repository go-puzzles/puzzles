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
	
	redis2 "github.com/go-puzzles/puzzles/dialer/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestClientCommandDo(t *testing.T) {
	client := NewRedisClient(redis2.DialRedisPool("localhost:6379", 12, 100))
	
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
	client := NewRedisClient(redis2.DialRedisPool("localhost:6379", 13, 100))
	
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
	client := NewRedisClient(redis2.DialRedisPool("localhost:6379", 12, 100))
	
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
