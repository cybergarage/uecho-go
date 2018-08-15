// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A MulticastListener represents a listener for MulticastServer.
type MulticastListener interface {
	protocol.MessageListener
}

// A MulticastServer represents a multicast server.
type MulticastServer struct {
	*Server
	Socket   *MulticastSocket
	Listener MulticastListener
}

// NewMulticastServer returns a new MulticastServer.
func NewMulticastServer() *MulticastServer {
	server := &MulticastServer{
		Server:   NewServer(),
		Socket:   NewMulticastSocket(),
		Listener: nil,
	}
	return server
}

// SetListener set a listener.
func (server *MulticastServer) SetListener(l UnicastListener) {
	server.Listener = l
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
		msg, err := server.Socket.ReadMessage()
		if err != nil {
			break
		}

		if server.Listener != nil {
			server.Listener.MessageReceived(msg)
		}
	}
}
