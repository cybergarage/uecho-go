// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"encoding/hex"
	"net"
	"strconv"

	"github.com/cybergarage/uecho-go/net/echonet/log"
)

// A UnicastUDPSocket represents a socket.
type UnicastUDPSocket struct {
	*UDPSocket
}

// NewUnicastUDPSocket returns a new UnicastUDPSocket.
func NewUnicastUDPSocket() *UnicastUDPSocket {
	sock := &UnicastUDPSocket{
		UDPSocket: NewUDPSocket(),
	}
	return sock
}

// Bind binds to Echonet multicast address.
func (sock *UnicastUDPSocket) Bind(ifi net.Interface, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr, err := GetInterfaceAddress(ifi)
	if err != nil {
		return err
	}

	boundAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp", boundAddr)
	if err != nil {
		return err
	}

	sock.SetBoundStatus(ifi, addr, port)

	return nil
}

func (sock *UnicastUDPSocket) outputWriteLog(logLevel log.LogLevel, msgTo string, msg string, msgSize int) {
	msgFrom, _ := sock.GetBoundIPAddr()
	outputSocketLog(logLevel, logSocketTypeUDPUnicast, logSocketDirectionWrite, msgFrom, msgTo, msg, msgSize)
}

// Write sends the specified bytes.
func (sock *UnicastUDPSocket) Write(addr string, port int, b []byte) (int, error) {
	toAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return 0, err
	}

	// Send from binding port

	if sock.Conn != nil {
		n, err := sock.Conn.WriteToUDP(b, toAddr)
		sock.outputWriteLog(log.LoggerLevelTrace, toAddr.String(), hex.EncodeToString(b), n)
		if err != nil {
			log.Error(err.Error())
		}
		return n, err
	}

	// Send from no binding port

	conn, err := net.DialUDP("udp", nil, toAddr)
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	n, err := conn.Write(b)
	sock.outputWriteLog(log.LoggerLevelTrace, toAddr.String(), hex.EncodeToString(b), n)
	if err != nil {
		log.Error(err.Error())
	}
	conn.Close()

	return n, err
}
