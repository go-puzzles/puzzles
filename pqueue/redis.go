// File:		redis.go
// Created by:	Hoven
// Created on:	2024-07-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import (
	"encoding/json"

	"github.com/go-puzzles/puzzles/predis"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

var _ Queue[*noItem] = (*RedisQueue[*noItem])(nil)

type noItem struct{}

func (i *noItem) Key() string {
	return ""
}

type Item interface {
	Key() string
}

type RedisQueue[T Item] struct {
	client *predis.RedisClient
	queue  string
}

func NewRedisQueue[T Item](pool *redis.Pool, queue string) *RedisQueue[T] {
	return &RedisQueue[T]{client: predis.NewRedisClient(pool), queue: queue}
}

func (q *RedisQueue[T]) Do(command string, args ...any) (reply any, err error) {
	return q.client.Do(command, args...)
}

func (q *RedisQueue[T]) Enqueue(item T) error {
	serializedValue, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = q.Do("RPUSH", q.queue, serializedValue)
	return err

}

func (q *RedisQueue[T]) parseQueueData(bulk []any) (T, error) {
	var ret T
	var zero T
	if err := json.Unmarshal(bulk[1].([]byte), &ret); err != nil {
		return zero, errors.Wrap(err, "redisDequeueDecode")
	}

	return ret, nil
}

func (q *RedisQueue[T]) Dequeue() (T, error) {
	var zero T

	bulks, err := redis.Values(q.Do("BLPOP", q.queue, 5))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return zero, QueueEmptyError
		}
		return zero, err
	}

	return q.parseQueueData(bulks)
}

func (q *RedisQueue[T]) IsEmpty() (bool, error) {
	length, err := q.size()
	if err != nil {
		return false, err
	}

	return length == 0, nil
}

func (q *RedisQueue[T]) size() (int, error) {
	length, err := redis.Int(q.Do("LLEN", q.queue))
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (q *RedisQueue[T]) Size() (int, error) {
	return q.size()
}
