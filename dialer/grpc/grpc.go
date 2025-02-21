package grpc

import (
	"context"
	"time"

	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/plog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialGrpc(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialGrpcWithTimeOut(10*time.Second, service, opts...)
}

func DialGrpcWithUnBlock(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dialGrpcWithTagContext(ctx, service, "", opts...)
}

func DialGrpcWithTag(service string, tag string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialGrpcWithTagTimeout(10*time.Second, service, tag, opts...)
}

func DialGrpcWithTimeOut(timeout time.Duration, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return DialGrpcWithContext(ctx, service, opts...)
}

func DialGrpcWithTagTimeout(timeout time.Duration, service, tag string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return DialGrpcWithTagContext(ctx, service, tag, opts...)
}

func DialGrpcWithContext(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return dialGrpcWithTagContext(ctx, service, "", opts...)
}

func DialGrpcWithTagContext(ctx context.Context, service string, tag string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return dialGrpcWithTagContext(ctx, service, tag, opts...)
}

// func dialGrpcWithTagContextUnblock(ctx context.Context, service string, tag string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
// 	options := append(opts, defaultGRPCDialOptions()...)
// 	options = append(options, grpc.WithBlock())
//
// 	address := discover.GetServiceFinder().GetAddressWithTag(service, tag)
// 	conn, err := grpc.DialContext(ctx, address, options...)
// 	if tag != "" {
// 		plog.Debugc(ctx, "dial grpc service %s with tag %s. Addr=%s", service, tag, address)
// 	} else {
// 		plog.Debugc(ctx, "dial grpc service %s. Addr=%s", service, address)
// 	}
// 	return conn, err
// }

func dialGrpcWithTagContext(ctx context.Context, service, tag string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	options := append(opts, defaultGRPCDialOptions()...)

	address := discover.GetServiceFinder().GetAddressWithTag(service, tag)

	conn, err := grpc.NewClient(
		address,
		options...,
	)

	if tag != "" {
		plog.Debugc(ctx, "dial grpc service %s with tag %s. Addr=%s", service, tag, address)
	} else {
		plog.Debugc(ctx, "dial grpc service %s. Addr=%s", service, address)
	}
	return conn, err
}

func defaultGRPCDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024 * 16)),
		grpc.WithChainUnaryInterceptor(unaryClientLoggerInterceptor()),
		grpc.WithChainStreamInterceptor(streamClientLoggerInterceptor()),
	}
}
