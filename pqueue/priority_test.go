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

func TestHighPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[*mockItem](WithPriorityMode(HighPriorityFirst))

	pq.Enqueue(&mockItem{priority: 1, value: "low"})
	pq.Enqueue(&mockItem{priority: 10, value: "high"})
	pq.Enqueue(&mockItem{priority: 3, value: "medium"})

	item, err := pq.Dequeue()
	if err != nil || item.Value() != "high" {
		t.Errorf("expected high priority item, got %v, error: %v", item, err)
	}
}

func TestLowPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[*mockItem](WithPriorityMode(LowPriorityFirst))

	pq.Enqueue(&mockItem{priority: 5, value: "high"})
	pq.Enqueue(&mockItem{priority: 1, value: "low"})
	pq.Enqueue(&mockItem{priority: 3, value: "medium"})

	item, err := pq.Dequeue()
	if err != nil || item.Value() != "low" {
		t.Errorf("expected low priority item, got %v, error: %v", item, err)
	}
}

func BenchmarkPriorityQueue(b *testing.B) {
	pq := NewPriorityQueue[*mockItem](WithPriorityMode(HighPriorityFirst))

	for i := 0; i < b.N; i++ {
		pq.Enqueue(&mockItem{priority: i, value: "item"})
	}

	for i := 0; i < b.N; i++ {
		pq.Dequeue()
	}
}
