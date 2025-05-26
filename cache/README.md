# cache

cache 是一个缓存接口和实现库，提供统一的缓存操作接口，支持内存缓存和TTL功能。

## 功能特性

- 统一的缓存接口设计
- 支持内存缓存实现
- 支持TTL（生存时间）功能
- 支持缓存创建器模式
- 自动序列化和反序列化
- 简单易用的API

## 基本使用

### 内存缓存

```go
package main

import (
    "fmt"
    "time"
    "github.com/go-puzzles/puzzles/cache"
)

func main() {
    // 创建内存缓存
    c := cache.NewMemoryCache()
    defer c.Close()
    
    // 设置缓存
    err := c.Set("key1", "value1")
    if err != nil {
        panic(err)
    }
    
    // 获取缓存
    var value string
    err = c.Get("key1", &value)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Value: %s\n", value)
    
    // 检查缓存是否存在
    exists := c.Exists("key1")
    fmt.Printf("Key exists: %t\n", exists)
    
    // 删除缓存
    c.Delete("key1")
}
```

### TTL缓存

```go
func ttlCacheExample() {
    // 创建支持TTL的内存缓存，清理间隔为1分钟
    c := cache.NewMemoryCacheWithTTL(time.Minute)
    defer c.Close()
    
    // 设置带TTL的缓存
    err := c.SetWithTTL("session:123", "user_data", time.Minute*30)
    if err != nil {
        panic(err)
    }
    
    // 获取缓存
    var sessionData string
    err = c.Get("session:123", &sessionData)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Session data: %s\n", sessionData)
}
```

## 缓存创建器模式

### GetOrCreate

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func cacheWithCreator() {
    c := cache.NewMemoryCache()
    defer c.Close()
    
    // 获取或创建缓存
    var user User
    err := c.GetOrCreate("user:123", func() (any, error) {
        // 模拟从数据库获取用户信息
        return User{
            ID:   123,
            Name: "张三",
            Age:  25,
        }, nil
    }, &user)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User: %+v\n", user)
}
```

### GetOrCreateWithTTL

```go
func cacheWithTTLCreator() {
    c := cache.NewMemoryCacheWithTTL(time.Minute)
    defer c.Close()
    
    // 获取或创建带TTL的缓存
    var user User
    err := c.GetOrCreateWithTTL("user:456", func() (any, error) {
        // 模拟从数据库获取用户信息
        return getUserFromDatabase(456)
    }, &user, time.Hour)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User with TTL: %+v\n", user)
}

func getUserFromDatabase(id int) (User, error) {
    // 模拟数据库查询
    return User{
        ID:   id,
        Name: "李四",
        Age:  30,
    }, nil
}
```

## 复杂数据类型缓存

### 结构体缓存

```go
type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
}

func structCacheExample() {
    c := cache.NewMemoryCache()
    defer c.Close()
    
    // 缓存结构体
    product := Product{
        ID:          1,
        Name:        "笔记本电脑",
        Price:       5999.99,
        Description: "高性能笔记本电脑",
    }
    
    err := c.Set("product:1", product)
    if err != nil {
        panic(err)
    }
    
    // 获取结构体
    var cachedProduct Product
    err = c.Get("product:1", &cachedProduct)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Product: %+v\n", cachedProduct)
}
```

### 切片缓存

```go
func sliceCacheExample() {
    c := cache.NewMemoryCache()
    defer c.Close()
    
    // 缓存切片
    users := []User{
        {ID: 1, Name: "张三", Age: 25},
        {ID: 2, Name: "李四", Age: 30},
        {ID: 3, Name: "王五", Age: 35},
    }
    
    err := c.Set("users:active", users)
    if err != nil {
        panic(err)
    }
    
    // 获取切片
    var cachedUsers []User
    err = c.Get("users:active", &cachedUsers)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Users: %+v\n", cachedUsers)
}
```

### Map缓存

```go
func mapCacheExample() {
    c := cache.NewMemoryCache()
    defer c.Close()
    
    // 缓存Map
    config := map[string]interface{}{
        "app_name":    "MyApp",
        "version":     "1.0.0",
        "debug":       true,
        "max_users":   1000,
        "timeout":     30.5,
    }
    
    err := c.Set("app:config", config)
    if err != nil {
        panic(err)
    }
    
    // 获取Map
    var cachedConfig map[string]interface{}
    err = c.Get("app:config", &cachedConfig)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Config: %+v\n", cachedConfig)
}
```

## 错误处理

### 缓存未命中处理

```go
func handleCacheMiss() {
    c := cache.NewMemoryCache()
    defer c.Close()
    
    var value string
    err := c.Get("nonexistent", &value)
    if err != nil {
        fmt.Printf("缓存未命中: %v\n", err)
        
        // 从其他数据源获取数据
        value = "从数据库获取的数据"
        
        // 设置到缓存中
        c.Set("nonexistent", value)
    }
    
    fmt.Printf("Value: %s\n", value)
}
```

### 创建器错误处理

```go
func handleCreatorError() {
    c := cache.NewMemoryCache()
    defer c.Close()
    
    var user User
    err := c.GetOrCreate("user:999", func() (any, error) {
        // 模拟创建器返回错误
        return nil, fmt.Errorf("用户不存在")
    }, &user)
    
    if err != nil {
        fmt.Printf("创建器错误: %v\n", err)
        return
    }
    
    fmt.Printf("User: %+v\n", user)
}
```

## 接口定义

### Cache 接口

```go
type Cache interface {
    Get(key string, out any) error
    Set(key string, value any) error
    GetOrCreate(key string, creater Creater, out any) error
    Exists(key string) bool
    Delete(key string)
    Close() error
}
```

### CacheWithTTL 接口

```go
type CacheWithTTL interface {
    Cache
    GetOrCreateWithTTL(key string, creater Creater, out any, ttl time.Duration) error
    SetWithTTL(key string, value any, ttl time.Duration) error
}
```

### Creater 函数类型

```go
type Creater func() (any, error)
```

## 实现类型

### 内存缓存

- `NewMemoryCache()`: 创建基本内存缓存
- `NewMemoryCacheWithTTL(compactInterval time.Duration)`: 创建支持TTL的内存缓存，compactInterval参数设置清理过期缓存的间隔时间

## 最佳实践

1. **合理设置TTL**：根据数据的时效性设置合适的过期时间
2. **错误处理**：正确处理缓存未命中和创建器错误
3. **内存管理**：定期清理过期的缓存项
4. **键命名规范**：使用有意义的键名规范
5. **数据序列化**：确保缓存的数据可以正确序列化和反序列化
6. **并发安全**：在多协程环境中注意缓存的并发安全性

## 注意事项

- 内存缓存在程序重启后会丢失
- 大量缓存数据可能占用较多内存
- TTL缓存会自动清理过期项
- 缓存的数据必须是可序列化的
- 在高并发场景下注意缓存的性能影响 