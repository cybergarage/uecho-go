// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package transport

import (
	"errors"
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
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
		return errors.New(errorAvailableAddressNotFound)
	}

	if err != nil {
		return fmt.Errorf("%w (%s)", err, ifi.Name)
	}

	sock.SetBoundStatus(ifi, ifaddr, Port)
	sock.Conn.SetReadBuffer(sock.ReadBufferSize())

	rawConn, err := sock.Conn.SyscallConn()
	if err != nil {
		return err
	}
	fdCh := make(chan uintptr, 1)
	err = rawConn.Control(func(fd uintptr) {
		fdCh <- fd
	})
	if err != nil {
		return err
	}
	fd := <-fdCh
	err = sock.SetReuseAddr(fd, true)
	if err != nil {
		return err
	}

	pc := ipv4.NewPacketConn(sock.Conn)

	if err := pc.JoinGroup(ifi, &net.UDPAddr{IP: net.ParseIP(MulticastIPv4Address), Port: Port}); err != nil {
		return err
	}

	if loop, err := pc.MulticastLoopback(); err == nil {
		if !loop {
			if err := pc.SetMulticastLoopback(true); err != nil {
				return err
			}
		}
	}

	return nil
}
