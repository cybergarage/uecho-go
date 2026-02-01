// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package transport

import (
	"fmt"
	"net"
)

// Listen listens the Ethonet multicast address with the specified interface.
func (sock *MulticastSocket) Listen(ifi *net.Interface, ipaddr string, port int) error {
	ipv4Addr := &net.UDPAddr{IP: net.ParseIP(ipaddr), Port: port}
	conn, err := net.ListenUDP("udp4", ipv4Addr)
	if err != nil {
		return fmt.Errorf("%w (%s %s %d)", err, ifi.Name, ipaddr, port)
	}
	sock.Conn = conn

	return nil
}
