// File:		log.go
// Created by:	Hoven
// Created on:	2025-01-08
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package grpcpuzzle

import (
	"context"
	"strings"

	"github.com/go-puzzles/puzzles/plog"
	"google.golang.org/grpc"
)

func unaryServerLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
		return handler(ctx, req)
	}

	prefix := strings.TrimPrefix(info.FullMethod, "/")
	ctx = plog.With(ctx, prefix)

	td := plog.TimeFuncDuration()
	ret, err := handler(ctx, req)
	duration := td()
	if err != nil {
		plog.Debugc(ctx, "Failed to handle method %s handle_time=%s handle_err=%s", info.FullMethod, duration, err)
	} else {
		plog.Debugc(ctx, "Succeed to handle method %s handle_time=%s", info.FullMethod, duration)
	}

	return ret, err
}

func StreamServerLoggerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
		return handler(srv, ss)
	}

	ctx := ss.Context()

	prefix := strings.TrimPrefix(info.FullMethod, "/")
	ctx = plog.With(ctx, prefix)
	ss = newInjectServerStream(ctx, ss)

	td := plog.TimeFuncDuration()
	err := handler(srv, ss)
	duration := td()
	if err != nil {
		plog.Debugc(ctx, "Failed to handle method %s handle_time=%s handle_err=%s", info.FullMethod, duration, err)
	} else {
		plog.Debugc(ctx, "Succeed to handle method %s handle_time=%s", info.FullMethod, duration)
	}

	return err
}
