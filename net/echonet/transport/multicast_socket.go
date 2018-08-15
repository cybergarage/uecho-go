// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"errors"
	"fmt"
	"net"
	"syscall"
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

// Bind binds to Echonet multicast address.
func (sock *MulticastSocket) Bind(ifi net.Interface) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", MulticastAddress, UDPPort))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenMulticastUDP("udp", &ifi, addr)
	if err != nil {
		return fmt.Errorf("%s (%s)", err.Error(), ifi.Name)
	}

	fd, err := sock.GetFD()
	if err != nil {
		return err
	}

	err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return err
	}

	_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
	if err != nil {
		return err
	}

	sock.Interface = ifi

	return nil
}

// Write sends the specified bytes.
func (sock *MulticastSocket) Write(b []byte) (int, error) {
	if sock.Conn == nil {
		return 0, errors.New(errorSocketIsClosed)
	}

	addr, err := net.ResolveUDPAddr("udp", MulticastAddress)
	if err != nil {
		return 0, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return conn.Write(b)
}
