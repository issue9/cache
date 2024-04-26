// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package cachetest 缓存的测试用例
package cachetest

import (
	"time"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cache"
)

// Counter 测试计数器
func Counter(a *assert.Assertion, d cache.Driver) {
	n, set, err := d.Counter("v1", time.Second)
	a.NotError(err).Zero(n).NotNil(set)

	v1, err := set(0)
	a.NotError(err).Equal(v1, 0)

	v1, err = set(5)
	a.NotError(err).Equal(v1, 5)

	v1, err = set(-3)
	a.NotError(err).Equal(v1, 2)
	v1, err = set(-1)
	a.NotError(err).Equal(v1, 1)

	a.True(d.Exists("v1"))

	a.NotError(d.Delete("v1"))
	a.False(d.Exists("v1"))

	v2, err := set(3) // 已经被删除
	a.ErrorIs(err, cache.ErrCacheMiss()).Zero(v2)
	v2, err = set(3) // 已经被删除
	a.ErrorIs(err, cache.ErrCacheMiss()).Zero(v2)

	// 多个 Counter 指向同一个 key

	n1, set1, err := d.Counter("v3", time.Second)
	a.NotError(err).Zero(n1).NotNil(set1)
	v1, err = set1(5)
	a.NotError(err).Equal(v1, 5)

	n2, set2, err := d.Counter("v3", time.Second)
	a.NotError(err).Equal(n2, 5)
	v2, err = set2(-5)
	a.NotError(err).Equal(v2, 0)
}

// Basic 测试基本功能
func Basic(a *assert.Assertion, c cache.Driver) {
	// driver
	a.NotNil(c.Driver())

	var v string
	err := c.Get("not_exists", &v)
	a.ErrorIs(err, cache.ErrCacheMiss(), "找到了一个并不存在的值").
		Zero(v, "查找一个并不存在的值，且有返回。")

	a.Nil(c.Touch("not_exists", time.Second))

	a.NotError(c.Set("k1", 123, cache.Forever))
	var num int
	err = c.Get("k1", &num)
	a.NotError(err, "Forever 返回未知错误 %s", err).
		Equal(num, 123).
		NotError(c.Touch("k1", time.Second))

	now := time.Now()
	a.NotError(c.Set("t1", now, cache.Forever))
	var t time.Time
	err = c.Get("t1", &t)
	a.NotError(err, "Forever 返回未知错误 %s", err).
		NotZero(t)

	// 重新设置 k1
	a.NotError(c.Set("k1", uint(789), time.Minute))
	var unum uint
	err = c.Get("k1", &unum)
	a.NotError(err, "1*time.Hour 的值 k1 返回错误信息 %s", err).
		Equal(unum, 789, "无法正常获取重新设置之后 k1 的值 v1=%s, v2=%d", v, 789)

	// 被 delete 删除
	a.NotError(c.Delete("k1"))
	err = c.Get("k1", &unum)
	a.Equal(err, cache.ErrCacheMiss(), "k1 并未被回收").
		Zero(v, "被删除之后值并未为空：%+v", v)

	// 超时被回收
	a.NotError(c.Set("k1", 123, time.Second))
	a.NotError(c.Set("k2", 456, time.Second))
	a.NotError(c.Set("k3", 789, time.Second))
	time.Sleep(2 * time.Second)
	a.False(c.Exists("k1"), "k1 超时且未被回收")
	a.False(c.Exists("k2"), "k2 超时且未被回收")
	a.False(c.Exists("k3"), "k3 超时且未被回收")

	// Clean
	a.NotError(c.Set("k1", 123, time.Second))
	a.NotError(c.Set("k2", 456, time.Second))
	a.NotError(c.Set("k3", 789, time.Second))
	a.NotError(c.Clean())
	a.False(c.Exists("k1"), "clean 之后 k1 依然存在").
		False(c.Exists("k2"), "clean 之后 k2 依然存在").
		False(c.Exists("k3"), "clean 之后 k3 依然存在")
}

type object struct {
	Name string
	age  int
}

// Object 测试对象的缓存
func Object(a *assert.Assertion, c cache.Driver) {
	obj := &object{Name: "test", age: 5}

	a.NotError(c.Set("obj", obj, cache.Forever))
	var v object
	err := c.Get("obj", &v)
	a.NotError(err).Equal(&v, &object{Name: "test"}) // 私有字段，无法解码
}
