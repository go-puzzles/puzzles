# plog

一个简单易用的 Go 日志库，支持多种日志级别、上下文日志记录和文件输出。

## 特性

- 🚀 简单易用的 API
- 📊 支持多种日志级别 (Debug, Info, Warn, Error, Fatal)
- 🎯 支持格式化日志和上下文日志
- 📁 支持日志文件输出和自动轮转
- 🔧 可插拔的日志实现 (默认实现 + slog 实现)
- 🏷️ 支持上下文中的键值对存储
- ⚡ 高性能，低内存占用

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "github.com/go-puzzles/puzzles/plog"
)

func main() {
    // 基本日志记录
    plog.Infof("这是一条信息日志")
    plog.Errorf("这是一条错误日志: %v", "错误信息")
    
    // 上下文日志记录
    ctx := context.Background()
    plog.Infoc(ctx, "这是一条上下文日志: %v", 123)
}
```

### 日志级别

```go
// 设置日志级别
plog.Enable(level.LevelDebug)

// 检查是否启用 Debug 级别
if plog.IsDebug() {
    plog.Debugf("Debug 模式已启用")
}

// 支持的日志级别
plog.Debugf("调试日志")   // -4
plog.Infof("信息日志")    // 0
plog.Warnf("警告日志")    // 4
plog.Errorf("错误日志")   // 8
plog.Fatalf("致命日志")   // 程序会退出
```

### 文件输出

```go
// 配置日志文件
logConfig := &plog.LogConfig{
    Filename:   "app.log",     // 日志文件名
    MaxSize:    10,            // 最大文件大小 (MB)
    MaxBackups: 3,             // 保留的备份文件数
    MaxAge:     28,            // 文件保留天数
    LocalTime:  true,          // 使用本地时间
    Compress:   true,          // 压缩备份文件
}

// 设置默认值
logConfig.SetDefault()

// 启用文件输出
plog.EnableLogToFile(logConfig)
```

### 使用 slog 后端

```go
// 切换到 slog 实现
plog.SetSlog()

// 正常使用所有 API
plog.Infof("使用 slog 后端记录日志")
```

### 上下文增强

使用 `With` 函数可以在上下文中添加键值对信息：

```go
// 添加分组
ctx := plog.With(context.Background(), "用户操作")

// 添加键值对
ctx = plog.With(ctx, "用户ID", "12345")
ctx = plog.With(ctx, "操作", "登录", "IP", "192.168.1.1")

// 使用增强的上下文记录日志
plog.Infoc(ctx, "用户登录成功")
```

### 错误处理

```go
err := someFunction()
// 如果 err 不为 nil，记录错误并 panic
plog.PanicError(err, "执行某个函数时出错")
```

### 自定义 Logger

```go
// 获取当前 Logger
logger := plog.GetLogger()

// 设置自定义 Logger
customLogger := log.New(log.WithCalldepth(3))
plog.SetLogger(customLogger)

// 设置输出目标
plog.SetOutput(os.Stdout)
```

## API 参考

### 全局函数

#### 配置函数

- `SetSlog()` - 切换到 slog 实现
- `GetLogger() Logger` - 获取当前 Logger 实例
- `SetLogger(l Logger)` - 设置自定义 Logger
- `SetOutput(w io.Writer)` - 设置输出目标
- `Enable(l level.Level)` - 设置日志级别
- `EnableLogToFile(jackLog *LogConfig)` - 启用文件输出
- `IsDebug() bool` - 检查是否启用 Debug 级别

#### 格式化日志函数

- `Debugf(msg string, v ...any)` - Debug 级别格式化日志
- `Infof(msg string, v ...any)` - Info 级别格式化日志
- `Warnf(msg string, v ...any)` - Warn 级别格式化日志
- `Errorf(msg string, v ...any)` - Error 级别格式化日志
- `Fatalf(msg string, v ...any)` - Fatal 级别格式化日志

#### 上下文日志函数

- `Debugc(ctx context.Context, msg string, v ...any)` - Debug 级别上下文日志
- `Infoc(ctx context.Context, msg string, v ...any)` - Info 级别上下文日志
- `Warnc(ctx context.Context, msg string, v ...any)` - Warn 级别上下文日志
- `Errorc(ctx context.Context, msg string, v ...any)` - Error 级别上下文日志
- `Fatalc(ctx context.Context, msg string, v ...any)` - Fatal 级别上下文日志

#### 上下文增强函数

- `With(c context.Context, msg string, v ...any) context.Context` - 添加日志上下文信息
- `WithLogger(c context.Context, w io.Writer) context.Context` - 在上下文中设置 Logger

#### 工具函数

- `PanicError(err error, v ...any)` - 如果错误不为 nil 则记录并 panic

### 日志级别

```go
const (
    LevelDebug Level = -4
    LevelInfo  Level = 0
    LevelWarn  Level = 4
    LevelError Level = 8
)
```

### LogConfig 结构

基于 `lumberjack.Logger` 的日志文件配置：

```go
type LogConfig struct {
    Filename   string // 日志文件路径
    MaxSize    int    // 最大文件大小 (MB)
    MaxAge     int    // 文件保留天数
    MaxBackups int    // 保留的备份文件数
    LocalTime  bool   // 是否使用本地时间
    Compress   bool   // 是否压缩备份文件
}
```

## 许可证

本项目采用 MIT 许可证。详见 LICENSE 文件。
