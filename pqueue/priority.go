// File:		priority.go
// Created by:	Hoven
// Created on:	2024-11-15
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import (
	"container/heap"
	"errors"
	"sync"
)

// ErrEmpty is returned for queues with no items
var ErrEmpty = errors.New("queue is empty")

type PriorityMode int

const (
	HighPriorityFirst PriorityMode = iota
	LowPriorityFirst
)

type priorityItem[T any] struct {
	item     T
	index    int
	priority int
}

type queue[T any] struct {
	q            []*priorityItem[T]
	priorityMode PriorityMode
}

func (q queue[T]) Len() int { return len(q.q) }

func (q queue[T]) Less(i, j int) bool {
	if q.priorityMode == HighPriorityFirst {
		return q.q[i].priority > q.q[j].priority
	}
	return q.q[i].priority < q.q[j].priority
}

func (q queue[T]) Swap(i, j int) {
	q.q[i], q.q[j] = q.q[j], q.q[i]
	q.q[i].index = i
	q.q[j].index = j
}

func (q *queue[T]) Push(a any) {
	n := len(q.q)
	item := a.(*priorityItem[T])
	item.index = n
	q.q = append(q.q, item)
}

func (q *queue[T]) Pop() any {
	old := q.q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	q.q = old[0 : n-1]
	return item
}

type PriorityQueue[T any] struct {
	lock sync.RWMutex
	data queue[T]
	opts *priorityOpts
}

type priorityOpts struct {
	priorityMode PriorityMode
}

type PriorityQueueOption func(*priorityOpts)

func WithPriorityMode(mode PriorityMode) PriorityQueueOption {
	return func(pq *priorityOpts) {
		pq.priorityMode = mode
	}
}

func NewPriorityQueue[T any](opts ...PriorityQueueOption) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		opts: &priorityOpts{},
		data: queue[T]{
			q: make([]*priorityItem[T], 0),
		},
	}
	for _, opt := range opts {
		opt(pq.opts)
	}

	pq.data.priorityMode = pq.opts.priorityMode
	return pq
}

func (pq *PriorityQueue[T]) Len() int {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return pq.data.Len()
}

func (pq *PriorityQueue[T]) Pop() (T, error) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if pq.data.Len() == 0 {
		var zero T
		return zero, ErrEmpty
	}

	item := heap.Pop(&pq.data).(*priorityItem[T])
	return item.item, nil
}

func (pq *PriorityQueue[T]) Push(i T, priority int) error {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	pi := &priorityItem[T]{
		item:     i,
		priority: priority,
	}
	heap.Push(&pq.data, pi)
	return nil
}
