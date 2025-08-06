// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

func propertyMapCodeToFormat2(propCode PropertyCode) (int, int, bool) {
	if (propCode < PropertyCodeMin) || (PropertyCodeMax < propCode) {
		return 0, 0, false
	}
	// 1 <= propCodeIdx <= 16
	propCodeIdx := int(((propCode - PropertyCodeMin) & 0x0F)) + 1
	// 0 <= propCodeIdx <= 7
	propCodeBit := (((int(propCode-PropertyCodeMin) & 0xF0) >> 4) & 0x0F)
	return propCodeIdx, propCodeBit, true
}

func isPropertyMapDescriptionFormat1(n int) bool {
	return (n <= PropertyMapFormat1MaxSize)
}

func isPropertyMapDescriptionFormat2(n int) bool {
	return !isPropertyMapDescriptionFormat1(n)
}

func propertyMapFormat2BitToCode(row int, bit int) PropertyCode {
	// 0 <= bit <= 7
	code := (0x10 * bit) + PropertyCodeMin
	// 0 <= row <= 15
	code += row
	return PropertyCode(code)
}

func propertyMapFormat2ByteToCodes(row int, b byte) []PropertyCode {
	codes := make([]PropertyCode, 0)
	for n := range 8 {
		bit := byte((0x01 << n) & 0x0F)
		if (b & bit) == 0 {
			continue
		}
		codes = append(codes, propertyMapFormat2BitToCode(row, n))
	}
	return codes
}
