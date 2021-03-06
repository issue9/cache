// SPDX-License-Identifier: MIT

// Package memcache memcached 客户端的 cache 接口实现
package memcache

import (
	gm "github.com/bradfitz/gomemcache/memcache"

	"github.com/issue9/cache"
)

// memcache 实现了 memcache 的 Cache 接口
type memcache struct {
	client *gm.Client
}

// NewFromServers 声明一个新的 Memcache 实例
func NewFromServers(addr ...string) cache.Cache {
	return New(gm.New(addr...))
}

// New 声明一个新的 Memcache 实例
func New(client *gm.Client) cache.Cache {
	return &memcache{
		client: client,
	}
}

func (mem *memcache) Get(key string) (val interface{}, err error) {
	item, err := mem.client.Get(key)
	if err == gm.ErrCacheMiss {
		return nil, cache.ErrCacheMiss
	} else if err != nil {
		return nil, err
	}

	if err := cache.GoDecode(item.Value, &val); err != nil {
		return nil, err
	}

	return val, nil
}

func (mem *memcache) Set(key string, val interface{}, seconds int) error {
	bs, err := cache.GoEncode(&val)
	if err != nil {
		return err
	}

	return mem.client.Set(&gm.Item{
		Key:        key,
		Value:      bs,
		Expiration: int32(seconds),
	})
}

func (mem *memcache) Delete(key string) error {
	return mem.client.Delete(key)
}

func (mem *memcache) Exists(key string) bool {
	_, found := mem.Get(key)
	return found != cache.ErrCacheMiss
}

func (mem *memcache) Clear() error {
	return mem.client.DeleteAll()
}

func (mem *memcache) Close() error {
	// NOTE: 关闭服务，不能清除服务器的内容
	mem.client = nil
	return nil
}
