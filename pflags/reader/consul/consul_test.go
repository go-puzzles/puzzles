// File:		consul_test.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPossibleTags(t *testing.T) {
	t.Run("listPossibleTags-vx.x.x", func(t *testing.T) {
		tag := "v1.2.4"
		possibleTags := listPossibleTags(tag)

		assert.ElementsMatch(t, possibleTags, []string{"v1.2.4", "v1.2.3", "v1.2.2", "v1.2.1", "v1.2.0", "v1.2"})
	})

	t.Run("listPossibleTags-vx.x.0", func(t *testing.T) {
		tag := "v1.2.0"
		possibleTags := listPossibleTags(tag)

		assert.ElementsMatch(t, possibleTags, []string{"v1.2.0", "v1.2"})
	})

	t.Run("listPossibleTags-vx.x", func(t *testing.T) {
		tag := "v1.2.0"
		possibleTags := listPossibleTags(tag)

		assert.ElementsMatch(t, possibleTags, []string{"v1.2.0", "v1.2"})
	})
}
