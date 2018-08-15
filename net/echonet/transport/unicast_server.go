// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A UnicastListener represents a listener for UnicastServer.
type UnicastListener interface {
	protocol.MessageListener
}

// A UnicastServer represents a unicast server.
type UnicastServer struct {
	*Server
	Socket   *UnicastSocket
	Listener UnicastListener
}

// NewUnicastServer returns a new UnicastServer.
func NewUnicastServer() *UnicastServer {
	server := &UnicastServer{
		Server:   NewServer(),
		Socket:   NewUnicastSocket(),
		Listener: nil,
	}
	return server
}

// SetListener set a listener.
func (server *UnicastServer) SetListener(l UnicastListener) {
	server.Listener = l
}

// Start starts this server.
func (server *UnicastServer) Start(ifi net.Interface, port int) error {
	err := server.Socket.Bind(ifi, port)
	if err != nil {
		return err
	}
	server.Interface = ifi
	go handleUnicastConnection(server)
	return nil
}

// Stop stops this server.
func (server *UnicastServer) Stop() error {
	err := server.Socket.Close()
	if err != nil {
		return err
	}
	return nil
}

func handleUnicastConnection(server *UnicastServer) {
	for {
		msg, err := server.Socket.ReadMessage()
		if err != nil {
			break
		}

		if server.Listener != nil {
			server.Listener.MessageReceived(msg)
		}
	}
}
