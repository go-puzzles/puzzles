// File:		dedup.go
// Created by:	Hoven
// Created on:	2024-11-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

func Dedup[T comparable](slice []T) []T {
	if len(slice) <= 50 {
		return dupSmall(slice)
	}
	return dupLarge(slice)
}

func dupSmall[T comparable](slice []T) []T {
	idx := 0
	for _, s := range slice {
		var j int
		for j = 0; j < idx; j++ {
			if slice[j] == s {
				break
			}
		}
		if j >= idx {
			slice[idx] = s
			idx++
		}
	}
	return slice[:idx]
}

// dupLarge is the hashmap version of DupStrings with O(n) algorithm.
func dupLarge[T comparable](slice []T) []T {
	m := map[T]struct{}{}
	idx := 0
	for i, s := range slice {
		if _, hit := m[s]; hit {
			continue
		} else {
			m[s] = struct{}{}
			slice[idx] = slice[i]
			idx++
		}
	}
	return slice[:idx]
}
