// File:		redis.go
// Created by:	Hoven
// Created on:	2024-07-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pqueue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-puzzles/puzzles/goredis"
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
	client *goredis.PuzzleRedisClient
	queue  string
	ctx    context.Context
}

func NewRedisQueue[T Item](addr string, db int, queue string) *RedisQueue[T] {
	return &RedisQueue[T]{
		client: goredis.NewRedisClient(addr, db),
		queue:  queue,
		ctx:    context.Background(),
	}
}

func NewRedisQueueWithClient[T Item](client *goredis.PuzzleRedisClient, queue string) *RedisQueue[T] {
	return &RedisQueue[T]{
		client: client,
		queue:  queue,
		ctx:    context.Background(),
	}
}

func (q *RedisQueue[T]) Enqueue(item T) error {
	serializedValue, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return q.client.RPush(q.ctx, q.queue, serializedValue).Err()
}

func (q *RedisQueue[T]) parseQueueData(data string) (T, error) {
	var ret T
	var zero T
	if err := json.Unmarshal([]byte(data), &ret); err != nil {
		return zero, errors.Wrap(err, "redisDequeueDecode")
	}

	return ret, nil
}

func (q *RedisQueue[T]) Dequeue() (T, error) {
	var zero T

	// 使用BLPop阻塞获取数据，超时时间5秒
	result, err := q.client.BLPop(q.ctx, 5*time.Second, q.queue).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return zero, QueueEmptyError
		}
		return zero, err
	}

	// BLPop返回一个字符串切片，第一个元素是键名，第二个元素是值
	if len(result) != 2 {
		return zero, errors.New("invalid redis response")
	}

	return q.parseQueueData(result[1])
}

func (q *RedisQueue[T]) IsEmpty() (bool, error) {
	length, err := q.size()
	if err != nil {
		return false, err
	}

	return length == 0, nil
}

func (q *RedisQueue[T]) size() (int, error) {
	length, err := q.client.LLen(q.ctx, q.queue).Result()
	if err != nil {
		return 0, err
	}

	return int(length), nil
}

func (q *RedisQueue[T]) Size() (int, error) {
	return q.size()
}
