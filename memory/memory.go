// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import (
	"sync"
	"time"

	"github.com/issue9/cache"
)

// Memory 内存类型的缓存
type Memory struct {
	lock   sync.RWMutex
	items  map[string]*item
	size   int
	ticker *time.Ticker
	done   chan struct{}
}

type item struct {
	val    interface{}
	expire time.Time // 过期的时间
}

// New 声明一个内存缓存。
//
// size 表示初始时的记录数量；
// gcdur 表示执行回收操作的间隔。
func New(size int, gcdur time.Duration) *Memory {
	mem := &Memory{
		items:  make(map[string]*item, size),
		size:   size,
		ticker: time.NewTicker(gcdur),
		done:   make(chan struct{}, 1),
	}

	go func(mem *Memory) {
		for {
			select {
			case <-mem.ticker.C:
				mem.gc()
			case <-mem.done:
				return
			}
		}
	}(mem)

	return mem
}

// Get 获取缓存项。
func (mem *Memory) Get(key string) (interface{}, bool) {
	i, found := mem.findItem(key)
	if !found {
		return nil, false
	}

	return i.val, true
}

// findItem
func (mem *Memory) findItem(key string) (*item, bool) {
	mem.lock.RLock()
	i, found := mem.items[key]
	mem.lock.RUnlock()

	if !found {
		return nil, false
	}

	// 已经过期
	if i.expire.Before(time.Now()) {
		go mem.Delete(key)
		return nil, false
	}

	return i, true
}

// Set 设置或是添加缓存项。
func (mem *Memory) Set(key string, val interface{}, timeout time.Duration) error {
	mem.lock.Lock()
	defer mem.lock.Unlock()

	mem.items[key] = &item{
		val:    val,
		expire: time.Now().Add(timeout),
	}

	return nil
}

// Delete 删除一个缓存项。
func (mem *Memory) Delete(key string) error {
	mem.lock.Lock()
	delete(mem.items, key)
	mem.lock.Unlock()

	return nil
}

// Exists 判断一个缓存项是否存在
func (mem *Memory) Exists(key string) bool {
	mem.lock.RLock()
	_, exists := mem.items[key]
	mem.lock.RUnlock()

	return exists
}

// Incr 增加计数
func (mem *Memory) Incr(key string) error {
	item, found := mem.findItem(key)
	if !found {
		return cache.ErrKeyNotExists
	}

	switch v := item.val.(type) {
	case int:
		item.val = v + 1
	case int64:
		item.val = v + 1
	case int32:
		item.val = v + 1
	case int16:
		item.val = v + 1
	case int8:
		item.val = v + 1
	case uint:
		item.val = v + 1
	case uint64:
		item.val = v + 1
	case uint32:
		item.val = v + 1
	case uint16:
		item.val = v + 1
	case uint8:
		item.val = v + 1
	}

	return nil
}

// Decr 减小计数
func (mem *Memory) Decr(key string) error {
	item, found := mem.findItem(key)
	if !found {
		return cache.ErrKeyNotExists
	}

	switch v := item.val.(type) {
	case int:
		item.val = v - 1
	case int64:
		item.val = v - 1
	case int32:
		item.val = v - 1
	case int16:
		item.val = v - 1
	case int8:
		item.val = v - 1
	case uint:
		if v < 1 {
			return cache.ErrUintNotAllowLessThanZero
		}
		item.val = v - 1
	case uint64:
		if v < 1 {
			return cache.ErrUintNotAllowLessThanZero
		}
		item.val = v - 1
	case uint32:
		if v < 1 {
			return cache.ErrUintNotAllowLessThanZero
		}
		item.val = v - 1
	case uint16:
		if v < 1 {
			return cache.ErrUintNotAllowLessThanZero
		}
		item.val = v - 1
	case uint8:
		if v < 1 {
			return cache.ErrUintNotAllowLessThanZero
		}
		item.val = v - 1
	}

	return nil
}

// Clear 清除所有的缓存内容
func (mem *Memory) Clear() error {
	mem.lock.Lock()
	mem.items = make(map[string]*item, mem.size)
	mem.lock.Unlock()

	return nil
}

// Close 关闭整个缓存系统
func (mem *Memory) Close() error {
	mem.ticker.Stop()
	mem.items = nil
	close(mem.done)

	return nil
}

func (mem *Memory) gc() {
	now := time.Now()

	mem.lock.Lock()
	defer mem.lock.Unlock()

	for key, item := range mem.items {
		if item.expire.Before(now) {
			delete(mem.items, key)
		}
	}
}
