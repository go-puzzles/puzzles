// File:		logger.go
// Created by:	Hoven
// Created on:	2024-09-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/plog/log"
)

func LoggerMiddleware(loggers ...plog.Logger) gin.HandlerFunc {
	var logger plog.Logger
	if len(loggers) == 0 {
		logger = log.New()
	} else {
		logger = loggers[0]
	}
	return func(c *gin.Context) {
		start := time.Now()

		clientIp := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		statusCode := c.Writer.Status()
		spendTime := time.Since(start)

		var logFunc func(ctx context.Context, msg string, v ...any)
		switch {
		case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
			logFunc = logger.Infoc
		case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest:
			logFunc = logger.Warnc
		case statusCode >= http.StatusBadRequest && statusCode <= http.StatusNetworkAuthenticationRequired:
			logFunc = logger.Warnc
		default:
			logFunc = logger.Errorc
		}

		args := []any{
			statusCode,
			spendTime,
			clientIp,
			method,
			path,
		}

		logFunc(c, logMsg, args...)
	}
}

const (
	logMsg = "| %v | %13v | %15s | %4v | %#v"
)
