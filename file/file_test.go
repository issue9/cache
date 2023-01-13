// SPDX-License-Identifier: MIT

package file

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v3"

	"github.com/issue9/cache"
	"github.com/issue9/cache/internal/testcase"
)

var _ cache.Cache = &file{}

func TestFile(t *testing.T) {
	a := assert.New(t, false)

	c := New("./testdir", 500*time.Millisecond, log.New(os.Stderr, "", 0))
	a.NotNil(c)

	testcase.Test(a, c)
}

func TestFile_Close(t *testing.T) {
	a := assert.New(t, false)

	c := New("./testdir", 500*time.Millisecond, log.New(os.Stderr, "", 0))
	a.NotNil(c)
	a.NotError(c.Set("key", "val", cache.Forever))
	a.NotError(c.Close())

	c = New("./testdir", 500*time.Millisecond, log.New(os.Stderr, "", 0))
	a.NotNil(c)
	var val string
	a.NotError(c.Get("key", &val)).Equal(val, "val")
}
