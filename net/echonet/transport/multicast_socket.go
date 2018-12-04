// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"strconv"
	//"golang.org/x/net/ipv4"
)

// A MulticastSocket represents a socket.
type MulticastSocket struct {
	*UDPSocket
}

// NewMulticastSocket returns a new MulticastSocket.
func NewMulticastSocket() *MulticastSocket {
	sock := &MulticastSocket{
		UDPSocket: NewUDPSocket(),
	}
	return sock
}

// ListenMulticastUDP listens the Ethonet multicast address with the specified interface.
/*
func (sock *MulticastSocket) ListenMulticastUDP(ifi *net.Interface) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(MulticastAddress, strconv.Itoa(UDPPort)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("%s (%s)", err.Error(), ifi.Name)
	}

	pktConn := ipv4.NewPacketConn(conn)

	err = pktConn.pktConn(ifi, addr)
	if err != nil {
		return err
	}

	err = pc.JoinGroup(ifi, addr)
	if err != nil {
		return err
	}

	err = pktConn.SetMulticastLoopback(true)
	if err != nil {
		return err
	}

	return nil
}
*/

// Bind binds to the Echonet multicast address with the specified interface.
func (sock *MulticastSocket) Bind(ifi *net.Interface) error {
	err := sock.Close()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(MulticastAddress, strconv.Itoa(UDPPort)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenMulticastUDP("udp", ifi, addr)
	//err = sock.ListenMulticastUDP(ifi)
	if err != nil {
		return fmt.Errorf("%s (%s)", err.Error(), ifi.Name)
	}

	sock.Conn.SetReadBuffer(sock.GetReadBufferSize())

	f, err := sock.Conn.File()
	if err != nil {
		return err
	}
	err = sock.SetReuseAddr(f, true)
	if err != nil {
		return err
	}

	sock.SetBoundStatus(ifi, MulticastAddress, UDPPort)

	return nil
}
