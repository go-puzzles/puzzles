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
	"time"

	"github.com/go-puzzles/puzzles/goredis"
)

type Person struct {
	Name string
	Age  int
}

func popValue(client *goredis.PuzzleRedisClient, ret any) {
	err := client.RPopValue(context.Background(), "test-push", ret)
	if err != nil {
		panic(err)
	}
}

func main() {
	client := goredis.NewRedisClient("localhost:6379", 0)
	err := client.LPushValue(
		context.Background(),
		"test-push",
		"this is a string",
		123,
		true,
		123.456,
		time.Now(),
		time.Second*30,
		Person{Name: "John Doe", Age: 30},
	)
	if err != nil {
		panic(err)
	}
	defer client.Del(context.Background(), "test-push")

	var stringRet string
	popValue(client, &stringRet)
	fmt.Printf("stringRet: %v\n", stringRet)

	var intRet int
	popValue(client, &intRet)
	fmt.Printf("intRet: %v\n", intRet)

	var boolRet bool
	popValue(client, &boolRet)
	fmt.Printf("boolRet: %v\n", boolRet)

	var float64Ret float64
	popValue(client, &float64Ret)
	fmt.Printf("float64Ret: %v\n", float64Ret)

	var timeRet time.Time
	popValue(client, &timeRet)
	fmt.Printf("timeRet: %v\n", timeRet)

	var timeDurationRet time.Duration
	popValue(client, &timeDurationRet)
	fmt.Printf("timeDurationRet: %v\n", timeDurationRet)

	var personRet Person
	popValue(client, &personRet)
	fmt.Printf("personRet: %v\n", personRet)
}
