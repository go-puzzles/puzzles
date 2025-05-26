# pgin

pgin 是对 Gin Web 框架的扩展库，提供了统一的响应格式、强大的泛型处理器、便捷的数据绑定、会话管理、日志记录等功能，让 Gin 开发更加便捷和规范。

## 功能特性

- 统一的响应格式封装
- 强大的泛型处理器系统
- 自动数据绑定和验证
- 便捷的会话管理
- 丰富的中间件支持
- 集成的日志记录
- 扩展的引擎功能

## 主要组件

### 响应封装 (resp.go)

提供统一的API响应格式：

```go
import "github.com/go-puzzles/puzzles/pgin"

// 响应结构
type Ret struct {
    Code    int `json:"code"`
    Data    any `json:"data,omitempty"`
    Message any `json:"message,omitempty"`
}

func handler(c *gin.Context) {
    // 成功响应
    pgin.ReturnSuccess(c, data)
    
    // 错误响应
    pgin.ReturnError(c, 400, "参数错误")
    
    // 手动构造响应
    c.JSON(200, pgin.SuccessRet(data))
    c.JSON(400, pgin.ErrorRet(400, "错误信息"))
}
```

### 泛型处理器 (handler.go)

强大的泛型处理器系统，自动处理数据绑定和响应：

```go
type UserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"min=1,max=120"`
}

type UserResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// 请求处理器 - 自动绑定和验证请求参数
router.GET("/users", pgin.RequestHandler(func(c *gin.Context, req *UserRequest) {
    // req 已自动绑定和验证
    c.JSON(200, gin.H{"message": "Hello " + req.Name})
}))

// 请求响应处理器 - 自动处理请求和响应
router.POST("/users", pgin.RequestResponseHandler(func(c *gin.Context, req *UserRequest) (resp *UserResponse, err error) {
    if req.Name == "error" {
        return nil, pgin.PackError(4001, "用户创建失败")
    }
    return &UserResponse{
        ID:   1,
        Name: req.Name,
    }, nil
}))

// 响应处理器 - 只处理响应
router.GET("/config", pgin.ResponseHandler(func(c *gin.Context) (resp *ConfigResponse, err error) {
    return getConfig()
}))

// 错误处理器 - 只返回错误或成功
router.DELETE("/users/:id", pgin.ErrorReturnHandler(func(c *gin.Context) error {
    return deleteUser(c.Param("id"))
}))
```

### 模型处理器

基于接口的模型处理器：

```go
type UserHandler struct {
    Name string `form:"name"`
}

type UserResponse struct {
    Message string `json:"message"`
}

func (h UserHandler) Handle(c *gin.Context) (resp *UserResponse, err error) {
    return &UserResponse{
        Message: "Hello, " + h.Name,
    }, nil
}

// 挂载模型处理器
router.GET("/model/user", pgin.MountHandler[UserHandler]())
```

### 数据绑定 (binding.go)

自动数据绑定，支持多种绑定策略：

```go
// 自动绑定 Header、URL 参数、Query 参数和 Body
err := pgin.ParseRequestParams(c, &request)

// 验证请求参数
err := pgin.ValidateRequestParams(&request)

// 泛型处理器会自动执行这些步骤
```

### 错误处理 (error.go)

集成 perror 包的错误处理：

```go
// 创建错误
err := pgin.PackError(4001, "用户不存在")

// 在处理器中抛出错误
func handler(c *gin.Context, req *Request) (resp *Response, err error) {
    if someCondition {
        return nil, pgin.PackError(4001, "业务错误")
    }
    return &Response{}, nil
}
```

### 会话管理 (session.go)

基于 gin-contrib/sessions 的会话管理：

```go
// 初始化不同类型的存储
cookieStore := pgin.InitCookieStore()
memStore := pgin.InitMemStore()
redisStore := pgin.InitRedisStore(redisPool)

// 添加会话中间件
router.Use(pgin.NewSession("session_name", cookieStore))

func handler(c *gin.Context) {
    // 获取会话
    session := pgin.GetSession(c)
    
    // 设置值
    session.Set("user_id", "123")
    session.Save()
    
    // 获取值
    userID := session.Get("user_id")
}
```

### 日志记录 (logger.go)

集成的日志中间件：

```go
// 使用默认日志中间件
router.Use(pgin.LoggerMiddleware())

// 使用自定义日志器
router.Use(pgin.LoggerMiddleware(customLogger))
```

### 请求日志中间件 (middleware.go)

```go
// 记录请求详情
router.Use(pgin.LoggingRequest(true)) // true 表示记录 header

// 重用请求体
router.Use(pgin.ReuseBody())
```

### 引擎扩展 (engine.go)

扩展的 Gin 引擎功能：

```go
// 创建默认引擎（带日志和恢复中间件）
engine := pgin.Default()

// 创建标准服务器引擎
engine := pgin.NewStandardServerHandler()

// 创建带选项的服务器引擎
engine := pgin.NewServerHandlerWithOptions(
    pgin.WithRouters("/api", userRouter, orderRouter),
    pgin.WithMiddlewares(corsMiddleware),
    pgin.WithLoggingRequest(true),
    pgin.WithReuseBody(),
    pgin.WithServiceName("my-service"),
    pgin.WithHiddenRoutesLog(),
)
```

## 完整示例

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/go-puzzles/puzzles/pgin"
)

type UserRequest struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age" binding:"min=1,max=120"`
}

type UserResponse struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // 创建引擎
    r := pgin.Default()
    
    // 添加会话支持
    r.Use(pgin.NewSession("app_session", pgin.InitCookieStore()))
    
    // API 路由
    api := r.Group("/api")
    {
        api.GET("/users/:id", pgin.RequestHandler(getUser))
        api.POST("/users", pgin.RequestResponseHandler(createUser))
        api.PUT("/users/:id", pgin.RequestWithErrorHandler(updateUser))
        api.DELETE("/users/:id", pgin.ErrorReturnHandler(deleteUser))
    }
    
    r.Run(":8080")
}

func getUser(c *gin.Context, req *UserRequest) {
    // 获取路径参数
    id := c.Param("id")
    
    // 设置会话
    session := pgin.GetSession(c)
    session.Set("last_viewed_user", id)
    session.Save()
    
    // 返回响应
    c.JSON(200, pgin.SuccessRet(gin.H{"id": id}))
}

func createUser(c *gin.Context, req *UserRequest) (*UserResponse, error) {
    // 业务逻辑
    if req.Name == "admin" {
        return nil, pgin.PackError(4001, "用户名不可用")
    }
    
    return &UserResponse{
        ID:   1,
        Name: req.Name,
        Age:  req.Age,
    }, nil
}

func updateUser(c *gin.Context, req *UserRequest) error {
    id := c.Param("id")
    // 更新逻辑
    return nil // 成功时返回 nil
}

func deleteUser(c *gin.Context) error {
    id := c.Param("id")
    // 删除逻辑
    return nil // 成功时返回 nil
}
```

## 使用建议

- 优先使用泛型处理器，减少重复的绑定和响应代码
- 合理使用错误处理，统一错误码管理
- 根据需要选择合适的会话存储方式
- 在生产环境中使用合适的日志配置
- 利用引擎选项功能简化初始化代码
