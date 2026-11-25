// SPDX-FileCopyrightText: 2024-2026 caixw
//
// SPDX-License-Identifier: MIT

package cachetest

import (
	"testing"
	"time"

	"github.com/issue9/assert/v6"

	"github.com/issue9/cache"
)

// BenchCounter 测试计数器的性能
func BenchCounter(b *testing.B, d cache.Driver) {
	a := assert.New(b, false)
	c, set, found, err := d.Counter("v1", cache.Forever)
	a.NotError(err).Zero(c).NotNil(set).False(found)

	b.Run("incr", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		for b.Loop() {
			_, err = set(1)
		}
		a.NotError(err)
	})

	b.Run("decr", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		for b.Loop() {
			_, err = set(-1)
		}
		a.NotError(err)
	})
}

// BenchBasic 测试基本功能的性能
func BenchBasic(b *testing.B, c cache.Driver) {
	a := assert.New(b, false)

	b.Run("set-string", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Set("s1", "str", cache.Forever)
		}
		a.NotError(err)
	})
	b.Run("get-string", func(b *testing.B) {
		var v string
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Get("s1", &v)
		}
		a.NotError(err)
	})

	b.Run("set-int", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Set("i1", 123, cache.Forever)
		}
		a.NotError(err)
	})

	b.Run("get-int", func(b *testing.B) {
		var v int
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Get("i1", &v)
		}
		a.NotError(err)
	})
}

// BenchObject 测试对象的缓存性能
func BenchObject(b *testing.B, c cache.Driver) {
	a := assert.New(b, false)

	b.Run("set-time", func(b *testing.B) {
		now := time.Now()
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Set("t1", now, cache.Forever)
		}
		a.NotError(err)
	})

	b.Run("get-time", func(b *testing.B) {
		var v time.Time
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Get("t1", &v)
		}
		a.NotError(err)
	})

	b.Run("set-object", func(b *testing.B) {
		obj := &object{Name: "x"}
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Set("o1", obj, cache.Forever)
		}
		a.NotError(err)
	})
	b.Run("get-object", func(b *testing.B) {
		var obj object
		b.ReportAllocs()
		var err error
		for b.Loop() {
			err = c.Get("o1", &obj)
		}
		a.NotError(err)
	})
}
