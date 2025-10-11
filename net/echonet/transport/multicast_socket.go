// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows
// +build !windows

package transport

import (
	"fmt"
	"net"
)

// A MulticastSocket represents a socket.
type MulticastSocket struct {
	*UDPSocket
}

// NewMulticastSocket returns a new MulticastSocket.
func NewMulticastSocket() *MulticastSocket {
	sock := &MulticastSocket{
		UDPSocket: NewUDPSocket(),
	}
	return sock
}

// Bind binds to the Echonet multicast address with the specified interface.
func (sock *MulticastSocket) Bind(ifi *net.Interface, ifaddr string) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	switch {
	case IsIPv4Address(ifaddr):
		err = sock.Listen(ifi, MulticastIPv4Address, Port)
	case IsIPv6Address(ifaddr):
		err = sock.Listen(ifi, MulticastIPv6Address, Port)
	default:
		return errAvailableAddressNotFound
	}

	if err != nil {
		return fmt.Errorf("%w (%s)", err, ifi.Name)
	}

	sock.SetBoundStatus(ifi, ifaddr, Port)
	sock.Conn.SetReadBuffer(sock.ReadBufferSize())

	rawConn, err := sock.Conn.SyscallConn()
	if err != nil {
		sock.Close()
		return err
	}

	err = sock.SetReuseAddr(rawConn, true)
	if err != nil {
		return err
	}

	err = sock.SetMulticastLoop(rawConn, ifaddr, true)
	if err != nil {
		return err
	}

	return nil
}
