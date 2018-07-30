// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"syscall"
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

	defer conn.Close()

	return conn.Write(b)
}
