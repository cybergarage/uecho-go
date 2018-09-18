// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type testMessageManager struct {
	*MessageManager
	lastNotificationMessage *protocol.Message
}

// NewMessageManager returns a new message manager.
func newTestMessageManager() *testMessageManager {
	mgr := &testMessageManager{
		MessageManager:          NewMessageManager(),
		lastNotificationMessage: nil,
	}
	return mgr
}

func (mgr *testMessageManager) ProtocolMessageReceived(msg *protocol.Message) {
	if !msg.IsESV(protocol.ESVNotificationRequest) {
		return
	}
	mgr.lastNotificationMessage = msg
}

func newTestMessage(tid uint) (*protocol.Message, error) {
	tidBytes := make([]byte, 2)
	encoding.IntegerToByte(tid, tidBytes)

	testMessageBytes := []byte{
		protocol.EHD1Echonet,
		protocol.EHD2Format1,
		tidBytes[0], tidBytes[1],
		0xA0, 0xB0, 0xC0,
		0xD0, 0xE0, 0xF0,
		protocol.ESVNotificationRequest,
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

	err = mgr.AnnounceMessage(msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	if mgr.lastNotificationMessage != nil {
		if bytes.Compare(msg.Bytes(), mgr.lastNotificationMessage.Bytes()) != 0 {
			t.Errorf("%s != %s", msg, mgr.lastNotificationMessage)
		}
	} else {
		t.Error("")
	}

	err = mgr.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMulticastAndUnicastMessaging(t *testing.T) {
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

		err = srcMgr.AnnounceMessage(msg)
		if err != nil {
			t.Error(err)
			continue
		}

		time.Sleep(time.Second)

		dstMsg := dstMgr.lastNotificationMessage
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
		_, err = srcMgr.SendMessage(dstAddr, dstPort, msg)
		if err != nil {
			t.Error(err)
			continue
		}

		time.Sleep(time.Second)

		dstMsg := dstMgr.lastNotificationMessage
		if dstMsg == nil {
			t.Error("")
		}

		if bytes.Compare(msg.Bytes(), dstMsg.Bytes()) != 0 {
			t.Errorf("%s != %s", msg, dstMsg)
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
