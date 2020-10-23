// SPDX-License-Identifier: MIT

// Package cache 统一的缓存接口
package cache

import (
	"errors"
	"time"
)

// Forever 永不过时
const Forever = 0

var (
	// ErrCacheMiss 当不存在此缓存项时返回的错误
	ErrCacheMiss = errors.New("cache: 未找到缓存项")

	// ErrOverflow 当 Incr 和 Decr 操作超出数值范围时返回此错误
	ErrOverflow = errors.New("cache: 数值溢出")

	// ErrInvalidType 无效的类型
	//
	// Incr 和 Decr 对类型有要求，当不符合时返回此错误。
	ErrInvalidType = errors.New("cache: 类型不正确")
)

// Cache 一个统一的缓存接口
type Cache interface {
	// 获取缓存项
	Get(key string) (val interface{}, found bool)

	// 设置或是添加缓存项
	Set(key string, val interface{}, timeout time.Duration) error

	// 删除一个缓存项
	Delete(key string) error

	// 判断一个缓存项是否存在
	Exists(key string) bool

	// 清除所有的缓存内容
	Clear() error

	// 关闭整个缓存系统
	Close() error
}
