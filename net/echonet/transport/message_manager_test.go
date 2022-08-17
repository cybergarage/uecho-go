// Copyright 2018 Satoshi Konno. All rights reserved.
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

// nolint: staticcheck
func (mgr *testMessageManager) ProtocolMessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	// log.Trace("ProtocolMessageReceived (R) : %s", msg.String())

	if isTestMessage(msg) {
		copyMsg, err := protocol.NewMessageWithMessage(msg)
		if err == nil {
			// log.Trace("lastNotificationMessage (O) : %s", copyMsg.String())
			mgr.lastNotificationMessage = copyMsg
		} else {
			log.Error("ProtocolMessageReceived (X) : %s", msg.String())
		}
	} else {
		// log.Trace("ProtocolMessageReceived (-) : %s", msg.String())
	}

	return nil, nil
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

		t.Run(fmt.Sprintf("Unicast:%d->%d", srcMgr.GetPort(), dstMgr.GetPort()), func(t *testing.T) {
			msg, err := newTestMessage(uint(n))
			if err != nil {
				t.Error(err)
				return
			}

			dstPort := dstMgr.GetPort()
			dstAddrs := dstMgr.GetBoundAddresses()
			if len(dstAddrs) == 0 {
				t.Errorf("Not found available interfaces ")
				return
			}

			dstAddr := dstAddrs[0]
			_, err = srcMgr.SendMessage(dstAddr, dstPort, msg)
			if err != nil {
				t.Error(err)
				return
			}

			time.Sleep(time.Second)

			dstLastMsg := dstMgr.lastNotificationMessage
			if dstLastMsg == nil {
				t.Errorf("%s != (nil)", msg)
				return
			}

			// log.Trace("CMP(U) : %s ?= %s", msg.String(), dstLastMsg.String())

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
		})
	}
}

func testUnicastMessagingWithConfig(t *testing.T, conf *Config, checkSourcePort bool) {
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

	for _, mgr := range mgrs {
		err := mgr.Start()
		if err != nil {
			t.Error(err)
			return
		}
	}

	time.Sleep(time.Second)

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

func TestUnicastMessaging(t *testing.T) {
	log.SetStdoutDebugEnbled(true)
	defer log.SetStdoutDebugEnbled(false)

	t.Run("Default", func(t *testing.T) {
		conf := newTestDefaultConfig()
		testUnicastMessagingWithConfig(t, conf, true)
	})

	t.Run("TCPEnabled", func(t *testing.T) {
		conf := newTestDefaultConfig()
		conf.SetTCPEnabled(true)
		testUnicastMessagingWithConfig(t, conf, false)
	})
}
