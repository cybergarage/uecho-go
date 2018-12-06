// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package transport

import (
	"fmt"
	"net"
	"strconv"

	"golang.org/x/net/ipv4"
)

// Listen listens the Ethonet multicast address with the specified interface.
func (sock *MulticastSocket) Listen(ifi *net.Interface) error {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(MulticastAddress, strconv.Itoa(UDPPort)))
	if err != nil {
		return err
	}

	sock.Conn, err = net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("%s (%s)", err.Error(), ifi.Name)
	}

	pktConn := ipv4.NewPacketConn(sock.Conn)

	err = pktConn.SetMulticastLoopback(true)
	if err != nil {
		return err
	}

	err = pktConn.JoinGroup(ifi, addr)
	if err != nil {
		return err
	}

	return nil
}
