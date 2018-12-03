// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
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
	//log.Trace("ProtocolMessageReceived (R) : %s", msg.String())

	if isTestMessage(msg) {
		copyMsg, err := protocol.NewMessageWithMessage(msg)
		if err == nil {
			//log.Trace("ProtocolMessageReceived (U) : %s", copyMsg.String())
			server.lastMessage = copyMsg
		}
	}

	return nil, nil
}

func testMulticastServerWithInterface(t *testing.T, ifi *net.Interface) {
	server := newTestMulticastServer()

	// Start server

	err := server.Start(ifi)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	// Send a test message

	now := time.Now()
	msg, err := newTestMessage(uint(now.Unix()))

	sock := NewUnicastUDPSocket()
	nSent, err := sock.SendMessage(MulticastAddress, Port, msg)
	if err != nil {
		t.Error(err)
	}

	msgBytes := msg.Bytes()
	if nSent != len(msgBytes) {
		t.Errorf("%d != %d", nSent, len(msgBytes))
	}

	// Wait a test message

	time.Sleep(time.Second)

	if !msg.Equals(server.lastMessage) {
		t.Errorf("%s != %s", msg.String(), server.lastMessage.String())
	}

	// Stop server

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestMulticastServerWithInterface(t *testing.T) {
	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}

	testMulticastServerWithInterface(t, ifs[0])
}

func TestMulticastServerWithNoInterface(t *testing.T) {
	testMulticastServerWithInterface(t, nil)
}
