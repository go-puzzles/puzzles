// File:		engine.go
// Created by:	Hoven
// Created on:	2024-09-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/snail"
)

func init() {
	snail.RegisterObject("ginModeSet", func() error {
		if plog.IsDebug() {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		return nil
	})
}

func Default(opts ...gin.OptionFunc) *gin.Engine {
	engine := gin.New()
	engine.Use(LoggerMiddleware(), gin.Recovery())
	return engine.With(opts...)
}

func NewServerHandler() *gin.Engine {
	engine := gin.New()

	// default health check
	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return engine
}

func NewStandardServerHandler() *gin.Engine {
	engine := NewServerHandler()
	engine.Use(LoggerMiddleware(), gin.Recovery())
	return engine
}

func NewServerHandlerWithOptions(opts ...Option) *gin.Engine {
	engine := NewStandardServerHandler()
	for _, opt := range opts {
		opt(engine)
	}
	return engine
}

type Option func(*gin.Engine)

type Router interface {
	Init(gin.IRouter)
}

func WithRouters(basePath string, routers ...Router) Option {
	return func(e *gin.Engine) {
		group := e.Group(basePath)
		for _, router := range routers {
			router.Init(group)
		}
	}
}

func WithMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return func(e *gin.Engine) {
		for _, m := range middlewares {
			e.Use(m)
		}
	}
}

func WithLoggingRequest(header bool) Option {
	return func(e *gin.Engine) {
		e.Use(LoggingRequest(header))
	}
}

func WithReuseBody() Option {
	return func(e *gin.Engine) {
		e.Use(ReuseBody())
	}
}

func WithServiceName(name string) Option {
	return func(engine *gin.Engine) {
		engine.Use(func(c *gin.Context) {
			c.Set("service", name)
		})
	}
}

func WithHiddenRoutesLog() Option {
	return func(engine *gin.Engine) {
		gin.DefaultWriter = io.Discard
	}
}
