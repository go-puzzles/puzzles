# poss

poss 是对象存储服务 (OSS) 操作接口包，提供统一的文件上传、下载、管理等功能。目前支持 MinIO 对象存储。

## 功能特性

- 统一的对象存储接口
- 支持文件上传和下载
- 支持文件存在性检查
- 支持预签名URL生成
- 支持代理访问
- 扩展性好，易于添加其他存储后端

## 接口定义

```go
type IOSS interface {
    // 上传文件
    UploadFile(ctx context.Context, size int64, dir, objName string, obj io.Reader, tags map[string]string) (uri string, err error)
    
    // 下载文件
    GetFile(ctx context.Context, objName string, w io.Writer) error
    
    // 检查文件是否存在
    CheckFileExists(ctx context.Context, objName string) (bool, error)
    
    // 生成预签名获取对象URL
    PresignedGetObject(ctx context.Context, objName string, expires time.Duration) (*url.URL, error)
    
    // 代理预签名获取对象
    ProxyPresignedGetObject(objName string, rw http.ResponseWriter, req *http.Request)
}
```

## MinIO 实现

当前提供了 MinIO 的实现，位于 `minio/` 子包中。

### 基本使用

```go
package main

import (
    "context"
    "strings"
    "github.com/go-puzzles/puzzles/poss"
    "github.com/go-puzzles/puzzles/poss/minio"
)

func main() {
    // 创建 MinIO 客户端
    client, err := minio.NewMinIOClient(minio.Config{
        Endpoint:        "localhost:9000",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        UseSSL:          false,
        BucketName:      "my-bucket",
    })
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // 上传文件
    content := strings.NewReader("Hello, World!")
    uri, err := client.UploadFile(ctx, int64(content.Len()), "test", "hello.txt", content, nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("文件上传成功: %s\n", uri)
    
    // 检查文件是否存在
    exists, err := client.CheckFileExists(ctx, "test/hello.txt")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("文件存在: %t\n", exists)
    
    // 生成预签名URL
    url, err := client.PresignedGetObject(ctx, "test/hello.txt", time.Hour)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("预签名URL: %s\n", url.String())
}
```

## 配置

### MinIO 配置

```go
type Config struct {
    Endpoint        string // MinIO 服务端点
    AccessKeyID     string // 访问密钥ID
    SecretAccessKey string // 访问密钥
    UseSSL          bool   // 是否使用SSL
    BucketName      string // 存储桶名称
}
```

## 扩展其他存储后端

要添加新的存储后端支持，只需实现 `IOSS` 接口：

```go
type MyStorageClient struct {
    // 你的配置字段
}

func (c *MyStorageClient) UploadFile(ctx context.Context, size int64, dir, objName string, obj io.Reader, tags map[string]string) (uri string, err error) {
    // 实现上传逻辑
    return "", nil
}

// 实现其他接口方法...
```

## 注意事项

- 确保存储服务可访问
- 注意文件路径和权限设置
- 预签名URL有时效性限制
- 大文件上传建议使用分片上传（后续版本支持） 