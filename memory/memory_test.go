// SPDX-License-Identifier: MIT

package memory

import (
	"testing"
	"time"

	"github.com/issue9/assert/v3"

	"github.com/issue9/cache"
	"github.com/issue9/cache/internal/testcase"
)

var _ cache.Cache = &memory{}

func TestMemory(t *testing.T) {
	a := assert.New(t, false)
	c := New(500 * time.Millisecond)
	a.NotNil(c)

	testcase.Test(a, c)
}
