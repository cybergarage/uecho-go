// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"encoding/hex"
	"net"
	"strconv"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
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
func (sock *UnicastUDPSocket) Bind(ifi *net.Interface, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr := ""
	if ifi != nil {
		addr, err = GetInterfaceAddress(ifi)
		if err != nil {
			return err
		}
	}

	boundAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp", boundAddr)
	if err != nil {
		return err
	}

	f, err := sock.Conn.File()
	if err != nil {
		sock.Close()
		return err
	}
	err = sock.SetReuseAddr(f, true)
	if err != nil {
		sock.Close()
		return err
	}

	sock.SetBoundStatus(ifi, addr, port)

	return nil
}

func (sock *UnicastUDPSocket) outputWriteLog(logLevel log.LogLevel, msgTo string, msg string, msgSize int) {
	msgFrom, _ := sock.GetBoundIPAddr()
	outputSocketLog(logLevel, logSocketTypeUDPUnicast, logSocketDirectionWrite, msgFrom, msgTo, msg, msgSize)
}

// SendBytes sends the specified bytes.
func (sock *UnicastUDPSocket) SendBytes(addr string, port int, b []byte) (int, error) {
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

// SendMessage send a message to the destination address.
func (sock *UnicastUDPSocket) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	return sock.SendBytes(addr, port, msg.Bytes())
}

// ResponseMessageForRequestMessage sends a specified response message to the request node
func (sock *UnicastUDPSocket) ResponseMessageForRequestMessage(reqMsg *protocol.Message, resMsg *protocol.Message) error {
	dstAddr := reqMsg.From.IP.String()
	dstPort := reqMsg.From.Port
	_, err := sock.SendMessage(dstAddr, dstPort, resMsg)
	return err
}
