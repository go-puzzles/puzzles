// File:		random.go
// Created by:	Hoven
// Created on:	2024-09-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

import (
	crand "crypto/rand"
	"io"
	"math"
	"time"
	"unsafe"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/rand"
)

const (
	MaximumCapacity = math.MaxInt>>1 + 1
	Numeral         = "0123456789"
	LowwerLetters   = "abcdefghijklmnopqrstuvwxyz"
	UpperLetters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SymbolChars     = "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	AllChars        = Numeral + LowwerLetters + UpperLetters + SymbolChars
)

var rn *rand.Rand

func init() {
	rn = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
}

// RandBool generates a random boolean value (true or false).
func RandBool() bool {
	return rn.Intn(2) == 1
}

// RandBoolSlice generates a random boolean slice of specified length.
func RandBoolSlice(length int) []bool {
	if length <= 0 {
		return []bool{}
	}

	result := make([]bool, length)
	for i := range result {
		result[i] = RandBool()
	}

	return result
}

// RandInt generate random int between [min, max).
func RandInt(min, max int) int {
	if min == max {
		return min
	}

	if max < min {
		min, max = max, min
	}

	if min == 0 && max == math.MaxInt {
		return rn.Int()
	}

	return rn.Intn(max-min) + min
}

// RandIntSlice generates a slice of random integers.
// The generated integers are between min and max (exclusive).
func RandIntSlice(length, min, max int) []int {
	if length <= 0 || min > max {
		return []int{}
	}

	result := make([]int, length)
	for i := range result {
		result[i] = RandInt(min, max)
	}

	return result
}

// RandUniqueIntSlice generate a slice of random int of length that do not repeat.
func RandUniqueIntSlice(length, min, max int) []int {
	if min > max {
		return []int{}
	}
	if length > max-min {
		length = max - min
	}

	nums := make([]int, length)
	used := make(map[int]struct{}, length)
	for i := 0; i < length; {
		r := RandInt(min, max)
		if _, use := used[r]; use {
			continue
		}
		used[r] = struct{}{}
		nums[i] = r
		i++
	}

	return nums
}

func RoundToFloat[T constraints.Float | constraints.Integer](x T, n int) float64 {
	tmp := math.Pow(10.0, float64(n))
	x *= T(tmp)
	r := math.Round(float64(x))
	return r / tmp
}

// RandFloat generate random float64 number between [min, max) with specific precision.
func RandFloat(min, max float64, precision int) float64 {
	if min == max {
		return min
	}

	if max < min {
		min, max = max, min
	}

	n := rn.Float64()*(max-min) + min

	return RoundToFloat(n, precision)
}

// RandFloats generate a slice of random float64 numbers of length that do not repeat.
func RandFloats(length int, min, max float64, precision int) []float64 {
	nums := make([]float64, length)
	used := make(map[float64]struct{}, length)
	for i := 0; i < length; {
		r := RandFloat(min, max, precision)
		if _, use := used[r]; use {
			continue
		}
		used[r] = struct{}{}
		nums[i] = r
		i++
	}

	return nums
}

// RandBytes generate random byte slice.
func RandBytes(length int) []byte {
	if length < 1 {
		return []byte{}
	}
	b := make([]byte, length)

	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return nil
	}

	return b
}

// RandString generate random alphabeta string of specified length.
func RandString(length int) string {
	return random(Letters, length)
}

// RandString generate a slice of random string of length strLen based on charset.
// chartset should be one of the following: random.Numeral, random.LowwerLetters, random.UpperLetters
// random.Letters, random.SymbolChars, random.AllChars. or a combination of them.
func RandStringSlice(charset string, sliceLen, strLen int) []string {
	if sliceLen <= 0 || strLen <= 0 {
		return []string{}
	}

	result := make([]string, sliceLen)

	for i := range result {
		result[i] = random(charset, strLen)
	}

	return result
}

// RandFromGivenSlice generate a random element from given slice.
func RandFromGivenSlice[T any](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	return slice[rn.Intn(len(slice))]
}

// RandSliceFromGivenSlice generate a random slice of length num from given slice.
//   - If repeatable is true, the generated slice may contain duplicate elements.
func RandSliceFromGivenSlice[T any](slice []T, num int, repeatable bool) []T {
	if num <= 0 || len(slice) == 0 {
		return slice
	}

	if !repeatable && num > len(slice) {
		num = len(slice)
	}

	result := make([]T, num)
	if repeatable {
		for i := range result {
			result[i] = slice[rn.Intn(len(slice))]
		}
	} else {
		shuffled := make([]T, len(slice))
		copy(shuffled, slice)
		rn.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		result = shuffled[:num]
	}
	return result
}

// RandUpper generate a random upper case string of specified length.
func RandUpper(length int) string {
	return random(UpperLetters, length)
}

// RandLower generate a random lower case string of specified length.
func RandLower(length int) string {
	return random(LowwerLetters, length)
}

// RandNumeral generate a random numeral string of specified length.
func RandNumeral(length int) string {
	return random(Numeral, length)
}

// RandNumeralOrLetter generate a random numeral or alpha string of specified length.
func RandNumeralOrLetter(length int) string {
	return random(Numeral+Letters, length)
}

// RandSymbolChar generate a random symbol char of specified length.
// symbol chars: !@#$%^&*()_+-=[]{}|;':\",./<>?.
func RandSymbolChar(length int) string {
	return random(SymbolChars, length)
}

// nearestPowerOfTwo 返回一个大于等于cap的最近的2的整数次幂，参考java8的hashmap的tableSizeFor函数
func nearestPowerOfTwo(cap int) int {
	n := cap - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return 1
	} else if n >= MaximumCapacity {
		return MaximumCapacity
	}
	return n + 1
}

// random generate a random string based on given string range.
func random(s string, length int) string {
	bytes := make([]byte, length)
	strLength := len(s)
	if strLength <= 0 {
		return ""
	} else if strLength == 1 {
		for i := 0; i < length; i++ {
			bytes[i] = s[0]
		}
		return *(*string)(unsafe.Pointer(&bytes))
	}
	// how many bits are needed to represent the character 's'
	letterIdBits := int(math.Log2(float64(nearestPowerOfTwo(strLength))))
	var letterIdMask int64 = 1<<letterIdBits - 1

	letterIdMax := 63 / letterIdBits
	for i, cache, remain := length-1, rn.Int63(), letterIdMax; i >= 0; {
		// check if the random number generator has exhausted all random numbers
		if remain == 0 {
			cache, remain = rn.Int63(), letterIdMax
		}

		// if s is not an integer multiple of 2, idx may exceed the length of s
		if idx := int(cache & letterIdMask); idx < strLength {
			bytes[i] = s[idx]
			i--
		}

		// go to the next group of letterIdBits bits to perform string retrieval operation
		cache >>= letterIdBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&bytes))
}
