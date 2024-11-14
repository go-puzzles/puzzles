// File:		main.go
// Created by:	Hoven
// Created on:	2024-11-14
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	redisDialer "github.com/go-puzzles/puzzles/dialer/redis"
	"github.com/go-puzzles/puzzles/predis"
)

func main() {
	client := predis.NewRedisClient(redisDialer.DialRedisPool("localhost:6379", 12, 100))
	client.Set("test", "1111")
}
