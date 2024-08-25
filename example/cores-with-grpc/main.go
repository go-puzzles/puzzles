package main

import (
	"github.com/go-puzzles/puzzles/cores"
	grpcpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/grpc-puzzle"
	grpcuipuzzle "github.com/go-puzzles/puzzles/cores/puzzles/grpcui-puzzle"
	"github.com/go-puzzles/puzzles/example/cores-with-grpc/examplepb"
	srv "github.com/go-puzzles/puzzles/example/cores-with-grpc/service"
	"github.com/go-puzzles/puzzles/example/cores-with-grpc/testpb"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/level"
	"google.golang.org/grpc"
)

func main() {
	plog.Enable(level.LevelDebug)
	
	example := srv.NewExampleService()
	test := srv.NewTestService()
	
	srv := cores.NewPuzzleCore(
		grpcuipuzzle.WithCoreGrpcUI(),
		grpcpuzzle.WithCoreGrpcPuzzle(func(srv *grpc.Server) {
			examplepb.RegisterExampleHelloServiceServer(srv, example)
			testpb.RegisterExampleHelloServiceServer(srv, test)
		}),
	)
	
	if err := cores.Start(srv, 0); err != nil {
		panic(err)
	}
}
