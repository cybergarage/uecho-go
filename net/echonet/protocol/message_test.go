// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"bytes"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage()
	_, err := NewMessageWithMessage(msg)
	if err != nil {
		t.Error(err)
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
