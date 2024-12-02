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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Str string

func (s Str) Key() string {
	return string(s)
}

type RedisQueueTestSuite struct {
	suite.Suite
	queue *RedisQueue[Str]
}

func (s *RedisQueueTestSuite) SetupSuite() {
	s.queue = NewRedisQueue[Str]("localhost:6379", 0, "test_queue")
}

func (s *RedisQueueTestSuite) TearDownTest() {
	// Clear all data in the queue
	for {
		_, err := s.queue.Dequeue()
		if errors.Is(err, QueueEmptyError) {
			break
		}
	}
}

func (s *RedisQueueTestSuite) TestEnqueueDequeue() {
	testCases := []struct {
		name  string
		value Str
	}{
		{"basic string", "hello"},
		{"unicode string", "Hello World"},
		{"special chars", "!@#$%^&*()"},
		{"中文字符串", "你好世界"},
		{"empty string", ""},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Test enqueue
			err := s.queue.Enqueue(tc.value)
			s.NoError(err)

			// Test dequeue
			val, err := s.queue.Dequeue()
			s.NoError(err)
			s.Equal(tc.value, val)
		})
	}
}

func (s *RedisQueueTestSuite) TestEmptyQueue() {
	// Ensure queue is empty
	s.TearDownTest()

	// Test dequeue on empty queue
	_, err := s.queue.Dequeue()
	s.True(errors.Is(err, QueueEmptyError))

	// Test IsEmpty method
	empty, err := s.queue.IsEmpty()
	s.NoError(err)
	s.True(empty)
}

func (s *RedisQueueTestSuite) TestQueueSize() {
	// Ensure queue is empty
	s.TearDownTest()

	testData := []Str{"one", "two", "three"}

	// Test size changes after enqueue
	for i, data := range testData {
		err := s.queue.Enqueue(data)
		s.NoError(err)

		size, err := s.queue.Size()
		s.NoError(err)
		s.Equal(i+1, size)
	}

	// Test size changes after dequeue
	for i := len(testData) - 1; i >= 0; i-- {
		_, err := s.queue.Dequeue()
		s.NoError(err)

		size, err := s.queue.Size()
		s.NoError(err)
		s.Equal(i, size)
	}
}

func (s *RedisQueueTestSuite) TestDequeueTimeout() {
	// Ensure queue is empty
	s.TearDownTest()

	// Start a goroutine that enqueues after a delay
	go func() {
		time.Sleep(2 * time.Second)
		s.queue.Enqueue("delayed")
	}()

	// Test blocking wait
	start := time.Now()
	val, err := s.queue.Dequeue()
	duration := time.Since(start)

	s.NoError(err)
	s.Equal(Str("delayed"), val)
	s.True(duration >= 2*time.Second)
}

func (s *RedisQueueTestSuite) TestConcurrentOperations() {
	// Ensure queue is empty
	s.TearDownTest()

	const goroutines = 10
	done := make(chan bool)

	// Concurrent enqueue
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			err := s.queue.Enqueue(Str(string(rune('A' + id))))
			s.NoError(err)
			done <- true
		}(i)
	}

	// Wait for all enqueue operations to complete
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// Verify queue size
	size, err := s.queue.Size()
	s.NoError(err)
	s.Equal(goroutines, size)

	// 添加互斥锁保护 results map
	var mu sync.Mutex
	results := make(map[string]bool)

	// Concurrent dequeue
	for i := 0; i < goroutines; i++ {
		go func() {
			val, err := s.queue.Dequeue()
			s.NoError(err)

			// 使用互斥锁保护 map 写入
			mu.Lock()
			results[string(val)] = true
			mu.Unlock()

			done <- true
		}()
	}

	// Wait for all dequeue operations to complete
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// Verify all data was processed correctly
	s.Equal(goroutines, len(results))
}

func TestRedisQueue(t *testing.T) {
	suite.Run(t, new(RedisQueueTestSuite))
}

// Helper function for testing
func assertQueueEmpty(t *testing.T, q *RedisQueue[Str]) {
	empty, err := q.IsEmpty()
	assert.NoError(t, err)
	assert.True(t, empty)
}
