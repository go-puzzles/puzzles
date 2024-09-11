// File:		memory.go
// Created by:	Hoven
// Created on:	2024-09-11
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package cache

import (
	"errors"
	"sync"
)

type MemoryCache struct {
	lock     sync.RWMutex
	payloads map[string]payload
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		payloads: make(map[string]payload),
	}
}

var _ Cache = (*MemoryCache)(nil)

func (mc *MemoryCache) Get(key string, out any) error {
	mc.lock.RLock()
	p, ok := mc.payloads[key]
	mc.lock.RUnlock()
	if ok {
		return p.Get(out)
	}
	return errors.New("Not found")
}

func (mc *MemoryCache) Set(key string, value any) error {
	p := payload{Content: value}
	mc.lock.Lock()
	mc.payloads[key] = p
	mc.lock.Unlock()

	return nil
}

func (mc *MemoryCache) GetOrCreate(key string, creater Creater, out any) error {
	mc.lock.RLock()
	p, ok := mc.payloads[key]
	mc.lock.RUnlock()
	if ok {
		return p.Get(out)
	}

	p = newPayload(creater())
	mc.lock.Lock()
	mc.payloads[key] = p
	mc.lock.Unlock()

	return p.Get(out)
}

func (mc *MemoryCache) Exists(key string) bool {
	mc.lock.RLock()
	_, ok := mc.payloads[key]
	mc.lock.RUnlock()
	return ok
}

func (mc *MemoryCache) Delete(key string) {
	mc.lock.Lock()
	delete(mc.payloads, key)
	mc.lock.Unlock()
}

func (mc *MemoryCache) Close() error {
	return nil
}
