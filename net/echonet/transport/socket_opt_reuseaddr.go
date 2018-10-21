// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux windows

package transport

import (
	"os"
	"syscall"
)

// SetReuseAddr sets a flag to SO_REUSEADDR and SO_REUSEPORT
func (sock *Socket) SetReuseAddr(file *os.File, flag bool) error {
	fd := file.Fd()

	opt := 0
	if flag {
		opt = 1
	}

	err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, opt)
	if err != nil {
		return err
	}

	return nil
}
