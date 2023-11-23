// SPDX-License-Identifier: MIT

package memcache

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/issue9/cache"
	"github.com/issue9/cache/cachetest"
)

var _ cache.Cache = &memcacheDriver{}

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
