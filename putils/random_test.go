// File:		random_test.go
// Created by:	Hoven
// Created on:	2024-09-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

import "testing"

func TestRandomStr(t *testing.T) {

	t.Run("random upper", func(t *testing.T) {
		t.Log(RandUpper(8))
	})
}
