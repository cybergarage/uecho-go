// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows
// +build !windows

package transport

import (
	"syscall"
)

// SetReuseAddr sets a flag to SO_REUSEADDR and SO_REUSEPORT.
// nolint: nosnakecase
func (sock *Socket) SetReuseAddr(fd uintptr, flag bool) error {
	opt := 0
	if flag {
		opt = 1
	}

	err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, opt)
	if err != nil {
		return err
	}

	return nil
}
