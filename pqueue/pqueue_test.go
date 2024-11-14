// File:		pqueue_test.go
// Created by:	Hoven
// Created on:	2024-07-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-puzzles/puzzles/predis"
)

func TestMain(m *testing.M) {
	memoryQueue = NewMemoryQueue[string]()

	go func() {
		for i := range 100 {
			memoryQueue.Enqueue(fmt.Sprintf("this is index %v", i))
			time.Sleep(time.Second)
		}
	}()

	redisConf := &predis.RedisConf{
		Server: "localhost:6379",
		Db:     12,
	}

	redisQueue = NewRedisQueue[Str](redisConf.DialRedisPool(), "redis-test-queue")
	go func() {
		for i := range 100 {
			if err := redisQueue.Enqueue(Str(fmt.Sprintf("this is redis index %v", i))); err != nil {
				fmt.Println("Redis Enqueue error", err)
				continue
			}
			time.Sleep(time.Second)
		}
	}()

	m.Run()
}
