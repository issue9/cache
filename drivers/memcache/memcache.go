// SPDX-FileCopyrightText: 2017-2026 caixw
//
// SPDX-License-Identifier: MIT

// Package memcache 适配 memcached 的实现
package memcache

import (
	"errors"
	"math"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/issue9/cache"
	"github.com/issue9/cache/drivers"
)

type memcacheDriver struct {
	client *memcache.Client
}

// New 声明基于 [memcached] 的缓存系统
//
// NOTE: Swap 和 Pop 当前实例并非是原子操作。
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

	return drivers.Unmarshal(item.Value, val)
}

func (d *memcacheDriver) Set(key string, val any, ttl time.Duration) error {
	expiration := calcTTL(ttl)

	bs, err := drivers.Marshal(val)
	if err != nil {
		return err
	}

	return d.client.Set(&memcache.Item{
		Key:        key,
		Value:      bs,
		Expiration: expiration,
	})
}

func (d *memcacheDriver) Delete(key string) error {
	if err := d.client.Delete(key); !errors.Is(err, memcache.ErrCacheMiss) {
		return err
	}
	return nil
}

func (d *memcacheDriver) Swap(key string, old, new any, ttl time.Duration) (err error) {
	defer func() {
		if errors.Is(err, memcache.ErrCacheMiss) {
			err = cache.ErrCacheMiss()
		}
	}()

	var item *memcache.Item
	if item, err = d.client.Get(key); err != nil {
		return
	}

	oldBytes := item.Value

	if item.Value, err = drivers.Marshal(new); err != nil {
		return
	}
	item.Expiration = calcTTL(ttl)

	if err = d.client.CompareAndSwap(item); err != nil {
		return
	}

	err = drivers.Unmarshal(oldBytes, old) // 赋值给 err，交由 defer 统一处理
	return err
}

func (d *memcacheDriver) Pop(key string, old any) error {
	if err := d.Get(key, old); err != nil {
		return err
	}
	return d.Delete(key)
}

func (d *memcacheDriver) Exists(key string) bool {
	_, err := d.client.Get(key)
	return err == nil || !errors.Is(err, memcache.ErrCacheMiss)
}

func (d *memcacheDriver) Clean() error { return d.client.DeleteAll() }

func (d *memcacheDriver) Close() error { return d.client.Close() }

func (d *memcacheDriver) Ping() error { return d.client.Ping() }

func (d *memcacheDriver) Touch(key string, ttl time.Duration) (err error) {
	if err = d.client.Touch(key, calcTTL(ttl)); errors.Is(err, memcache.ErrCacheMiss) {
		err = nil
	}
	return err
}

func (d *memcacheDriver) Counter(key string, ttl time.Duration) (n uint64, f cache.SetCounterFunc, exist bool, err error) {
	if n, err = cache.Get[uint64](d, key); errors.Is(err, cache.ErrCacheMiss()) {
		err = d.Set(key, 0, ttl)
		n = 0
	} else {
		exist = true
	}
	if err != nil {
		return 0, nil, false, err
	}

	return n, func(n int) (uint64, error) {
		switch {
		default: // n == 0
			return cache.Get[uint64](d, key)
		case n > 0:
			v, err := d.client.Increment(key, uint64(n))
			if err == nil {
				err = d.client.Touch(key, calcTTL(ttl))
			}

			if errors.Is(err, memcache.ErrCacheMiss) {
				return 0, cache.ErrCacheMiss()
			}
			return v, err
		case n < 0:
			v, err := d.client.Decrement(key, uint64(-n))
			if err == nil {
				err = d.client.Touch(key, calcTTL(ttl))
			}

			if errors.Is(err, memcache.ErrCacheMiss) {
				return 0, cache.ErrCacheMiss()
			}
			return v, err
		}
	}, exist, nil
}

const month = 30 * 24 * 60 * 60 * time.Second

func calcTTL(ttl time.Duration) int32 {
	if ttl == 0 {
		return cache.Forever
	}

	if ttl > month {
		expiresAt := time.Now().Add(ttl)
		seconds := expiresAt.Unix()
		if expiresAt.Nanosecond() > 0 {
			seconds++
		}
		if seconds > math.MaxInt32 {
			panic("TTL 转换为秒数时不能大于 math.MaxInt32")
		}
		return int32(seconds)
	}

	seconds := int64(ttl / time.Second)
	if ttl%time.Second != 0 {
		seconds++
	}
	return int32(seconds)
}
