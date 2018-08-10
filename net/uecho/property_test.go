// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"testing"
)

func TestNewProperty(t *testing.T) {
	NewProperty()
	NewPropertyWithCode(0)
}
