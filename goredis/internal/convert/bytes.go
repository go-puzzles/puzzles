// File:		bytes.go
// Created by:	Hoven
// Created on:	2024-12-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package convert

import (
	"unsafe"
)

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
