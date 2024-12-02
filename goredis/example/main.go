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

	"github.com/go-puzzles/puzzles/predis"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	client := predis.NewRedisClient("localhost:6379", 1)
	err := client.SetValue(context.Background(), "test-person", &Person{"Hoven", 16}, 0)
	if err != nil {
		panic(err)
	}

	err = client.SetValue(context.Background(), "test-str", "hoven", 0)
	if err != nil {
		panic(err)
	}

	err = client.SetValue(context.Background(), "test-int", 18, 0)
	if err != nil {
		panic(err)
	}
}
