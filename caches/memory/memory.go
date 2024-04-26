// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package memory 基于内存的实现
package memory

import (
	"strconv"
	"sync"
	"time"

	"github.com/issue9/cache"
	"github.com/issue9/cache/caches"
)

type memoryDriver struct {
	items *sync.Map
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

func (d *memoryDriver) Ping() error { return nil }

func (d *memoryDriver) Touch(key string, ttl time.Duration) error {
	if i, found := d.findItem(key); found {
		i.expire = time.Now().Add(i.dur)
	}
	return nil
}

func (d *memoryDriver) Counter(key string, ttl time.Duration) (n uint64, f cache.SetCounterFunc, err error) {
	i, loaded := d.items.LoadOrStore(key, &item{
		val:    []byte(strconv.FormatUint(0, 10)),
		dur:    ttl,
		expire: time.Now().Add(ttl),
	})
	if loaded {
		if n, err = strconv.ParseUint(string(i.(*item).val), 10, 64); err != nil {
			return 0, nil, err
		}
	} else {
		n = 0
	}

	var locker sync.Mutex

	return n, func(n int) (uint64, error) {
		locker.Lock()
		defer locker.Unlock()

		var num uint64
		if i, found := d.items.Load(key); found {
			num, err = strconv.ParseUint(string(i.(*item).val), 10, 64)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, cache.ErrCacheMiss()
		}

		switch {
		case n == 0:
			return num, nil
		case n > 0:
			num += uint64(n)
		case n < 0:
			n = -n
			if uint64(n) >= num {
				num = 0
			} else {
				num -= uint64(n)
			}
		}

		d.items.Store(key, &item{
			val:    []byte(strconv.FormatUint(num, 10)),
			dur:    ttl,
			expire: time.Now().Add(ttl),
		})

		return num, nil
	}, nil
}
