// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package memory 基于内存的实现
package memory

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/issue9/cache"
	"github.com/issue9/cache/caches"
)

type memoryDriver struct {
	items *sync.Map
}

type memoryCounter struct {
	driver  *memoryDriver
	key     string
	expires time.Duration
	locker  sync.Mutex
}

type item struct {
	val    []byte
	dur    time.Duration
	expire time.Time // 过期的时间
}

// New 声明一个内存缓存
//
// [cache.Driver.Driver] 的返回类型为 [sync.Map]。
func New() cache.Driver {
	mem := &memoryDriver{
		items: &sync.Map{},
	}
	return mem
}

func (d *memoryDriver) Get(key string, v any) error {
	if item, found := d.findItem(key); found {
		return caches.Unmarshal(item.val, v)
	}
	return cache.ErrCacheMiss()
}

func (d *memoryDriver) findItem(key string) (*item, bool) {
	i, found := d.items.Load(key)
	if !found {
		return nil, false
	}

	ii := i.(*item)
	if ii.dur > 0 && ii.expire.Before(time.Now()) {
		d.items.Delete(key)
		return nil, false
	}

	return ii, true
}

func (d *memoryDriver) Set(key string, val any, ttl time.Duration) error {
	i, found := d.findItem(key)
	if !found {
		bs, err := caches.Marshal(val)
		if err != nil {
			return err
		}

		d.items.Store(key, &item{
			val:    bs,
			dur:    ttl,
			expire: time.Now().Add(ttl),
		})
		return nil
	}

	bs, err := caches.Marshal(val)
	if err == nil {
		i.expire = time.Now().Add(i.dur)
		i.val = bs
	}
	return err
}

func (d *memoryDriver) Delete(key string) error {
	d.items.Delete(key)
	return nil
}

func (d *memoryDriver) Exists(key string) bool {
	_, found := d.findItem(key)
	return found
}

func (d *memoryDriver) Clean() error {
	d.items.Range(func(key, val any) bool {
		d.items.Delete(key)
		return true
	})
	return nil
}

func (d *memoryDriver) Close() error { return d.Clean() }

func (d *memoryDriver) Driver() any { return d.items }

func (d *memoryDriver) Ping(context.Context) error { return nil }

func (d *memoryDriver) Counter(key string, ttl time.Duration) (cache.Counter, error) {
	d.items.LoadOrStore(key, &item{
		val:    []byte(strconv.FormatUint(0, 10)),
		dur:    ttl,
		expire: time.Now().Add(ttl),
	})

	return &memoryCounter{
		driver:  d,
		key:     key,
		expires: ttl,
	}, nil
}

func (c *memoryCounter) Incr(n uint64) (uint64, error) {
	c.locker.Lock()
	defer c.locker.Unlock()

	v, err := c.Value()
	if err != nil {
		return 0, err
	}

	v += n
	c.driver.items.Store(c.key, &item{
		val:    []byte(strconv.FormatUint(v, 10)),
		dur:    c.expires,
		expire: time.Now().Add(c.expires),
	})
	return v, nil
}

func (c *memoryCounter) Decr(n uint64) (uint64, error) {
	c.locker.Lock()
	defer c.locker.Unlock()

	v, err := c.Value()
	if err != nil {
		return 0, err
	}
	if n > v {
		v = 0
	} else {
		v -= n
	}
	c.driver.items.Store(c.key, &item{
		val:    []byte(strconv.FormatUint(v, 10)),
		dur:    c.expires,
		expire: time.Now().Add(c.expires),
	})
	return v, nil
}

func (c *memoryCounter) Value() (uint64, error) {
	if i, found := c.driver.items.Load(c.key); found {
		return strconv.ParseUint(string(i.(*item).val), 10, 64)
	}
	return 0, cache.ErrCacheMiss()
}

func (c *memoryCounter) Delete() error { return c.driver.Delete(c.key) }
