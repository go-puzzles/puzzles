# goredis

goredis 是一个基于 go-redis 库封装的 Redis 客户端工具库，提供了便捷的 Redis 操作接口和高级功能。

## 功能特性

- 简化的 Redis 客户端配置和创建
- 支持自动类型转换的值存储和读取
- 提供分布式锁功能
- 支持列表操作（支持自动类型转换）
- 集成原生 go-redis 的所有功能

## 安装

```bash
go get github.com/go-puzzles/puzzles/goredis
```

## 基本使用

### 创建客户端

```go
package main

import (
    "context"
    "fmt"
    "github.com/go-puzzles/puzzles/goredis"
)

func main() {
    // 创建 Redis 客户端（无认证）
    client := goredis.NewRedisClient("localhost:6379", 0)
    defer client.Close()
    
    // 创建 Redis 客户端（带认证）
    authClient := goredis.NewRedisClientWithAuth("localhost:6379", 0, "username", "password")
    defer authClient.Close()
    
    ctx := context.Background()
    
    // 基本操作
    err := client.Set(ctx, "key", "value", 0)
    if err != nil {
        panic(err)
    }
    
    val, err := client.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Value: %s\n", val)
}
```

### 使用配置结构体

```go
// 使用配置结构体
conf := &goredis.RedisConf{
    Server:   "localhost:6379",
    Db:       0,
    Username: "",
    Password: "",
}
conf.SetDefault() // 设置默认值

client := conf.DialRedisClient()
defer client.Close()
```

## 高级功能

### 自动类型转换

`goredis` 提供了 `SetValue` 和 `GetValue` 方法，支持多种 Go 类型的自动转换：

```go
import (
    "context"
    "time"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

ctx := context.Background()

// 支持的基本类型
err := client.SetValue(ctx, "string_key", "hello", time.Hour)
err = client.SetValue(ctx, "int_key", 42, time.Hour)
err = client.SetValue(ctx, "bool_key", true, time.Hour)
err = client.SetValue(ctx, "float_key", 3.14, time.Hour)
err = client.SetValue(ctx, "time_key", time.Now(), time.Hour)
err = client.SetValue(ctx, "duration_key", time.Second*30, time.Hour)

// 复杂类型（自动 JSON 序列化）
user := User{ID: 1, Name: "张三", Age: 25}
err = client.SetValue(ctx, "user:1", user, time.Hour)

// 获取值（必须传入指针）
var stringVal string
err = client.GetValue(ctx, "string_key", &stringVal)

var intVal int
err = client.GetValue(ctx, "int_key", &intVal)

var userVal User
err = client.GetValue(ctx, "user:1", &userVal)

fmt.Printf("User: %+v\n", userVal)
```

支持的类型：
- 基本类型：`string`、`[]byte`、`int`、`int64`、`float32`、`float64`、`bool`
- 时间类型：`time.Time`、`time.Duration`
- 复杂类型：通过 JSON 序列化/反序列化

### 分布式锁

```go
import (
    "context"
    "time"
)

ctx := context.Background()

// 尝试获取锁
err := client.TryLock(ctx, "resource:123", time.Minute)
if err != nil {
    if errors.Is(err, goredis.ErrLockAcquireFailed) {
        fmt.Println("锁已被其他进程持有")
        return
    }
    panic(err)
}

// 执行业务逻辑
fmt.Println("执行关键业务逻辑...")

// 释放锁
err = client.Unlock(ctx, "resource:123")
if err != nil {
    panic(err)
}
```

### 带超时的分布式锁

```go
// 在指定时间内尝试获取锁
err := client.TryLockWithTimeout(ctx, "resource:123", time.Minute, 5*time.Second)
if err != nil {
    if errors.Is(err, goredis.ErrLockTimeout) {
        fmt.Println("获取锁超时")
        return
    }
    panic(err)
}

// 执行业务逻辑
fmt.Println("执行关键业务逻辑...")

// 释放锁
err = client.Unlock(ctx, "resource:123")
```

