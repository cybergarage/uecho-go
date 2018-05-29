// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"bytes"
	"testing"
)

var testMessageBytes = []byte{
	EHD1,
	EHD2,
	0x00, 0x00,
	0xA0, 0xB0, 0xC0,
	0xD0, 0xE0, 0xF0,
	ESVReadRequest,
	3,
	1, 1, 'a',
	2, 2, 'b', 'c',
	3, 3, 'c', 'd', 'e',
}

func TestNewMessage(t *testing.T) {
	NewMessage()
}

func TestParseMessage(t *testing.T) {

	msg := NewMessage()
	err := msg.Parse(testMessageBytes)
	if err != nil {
		t.Error(err)
		return
	}

	if msg.GetTID() != 0 {
		t.Errorf("%d != %d", msg.GetTID(), 0)
	}

	if msg.GetSourceObjectCode() != 0xA0B0C0 {
		t.Errorf("%03X != %03X", msg.GetSourceObjectCode(), 0xA0B0C0)
	}

	if msg.GetDestinationObjectCode() != 0xD0E0F0 {
		t.Errorf("%03X != %03X", msg.GetDestinationObjectCode(), 0xD0E0F0)
	}

	if msg.GetESV() != ESVReadRequest {
		t.Errorf("%03X != %03X", msg.GetESV(), ESVReadRequest)
	}

	if msg.GetOPC() != 3 {
		t.Errorf("%d != %d", msg.GetESV(), 3)
	}

	for n := 1; n <= int(msg.GetOPC()); n++ {
		prop := msg.GetProperty(n - 1)
		if prop == nil {
			t.Errorf("%d", n)
		}
		if prop.Code != byte(n) {
			t.Errorf("%d != %d", prop.Code, n)
		}
		if len(prop.Data) != n {
			t.Errorf("%d != %d", len(prop.Data), n)
		}
		for i := 0; i < len(prop.Data); i++ {
			dataByte := byte('a' + (n - 1) + i)
			if prop.Data[i] != dataByte {
				t.Errorf("%d != %d", prop.Data[i], dataByte)
			}
		}
	}
}

func TestEncodeMessage(t *testing.T) {

	msg := NewMessage()
	err := msg.Parse(testMessageBytes)
	if err != nil {
		t.Error(err)
		return
	}

	msgBytes := msg.Bytes()
	if bytes.Compare(testMessageBytes, msgBytes) != 0 {
		t.Errorf("%s != %s", string(msgBytes), string(testMessageBytes))
	}
}
