// SPDX-FileCopyrightText: 2017-2026 caixw
//
// SPDX-License-Identifier: MIT

package memcache

import (
	"testing"
	"time"

	"github.com/issue9/assert/v6"

	"github.com/issue9/cache"
	"github.com/issue9/cache/cachetest"
)

var _ cache.Cache = &memcacheDriver{}

func TestCalcTTL(t *testing.T) {
	a := assert.New(t, false)

	a.PanicString(func() {
		calcTTL(time.Duration(1<<63 - 1))
	}, "TTL 转换为秒数时不能大于 math.MaxInt32")

	tests := []struct {
		name string
		ttl  time.Duration
		want int32
	}{
		{name: "forever", ttl: 0, want: 0},
		{name: "nano second", ttl: time.Nanosecond, want: 1},
		{name: "1 second", ttl: time.Second, want: 1},
		{name: "fractional seconds", ttl: 1500 * time.Millisecond, want: 2},
		{name: "month boundary", ttl: month, want: int32(month / time.Second)},
		{name: "round up to month boundary", ttl: month - time.Nanosecond, want: int32(month / time.Second)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := a.New(t)
			a.Equal(tt.want, calcTTL(tt.ttl))
		})
	}
}

func BenchmarkMemcache(b *testing.B) {
	a := assert.New(b, false)
	c := New("localhost:11211")
	a.NotNil(c)

	cachetest.BenchCounter(b, c)
	cachetest.BenchBasic(b, c)
	cachetest.BenchObject(b, c)
}

func TestMemcache(t *testing.T) {
	a := assert.New(t, false)

	c := New("localhost:11211")
	a.NotNil(c)

	cachetest.Basic(a, c)
	cachetest.Object(a, c)
	cachetest.Counter(a, c)

	a.NotError(c.Close())
}

func TestMemcache_Close(t *testing.T) {
	a := assert.New(t, false)

	c := New("localhost:11211")
	a.NotNil(c)
	a.NotError(c.Set("key", "val", cache.Forever))
	a.NotError(c.Close())

	c = New("localhost:11211")
	a.NotNil(c)
	var val string
	a.NotError(c.Get("key", &val)).Equal(val, "val")
}
