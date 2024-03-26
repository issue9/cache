// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package caches 内置的缓存接口实现
package caches

import (
	"bytes"
	"encoding/gob"
)

// NOTE: 如果支持多种编码方式，为了编码和解码是同一种方式，需要对象同时实现编码和解码，比如
//  type Serializer interface {
//      encoding.TextMarshaler
//      encoding.TextUnarshaler
//  }
// 但是大部分实现者的对 Marshal 和 Unmarshal 的 receiver 是不一样的，
// 导致在 Marshal 和 Unmarshal 对是否实现接口的判断结果也不同。

// Marshal 序列化对象
//
// 这是 [cache.Cache] 存储对象时的转换方法，采用 GOB 编码，
// 如需要自定义，可实现 [gob.GobEncoder] 和 [gob.GobDecoder]。
//
// 大部分时候 [cache.Driver] 的实现者直接调用此方法即可。
func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(bs []byte, v any) error { return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(v) }
