// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"testing"
)

const (
	errorMandatoryPropertyNotFound     = "mandatory property (%0X) not found"
	errorMandatoryPropertyDataNotFound = "mandatory property data (%0X) not found"
)

func TestSuperObject(t *testing.T) {
	obj := NewSuperObject()

	mandatoryPropertyCodes := []PropertyCode{
		ProfileManufacturerCode,
	}

	for _, propCode := range mandatoryPropertyCodes {
		t.Run(fmt.Sprintf("%02X", propCode), func(t *testing.T) {
			if !obj.HasProperty(propCode) {
				t.Errorf(errorMandatoryPropertyNotFound, propCode)
			}
		})
	}

	testObjectPropertyMaps(t, obj.Object)
}

func testObjectPropertyMaps(t *testing.T, obj *Object) {
	t.Helper()

	propMapCodes := []PropertyCode{
		ProfileGetPropertyMap,
		ProfileSetPropertyMap,
		ProfileAnnoPropertyMap,
	}

	for _, propCode := range propMapCodes {
		t.Run(fmt.Sprintf("%02X", propCode), func(t *testing.T) {
			prop, ok := obj.FindProperty(propCode)
			if !ok {
				t.Errorf(errorMandatoryPropertyNotFound, propCode)
			}
			propData := prop.Data()
			if len(propData) < 1 {
				t.Errorf(errorMandatoryPropertyDataNotFound, propCode)
			}
			propCnt := int(propData[0])
			if propCnt <= PropertyMapFormat1MaxSize {
				for n, code := range propData {
					if n == 0 {
						continue
					}
					if code == 0x00 {
						t.Errorf("[%03d] %02X", n, code)
					}
				}
			}
		})
	}
}

func TestPropertyMapFormat2(t *testing.T) {
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
		propCodeIdx, propCodeBit, ok := propCodeToFormat2(prop.code)
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
