// SPDX-License-Identifier: MIT

// Package cache 统一的缓存接口
package cache

import "errors"

// 常用时间的定义
const (
	Forever     = 0 //  永不过时
	OneMinute   = 60
	FiveMinutes = 5 * OneMinute
	TenMinutes  = 10 * OneMinute
	HalfHour    = 30 * OneMinute
	OneHour     = 60 * OneMinute
	HalfDay     = 12 * OneHour
	OneDay      = 24 * OneHour
	OneWeek     = 7 * OneDay
	ThirtyDays  = 30 * OneDay // 30 天
	SixtyDays   = 60 * OneDay // 60 天
	NinetyDays  = 90 * OneDay // 90 天
)

// ErrCacheMiss 当不存在缓存项时返回的错误
var ErrCacheMiss = errors.New("cache: 未找到缓存项")

// ErrInvalidKey key 的格式无效
//
// 部分适配器对 key 可能是有特殊要求的，
// 比如在文件系统中，可能会不允许在 key 中包含 .. 或是 / 等，碰到此类情况，可返回此错误信息。
var ErrInvalidKey = errors.New("cache: key 的格式不符合要求")

// Cache 缓存的访问和控制接口
type Cache interface {
	Access

	// 清除所有的缓存内容
	Clear() error

	// 关闭整个缓存系统
	Close() error
}

// Access 缓存内容的访问接口
type Access interface {
	// 获取缓存项
	//
	// 当前不存在时，返回 ErrCacheMiss 错误。
	Get(key string) (interface{}, error)

	// 设置或是添加缓存项
	//
	// seconds 表示过了该时间，缓存项将被回收。如果该值为 0，该值永远不会回收。
	Set(key string, val interface{}, seconds int) error

	// 删除一个缓存项
	Delete(key string) error

	// 判断一个缓存项是否存在
	Exists(key string) bool
}
