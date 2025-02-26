// SPDX-FileCopyrightText: 2017-2025 caixw
//
// SPDX-License-Identifier: MIT

package cache_test

import (
	"testing"
	"time"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cache"
	"github.com/issue9/cache/caches/memory"
)

func TestGetOrInit(t *testing.T) {
	a := assert.New(t, false)

	d := memory.New()
	a.NotNil(d)

	var v1 string
	err := cache.GetOrInit(d, "key", &v1, time.Second, func(v *string) error { *v = "5"; return nil })
	a.NotError(err).
		Equal(v1, "5")

	a.NotError(d.Set("key", "10", cache.Forever))

	var v2 string
	err = cache.GetOrInit(d, "key", &v2, time.Second, func(v *string) error { *v = "10"; return nil })
	a.NotError(err).
		Equal(v2, "10")
}

func TestGet(t *testing.T) {
	a := assert.New(t, false)

	d := memory.New()
	a.NotNil(d)

	v1, err := cache.Get[string](d, "v1")
	a.Equal(err, cache.ErrCacheMiss()).Empty(v1)

	a.NotError(d.Set("v1", "string", cache.Forever))
	v2, err := cache.Get[string](d, "v1")
	a.NotError(err).Equal(v2, "string")
}
