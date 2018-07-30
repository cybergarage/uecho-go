// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	Port         uint
	multicastMgr *MulticastManager
	unicastMgr   *UnicastManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		Port:         UDPPort,
		multicastMgr: NewMulticastManager(),
		unicastMgr:   NewUnicastManager(),
	}
	return mgr
}

// SetPort sets a listen port.
func (mgr *MessageManager) SetPort(port int) {
	mgr.unicastMgr.SetPort(port)
}

// SetMessageListener set a listener to all managers.
func (mgr *MessageManager) SetMessageListener(l protocol.MessageListener) {
	mgr.multicastMgr.SetListener(l)
	mgr.unicastMgr.SetListener(l)
}

// SendMessage send a message to the destination address.
func (mgr *MessageManager) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	return mgr.unicastMgr.Write(addr, port, msg.Bytes())
}

// NofityMessage sends a message to the multicast address.
func (mgr *MessageManager) NofityMessage(msg *protocol.Message) error {
	_, err := mgr.SendMessage(MulticastAddress, UDPPort, msg)
	return err
}

// Start starts all transport managers.
func (mgr *MessageManager) Start() error {
	err := mgr.unicastMgr.Start()
	if err != nil {
		mgr.Stop()
		return err
	}

	err = mgr.multicastMgr.Start()
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
