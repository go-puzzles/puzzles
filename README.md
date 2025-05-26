# GO-PUZZLES

本项目是一个 Go  微服务项目工具库及脚手架，包含了在开发过程中常用的各种功能模块。

## 模块概览

以下是 `go-puzzles` 包含的主要模块及其功能简介：

* [`cores`](cores/README.md):  脚手架核心功能。
* [`dialer`](dialer/README.md): 网络连接相关工具 (例如：TCP/UDP dialer)。
* [`example`](#example): 示例代码或用法演示。
* [`pgin`](pgin/README.md): 与 [Gin](https://gin-gonic.com/zh-cn/docs/) Web 框架相关的库或工具。
* [`plog`](plog/README.md): 日志记录相关的工具。
* [`pgorm`](pgorm/README.md): MySQL 数据库 ORM 操作相关的工具库。
* [`goredis`](goredis/README.md): Redis 操作相关的工具库。
* [`predis`](predis/README.md): Redis 客户端工具库，提供高级操作功能。
* [`snail`](snail/README.md):  延迟函数。
* [`pflags`](pflags/README.md): 命令行参数或配置项处理相关的工具。
* [`putils`](putils/README.md): 通用工具函数集合。
* [`dice`](dice/README.md): 基于权重的随机抽取工具，支持不放回抽样模式。
* [`poss`](poss/README.md): 对象存储服务 (OSS) 操作接口，支持文件上传、下载等功能。
* [`perror`](perror/README.md): 错误处理工具库，提供错误码、错误原因追踪等功能。
* [`penum`](penum/README.md): 枚举类型生成工具，支持自动生成字符串和数字类型的枚举。
* [`pqueue`](pqueue/README.md): 队列数据结构实现，包括内存队列、优先队列和Redis队列。
* [`cache`](cache/README.md): 缓存接口和实现，支持内存缓存和TTL功能。

## 使用示例

### `cores` - 基础服务框架

`cores` 提供了一个基础的服务启动和管理框架，集成了多种常用功能：

#### 主要特性

* **优雅退出**: 捕获系统信号并安全关闭所有资源
* **HTTP服务**: 内置HTTP服务器支持，可挂载自定义处理器
* **后台Worker**: 支持添加多个后台工作协程
* **Pprof支持**: 内置性能分析工具
* **链路追踪**: 支持Jaeger/OpenTelemetry链路追踪
* **Sentry监控**: 集成Sentry错误捕获和监控
* **Kafka日志**: 支持将日志输出到Kafka

#### 使用方法

```go
package main

import (
 "context"
 "fmt"
 "time"

 "github.com/go-puzzles/puzzles/cores"
 "github.com/go-puzzles/puzzles/plog"
 "github.com/go-puzzles/puzzles/pflags"
)

var (
 port = pflags.Int("port", 8080, "Service port")
)

func main() {
 pflags.Parse()

 // 创建一个新的 Cores 服务实例
 srv := cores.NewCores(
  // 设置服务别名
  cores.WithServiceAlias("my-service"),
  // 启用 pprof 性能分析接口
  cores.WithPprof(),
  // 启用 Sentry 监控
  cores.WithSentryMonitor(),
  // 启用 Kafka 日志
  cores.WithKafkaLog(),
  // 启用链路追踪
  cores.WithTracing(),
  // 配置HTTP处理器
  cores.WithHttpHandler("/api", myHttpHandler),
  // 启用跨域支持
  cores.WithHttpCORS(),
  // 添加一个后台 worker 任务
  cores.WithWorker(func(ctx context.Context) error {
   ticker := time.NewTicker(time.Second * 5)
   defer ticker.Stop()

   plog.Infof("Background worker started.")
   for {
    select {
    case <-ctx.Done(): // 监听退出信号
     plog.Infof("Background worker stopping: %v", ctx.Err())
     return ctx.Err()
    case t := <-ticker.C:
     fmt.Printf("Worker is running at %s\n", t.Format(time.RFC3339))
     // 在这里执行你的后台任务逻辑
    }
   }
  }),
  // 设置默认等待时间
  cores.WithDefaultMaxWait(time.Second * 10),
  // 等待所有worker完成后再退出
  cores.WithWaitAllDone(),
 )

 // 启动服务，监听指定端口
 // Start 会阻塞直到服务退出
 plog.Infof("Starting server on port %d", port())
 err := cores.Start(srv, port())
 if err != nil {
  plog.Panicf("Failed to start server: %v", err)
 }

 plog.Infof("Server stopped gracefully.")
}
```

### `pflags` - 配置加载

详细文档请参阅 [pflags/README.md](pflags/README.md)。

### `plog` - 日志库

详细文档请参阅 [plog/README.md](plog/README.md)。

### `goredis` - Redis客户端

详细文档请参阅 [goredis/README.md](goredis/README.md)。

### `pgorm` - MySQL数据库ORM

详细文档请参阅 [pgorm/README.md](pgorm/README.md)。

### `dice` - 权重随机抽取

基于权重的随机抽取工具，支持不放回抽样模式，适用于抽奖、抽卡等场景。

```go
import "github.com/go-puzzles/puzzles/dice"

// 创建一个权重骰子
weights := []int{10, 20, 30, 40} // 权重分别为10, 20, 30, 40
dice := dice.NewDice(weights)

// 不放回抽样 - 适用于抽奖场景
for {
    n := dice.Next()
    if n == -1 {
        break // 所有奖品已抽完
    }
    fmt.Printf("抽中第%d号奖品\n", n)
}

// 重置后可以重新开始
dice.Reset()
```

### `perror` - 错误处理

提供错误码、错误原因追踪等功能的错误处理工具。

```go
import "github.com/go-puzzles/puzzles/perror"

// 创建带错误码的错误
err := perror.PackError(perror.CodeInvalidInput, "用户名不能为空")

// 包装现有错误
originalErr := errors.New("数据库连接失败")
wrappedErr := perror.WrapError(perror.CodeUnknown, originalErr, "用户查询失败")

// 获取错误码
code := perror.GetErrorCode(err)
```

### `pqueue` - 队列实现

提供多种队列实现，包括内存队列、优先队列和Redis队列。

```go
import "github.com/go-puzzles/puzzles/pqueue"

// 内存队列
queue := pqueue.NewMemoryQueue[string]()
queue.Enqueue("item1")
item, err := queue.Dequeue()

// 优先队列
priorityQueue := pqueue.NewPriorityQueue[string]()
// 使用示例请参考具体实现
```

### `cache` - 缓存接口

提供缓存接口和实现，支持内存缓存和TTL功能。

```go
import "github.com/go-puzzles/puzzles/cache"

// 内存缓存
cache := cache.NewMemoryCache()

// 设置缓存
cache.Set("key", "value")

// 获取缓存
var value string
err := cache.Get("key", &value)

// 获取或创建缓存
err = cache.GetOrCreate("key", func() (any, error) {
    return "computed value", nil
}, &value)
```

### `poss` - 对象存储服务

提供对象存储服务(OSS)操作接口，支持文件上传、下载等功能，目前支持MinIO。

```go
import "github.com/go-puzzles/puzzles/poss"

// 使用示例请参考具体实现
```

### `penum` - 枚举类型生成

枚举类型生成工具，支持自动生成字符串和数字类型的枚举。

```go
import "github.com/go-puzzles/puzzles/penum"

// 定义枚举结构
type Status struct {
    Pending   string
    Running   int
    Completed string
}

// 生成枚举
status := penum.New[Status]()
// status.Pending = "Pending"
// status.Running = 1
// status.Completed = "Completed"
```

### `snail` - 延迟初始化

延迟初始化工具，用于注册和管理需要延迟执行的初始化函数。

```go
import "github.com/go-puzzles/puzzles/snail"

// 注册延迟初始化对象
snail.RegisterObject("database", func() error {
    // 数据库初始化逻辑
    return nil
})

// 执行所有注册的初始化函数
snail.Init()
```

### `putils` - 通用工具函数

通用工具函数集合，包含字符串处理、时间处理、随机数生成、网络工具等。

```go
import "github.com/go-puzzles/puzzles/putils"

// 使用示例请参考具体实现
```

### `dialer` - 网络连接工具

网络连接相关工具，提供数据库连接(MySQL、SQLite)、gRPC连接等功能。

```go
import "github.com/go-puzzles/puzzles/dialer"

// 使用示例请参考具体实现
```

### `pgin` - Gin框架扩展

与Gin Web框架相关的扩展库，提供中间件、错误处理、响应封装等功能。

```go
import "github.com/go-puzzles/puzzles/pgin"

// 使用示例请参考具体实现
```
