// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
func (sock *MulticastSocket) Bind(ifi *net.Interface) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	err = sock.Listen(ifi)
	if err != nil {
		return fmt.Errorf("%s (%s)", err.Error(), ifi.Name)
	}

	sock.Conn.SetReadBuffer(sock.GetReadBufferSize())

	f, err := sock.Conn.File()
	if err != nil {
		return err
	}

	defer f.Close()

	err = sock.SetReuseAddr(f, true)
	if err != nil {
		return err
	}

	err = sock.SetMulticastLoop(f, true)
	if err != nil {
		return err
	}

	sock.SetBoundStatus(ifi, MulticastAddress, UDPPort)

	return nil
}
