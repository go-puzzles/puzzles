# PQueue - Go 通用队列库

PQueue 是一个轻量级、灵活的 Go 队列库，提供多种队列实现方式，包括内存队列、优先级队列和 Redis 队列。

## 功能特性

- 🚀 **内存队列**: 基于切片的简单 FIFO 队列实现
- ⭐ **优先级队列**: 支持高/低优先级模式的堆实现队列
- 📦 **Redis 队列**: 基于 Redis 的分布式队列，支持持久化
- 🔒 **线程安全**: 优先级队列内置并发控制
- 🎯 **泛型支持**: 使用 Go 1.18+ 泛型，类型安全
- 🔌 **统一接口**: 所有队列实现都遵循相同的 `Queue` 接口

## 安装

```bash
go get github.com/go-puzzles/puzzles/pqueue
```

## 快速开始

### 基本队列接口

所有队列实现都遵循统一的 `Queue` 接口：

```go
type Queue[T any] interface {
    Enqueue(value T) error
    Dequeue() (T, error)
    IsEmpty() (bool, error)
    Size() (int, error)
}
```

### 内存队列

```go
package main

import (
    "fmt"
    "github.com/go-puzzles/puzzles/pqueue"
)

func main() {
    q := pqueue.NewMemoryQueue[string]()
    
    // 入队
    q.Enqueue("第一个元素")
    q.Enqueue("第二个元素")
    
    // 获取队列大小
    size, _ := q.Size()
    fmt.Printf("队列大小: %d\n", size)
    
    // 出队
    item, err := q.Dequeue()
    if err != nil {
        fmt.Printf("出队错误: %v\n", err)
    } else {
        fmt.Printf("出队元素: %s\n", item)
    }
    
    // 检查是否为空
    empty, _ := q.IsEmpty()
    fmt.Printf("队列是否为空: %t\n", empty)
}
```

### 优先级队列

```go
package main

import (
    "fmt"
    "github.com/go-puzzles/puzzles/pqueue"
)

// 实现 PriorityItem 接口
type Task struct {
    name     string
    priority int
}

func (t *Task) Priority() int {
    return t.priority
}

func main() {
    // 创建高优先级优先的队列
    pq := pqueue.NewPriorityQueue[*Task](
        pqueue.WithPriorityMode(pqueue.HighPriorityFirst),
    )
    
    // 添加任务
    pq.Enqueue(&Task{name: "低优先级任务", priority: 1})
    pq.Enqueue(&Task{name: "高优先级任务", priority: 10})
    pq.Enqueue(&Task{name: "中优先级任务", priority: 5})
    
    // 按优先级顺序出队
    for {
        empty, _ := pq.IsEmpty()
        if empty {
            break
        }
        
        task, err := pq.Dequeue()
        if err != nil {
            break
        }
        fmt.Printf("执行任务: %s (优先级: %d)\n", task.name, task.priority)
    }
}
```

### Redis 队列

```go
package main

import (
    "fmt"
    "github.com/go-puzzles/puzzles/pqueue"
)

// 实现 Item 接口
type Job struct {
    ID   string `json:"id"`
    Data string `json:"data"`
}

func (j *Job) Key() string {
    return j.ID
}

func main() {
    // 创建 Redis 队列
    rq := pqueue.NewRedisQueue[*Job]("localhost:6379", 0, "job_queue")
    
    // 入队任务
    job := &Job{ID: "job_001", Data: "处理数据"}
    err := rq.Enqueue(job)
    if err != nil {
        fmt.Printf("入队失败: %v\n", err)
        return
    }
    
    // 出队任务（阻塞等待）
    processedJob, err := rq.Dequeue()
    if err != nil {
        fmt.Printf("出队失败: %v\n", err)
        return
    }
    
    fmt.Printf("处理任务: %s - %s\n", processedJob.ID, processedJob.Data)
}
```

## API 文档

### 内存队列 (MemoryQueue)

- `NewMemoryQueue[T any]() *MemoryQueue[T]`: 创建新的内存队列

### 优先级队列 (PriorityQueue)

- `NewPriorityQueue[T PriorityItem](opts ...PriorityQueueOption) *PriorityQueue[T]`: 创建优先级队列
- `WithPriorityMode(mode PriorityMode)`: 设置优先级模式
  - `HighPriorityFirst`: 高优先级优先
  - `LowPriorityFirst`: 低优先级优先

#### PriorityItem 接口

需要实现的接口：

```go
type PriorityItem interface {
    Priority() int
}
```

### Redis 队列 (RedisQueue)

- `NewRedisQueue[T Item](addr string, db int, queue string) *RedisQueue[T]`: 创建 Redis 队列
- `NewRedisQueueWithClient[T Item](client *goredis.PuzzleRedisClient, queue string) *RedisQueue[T]`: 使用现有客户端创建队列

#### Item 接口

需要实现的接口：

```go
type Item interface {
    Key() string
}
```

## 错误处理

- `QueueEmptyError`: 队列为空时的错误
- `ErrEmpty`: 优先级队列为空时的错误

## 并发安全

- **内存队列**: 非线程安全，需要外部同步
- **优先级队列**: 内置读写锁，线程安全
- **Redis 队列**: 基于 Redis 原子操作，天然支持并发

## 性能考虑

- **内存队列**: 最高性能，但仅限单机使用
- **优先级队列**: 堆操作复杂度 O(log n)，适合中等规模数据
- **Redis 队列**: 网络开销，但支持分布式和持久化

## 使用场景

### 内存队列

- 单机应用的简单任务队列
- 高性能要求的场景
- 临时数据缓存

### 优先级队列

- 任务调度系统
- 事件处理优先级排序
- 算法中的优先级数据结构

### Redis 队列

- 分布式系统任务队列
- 需要持久化的任务队列
- 跨服务通信

## 许可证

Copyright (c) 2024 Example Corp. All rights reserved.
