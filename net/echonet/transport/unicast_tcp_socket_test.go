// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"testing"
)

const (
	testUnicastTCPSocketPort = 32001
)

func TestUnicastTCPSocketOpenClose(t *testing.T) {
	ifis, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}
	for _, ifi := range ifis {
		ifaddrs, err := GetInterfaceAddresses(ifi)
		if err != nil {
			t.Error(err)
			continue
		}
		for _, ifaddr := range ifaddrs {
			t.Run(fmt.Sprintf("%s:%s", ifi.Name, ifaddr), func(t *testing.T) {
				sock := NewUnicastTCPSocket()
				err = sock.Bind(ifi, ifaddr, testUnicastTCPSocketPort)
				if err != nil {
					t.Error(err)
					return
				}
				err = sock.Close()
				if err != nil {
					t.Error(err)
				}
			})
		}
	}
}
