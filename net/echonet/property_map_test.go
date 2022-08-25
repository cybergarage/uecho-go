// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"testing"
)

func TestNewPropertyMap(t *testing.T) {
	NewPropertyMap()
}

func TestPropertyMapCodeToFormat2(t *testing.T) {
	props := []struct {
		code        PropertyCode
		expectedRow int
		expectedBit int
	}{
		{code: 0x80, expectedRow: 1, expectedBit: 0},
		{code: 0x81, expectedRow: 2, expectedBit: 0},
		{code: 0x82, expectedRow: 3, expectedBit: 0},
		{code: 0x83, expectedRow: 4, expectedBit: 0},
		{code: 0x84, expectedRow: 5, expectedBit: 0},
		{code: 0x85, expectedRow: 6, expectedBit: 0},
		{code: 0x86, expectedRow: 7, expectedBit: 0},
		{code: 0x87, expectedRow: 8, expectedBit: 0},
		{code: 0x88, expectedRow: 9, expectedBit: 0},
		{code: 0x89, expectedRow: 10, expectedBit: 0},
		{code: 0x8A, expectedRow: 11, expectedBit: 0},
		{code: 0x8B, expectedRow: 12, expectedBit: 0},
		{code: 0x8C, expectedRow: 13, expectedBit: 0},
		{code: 0x8D, expectedRow: 14, expectedBit: 0},
		{code: 0x8E, expectedRow: 15, expectedBit: 0},
		{code: 0x8F, expectedRow: 16, expectedBit: 0},
		{code: 0x90, expectedRow: 1, expectedBit: 1},
		{code: 0x91, expectedRow: 2, expectedBit: 1},
		{code: 0x92, expectedRow: 3, expectedBit: 1},
		{code: 0x93, expectedRow: 4, expectedBit: 1},
		{code: 0x94, expectedRow: 5, expectedBit: 1},
		{code: 0x95, expectedRow: 6, expectedBit: 1},
		{code: 0x96, expectedRow: 7, expectedBit: 1},
		{code: 0x97, expectedRow: 8, expectedBit: 1},
		{code: 0x98, expectedRow: 9, expectedBit: 1},
		{code: 0x99, expectedRow: 10, expectedBit: 1},
		{code: 0x9A, expectedRow: 11, expectedBit: 1},
		{code: 0x9B, expectedRow: 12, expectedBit: 1},
		{code: 0x9C, expectedRow: 13, expectedBit: 1},
		{code: 0x9D, expectedRow: 14, expectedBit: 1},
		{code: 0x9E, expectedRow: 15, expectedBit: 1},
		{code: 0x9F, expectedRow: 16, expectedBit: 1},
		{code: 0xC0, expectedRow: 1, expectedBit: 4},
		{code: 0xC1, expectedRow: 2, expectedBit: 4},
		{code: 0xC2, expectedRow: 3, expectedBit: 4},
		{code: 0xC3, expectedRow: 4, expectedBit: 4},
		{code: 0xC4, expectedRow: 5, expectedBit: 4},
		{code: 0xC5, expectedRow: 6, expectedBit: 4},
		{code: 0xC6, expectedRow: 7, expectedBit: 4},
		{code: 0xC7, expectedRow: 8, expectedBit: 4},
		{code: 0xC8, expectedRow: 9, expectedBit: 4},
		{code: 0xC9, expectedRow: 10, expectedBit: 4},
		{code: 0xCA, expectedRow: 11, expectedBit: 4},
		{code: 0xCB, expectedRow: 12, expectedBit: 4},
		{code: 0xCC, expectedRow: 13, expectedBit: 4},
		{code: 0xCD, expectedRow: 14, expectedBit: 4},
		{code: 0xCE, expectedRow: 15, expectedBit: 4},
		{code: 0xCF, expectedRow: 16, expectedBit: 4},
		{code: 0xF0, expectedRow: 1, expectedBit: 7},
		{code: 0xF1, expectedRow: 2, expectedBit: 7},
		{code: 0xF2, expectedRow: 3, expectedBit: 7},
		{code: 0xF3, expectedRow: 4, expectedBit: 7},
		{code: 0xF4, expectedRow: 5, expectedBit: 7},
		{code: 0xF5, expectedRow: 6, expectedBit: 7},
		{code: 0xF6, expectedRow: 7, expectedBit: 7},
		{code: 0xF7, expectedRow: 8, expectedBit: 7},
		{code: 0xF8, expectedRow: 9, expectedBit: 7},
		{code: 0xF9, expectedRow: 10, expectedBit: 7},
		{code: 0xFA, expectedRow: 11, expectedBit: 7},
		{code: 0xFB, expectedRow: 12, expectedBit: 7},
		{code: 0xFC, expectedRow: 13, expectedBit: 7},
		{code: 0xFD, expectedRow: 14, expectedBit: 7},
		{code: 0xFE, expectedRow: 15, expectedBit: 7},
		{code: 0xFF, expectedRow: 16, expectedBit: 7},
	}

	for _, prop := range props {
		propCodeIdx, propCodeBit, ok := propertyMapCodeToFormat2(prop.code)
		if !ok {
			t.Errorf("%02X", prop.code)
		}
		expectedIdx := prop.expectedRow
		if propCodeIdx != expectedIdx {
			t.Errorf("%02X : idx %02X != %02X", prop.code, propCodeIdx, expectedIdx)
		}
		expectedBit := prop.expectedBit
		if propCodeBit != expectedBit {
			t.Errorf("%02X : bit %02X != %02X", prop.code, propCodeBit, expectedBit)
		}
	}
}

