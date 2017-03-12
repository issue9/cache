// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cache

import "time"

// Cache 一个统一的缓存接口
type Cache interface {
	// 获取缓存项。
	Get(key string) (val interface{}, found bool)

	// 设置或是添加缓存项。
	Set(key string, val interface{}, timeout time.Duration) error

	// 删除一个缓存项。
	Delete(key string) error

	// 判断一个缓存项是否存在
	Exists(key string) bool

	// 清除所有的缓存内容
	Clear() error

	// 关闭整个缓存系统
	Close() error
}
