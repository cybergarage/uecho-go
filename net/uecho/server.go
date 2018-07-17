// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
	"github.com/cybergarage/uecho-go/net/uecho/transport"
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

// PostAnnounce posts a message.
func (server *server) PostAnnounce(msg *protocol.Message) error {
	return server.SendMessageAll(msg)
}

// PostResponse posts a message to the specified node.
func (server *server) PostResponse(msg *protocol.Message) error {
	return server.SendMessageAll(msg)
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
