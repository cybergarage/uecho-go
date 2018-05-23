// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
)

// A UnicastSocket represents a socket.
type UnicastSocket struct {
	*UDPSocket
}

// NewUnicastSocket returns a new UnicastSocket.
func NewUnicastSocket() *UnicastSocket {
	sock := &UnicastSocket{}
	sock.UDPSocket = NewUDPSocket()
	return sock
}

// Bind binds to Echonet multicast address.
func (sock *UnicastSocket) Bind(ifi net.Interface, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr, err := GetInterfaceAddress(ifi)
	if err != nil {
		return err
	}

	bindAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp", bindAddr)
	if err != nil {
		return err
	}

	sock.Interface = ifi

	return nil
}

// Write sends the specified bytes.
func (sock *UnicastSocket) Write(addr string, port int, b []byte) (int, error) {
	toAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return 0, err
	}

	if sock.Conn != nil {
		return sock.Conn.WriteToUDP(b, toAddr)
	}

	conn, err := net.DialUDP("udp", nil, toAddr)
	if err != nil {
		return 0, err
	}

	return conn.Write(b)
}
