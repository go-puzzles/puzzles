// File:		middleware.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/plog"
)

func ReuseBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := bytes.Buffer{}
		c.Request.Body = io.NopCloser(io.TeeReader(c.Request.Body, &buf))
		c.Next()
		c.Request.Body = io.NopCloser(&buf)
	}
}

const maxBodyLen = 1024

func LoggingRequest(header bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		fields := map[string]interface{}{
			"method": c.Request.Method,
			"uri":    c.Request.URL.RequestURI(),
			"remote": c.Request.RemoteAddr,
			"body":   requestBody(c),
		}
		if header {
			fields["header"] = c.Request.Header
		}
		plog.Infoc(ctx, "incoming http request: %+v", fields)
		c.Next()
	}
}

func requestBody(c *gin.Context) string {
	if c.Request.Body == nil || c.Request.Body == http.NoBody {
		return ""
	}
	bodyData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Sprintf("read request body err: %s", err.Error())
	}
	_ = c.Request.Body.Close()
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyData))

	bodySize := len(bodyData)
	if bodySize > maxBodyLen {
		bodySize = maxBodyLen
	}
	return string(bodyData[:bodySize])
}
