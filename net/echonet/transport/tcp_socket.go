// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bufio"
	"encoding/hex"
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

	sock.SetBoundStatus(ifi, addr, port)

	return nil
}

// Close closes the current opened socket.
func (sock *TCPSocket) Close() error {
	if sock.Listener == nil {
		return nil
	}

	// FIXE : Hung up on go1.11 darwin/amd64
	/*
		err := sock.Listener.Close()
		if err != nil {
			return err
		}
	*/

	return nil
}

func (sock *TCPSocket) outputReadLog(logLevel log.LogLevel, msgFrom string, msg string, msgSize int) {
	msgTo, _ := sock.GetBoundIPAddr()
	outputSocketLog(logLevel, logSocketTypeTCPUnicast, logSocketDirectionRead, msgFrom, msgTo, msg, msgSize)
}

// ReadMessage reads a message from the current opened socket.
func (sock *TCPSocket) ReadMessage(clientConn net.Conn) (*protocol.Message, error) {
	retemoAddr := clientConn.RemoteAddr()

	reader := bufio.NewReader(clientConn)
	msg, err := protocol.NewMessageWithReader(reader)
	if err != nil {
		sock.outputReadLog(log.LoggerLevelError, retemoAddr.String(), "", 0)
		log.Error(err.Error())
		return nil, err
	}

	err = msg.From.ParseString(retemoAddr.String())
	if err != nil {
		sock.outputReadLog(log.LoggerLevelError, retemoAddr.String(), msg.String(), msg.Size())
		log.Error(err.Error())
		return nil, err
	}

	sock.outputReadLog(log.LoggerLevelTrace, retemoAddr.String(), msg.String(), msg.Size())

	return msg, nil
}

func (sock *TCPSocket) outputWriteLog(logLevel log.LogLevel, msgFrom string, msgTo string, msg string, msgSize int) {
	outputSocketLog(logLevel, logSocketTypeTCPUnicast, logSocketDirectionWrite, msgFrom, msgTo, msg, msgSize)
}

// Write sends the specified bytes.
func (sock *TCPSocket) Write(addr string, port int, b []byte, timeout time.Duration) (int, error) {
	toAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		sock.outputWriteLog(log.LoggerLevelError, "", toAddr.String(), hex.EncodeToString(b), 0)
		log.Error(err.Error())
		return 0, err
	}

	// Send from binding or any port

	/* Disable to send from listen port
	boundAddr, err := net.ResolveTCPAddr("tcp", sock.Listener.Addr().String())
	if err != nil {
		sock.outputWriteLog(log.LoggerLevelError, sock.Listener.Addr().String(), toAddr.String(), hex.EncodeToString(b), 0)
		log.Error(err.Error())
		return 0, err
	}
	fromAddrs := []*net.TCPAddr{boundAddr, nil}
	*/

	var lastError error

	fromAddrs := []*net.TCPAddr{nil}

	for _, fromAddr := range fromAddrs {
		var conn *net.TCPConn
		conn, lastError = net.DialTCP("tcp", fromAddr, toAddr)
		if lastError != nil {
			sock.outputWriteLog(log.LoggerLevelError, fromAddr.String(), toAddr.String(), hex.EncodeToString(b), 0)
			log.Error(lastError.Error())
			continue
		}

		localAddr := conn.LocalAddr()

		var nWrote int
		nWrote, lastError = conn.Write(b)
		if lastError != nil {
			sock.outputWriteLog(log.LoggerLevelError, localAddr.String(), toAddr.String(), hex.EncodeToString(b), 0)
			log.Error(lastError.Error())
			continue
		}

		sock.outputWriteLog(log.LoggerLevelTrace, localAddr.String(), toAddr.String(), hex.EncodeToString(b), nWrote)
		conn.Close()

		return nWrote, nil
	}

	return 0, lastError
}
