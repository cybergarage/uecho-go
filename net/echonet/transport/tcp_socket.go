// Copyright 2018 The uecho-go Authors. All rights reserved.
// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows
// +build !windows

package transport

import (
	"bufio"
	"encoding/hex"
	"net"
	"strconv"
	"time"

	"github.com/cybergarage/go-logger/log"
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
		Socket:   NewSocket(),
		Listener: nil,
		readBuf:  make([]byte, MaxPacketSize),
	}
	return sock
}

// Bind binds to Echonet multicast address.
func (sock *TCPSocket) Bind(ifi *net.Interface, ifaddr string, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	boundAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(ifaddr, strconv.Itoa(port)))
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

	defer f.Close()
	fd := f.Fd()

	err = sock.SetReuseAddr(fd, true)
	if err != nil {
		return err
	}

	sock.SetBoundStatus(ifi, ifaddr, port)

	return nil
}

// Close closes the current opened socket.
func (sock *TCPSocket) Close() error {
	if sock.Listener == nil {
		return nil
	}

	// FIXME : sock.Listener.Close() hung up on darwin
	/*
		err := sock.Listener.Close()
		if err != nil {
			return err
		}
	*/
	go sock.Listener.Close()
	time.Sleep(time.Millisecond * 100)

	sock.Listener = nil

	return nil
}

func (sock *TCPSocket) outputReadLog(logLevel log.Level, msgFrom string, msg string, msgSize int) {
	msgTo, _ := sock.IPAddr()
	outputSocketLog(logLevel, logSocketTypeTCPUnicast, logSocketDirectionRead, msgFrom, msgTo, msg, msgSize)
}

// ReadMessage reads a message from the current opened socket.
func (sock *TCPSocket) ReadMessage(conn net.Conn) (*protocol.Message, error) {
	retemoAddr := conn.RemoteAddr()

	reader := bufio.NewReader(conn)
	msg, err := protocol.NewMessageWithReader(reader)
	if err != nil {
		sock.outputReadLog(log.LevelError, retemoAddr.String(), "", 0)
		log.Error(err)
		return nil, err
	}

	err = msg.From.ParseString(retemoAddr.String())
	if err != nil {
		sock.outputReadLog(log.LevelError, retemoAddr.String(), msg.String(), msg.Size())
		log.Error(err)
		return nil, err
	}

	sock.outputReadLog(log.LevelTrace, retemoAddr.String(), msg.String(), msg.Size())

	return msg, nil
}

func (sock *TCPSocket) outputWriteLog(logLevel log.Level, msgFrom string, msgTo string, msg string, msgSize int) {
	outputSocketLog(logLevel, logSocketTypeTCPUnicast, logSocketDirectionWrite, msgFrom, msgTo, msg, msgSize)
}

// SendMessage sends a message to the destination address.
func (sock *TCPSocket) SendMessage(addr string, port int, msg *protocol.Message, timeout time.Duration) (int, error) {
	conn, nWorte, err := sock.dialAndWriteBytes(addr, port, msg.Bytes(), timeout)
	if conn != nil {
		conn.Close()
	}
	return nWorte, err
}

// PostMessage sends a message to the destination address.
func (sock *TCPSocket) PostMessage(addr string, port int, reqMsg *protocol.Message, timeout time.Duration) (*protocol.Message, error) {
	conn, _, err := sock.dialAndWriteBytes(addr, port, reqMsg.Bytes(), timeout)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		conn.Close()
		return nil, err
	}

	return sock.ReadMessage(conn)
}

// ResponseMessageForRequestMessage sends a specified response message to the request node.
func (sock *TCPSocket) ResponseMessageForRequestMessage(reqMsg *protocol.Message, resMsg *protocol.Message, timeout time.Duration) error {
	dstAddr := reqMsg.From.IP.String()
	dstPort := reqMsg.From.Port
	_, err := sock.SendMessage(dstAddr, dstPort, resMsg, timeout)
	return err
}

// ResponseMessageToConnection sends a response message to the specified connection.
func (sock *TCPSocket) ResponseMessageToConnection(conn *net.TCPConn, resMsg *protocol.Message) error {
	_, err := sock.writeBytesToConnection(conn, resMsg.Bytes())
	return err
}

// writeBytesToConnection sends the specified bytes to the specified connection.
func (sock *TCPSocket) writeBytesToConnection(conn *net.TCPConn, b []byte) (int, error) {
	toAddr := conn.RemoteAddr()
	localAddr := conn.LocalAddr()

	var nWrote int
	nWrote, err := conn.Write(b)
	if err != nil {
		sock.outputWriteLog(log.LevelError, localAddr.String(), toAddr.String(), hex.EncodeToString(b), 0)
		log.Error(err)
		return nWrote, err
	}

	sock.outputWriteLog(log.LevelTrace, localAddr.String(), toAddr.String(), hex.EncodeToString(b), nWrote)

	return nWrote, nil
}

// dialAndWriteBytes sends the specified bytes to the specified destination.
func (sock *TCPSocket) dialAndWriteBytes(addr string, port int, b []byte, timeout time.Duration) (*net.TCPConn, int, error) {
	toAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		sock.outputWriteLog(log.LevelError, "", toAddr.String(), hex.EncodeToString(b), 0)
		log.Error(err)
		return nil, 0, err
	}

	fromAddr, err := sock.IPAddr()
	if err != nil {
		return nil, 0, err
	}

	conn, err := net.DialTCP("tcp", nil, toAddr)
	if err != nil {
		sock.outputWriteLog(log.LevelError, fromAddr, toAddr.String(), hex.EncodeToString(b), 0)
		log.Error(err)
		return nil, 0, err
	}

	err = conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		conn.Close()
		return nil, 0, err
	}

	nWrote, err := sock.writeBytesToConnection(conn, b)
	if err != nil {
		conn.Close()
	}

	return conn, nWrote, err
}
