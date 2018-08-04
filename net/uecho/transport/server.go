// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"
)

// A Server represents a server.
type Server struct {
	Interface net.Interface
}

// NewServer returns a new UnicastServer.
func NewServer() *Server {
	server := &Server{}
	return server
}

// SetBoundInterface sets a bind interface.
func (server *Server) SetBoundInterface(i net.Interface) {
	server.Interface = i
}

// GetBoundInterface return a bind interface.
func (server *Server) GetBoundInterface() net.Interface {
	return server.Interface
}

// GetBoundAddresses returns the listen addresses.
func (server *Server) GetBoundAddresses() []net.Addr {
	boundAddrs := make([]net.Addr, 0)
	ifAddrs, err := server.Interface.Addrs()
	if err != nil {
		return boundAddrs
	}
	boundAddrs = append(boundAddrs, ifAddrs...)
	return boundAddrs
}
