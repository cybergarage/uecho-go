// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"errors"
	"fmt"
)

// ErrInvalid is returned when the value is invalid.
var ErrInvalid = errors.New("invalid")

var (
	errSocketClosed             = fmt.Errorf("%w: socket is closed", ErrInvalid)
	errTCPSocketDisabled        = fmt.Errorf("%w: TCP function is disabled", ErrInvalid)
	errAvailableAddressNotFound = fmt.Errorf("%w: no available address", ErrInvalid)
	errAvailableInterfaceFound  = fmt.Errorf("%w: no available interface", ErrInvalid)
	errUnicastServerNotRunning  = fmt.Errorf("%w: unicast server is not running", ErrInvalid)
)
