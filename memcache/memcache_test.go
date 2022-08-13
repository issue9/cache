// SPDX-License-Identifier: MIT

package memcache

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/issue9/cache"
	"github.com/issue9/cache/internal/testcase"
)

var _ cache.Cache = &memcache{}

func TestMemcache(t *testing.T) {
	a := assert.New(t, false)

	c := NewFromServers("localhost:11211")
	a.NotNil(c)

	testcase.Test(a, c)
}

func TestMemcache_Close(t *testing.T) {
	a := assert.New(t, false)

	c := NewFromServers("localhost:11211")
	a.NotNil(c)
	a.NotError(c.Set("key", "val", cache.Forever))
	a.NotError(c.Close())

	c = NewFromServers("localhost:11211")
	a.NotNil(c)
	val, err := c.Get("key")
	a.NotError(err).Equal(val, "val")
}
