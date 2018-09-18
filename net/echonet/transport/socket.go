// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"
)

// A Socket represents a socket.
type Socket struct {
	Interface net.Interface
	Port      int
}

// NewSocket returns a new UDPSocket.
func NewSocket() *Socket {
	sock := &Socket{
		Port: 0,
	}
	return sock
}

// IsBound returns truen whether the socket is bound, otherwise false.
func (sock *Socket) IsBound() bool {
	if sock.Port <= 0 {
		return false
	}
	return true
}
