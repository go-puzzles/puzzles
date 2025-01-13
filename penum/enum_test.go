// File:		enum_test.go
// Created by:	Hoven
// Created on:	2025-01-13
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package penum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumIntDefault(t *testing.T) {
	enumTestInt := New[struct {
		ZERO int
		ONE  int
		TWO  int
	}]()

	assert.Equal(t, 0, enumTestInt.ZERO)
	assert.Equal(t, 1, enumTestInt.ONE)
	assert.Equal(t, 2, enumTestInt.TWO)
}

func TestEnumIntUserDefine(t *testing.T) {
	enumTestInt := New[struct {
		ZERO int
		ONE  int
		TWO  int
	}](func(e *struct {
		ZERO int
		ONE  int
		TWO  int
	}) {
		e.ZERO = 100
		e.ONE = 101
		e.TWO = 102
	})

	assert.Equal(t, 100, enumTestInt.ZERO)
	assert.Equal(t, 101, enumTestInt.ONE)
	assert.Equal(t, 102, enumTestInt.TWO)
}

func TestEnumStringDefault(t *testing.T) {
	enumTestInt := New[struct {
		ZERO string
		ONE  string
		TWO  string
	}]()

	assert.Equal(t, "ZERO", enumTestInt.ZERO)
	assert.Equal(t, "ONE", enumTestInt.ONE)
	assert.Equal(t, "TWO", enumTestInt.TWO)
}
