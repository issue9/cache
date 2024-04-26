// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package memcache 适配 memcached 的实现
package memcache

import (
	"context"
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/issue9/cache"
	"github.com/issue9/cache/caches"
)

type memcacheDriver struct {
	client *memcache.Client
}

// New 声明基于 [memcached] 的缓存系统
//
// [cache.Driver.Driver] 的返回类型为 [memcache.Client]。
//
// [memcached]: https://memcached.org/
func New(addr ...string) cache.Driver {
	return &memcacheDriver{client: memcache.New(addr...)}
}

func (d *memcacheDriver) Get(key string, val any) error {
	item, err := d.client.Get(key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return cache.ErrCacheMiss()
	} else if err != nil {
		return err
	}

	return caches.Unmarshal(item.Value, val)
}

func (d *memcacheDriver) Set(key string, val any, ttl time.Duration) error {
	bs, err := caches.Marshal(val)
	if err != nil {
		return err
	}

	return d.client.Set(&memcache.Item{
		Key:        key,
		Value:      bs,
		Expiration: int32(ttl.Seconds()),
	})
}

func (d *memcacheDriver) Delete(key string) error {
	if err := d.client.Delete(key); !errors.Is(err, memcache.ErrCacheMiss) {
		return err
	}
	return nil
}

func (d *memcacheDriver) Exists(key string) bool {
	_, err := d.client.Get(key)
	return err == nil || !errors.Is(err, memcache.ErrCacheMiss)
}

func (d *memcacheDriver) Clean() error { return d.client.DeleteAll() }

func (d *memcacheDriver) Close() error { return d.client.Close() }

func (d *memcacheDriver) Driver() any { return d.client }

func (d *memcacheDriver) Ping(context.Context) error { return d.client.Ping() }

func (d *memcacheDriver) Counter(key string, ttl time.Duration) (n uint64, f cache.SetCounterFunc, err error) {
	t := int32(ttl.Seconds())

	if n, err = cache.Get[uint64](d, key); errors.Is(err, cache.ErrCacheMiss()) {
		err = d.Set(key, 0, ttl)
		n = 0
	}
	if err != nil {
		return 0, nil, err
	}

	return n, func(n int) (uint64, error) {
		switch {
		default: // n == 0
			return cache.Get[uint64](d, key)
		case n > 0:
			v, err := d.client.Increment(key, uint64(n))
			if err == nil && t > 0 {
				err = d.client.Touch(key, t)
			}

			if errors.Is(err, memcache.ErrCacheMiss) {
				return 0, cache.ErrCacheMiss()
			}
			return v, err
		case n < 0:
			nn := uint64(-n)
			v, err := d.client.Decrement(key, nn)
			if err == nil && t > 0 {
				err = d.client.Touch(key, t)
			}

			if errors.Is(err, memcache.ErrCacheMiss) {
				return 0, cache.ErrCacheMiss()
			}
			return v, err
		}
	}, nil
}
