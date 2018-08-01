// Copyright 2018 Satoshi Konno. All rights reserved.
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
	mgr := newTestMessageManager()
	mgr.SetMessageListener(mgr)

	err := mgr.Start()
	if err != nil {
		t.Error(err)
		return
	}

	msg, err := newTestMessage()
	if err != nil {
		t.Error(err)
		return
	}

	// Send a test message

	err = mgr.NotifyMessage(msg)
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

func TestNewMessageManagers(t *testing.T) {
	mgrs := []*testMessageManager{
		newTestMessageManager(),
		newTestMessageManager(),
	}

	for n, mgr := range mgrs {
		mgr.SetPort(UDPPort + n)
	}

	// Set the test listener only to the destination managers, mgrs[1]

	srcMgr := mgrs[0]
	dstMgr := mgrs[1]
	dstMgr.SetMessageListener(dstMgr)

	// Start managers

	for _, mgr := range mgrs {
		err := mgr.Start()
		if err != nil {
			t.Error(err)
			return
		}
	}

	// Send a test message from mgrs[0] to mgrs[1]

	msg, err := newTestMessage()
	if err != nil {
		t.Error(err)
		return
	}

	err = srcMgr.NotifyMessage(msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	if dstMgr.lastMessage == nil {
		t.Error("")
	}

	if bytes.Compare(msg.Bytes(), dstMgr.lastMessage.Bytes()) != 0 {
		t.Errorf("%s != %s", string(msg.Bytes()), string(dstMgr.lastMessage.Bytes()))
	}

	// Stop managers

	for _, mgr := range mgrs {
		err := mgr.Stop()
		if err != nil {
			t.Error(err)
			return
		}
	}
}
