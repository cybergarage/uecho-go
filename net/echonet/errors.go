// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"errors"
)

// ErrInvalid is returned when the value is invalid.
var ErrInvalid = errors.New("invalid")

// ErrNoData is returned when there is no data.
var ErrNoData = errors.New("no data")

// ErrNotFound is returned when the value is not found.
var ErrNotFound = errors.New("not found")

// ErrUnknown is returned when the value is unknown.
var ErrUnknown = errors.New("unknown")

// ErrTimeout is returned when the operation times out.
var ErrTimeout = errors.New("timeout")
