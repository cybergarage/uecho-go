// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux
// +build !bsd
// +build !freebsd
// +build !darwin
// +build !windows

package transport

import (
	"os"
)

// SetReuseAddr sets a flag to SO_REUSEADDR and SO_REUSEPORT
func (sock *Socket) SetReuseAddr(file *os.File, flag bool) error {
	return nil
}
