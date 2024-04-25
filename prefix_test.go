// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

package cache_test

import (
	"testing"

	"github.com/issue9/assert/v4"
	"github.com/issue9/cache"

	"github.com/issue9/cache/caches/memory"
)

func TestPrefix(t *testing.T) {
	a := assert.New(t, false)

	d, _ := memory.New()
	a.NotNil(d)

	p1 := cache.Prefix(d, "p1")
	p2 := cache.Prefix(p1, "p2")
	p2.Set("key", 5, cache.Forever)
	a.True(d.Exists("p1p2key"))

	p2.Delete("key")
	a.False(d.Exists("p1p2key"))
}
