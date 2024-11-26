// File:		convert.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

func Convert[S, T any](s []S, fn func(S) T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}
