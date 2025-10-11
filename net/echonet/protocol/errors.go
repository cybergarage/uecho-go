// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"errors"
)

// ErrNoData is returned when there is no data.
var ErrNoData = errors.New("no data")

// ErrInvalid is returned when the value is invalid.
var ErrInvalid = errors.New("invalid")

const (
	errorInvalidMessageSize   = "%w: message length : %d < %d"
	errorInvalidMessageHeader = "%w: message header [%d] : %02X != %02X"
	errorInvalidAddress       = "%w: address string : %s"
	errorInvalidObjectCodes   = "%w: object code : %s"
)
