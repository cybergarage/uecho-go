// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func TestGetAvailableInterfaces(t *testing.T) {
	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
	}
	if len(ifs) <= 0 {
		t.Errorf("available interface is not found")
	}
}
