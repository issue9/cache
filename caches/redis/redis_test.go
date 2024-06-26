// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

package redis

import (
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cache"
	"github.com/issue9/cache/cachetest"
)

var _ cache.Cache = &redisDriver{}

const redisURL = "redis://localhost:6379?dial_timeout=1&db=1&read_timeout=1&write_timeout=1"

func BenchmarkRedis(b *testing.B) {
	a := assert.New(b, false)
	c, err := NewFromURL(redisURL)
	a.NotError(err).NotNil(c)

	cachetest.BenchCounter(b, c)
	cachetest.BenchBasic(b, c)
	cachetest.BenchObject(b, c)
}

func TestRedis(t *testing.T) {
	a := assert.New(t, false)

	c, err := NewFromURL(redisURL)
	a.NotError(err).NotNil(c)

	cachetest.Basic(a, c)
	cachetest.Object(a, c)
	cachetest.Counter(a, c)

	a.NotError(c.Close())
}

func TestRedis_Close(t *testing.T) {
	a := assert.New(t, false)

	c, err := NewFromURL(redisURL)
	a.NotError(err).NotNil(c)
	a.NotError(c.Set("key", "val", cache.Forever))
	a.NotError(c.Close())

	c, err = NewFromURL(redisURL)
	a.NotError(err).NotNil(c)
	var val string
	a.NotError(c.Get("key", &val)).Equal(val, "val")
}
