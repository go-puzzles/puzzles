// File:		memory.go
// Created by:	Hoven
// Created on:	2024-07-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

var _ Queue[any] = (*MemoryQueue[any])(nil)

type MemoryQueue[T any] struct {
	items []T
}

func NewMemoryQueue[T any]() *MemoryQueue[T] {
	return &MemoryQueue[T]{items: make([]T, 0)}
}

func (q *MemoryQueue[T]) Enqueue(value T) error {
	q.items = append(q.items, value)
	return nil
}

func (q *MemoryQueue[T]) Dequeue() (T, error) {
	if len(q.items) == 0 {
		var zero T
		return zero, QueueEmptyError
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *MemoryQueue[T]) IsEmpty() (bool, error) {
	return len(q.items) == 0, nil
}

func (q *MemoryQueue[T]) Size() (int, error) {
	return len(q.items), nil
}
