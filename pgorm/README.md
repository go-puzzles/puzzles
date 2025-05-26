# pgorm

pgorm 是一个基于 GORM 的数据库 ORM 工具库，提供了 MySQL 和 SQLite 数据库的便捷操作接口，支持全局模型注册和统一的客户端管理。

## 功能特性

- 支持 MySQL 和 SQLite 数据库
- 支持全局模型注册管理
- 集成日志记录功能
- 简化的数据库连接配置
- 支持慢 SQL 日志记录
- 提供连接选项配置
- 支持数据库连接预检和自动迁移

## 基本使用

### MySQL 连接

#### 方式一：直接使用配置连接

```go
package main

import (
    "context"
    "github.com/go-puzzles/puzzles/pgorm"
)

func main() {
    conf := &pgorm.MysqlConfig{
        Instance: "localhost:3306",
        Database: "test",
        Username: "root",
        Password: "password",
    }

    db, err := conf.DialGorm()
    if err != nil {
        panic(err)
    }

    db.WithContext(context.Background())
}
```

#### 方式二：使用全局注册

```go
package main

import (
    "context"
    "github.com/go-puzzles/puzzles/pgorm"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name string
}

func (u *User) TableName() string {
    return "users"
}

func main() {
    conf := &pgorm.MysqlConfig{
        Instance: "localhost:3306",
        Database: "test",
        Username: "root",
        Password: "password",
    }
    
    // 注册模型到全局
    err := pgorm.RegisterSqlModelWithConf(conf, &User{})
    if err != nil {
        panic(err)
    }
    
    // 通过配置获取数据库连接
    db := pgorm.GetDbByConf(conf)
    
    // 或通过模型获取数据库连接
    db = pgorm.GetDbByModel(&User{})
    
    db.WithContext(context.Background())
}
```

#### 使用 DSN 连接

```go
dsnConf := &pgorm.MysqlDsn{
    DSN: "root:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
}

db, err := dsnConf.DialGorm()
if err != nil {
    panic(err)
}
```

### SQLite 连接

```go
conf := &pgorm.SqliteConfig{
    DbFile: "test.db",
}

db, err := conf.DialGorm()
if err != nil {
    panic(err)
}
```

## 连接选项配置

### 可用选项

```go
// 设置日志前缀
pgorm.WithLogPrefix("my-app")

// 忽略 RecordNotFound 错误日志
pgorm.WithDialIgnoreNotFound()

// 设置慢 SQL 阈值
pgorm.WithDialThreshold(500 * time.Millisecond)

// 使用示例
db, err := conf.DialGorm(
    pgorm.WithLogPrefix("user-service"),
    pgorm.WithDialIgnoreNotFound(),
    pgorm.WithDialThreshold(300 * time.Millisecond),
)
```

## 模型定义

### 定义模型接口

模型需要实现 `SqlModel` 接口：

```go
type SqlModel interface {
    TableName() string
}

type User struct {
    gorm.Model
    Name   string `gorm:"size:100;not null" json:"name"`
    Email  string `gorm:"size:100;unique;not null" json:"email"`
    Age    int    `gorm:"default:0" json:"age"`
    Status int    `gorm:"default:1" json:"status"`
}

func (u *User) TableName() string {
    return "users"
}
```

### 模型注册

```go
// 推荐：使用带错误返回的注册方法
err := pgorm.RegisterSqlModelWithConf(conf, &User{}, &Post{})
if err != nil {
    log.Fatal(err)
}

// 不推荐：直接 panic 的注册方法（已废弃）
pgorm.RegisterSqlModel(conf, &User{}, &Post{})
```

## 数据库操作

### 基本 CRUD

```go
// 创建
user := &User{
    Name:  "张三",
    Email: "zhangsan@example.com",
    Age:   25,
}

db := pgorm.GetDbByModel(user)
result := db.Create(user)
if result.Error != nil {
    panic(result.Error)
}

// 查询
var user User
db = pgorm.GetDbByModel(&user)
db.First(&user, 1) // 根据主键查询
db.Where("email = ?", "zhangsan@example.com").First(&user) // 条件查询

// 更新
db.Update("age", 26)
db.Updates(User{Age: 27, Status: 2})

// 删除
db.Delete(&user, 1)
```

### 事务处理

```go
db := pgorm.GetDbByConf(conf)

err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    
    return nil
})
```

## 数据库管理工具

### 连接预检

```go
// 检查数据库连接是否正常
err := pgorm.PrePing(conf)
if err != nil {
    log.Fatal("数据库连接失败:", err)
}
```

### 自动迁移

```go
// 自动迁移已注册的模型
err := pgorm.AutoMigrate(conf)
if err != nil {
    log.Fatal("自动迁移失败:", err)
}
```

## 工具函数

### 错误处理

```go
// 检查是否为记录未找到错误
isNotFound, err := pgorm.GormIsErrRecordNotFound(dbErr)
if err != nil {
    log.Error("数据库错误:", err)
} else if isNotFound {
    log.Info("记录未找到")
}
```

## 配置选项

### MySQL 配置

```go
type MysqlConfig struct {
    Instance string // 实例地址 (host:port)
    Database string // 数据库名
    Username string // 用户名
    Password string // 密码
}
```

### MySQL DSN 配置

```go
type MysqlDsn struct {
    DSN string // 完整的 DSN 连接字符串
}
```

### SQLite 配置

```go
type SqliteConfig struct {
    DbFile string // 数据库文件路径
}
```

### 连接选项

```go
type DialOption struct {
    LogPrefix            string        // 日志前缀
    IgnoreRecordNotFound bool          // 是否忽略记录未找到错误
    SlowThreshold        time.Duration // 慢 SQL 阈值
}
```

## 最佳实践

1. **模型注册**：在应用启动时统一注册所有模型
2. **错误处理**：使用 `RegisterSqlModelWithConf` 而不是 `RegisterSqlModel`
3. **连接管理**：利用全局注册机制管理多个数据库连接
4. **日志配置**：合理设置慢 SQL 阈值和日志前缀
5. **预检机制**：在应用启动时使用 `PrePing` 检查连接
6. **自动迁移**：在需要时使用 `AutoMigrate` 进行数据库迁移

## 注意事项

- 模型必须实现 `SqlModel` 接口
- 配置对象需要实现 `Config` 接口
- 全局注册的模型和配置具有唯一性约束
- 日志记录依赖 `plog` 包的 Debug 模式
- 连接池参数通过底层 `dialer` 包管理 