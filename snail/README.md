# snail

snail 是一个延迟初始化工具，用于注册和管理需要延迟执行的初始化函数。它允许你将初始化逻辑分散到各个模块中，然后在合适的时机统一执行。

## 功能特性

- 延迟初始化管理
- 支持多个初始化对象注册
- 统一的初始化执行
- 错误处理和日志记录
- 简单易用的API

## 基本使用

```go
package main

import (
    "fmt"
    "github.com/go-puzzles/puzzles/snail"
)

func main() {
    // 注册数据库初始化
    snail.RegisterObject("database", func() error {
        fmt.Println("初始化数据库连接...")
        // 数据库初始化逻辑
        return nil
    })
    
    // 注册Redis初始化
    snail.RegisterObject("redis", func() error {
        fmt.Println("初始化Redis连接...")
        // Redis初始化逻辑
        return nil
    })
    
    // 注册配置初始化
    snail.RegisterObject("config", func() error {
        fmt.Println("加载配置文件...")
        // 配置加载逻辑
        return nil
    })
    
    // 执行所有注册的初始化函数
    snail.Init()
    
    fmt.Println("所有组件初始化完成!")
}
```

## 使用场景

### 模块化初始化

在大型项目中，不同模块可能有各自的初始化逻辑：

```go
// 在数据库模块中
package database

import "github.com/go-puzzles/puzzles/snail"

func init() {
    snail.RegisterObject("mysql", func() error {
        // MySQL连接初始化
        return initMySQL()
    })
    
    snail.RegisterObject("redis", func() error {
        // Redis连接初始化
        return initRedis()
    })
}

func initMySQL() error {
    // MySQL初始化逻辑
    fmt.Println("MySQL连接已建立")
    return nil
}

func initRedis() error {
    // Redis初始化逻辑
    fmt.Println("Redis连接已建立")
    return nil
}
```

```go
// 在配置模块中
package config

import "github.com/go-puzzles/puzzles/snail"

func init() {
    snail.RegisterObject("app-config", func() error {
        return loadAppConfig()
    })
}

func loadAppConfig() error {
    // 加载应用配置
    fmt.Println("应用配置已加载")
    return nil
}
```

```go
// 在主程序中
package main

import (
    _ "myapp/database" // 触发init函数
    _ "myapp/config"   // 触发init函数
    "github.com/go-puzzles/puzzles/snail"
)

func main() {
    // 执行所有模块注册的初始化函数
    snail.Init()
    
    // 启动应用
    startApplication()
}
```

### 有序初始化

通过注册顺序控制初始化顺序：

```go
func setupApplication() {
    // 1. 首先初始化配置
    snail.RegisterObject("config", loadConfig)
    
    // 2. 然后初始化日志
    snail.RegisterObject("logger", initLogger)
    
    // 3. 再初始化数据库
    snail.RegisterObject("database", initDatabase)
    
    // 4. 最后初始化缓存
    snail.RegisterObject("cache", initCache)
    
    // 按注册顺序执行初始化
    snail.Init()
}

func loadConfig() error {
    fmt.Println("1. 加载配置")
    return nil
}

func initLogger() error {
    fmt.Println("2. 初始化日志系统")
    return nil
}

func initDatabase() error {
    fmt.Println("3. 连接数据库")
    return nil
}

func initCache() error {
    fmt.Println("4. 初始化缓存")
    return nil
}
```

### 错误处理

如果某个初始化函数返回错误，snail会记录错误并终止程序：

```go
func problematicInit() error {
    return errors.New("初始化失败")
}

func main() {
    snail.RegisterObject("problematic", problematicInit)
    
    // 这将会panic并输出错误信息
    snail.Init()
}
```

### 条件初始化

可以根据条件决定是否注册某些初始化函数：

```go
func conditionalSetup() {
    // 基础组件总是需要初始化
    snail.RegisterObject("config", loadConfig)
    snail.RegisterObject("logger", initLogger)
    
    // 根据配置决定是否初始化可选组件
    if isRedisEnabled() {
        snail.RegisterObject("redis", initRedis)
    }
    
    if isElasticsearchEnabled() {
        snail.RegisterObject("elasticsearch", initElasticsearch)
    }
    
    snail.Init()
}

func isRedisEnabled() bool {
    // 检查Redis是否启用
    return true
}

func isElasticsearchEnabled() bool {
    // 检查Elasticsearch是否启用
    return false
}
```

## API 参考

### RegisterObject
```go
func RegisterObject(name string, fn func() error)
```
注册一个初始化对象。

参数：
- `name`: 对象名称，用于标识和日志记录
- `fn`: 初始化函数，返回error表示初始化是否成功

### Init
```go
func Init()
```
执行所有注册的初始化函数。按注册顺序依次执行，如果某个初始化函数返回错误，会记录错误信息并panic。

## 最佳实践

1. **明确的命名**：为每个初始化对象使用清晰的名称
2. **合理的顺序**：按依赖关系顺序注册初始化函数
3. **错误处理**：在初始化函数中正确处理和返回错误
4. **模块化设计**：将相关的初始化逻辑组织在同一个模块中
5. **条件初始化**：根据配置或环境决定是否执行某些初始化

## 注意事项

- 初始化函数应该是幂等的（可重复执行）
- 避免在初始化函数中执行耗时操作
- 确保初始化函数的错误处理正确
- 注册顺序很重要，按依赖关系注册
- 初始化失败会导致程序panic，确保生产环境的稳定性 