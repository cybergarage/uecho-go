// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"encoding/hex"
	"errors"
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A UDPSocket represents a socket for UDP.
type UDPSocket struct {
	*Socket
	Conn    *net.UDPConn
	readBuf []byte
}

// NewUDPSocket returns a new UDPSocket.
func NewUDPSocket() *UDPSocket {
	sock := &UDPSocket{
		Socket:  NewSocket(),
		readBuf: make([]byte, MaxPacketSize),
	}
	return sock
}

// Close closes the current opened socket.
func (sock *UDPSocket) Close() error {
	if sock.Conn == nil {
		return nil
	}

	// FIXE : Hung up on go1.11 darwin/amd64
	/*
		err := sock.Conn.Close()
		if err != nil {
			return err
		}
	*/

	sock.Conn = nil
	sock.Port = 0
	sock.Interface = net.Interface{}

	return nil
}

func (sock *UDPSocket) outputReadLog(logLevel log.LogLevel, msgFrom string, msg string, msgSize int) {
	if sock.Conn == nil {
		return
	}
	outputSocketLog(logLevel, logSocketTypeUDP, logSocketDirectionRead, msgFrom, sock.Conn.LocalAddr().String(), msg, msgSize)
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
		sock.outputReadLog(log.LoggerLevelError, (*from).String(), hex.EncodeToString(sock.readBuf[:n]), n)
		return nil, err
	}

	msg.From.IP = (*from).IP
	msg.From.Port = (*from).Port
	msg.Interface = sock.Interface

	sock.outputReadLog(log.LoggerLevelTrace, msg.From.String(), msg.String(), msg.Size())

	return msg, nil
}
