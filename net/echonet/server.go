// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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

// Start starts the server.
func (server *server) Start() error {
	err := server.MessageManager.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the server.
func (server *server) Stop() error {
	err := server.MessageManager.Stop()
	if err != nil {
		return err
	}

	return nil
}
