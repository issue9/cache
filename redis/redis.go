// SPDX-License-Identifier: MIT

// Package redis redis 客户端的 cache 接口实现
package redis

import (
	"errors"

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

func (redis *redis) Get(key string, val interface{}) error {
	bs, err := redigo.Bytes(redis.conn.Do("GET", key))
	if errors.Is(err, redigo.ErrNil) {
		return cache.ErrCacheMiss
	} else if err != nil {
		return err
	}

	return cache.Unmarshal(bs, val)
}

func (redis *redis) Set(key string, val interface{}, seconds int) error {
	bs, err := cache.Marshal(val)
	if err != nil {
		return err
	}

	if seconds == 0 {
		_, err = redis.conn.Do("SET", key, string(bs))
		return err
	}

	_, err = redis.conn.Do("SET", key, string(bs), "EX", seconds)
	return err
}

func (redis *redis) Delete(key string) error {
	_, err := redis.conn.Do("DEL", key)
	return err
}

func (redis *redis) Exists(key string) bool {
	_, err := redis.conn.Do("GET", key)
	return err == nil || !errors.Is(err, redigo.ErrNil)
}

func (redis *redis) Clear() error {
	_, err := redis.conn.Do("FLUSHDB")
	return err
}

func (redis *redis) Close() error {
	// NOTE: 关闭服务，不能清除服务器的内容
	return redis.conn.Close()
}
