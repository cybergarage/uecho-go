// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"bytes"
	"testing"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage()
	_, err := NewMessageWithMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestNewFormat1TestMessage(t *testing.T) {
	tid := 100
	tidBytes := make([]byte, 2)
	encoding.IntegerToByte(uint(tid), tidBytes)

	opc := 3

	testMessageBytes := []byte{
		EHD1Echonet,
		EHD2Format1,
		tidBytes[0], tidBytes[1],
		0xA0, 0xB0, 0xC0,
		0xD0, 0xE0, 0xF0,
		ESVNotification,
		byte(opc),
		1, 1, 'a',
		2, 2, 'b', 'c',
		3, 3, 'c', 'd', 'e',
	}

	msg, err := NewMessageWithBytes(testMessageBytes)
	if err != nil {
		t.Error(err)
		return
	}

	if !msg.IsTID(uint(tid)) {
		t.Errorf("%d != %d", msg.GetTID(), tid)
	}

	if !msg.IsSourceObjectCode(0xA0B0C0) {
		t.Errorf("%X != %X", msg.GetSourceObjectCode(), 0xA0B0C0)
	}

	if !msg.IsDestinationObjectCode(0xD0E0F0) {
		t.Errorf("%X != %X", msg.GetDestinationObjectCode(), 0xD0E0F0)
	}

	if !msg.IsESV(ESVNotification) {
		t.Errorf("%X != %X", msg.GetESV(), ESVNotification)
	}

	if msg.GetOPC() != opc {
		t.Errorf("%d != %d", msg.GetOPC(), opc)
	}

	for n := 1; n <= opc; n++ {
		prop := msg.GetProperty(n - 1)
		if prop.Code != PropertyCode(n) {
			t.Errorf("%d != %d", prop.Code, n)
		}
		if prop.Size() != n {
			t.Errorf("%d != %d", prop.Size(), n)
		}

	}
}

func TestMessageAddProperty(t *testing.T) {
	msg := NewMessage()

	if msg.GetOPC() != 0 {
		t.Errorf("%d != %d", msg.GetOPC(), 0)
	}

	prop := NewProperty()
	err := msg.AddProperty(prop)
	if err != nil {
		t.Error(err)
		return
	}

	if msg.GetOPC() != 1 {
		t.Errorf("%d != %d", msg.GetOPC(), 1)
	}
}

func TestEncodeMessage(t *testing.T) {

	msg := NewMessage()
	err := msg.ParseBytes(testMessageBytes)
	if err != nil {
		t.Error(err)
		return
	}

	msgBytes := msg.Bytes()
	if bytes.Compare(testMessageBytes, msgBytes) != 0 {
		t.Errorf("%s != %s", string(msgBytes), string(testMessageBytes))
	}
}

func TestMessageEquals(t *testing.T) {

	msg1 := NewMessage()
	err := msg1.ParseBytes(testMessageBytes)
	if err != nil {
		t.Error(err)
		return
	}

	// nil message

	if msg1.Equals(nil) {
		t.Errorf("%s !=", msg1.String())
	}

	// Same message

	msg2 := NewMessage()
	err = msg2.ParseBytes(testMessageBytes)
	if err != nil {
		t.Error(err)
		return
	}

	if !msg1.Equals(msg2) {
		t.Errorf("%s != %s", msg1.String(), msg2.String())
	}
}
