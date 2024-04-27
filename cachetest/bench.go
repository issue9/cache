// SPDX-FileCopyrightText: 2024 caixw
//
// SPDX-License-Identifier: MIT

package cachetest

import (
	"testing"
	"time"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cache"
)

// BenchCounter 测试计数器的性能
func BenchCounter(b *testing.B, d cache.Driver) {
	a := assert.New(b, false)
	c, set, found, err := d.Counter("v1", cache.Forever)
	a.NotError(err).Zero(c).NotNil(set).False(found)

	b.Run("incr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := set(1)
			a.NotError(err)
		}
	})

	b.Run("decr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := set(-1)
			a.NotError(err)
		}
	})
}

// BenchBasic 测试基本功能的性能
func BenchBasic(b *testing.B, c cache.Driver) {
	a := assert.New(b, false)

	b.Run("set-string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a.NotError(c.Set("s1", "str", cache.Forever))
		}
	})
	b.Run("get-string", func(b *testing.B) {
		var v string
		for i := 0; i < b.N; i++ {
			a.NotError(c.Get("s1", &v))
		}
	})

	b.Run("set-int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a.NotError(c.Set("i1", 123, cache.Forever))
		}
	})
	b.Run("get-int", func(b *testing.B) {
		var v int
		for i := 0; i < b.N; i++ {
			a.NotError(c.Get("i1", &v))
		}
	})
}

// BenchObject 测试对象的缓存性能
func BenchObject(b *testing.B, c cache.Driver) {
	a := assert.New(b, false)

	b.Run("set-time", func(b *testing.B) {
		now := time.Now()
		for i := 0; i < b.N; i++ {
			a.NotError(c.Set("t1", now, cache.Forever))
		}
	})
	b.Run("get-time", func(b *testing.B) {
		var v time.Time
		for i := 0; i < b.N; i++ {
			a.NotError(c.Get("t1", &v))
		}
	})

	b.Run("set-object", func(b *testing.B) {
		obj := &object{Name: "x"}
		for i := 0; i < b.N; i++ {
			a.NotError(c.Set("o1", obj, cache.Forever))
		}
	})
	b.Run("get-object", func(b *testing.B) {
		var obj object
		for i := 0; i < b.N; i++ {
			a.NotError(c.Get("o1", &obj))
		}
	})
}
