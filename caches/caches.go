// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package caches 内置的缓存接口实现
package caches

import (
	"bytes"
	"encoding/gob"
)

// Marshal 序列化对象
//
// 这是 [cache.Cache] 存储对象时的转换方法，按以下顺序进行：
//   - 是否实现 [Serializer]；
//   - 是否同时实现了 [encoding.TextMarshaler] 和 [encoding.TextUnmarshaler]；
//   - 采用 gob 编码；
//
// [Unmarshal] 按同样的顺序执行。
//
// 大部分时候 [cache.Driver] 的实现者直接调用此方法即可，
// 如果需要自己实现，需要注意 [Serializer] 接口的判断。
func Marshal(v any) ([]byte, error) {
	switch m := v.(type) {
	case Serializer:
		return m.MarshalCache()
	case textSerializer:
		return m.MarshalText()
	default:
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)

		if err := enc.Encode(v); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

func Unmarshal(bs []byte, v any) error {
	switch u := v.(type) {
	case Serializer:
		return u.UnmarshalCache(bs)
	case textSerializer:
		return u.UnmarshalText(bs)
	default:
		return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(v)
	}
}
