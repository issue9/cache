// SPDX-FileCopyrightText: 2017-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package caches 内置的缓存接口实现
package caches

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"time"
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
// 这是 [github.com/issue9/cache.Cache] 存储对象时的转换方法，采用 GOB 编码，
// 如需要自定义，可实现 [gob.GobEncoder] 和 [gob.GobDecoder]。
//
// 大部分时候 [github.com/issue9/cache.Driver] 的实现者直接调用此方法即可。
func Marshal(val any) ([]byte, error) {
	switch v := val.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int:
		return []byte(strconv.Itoa(v)), nil
	case int8:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int16:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int32:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int64:
		return []byte(strconv.FormatInt(v, 10)), nil
	case uint:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint8:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint16:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint32:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint64:
		return []byte(strconv.FormatUint(v, 10)), nil
	case time.Time:
		return []byte(v.Format(time.RFC3339Nano)), nil
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(bs []byte, val any) (err error) {
	switch v := val.(type) {
	case *[]byte:
		*v = bs
		return nil
	case *string:
		*v = string(bs)
		return nil
	case *int:
		*v, err = strconv.Atoi(string(bs))
		return err
	case *int8:
		vv, err := strconv.ParseInt(string(bs), 10, 8)
		if err != nil {
			return err
		}
		*v = int8(vv)
		return nil
	case *int16:
		vv, err := strconv.ParseInt(string(bs), 10, 16)
		if err != nil {
			return err
		}
		*v = int16(vv)
		return nil
	case *int32:
		vv, err := strconv.ParseInt(string(bs), 10, 32)
		if err != nil {
			return err
		}
		*v = int32(vv)
		return nil
	case *int64:
		*v, err = strconv.ParseInt(string(bs), 10, 64)
		return err
	case *uint:
		vv, err := strconv.ParseUint(string(bs), 10, strconv.IntSize)
		if err != nil {
			return err
		}
		*v = uint(vv)
		return nil
	case *uint8:
		vv, err := strconv.ParseUint(string(bs), 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(vv)
		return nil
	case *uint16:
		vv, err := strconv.ParseUint(string(bs), 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(vv)
		return nil
	case *uint32:
		vv, err := strconv.ParseUint(string(bs), 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(vv)
		return nil
	case *uint64:
		vv, err := strconv.ParseUint(string(bs), 10, 64)
		if err != nil {
			return err
		}
		*v = uint64(vv)
		return nil
	case *time.Time:
		*v, err = time.Parse(time.RFC3339Nano, string(bs))
		return err
	}

	return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(val)
}
