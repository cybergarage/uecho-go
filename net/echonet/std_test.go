// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"testing"
)

func TestStandardDatabase(t *testing.T) {
	objCodes := []ObjectCode{
		0x000000,
		0x0EF000,
	}

	db := SharedStandardDatabase()
	for _, objCode := range objCodes {
		t.Run(fmt.Sprintf("%06X", objCode), func(t *testing.T) {
			_, ok := db.LookupObjectByCode(objCode)
			if !ok {
				t.Errorf("%06X object is not found", objCode)
			}
		})
	}
}
