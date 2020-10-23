// SPDX-License-Identifier: MIT

package redis

import (
	"testing"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/issue9/assert"

	"github.com/issue9/cache"
	"github.com/issue9/cache/internal/testcase"
)

var _ cache.Cache = &redis{}

func TestRedis(t *testing.T) {
	a := assert.New(t)

	options := []redigo.DialOption{
		redigo.DialConnectTimeout(time.Second),
		redigo.DialReadTimeout(time.Second),
		redigo.DialWriteTimeout(time.Second),
	}
	conn, err := redigo.Dial("tcp", "localhost:6379", options...)
	a.NotError(err).NotNil(conn)

	c := New(conn)
	a.NotNil(c)

	testcase.Test(a, c)
}
