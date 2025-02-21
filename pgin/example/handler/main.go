// File:		main.go
// Created by:	Hoven
// Created on:	2024-09-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
)

type CustomError struct {
	error
	code int
}

func (ce *CustomError) Error() string {
	return ce.error.Error()
}

func (ce *CustomError) Code() int {
	return ce.code
}

func (ce *CustomError) Cause() error {
	return ce.error
}

func (ce *CustomError) Message() string {
	return ce.error.Error()
}

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

	engine.GET("/hello/error", pgin.RequestResponseHandler(func(c *gin.Context, req *HelloRequest) (resp *HelloResponse, err error) {
		if req.Name == "error" {
			return nil, &CustomError{
				error: errors.New("custom error"),
				code:  401,
			}
		}

		return &HelloResponse{
			Message: "Hello, " + req.Name,
		}, nil
	}))

	engine.GET("/mount/test", pgin.MountHandler[MountTestHandler]())

	engine.Run(":8080")
}

type MountTestHandler struct {
	Name string `form:"name"`
}

type MountTestResponse struct {
	Message string `json:"message"`
}

func (m MountTestHandler) Handle(c *gin.Context) (resp *MountTestResponse, err error) {
	return &MountTestResponse{
		Message: "Hello, " + m.Name,
	}, nil
}
