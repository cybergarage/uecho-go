// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A TCPSocket represents a socket for TCP.
type TCPSocket struct {
	*Socket
	Conn    *net.TCPConn
	readBuf []byte
}

// NewTCPSocket returns a new TCPSocket.
func NewTCPSocket() *TCPSocket {
	sock := &TCPSocket{
		Socket:  NewSocket(),
		readBuf: make([]byte, MaxPacketSize),
	}
	return sock
}

// GetFD returns the file descriptor.
func (sock *TCPSocket) GetFD() (uintptr, error) {
	f, err := sock.Conn.File()
	if err != nil {
		return 0, err
	}
	return f.Fd(), nil

}

// Close closes the current opened socket.
func (sock *TCPSocket) Close() error {
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

// ReadMessage reads a message from the current opened socket.
func (sock *TCPSocket) ReadMessage() (*protocol.Message, error) {
	if sock.Conn == nil {
		return nil, errors.New(errorSocketIsClosed)
	}

	retemoAddr := sock.Conn.RemoteAddr()

	reader := bufio.NewReader(sock.Conn)
	msg, err := protocol.NewMessageWithReader(reader)
	if err != nil {
		if sock.Conn != nil {
			log.Error(fmt.Sprintf(logSocketReadFormat, sock.Conn.LocalAddr().String(), retemoAddr, 0, ""))
		}
		return nil, err
	}

	err = msg.From.ParseString(retemoAddr.String())
	if err != nil {
		return nil, err
	}

	if msg != nil && sock.Conn != nil {
		log.Trace(fmt.Sprintf(logSocketReadFormat, sock.Conn.LocalAddr().String(), retemoAddr, msg.Size(), msg.String()))
	}

	return msg, nil
}
