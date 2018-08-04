// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/uecho/encoding"
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

func newTestMessage(tid uint) (*protocol.Message, error) {
	tidBytes := make([]byte, 2)
	encoding.IntegerToByte(tid, tidBytes)

	testMessageBytes := []byte{
		protocol.EHD1,
		protocol.EHD2,
		tidBytes[0], tidBytes[1],
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

	msg, err := newTestMessage(0)
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

func TestMulticastMessaging(t *testing.T) {
	mgrs := []*testMessageManager{
		newTestMessageManager(),
		newTestMessageManager(),
	}

	for n, mgr := range mgrs {
		mgr.SetPort(UDPPort + n)
		mgr.SetMessageListener(mgr)
	}

	// Start managers

	for _, mgr := range mgrs {
		err := mgr.Start()
		if err != nil {
			t.Error(err)
			return
		}
	}

	// Send multicast messages, and check the received message

	srcMgrs := []*testMessageManager{mgrs[0], mgrs[1]}
	dstMgrs := []*testMessageManager{mgrs[1], mgrs[0]}

	for n := 0; n < len(srcMgrs); n++ {
		srcMgr := srcMgrs[n]
		dstMgr := dstMgrs[n]

		msg, err := newTestMessage(uint(rand.Uint32()))
		if err != nil {
			t.Error(err)
			continue
		}

		err = srcMgr.NotifyMessage(msg)
		if err != nil {
			t.Error(err)
			continue
		}

		time.Sleep(time.Second)

		dstMsg := dstMgr.lastMessage
		if dstMsg == nil {
			t.Error("")
		}

		if bytes.Compare(msg.Bytes(), dstMsg.Bytes()) != 0 {
			t.Errorf("%s != %s", string(msg.Bytes()), string(dstMsg.Bytes()))
		}

		srcPort := srcMgr.GetPort()
		msgPort := dstMsg.GetSourcePort()

		if srcPort != msgPort {
			t.Errorf("%d != %d", srcPort, msgPort)
		}
	}

	// Send unicast messages, and check the received message

	for n := 0; n < len(srcMgrs); n++ {
		srcMgr := srcMgrs[n]
		dstMgr := dstMgrs[n]

		msg, err := newTestMessage(uint(rand.Uint32()))
		if err != nil {
			t.Error(err)
			continue
		}

		dstPort := dstMgr.GetPort()
		dstAddrs := dstMgr.GetBoundAddresses()
		if len(dstAddrs) <= 0 {
			t.Errorf("Not found available interfaces ")
			continue
		}

		dstAddr := dstAddrs[0]
		_, err = srcMgr.SendMessage(dstAddr.String(), dstPort, msg)
		if err != nil {
			t.Error(err)
			continue
		}

		time.Sleep(time.Second)

		dstMsg := dstMgr.lastMessage
		if dstMsg == nil {
			t.Error("")
		}

		if bytes.Compare(msg.Bytes(), dstMsg.Bytes()) != 0 {
			t.Errorf("%s != %s", string(msg.Bytes()), string(dstMsg.Bytes()))
		}

		srcPort := srcMgr.GetPort()
		msgPort := dstMsg.GetSourcePort()

		if srcPort != msgPort {
			t.Errorf("%d != %d", srcPort, msgPort)
		}
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
