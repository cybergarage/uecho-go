// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"fmt"
	"net"
	"strconv"
)

const (
	errorInvalidAddress = "Invalid address string : %s"
)

// RemoteAddr represents the address of the message end point.
type RemoteAddr struct {
	IP   net.IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// NewRemoteAddr returns a new remote address.
func NewRemoteAddr() *RemoteAddr {
	addr := &RemoteAddr{}
	return addr
}

// NewRemoteAddrWithString returns a new remote address with the specified string.
func NewRemoteAddrWithString(addrString string) (*RemoteAddr, error) {
	addr := NewRemoteAddr()
	err := addr.ParseString(addrString)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// ParseString parses the specified address string.
func (addr *RemoteAddr) ParseString(addrStr string) error {
	hostStr, portStr, err := net.SplitHostPort(addrStr)
	if err != nil {
		return err
	}

	addr.IP = net.ParseIP(hostStr)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf(errorInvalidAddress, addrStr)
	}
	addr.Port = port

	return nil
}

// String returns the node string representation.
func (addr *RemoteAddr) String() string {
	return net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port))
}
