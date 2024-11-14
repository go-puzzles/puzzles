// File:		queue.go
// Created by:	Hoven
// Created on:	2024-07-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import "errors"

var (
	QueueEmptyError = errors.New("queue is empty")
)

type Queue[T any] interface {
	Enqueue(value T) error
	Dequeue() (T, error)
	IsEmpty() (bool, error)
	Size() (int, error)
}
