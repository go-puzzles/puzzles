// File:		float.go
// Created by:	Hoven
// Created on:	2024-12-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package convert

import "strconv"

func ParseFloat(b []byte, bitSize int) (float64, error) {
	return strconv.ParseFloat(BytesToString(b), bitSize)
}
