// File:		main.go
// Created by:	Hoven
// Created on:	2024-12-02
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"context"
	"fmt"

	"github.com/go-puzzles/puzzles/goredis"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	client := goredis.NewRedisClient("localhost:6379", 0)
	err := client.LPushValue(context.Background(), "test-push", "this is a string")
	if err != nil {
		panic(err)
	}

	var intRet int
	err = client.LPopValue(context.Background(), "test-push", &intRet)
	if err != nil {
		panic(err)
	}

	fmt.Println(intRet)
}
