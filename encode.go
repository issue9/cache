// SPDX-License-Identifier: MIT

package cache

import (
	"bytes"
	"encoding/gob"
)

// GoEncode 将 v 转换成 []byte
func GoEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GoDecode 将由 GoEncode 编码的内容解码至 v
func GoDecode(bs []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(v)
}
