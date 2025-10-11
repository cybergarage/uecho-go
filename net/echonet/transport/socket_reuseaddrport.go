// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (bsd || freebsd || darwin) && !windows
// +build bsd freebsd darwin
// +build !windows

package transport

import (
	"syscall"
)

// SetReuseAddr sets a flag to SO_REUSEADDR and SO_REUSEPORT.
// nolint: nosnakecase
func (sock *Socket) SetReuseAddr(rawConn syscall.RawConn, flag bool) error {
	opt := 0
	if flag {
		opt = 1
	}

	err := rawConn.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, opt)
	})
	if err != nil {
		return err
	}

	err = rawConn.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, opt)
	})

	if err != nil {
		return err
	}

	return nil
}
