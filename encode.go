// SPDX-License-Identifier: MIT

package cache

import (
	"bytes"
	"encoding/gob"
)

// GoEncode 将 v 转换成 gob 编码
func GoEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GoDecode 将 bs 内容以 gob 规则解码到 v
func GoDecode(bs []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(v)
}
