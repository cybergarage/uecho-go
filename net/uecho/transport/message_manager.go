// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	multicastMgr *MulticastManager
	unicastMgr   *UnicastManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		multicastMgr: NewMulticastManager(),
		unicastMgr:   NewUnicastManager(),
	}
	return mgr
}

// SetMessageListener set a listener to all managers.
func (mgr *MessageManager) SetMessageListener(l protocol.MessageListener) {
	mgr.multicastMgr.SetListener(l)
	mgr.unicastMgr.SetListener(l)
}

// SendMulticastMessage send a message to the multicast address.
func (mgr *MessageManager) SendMulticastMessage(msg *protocol.Message) error {
	addr, err := net.ResolveUDPAddr("udp", MULTICAST_ADDRESS)
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	_, err = c.Write(msg.Bytes())
	if err != nil {
		return err
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

	/* FIXME
	err = mgr.unicastMgr.Start()
	if err != nil {
		mgr.Stop()
		return err
	}
	*/

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
