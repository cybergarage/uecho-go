// Copyright 2018 The uecho-go Authors. All rights reserved.
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
	Conn           *net.UDPConn
	ReadBufferSize int
	ReadBuffer     []byte
}

// NewUDPSocket returns a new UDPSocket.
func NewUDPSocket() *UDPSocket {
	sock := &UDPSocket{
		Socket:         NewSocket(),
		ReadBufferSize: MaxPacketSize,
		ReadBuffer:     make([]byte, 0),
	}
	sock.SetReadBufferSize(MaxPacketSize)
	return sock
}

// SetReadBufferSize sets the read buffer size.
func (sock *UDPSocket) SetReadBufferSize(n int) {
	sock.ReadBufferSize = n
	sock.ReadBuffer = make([]byte, n)
}

// GetReadBufferSize returns the read buffer size.
func (sock *UDPSocket) GetReadBufferSize() int {
	return sock.ReadBufferSize
}

// Close closes the current opened socket.
func (sock *UDPSocket) Close() error {
	if sock.Conn == nil {
		return nil
	}

	err := sock.Conn.Close()
	if err != nil {
		//return err
	}

	sock.Socket.Close()

	return nil
}

func (sock *UDPSocket) outputReadLog(logLevel log.LogLevel, logType string, msgFrom string, msg string, msgSize int) {
	msgTo, _ := sock.GetBoundIPAddr()
	outputSocketLog(logLevel, logType, logSocketDirectionRead, msgFrom, msgTo, msg, msgSize)
}

// ReadMessage reads a message from the current opened socket.
func (sock *UDPSocket) ReadMessage() (*protocol.Message, error) {
	if sock.Conn == nil {
		return nil, errors.New(errorSocketClosed)
	}

	n, from, err := sock.Conn.ReadFromUDP(sock.ReadBuffer)
	if err != nil {
		return nil, err
	}

	msg, err := protocol.NewMessageWithBytes(sock.ReadBuffer[:n])
	if err != nil {
		sock.outputReadLog(log.LevelError, logSocketTypeUDPUnicast, (*from).String(), hex.EncodeToString(sock.ReadBuffer[:n]), n)
		log.Error(err.Error())
		return nil, err
	}

	msg.From.IP = (*from).IP
	msg.From.Port = (*from).Port
	msg.Interface = sock.Socket.BoundInterface

	return msg, nil
}
