// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package memcache 适配 memcached 的实现
package memcache

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/issue9/cache"
	"github.com/issue9/cache/caches"
)

type memcacheDriver struct {
	client *memcache.Client
}

type memcacheCounter struct {
	driver *memcacheDriver
	key    string
	ttl    int32
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

func (d *memcacheDriver) Counter(key string, ttl time.Duration) (cache.Counter, error) {
	if !d.Exists(key) {
		err := d.client.Set(&memcache.Item{
			Key:        key,
			Value:      []byte(strconv.FormatUint(0, 10)),
			Expiration: int32(ttl.Seconds()),
		})
		if err != nil {
			return nil, err
		}
	}

	return &memcacheCounter{
		driver: d,
		key:    key,
		ttl:    int32(ttl.Seconds()),
	}, nil
}

func (c *memcacheCounter) Incr(n uint64) (uint64, error) {
	v, err := c.driver.client.Increment(c.key, n)
	if err == nil && c.ttl > 0 {
		err = c.driver.client.Touch(c.key, c.ttl)
	}

	if errors.Is(err, memcache.ErrCacheMiss) {
		return 0, cache.ErrCacheMiss()
	}
	return v, err
}

func (c *memcacheCounter) Decr(n uint64) (uint64, error) {
	v, err := c.driver.client.Decrement(c.key, n)
	if err == nil && c.ttl > 0 {
		err = c.driver.client.Touch(c.key, c.ttl)
	}

	if errors.Is(err, memcache.ErrCacheMiss) {
		return 0, cache.ErrCacheMiss()
	}
	return v, err
}

func (c *memcacheCounter) Value() (uint64, error) {
	item, err := c.driver.client.Get(c.key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return 0, cache.ErrCacheMiss()
	} else if err != nil {
		return 0, err
	}

	v := string(item.Value)
	if v == "0 " { // 零值?
		return 0, nil
	}
	return strconv.ParseUint(v, 10, 64)
}

func (c *memcacheCounter) Delete() error {
	err := c.driver.client.Delete(c.key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil
	}
	return err
}
