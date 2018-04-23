// Copyright 2017 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"net"
)

// A MulticastListener represents a listener for MulticastServer.
type MulticastListener interface {
	MessageReceived(*Message)
}

// A MulticastServer represents a packet.
type MulticastServer struct {
	Socket    *MulticastSocket
	Listener  MulticastListener
	Interface net.Interface
}

// NewMulticastServer returns a new MulticastServer.
func NewMulticastServer() *MulticastServer {
	server := &MulticastServer{}
	server.Socket = NewMulticastSocket()
	server.Listener = nil
	return server
}

// Start starts this server.
func (server *MulticastServer) Start(ifi net.Interface) error {
	err := server.Socket.Bind(ifi)
	if err != nil {
		return err
	}
	server.Interface = ifi
	go handleMulticastConnection(server)
	return nil
}

// Stop stops this server.
func (server *MulticastServer) Stop() error {
	err := server.Socket.Close()
	if err != nil {
		return err
	}
	return nil
}

func handleMulticastConnection(server *MulticastServer) {
	for {
		msg, err := server.Socket.Read()
		if err != nil {
			break
		}

		if server.Listener != nil {
			server.Listener.MessageReceived(msg)
		}
	}
}