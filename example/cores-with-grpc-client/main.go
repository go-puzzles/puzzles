// File:		main.go
// Created by:	Hoven
// Created on:	2025-01-08
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"context"

	"github.com/go-puzzles/puzzles/dialer/grpc"
	"github.com/go-puzzles/puzzles/example/cores-with-grpc/examplepb"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/level"
)

func main() {
	plog.Enable(level.LevelDebug)

	conn, err := grpc.DialGrpc("localhost:21112")
	plog.PanicError(err)

	client := examplepb.NewExampleHelloServiceClient(conn)
	resp, err := client.SayHello(context.Background(), &examplepb.HelloRequest{Name: "super"})
	plog.PanicError(err)

	plog.Infof("%v", resp)
}
