// File:		enum.go
// Created by:	Hoven
// Created on:	2025-01-13
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package penum

import "reflect"

type EnumOption[T any] func(*T)

func New[T any](opts ...EnumOption[T]) T {
	var enum T

	enumType := reflect.TypeOf(enum)
	if enumType.Kind() != reflect.Struct {
		panic("enum type must be a struct")
	}

	enumValue := reflect.ValueOf(&enum).Elem()

	for i := 0; i < enumValue.NumField(); i++ {
		field := enumValue.Field(i)
		fieldType := enumType.Field(i)
		if !field.CanSet() {
			continue
		}
		switch k := field.Kind(); k {
		case reflect.String:
			field.SetString(fieldType.Name)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(i))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(uint64(i))
		default:
			panic("enum field must be string or integer")
		}
	}

	for _, opt := range opts {
		opt(&enum)
	}

	return enum
}
