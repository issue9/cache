// SPDX-License-Identifier: MIT

package memcache

import (
	"log"
	"os"
	"testing"

	gm "github.com/bradfitz/gomemcache/memcache"
	"github.com/issue9/assert"

	"github.com/issue9/cache"
	"github.com/issue9/cache/internal/testcase"
)

var _ cache.Cache = &memcache{}

func TestMemcache(t *testing.T) {
	a := assert.New(t)

	client := gm.New("localhost:11211")

	c := New(log.New(os.Stderr, "", 0), client)
	a.NotNil(c)

	testcase.Test(a, c)
}
