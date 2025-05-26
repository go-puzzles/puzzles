# perror

perror 是一个全面的错误处理解决方案，提供错误码、错误原因追踪和错误包装功能。它扩展了标准 error 接口，为错误处理提供了更多功能。

## 功能特性

- 支持错误码定义
- 支持错误原因链追踪
- 支持错误包装和嵌套
- 提供便捷的错误创建和检查函数
- 兼容标准 error 接口
- 支持结构化错误信息

## 预定义错误码

```go
const (
    CodeUnknown      = 1000 // 未知错误
    CodeInvalidInput = 1001 // 无效输入参数
    CodeNotFound     = 1002 // 资源未找到
    CodeUnauthorized = 1003 // 未授权访问
)
```

## 快速开始

```go
package main

import (
    "errors"
    "fmt"
    "github.com/go-puzzles/puzzles/perror"
)

func main() {
    // 创建带错误码的错误
    err := perror.PackError(perror.CodeInvalidInput, "用户名不能为空")
    fmt.Printf("错误: %s\n", err.Error())
    fmt.Printf("错误码: %d\n", err.Code())

    // 包装现有错误
    originalErr := errors.New("数据库连接失败")
    wrappedErr := perror.WrapError(perror.CodeUnknown, originalErr, "用户查询失败")
    fmt.Printf("包装后的错误: %s\n", wrappedErr.Error())
    fmt.Printf("详细信息: %s\n", wrappedErr.String())
}
```

## 基本使用

### 创建错误

```go
import (
    "errors"
    "github.com/go-puzzles/puzzles/perror"
)

// 创建带错误码的错误
err := perror.PackError(perror.CodeInvalidInput, "用户名不能为空")

// 创建带错误码和原因的错误
originalErr := errors.New("数据库连接失败")
err := perror.PackError(perror.CodeUnknown, "用户查询失败", originalErr)

// 包装现有错误
wrappedErr := perror.WrapError(perror.CodeUnknown, originalErr, "用户查询失败")
```

### 错误检查

```go
// 检查是否为 ErrorR 类型
if errR, ok := perror.AsErrorR(err); ok {
    fmt.Printf("错误码: %d\n", errR.Code())
    fmt.Printf("错误信息: %s\n", errR.Error())
    if cause := errR.Cause(); cause != nil {
        fmt.Printf("原因: %s\n", cause.Error())
    }
}

// 获取错误码
code := perror.GetErrorCode(err)
if code == perror.CodeInvalidInput {
    // 处理无效输入错误
}
```

## 接口定义

### ErrorCoder
提供错误码功能：

```go
type ErrorCoder interface {
    error
    Code() int
}
```

### ErrorCauser
提供错误原因追踪功能：

```go
type ErrorCauser interface {
    error
    Cause() error
}
```

### ErrorR
组合接口，提供完整的错误功能：

```go
type ErrorR interface {
    ErrorCoder
    ErrorCauser
    fmt.Stringer
}
```

## 详细示例

### Web API 错误处理

```go
package main

import (
    "errors"
    "fmt"
    "strings"
    "github.com/go-puzzles/puzzles/perror"
)

func validateUser(name, email string) error {
    if name == "" {
        return perror.PackError(perror.CodeInvalidInput, "用户名不能为空")
    }
    
    if email == "" {
        return perror.PackError(perror.CodeInvalidInput, "邮箱不能为空")
    }
    
    if !isValidEmail(email) {
        return perror.PackError(perror.CodeInvalidInput, "邮箱格式不正确")
    }
    
    return nil
}

func createUser(name, email string) error {
    // 验证用户输入
    if err := validateUser(name, email); err != nil {
        return err // 直接返回验证错误
    }
    
    // 模拟数据库操作
    if err := saveToDatabase(name, email); err != nil {
        return perror.WrapError(perror.CodeUnknown, err, "用户创建失败")
    }
    
    return nil
}

func handleCreateUser(name, email string) {
    if err := createUser(name, email); err != nil {
        code := perror.GetErrorCode(err)
        
        switch code {
        case perror.CodeInvalidInput:
            fmt.Printf("输入错误: %s\n", err.Error())
        case perror.CodeUnknown:
            fmt.Printf("系统错误: %s\n", err.Error())
            if errR, ok := perror.AsErrorR(err); ok {
                if cause := errR.Cause(); cause != nil {
                    fmt.Printf("详细错误: %s\n", cause.Error())
                }
            }
        default:
            fmt.Printf("未知错误: %s\n", err.Error())
        }
        return
    }
    
    fmt.Println("用户创建成功")
}

func isValidEmail(email string) bool {
    // 简单的邮箱格式验证
    return strings.Contains(email, "@")
}

func saveToDatabase(name, email string) error {
    // 模拟数据库操作可能出现的错误
    return errors.New("连接超时")
}
```

