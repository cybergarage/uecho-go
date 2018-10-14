// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/log"
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

func (mgr *testMessageManager) ProtocolMessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	log.Trace(fmt.Sprintf("ProtocolMessageReceived (R) : %s", msg.String()))
	localPort := mgr.GetPort()
	fromPort := msg.From.Port

	if localPort == fromPort {
		log.Trace(fmt.Sprintf("ProtocolMessageReceived (D) : %s", msg.String()))
		return nil, nil
	}

	if msg.IsESV(protocol.ESVWriteReadRequest) {
		copyMsg, err := protocol.NewMessageWithMessage(msg)
		if err == nil {
			mgr.lastNotificationMessage = copyMsg
			log.Trace(fmt.Sprintf("ProtocolMessageReceived (U) : %s", copyMsg.String()))
		}
	}
	return nil, nil
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

func testMulticastMessagingWithRunningManagers(t *testing.T, mgrs []*testMessageManager) {
	// Initialize managers

	for _, mgr := range mgrs {
		mgr.lastNotificationMessage = nil
	}

	srcMgrs := []*testMessageManager{mgrs[0], mgrs[1]}
	dstMgrs := []*testMessageManager{mgrs[1], mgrs[0]}

	for n := 0; n < len(srcMgrs); n++ {
		srcMgr := srcMgrs[n]
		dstMgr := dstMgrs[n]
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

		time.Sleep(time.Microsecond * 500)

		dstLastMsg := dstMgr.lastNotificationMessage
		if dstLastMsg == nil {
			t.Errorf("%s !=", msg.String())
			continue
		}

		if bytes.Compare(msg.Bytes(), dstLastMsg.Bytes()) != 0 {
			t.Errorf("%s != %s", msg.String(), dstLastMsg.String())
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
	// Initialize managers

	for _, mgr := range mgrs {
		mgr.lastNotificationMessage = nil
	}

	srcMgrs := []*testMessageManager{mgrs[0], mgrs[1]}
	dstMgrs := []*testMessageManager{mgrs[1], mgrs[0]}

	// Send unicast messages, and check the received message

	for n := 0; n < len(srcMgrs); n++ {
		srcMgr := srcMgrs[n]
		dstMgr := dstMgrs[n]
		dstMgr.lastNotificationMessage = nil

		msg, err := newTestMessage(uint(n))
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
			t.Errorf("%s != (nil)", msg)
			continue
		}

		if bytes.Compare(msg.Bytes(), dstMsg.Bytes()) != 0 {
			t.Errorf("%s != %s", msg, dstMsg)
		}

		if checkSourcePort {
			srcPort := srcMgr.GetPort()
			msgPort := dstMsg.GetSourcePort()

			if srcPort != msgPort {
				t.Errorf("%d -!-> %d", srcPort, msgPort)
			}
		} else {
			//t.Logf("Checking source port : %v", checkSourcePort)
		}
	}
}

func testMulticastAndUnicastMessagingWithConfig(t *testing.T, conf *Config, checkSourcePort bool) {
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
		log.Trace(fmt.Sprintf("mgr[%d] : %d", n, mgr.GetPort()))
	}

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

/*
func TestMulticastAndUnicastMessagingWithDefaultConfig(t *testing.T) {
	//log.SetStdoutDebugEnbled(true)
	conf := NewDefaultConfig()
	testMulticastAndUnicastMessagingWithConfig(t, conf, true)
}

func TestMulticastAndUnicastMessagingWithDisableTCPConfig(t *testing.T) {
	//log.SetStdoutDebugEnbled(true)
	conf := NewDefaultConfig()
	conf.SetTCPEnabled(false)
	testMulticastAndUnicastMessagingWithConfig(t, conf, true)
}

func TestMulticastAndUnicastMessagingWithEnableTCPConfig(t *testing.T) {
	//log.SetStdoutDebugEnbled(true)
	conf := NewDefaultConfig()
	conf.SetTCPEnabled(true)
	testMulticastAndUnicastMessagingWithConfig(t, conf, false)
}
*/
