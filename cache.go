// SPDX-License-Identifier: MIT

package cache

import (
	"errors"
	"time"
)

// 一些全局的错误信息
var (
	ErrKeyNotExists             = errors.New("不存在的项")
	ErrUintNotAllowLessThanZero = errors.New("uint 值不能小于 0")
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

	// 增加计数
	Incr(key string) error

	// 减小计数
	Decr(key string) error

	// 清除所有的缓存内容
	Clear() error

	// 关闭整个缓存系统
	Close() error
}
