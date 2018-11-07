// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

// server is an instance for Echonet node.
type server struct {
	*transport.MessageManager
}

// newServer returns a new server.
func newServer() *server {
	server := &server{
		MessageManager: transport.NewMessageManager(),
	}
	return server
}

// GetAddress returns a bound address.
func (server *server) GetAddress() string {
	addrs, err := server.GetBoundAddresses()
	if err != nil || len(addrs) <= 0 {
		return ""
	}
	return addrs[0]
}

// GetPort returns the bound port.
func (server *server) GetPort() int {
	port, err := server.GetBoundPort()
	if err != nil {
		return 0
	}
	return port
}

// Start starts the server.
func (server *server) Start() error {
	err := server.MessageManager.Start()
	if err != nil {
		return nil
	}

	return nil
}

// Stop stops the server.
func (server *server) Stop() error {
	err := server.MessageManager.Stop()
	if err != nil {
		return nil
	}

	return nil
}
