// File:		main.go
// Created by:	Hoven
// Created on:	2024-09-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"fmt"

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
	engine := gin.Default()
	engine.Use(pgin.NewSession("__test_session", pgin.InitCookieStore()))

	engine.GET("/hello", pgin.RequestHandler(func(c *gin.Context, req *HelloRequest) {
		session := pgin.GetSession(c)
		session.Set("username", req.Name)
		session.Save()

		c.JSON(200, gin.H{
			"message": "Hello, " + req.Name,
		})
	}))

	engine.GET("/hello/session", func(ctx *gin.Context) {
		session := pgin.GetSession(ctx)
		userName := session.Get("username")
		ctx.JSON(200, gin.H{
			"message": fmt.Sprintf("Your session username is: %v", userName.(string)),
		})
	})

	engine.Run(":8080")
}
