// File:		snail.go
// Created by:	Hoven
// Created on:	2024-11-17
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package snail

import (
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
)

type slowerObject struct {
	name string
	fn   func() error
}

var (
	objs = make([]*slowerObject, 0)
)

func RegisterObject(name string, fn func() error) {
	objs = append(objs, &slowerObject{
		name: name,
		fn:   fn,
	})
}

func Init() {
	for _, obj := range objs {
		if err := obj.fn(); err != nil {
			plog.PanicError(errors.Wrapf(err, "slower init obj: %v", obj.name))
		}
		plog.Debugf("slower init obj success! name=%v", obj.name)
	}
}
