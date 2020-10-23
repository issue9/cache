// SPDX-License-Identifier: MIT

// Package memory 以内存形式存储缓存内容
package memory

import (
	"sync"
	"time"

	"github.com/issue9/cache"
)

type memory struct {
	items  *sync.Map
	ticker *time.Ticker
	done   chan struct{}
}

type item struct {
	val    interface{}
	dur    time.Duration
	expire time.Time // 过期的时间
}

func (i *item) update(val interface{}) {
	i.val = val
	i.expire = time.Now().Add(i.dur)
}

func (i *item) isExpired(now time.Time) bool {
	return i.dur != cache.Forever && i.expire.Before(now)
}

// New 声明一个内存缓存
//
// size 表示初始时的记录数量；
// gc 表示执行回收操作的间隔。
func New(size int, gc time.Duration) cache.Cache {
	mem := &memory{
		items:  &sync.Map{},
		ticker: time.NewTicker(gc),
		done:   make(chan struct{}, 1),
	}

	go func(mem *memory) {
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

func (mem *memory) Get(key string) (interface{}, bool) {
	i, found := mem.findItem(key)
	if !found {
		return nil, false
	}

	return i.val, true
}

func (mem *memory) findItem(key string) (*item, bool) {
	i, found := mem.items.Load(key)
	if !found {
		return nil, false
	}
	return i.(*item), true
}

func (mem *memory) Set(key string, val interface{}, timeout time.Duration) error {
	i, found := mem.findItem(key)
	if !found {
		mem.items.Store(key, &item{
			val:    val,
			dur:    timeout,
			expire: time.Now().Add(timeout),
		})
		return nil
	}

	i.update(val)
	return nil
}

func (mem *memory) Delete(key string) error {
	mem.items.Delete(key)
	return nil
}

func (mem *memory) Exists(key string) bool {
	_, found := mem.items.Load(key)
	return found
}

func (mem *memory) Clear() error {
	mem.items = &sync.Map{}
	return nil
}

func (mem *memory) Close() error {
	mem.ticker.Stop()
	mem.items = nil
	close(mem.done)

	return nil
}

func (mem *memory) gc() {
	now := time.Now()

	mem.items.Range(func(key, val interface{}) bool {
		if v := val.(*item); v.isExpired(now) {
			mem.items.Delete(key)
		}
		return true
	})
}
