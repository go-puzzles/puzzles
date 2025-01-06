// File:		resp.go
// Created by:	Hoven
// Created on:	2024-09-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import "github.com/gin-gonic/gin"

type Ret struct {
	Code    int `json:"code"`
	Data    any `json:"data,omitempty"`
	Message any `json:"message,omitempty"`
}

func SuccessRet(data any) *Ret {
	return &Ret{Code: 200, Data: data, Message: "success"}
}

func ErrorRet(code int, message any) *Ret {
	var msg any
	switch m := message.(type) {
	case string:
		msg = m
	case error:
		msg = m.Error()
	default:
		msg = m
	}

	return &Ret{Code: code, Data: nil, Message: msg}
}

func ReturnSuccess(c *gin.Context, data any) {
	c.JSON(200, SuccessRet(data))
}

func ReturnError(c *gin.Context, code int, message any) {
	c.JSON(code, ErrorRet(code, message))
}
