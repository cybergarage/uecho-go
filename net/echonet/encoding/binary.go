// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding

// IntegerToByte converts a specified integer to bytes.
func IntegerToByte(v uint, b []byte) {
	byteSize := len(b)
	for n := 0; n < byteSize; n++ {
		idx := ((byteSize - 1) - n)
		b[idx] = byte((v >> (uint(n) * 8)) & 0xFF)
	}
}

// ByteToInteger converts specified bytes to a integer.
func ByteToInteger(b []byte) uint {
	var v uint
	byteSize := len(b)
	for n := 0; n < byteSize; n++ {
		idx := ((byteSize - 1) - n)
		v += uint((uint(b[idx]) << (uint(n) * 8)))
	}
	return v
}
