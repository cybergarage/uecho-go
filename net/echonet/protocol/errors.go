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
	errInvalidMessageSize   = "%w: message length : %d < %d"
	errInvalidMessageHeader = "%w: message header [%d] : %02X != %02X"
	errInvalidAddress       = "%w: address string : %s"
	errInvalidObjectCodes   = "%w: object code : %s"
)
