// Copyright (C) 2018 Satoshi Konno. All rights reserved.
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
