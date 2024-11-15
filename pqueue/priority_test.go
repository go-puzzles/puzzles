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
	pq := NewPriorityQueue(WithPriorityMode(HighPriorityFirst))

	pq.Push(&mockItem{priority: 1, value: "low"})
	pq.Push(&mockItem{priority: 10, value: "high"})
	pq.Push(&mockItem{priority: 3, value: "medium"})

	item, err := pq.Pop()
	if err != nil || item.Value() != "high" {
		t.Errorf("expected high priority item, got %v, error: %v", item, err)
	}
}

func TestLowPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue(WithPriorityMode(LowPriorityFirst))

	pq.Push(&mockItem{priority: 5, value: "high"})
	pq.Push(&mockItem{priority: 1, value: "low"})
	pq.Push(&mockItem{priority: 3, value: "medium"})

	item, err := pq.Pop()
	if err != nil || item.Value() != "low" {
		t.Errorf("expected low priority item, got %v, error: %v", item, err)
	}
}

func BenchmarkPriorityQueue(b *testing.B) {
	pq := NewPriorityQueue(WithPriorityMode(HighPriorityFirst))

	for i := 0; i < b.N; i++ {
		pq.Push(&mockItem{priority: i, value: "item"})
	}

	for i := 0; i < b.N; i++ {
		pq.Pop()
	}
}
