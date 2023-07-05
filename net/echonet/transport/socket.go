// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"strconv"
	"syscall"
)

// A Socket represents a socket.
type Socket struct {
	interfac *net.Interface
	port     int
	address  string
}

// NewSocket returns a new UDPSocket.
func NewSocket() *Socket {
	sock := &Socket{
		interfac: nil,
		port:     0,
		address:  "",
	}
	sock.Close()
	return sock
}

// Close initialize this socket.
func (sock *Socket) Close() {
	sock.interfac = nil
	sock.address = ""
	sock.port = 0
}

// SetBoundStatus sets the bound interface, port, and address.
func (sock *Socket) SetBoundStatus(i *net.Interface, addr string, port int) {
	sock.interfac = i
	sock.address = addr
	sock.port = port
}

// IsBound returns true whether the socket is bound, otherwise false.
func (sock *Socket) IsBound() bool {
	return sock.port != 0
}

// Port returns the bound port.
func (sock *Socket) Port() (int, error) {
	if !sock.IsBound() {
		return 0, fmt.Errorf(errorSocketClosed)
	}
	return sock.port, nil
}

// Interface returns the bound interface.
func (sock *Socket) Interface() (*net.Interface, error) {
	if !sock.IsBound() {
		return nil, fmt.Errorf(errorSocketClosed)
	}
	return sock.interfac, nil
}

// Address returns the bound address.
func (sock *Socket) Address() (string, error) {
	if !sock.IsBound() {
		return "", fmt.Errorf(errorSocketClosed)
	}
	return sock.address, nil
}

// IPAddr returns the bound address.
func (sock *Socket) IPAddr() (string, error) {
	port, err := sock.Port()
	if err != nil {
		return "", err
	}

	addr, err := sock.Address()
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(addr, strconv.Itoa(port)), nil
}

// SetMulticastLoop sets a flag to IP_MULTICAST_LOOP.
// nolint: nosnakecase
func (sock *Socket) SetMulticastLoop(fd uintptr, addr string, flag bool) error {
	opt := 0
	if flag {
		opt = 1
	}

	if IsIPv6Address(addr) {
		return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_MULTICAST_LOOP, opt)
	}
	return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MULTICAST_LOOP, opt)
}
