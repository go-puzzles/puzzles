// File:		iter_function.go
// Created by:	Hoven
// Created on:	2024-08-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

import "iter"

func MapIter[inputs ~[]E, E any, U any](arr inputs, fn func(E) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for _, item := range arr {
			if mapped := fn(item); !yield(mapped) {
				return
			}
		}
	}
}

func FilterIter[inputs ~[]E, E any](arr inputs, fn func(E) bool) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, item := range arr {
			if !fn(item) {
				continue
			}

			if !yield(item) {
				return
			}
		}
	}
}

func ZipIter[T1 any, T2 any, E [2]any](arr1 []T1, arrs []T2) iter.Seq[E] {
	return func(yield func(E) bool) {
		length := len(arr1)
		if len(arrs) < length {
			length = len(arrs)
		}

		for i := 0; i < length; i++ {
			zipped := [2]any{arr1[i], arrs[i]}
			if !yield(zipped) {
				return
			}
		}
	}
}
