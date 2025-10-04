// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A MulticastServer represents a multicast server.
type MulticastServer struct {
	*Server
	*MulticastSocket

	Channel       chan any
	Handler       MulticastHandler
	UnicastServer *UnicastServer
}

// NewMulticastServer returns a new MulticastServer.
func NewMulticastServer() *MulticastServer {
	server := &MulticastServer{
		Server:          NewServer(),
		MulticastSocket: NewMulticastSocket(),
		Channel:         nil,
		Handler:         nil,
		UnicastServer:   nil,
	}
	return server
}

// SetHandler set a listener.
func (server *MulticastServer) SetHandler(l MulticastHandler) {
	server.Handler = l
}

// SetUnicastServer set a unicast server to response received messages.
func (server *MulticastServer) SetUnicastServer(s *UnicastServer) {
	server.UnicastServer = s
}

// Start starts this server.
func (server *MulticastServer) Start(ifi *net.Interface, ifaddr string) error {
	if err := server.MulticastSocket.Bind(ifi, ifaddr); err != nil {
		return err
	}
	server.Channel = make(chan any)
	go handleMulticastConnection(server, server.Channel)
	return nil
}

// Stop stops this server.
func (server *MulticastServer) Stop() error {
	if err := server.MulticastSocket.Close(); err != nil {
		return err
	}
	return nil
}

func handleMulticastRequestMessage(server *MulticastServer, reqMsg *protocol.Message) {
	server.MulticastSocket.outputReadLog(log.LevelTrace, logSocketTypeUDPMulticast, reqMsg.From.String(), reqMsg.String(), reqMsg.Size())

	if server.Handler == nil {
		return
	}

	resMsg, err := server.Handler.ProtocolMessageReceived(reqMsg)
	if server.UnicastServer == nil || err != nil || resMsg == nil {
		return
	}

	server.UnicastServer.UDPSocket.ResponseMessageForRequestMessage(reqMsg, resMsg)
}

func handleMulticastConnection(server *MulticastServer, cancel chan any) {
	defer server.MulticastSocket.Close()
	for {
		select {
		case <-cancel:
			return
		default:
			reqMsg, err := server.MulticastSocket.ReadMessage()
			if err != nil {
				return
			}
			reqMsg.SetPacketType(protocol.MulticastPacket)

			go handleMulticastRequestMessage(server, reqMsg)
		}
	}
}
