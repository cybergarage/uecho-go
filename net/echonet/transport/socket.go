// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
)

// A Socket represents a socket.
type Socket struct {
	BoundInterface net.Interface
	BoundPort      int
	BoundAddress   string
}

// NewSocket returns a new UDPSocket.
func NewSocket() *Socket {
	sock := &Socket{}
	sock.Close()
	return sock
}

// Close initialize this socket.
func (sock *Socket) Close() {
	sock.BoundInterface = net.Interface{}
	sock.BoundAddress = ""
	sock.BoundPort = 0
}

// SetBoundStatus sets the bound interface, port, and address.
func (sock *Socket) SetBoundStatus(i net.Interface, addr string, port int) {
	sock.BoundInterface = i
	sock.BoundAddress = addr
	sock.BoundPort = port
}

// IsBound returns true whether the socket is bound, otherwise false.
func (sock *Socket) IsBound() bool {
	if sock.BoundPort <= 0 {
		return false
	}
	return true
}

// GetBoundPort returns the bound port.
func (sock *Socket) GetBoundPort() (int, error) {
	if !sock.IsBound() {
		return 0, fmt.Errorf(errorSocketClosed)
	}
	return sock.BoundPort, nil
}

// GetBoundInterface returns the bound interface.
func (sock *Socket) GetBoundInterface() (net.Interface, error) {
	if !sock.IsBound() {
		return net.Interface{}, fmt.Errorf(errorSocketClosed)
	}
	return sock.BoundInterface, nil
}

// GetBoundAddr returns the bound address
func (sock *Socket) GetBoundAddr() (string, error) {
	if !sock.IsBound() {
		return "", fmt.Errorf(errorSocketClosed)
	}

	return sock.BoundAddress, nil
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

// SetReuseAddr sets a flag to SO_REUSEADDR and SO_REUSEPORT
func (sock *Socket) SetReuseAddr(file *os.File, flag bool) error {
	fd := file.Fd()

	opt := 0
	if flag {
		opt = 1
	}

	err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, opt)
	if err != nil {
		return err
	}

	// Disable for Linux platrorms
	//_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, opt)
	//if err != nil {
	//	return err
	//}

	return nil
}
