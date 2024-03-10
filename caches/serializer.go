// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

package caches

import "encoding"

// Serializer 缓存系统存取数据时采用的序列化方法
//
// 如果你存储的对象实现了该接口，那么在存取数据时，会采用此方法将对象进行编解码。
// 否则会采用默认的方法进行编辑码。
//
// 实现 Serializer 可以拥有更高效的转换效率，以及一些默认行为不可实现的功能，
// 比如需要对拥有不可导出的字段进行编解码。
//
// NOTE: 必须要有自定义的序列化接口，而不是直接采用 [encoding.TextMarshaler]
// 防止无意中改变了 JSON 的编码方式。
type Serializer interface {
	MarshalCache() ([]byte, error)
	UnmarshalCache([]byte) error
}

type textSerializer interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}
