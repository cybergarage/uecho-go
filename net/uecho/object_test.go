// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"testing"
)

const (
	errorMandatoryPropertyNotFound = "Mandatory Property (%0X) Not Found"
	errorInvalidGroupClassCode     = "Invalid Group Class Code (%0X)"
)

func TestNewObject(t *testing.T) {
	NewObject()
}
