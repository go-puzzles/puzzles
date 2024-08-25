package service

import (
	"context"
	"fmt"
	
	"github.com/go-puzzles/puzzles/example/cores-with-grpc/examplepb"
	"github.com/go-puzzles/puzzles/example/cores-with-grpc/testpb"
	"google.golang.org/grpc/metadata"
)

type ExampleService struct {
	examplepb.UnimplementedExampleHelloServiceServer
}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (es *ExampleService) SayHello(ctx context.Context, in *examplepb.HelloRequest) (*examplepb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md, ok)
	
	return &examplepb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}

type TestService struct {
	testpb.UnimplementedExampleHelloServiceServer
}

func NewTestService() *TestService {
	return &TestService{}
}

func (es *TestService) SayHello(ctx context.Context, in *testpb.HelloRequest) (*testpb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(md, ok)
	
	return &testpb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}
