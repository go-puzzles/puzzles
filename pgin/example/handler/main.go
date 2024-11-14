// File:		main.go
// Created by:	Hoven
// Created on:	2024-09-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
)

type HelloRequest struct {
	Name string `form:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	engine := pgin.Default()

	engine.GET("/hello", pgin.RequestHandler(func(c *gin.Context, req *HelloRequest) {
		c.JSON(200, gin.H{
			"message": "Hello, " + req.Name,
		})
	}))

	engine.GET("/hello/response", pgin.RequestResponseHandler(func(c *gin.Context, req *HelloRequest) (resp *HelloResponse, err error) {
		if req.Name == "error" {
			return nil, pgin.PackError(4444, "test error")
		}
		return &HelloResponse{
			Message: "Hello, " + req.Name,
		}, nil
	}))

	engine.Run(":8080")
}