### 列表操作（支持类型转换）

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

ctx := context.Background()

// 左侧推入多种类型
err := client.LPushValue(ctx, "test_list", 
    "字符串", 
    123, 
    true, 
    3.14, 
    time.Now(),
    Person{Name: "张三", Age: 30},
)

// 右侧推入
err = client.RPushValue(ctx, "test_list", "新元素")

// 弹出并自动转换类型
var stringVal string
err = client.LPopValue(ctx, "test_list", &stringVal)

var person Person
err = client.RPopValue(ctx, "test_list", &person)

// 范围获取
var results []interface{}
err = client.RangeValue(ctx, "test_list", 0, -1, &results)

// 从右边开始的范围获取
var rightResults []string
err = client.RRangeValue(ctx, "test_list", 0, 2, &rightResults)
```

### 删除操作

```go
// 删除键
err := client.DeleteValue(ctx, "key")
if err != nil {
    panic(err)
}
```

## 错误处理

库定义了以下错误类型：

```go
var (
    ErrLockAcquireFailed = errors.New("failed to acquire lock")
    ErrLockNotFound      = errors.New("lock not found")  
    ErrLockReleaseFailed = errors.New("failed to release lock")
    ErrLockTimeout       = errors.New("lock timeout")
)
```

错误处理示例：

```go
import "errors"

// 处理锁相关错误
err := client.TryLock(ctx, "resource", time.Minute)
if err != nil {
    switch {
    case errors.Is(err, goredis.ErrLockAcquireFailed):
        fmt.Println("锁已被占用")
    case errors.Is(err, goredis.ErrLockTimeout):
        fmt.Println("获取锁超时")
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}

// 处理键不存在的情况
var value string
err = client.GetValue(ctx, "nonexistent", &value)
if err != nil {
    if errors.Is(err, redis.Nil) {
        fmt.Println("键不存在")
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 完整示例

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/go-puzzles/puzzles/goredis"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // 创建客户端
    client := goredis.NewRedisClient("localhost:6379", 0)
    defer client.Close()
    
    ctx := context.Background()
    
    // 存储不同类型的值
    client.SetValue(ctx, "str", "hello", time.Hour)
    client.SetValue(ctx, "num", 42, time.Hour)
    client.SetValue(ctx, "person", Person{Name: "张三", Age: 30}, time.Hour)
    
    // 获取值
    var str string
    var num int
    var person Person
    
    client.GetValue(ctx, "str", &str)
    client.GetValue(ctx, "num", &num)
    client.GetValue(ctx, "person", &person)
    
    fmt.Printf("字符串: %s\n", str)
    fmt.Printf("数字: %d\n", num)
    fmt.Printf("对象: %+v\n", person)
    
    // 使用分布式锁
    err := client.TryLock(ctx, "critical_section", time.Minute)
    if err == nil {
        fmt.Println("获取锁成功，执行关键代码...")
        time.Sleep(time.Second)
        client.Unlock(ctx, "critical_section")
        fmt.Println("释放锁")
    }
    
    // 列表操作
    client.LPushValue(ctx, "list", "a", "b", "c")
    
    var item string
    client.RPopValue(ctx, "list", &item)
    fmt.Printf("弹出的元素: %s\n", item)
}
```

## 注意事项

1. **类型安全**：`GetValue` 方法要求传入指针类型
2. **锁的生命周期**：获取锁后必须释放，建议使用 defer
3. **错误处理**：正确处理各种错误类型
4. **资源清理**：使用完客户端后调用 `Close()` 方法
5. **并发安全**：客户端是并发安全的，可以在多个 goroutine 中使用

## 与原生 go-redis 的兼容性

`PuzzleRedisClient` 嵌入了 `*redis.Client`，因此可以直接使用 go-redis 的所有原生方法：

```go
// 可以直接使用原生方法
result, err := client.Get(ctx, "key").Result()
err = client.HSet(ctx, "hash", "field", "value").Err()
pubsub := client.Subscribe(ctx, "channel")
``` 