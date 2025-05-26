# dialer

dialer 是一个网络连接工具包，提供数据库连接(MySQL、SQLite、Redis)、gRPC连接等功能的统一接口和配置管理，支持服务发现机制。

## 功能特性

- 统一的连接配置管理
- 支持多种数据库连接 (MySQL、SQLite、Redis)
- 支持 gRPC 客户端连接
- 集成服务发现机制
- 提供连接池配置选项
- 集成 GORM 配置
- 简化的连接建立流程

## 基本配置

## 数据库连接

### MySQL 连接

```go
package main

import (
    "github.com/go-puzzles/puzzles/dialer"
    "github.com/go-puzzles/puzzles/dialer/mysql"
)

func connectMySQL() {
    // 使用服务发现连接 MySQL
    db, err := mysql.DialMysqlGorm("mysql-service", 
        dialer.WithAuth("root", "password"),
        dialer.WithDBName("myapp"),
    )
    if err != nil {
        panic(err)
    }
    defer func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }()
    
    // 配置连接池
    sqlDB, _ := db.DB()
    dialer.ConfigDB(sqlDB) // 使用默认连接池配置
    
    // 使用数据库
    var users []User
    db.Find(&users)
}

func connectMySQLWithDSN() {
    // 直接使用 DSN 连接
    dsn := "root:password@tcp(localhost:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := mysql.DialMysqlGormWithDSN(dsn,
        dialer.WithAuth("root", "password"),
    )
    if err != nil {
        panic(err)
    }
    defer func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }()
}
```

### SQLite 连接

```go
import "github.com/go-puzzles/puzzles/dialer/sqlite"

func connectSQLite() {
    // 基本连接
    db, err := sqlite.DialSqlLiteGorm("app.db",
        dialer.WithSqliteCache(), // 启用缓存模式
    )
    if err != nil {
        panic(err)
    }
    defer func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }()
    
    // 自动迁移
    db.AutoMigrate(&User{}, &Post{})
}
```

### Redis 连接

```go
import "github.com/go-puzzles/puzzles/dialer/redis"

func connectRedis() {
    // 创建 Redis 连接池
    pool := redis.DialRedisPool("redis-service", 0, 10, "password")
    defer pool.Close()
    
    // 获取连接
    conn := pool.Get()
    defer conn.Close()
    
    // 使用连接
    _, err := conn.Do("SET", "key", "value")
    if err != nil {
        panic(err)
    }
    
    result, err := conn.Do("GET", "key")
    if err != nil {
        panic(err)
    }
    fmt.Println("Result:", result)
}
```

## gRPC 连接

### 客户端连接

```go
import (
    "context"
    "time"
    "github.com/go-puzzles/puzzles/dialer/grpc"
    "google.golang.org/grpc"
)

func connectGRPC() {
    // 使用服务发现创建gRPC连接
    conn, err := grpc.DialGrpc("user-service")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    // 创建客户端
    client := pb.NewUserServiceClient(conn)
    
    // 调用服务
    ctx := context.Background()
    resp, err := client.GetUser(ctx, &pb.GetUserRequest{
        Id: 123,
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User: %+v\n", resp.User)
}

func connectGRPCWithTimeout() {
    // 带超时的连接
    conn, err := grpc.DialGrpcWithTimeOut(5*time.Second, "user-service")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}

func connectGRPCWithTag() {
    // 带标签的服务发现
    conn, err := grpc.DialGrpcWithTag("user-service", "v2")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}
```

## 高级配置

### 自定义 GORM 配置

```go
func customGormConfig() {
    opt := dialer.PackDialOption(
        dialer.WithAuth("root", "password"),
        dialer.WithDBName("myapp"),
    )
    
    // 使用自定义GORM配置
    config := &gorm.Config{
        PrepareStmt:                              true,
        DisableForeignKeyConstraintWhenMigrating: true,
        Logger:                                   logger.Default.LogMode(logger.Info),
    }
    
    // 或者使用默认配置
    defaultConfig := dialer.DefaultGormConfig(opt)
    
    // 直接使用 gorm.Open
    dsn := "root:password@tcp(localhost:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), config)
    if err != nil {
        panic(err)
    }
}
```

### 连接池优化

```go
func optimizeConnectionPool() {
    db, err := mysql.DialMysqlGorm("mysql-service",
        dialer.WithAuth("root", "password"),
        dialer.WithDBName("myapp"),
    )
    if err != nil {
        panic(err)
    }
    
    // 获取底层sql.DB
    sqlDB, err := db.DB()
    if err != nil {
        panic(err)
    }
    
    // 自定义连接池配置
    sqlDB.SetMaxIdleConns(20)               // 最大空闲连接数
    sqlDB.SetMaxOpenConns(200)              // 最大打开连接数
    sqlDB.SetConnMaxLifetime(time.Hour * 2) // 连接最大生存时间
    
    // 或使用默认配置
    dialer.ConfigDB(sqlDB)
}
```

### 多数据库连接

