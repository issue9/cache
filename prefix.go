// SPDX-License-Identifier: MIT

package cache

type prefix struct {
	prefix string
	access Access
}

// Prefix 生成一个带有统一前缀名称的缓存访问对象
//
// c := memory.New(...)
// p := cache.Prefix("prefix_", c)
// p.Get("k1") // 相当于 c.Get("prefix_k1")
func Prefix(p string, a Access) Access {
	return &prefix{
		prefix: p,
		access: a,
	}
}

func (p *prefix) Get(key string, v interface{}) error {
	return p.access.Get(p.prefix+key, v)
}

func (p *prefix) Set(key string, val interface{}, seconds int) error {
	return p.access.Set(p.prefix+key, val, seconds)
}

func (p *prefix) Delete(key string) error { return p.access.Delete(p.prefix + key) }

func (p *prefix) Exists(key string) bool { return p.access.Exists(p.prefix + key) }
