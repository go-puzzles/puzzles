// File:		handler.go
// Created by:	Hoven
// Created on:	2024-09-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import (
	"errors"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-puzzles/puzzles/plog"
)

var (
	bindLoop = []bindStrategy{
		&headerBind{},
		&urlBind{},
		&queryBind{},
	}
)

func ParseRequestParams(c *gin.Context, obj any) (err error) {
	for _, b := range bindLoop {
		if !b.Need(c) {
			continue
		}
		
		err = b.Bind(c, obj)
		if err == nil {
			continue
		}
		
		switch err.(type) {
		case validator.ValidationErrors:
			err = nil
		case binding.SliceValidationError:
			err = nil
		default:
			return err
		}
	}
	
	return c.Bind(obj)
}

func ValidateRequestParams(obj any) (err error) {
	return binding.Validator.ValidateStruct(obj)
}

type requestHandler[Q any] func(c *gin.Context, req *Q)

func RequestHandler[Q any](fn requestHandler[Q]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		requestPtr := new(Q)
		
		if err = ParseRequestParams(c, requestPtr); err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}
		
		if err = ValidateRequestParams(requestPtr); err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}
		
		fn(c, requestPtr)
	}
}

type requestResponseHandler[Q any, P any] func(c *gin.Context, req *Q) (resp *P, err error)

func RequestResponseHandler[Q any, P any](fn requestResponseHandler[Q, P]) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPtr := new(Q)
		var err error
		
		if err = ParseRequestParams(c, requestPtr); err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}
		
		if err = ValidateRequestParams(requestPtr); err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}
		
		resp, err := fn(c, requestPtr)
		if err != nil {
			parseError(c, err)
			return
		}
		
		c.JSON(http.StatusOK, SuccessRet(resp))
	}
}

type responseHandler[P any] func(c *gin.Context) (resp *P, err error)

func ResponseHandler[P any](fn responseHandler[P]) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := fn(c)
		if err != nil {
			parseError(c, err)
			return
		}
		
		c.JSON(http.StatusOK, SuccessRet(resp))
	}
}

type errorReturnHandler func(c *gin.Context) (err error)

func ErrorReturnHandler(fn errorReturnHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := fn(c)
		if err != nil {
			parseError(c, err)
			return
		}
		
		c.JSON(http.StatusOK, SuccessRet(nil))
	}
}

func parseError(c *gin.Context, err error) {
	var (
		ie       *internalError
		code     int
		respCode int
		message  string
	)
	
	if errors.As(err, &ie) {
		code = ie.Code()
		respCode = ie.Code()
		message = ie.Message()
	} else {
		code = c.Writer.Status()
		respCode = code
		message = err.Error()
	}
	
	if http.StatusText(code) == "" {
		code = http.StatusBadRequest
	}
	
	c.JSON(code, ErrorRet(respCode, message))
	plog.Errorf("handle request: %s error: %v", c.Request.URL.Path, err)
}

type requestWithErrorHandler[Q any] func(c *gin.Context, req *Q) (err error)

func RequestWithErrorHandler[Q any](fn requestWithErrorHandler[Q]) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPtr := new(Q)
		var err error
		
		if err = ParseRequestParams(c, requestPtr); err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}
		
		if err = ValidateRequestParams(requestPtr); err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}
		
		if err := fn(c, requestPtr); err != nil {
			parseError(c, err)
			return
		}
		
		c.JSON(http.StatusOK, SuccessRet(nil))
	}
}
