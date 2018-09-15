// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"net"
	"strconv"
	"testing"
)

func TestNewRemoteAddr(t *testing.T) {
	NewRemoteAddr()
}

func TestNewRemoteAddrParse(t *testing.T) {
	testHosts := []string{
		"192.0.2.1",
		"2001:db8::1",
	}
	testPorts := []int{
		25,
		80,
	}

	for n := 0; n < len(testHosts); n++ {
		testHost := testHosts[n]
		testPort := testPorts[n]

		testAddr := net.JoinHostPort(testHost, strconv.Itoa(testPort))

		addr, err := NewRemoteAddrWithString(testAddr)
		if err != nil {
			t.Error(err)
			continue
		}

		if addr.IP.String() != testHost {
			t.Errorf("%s != %s", addr.IP.String(), testHost)
		}

		if addr.Port != testPort {
			t.Errorf("%d != %d", addr.Port, testPort)
		}

		if addr.String() != testAddr {
			t.Errorf("%s != %s", addr.String(), testAddr)
		}
	}
}
