// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

package cache

import "time"

type prefix struct {
	prefix string
	cache  Cache
}

// Prefix 生成一个带有统一前缀名称的缓存访问对象
//
//	c := memory.New(...)
//	p := cache.Prefix(c, "prefix_")
//	p.Get("k1") // 相当于 c.Get("prefix_k1")
func Prefix(a Cache, p string) Cache {
	if pp, ok := a.(*prefix); ok {
		return &prefix{prefix: pp.prefix + p, cache: pp.cache}
	}
	return &prefix{prefix: p, cache: a}
}

func (p *prefix) Get(key string, v any) error { return p.cache.Get(p.prefix+key, v) }

func (p *prefix) Set(key string, val any, seconds time.Duration) error {
	return p.cache.Set(p.prefix+key, val, seconds)
}

func (p *prefix) Delete(key string) error { return p.cache.Delete(p.prefix + key) }

func (p *prefix) Exists(key string) bool { return p.cache.Exists(p.prefix + key) }

func (p *prefix) Touch(key string, ttl time.Duration) error { return p.cache.Touch(p.prefix+key, ttl) }

func (p *prefix) Counter(key string, ttl time.Duration) (uint64, SetCounterFunc, bool, error) {
	return p.cache.Counter(p.prefix+key, ttl)
}
