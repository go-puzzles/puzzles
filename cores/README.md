# go-puzzles/cores
简单易用、足够轻量、性能好的 Golang 服务、任务的管理、监控核心 

Easy to use, light enough, good performance Golang worker or service manager and monitor core library

## 特性
简单易用、足够轻量，避免过多的外部依赖

目前实现了以下特性：
- 任务管理
- 定时任务
- 守护进程任务
- 优雅终止
- 服务发现
- 服务注册

支持各种外部扩展:
- httpServer
- grpcServer
- gprcuiHandler

## 快速上手

### 安装
```shell
go get github.com/go-puzzles/puzzles/cores
```

### Worker
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-puzzles/puzzles/cores"
)


func main() {
	core := cores.NewPuzzleCore(
		cores.WithWorker(func(ctx context.Context) error {
			t := time.NewTicker(time.Second * 3)
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-t.C:
				}

				fmt.Println("hello world")
			}
		}),
	)

	cores.Run(core)
}
```

### DaemonWorker
`DaemonWorker` 若返回了错误，则整个服务都将停止，
但是 `cores.WithWorker` 则不会

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-puzzles/puzzles/cores"
)


func main() {
	core := cores.NewPuzzleCore(
		cores.WithDaemonWorker(func(ctx context.Context) error {
			t := time.NewTicker(time.Second * 3)
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-t.C:
				}

				fmt.Println("hello world")
			}
		}),
	)

	cores.Run(core)
}
```

## CronWorker
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-puzzles/puzzles/cores"
)


func main() {
	core := cores.NewPuzzleCore(
		cores.WithCronWorker("0 */1 * * *", func(ctx context.Context) error {
		    fmt.Println("hello world")	
		}),
	)

	cores.Run(core)
}
```


### Http服务 
```go
package main

import (
	"net/http"
	
	"github.com/go-puzzles/puzzles/cores"
	httppuzzle "github.com/go-puzzles/puzzles/cores/puzzles/http-puzzle"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/mux"
)

func main() {
	pflags.Parse()
	router := mux.NewRouter()
	
	router.Path("/hello").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	
	core := cores.NewPuzzleCore(
		httppuzzle.WithCoreHttpPuzzle("/api", router),
	)
	
    if err := cores.Start(core, port()); err != nil {
        panic(err)
    }
}
```

### Grpc服务
```go
package main

import (
	"github.com/go-puzzles/puzzles/cores"
	grpcpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/grpc-puzzle"
	grpcuipuzzle "github.com/go-puzzles/puzzles/cores/puzzles/grpcui-puzzle"
	"github.com/go-puzzles/puzzles/example/cores-with-grpc/examplepb"
	srv "github.com/go-puzzles/puzzles/example/cores-with-grpc/service"
	"github.com/go-puzzles/puzzles/example/cores-with-grpc/testpb"
	"google.golang.org/grpc"
)

func main() {
	example := srv.NewExampleService()
	test := srv.NewTestService()

	srv := cores.NewPuzzleCore(	
		grpcpuzzle.WithCoreGrpcPuzzle(func(srv *grpc.Server) {
			examplepb.RegisterExampleHelloServiceServer(srv, example)
			testpb.RegisterExampleHelloServiceServer(srv, test)
		}),
	)

	if err := cores.Start(srv, 0); err != nil {
		panic(err)
	}
}
```

### 开启GRPCUI
```go
srv := cores.NewPuzzleCore(	
    grpcuipuzzle.WithCoreGrpcUI(),
	grpcpuzzle.WithCoreGrpcPuzzle(func(srv *grpc.Server) {
		examplepb.RegisterExampleHelloServiceServer(srv, example)
		testpb.RegisterExampleHelloServiceServer(srv, test)
	}),
)
```

### Consul服务注册
```go
package main

import (
	"github.com/go-puzzles/puzzles/cores"
	consulpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/consul-puzzle"
)

func main() {
	pflags.Parse(
		pflags.WithConsulEnable(),
	)

	core := cores.NewPuzzleCore(
		cores.WithService(pflags.GetServiceName()),
		consulpuzzle.WithConsulRegsiter(),
	)

	cores.Start(core, port())
}
```

### Consul服务发现
```go
package main

import (
	"github.com/go-puzzles/puzzles/cores/discover"
)


func main() {
    // ....
    discover.GetServiceFinder().GetAddress("serviceName")
    discover.GetServiceFinder().GetAddressWithTag("serviceName", "v.0.0")
    // ....
}
```

## go-puzzles其他工具集
[plog日志工具](https://github.com/go-puzzles/puzzles/plog)

[pflags flags工具](https://github.com/go-puzzles/puzzles/pflags)

[pgorm数据库orm工具](https://github.com/go-puzzles/puzzles/pgorm)

更多请详见: [go-puzzles](https://github.com/go-puzzles/puzzles)
