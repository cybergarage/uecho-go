// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bufio"
	"encoding/hex"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A TCPSocket represents a socket for TCP.
type TCPSocket struct {
	*Socket
	Conn     *net.TCPConn
	Listener net.Listener
	readBuf  []byte
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

// Bind binds to Echonet multicast address.
func (sock *TCPSocket) Bind(ifi net.Interface, port int) error {
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

	l, err := net.ListenTCP("tcp", boundAddr)
	if err != nil {
		return err
	}

	sock.Port = port
	sock.Listener = l
	sock.Interface = ifi

	return nil
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
	sock.Listener = nil
	sock.Port = 0
	sock.Interface = net.Interface{}

	return nil
}

func (sock *TCPSocket) outputReadLog(logLevel log.LogLevel, msgFrom string, msg string, msgSize int) {
	if sock.Conn == nil {
		return
	}
	outputSocketLog(logLevel, logSocketTypeTCP, logSocketDirectionRead, msgFrom, sock.Conn.LocalAddr().String(), msg, msgSize)
}

// ReadMessage reads a message from the current opened socket.
func (sock *TCPSocket) ReadMessage(clientConn net.Conn) (*protocol.Message, error) {
	if sock.Conn == nil {
		return nil, errors.New(errorSocketIsClosed)
	}

	retemoAddr := clientConn.RemoteAddr()

	reader := bufio.NewReader(clientConn)
	msg, err := protocol.NewMessageWithReader(reader)
	if err != nil {
		sock.outputReadLog(log.LoggerLevelError, retemoAddr.String(), "", 0)
		return nil, err
	}

	err = msg.From.ParseString(retemoAddr.String())
	if err != nil {
		return nil, err
	}

	sock.outputReadLog(log.LoggerLevelTrace, retemoAddr.String(), msg.String(), msg.Size())

	return msg, nil
}

func (sock *TCPSocket) outputWriteLog(logLevel log.LogLevel, msgTo string, msg string, msgSize int) {
	localAddr := ""
	if sock.Conn != nil {
		localAddr = sock.Conn.LocalAddr().String()
	}
	outputSocketLog(logLevel, logSocketTypeTCP, logSocketDirectionWrite, localAddr, msgTo, msg, msgSize)
}

// Write sends the specified bytes.
func (sock *TCPSocket) Write(addr string, port int, b []byte, timeout time.Duration) (int, error) {
	toAddr := net.JoinHostPort(addr, strconv.Itoa(port))
	/*
		toAddr, err := net.ResolveIPAddr("tcp", toAddr)
		if err != nil {
			sock.outputWriteLog(log.LoggerLevelError, toAddr, hex.EncodeToString(b), 0)
			return 0, err
		}
	*/

	// Send from no binding port

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", toAddr)
	//conn, err := net.DialTCP("tcp", nil, toAddr)
	if err != nil {
		sock.outputWriteLog(log.LoggerLevelError, toAddr, hex.EncodeToString(b), 0)
		return 0, err
	}

	n, err := conn.Write(b)
	sock.outputWriteLog(log.LoggerLevelTrace, toAddr, hex.EncodeToString(b), n)
	conn.Close()

	return n, err
}
