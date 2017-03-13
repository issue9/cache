// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memcache

import (
	"errors"
	"log"
	"time"

	gm "github.com/bradfitz/gomemcache/memcache"
)

// Memcache 实现了 memcache 的 cache 接口。
type Memcache struct {
	errlog *log.Logger
	client *gm.Client
}

// New 声明一个新的 Memcache 实例。
func New(errlog *log.Logger, servers ...string) *Memcache {
	return &Memcache{
		errlog: errlog,
		client: gm.New(servers...),
	}
}

// Get 获取缓存项。
//
// memcache 只能返回 []byte，用户得自行转换其类型
func (mem *Memcache) Get(key string) (interface{}, bool) {
	item, err := mem.client.Get(key)
	if err != nil {
		mem.errlog.Println(err)
		return nil, false
	}

	return item.Value, true
}

// Set 设置或是添加缓存项。
//
// memcache 中的 val 只能是 []byte、string、[]rune 三种类型，用户得自行转换其类型
func (mem *Memcache) Set(key string, val interface{}, timeout time.Duration) error {
	var v []byte

	switch vv := val.(type) {
	case string:
		v = []byte(vv)
	case []byte:
		v = vv
	case []rune:
		v = []byte(string(vv))
	default:
		return errors.New("不允许的 val 类型")
	}

	expire := int32(timeout.Seconds())

	return mem.client.Set(&gm.Item{
		Key:        key,
		Value:      v,
		Expiration: expire,
	})
}

// Delete 删除一个缓存项。
func (mem *Memcache) Delete(key string) error {
	return mem.client.Delete(key)
}

// Exists 判断一个缓存项是否存在
func (mem *Memcache) Exists(key string) bool {
	_, found := mem.Get(key)
	return found
}

// Clear 清除所有的缓存内容
func (mem *Memcache) Clear() error {
	return mem.client.DeleteAll()
}

// Close 关闭整个缓存系统
func (mem *Memcache) Close() error {
	mem.Clear()
	mem.client = nil
	return nil
}