### 业务逻辑错误处理

```go
package main

import (
    "errors"
    "fmt"
    "github.com/go-puzzles/puzzles/perror"
)

const (
    // 业务错误码
    CodeUserNotFound      = 2001
    CodeUserExists        = 2002
    CodeInsufficientFunds = 2003
)

type Account struct {
    ID      string
    Balance float64
}

func transferMoney(fromID, toID string, amount float64) error {
    // 检查发送者账户
    fromAccount, err := getAccount(fromID)
    if err != nil {
        return perror.WrapError(CodeUserNotFound, err, "发送者账户不存在")
    }
    
    // 检查接收者账户
    toAccount, err := getAccount(toID)
    if err != nil {
        return perror.WrapError(CodeUserNotFound, err, "接收者账户不存在")
    }
    
    // 检查余额
    if fromAccount.Balance < amount {
        return perror.PackError(CodeInsufficientFunds, "余额不足")
    }
    
    // 执行转账
    if err := executeTransfer(fromAccount, toAccount, amount); err != nil {
        return perror.WrapError(perror.CodeUnknown, err, "转账执行失败")
    }
    
    return nil
}

func handleTransfer(fromID, toID string, amount float64) {
    if err := transferMoney(fromID, toID, amount); err != nil {
        code := perror.GetErrorCode(err)
        
        switch code {
        case CodeUserNotFound:
            fmt.Printf("用户不存在: %s\n", err.Error())
        case CodeInsufficientFunds:
            fmt.Printf("余额不足: %s\n", err.Error())
        default:
            fmt.Printf("转账失败: %s\n", err.Error())
            // 记录详细错误日志
            if errR, ok := perror.AsErrorR(err); ok {
                fmt.Printf("详细信息: %s\n", errR.String())
            }
        }
        return
    }
    
    fmt.Println("转账成功")
}

func getAccount(id string) (*Account, error) {
    // 模拟账户查询
    if id == "" {
        return nil, errors.New("账户ID不能为空")
    }
    
    // 模拟账户不存在的情况
    if id == "nonexistent" {
        return nil, errors.New("账户不存在")
    }
    
    return &Account{
        ID:      id,
        Balance: 1000.0,
    }, nil
}

func executeTransfer(from, to *Account, amount float64) error {
    // 模拟转账操作可能出现的错误
    if amount <= 0 {
        return errors.New("转账金额必须大于0")
    }
    
    // 模拟系统错误
    return errors.New("系统繁忙，请稍后重试")
}
```

## API 参考

### PackError
```go
func PackError(code int, vals ...any) ErrorR
```
创建一个新的错误，支持传入错误码、消息和原因错误。

### WrapError
```go
func WrapError(code int, err error, message string) ErrorR
```
包装一个现有错误，添加错误码和新的消息。

### AsErrorR
```go
func AsErrorR(err error) (ErrorR, bool)
```
尝试将标准 error 转换为 ErrorR 接口。

### GetErrorCode
```go
func GetErrorCode(err error) int
```
获取错误的错误码，如果不是 ErrorR 类型则返回 CodeUnknown。

## 最佳实践

1. **定义明确的错误码**：为不同类型的错误定义专门的错误码
2. **保留错误链**：使用 WrapError 保留原始错误信息
3. **分层错误处理**：在不同层级添加适当的上下文信息
4. **统一错误响应**：在 API 层统一处理和响应错误
5. **记录详细日志**：使用 String() 方法记录完整的错误信息

## 注意事项

- 错误码应该有明确的语义
- 避免过度嵌套错误链
- 在记录日志时使用 String() 方法获取完整信息
- 在返回给用户时使用 Error() 方法获取简洁信息 