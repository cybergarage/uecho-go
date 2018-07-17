// Copyright 2018 The uecho-go Authors. All rights reserved.
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
	sock.readBuf = make([]byte, MaxPacketSize)
	return sock
}

// GetFD returns the file descriptor.
func (sock *UDPSocket) GetFD() (uintptr, error) {
	f, err := sock.Conn.File()
	if err != nil {
		return 0, err
	}
	return f.Fd(), nil
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

// ReadBytes reads received bytes from the current opened socket.
func (sock *UDPSocket) ReadBytes() ([]byte, error) {
	if sock.Conn == nil {
		return nil, errors.New(errorSocketIsClosed)
	}

	n, _, err := sock.Conn.ReadFromUDP(sock.readBuf)
	if err != nil {
		return nil, err
	}

	return sock.readBuf[:n], nil
}

// ReadMessage reads a message from the current opened socket.
func (sock *UDPSocket) ReadMessage() (*protocol.Message, error) {
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
