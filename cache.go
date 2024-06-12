// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

//go:generate web locale -l=und -m -f=yaml ./
//go:generate web update-locale -src=./locales/und.yaml -dest=./locales/cmn-Hans.yaml

// Package cache 统一的缓存系统接口
package cache

import (
	"errors"
	"time"

	"github.com/issue9/localeutil"
)

var errCacheMiss = localeutil.Error("cache miss")

// Forever 永不过时
const Forever = 0

// Cache 缓存内容的访问接口
type Cache interface {
	// Get 获取缓存项
	//
	// 当前不存在时，返回 [ErrCacheMiss] 错误。
	// key 为缓存项的唯一 ID；
	// v 为缓存写入的地址，应该始终为指针类型；
	Get(key string, v any) error

	// Set 设置或是添加缓存项
	//
	// key 表示保存该数据的唯一 ID；
	// val 表示保存的数据对象，如果是结构体，则会调用 gob 包进行序列化。
	// ttl 表示过了该时间，缓存项将被回收。如果该值为 [Forever]，该值永远不会回收。
	Set(key string, val any, ttl time.Duration) error

	// Delete 删除一个缓存项
	//
	// 如果该项目不存在，则返回 nil。
	Delete(string) error

	// Exists 判断一个缓存项是否存在
	Exists(string) bool

	// Touch 重新设置缓存项的过期时间
	//
	// 如果不存在该项，不会返回 [ErrCacheMiss]，而是 nil。
	Touch(key string, ttl time.Duration) error

	// Counter 从 key 指向的值初始化一个计数器操作接口
	//
	// key 表示计数器在缓存中的名称，如果已经存在同名值，将采用该值，否则初始化为零。
	// 如果 key 指定的值无法被当作数值操作，将在 f 的调用中返回相应在的错误。
	//
	// n 表示当前的数值；
	// f 表示对数据进行操作的函数；
	// exist 表示该元素原来是否就存在；
	Counter(key string, ttl time.Duration) (n uint64, f SetCounterFunc, exist bool, err error)
}

// SetCounterFunc 为计数器增加数值的函数原型
//
// n 为增加的数值，如果为负数，则表示减少。
// 返回的数值为操作完成之后的数值，如果数值被删除，则返回 [ErrCacheMiss] 错误。
//
// 调用此方法也会更新元素的 TTL 值。
type SetCounterFunc = func(n int) (uint64, error)

// Cleanable 可清除所有缓存内容的接口
type Cleanable interface {
	Cache

	// Clean 清除所有的缓存内容
	Clean() error
}

// Driver 所有缓存驱动需要实现的接口
//
// 对于数据的序列化相关操作可直接调用 [caches.Marshal] 和 [caches.Unmarshal] 进行处理。
// 新的驱动可以采用 [github.com/issue9/cache/cachetest] 对接口进行测试，看是否符合要求。
type Driver interface {
	Cleanable

	// Ping 检测连接是否依然有效
	Ping() error

	// Close 关闭客户端
	Close() error

	// Driver 关联的底层驱动实例
	Driver() any
}

// ErrCacheMiss 当不存在缓存项时返回的错误
func ErrCacheMiss() error { return errCacheMiss }

// GetOrInit 获取缓存项
//
// 在缓存不存在时，会尝试调用 init 初始化，并调用 [Cache.Set] 存入缓存。
//
// key 和 v 相当于调用 [Cache.Get] 的参数；
// 如果 [Cache.Get] 返回 [ErrCacheMiss]，那么将调用 init 方法初始化值并写入缓存，
// 最后再调用 [Cache.Get] 返回值。
func GetOrInit[T any](cache Cache, key string, v *T, ttl time.Duration, init func() (T, error)) error {
	switch err := cache.Get(key, v); {
	case err == nil:
		return nil
	case errors.Is(err, ErrCacheMiss()):
		val, err := init()
		if err != nil {
			return err
		}
		if err := cache.Set(key, val, ttl); err != nil {
			return err
		}
		return cache.Get(key, v) // 依然有可能返回 [ErrCacheMiss]
	default:
		return err
	}
}

// Get [cache.Get] 的泛型版本
func Get[T any](cache Cache, key string) (T, error) {
	var val T
	err := cache.Get(key, &val)
	return val, err
}
