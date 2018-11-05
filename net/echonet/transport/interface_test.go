// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func TestGetAvailableInterfaces(t *testing.T) {
	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
	}
	if len(ifs) <= 0 {
		t.Errorf("available interface is not found")
	}
}

func TestGetAvailableAddresses(t *testing.T) {
	addrs, err := GetAvailableAddresses()
	if err != nil {
		t.Error(err)
	}
	if len(addrs) <= 0 {
		t.Errorf("available address is not found")
	}
}

func TestIPv6Addresses(t *testing.T) {
	addrs := []string{
		"::1",
		"fe80::1875:6549:801:d487",
	}

	for n, addr := range addrs {
		if !IsIPv6Address(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}
}

func TestIPv4Addresses(t *testing.T) {
	addrs := []string{
		"127.0.0.1",
		"192.168.0.1",
	}

	for n, addr := range addrs {
		if IsIPv6Address(addr) {
			t.Errorf("[%d] %s", n, addr)
		}
	}
}
