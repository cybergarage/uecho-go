// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package transport

import (
	"net"
	"strconv"

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
func (sock *UnicastUDPSocket) Bind(ifi *net.Interface, ifaddr string, port int) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	boundAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(ifaddr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp", boundAddr)
	if err != nil {
		return err
	}

	sock.SetBoundStatus(ifi, ifaddr, port)

	return nil
}

// ResponseMessageForRequestMessage sends a specified response message to the request node.
func (sock *UnicastUDPSocket) ResponseMessageForRequestMessage(reqMsg *protocol.Message, resMsg *protocol.Message) error {
	dstAddr := reqMsg.From.IP.String()
	dstPort := reqMsg.From.Port
	_, err := sock.SendMessage(dstAddr, dstPort, resMsg)
	return err
}
