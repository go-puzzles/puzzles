// File:		log.go
// Created by:	Hoven
// Created on:	2025-01-08
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package grpc

import (
	"context"

	"github.com/go-puzzles/puzzles/plog"
	"google.golang.org/grpc"
)

func unaryClientLoggerInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !plog.IsDebug() {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		td := plog.TimeFuncDuration()
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			plog.Debugc(ctx, "Failed to invoke method %s invoke_time=%s invoke_err=%s", method, td(), err)
		} else {
			plog.Debugc(ctx, "Succeed to invoke method %s invoke_time=%s", method, td())
		}
		return err
	}
}

func streamClientLoggerInterceptor() func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if !plog.IsDebug() {
			return streamer(ctx, desc, cc, method, opts...)
		}

		td := plog.TimeFuncDuration()
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			plog.Debugc(ctx, "Failed to invoke method %s invoke_time=%s invoke_err=%s", method, td(), err)
		} else {
			plog.Debugc(ctx, "Succeed to invoke method %s invoke_time=%s", method, td())
		}
		return cs, err
	}
}
