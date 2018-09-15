// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"github.com/cybergarage/uecho-go/net/echonet/log"
)

// A UnicastTCPSocket represents a socket.
type UnicastTCPSocket struct {
	*TCPSocket
}

// NewUnicastTCPSocket returns a new UnicastTCPSocket.
func NewUnicastTCPSocket() *UnicastTCPSocket {
	sock := &UnicastTCPSocket{
		TCPSocket: NewTCPSocket(),
	}
	return sock
}

// Bind binds to Echonet multicast address.
func (sock *UnicastTCPSocket) Bind(ifi net.Interface, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr, err := GetInterfaceAddress(ifi)
	if err != nil {
		return err
	}

	boundAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	_, err = net.ListenTCP("tcp", boundAddr)
	if err != nil {
		return err
	}

	sock.Interface = ifi

	return nil
}

// Write sends the specified bytes.
func (sock *UnicastTCPSocket) Write(addr string, port int, b []byte) (int, error) {
	toAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return 0, err
	}

	// Send from no binding port

	conn, err := net.DialTCP("tcp", nil, toAddr)
	if err != nil {
		return 0, err
	}

	n, err := conn.Write(b)
	log.Trace(fmt.Sprintf(logSocketWriteFormat, conn.LocalAddr().String(), toAddr.String(), n, hex.EncodeToString(b)))
	conn.Close()

	return n, err
}
