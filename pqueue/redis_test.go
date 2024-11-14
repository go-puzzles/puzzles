// File:		redis_test.go
// Created by:	Hoven
// Created on:	2024-07-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Str string

func (s Str) Key() string {
	return string(s)
}

var (
	redisQueue *RedisQueue[Str]
)

func TestRedisQueueEnqueue(t *testing.T) {
	for {
		val, err := redisQueue.Dequeue()
		if err != nil {
			if errors.Is(err, QueueEmptyError) {
				time.Sleep(4)
				continue
			}
			assert.Nil(t, err)
		}

		fmt.Println(val)
		time.Sleep(4 * time.Second)
	}
}
