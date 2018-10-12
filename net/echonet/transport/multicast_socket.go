// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"strconv"
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
func (sock *MulticastSocket) Bind(ifi *net.Interface) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(MulticastAddress, strconv.Itoa(UDPPort)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		return fmt.Errorf("%s (%s)", err.Error(), ifi.Name)
	}

	f, err := sock.Conn.File()
	if err != nil {
		return err
	}
	err = sock.SetReuseAddr(f, true)
	if err != nil {
		return err
	}

	sock.SetBoundStatus(ifi, MulticastAddress, UDPPort)

	return nil
}
