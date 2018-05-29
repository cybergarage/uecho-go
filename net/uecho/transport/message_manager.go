// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"
	"time"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	multicastMgr *MulticastServerManager
	unicastMgr   *UnicastServerManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		multicastMgr: NewMulticastServerManager(),
		unicastMgr:   NewUnicastServerManager(),
	}
	return mgr
}

// SendMulticastMessage send a message to the multicast address.
func (mgr *MessageManager) SendMulticastMessage(msg protocol.Message) error {
	addr, err := net.ResolveUDPAddr("udp", MULTICAST_ADDRESS)
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	
	for {
		c.Write([]byte("hello, world\n"))
		time.Sleep(1 * time.Second)
	}

	return nil
}

// Start starts all transport managers.
func (mgr *MessageManager) Start() error {
	err := mgr.multicastMgr.Start()
	if err != nil {
		mgr.Stop()
		return err
	}

	err = mgr.unicastMgr.Start()
	if err != nil {
		mgr.Stop()
		return err
	}
	return nil
}

// Stop stops all transport managers.
func (mgr *MessageManager) Stop() error {
	err := mgr.multicastMgr.Stop()
	if err != nil {
		return err
	}

	err = mgr.unicastMgr.Stop()
	if err != nil {
		return err
	}

	return nil
}
