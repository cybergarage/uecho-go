// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	NewMessage()
}

func TestParseMessage(t *testing.T) {
	testMessageBytes := []byte{
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

	/*
	  for (int n=1; n<=uecho_message_getopc(msg); n++) {
	    uEchoProperty *prop = uecho_message_getproperty(msg, (n-1));
	    BOOST_CHECK(prop);
	    BOOST_CHECK_EQUAL(uecho_property_getcode(prop), n);
	    BOOST_CHECK_EQUAL(uecho_property_getdatasize(prop), n);
	    byte *data = uecho_property_getdata(prop);
	    BOOST_CHECK(data);
	    for (int i=0; i<uecho_property_getdatasize(prop); i++) {
	      BOOST_CHECK_EQUAL(data[i], 'a' + (n-1) + i);
	    }
	  }

	  uecho_message_delete(msg);
	*/
}
