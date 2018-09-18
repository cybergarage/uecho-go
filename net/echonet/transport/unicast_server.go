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
	TCPSocket *UnicastTCPSocket
	UDPSocket *UnicastUDPSocket
	Listener  UnicastListener
}

// NewUnicastServer returns a new UnicastServer.
func NewUnicastServer() *UnicastServer {
	server := &UnicastServer{
		Server:    NewServer(),
		TCPSocket: NewUnicastTCPSocket(),
		UDPSocket: NewUnicastUDPSocket(),
		Listener:  nil,
	}
	return server
}

// SetListener set a listener.
func (server *UnicastServer) SetListener(l UnicastListener) {
	server.Listener = l
}

// Start starts this server.
func (server *UnicastServer) Start(ifi net.Interface, port int) error {
	err := server.TCPSocket.Bind(ifi, port)
	if err != nil {
		return err
	}

	err = server.UDPSocket.Bind(ifi, port)
	if err != nil {
		server.TCPSocket.Close()
		return err
	}

	server.Interface = ifi
	go handleUnicastUDPConnection(server)

	return nil
}

// Stop stops this server.
func (server *UnicastServer) Stop() error {
	var lastErr error

	err := server.TCPSocket.Close()
	if err != nil {
		lastErr = err
	}

	err = server.UDPSocket.Close()
	if err != nil {
		lastErr = err
	}

	return lastErr
}

func handleUnicastUDPConnection(server *UnicastServer) {
	for {
		msg, err := server.UDPSocket.ReadMessage()
		if err != nil {
			break
		}

		if server.Listener != nil {
			server.Listener.ProtocolMessageReceived(msg)
		}
	}
}

func handleUnicastTCPListener(server *UnicastServer) {
	for {
		conn, err := server.TCPSocket.Listener.Accept()
		if err != nil {
			break
		}

		go handleUnicastTCPConnection(server, conn)
	}
}

func handleUnicastTCPConnection(server *UnicastServer, conn net.Conn) {
	msg, err := server.TCPSocket.ReadMessage(conn)
	if err != nil {
		return
	}

	conn.Close()

	if server.Listener != nil {
		server.Listener.ProtocolMessageReceived(msg)
	}
}
