// File:		error.go
// Created by:	Hoven
// Created on:	2024-09-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package pgin

import (
	"github.com/go-puzzles/puzzles/perror"
)

func PackError(code int, vals ...any) perror.ErrorR {
	return perror.PackError(code, vals...)
}
