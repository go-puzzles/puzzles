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

type PriorityItem interface {
	Priority() int
	Value() any
}

type priorityItem struct {
	index int
	Item  PriorityItem
}

type queue struct {
	q            []*priorityItem
	priorityMode PriorityMode
}

func (q queue) Len() int { return len(q.q) }

func (q queue) Less(i, j int) bool {
	if q.priorityMode == HighPriorityFirst {
		return q.q[i].Item.Priority() > q.q[j].Item.Priority()
	}
	return q.q[i].Item.Priority() < q.q[j].Item.Priority()
}

func (q queue) Swap(i, j int) {
	q.q[i], q.q[j] = q.q[j], q.q[i]
	q.q[i].index = i
	q.q[j].index = j
}

func (q *queue) Push(a any) {
	n := len(q.q)
	item := a.(*priorityItem)
	item.index = n
	q.q = append(q.q, item)
}

func (q *queue) Pop() any {
	old := q.q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	q.q = old[0 : n-1]
	return item
}

type PriorityQueue struct {
	lock         sync.RWMutex
	data         queue
	priorityMode PriorityMode
}

type PriorityQueueOption func(*PriorityQueue)

func WithPriorityMode(mode PriorityMode) PriorityQueueOption {
	return func(pq *PriorityQueue) {
		pq.priorityMode = mode
	}
}

func NewPriorityQueue(opts ...PriorityQueueOption) *PriorityQueue {
	pq := &PriorityQueue{
		data: queue{
			q: make([]*priorityItem, 0),
		},
	}
	for _, opt := range opts {
		opt(pq)
	}

	pq.data.priorityMode = pq.priorityMode
	return pq
}

func (pq *PriorityQueue) Len() int {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return pq.data.Len()
}

func (pq *PriorityQueue) Pop() (PriorityItem, error) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if pq.data.Len() == 0 {
		return nil, ErrEmpty
	}

	item := heap.Pop(&pq.data).(*priorityItem)
	return item.Item, nil
}

func (pq *PriorityQueue) Push(i PriorityItem) error {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	pi := &priorityItem{
		Item: i,
	}
	heap.Push(&pq.data, pi)
	return nil
}
