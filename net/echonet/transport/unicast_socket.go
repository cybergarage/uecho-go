// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"syscall"

	"github.com/cybergarage/uecho-go/net/echonet/log"
)

// A UnicastSocket represents a socket.
type UnicastSocket struct {
	*UDPSocket
}

// NewUnicastSocket returns a new UnicastSocket.
func NewUnicastSocket() *UnicastSocket {
	sock := &UnicastSocket{
		UDPSocket: NewUDPSocket(),
	}
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

	boundAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp", boundAddr)
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

	// Disable for Linux platrorms
	//_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
	//if err != nil {
	//	return err
	//}

	sock.Interface = ifi

	return nil
}

// Write sends the specified bytes.
func (sock *UnicastSocket) Write(addr string, port int, b []byte) (int, error) {
	toAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return 0, err
	}

	// Send from binding port

	if sock.Conn != nil {
		n, err := sock.Conn.WriteToUDP(b, toAddr)
		log.Trace(fmt.Sprintf(logSocketWriteFormat, sock.Conn.LocalAddr().String(), toAddr.String(), n, hex.EncodeToString(b)))
		return n, err
	}

	// Send from no binding port

	conn, err := net.DialUDP("udp", nil, toAddr)
	if err != nil {
		return 0, err
	}

	n, err := conn.Write(b)
	log.Trace(fmt.Sprintf(logSocketWriteFormat, conn.LocalAddr().String(), toAddr.String(), n, hex.EncodeToString(b)))
	conn.Close()

	return n, err
}
