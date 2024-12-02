// File:		priority_test.go
// Created by:	Hoven
// Created on:	2024-11-15
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockItem struct {
	priority int
	value    string
}

func (m *mockItem) Priority() int {
	return m.priority
}

func (m *mockItem) Value() any {
	return m.value
}

// TestHighPriorityQueueOperations tests high priority queue operations
func TestHighPriorityQueueOperations(t *testing.T) {
	pq := NewPriorityQueue[*mockItem](WithPriorityMode(HighPriorityFirst))

	// Test enqueue
	pq.Enqueue(&mockItem{priority: 1, value: "low"})
	pq.Enqueue(&mockItem{priority: 10, value: "high"})
	pq.Enqueue(&mockItem{priority: 3, value: "medium"})

	// Test size
	size, err := pq.Size()
	assert.Nil(t, err)
	assert.Equal(t, 3, size)

	// Test dequeue order
	expected := []string{"high", "medium", "low"}
	for _, exp := range expected {
		item, err := pq.Dequeue()
		assert.Nil(t, err)
		assert.Equal(t, exp, item.Value())
	}

	// Test empty queue
	_, err = pq.Dequeue()
	assert.ErrorIs(t, err, ErrEmpty)
}

// TestLowPriorityQueueOperations tests low priority queue operations
func TestLowPriorityQueueOperations(t *testing.T) {
	pq := NewPriorityQueue[*mockItem](WithPriorityMode(LowPriorityFirst))

	// Test enqueue with same priorities
	pq.Enqueue(&mockItem{priority: 1, value: "first"})
	pq.Enqueue(&mockItem{priority: 1, value: "second"})

	// Test FIFO behavior for same priority
	item, _ := pq.Dequeue()
	assert.Equal(t, "first", item.Value())
	item, _ = pq.Dequeue()
	assert.Equal(t, "second", item.Value())
}

// TestPriorityQueueEmpty tests empty queue operations
func TestPriorityQueueEmpty(t *testing.T) {
	pq := NewPriorityQueue[*mockItem](WithPriorityMode(HighPriorityFirst))

	// Test IsEmpty on new queue
	empty, err := pq.IsEmpty()
	assert.Nil(t, err)
	assert.True(t, empty)

	// Add and remove an item
	pq.Enqueue(&mockItem{priority: 1, value: "test"})
	item, err := pq.Dequeue()
	assert.Nil(t, err)
	assert.Equal(t, "test", item.Value())

	// Verify queue is empty again
	empty, err = pq.IsEmpty()
	assert.Nil(t, err)
	assert.True(t, empty)
}

// BenchmarkPriorityQueue benchmarks priority queue operations
func BenchmarkPriorityQueue(b *testing.B) {
	pq := NewPriorityQueue[*mockItem](WithPriorityMode(HighPriorityFirst))
	b.Run("Enqueue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pq.Enqueue(&mockItem{priority: i, value: "item"})
		}
	})

	b.Run("Dequeue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pq.Dequeue()
		}
	})
}
