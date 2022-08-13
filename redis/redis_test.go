// SPDX-License-Identifier: MIT

package redis

import (
	"testing"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/issue9/assert/v3"

	"github.com/issue9/cache"
	"github.com/issue9/cache/internal/testcase"
)

var _ cache.Cache = &redis{}

func TestRedis(t *testing.T) {
	a := assert.New(t, false)

	c := New(dial(a))
	a.NotNil(c)

	testcase.Test(a, c)
}

func TestRedis_Close(t *testing.T) {
	a := assert.New(t, false)

	c := New(dial(a))
	a.NotNil(c)
	a.NotError(c.Set("key", "val", cache.Forever))
	a.NotError(c.Close())

	c = New(dial(a))
	a.NotNil(c)
	val, err := c.Get("key")
	a.NotError(err).Equal(val, "val")
}

func dial(a *assert.Assertion) redigo.Conn {
	options := []redigo.DialOption{
		redigo.DialConnectTimeout(time.Second),
		redigo.DialReadTimeout(time.Second),
		redigo.DialWriteTimeout(time.Second),
	}
	conn, err := redigo.Dial("tcp", "localhost:6379", options...)
	a.NotError(err).NotNil(conn)

	return conn
}
