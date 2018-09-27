// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A UnicastHandler represents a listener for UnicastServer.
type UnicastHandler interface {
	protocol.MessageHandler
}

// A UnicastServer represents a unicast server.
type UnicastServer struct {
	*UnicastConfig
	*Server
	TCPSocket *UnicastTCPSocket
	UDPSocket *UnicastUDPSocket
	Handler   UnicastHandler
}

// NewUnicastServer returns a new UnicastServer.
func NewUnicastServer() *UnicastServer {
	server := &UnicastServer{
		UnicastConfig: NewDefaultUnicastConfig(),
		Server:        NewServer(),
		TCPSocket:     NewUnicastTCPSocket(),
		UDPSocket:     NewUnicastUDPSocket(),
		Handler:       nil,
	}
	return server
}

// SetHandler set a listener.
func (server *UnicastServer) SetHandler(l UnicastHandler) {
	server.Handler = l
}

// SendMessage send a message to the destination address.
func (server *UnicastServer) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	if server.IsTCPEnabled() {
		n, err := server.TCPSocket.SendMessage(addr, port, msg, server.GetConnectionTimeout())
		if err == nil {
			return n, nil
		}
	}

	return server.UDPSocket.SendMessage(addr, port, msg)
}

// AnnounceMessage sends a message to the multicast address.
func (server *UnicastServer) AnnounceMessage(addr string, port int, msg *protocol.Message) error {
	_, err := server.UDPSocket.SendBytes(addr, port, msg.Bytes())
	return err
}

// ResponseMessageForRequestMessage sends a specified response message to the request node
func (server *UnicastServer) ResponseMessageForRequestMessage(reqMsg *protocol.Message, resMsg *protocol.Message) error {
	return server.UDPSocket.ResponseMessageForRequestMessage(reqMsg, resMsg)
}

// Start starts this server.
func (server *UnicastServer) Start(ifi net.Interface, port int) error {
	err := server.UDPSocket.Bind(ifi, port)
	if err != nil {
		server.TCPSocket.Close()
		return err
	}
	go handleUnicastUDPConnection(server)

	if server.IsTCPEnabled() {
		err := server.TCPSocket.Bind(ifi, port)
		if err != nil {
			return err
		}
	}
	go handleUnicastTCPHandler(server)

	server.Interface = ifi

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
		reqMsg, err := server.UDPSocket.ReadMessage()
		if err != nil {
			break
		}

		server.UDPSocket.outputReadLog(log.LoggerLevelTrace, logSocketTypeUDPUnicast, reqMsg.From.String(), reqMsg.String(), reqMsg.Size())

		if server.Handler == nil {
			continue
		}

		resMsg, err := server.Handler.ProtocolMessageReceived(reqMsg)
		if err != nil || resMsg == nil {
			continue
		}

		server.UDPSocket.ResponseMessageForRequestMessage(reqMsg, resMsg)
	}
}

func handleUnicastTCPHandler(server *UnicastServer) {
	for {
		conn, err := server.TCPSocket.Listener.AcceptTCP()
		if err != nil {
			break
		}

		go handleUnicastTCPConnection(server, conn)
	}
}

func handleUnicastTCPConnection(server *UnicastServer, conn *net.TCPConn) {
	reqMsg, err := server.TCPSocket.ReadMessage(conn)
	if err != nil {
		return
	}

	defer conn.Close()

	if server.Handler != nil {
		resMsg, err := server.Handler.ProtocolMessageReceived(reqMsg)
		if err != nil && resMsg != nil {
			server.TCPSocket.ResponseMessage(conn, resMsg)
		}
	}
}
