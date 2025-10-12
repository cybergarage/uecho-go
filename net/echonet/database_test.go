// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"testing"
)

func TestStandardDatabase(t *testing.T) {
	db := SharedStandardDatabase()

	// Check manufactures

	manCodes := []ManufactureCode{
		0x000005,
	}

	for _, manCode := range manCodes {
		t.Run(fmt.Sprintf("%06X", manCode), func(t *testing.T) {
			_, ok := db.LookupManufacture(manCode)
			if !ok {
				t.Errorf("%06X manufacture is not found", manCode)
			}
		})
	}

	// Check some standard objects

	objCodes := []ObjectCode{
		0x000000,
		0x0EF000,
		0x029101, // Mono functional lighting
	}

	superObj := db.SuperObject()
	if superObj == nil {
		t.Error("Super object is nil")
	}

	nodeProf := db.NodeProfile()
	if nodeProf == nil {
		t.Error("Node profile object is nil")
	}
	for _, objCode := range objCodes {
		t.Run(fmt.Sprintf("%06X", objCode), func(t *testing.T) {
			_, ok := db.LookupObject(objCode)
			if !ok {
				t.Errorf("%06X object is not found", objCode)
			}
		})
	}
}
