// Copyright 2018 Satoshi Konno. All rights reserved.
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
func (sock *Socket) GetBoundAddr() (net.Addr, error) {
	if !sock.IsBound() {
		return nil, fmt.Errorf(errorSocketIsClosed)
	}

	addrs, err := sock.Interface.Addrs()
	if err != nil {
		return nil, err
	}

	if len(addrs) <= 0 {
		return nil, fmt.Errorf(errorSocketIsClosed)
	}

	return addrs[0], nil
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

	return net.JoinHostPort(addr.String(), strconv.Itoa(port)), nil
}
