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

func WithKMPV2() stringSearchOptFunc {
	return func(sso *stringSearchOpt) {
		sso.engine = &KMPSearchEngineV2{}
	}
}

func WithBM() stringSearchOptFunc {
	return func(sso *stringSearchOpt) {
		sso.engine = &BMSearchEngine{}
	}
}

func WithBruteForce() stringSearchOptFunc {
	return func(sso *stringSearchOpt) {
		sso.engine = &BruteForceSearchEngine{}
	}
}

func StringSearch(text, pattern string, opts ...stringSearchOptFunc) int {
	var engine SearchEngine = &BMSearchEngine{}
	if len(text) < 512 {
		engine = &BruteForceSearchEngine{}
	}

	opt := &stringSearchOpt{
		engine: engine,
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

type KMPSearchEngineV2 struct{}

func (k *KMPSearchEngineV2) buildPrefixTable(pattern string) []int {
	if len(pattern) < 2 {
		return []int{-1}
	}
	if len(pattern) < 3 {
		return []int{-1, 0}
	}

	next := make([]int, len(pattern))
	next[0] = -1
	next[1] = 0

	i := 2
	j := 0

	for i < len(pattern) {
		if pattern[j] == pattern[i-1] {
			next[i] = next[i-1] + 1
			j++
			i++
		} else if j > 0 {
			j = next[j]
		} else {
			next[i] = 0
			i++
		}
	}

	return next
}

func (k *KMPSearchEngineV2) Search(text, pattern string) int {
	if len(text) == 0 || len(pattern) == 0 {
		return -1
	}

	next := k.buildPrefixTable(pattern)

	i, j := 0, 0

	for i < len(text) && j < len(pattern) {
		if text[i] == pattern[j] {
			i++
			j++
		} else if j == 0 {
			i++
		} else {
			j = next[j]
		}
	}

	if j == len(pattern) {
		return i - j
	}
	return -1
}

type BMSearchEngine struct{}

func (bm *BMSearchEngine) generateBadCharTable(pattern string) []int {
	const alphabetSize = 256
	badCharTable := make([]int, alphabetSize)
	patternLength := len(pattern)

	for i := range badCharTable {
		badCharTable[i] = patternLength
	}

	for i := 0; i < patternLength-1; i++ {
		badCharTable[pattern[i]] = patternLength - 1 - i
	}

	return badCharTable
}

func (bm *BMSearchEngine) generateGoodSuffixTable(pattern string) []int {
	patternLength := len(pattern)
	goodSuffixTable := make([]int, patternLength)
	suffixes := make([]int, patternLength)

	suffixes[patternLength-1] = patternLength
	g := patternLength - 1
	f := 0
	for i := patternLength - 2; i >= 0; i-- {
		if i > g && suffixes[i+patternLength-1-f] < i-g {
			suffixes[i] = suffixes[i+patternLength-1-f]
		} else {
			if i < g {
				g = i
			}
			f = i
			for g >= 0 && pattern[g] == pattern[g+patternLength-1-f] {
				g--
			}
			suffixes[i] = f - g
		}
	}

	for i := range goodSuffixTable {
		goodSuffixTable[i] = patternLength
	}

	j := 0
	for i := patternLength - 1; i >= 0; i-- {
		if suffixes[i] == i+1 {
			for j < patternLength-1-i {
				if goodSuffixTable[j] == patternLength {
					goodSuffixTable[j] = patternLength - 1 - i
				}
				j++
			}
		}
	}

	for i := 0; i <= patternLength-2; i++ {
		goodSuffixTable[patternLength-1-suffixes[i]] = patternLength - 1 - i
	}

	return goodSuffixTable
}

func (bm *BMSearchEngine) Search(text, pattern string) int {
	badCharTable := bm.generateBadCharTable(pattern)
	goodSuffixTable := bm.generateGoodSuffixTable(pattern)

	textLength := len(text)
	patternLength := len(pattern)

	/*
		i
		BBC ABCDAB ABCDABCDABDE
		ABCDABD
			  j
		bs = 4 - (7-1-6) = 4

		    i
		BBC ABCDAB ABCDABCDABDE
		    ABCDABD
			      j
		bs = 7 - (7-1-6) = 7

		           i
		BBC ABCDAB ABCDABCDABDE
		           ABCDABD
			             j
		bs = 4 - (7-1-6) = 4

		               i
		BBC ABCDAB ABCDABCDABDE
		               ABCDABD
			                 j
	*/

	for i := 0; i <= textLength-patternLength; {
		j := patternLength - 1

		// check if the current character(text[i+j]) in text
		// matches the last character of pattern
		for j >= 0 && pattern[j] == text[i+j] {
			j--
		}

		if j < 0 {
			// match the pattern exactly
			return i
		} else {
			// character not match: calc badCharShift and goodSuffixShift
			// and choice the maximum value to shift
			badCharShift := badCharTable[text[i+j]] - (patternLength - 1 - j)
			goodSuffixShift := goodSuffixTable[j]
			i += max(badCharShift, goodSuffixShift)
		}
	}

	return -1
}

type BruteForceSearchEngine struct{}

func (bf *BruteForceSearchEngine) Search(text, pattern string) int {
	n := len(text)
	m := len(pattern)

	for i := 0; i <= n-m; i++ {
		j := 0
		for j < m && text[i+j] == pattern[j] {
			j++
		}
		if j == m {
			return i
		}
	}
	return -1
}
