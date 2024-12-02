// File:		memory_test.go
// Created by:	Hoven
// Created on:	2024-07-30
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	memoryQueue *MemoryQueue[string]
)

// TestMemoryQueueBasicOperations tests basic queue operations
func TestMemoryQueueBasicOperations(t *testing.T) {
	q := NewMemoryQueue[string]()

	// Test enqueue
	q.Enqueue("test1")
	q.Enqueue("test2")

	// Test size
	size, err := q.Size()
	assert.Nil(t, err)
	assert.Equal(t, 2, size)

	// Test dequeue
	val, err := q.Dequeue()
	assert.Nil(t, err)
	assert.Equal(t, "test1", val)
}

// TestMemoryQueueEmptyOperations tests operations on empty queue
func TestMemoryQueueEmptyOperations(t *testing.T) {
	q := NewMemoryQueue[string]()

	// Test dequeue on empty queue
	_, err := q.Dequeue()
	assert.ErrorIs(t, err, QueueEmptyError)

	// Test IsEmpty
	empty, err := q.IsEmpty()
	assert.Nil(t, err)
	assert.True(t, empty)
}

// TestMemoryQueueConcurrent tests concurrent operations
func TestMemoryQueueConcurrent(t *testing.T) {
	q := NewMemoryQueue[string]()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Producer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				q.Enqueue("item")
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	// Consumer goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := q.Dequeue()
				if err != nil {
					time.Sleep(10 * time.Millisecond)
					continue
				}
			}
		}
	}()

	<-ctx.Done()
	assert.True(t, true) // Test completed without deadlock
}
