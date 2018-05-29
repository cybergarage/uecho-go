// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bytes"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

type testMessageManager struct {
	*MessageManager
	lastMessage *protocol.Message
}

// NewMessageManager returns a new message manager.
func newTestMessageManager() *testMessageManager {
	mgr := &testMessageManager{
		MessageManager: NewMessageManager(),
		lastMessage:    nil,
	}
	return mgr
}

func (mgr *testMessageManager) MessageReceived(msg *protocol.Message) {
	mgr.lastMessage = msg
}

func newTestMessage() (*protocol.Message, error) {
	testMessageBytes := []byte{
		protocol.EHD1,
		protocol.EHD2,
		0x00, 0x00,
		0xA0, 0xB0, 0xC0,
		0xD0, 0xE0, 0xF0,
		protocol.ESVReadRequest,
		3,
		1, 1, 'a',
		2, 2, 'b', 'c',
		3, 3, 'c', 'd', 'e',
	}

	return protocol.NewMessageWithBytes(testMessageBytes)
}

func TestNewMessageManager(t *testing.T) {
	msg, err := newTestMessage()
	if err != nil {
		t.Error(err)
		return
	}

	mgr := newTestMessageManager()
	mgr.SetMessageListener(mgr)

	err = mgr.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = mgr.SendMulticastMessage(msg)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second)
	if mgr.lastMessage == nil {
		t.Error("")
	}
	if bytes.Compare(msg.Bytes(), mgr.lastMessage.Bytes()) != 0 {
		t.Errorf("%s != %s", string(msg.Bytes()), string(mgr.lastMessage.Bytes()))
	}

	err = mgr.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