func TestPropertyMapFormat2ToCode(t *testing.T) {
	props := []struct {
		row          int
		bit          int
		expectedCode PropertyCode
	}{
		{row: 0, bit: 0, expectedCode: 0x80},
		{row: 8, bit: 0, expectedCode: 0x88},
		{row: 15, bit: 0, expectedCode: 0x8F},
		{row: 0, bit: 4, expectedCode: 0xC0},
		{row: 8, bit: 4, expectedCode: 0xC8},
		{row: 15, bit: 4, expectedCode: 0xCF},
		{row: 0, bit: 7, expectedCode: 0xF0},
		{row: 8, bit: 7, expectedCode: 0xF8},
		{row: 15, bit: 7, expectedCode: 0xFF},
	}

	for _, prop := range props {
		code := propertyMapFormat2BitToCode(prop.row, prop.bit)
		if code != prop.expectedCode {
			t.Errorf("%02X != %02X", code, prop.expectedCode)
		}
	}
}

func TestObjectPropertyMap(t *testing.T) {
	objs := []*Object{
		NewSuperObject().Object,
		// NewLocalNodeProfile().SuperObject.Object,
		NewStandardDeviceWithCode(0x03CE).SuperObject.Object,
	}
	// objCodes := []ObjectCode{SuperObjectCode}
	for _, obj := range objs {
		t.Run(fmt.Sprintf("%06X", obj.Code()), func(t *testing.T) {
			propMapCodes := []PropertyCode{ObjectGetPropertyMap, ObjectSetPropertyMap, ObjectAnnoPropertyMap}
			for _, propMapCode := range propMapCodes {
				t.Run(fmt.Sprintf("%02X", propMapCode), func(t *testing.T) {
					expectedCodes := make([]PropertyCode, 0)
					for _, prop := range obj.Properties() {
						switch propMapCode {
						case ObjectGetPropertyMap:
							if prop.IsReadable() {
								expectedCodes = append(expectedCodes, prop.Code())
							}
						case ObjectSetPropertyMap:
							if prop.IsWritable() {
								expectedCodes = append(expectedCodes, prop.Code())
							}
						case ObjectAnnoPropertyMap:
							if prop.IsAnnounceable() {
								expectedCodes = append(expectedCodes, prop.Code())
							}
						}
					}
					prop, ok := obj.FindProperty(propMapCode)
					if !ok {
						t.Errorf("%02X is not found", propMapCode)
						return
					}
					propCodes, err := prop.PropertyMapData()
					if err != nil {
						t.Error(err)
						return
					}
					propMapEquals := func(m, o []PropertyCode) bool {
						if len(m) != len(o) {
							return false
						}
						for _, mcode := range m {
							hasCode := false
							for _, ocode := range o {
								if ocode == mcode {
									hasCode = true
									break
								}
							}
							if !hasCode {
								return false
							}
						}
						return true
					}
					if !propMapEquals(propCodes, expectedCodes) {
						t.Errorf("%v != %v", propCodes, expectedCodes)
					}
				})
			}
		})
	}
}
