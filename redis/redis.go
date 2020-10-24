// SPDX-License-Identifier: MIT

// Package redis redis 客户端的 cache 接口实现
package redis

import (
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/issue9/cache"
)

type redis struct {
	conn redigo.Conn
}

// New 返回 redis 的缓存实现
func New(conn redigo.Conn) cache.Cache {
	return &redis{
		conn: conn,
	}
}

func (redis *redis) Get(key string) (val interface{}, err error) {
	bs, err := redigo.Bytes(redis.conn.Do("GET", key))
	if err == redigo.ErrNil {
		return nil, cache.ErrCacheMiss
	} else if err != nil {
		return nil, err
	}

	if err := cache.GoDecode(bs, &val); err != nil {
		return nil, err
	}

	return val, nil
}

func (redis *redis) Set(key string, val interface{}, timeout time.Duration) error {
	bs, err := cache.GoEncode(&val)
	if err != nil {
		return err
	}

	if timeout == cache.Forever {
		_, err = redis.conn.Do("SET", key, string(bs))
		return err
	}

	exp := int32(timeout.Seconds())
	if exp < 1 {
		exp = 1
	}
	_, err = redis.conn.Do("SET", key, string(bs), "EX", exp)
	return err
}

func (redis *redis) Delete(key string) error {
	_, err := redis.conn.Do("DEL", key)
	return err
}

func (redis *redis) Exists(key string) bool {
	_, found := redis.Get(key)
	return found != cache.ErrCacheMiss
}

func (redis *redis) Clear() error {
	_, err := redis.conn.Do("FLUSHDB")
	return err
}

func (redis *redis) Close() error {
	// NOTE: 关闭服务，不能清除服务器的内容
	return redis.conn.Close()
}
