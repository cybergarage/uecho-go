// Copyright 2017 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"errors"
	"net"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	errorSocketIsClosed = "Socket is closed"
)

// A UDPSocket represents a socket for UDP.
type UDPSocket struct {
	Conn      *net.UDPConn
	readBuf   []byte
	Interface net.Interface
}

// NewUDPSocket returns a new UDPSocket.
func NewUDPSocket() *UDPSocket {
	sock := &UDPSocket{}
	sock.readBuf = make([]byte, MAX_PACKET_SIZE)
	return sock
}

// Close closes the current opened socket.
func (sock *UDPSocket) Close() error {
	if sock.Conn == nil {
		return nil
	}
	err := sock.Conn.Close()
	if err != nil {
		return err
	}

	sock.Conn = nil
	sock.Interface = net.Interface{}

	return nil
}

// Read reads from the current opened socket.
func (sock *UDPSocket) Read() (*protocol.Message, error) {
	if sock.Conn == nil {
		return nil, errors.New(errorSocketIsClosed)
	}

	n, from, err := sock.Conn.ReadFromUDP(sock.readBuf)
	if err != nil {
		return nil, err
	}

	msg, err := protocol.NewMessageWithBytes(sock.readBuf[:n])
	if err != nil {
		return nil, err
	}

	msg.From = *from
	msg.Interface = sock.Interface

	return msg, nil
}
