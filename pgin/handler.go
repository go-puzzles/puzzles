// File:		handler.go
// Created by:	Hoven
// Created on:	2024-09-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import (
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-puzzles/puzzles/perror"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
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

	err = c.ShouldBind(obj)
	if err != nil && !errors.Is(err, io.EOF) {
		return errors.Wrap(err, "bindJson")
	}

	return nil
}

func ValidateRequestParams(obj any) (err error) {
	return binding.Validator.ValidateStruct(obj)
}

type requestHandler[Q any] func(c *gin.Context, req *Q)

func bindAndValidate[Q any](c *gin.Context) (*Q, error) {
	requestPtr := new(Q)

	if err := ParseRequestParams(c, requestPtr); err != nil {
		plog.Errorc(c, "parse request params failed: %v", err)
		return nil, err
	}

	if err := ValidateRequestParams(requestPtr); err != nil {
		return nil, err
	}

	return requestPtr, nil
}

func isValidHTTPStatusCode(code int) bool {
	return code >= 100 && code < 600 && http.StatusText(code) != ""
}

func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	parseError(c, err)
}

func RequestHandler[Q any](fn requestHandler[Q]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := bindAndValidate[Q](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}

		fn(c, req)
	}
}

type requestResponseHandler[Q any, P any] func(c *gin.Context, req *Q) (resp *P, err error)

func RequestResponseHandler[Q any, P any](fn requestResponseHandler[Q, P]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := bindAndValidate[Q](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorRet(http.StatusBadRequest, err.Error()))
			return
		}

		resp, err := fn(c, req)
		if err != nil {
			handleError(c, err)
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
		httpCode = http.StatusBadRequest
		errCode  int
		message  string
	)

	switch e := err.(type) {
	case perror.ErrorR:
		errCode = e.Code()
		message = e.String()
	default:
		status := c.Writer.Status()
		if status != http.StatusOK {
			errCode = status
		}
		message = err.Error()
	}

	if isValidHTTPStatusCode(errCode) {
		httpCode = errCode
	}

	c.JSON(httpCode, ErrorRet(errCode, message))
	plog.Errorf("handle request: %s error: %v", c.Request.URL.Path, err)
}

type requestWithErrorHandler[Q any] func(c *gin.Context, req *Q) (err error)

func RequestWithErrorHandler[Q any](fn requestWithErrorHandler[Q]) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPtr := new(Q)
		var err error

		if err = ParseRequestParams(c, requestPtr); err != nil {
			plog.Errorc(c, "parse request params failed: %v", err)
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

type ModelHandler[R any] interface {
	Handle(c *gin.Context) (resp *R, err error)
}

func MountHandler[MH ModelHandler[R], R any]() gin.HandlerFunc {
	// r is a pointer of ModelHandler like *ModelHandler
	r := new(MH)
	to := reflect.TypeOf(r)

	// if depth == 1, it means MountHandler[MountTestHandler]()
	// if depth == 2, it means MountHandler[*MountTestHandler]()
	depth := 0

	for to.Kind() == reflect.Pointer {
		to = to.Elem()
		depth++
	}

	if depth != 1 {
		panic("MountHandler[]() Generic types should not be pointer types")
	}

	return RequestResponseHandler(func(c *gin.Context, req *MH) (resp *R, err error) {
		resp, err = (*req).Handle(c)
		return
	})
}
