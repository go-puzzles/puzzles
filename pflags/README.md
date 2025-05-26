# pflags

pflags 是一个强大的配置管理工具，支持从多种源加载配置参数，包括命令行参数、环境变量、配置文件等。

## 功能特性

- 支持多种数据类型：string、int、bool、float64、duration、slice等
- 支持结构体配置映射
- 支持配置文件监听和热更新
- 支持环境变量绑定
- 支持命令行参数解析
- 支持配置验证

## 基本使用

```go
package main

import (
    "fmt"
    "github.com/go-puzzles/puzzles/pflags"
)

var (
    port = pflags.Int("port", 8080, "服务端口")
    debug = pflags.Bool("debug", false, "是否开启调试模式")
    name = pflags.String("name", "myapp", "应用名称")
)

func main() {
    pflags.Parse()
    
    fmt.Printf("Port: %d\n", port())
    fmt.Printf("Debug: %t\n", debug())
    fmt.Printf("Name: %s\n", name())
}
```

## 结构体配置

```go
type Config struct {
    Port    int    `flag:"port" default:"8080" usage:"服务端口"`
    Debug   bool   `flag:"debug" default:"false" usage:"调试模式"`
    Name    string `flag:"name" default:"myapp" usage:"应用名称"`
}

var config = pflags.Struct("config", &Config{}, "struct config")

func main() {
    pflags.Parse()
    
    cfg := &Config{} 
    err := config(cfg)
    fmt.Printf("Config: %+v\n", cfg)
}
```

## 更多示例

详细使用示例请参考 [example](example/) 目录。
