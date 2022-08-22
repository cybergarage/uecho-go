// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
)

func TestNewProfile(t *testing.T) {
	prof := NewProfile()

	if !prof.IsProfile() {
		t.Errorf(errorInvalidGroupClassCode, prof.ClassGroupCode())
	}

	mandatoryPropertyCodes := []PropertyCode{
		ProfileManufacturerCode,
		ProfileGetPropertyMap,
		ProfileSetPropertyMap,
		ProfileAnnoPropertyMap,
	}

	for _, propCode := range mandatoryPropertyCodes {
		if !prof.HasProperty(propCode) {
			t.Errorf(errorMandatoryPropertyNotFound, propCode)
		}
	}
}
