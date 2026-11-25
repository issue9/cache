// SPDX-FileCopyrightText: 2017-2026 caixw
//
// SPDX-License-Identifier: MIT

// Package memory 基于内存的实现
package memory

import (
	"strconv"
	"sync"
	"time"

	"github.com/issue9/cache"
	"github.com/issue9/cache/drivers"
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
// 虽是基于内存的缓存，但是操作的每一份数据都是由经由
// [drivers.Marshal] 和 [drivers.Unmarshal] 编码的数据副本，
// 不会相互影响。
func New() cache.Driver { return &memoryDriver{items: &sync.Map{}} }

func (d *memoryDriver) Get(key string, v any) error {
	if item, found := d.findItem(key); found {
		return drivers.Unmarshal(item.val, v)
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
	bs, err := drivers.Marshal(val)
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

func (d *memoryDriver) Delete(key string) error {
	d.items.Delete(key)
	return nil
}

func (d *memoryDriver) Swap(key string, old, new any, ttl time.Duration) error {
	bs, err := drivers.Marshal(new)
	if err != nil {
		return err
	}

	i, found := d.items.Swap(key, &item{
		val:    bs,
		dur:    ttl,
		expire: time.Now().Add(ttl),
	})

	if found {
		ii := i.(*item)
		if ii.dur > 0 && ii.expire.Before(time.Now()) {
			d.items.Delete(key) // 如果加载项已经过期，那么应该是不能保存的，所以要删除已保存的项。
			return cache.ErrCacheMiss()
		}
		return drivers.Unmarshal(ii.val, old)
	}

	return cache.ErrCacheMiss()
}

func (d *memoryDriver) Pop(key string, val any) error {
	if i, found := d.items.LoadAndDelete(key); found {
		ii := i.(*item)
		if ii.dur > 0 && ii.expire.Before(time.Now()) {
			d.items.Delete(key) // 如果加载项已经过期，那么应该是不能保存的，所以要删除已保存的项。
			return cache.ErrCacheMiss()
		}
		return drivers.Unmarshal(ii.val, val)
	}
	return cache.ErrCacheMiss()
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

func (d *memoryDriver) Ping() error { return nil }

func (d *memoryDriver) Touch(key string, ttl time.Duration) error {
	if i, found := d.findItem(key); found {
		i.dur = ttl
		i.expire = time.Now().Add(ttl)
		d.items.Store(key, i)
	}
	return nil
}

func (d *memoryDriver) Counter(key string, ttl time.Duration) (n uint64, f cache.SetCounterFunc, exist bool, err error) {
	i, loaded := d.items.LoadOrStore(key, &item{
		val:    []byte(strconv.FormatUint(0, 10)),
		dur:    ttl,
		expire: time.Now().Add(ttl),
	})
	if loaded {
		ii := i.(*item)
		if ii.dur > 0 && ii.expire.Before(time.Now()) {
			d.items.Delete(key)
		} else {
			exist = true
			if n, err = strconv.ParseUint(string(ii.val), 10, 64); err != nil {
				return 0, nil, false, err
			}
		}
	}

	var locker sync.Mutex

	return n, func(n int) (uint64, error) {
		locker.Lock()
		defer locker.Unlock()

		var num uint64
		if i, found := d.items.Load(key); found {
			ii := i.(*item)
			if ii.dur > 0 && ii.expire.Before(time.Now()) {
				d.items.Delete(key)
				return 0, cache.ErrCacheMiss()
			}

			if num, err = strconv.ParseUint(string(ii.val), 10, 64); err != nil {
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
	}, exist, nil
}
