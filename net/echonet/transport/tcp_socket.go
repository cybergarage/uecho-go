// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A TCPSocket represents a socket for TCP.
type TCPSocket struct {
	*Socket
	Listener *net.TCPListener
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

	sock.Listener, err = net.ListenTCP("tcp", boundAddr)
	if err != nil {
		return err
	}

	f, err := sock.Listener.File()
	if err != nil {
		return err
	}
	err = sock.SetReuseAddr(f, true)
	if err != nil {
		return err
	}

	sock.Port = port
	sock.Interface = ifi

	return nil
}

// Close closes the current opened socket.
func (sock *TCPSocket) Close() error {
	if sock.Listener == nil {
		return nil
	}

	err := sock.Listener.Close()
	if err != nil {
		return err
	}

	sock.Listener = nil
	sock.Port = 0
	sock.Interface = net.Interface{}

	return nil
}

func (sock *TCPSocket) outputReadLog(logLevel log.LogLevel, msgFrom string, msg string, msgSize int) {
	if sock.Listener == nil {
		return
	}
	outputSocketLog(logLevel, logSocketTypeTCP, logSocketDirectionRead, msgFrom, sock.Listener.Addr().String(), msg, msgSize)
}

// ReadMessage reads a message from the current opened socket.
func (sock *TCPSocket) ReadMessage(clientConn net.Conn) (*protocol.Message, error) {
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
	msgFrom, _ := sock.GetBoundIPAddr()
	outputSocketLog(logLevel, logSocketTypeTCP, logSocketDirectionWrite, msgFrom, msgTo, msg, msgSize)
}

// Write sends the specified bytes.
func (sock *TCPSocket) Write(addr string, port int, b []byte, timeout time.Duration) (int, error) {
	toAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		sock.outputWriteLog(log.LoggerLevelError, toAddr.String(), hex.EncodeToString(b), 0)
		return 0, err
	}

	// Send from binding port

	boundAddr, err := sock.GetBoundIPAddr()
	if err == nil {
		fromAddr, err := net.ResolveTCPAddr("tcp", boundAddr)
		if err == nil {
			conn, err := net.DialTCP("tcp", fromAddr, toAddr)
			if err == nil {
				n, err := conn.Write(b)
				sock.outputWriteLog(log.LoggerLevelTrace, toAddr.String(), hex.EncodeToString(b), n)
				conn.Close()
				if err == nil {
					return n, nil
				}
			} else {
				log.Error(fmt.Sprintf("%s", err))
			}
		}
	}

	// Send from no binding port

	conn, err := net.DialTCP("tcp", nil, toAddr)
	if err != nil {
		sock.outputWriteLog(log.LoggerLevelError, toAddr.String(), hex.EncodeToString(b), 0)
		return 0, err
	}

	n, err := conn.Write(b)
	sock.outputWriteLog(log.LoggerLevelTrace, toAddr.String(), hex.EncodeToString(b), n)
	conn.Close()

	return n, err
}
