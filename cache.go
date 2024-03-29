// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

//go:generate web locale -l=und -m -f=yaml ./
//go:generate web update-locale -src=./locales/und.yaml -dest=./locales/zh-CN.yaml

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
	//
	// NOTE: 不能正确获取由 [Cache.Counter] 设置的值，[Cache.Counter]
	// 的实现是基于缓存系统原生的功能，存储方式与当前的实现可能是不同的。
	Get(key string, v any) error

	// Set 设置或是添加缓存项
	//
	// key 表示保存该数据的唯一 ID；
	// val 表示保存的数据对象，如果是结构体，需要所有的字段都是公开的或是实现了
	// [Serializer] 接口，否则在缓存过程中将失去这些非公开的字段。
	// ttl 表示过了该时间，缓存项将被回收。如果该值为 [Forever]，该值永远不会回收。
	Set(key string, val any, ttl time.Duration) error

	// Delete 删除一个缓存项
	Delete(string) error

	// Exists 判断一个缓存项是否存在
	Exists(string) bool

	// Counter 返回计数器操作接口
	//
	// val 和 ttl 表示在该计数器不存在时，初始化的值以及回收时间。
	Counter(key string, val uint64, ttl time.Duration) Counter
}

// Counter 计数器需要实现的接口
//
// [Cache] 支持自定义的序列化接口，但是对于自增等纯数值操作，
// 各个缓存服务都实现自有的快捷操作，无法适用自定义的序列化。
//
// 由 Counter 设置的值，可能无法由 [Cache.Get] 读取到正确的值，
// 但是可以由 [Cache.Exists]、[Cache.Delete] 和 [Cache.Set] 进行相应的操作。
//
// 各个驱动对自增值的类型定义是不同的，
// 只有在 [0, math.MaxInt32) 范围内的数值是安全的。
type Counter interface {
	// Incr 增加计数并返回增加后的值
	Incr(uint64) (uint64, error)

	// Decr 减少数值并返回减少后的值
	Decr(uint64) (uint64, error)

	// Value 返回该计数器的当前值
	Value() (uint64, error)

	// Delete 删除当前的计数器
	//
	// 当计数器不存在时，不应该返回错误。
	Delete() error
}

// Cleanable 可清除所有缓存内容的接口
type Cleanable interface {
	Cache

	// Clean 清除所有的缓存内容
	Clean() error
}

// Driver 所有缓存驱动需要实现的接口
//
// 对于数据的序列化相关操作可直接调用 [caches.Marshal] 和 [caches.Unmarshal]
// 进行处理，如果需要自行处理，需要对实现 [caches.Serializer] 接口的数据进行处理。
//
// 新的驱动可以采用 [github.com/issue9/cache/cachetest] 对接口进行测试，看是否符合要求。
type Driver interface {
	Cleanable

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
