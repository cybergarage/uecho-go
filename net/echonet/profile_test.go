// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
)

func TestNewProfile(t *testing.T) {
	prop := NewProfile()

	if !prop.IsProfile() {
		t.Errorf(errorInvalidGroupClassCode, prop.GetClassGroupCode())
	}

	mandatoryPropertyCodes := []PropertyCode{
		ProfileManufacturerCode,
		ProfileGetPropertyMap,
		ProfileSetPropertyMap,
		ProfileAnnoPropertyMap,
	}

	for _, propCode := range mandatoryPropertyCodes {
		if !prop.HasProperty(propCode) {
			t.Errorf(errorMandatoryPropertyNotFound, propCode)
		}
	}
}
