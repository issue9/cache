// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package memory

import (
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/cache"
)

var _ cache.Cache = &Memory{}

func TestMemory_Get_Set(t *testing.T) {
	a := assert.New(t)
	mem := New(10, 50*time.Millisecond)
	a.NotNil(mem)

	v, found := mem.Get("not exists")
	a.False(found).Nil(v)

	mem.Set("k1", 123, time.Millisecond*10)
	v, found = mem.Get("k1")
	a.True(found).Equal(v, 123)

	// 超时
	time.Sleep(15 * time.Millisecond)
	v, found = mem.Get("k1")
	a.False(found).Nil(v)

	// 超时被 gc 清除
	mem.Set("k1", 123, time.Millisecond*10)
	mem.Set("k2", 123, time.Millisecond*10)
	mem.Set("k3", 123, time.Millisecond*10)
	time.Sleep(60 * time.Millisecond)
	a.Equal(len(mem.items), 0)
	a.False(mem.Exists("k1"))

	// Clear
	mem.Set("k1", 123, time.Millisecond*10)
	mem.Set("k2", 123, time.Millisecond*10)
	mem.Set("k3", 123, time.Millisecond*10)
	a.Equal(len(mem.items), 3)
	a.NotError(mem.Clear())
	a.Equal(len(mem.items), 0)

	// Close
	mem.Set("k1", 123, time.Millisecond*10)
	a.NotError(mem.Close())
	a.Equal(len(mem.items), 0)
}
