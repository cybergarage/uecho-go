// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"
)

// A Server represents a server.
type Server struct {
	Interface *net.Interface
}

// NewServer returns a new UnicastServer.
func NewServer() *Server {
	server := &Server{
		Interface: nil,
	}
	return server
}

// SetBoundInterface sets a bind interface.
func (server *Server) SetBoundInterface(i *net.Interface) {
	server.Interface = i
}

// GetBoundInterface return a bind interface.
func (server *Server) GetBoundInterface() *net.Interface {
	return server.Interface
}

// GetBoundAddresses returns the listen addresses.
func (server *Server) GetBoundAddresses() []string {
	boundAddrs := make([]string, 0)
	ifAddr, err := GetInterfaceAddress(server.Interface)
	if err != nil {
		return boundAddrs
	}
	boundAddrs = append(boundAddrs, ifAddr)
	return boundAddrs
}
