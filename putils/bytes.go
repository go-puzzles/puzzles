// File:		bytes.go
// Created by:	Hoven
// Created on:	2025-02-21
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func Md5(src any) []byte {
	h := md5.New()

	switch val := src.(type) {
	case []byte:
		h.Write(val)
	case string:
		h.Write([]byte(val))
	default:
		h.Write([]byte(fmt.Sprint(src)))
	}

	bs := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(bs)))
	hex.Encode(dst, bs)
	return dst
}

func ShortMd5(src any) []byte {
	return Md5(src)[8:24]
}

func AppendAny(dst []byte, v any) []byte {
	if v == nil {
		return append(dst, "<nil>"...)
	}

	switch val := v.(type) {
	case []byte:
		dst = append(dst, val...)
	case string:
		dst = append(dst, val...)
	case int:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int8:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int16:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int32:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int64:
		dst = strconv.AppendInt(dst, val, 10)
	case uint:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint8:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint16:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint32:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint64:
		dst = strconv.AppendUint(dst, val, 10)
	case float32:
		dst = strconv.AppendFloat(dst, float64(val), 'f', -1, 32)
	case float64:
		dst = strconv.AppendFloat(dst, val, 'f', -1, 64)
	case bool:
		dst = strconv.AppendBool(dst, val)
	case time.Time:
		dst = val.AppendFormat(dst, time.RFC3339)
	case time.Duration:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case error:
		dst = append(dst, val.Error()...)
	case fmt.Stringer:
		dst = append(dst, val.String()...)
	default:
		dst = append(dst, fmt.Sprint(v)...)
	}
	return dst
}
