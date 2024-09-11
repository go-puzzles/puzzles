// File:		memory_ttl.go
// Created by:	Hoven
// Created on:	2024-09-11
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package cache

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var _ CacheWithTTL = (*MemoryCacheWithTTL)(nil)

type MemoryCacheWithTTL struct {
	lock            sync.RWMutex
	payloads        map[string]payloadWithExpire
	compactCancel   func()
	compactInterval time.Duration
}

type payloadWithExpire struct {
	payload  payload
	expireAt time.Time
}

func (p *payloadWithExpire) IsValid() bool {
	if p == nil {
		return false
	}
	if p.expireAt.IsZero() {
		// No expire
		return true
	}
	return time.Now().Before(p.expireAt)
}

func (p *payloadWithExpire) Get(out any) error {
	return p.payload.Get(out)
}

func NewMemoryCacheWithTTL(compactInterval time.Duration) *MemoryCacheWithTTL {
	mc := &MemoryCacheWithTTL{
		payloads: make(map[string]payloadWithExpire),
	}
	ctx, cancel := context.WithCancel(context.Background())
	mc.compactCancel = cancel
	go mc.runCompact(ctx, compactInterval)
	return mc
}

func (mc *MemoryCacheWithTTL) runCompact(ctx context.Context, compactInterval time.Duration) {
	ticker := time.NewTicker(compactInterval)
	defer ticker.Stop()
	for range ticker.C {
		if ctx.Err() != nil {
			return
		}
		mc.lock.Lock()
		for key := range mc.payloads {
			p := mc.payloads[key]
			if p.expireAt.IsZero() {
				// no exipre
				continue
			}
			if time.Now().After(p.expireAt) {
				delete(mc.payloads, key)
			}
		}
		mc.lock.Unlock()
	}
}

func (mc *MemoryCacheWithTTL) Get(key string, out any) error {
	mc.lock.RLock()
	p, ok := mc.payloads[key]
	mc.lock.RUnlock()
	if !ok || !p.IsValid() {
		return errors.New("Not found")
	}
	return p.Get(out)
}

func (mc *MemoryCacheWithTTL) Set(key string, value any) error {
	return mc.SetWithTTL(key, value, 0)
}

func (mc *MemoryCacheWithTTL) SetWithTTL(key string, value any, ttl time.Duration) error {
	p := payload{Content: value}
	mc.lock.Lock()

	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}
	mc.payloads[key] = payloadWithExpire{
		payload:  p,
		expireAt: expireAt,
	}
	mc.lock.Unlock()
	return nil
}

func (mc *MemoryCacheWithTTL) GetOrCreate(key string, creator Creater, out any) error {
	return mc.GetOrCreateWithTTL(key, creator, out, 0)
}

func (mc *MemoryCacheWithTTL) GetOrCreateWithTTL(key string, creator Creater, out any, ttl time.Duration) error {
	mc.lock.RLock()
	wrapPayload, ok := mc.payloads[key]
	mc.lock.RUnlock()
	if ok && wrapPayload.IsValid() {
		return wrapPayload.Get(out)
	}

	p := newPayload(creator())
	mc.lock.Lock()

	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}
	mc.payloads[key] = payloadWithExpire{
		payload:  p,
		expireAt: expireAt,
	}
	mc.lock.Unlock()

	return p.Get(out)
}

func (mc *MemoryCacheWithTTL) Exists(key string) bool {
	mc.lock.RLock()
	_, ok := mc.payloads[key]
	mc.lock.RUnlock()
	return ok
}

func (mc *MemoryCacheWithTTL) Delete(key string) {
	mc.lock.Lock()
	delete(mc.payloads, key)
	mc.lock.Unlock()
}

func (mc *MemoryCacheWithTTL) Close() error {
	mc.compactCancel()
	mc.lock.Lock()
	mc.payloads = nil
	mc.lock.Unlock()
	return nil
}
