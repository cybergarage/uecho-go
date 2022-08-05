// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bytes"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type testMessageManager struct {
	*MessageManager
	FromPort                int
	FromPacketType          int
	lastNotificationMessage *protocol.Message
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
		protocol.ESVWriteReadRequest,
		3,
		1, 1, 'a',
		2, 2, 'b', 'c',
		3, 3, 'c', 'd', 'e',
	}

	return protocol.NewMessageWithBytes(testMessageBytes)
}

func isTestMessage(msg *protocol.Message) bool {
	return msg.IsESV(protocol.ESVWriteReadRequest)
}

// NewMessageManager returns a new message manager.
func newTestMessageManager() *testMessageManager {
	mgr := &testMessageManager{
		MessageManager:          NewMessageManager(),
		FromPort:                0,
		FromPacketType:          protocol.UnknownPacket,
		lastNotificationMessage: nil,
	}
	return mgr
}

func (mgr *testMessageManager) ProtocolMessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	// log.Trace("ProtocolMessageReceived (R) : %s", msg.String())

	if isTestMessage(msg) {
		copyMsg, err := protocol.NewMessageWithMessage(msg)
		if err == nil {
			// log.Trace("ProtocolMessageReceived (U) : %s", copyMsg.String())
			mgr.lastNotificationMessage = copyMsg
		}
	}

	return nil, nil
}

func testMulticastMessagingWithRunningManagers(t *testing.T, mgrs []*testMessageManager) {
	t.Helper()

	srcMgrs := []*testMessageManager{mgrs[0], mgrs[1]}
	dstMgrs := []*testMessageManager{mgrs[1], mgrs[0]}

	for n := 0; n < len(srcMgrs); n++ {
		srcMgr := srcMgrs[n]
		srcMgr.FromPacketType = protocol.UnknownPacket

		dstMgr := dstMgrs[n]
		dstMgr.FromPort = srcMgr.GetPort()
		dstMgr.FromPacketType = protocol.MulticastPacket
		dstMgr.lastNotificationMessage = nil

		msg, err := newTestMessage(uint(n | 0xF0))
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

		dstLastMsg := dstMgr.lastNotificationMessage
		if dstLastMsg == nil {
			t.Errorf("%s !=", msg.String())
			continue
		}

		log.Trace("CMP(M) : %s ?= %s", msg.String(), dstLastMsg.String())

		if !msg.Equals(dstLastMsg) {
			log.Trace("CMP(M) : %s != %s", msg.String(), dstLastMsg.String())
			t.Errorf("CMP(M) : %s != %s", msg.String(), dstLastMsg.String())
			continue
		}

		srcPort := srcMgr.GetPort()
		msgPort := dstLastMsg.GetSourcePort()

		if srcPort != msgPort {
			t.Errorf("%d -!-> %d", srcPort, msgPort)
		}
	}
}

func testUnicastMessagingWithRunningManagers(t *testing.T, mgrs []*testMessageManager, checkSourcePort bool) {
	t.Helper()

	srcMgrs := []*testMessageManager{mgrs[0], mgrs[1]}
	dstMgrs := []*testMessageManager{mgrs[1], mgrs[0]}

	// Send unicast messages, and check the received message

	for n := 0; n < len(srcMgrs); n++ {
		srcMgr := srcMgrs[n]
		srcMgr.FromPacketType = protocol.UnknownPacket

		dstMgr := dstMgrs[n]
		dstMgr.FromPort = srcMgr.GetPort()
		dstMgr.FromPacketType = protocol.UnicastPacket
		dstMgr.lastNotificationMessage = nil

		msg, err := newTestMessage(uint(n))
		if err != nil {
			t.Error(err)
			continue
		}

		dstPort := dstMgr.GetPort()
		dstAddrs, err := dstMgr.GetBoundAddresses()
		if err != nil {
			t.Error(err)
		}
		if len(dstAddrs) == 0 {
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

		dstLastMsg := dstMgr.lastNotificationMessage
		if dstLastMsg == nil {
			t.Errorf("%s != (nil)", msg)
			continue
		}

		log.Trace("CMP(U) : %s ?= %s", msg.String(), dstLastMsg.String())

		if !bytes.Equal(msg.Bytes(), dstLastMsg.Bytes()) {
			log.Error("CMP(U) : %s != %s", msg.String(), dstLastMsg.String())
			t.Errorf("CMP(U) : %s != %s", msg, dstLastMsg)
		}

		if checkSourcePort {
			srcPort := srcMgr.GetPort()
			msgPort := dstLastMsg.GetSourcePort()

			if srcPort != msgPort {
				t.Errorf("%d -!-> %d", srcPort, msgPort)
			}
		} else {
			t.Logf("Checking source port : %v", checkSourcePort)
		}
	}
}

func testMulticastAndUnicastMessagingWithConfig(t *testing.T, conf *Config, checkSourcePort bool) {
	t.Helper()

	mgrs := []*testMessageManager{
		newTestMessageManager(),
		newTestMessageManager(),
	}

	for n, mgr := range mgrs {
		mgr.SetConfig(conf)
		mgr.SetPort(UDPPort + n)
		mgr.SetMessageHandler(mgr)
	}

	// Start managers

	for n, mgr := range mgrs {
		err := mgr.Start()
		if err != nil {
			t.Error(err)
			return
		}
		log.Trace("mgr[%d] : %d", n, mgr.GetPort())
	}

	time.Sleep(time.Second)

	// Send multicast messages, and check the received message

	testMulticastMessagingWithRunningManagers(t, mgrs)

	// Send unicast messages, and check the received message

	testUnicastMessagingWithRunningManagers(t, mgrs, checkSourcePort)

	// Stop managers

	for _, mgr := range mgrs {
		err := mgr.Stop()
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestMulticastAndUnicastMessagingWithDefaultConfig(t *testing.T) {
	// log.SetStdoutDebugEnbled(true)
	conf := newTestDefaultConfig()
	testMulticastAndUnicastMessagingWithConfig(t, conf, true)
}

func TestMulticastAndUnicastMessagingWithDisableTCPConfig(t *testing.T) {
	// log.SetStdoutDebugEnbled(true)
	conf := newTestDefaultConfig()
	conf.SetTCPEnabled(false)
	testMulticastAndUnicastMessagingWithConfig(t, conf, true)
}

func TestMulticastAndUnicastMessagingWithEnableTCPConfig(t *testing.T) {
	// log.SetStdoutDebugEnbled(true)
	conf := newTestDefaultConfig()
	conf.SetTCPEnabled(true)
	testMulticastAndUnicastMessagingWithConfig(t, conf, false)
}
