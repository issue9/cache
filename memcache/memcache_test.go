// SPDX-License-Identifier: MIT

package memcache

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/cache"
	"github.com/issue9/cache/internal/testcase"
)

var _ cache.Cache = &memcache{}

func TestMemcache(t *testing.T) {
	a := assert.New(t)

	c := NewFromServers("localhost:11211")
	a.NotNil(c)

	testcase.Test(a, c)
}
