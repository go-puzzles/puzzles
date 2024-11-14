// File:		error.go
// Created by:	Hoven
// Created on:	2024-09-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import "fmt"

const (
	defaultErrorMessage = "An error occurred"
)

type Error interface {
	error

	Code() int
	Cause() error
	Message() string
	String() string
}

type internalError struct {
	error
	code    int
	message string
	cause   error
}

func PackError(code int, vals ...any) Error {
	e := &internalError{
		code: code,
	}

	for _, v := range vals {
		switch t := v.(type) {
		case string:
			e.message = t
		case error:
			e.cause = t
		default:
		}
	}

	if e.message == "" {
		e.message = defaultErrorMessage
	}

	return e
}

func (e *internalError) Error() string {
	if e.cause == nil {
		return e.message
	}

	return fmt.Sprintf("%s (caused by: %s)", e.message, e.cause.Error())
}

func (e *internalError) Code() int {
	return e.code
}

func (e *internalError) Cause() error {
	return e.cause
}

func (e *internalError) Message() string {
	return e.message
}

func (e *internalError) String() string {
	return e.Error()
}
