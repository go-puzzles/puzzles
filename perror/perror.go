// File:		perror.go
// Created by:	Hoven
// Created on:	2025-02-21
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

// Package perror provides a comprehensive error handling solution with error codes,
// error causes tracking, and error wrapping capabilities. It implements a custom
// error interface that extends the standard error interface with additional
// functionality for error codes and error cause chains.
package perror

import "fmt"

const (
	defaultErrorMessage = "An error occurred"
	// Predefined error codes
	CodeUnknown      = 1000 // Unknown error
	CodeInvalidInput = 1001 // Invalid input parameters
	CodeNotFound     = 1002 // Resource not found
	CodeUnauthorized = 1003 // Unauthorized access
)

type ErrorCoder interface {
	error
	Code() int
}

type ErrorCauser interface {
	error
	Cause() error
}

type ErrorR interface {
	ErrorCoder
	ErrorCauser
	fmt.Stringer
}

type internalError struct {
	code    int
	message string
	cause   error
}

func PackError(code int, vals ...any) ErrorR {
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

func (e *internalError) String() string {
	if e.cause == nil {
		return fmt.Sprintf("error code %d: %s", e.code, e.message)
	}
	return fmt.Sprintf("error code %d: %s (caused by: %s)", e.code, e.message, e.cause)
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

// AsErrorR attempts to convert an error to ErrorR interface
func AsErrorR(err error) (ErrorR, bool) {
	if err == nil {
		return nil, false
	}
	if e, ok := err.(ErrorR); ok {
		return e, true
	}
	return nil, false
}

// GetErrorCode returns the error code from an error
// If the error is not ErrorR type, returns CodeUnknown
func GetErrorCode(err error) int {
	if e, ok := AsErrorR(err); ok {
		return e.Code()
	}
	return CodeUnknown
}

// WrapError wraps an existing error with a new message and error code
func WrapError(code int, err error, message string) ErrorR {
	return PackError(code, message, err)
}
