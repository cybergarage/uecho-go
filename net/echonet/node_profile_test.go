// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"testing"
)

func TestNewLocalNodeProfile(t *testing.T) {
	prof := NewLocalNodeProfile()

	if !prof.IsProfile() {
		t.Errorf(errorInvalidGroupClassCode, prof.ClassGroupCode())
	}

	mandatoryPropertyCodes := []PropertyCode{
		// Profile
		ProfileManufacturerCode,
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
		t.Run(fmt.Sprintf("%02X", propCode), func(t *testing.T) {
			if !prof.HasProperty(propCode) {
				t.Errorf(errorMandatoryPropertyNotFound, propCode)
			}
		})
	}

	testObjectPropertyMaps(t, prof.SuperObject.Object)
}
