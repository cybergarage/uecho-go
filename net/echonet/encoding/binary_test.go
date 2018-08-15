// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding

import (
	"testing"
)

func TestBinaryEncoding(t *testing.T) {
	var n uint

	intBytes := make([]byte, 1)
	for n = 0; n <= 0xFF; n++ {
		IntegerToByte(n, intBytes)
		if n != ByteToInteger(intBytes) {
			t.Errorf("[1:%d] : %d != %d", n, ByteToInteger(intBytes), n)
		}
	}

	intBytes = make([]byte, 2)
	for n = 0; n <= 0xFFFF; n += (0xFFFF / 0xFF) {
		IntegerToByte(n, intBytes)
		if n != ByteToInteger(intBytes) {
			t.Errorf("[2:%d] : %d != %d", n, ByteToInteger(intBytes), n)
		}
	}

	intBytes = make([]byte, 3)
	for n = 0; n <= 0xFFFFFF; n += (0xFFFFFF / 0xFF) {
		IntegerToByte(n, intBytes)
		if n != ByteToInteger(intBytes) {
			t.Errorf("[3:%d] : %d != %d", n, ByteToInteger(intBytes), n)
		}
	}

	intBytes = make([]byte, 4)
	for n = 0; n < 0xFFFFFFFF; n += (0xFFFFFFFF / 0xFF) {
		IntegerToByte(n, intBytes)
		if n != ByteToInteger(intBytes) {
			t.Errorf("[4:%d] : %d != %d", n, ByteToInteger(intBytes), n)
		}
	}
}
