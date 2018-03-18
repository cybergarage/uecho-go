// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"errors"
	"net"
)

const (
	errorSocketIsClosed = "Socket is closed"
)

// UDPServer represents a UDP server.
type UDPServer struct {
	Conn      *net.UDPConn
	readBuf   []byte
	Interface net.Interface
}

// newUDPServer returns a new UDPServer.
func newUDPServer() *UDPServer {
	uppSock := &UDPServer{}
	uppSock.readBuf = make([]byte, uechoMaxPacketSize)
	return uppSock
}

// Close closes the current socket.
func (server *UDPServer) Close() error {
	if server.Conn == nil {
		return nil
	}
	err := server.Conn.Close()
	if err != nil {
		return err
	}

	server.Conn = nil
	server.Interface = net.Interface{}

	return nil
}

// Read reads data from the current socket.
func (server *UDPServer) Read() (*Message, error) {
	if server.Conn == nil {
		return nil, errors.New(errorSocketIsClosed)
	}

	n, from, err := server.Conn.ReadFromUDP(server.readBuf)
	if err != nil {
		return nil, err
	}

	msg := NewMessage()
	err = msg.Parse(server.readBuf[:n])
	if err != nil {
		return nil, err
	}

	msg.From = *from
	msg.Interface = server.Interface

	return msg, nil
}
