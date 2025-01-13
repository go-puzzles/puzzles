// File:		inject.go
// Created by:	Hoven
// Created on:	2025-01-09
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package grpcpuzzle

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type injectServerStream struct {
	ctx context.Context
	ss  grpc.ServerStream
}

func newInjectServerStream(ctx context.Context, ss grpc.ServerStream) *injectServerStream {
	return &injectServerStream{
		ss:  ss,
		ctx: ctx,
	}
}

func (ss *injectServerStream) SetHeader(md metadata.MD) error {
	return ss.ss.SetHeader(md)
}

func (ss *injectServerStream) SendHeader(md metadata.MD) error {
	return ss.ss.SendHeader(md)
}

func (ss *injectServerStream) SetTrailer(md metadata.MD) {
	ss.ss.SetTrailer(md)
}

func (ss *injectServerStream) Context() context.Context {
	return ss.ctx
}

func (ss *injectServerStream) SendMsg(m interface{}) error {
	return ss.ss.SendMsg(m)
}

func (ss *injectServerStream) RecvMsg(m interface{}) error {
	return ss.ss.RecvMsg(m)
}
