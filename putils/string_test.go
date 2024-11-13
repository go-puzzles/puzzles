// File:		string_test.go
// Created by:	Hoven
// Created on:	2024-11-13
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/go-puzzles/puzzles/putils"
	"github.com/stretchr/testify/assert"
)

var (
	text    string
	pattern = "The heap target should arguably only include the scannable live heap as a closer"
)

func TestMain(m *testing.M) {
	var err error
	text, err = requestStringData()
	if err != nil {
		fmt.Printf("failed to request data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("text string length: ", len(text))
	fmt.Println("pattern string length: ", len(pattern))

	m.Run()
}

func requestStringData() (string, error) {
	resp, err := http.DefaultClient.Get("https://go.dev/doc/gc-guide")
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func TestStringSearch(t *testing.T) {
	t.Run("testMatch", func(t *testing.T) {
		i := putils.StringSearch(text, pattern, putils.WithEngine(&putils.BMSearchEngine{}))
		if !assert.NotEqual(t, i, -1) {
			t.Error("expected pattern to be found")
			return
		}
		assert.Equal(t, text[i:i+len(pattern)], pattern)
	})

	t.Run("testNotMatch", func(t *testing.T) {
		i := putils.StringSearch(text, "gopuzzles")
		assert.Equal(t, i, -1)
	})
}

func BenchmarkStringSearchWithKMP(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		i := putils.StringSearch(text, pattern, putils.WithEngine(&putils.KMPSearchEngine{}))
		if !assert.Equal(b, text[i:i+len(pattern)], pattern) {
			b.Error("expected pattern to be found")
			break
		}
	}
}

func BenchmarkStringSearchWithKMPV2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		i := putils.StringSearch(text, pattern, putils.WithEngine(&putils.KMPSearchEngineV2{}))
		if !assert.Equal(b, text[i:i+len(pattern)], pattern) {
			b.Error("expected pattern to be found")
			break
		}
	}
}

func BenchmarkStringSearchWithBM(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		i := putils.StringSearch(text, pattern, putils.WithEngine(&putils.BMSearchEngine{}))
		if !assert.Equal(b, text[i:i+len(pattern)], pattern) {
			b.Error("expected pattern to be found")
			break
		}
	}
}

func BenchmarkStringSearchWithBruteForce(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		i := putils.StringSearch(text, pattern, putils.WithEngine(&putils.BruteForceSearchEngine{}))
		if !assert.Equal(b, text[i:i+len(pattern)], pattern) {
			b.Error("expected pattern to be found")
			break
		}
	}
}
