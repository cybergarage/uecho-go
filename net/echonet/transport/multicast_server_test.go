// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type testMulticastServer struct {
	*MulticastServer
	lastMessage *protocol.Message
}

// NewMessageManager returns a new message manager.
func newTestMulticastServer() *testMulticastServer {
	server := &testMulticastServer{
		MulticastServer: NewMulticastServer(),
		lastMessage:     nil,
	}
	server.SetHandler(server)
	return server
}

func (server *testMulticastServer) ProtocolMessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	if isTestMessage(msg) {
		copyMsg, err := protocol.NewMessageWithMessage(msg)
		if err == nil {
			server.lastMessage = copyMsg
		}
	}

	return nil, nil
}

func testMulticastServerWithInterface(t *testing.T, ifi *net.Interface, ifaddr string) {
	t.Helper()

	server := newTestMulticastServer()

	// Start server

	err := server.Start(ifi, ifaddr)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)

	// Send a test message

	now := time.Now()
	msg, err := newTestMessage(uint(now.Unix()))
	if err != nil {
		t.Error(err)
		return
	}

	sock := NewUnicastUDPSocket()
	toAddr := MulticastIPv4Address
	if IsIPv6Address(ifaddr) {
		toAddr = MulticastIPv6Address
	}
	nSent, err := sock.SendMessage(toAddr, Port, msg)
	if err != nil {
		t.Error(err)
	}

	if msgBytes := msg.Bytes(); nSent != len(msgBytes) {
		t.Errorf("%d != %d", nSent, len(msgBytes))
		return
	}

	// Wait a test message

	time.Sleep(time.Second)

	if !msg.Equals(server.lastMessage) {
		ifi, _ := server.MulticastServer.Socket.GetBoundInterface()
		t.Errorf("%v", ifi)
		t.Errorf("%s != %s", msg.String(), server.lastMessage.String())
	}

	// Stop server

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestMulticastServerWithInterface(t *testing.T) {
	ifis, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}

	for _, ifi := range ifis {
		ifaddrs, err := GetInterfaceAddresses(ifi)
		if err != nil {
			t.Error(err)
			continue
		}
		for _, ifaddr := range ifaddrs {
			t.Run(fmt.Sprintf("%s:%s", ifi.Name, ifaddr), func(t *testing.T) {
				testMulticastServerWithInterface(t, ifi, ifaddr)
			})
		}
	}
}
