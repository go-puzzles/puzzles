// File:		int.go
// Created by:	Hoven
// Created on:	2024-12-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package convert

import "strconv"

func Atoi(b []byte) (int, error) {
	return strconv.Atoi(BytesToString(b))
}

func ParseInt(b []byte, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(BytesToString(b), base, bitSize)
}

func ParseUint(b []byte, base int, bitSize int) (uint64, error) {
	return strconv.ParseUint(BytesToString(b), base, bitSize)
}
