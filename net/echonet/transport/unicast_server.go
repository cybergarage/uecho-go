// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A UnicastServer represents a unicast server.
type UnicastServer struct {
	*UnicastConfig
	*Server
	TCPSocket  *UnicastTCPSocket
	TCPChannel chan interface{}
	UDPSocket  *UnicastUDPSocket
	UDPChannel chan interface{}
	Handler    UnicastHandler
}

// NewUnicastServer returns a new UnicastServer.
func NewUnicastServer() *UnicastServer {
	server := &UnicastServer{
		UnicastConfig: NewDefaultUnicastConfig(),
		Server:        NewServer(),
		TCPSocket:     NewUnicastTCPSocket(),
		TCPChannel:    nil,
		UDPSocket:     NewUnicastUDPSocket(),
		UDPChannel:    nil,
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
func (server *UnicastServer) AnnounceMessage(msg *protocol.Message) error {
	return server.UDPSocket.AnnounceMessage(msg)
}

// Start starts this server.
func (server *UnicastServer) Start(ifi *net.Interface, ifaddr string, port int) error {
	err := server.UDPSocket.Bind(ifi, ifaddr, port)
	if err != nil {
		server.TCPSocket.Close()
		return err
	}
	server.UDPChannel = make(chan interface{})
	go handleUnicastUDPConnection(server, server.UDPChannel)

	if server.IsTCPEnabled() {
		err := server.TCPSocket.Bind(ifi, ifaddr, port)
		if err != nil {
			return err
		}
		server.TCPChannel = make(chan interface{})
		go handleUnicastTCPListener(server, server.TCPChannel)
	}

	server.TCPSocket.SetBoundStatus(ifi, ifaddr, port)
	server.UDPSocket.SetBoundStatus(ifi, ifaddr, port)

	return nil
}

// Stop stops this server.
func (server *UnicastServer) Stop() error {
	var lastErr error

	close(server.UDPChannel)
	err := server.UDPSocket.Close()
	if err != nil {
		lastErr = err
	}

	if server.IsTCPEnabled() {
		close(server.TCPChannel)
		err := server.TCPSocket.Close()
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func handleUnicastUDPRequestMessage(server *UnicastServer, reqMsg *protocol.Message) {
	server.UDPSocket.outputReadLog(log.LevelTrace, logSocketTypeUDPUnicast, reqMsg.From.String(), reqMsg.String(), reqMsg.Size())

	if server.Handler == nil {
		return
	}

	resMsg, err := server.Handler.ProtocolMessageReceived(reqMsg)
	if err != nil || resMsg == nil {
		return
	}

	server.UDPSocket.ResponseMessageForRequestMessage(reqMsg, resMsg)
}

func handleUnicastUDPConnection(server *UnicastServer, cancel chan interface{}) {
	for {
		select {
		case <-cancel:
			return
		default:
			reqMsg, err := server.UDPSocket.ReadMessage()
			if err != nil {
				break
			}
			reqMsg.SetPacketType(protocol.UDPUnicastPacket)

			go handleUnicastUDPRequestMessage(server, reqMsg)
		}
	}
}

func handleUnicastTCPConnection(server *UnicastServer, conn *net.TCPConn) {
	defer conn.Close()

	reqMsg, err := server.TCPSocket.ReadMessage(conn)
	if err != nil {
		return
	}
	reqMsg.SetPacketType(protocol.TCPUnicastPacket)

	if server.Handler == nil {
		return
	}

	resMsg, err := server.Handler.ProtocolMessageReceived(reqMsg)
	if err != nil || resMsg == nil {
		return
	}

	server.TCPSocket.ResponseMessageToConnection(conn, resMsg)
}

func handleUnicastTCPListener(server *UnicastServer, cancel chan interface{}) {
	for {
		select {
		case <-cancel:
			return
		default:
			conn, err := server.TCPSocket.Listener.AcceptTCP()
			if err != nil {
				break
			}

			go handleUnicastTCPConnection(server, conn)
		}
	}
}
