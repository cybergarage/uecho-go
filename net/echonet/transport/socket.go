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
	ReadBufferSize int
	ReadBuffer     []byte
	BoundInterface *net.Interface
	BoundPort      int
	BoundAddress   string
}

// NewSocket returns a new UDPSocket.
func NewSocket() *Socket {
	sock := &Socket{
		ReadBufferSize: MaxPacketSize,
		ReadBuffer:     make([]byte, 0),
	}
	sock.SetReadBufferSize(MaxPacketSize)
	sock.Close()
	return sock
}

// Close initialize this socket.
func (sock *Socket) Close() {
	sock.BoundInterface = nil
	sock.BoundAddress = ""
	sock.BoundPort = 0
}

// SetReadBufferSize sets the read buffer size.
func (sock *Socket) SetReadBufferSize(n int) {
	sock.ReadBufferSize = n
	sock.ReadBuffer = make([]byte, n)
}

// GetReadBufferSize returns the read buffer size.
func (sock *Socket) GetReadBufferSize() int {
	return sock.ReadBufferSize
}

// SetBoundStatus sets the bound interface, port, and address.
func (sock *Socket) SetBoundStatus(i *net.Interface, addr string, port int) {
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
func (sock *Socket) GetBoundInterface() (*net.Interface, error) {
	if !sock.IsBound() {
		return nil, fmt.Errorf(errorSocketClosed)
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
