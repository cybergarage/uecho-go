// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"testing"
)

const (
	errMandatoryPropertyNotFound     = "%s: mandatory property (%0X)"
	errMandatoryPropertyDataNotFound = "%s, mandatory property data (%0X)"
)

func TestSuperObject(t *testing.T) {
	obj := NewSuperObject()

	mandatoryPropertyCodes := []PropertyCode{
		ProfileManufacturerCode,
	}

	for _, propCode := range mandatoryPropertyCodes {
		t.Run(fmt.Sprintf("%02X", propCode), func(t *testing.T) {
			if !obj.HasProperty(propCode) {
				t.Errorf(errMandatoryPropertyNotFound, ErrNotFound, propCode)
			}
		})
	}

	testObjectPropertyMaps(t, obj)
}

func testObjectPropertyMaps(t *testing.T, obj Object) {
	t.Helper()

	propMapCodes := []PropertyCode{
		ProfileGetPropertyMap,
		ProfileSetPropertyMap,
		ProfileAnnoPropertyMap,
	}

	for _, propCode := range propMapCodes {
		t.Run(fmt.Sprintf("%02X", propCode), func(t *testing.T) {
			prop, ok := obj.LookupProperty(propCode)
			if !ok {
				t.Errorf(errMandatoryPropertyNotFound, ErrNotFound, propCode)
			}
			propData := prop.Data()
			if len(propData) < 1 {
				t.Errorf(errMandatoryPropertyDataNotFound, ErrNotFound, propCode)
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
