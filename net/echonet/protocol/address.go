// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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

// Address represents the address of the message end point.
type Address struct {
	IP   net.IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// NewAddress returns a new remote address.
func NewAddress() *Address {
	addr := &Address{}
	return addr
}

// NewAddressWithString returns a new remote address with the specified string.
func NewAddressWithString(addrString string) (*Address, error) {
	addr := NewAddress()
	err := addr.ParseString(addrString)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// ParseString parses the specified address string.
func (addr *Address) ParseString(addrStr string) error {
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
func (addr *Address) String() string {
	return net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port))
}
