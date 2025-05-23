// File:		scan.go
// Created by:	Hoven
// Created on:	2024-12-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package goredis

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/go-puzzles/puzzles/goredis/internal/convert"
)

func makeSliceNextElemFunc(v reflect.Value) func() reflect.Value {
	elemType := v.Type().Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		return func() reflect.Value {
			if v.Len() < v.Cap() {
				v.Set(v.Slice(0, v.Len()+1))
				elem := v.Index(v.Len() - 1)
				if elem.IsNil() {
					elem.Set(reflect.New(elemType))
				}
				return elem.Elem()
			}

			elem := reflect.New(elemType)
			v.Set(reflect.Append(v, elem))
			return elem.Elem()
		}
	}

	zero := reflect.Zero(elemType)
	return func() reflect.Value {
		if v.Len() < v.Cap() {
			v.Set(v.Slice(0, v.Len()+1))
			return v.Index(v.Len() - 1)
		}

		v.Set(reflect.Append(v, zero))
		return v.Index(v.Len() - 1)
	}
}

func scanRedisSlice(data []string, v reflect.Value) error {
	next := makeSliceNextElemFunc(v)
	for i, s := range data {
		elem := next()
		if err := scan([]byte(s), elem.Addr().Interface()); err != nil {
			err = fmt.Errorf("redis: ScanSlice index=%d value=%q failed: %w", i, s, err)
			return err
		}
	}

	return nil
}

func scan(b []byte, v any) (err error) {
	switch v := v.(type) {
	case nil:
		return fmt.Errorf("goredis: Scan(nil)")
	case *string:
		*v = convert.BytesToString(b)
		return nil
	case *[]byte:
		*v = b
		return nil
	case *int:
		*v, err = convert.Atoi(b)
		return err
	case *int8:
		n, err := convert.ParseInt(b, 10, 8)
		if err != nil {
			return err
		}
		*v = int8(n)
		return nil
	case *int16:
		n, err := convert.ParseInt(b, 10, 16)
		if err != nil {
			return err
		}
		*v = int16(n)
		return nil
	case *int32:
		n, err := convert.ParseInt(b, 10, 32)
		if err != nil {
			return err
		}
		*v = int32(n)
		return nil
	case *int64:
		n, err := convert.ParseInt(b, 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *uint:
		n, err := convert.ParseUint(b, 10, 64)
		if err != nil {
			return err
		}
		*v = uint(n)
		return nil
	case *uint8:
		n, err := convert.ParseUint(b, 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(n)
		return nil
	case *uint16:
		n, err := convert.ParseUint(b, 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(n)
		return nil
	case *uint32:
		n, err := convert.ParseUint(b, 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(n)
		return nil
	case *uint64:
		n, err := convert.ParseUint(b, 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *float32:
		n, err := convert.ParseFloat(b, 32)
		if err != nil {
			return err
		}
		*v = float32(n)
		return nil
	case *float64:
		*v, err = convert.ParseFloat(b, 64)
		if err != nil {
			return err
		}
		return nil
	case *bool:
		*v, err = convert.ParseBool(b)
		if err != nil {
			return err
		}
		return nil
	case *time.Time:
		*v, err = time.Parse(time.RFC3339Nano, convert.BytesToString(b))
		if err != nil {
			return err
		}
		return nil
	case *time.Duration:
		n, err := convert.ParseInt(b, 10, 64)
		if err != nil {
			return err
		}
		*v = time.Duration(n)
		return nil
	case encoding.BinaryUnmarshaler:
		return v.UnmarshalBinary(b)
	case *net.IP:
		*v = b
		return nil
	default:
		return json.Unmarshal(b, v)
	}
}