```go
type DatabaseManager struct {
    primaryDB   *gorm.DB
    secondaryDB *gorm.DB
    cacheDB     *gorm.DB
    redisPool   *redis.Pool
}

func setupMultipleDatabases() *DatabaseManager {
    // 主数据库 (MySQL)
    primaryDB, err := mysql.DialMysqlGorm("primary-mysql",
        dialer.WithAuth("root", "primary_password"),
        dialer.WithDBName("primary_db"),
    )
    if err != nil {
        panic(err)
    }
    
    // 从数据库 (MySQL)
    secondaryDB, err := mysql.DialMysqlGorm("secondary-mysql",
        dialer.WithAuth("root", "secondary_password"),
        dialer.WithDBName("secondary_db"),
    )
    if err != nil {
        panic(err)
    }
    
    // 缓存数据库 (SQLite)
    cacheDB, err := sqlite.DialSqlLiteGorm("cache.db",
        dialer.WithSqliteCache(),
    )
    if err != nil {
        panic(err)
    }
    
    // Redis 连接池
    redisPool := redis.DialRedisPool("redis-service", 0, 10, "redis_password")
    
    return &DatabaseManager{
        primaryDB:   primaryDB,
        secondaryDB: secondaryDB,
        cacheDB:     cacheDB,
        redisPool:   redisPool,
    }
}
```

## 服务发现

dialer 集成了服务发现机制，可以通过服务名而不是具体地址来连接服务：

```go
// MySQL 服务发现
db, err := mysql.DialMysqlGorm("mysql-service", opts...)

// gRPC 服务发现
conn, err := grpc.DialGrpc("user-service")

// 带标签的服务发现
conn, err := grpc.DialGrpcWithTag("user-service", "production")

// Redis 服务发现
pool := redis.DialRedisPool("redis-service", 0, 10)
```

如果服务发现失败，会回退到直接使用提供的地址作为连接地址。

## 错误处理和重试

```go
import "time"

func connectWithRetry() *gorm.DB {
    var db *gorm.DB
    var err error
    
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        db, err = mysql.DialMysqlGorm("mysql-service",
            dialer.WithAuth("root", "password"),
            dialer.WithDBName("myapp"),
        )
        if err == nil {
            break
        }
        
        fmt.Printf("连接失败，第 %d 次重试: %v\n", i+1, err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    if err != nil {
        panic(fmt.Sprintf("经过 %d 次重试后仍无法连接数据库: %v", maxRetries, err))
    }
    
    return db
}
```

## 连接健康检查

```go
func healthCheck(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    // 检查连接是否正常
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("数据库连接异常: %w", err)
    }
    
    // 检查连接池状态
    stats := sqlDB.Stats()
    fmt.Printf("连接池状态: 打开连接数=%d, 使用中连接数=%d, 空闲连接数=%d\n",
        stats.OpenConnections, stats.InUse, stats.Idle)
    
    return nil
}
```

## 配置选项参考

### DialOption 配置

```go
type DialOption struct {
    User        string              // 用户名
    Password    string              // 密码
    DBName      string              // 数据库名
    SqliteCache bool                // SQLite缓存模式
    Logger      logger.Interface    // 日志记录器
}
```

### 选项函数

- `WithAuth(user, pwd string)`: 设置认证信息
- `WithDBName(db string)`: 设置数据库名
- `WithSqliteCache()`: 启用SQLite缓存模式
- `WithLogger(log logger.Interface)`: 设置自定义日志记录器

### 工具函数

- `PackDialOption(opts ...OptionFunc) *DialOption`: 组合配置选项
- `ConfigDB(sqlDB *sql.DB)`: 配置数据库连接池
- `DefaultGormConfig(opt *DialOption) *gorm.Config`: 获取默认GORM配置

## API 参考

### MySQL 连接器

- `mysql.DialMysqlGorm(service string, opts ...dialer.OptionFunc) (*gorm.DB, error)`: 使用服务发现连接 MySQL
- `mysql.DialMysqlGormWithDSN(dsn string, opts ...dialer.OptionFunc) (*gorm.DB, error)`: 使用 DSN 连接 MySQL
- `mysql.DialMysql(service string, opts ...dialer.OptionFunc) (*sql.DB, error)`: 原生 SQL 连接
- `mysql.DialMysqlX(service string, opts ...dialer.OptionFunc) (*sqlx.DB, error)`: sqlx 连接

### SQLite 连接器

- `sqlite.DialSqlLiteGorm(dbFile string, opts ...dialer.OptionFunc) (*gorm.DB, error)`: 连接 SQLite

### Redis 连接器

- `redis.DialRedisPool(addr string, db int, maxIdle int, password ...string) *redis.Pool`: 创建 Redis 连接池

### gRPC 连接器

- `grpc.DialGrpc(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error)`: 基本连接
- `grpc.DialGrpcWithTimeOut(timeout time.Duration, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error)`: 带超时连接
- `grpc.DialGrpcWithTag(service string, tag string, opts ...grpc.DialOption) (*grpc.ClientConn, error)`: 带标签连接
- `grpc.DialGrpcWithContext(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error)`: 带上下文连接

## 最佳实践

1. **连接池配置**：根据应用负载合理设置连接池参数
2. **错误处理**：实现连接重试和降级机制
3. **健康检查**：定期检查数据库连接状态
4. **资源管理**：及时关闭不再使用的连接
5. **日志记录**：配置合适的日志级别用于调试
6. **安全性**：避免在代码中硬编码数据库密码
7. **服务发现**：在微服务环境中优先使用服务名而非硬编码地址

## 注意事项

- 确保数据库服务可访问
- 正确配置防火墙和网络设置
- 注意连接池大小与应用负载的匹配
- 在生产环境中使用连接池
- 定期监控连接池状态和性能指标
- 服务发现依赖 `discover` 组件，确保其正确配置
- SQLite 缓存模式使用内存数据库，重启后数据会丢失 