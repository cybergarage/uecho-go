// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
)

func TestNewLocalNodeProfile(t *testing.T) {
	prop := NewLocalNodeProfile()

	if !prop.IsProfile() {
		t.Errorf(errorInvalidGroupClassCode, prop.GetClassGroupCode())
	}

	mandatoryPropertyCodes := []PropertyCode{
		// Profile
		ProfileManufacturerCode,
		ProfileGetPropertyMap,
		ProfileSetPropertyMap,
		ProfileAnnoPropertyMap,
		// Node Profile
		NodeProfileClassOperatingStatus,
		NodeProfileClassVersionInformation,
		NodeProfileClassIdentificationNumber,
		NodeProfileClassNumberOfSelfNodeInstances,
		NodeProfileClassNumberOfSelfNodeClasses,
		NodeProfileClassInstanceListNotification,
		NodeProfileClassSelfNodeInstanceListS,
		NodeProfileClassSelfNodeClassListS,
	}

	for _, propCode := range mandatoryPropertyCodes {
		if !prop.HasProperty(propCode) {
			t.Errorf(errorMandatoryPropertyNotFound, propCode)
		}
	}
}
