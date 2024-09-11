// File:		cache.go
// Created by:	Hoven
// Created on:	2024-09-11
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package cache

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Creater func() (any, error)

type Cache interface {
	Get(key string, out any) error
	Set(key string, value any) error
	GetOrCreate(key string, creater Creater, out any) error
	Exists(key string) bool
	Delete(key string)
	Close() error
}

type CacheWithTTL interface {
	Cache
	GetOrCreateWithTTL(key string, creater Creater, out any, ttl time.Duration) error
	SetWithTTL(key string, value any, ttl time.Duration) error
}

type payload struct {
	Content any    `json:"content"`
	Error   string `json:"error,omitempty"`
}

func (p payload) Get(out any) error {
	if p.Error != "" {
		return errors.New(p.Error)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  out,
		TagName: "json",
	})
	if err != nil {
		return errors.Wrap(err, "mapstructure.NewDecoder")
	}

	return decoder.Decode(p.Content)
}

func newPayload(content any, err error) payload {
	if err != nil {
		return payload{Error: err.Error()}
	}
	return payload{Content: content}
}
