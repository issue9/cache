// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package redis 适配 redis 的实现
package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/issue9/cache"
	"github.com/issue9/cache/caches"
)

type redisDriver struct {
	conn         *redis.Client
	decrByScript *redis.Script
}

type redisCounter struct {
	driver *redisDriver
	key    string
	ttl    time.Duration
}

// redis 处理 DECRBY 的事务脚本
const redisDecrByScript = `local cnt = redis.call('DECRBY', KEYS[1], ARGV[1])
if cnt < 0 then
    redis.call('SET', KEYS[1], '0')
end
return (cnt < 0 and 0 or cnt)`

// NewFromURL 声明基于 [redis] 的缓存系统
//
// url 为符合 [Redis URI scheme] 的字符串。
// [cache.Driver.Driver] 的返回类型为 [redis.Client]。
//
// [Redis URI scheme]: https://www.iana.org/assignments/uri-schemes/prov/redis
// [redis]: https://redis.io/
func NewFromURL(url string) (cache.Driver, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return New(redis.NewClient(opt)), nil
}

// New 声明基于 redis 的缓存系统
func New(c *redis.Client) cache.Driver {
	return &redisDriver{
		conn:         c,
		decrByScript: redis.NewScript(redisDecrByScript),
	}
}

func (d *redisDriver) Get(key string, val any) error {
	bs, err := d.conn.Get(context.Background(), key).Bytes()
	if errors.Is(err, redis.Nil) {
		return cache.ErrCacheMiss()
	} else if err != nil {
		return err
	}

	return caches.Unmarshal(bs, val)
}

func (d *redisDriver) Set(key string, val any, ttl time.Duration) error {
	bs, err := caches.Marshal(val)
	if err != nil {
		return err
	}
	return d.conn.Set(context.Background(), key, bs, ttl).Err()
}

func (d *redisDriver) Delete(key string) error { return d.conn.Del(context.Background(), key).Err() }

func (d *redisDriver) Exists(key string) bool {
	rslt, err := d.conn.Exists(context.Background(), key).Result()
	return err == nil && rslt > 0
}

func (d *redisDriver) Clean() error { return d.conn.FlushDB(context.Background()).Err() }

func (d *redisDriver) Close() error { return d.conn.Close() }

func (d *redisDriver) Driver() any { return d.conn }

func (d *redisDriver) Counter(key string, ttl time.Duration) (cache.Counter, error) {
	if err := d.conn.SetNX(context.Background(), key, 0, ttl).Err(); err != nil {
		return nil, err
	}

	return &redisCounter{
		driver: d,
		key:    key,
		ttl:    ttl,
	}, nil
}

func (c *redisCounter) Incr(n uint64) (uint64, error) {
	if !c.driver.Exists(c.key) {
		return 0, cache.ErrCacheMiss()
	}

	rslt, err := c.driver.conn.IncrBy(context.Background(), c.key, int64(n)).Result()
	if err == nil && c.ttl > 0 {
		_, err = c.driver.conn.Expire(context.Background(), c.key, c.ttl).Result()
	}
	return uint64(rslt), err
}

func (c *redisCounter) Decr(n uint64) (uint64, error) {
	if !c.driver.Exists(c.key) {
		return 0, cache.ErrCacheMiss()
	}

	v, err := c.driver.decrByScript.Run(context.Background(), c.driver.conn, []string{c.key}, int64(n)).Int64()
	if err == nil && c.ttl > 0 {
		_, err = c.driver.conn.Expire(context.Background(), c.key, c.ttl).Result()
	}
	return uint64(v), err
}

func (c *redisCounter) Value() (uint64, error) {
	s, err := c.driver.conn.Get(context.Background(), c.key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, cache.ErrCacheMiss()
	} else if err != nil {
		return 0, err
	}
	return strconv.ParseUint(s, 10, 64)
}

func (c *redisCounter) Delete() error { return c.driver.conn.Del(context.Background(), c.key).Err() }
