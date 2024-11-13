// File:		string.go
// Created by:	Hoven
// Created on:	2024-11-13
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

type SearchEngine interface {
	Search(text, pattern string) int
}

type stringSearchOpt struct {
	engine SearchEngine
}

type stringSearchOptFunc func(*stringSearchOpt)

func WithEngine(engine SearchEngine) stringSearchOptFunc {
	return func(sso *stringSearchOpt) {
		sso.engine = engine
	}
}

func WithKMP() stringSearchOptFunc {
	return func(sso *stringSearchOpt) {
		sso.engine = &KMPSearchEngine{}
	}
}

func StringSearch(text, pattern string, opts ...stringSearchOptFunc) int {
	opt := &stringSearchOpt{
		engine: &KMPSearchEngine{},
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	return opt.engine.Search(text, pattern)
}

type KMPSearchEngine struct{}

func (k *KMPSearchEngine) buildPrefixTable(pattern string) []int {
	m := len(pattern)
	prefixTable := make([]int, m)
	j := 0

	for i := 1; i < m; i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = prefixTable[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		prefixTable[i] = j
	}

	return prefixTable

}
func (k *KMPSearchEngine) Search(text, pattern string) int {
	next := k.buildPrefixTable(pattern)
	for ti, pi := 0, 0; ti < len(text); ti++ {
		// has match first character
		// when internal characters do not match,
		// use prefixTable to trace back
		for pi > 0 && text[ti] != pattern[pi] {
			pi = next[pi-1]
		}

		// if the first character of pattern is matched，continue further
		if text[ti] == pattern[pi] {
			pi++
		}

		if pi == len(pattern) {
			return ti - pi + 1
		}
	}

	return -1
}
