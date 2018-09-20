// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"strconv"
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

// IsBound returns true whether the socket is bound, otherwise false.
func (sock *Socket) IsBound() bool {
	if sock.Port <= 0 {
		return false
	}
	return true
}

// GetBoundPort returns the bound port
func (sock *Socket) GetBoundPort() (int, error) {
	if !sock.IsBound() {
		return 0, fmt.Errorf(errorSocketIsClosed)
	}
	return sock.Port, nil
}

// GetBoundAddr returns the bound address
func (sock *Socket) GetBoundAddr() (string, error) {
	if !sock.IsBound() {
		return "", fmt.Errorf(errorSocketIsClosed)
	}

	addr, err := GetInterfaceAddress(sock.Interface)
	if err != nil {
		return "", err
	}

	return addr, nil
}

// GetBoundIPAddr returns the bound address
func (sock *Socket) GetBoundIPAddr() (string, error) {
	port, err := sock.GetBoundPort()
	if err != nil {
		return "", err
	}

	addr, err := sock.GetBoundAddr()
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(addr, strconv.Itoa(port)), nil
}
