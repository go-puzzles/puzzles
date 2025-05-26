# penum

penum 是一个枚举类型生成工具，支持自动生成字符串和数字类型的枚举。它通过反射自动为结构体字段设置合适的枚举值，简化了枚举类型的定义和使用。

## 功能特性

- 自动生成枚举值
- 支持字符串和数字类型枚举
- 基于反射的智能类型处理
- 支持自定义选项配置
- 简单易用的API
- 类型安全的枚举操作

## 基本使用

### 字符串枚举

```go
package main

import (
    "fmt"
    "github.com/go-puzzles/puzzles/penum"
)

// 定义状态枚举
type Status struct {
    Pending   string
    Running   string
    Completed string
    Failed    string
}

func main() {
    // 生成枚举，字段名将作为字符串值
    status := penum.New[Status]()
    
    fmt.Println(status.Pending)   // "Pending"
    fmt.Println(status.Running)   // "Running"
    fmt.Println(status.Completed) // "Completed"
    fmt.Println(status.Failed)    // "Failed"
}
```

### 数字枚举

```go
// 定义优先级枚举
type Priority struct {
    Low    int
    Medium int
    High   int
    Urgent int
}

func main() {
    priority := penum.New[Priority]()
    
    fmt.Println(priority.Low)    // 0
    fmt.Println(priority.Medium) // 1
    fmt.Println(priority.High)   // 2
    fmt.Println(priority.Urgent) // 3
}
```

### 混合类型枚举

```go
// 混合不同数据类型
type TaskType struct {
    ID          int
    Name        string
    Code        uint
    Description string
}

func main() {
    taskType := penum.New[TaskType]()
    
    fmt.Println(taskType.ID)          // 0
    fmt.Println(taskType.Name)        // "Name"
    fmt.Println(taskType.Code)        // 2
    fmt.Println(taskType.Description) // "Description"
}
```

## 实际应用场景

### HTTP 状态码枚举

```go
type HTTPStatus struct {
    OK                  int
    BadRequest          int
    Unauthorized        int
    Forbidden           int
    NotFound            int
    InternalServerError int
}

func getHTTPStatusEnum() HTTPStatus {
    return penum.New[HTTPStatus](func(status *HTTPStatus) {
        // 可以通过选项函数自定义值
        status.OK = 200
        status.BadRequest = 400
        status.Unauthorized = 401
        status.Forbidden = 403
        status.NotFound = 404
        status.InternalServerError = 500
    })
}
```

### 用户角色枚举

```go
type UserRole struct {
    Guest  string
    User   string
    Admin  string
    Super  string
}

func getUserRoles() UserRole {
    return penum.New[UserRole]()
}

// 使用示例
func checkPermission(userRole string) bool {
    roles := getUserRoles()
    
    switch userRole {
    case roles.Admin, roles.Super:
        return true
    case roles.User:
        return false
    default:
        return false
    }
}
```

### 订单状态枚举

```go
type OrderStatus struct {
    Created   string
    Paid      string
    Shipped   string
    Delivered string
    Cancelled string
    Refunded  string
}

type Order struct {
    ID     int
    Status string
    Amount float64
}

func processOrder(order *Order) {
    status := penum.New[OrderStatus]()
    
    switch order.Status {
    case status.Created:
        fmt.Println("处理新订单")
    case status.Paid:
        fmt.Println("订单已支付，准备发货")
    case status.Shipped:
        fmt.Println("订单已发货")
    case status.Delivered:
        fmt.Println("订单已送达")
    case status.Cancelled:
        fmt.Println("订单已取消")
    case status.Refunded:
        fmt.Println("订单已退款")
    }
}
```

### 数据库操作枚举

```go
type DBOperation struct {
    Create string
    Read   string
    Update string
    Delete string
}

func auditLog(operation string, table string, userID int) {
    ops := penum.New[DBOperation]()
    
    var action string
    switch operation {
    case ops.Create:
        action = "创建"
    case ops.Read:
        action = "查询"
    case ops.Update:
        action = "更新"
    case ops.Delete:
        action = "删除"
    default:
        action = "未知操作"
    }
    
    fmt.Printf("用户 %d 对表 %s 执行了 %s 操作\n", userID, table, action)
}
```

## 高级功能

### 自定义选项

```go
type CustomEnum struct {
    First  string
    Second string
    Third  string
}

func main() {
    // 使用选项函数自定义枚举值
    enum := penum.New[CustomEnum](func(e *CustomEnum) {
        e.First = "自定义第一个"
        e.Second = "自定义第二个"
        e.Third = "自定义第三个"
    })
    
    fmt.Println(enum.First)  // "自定义第一个"
    fmt.Println(enum.Second) // "自定义第二个"
    fmt.Println(enum.Third)  // "自定义第三个"
}
```

### 枚举验证

```go
type Color struct {
    Red   string
    Green string
    Blue  string
}

var colors = penum.New[Color]()

func isValidColor(color string) bool {
    return color == colors.Red || 
           color == colors.Green || 
           color == colors.Blue
}

func getAllColors() []string {
    return []string{
        colors.Red,
        colors.Green,
        colors.Blue,
    }
}
```

### 枚举映射

```go
type ErrorCode struct {
    Success       int
    InvalidInput  int
    NotFound      int
    ServerError   int
}

func getErrorMessage(code int) string {
    errorCodes := penum.New[ErrorCode](func(e *ErrorCode) {
        e.Success = 0
        e.InvalidInput = 1001
        e.NotFound = 1002
        e.ServerError = 1003
    })
    
    messages := map[int]string{
        errorCodes.Success:      "操作成功",
        errorCodes.InvalidInput: "输入参数无效",
        errorCodes.NotFound:     "资源未找到",
        errorCodes.ServerError:  "服务器内部错误",
    }
    
    if msg, exists := messages[code]; exists {
        return msg
    }
    return "未知错误"
}
```

## API 参考

### New[T any](opts ...EnumOption[T]) T

创建一个新的枚举实例。

类型参数：

- `T`: 枚举结构体类型

参数：

- `opts`: 可选的配置函数

返回：

- `T`: 初始化后的枚举实例

### EnumOption[T any]

枚举配置选项函数类型：

```go
type EnumOption[T any] func(*T)
```

## 支持的数据类型

- `string`: 字段名作为字符串值
- `int`, `int8`, `int16`, `int32`, `int64`: 基于字段在结构体中的索引位置，从0开始递增
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`: 基于字段在结构体中的索引位置，从0开始递增

## 最佳实践

1. **命名规范**：使用清晰的字段名，因为它们会成为枚举值
2. **类型一致性**：在同一个枚举中尽量使用相同类型的字段
3. **文档注释**：为枚举类型和字段添加适当的注释
4. **验证函数**：提供枚举值验证函数
5. **常量导出**：考虑将枚举实例定义为包级别的常量，避免重复创建
6. **性能考虑**：避免在函数内部重复调用 `penum.New[T]()`，应该在包级别创建一次并复用

## 注意事项

- 枚举结构体必须是结构体类型，不能是其他类型
- 字段必须是可设置的（导出字段）
- 不支持复合类型（slice、map、struct等）
- 字段的零值会被枚举值覆盖
- 使用选项函数可以覆盖默认的枚举值
