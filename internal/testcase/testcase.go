// SPDX-License-Identifier: MIT

// Package testcase 提供测试用例
package testcase

import (
	"time"

	"github.com/issue9/assert"

	"github.com/issue9/cache"
)

// Test 测试 Cache 基本功能
func Test(a *assert.Assertion, c cache.Cache) {
	v, err := c.Get("not_exists")
	a.Equal(err, cache.ErrCacheMiss, "找到了一个并不存在的值").
		Nil(v, "查找一个并不存在的值，且有返回。")

	a.NotError(c.Set("k1", 123, cache.Forever))
	v, err = c.Get("k1")
	a.NotError(err, "Forever 返回未知错误 %s", err).
		Equal(v, 123, "无法正常获取 k1 的值")

	// 重设置 k1
	a.NotError(c.Set("k1", uint(789), 1*time.Hour))
	v, err = c.Get("k1")
	a.NotError(err, "1*time.Hover 的值 k1 返回错误信息 %s", err).
		Equal(v, 789, "无法正常获取 k1 的值")

	// 被 delete 删除
	a.NotError(c.Delete("k1"))
	v, err = c.Get("k1")
	a.Equal(err, cache.ErrCacheMiss, "k1 并未被回收").
		Nil(v, "被删除之后值并未为空：%+v", v)

	// 超时被回收
	a.NotError(c.Set("k1", 123, time.Millisecond*10))
	a.NotError(c.Set("k2", 456, time.Millisecond*10))
	a.NotError(c.Set("k3", 789, time.Millisecond*10))
	time.Sleep(1*time.Second + 500*time.Microsecond)
	a.False(c.Exists("k1"), "k1 超时且未被回收")
	a.False(c.Exists("k2"), "k2 超时且未被回收")
	a.False(c.Exists("k3"), "k3 超时且未被回收")

	// Clear
	a.NotError(c.Set("k1", 123, time.Millisecond*10))
	a.NotError(c.Set("k2", 456, time.Millisecond*10))
	a.NotError(c.Set("k3", 789, time.Millisecond*10))
	a.NotError(c.Clear())
	a.False(c.Exists("k1"), "clear 之后 k1 依然存在").
		False(c.Exists("k2"), "clear 之后 k2 依然存在").
		False(c.Exists("k3"), "clear 之后 k3 依然存在")

	// Close
	a.NotError(c.Set("k1", 123, time.Millisecond*10))
	a.NotError(c.Close())
}

// TestObject 测试对象的存储
func TestObject(a *assert.Assertion, c cache.Cache) {
	type o struct {
		Name string
		age  int
	}

	obj := &o{Name: "test", age: 5}
	obj2 := &o{Name: "test", age: 5}

	a.NotError(c.Set("obj", obj, cache.Forever))
	v, err := c.Get("obj")
	a.NotError(err).
		Equal(v, obj2, "obj not equal\nv1:%+v\nv2:%+v\n", v, obj2)
}
