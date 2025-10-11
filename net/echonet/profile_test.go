// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
)

func TestNewProfile(t *testing.T) {
	prof := NewProfile()

	if !prof.IsProfile() {
		t.Errorf(errInvalidGroupClassCode, ErrInvalid, prof.ClassGroupCode())
	}

	mandatoryPropertyCodes := []PropertyCode{
		ProfileManufacturerCode,
		ProfileGetPropertyMap,
		ProfileSetPropertyMap,
		ProfileAnnoPropertyMap,
	}

	for _, propCode := range mandatoryPropertyCodes {
		if !prof.HasProperty(propCode) {
			t.Errorf(errMandatoryPropertyNotFound, ErrNotFound, propCode)
		}
	}
}
