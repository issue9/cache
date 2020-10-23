// SPDX-License-Identifier: MIT

package memcache

import (
	"log"
	"time"

	gm "github.com/bradfitz/gomemcache/memcache"

	"github.com/issue9/cache"
)

// memcache 实现了 memcache 的 Cache 接口
type memcache struct {
	errlog *log.Logger
	client *gm.Client
}

// New 声明一个新的 Memcache 实例
func New(errlog *log.Logger, client *gm.Client) cache.Cache {
	return &memcache{
		errlog: errlog,
		client: client,
	}
}

func (mem *memcache) Get(key string) (val interface{}, found bool) {
	item, err := mem.client.Get(key)
	if err == gm.ErrCacheMiss {
		return nil, false
	} else if err != nil {
		mem.errlog.Println(err)
		return nil, false
	}

	if err := cache.GoDecode(item.Value, &val); err != nil {
		mem.errlog.Println(err)
		return nil, false
	}

	return val, true
}

func (mem *memcache) Set(key string, val interface{}, timeout time.Duration) error {
	bs, err := cache.GoEncode(&val)
	if err != nil {
		return err
	}

	exp := int32(timeout.Seconds())
	if exp < 1 {
		exp = 1
	}

	return mem.client.Set(&gm.Item{
		Key:        key,
		Value:      bs,
		Expiration: exp,
	})
}

func (mem *memcache) Delete(key string) error {
	return mem.client.Delete(key)
}

func (mem *memcache) Exists(key string) bool {
	_, found := mem.Get(key)
	return found
}

func (mem *memcache) Clear() error {
	return mem.client.DeleteAll()
}

func (mem *memcache) Close() error {
	mem.Clear()
	mem.client = nil
	return nil
}
