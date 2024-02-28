// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

package memory

import (
	"testing"
	"time"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cache"
	"github.com/issue9/cache/cachetest"
)

var _ cache.Cache = &memoryDriver{}

func TestMemory(t *testing.T) {
	a := assert.New(t, false)

	c, gc := New()
	a.NotNil(c)

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for now := range ticker.C {
			gc(now)
		}
	}()

	cachetest.Basic(a, c)
	cachetest.Object(a, c)
	cachetest.Counter(a, c)

	a.NotError(c.Close())
	ticker.Stop()
}
