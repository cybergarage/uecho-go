// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

func propCodeToFormat2(propCode PropertyCode) (int, int, bool) {
	if (propCode < PropertyCodeMin) || (PropertyCodeMax < propCode) {
		return 0, 0, false
	}
	propCodeIdx := int(((propCode - PropertyCodeMin) & 0x0F)) + 1
	propCodeBit := (((int(propCode-PropertyCodeMin) & 0xF0) >> 4) & 0x0F)
	return propCodeIdx, propCodeBit, true
}

func isPropertyMapDescriptionFormat2(n int) bool {
	return (n <= PropertyMapFormat1MaxSize)
}
