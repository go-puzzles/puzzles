// File:		memory_test.go
// Created by:	Hoven
// Created on:	2024-07-30
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

var (
	memoryQueue *MemoryQueue[string]
)

func TestMemoryQueueEnqueue(t *testing.T) {
	for {
		val, err := memoryQueue.Dequeue()
		if err != nil {
			if errors.Is(err, QueueEmptyError) {
				time.Sleep(2)
				continue
			}
			assert.Nil(t, err)
		}

		fmt.Println(val)
	}
}
