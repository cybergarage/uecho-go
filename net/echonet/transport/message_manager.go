// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	*Config
	Port            uint
	messageListener protocol.MessageListener
	multicastMgr    *MulticastManager
	unicastMgr      *UnicastManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		Config:          NewDefaultConfig(),
		Port:            UDPPort,
		messageListener: nil,
		multicastMgr:    NewMulticastManager(),
		unicastMgr:      NewUnicastManager(),
	}
	return mgr
}

// SetPort sets a listen port.
func (mgr *MessageManager) SetPort(port int) {
	mgr.unicastMgr.SetPort(port)
}

// GetPort returns the listen port.
func (mgr *MessageManager) GetPort() int {
	return mgr.unicastMgr.GetPort()
}

// SetMessageListener set a listener to all managers.
func (mgr *MessageManager) SetMessageListener(l protocol.MessageListener) {
	mgr.multicastMgr.SetListener(l)
	mgr.unicastMgr.SetListener(l)
	mgr.messageListener = l
}

// GetMessageListener returns the listener of the manager.
func (mgr *MessageManager) GetMessageListener() protocol.MessageListener {
	return mgr.messageListener
}

// GetBoundAddresses returns the listen addresses.
func (mgr *MessageManager) GetBoundAddresses() []string {
	boundAddrs := make([]string, 0)
	boundAddrs = append(boundAddrs, mgr.unicastMgr.GetBoundAddresses()...)
	return boundAddrs
}

// GetBoundInterfaces returns the listen interfaces.
func (mgr *MessageManager) GetBoundInterfaces() []net.Interface {
	boundIfs := make([]net.Interface, 0)
	boundIfs = append(boundIfs, mgr.unicastMgr.GetBoundInterfaces()...)
	return boundIfs
}

// SendMessage send a message to the destination address.
func (mgr *MessageManager) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	return mgr.unicastMgr.Write(addr, port, msg.Bytes())
}

// AnnounceMessage sends a message to the multicast address.
func (mgr *MessageManager) AnnounceMessage(msg *protocol.Message) error {
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

	mgr.SetPort(mgr.unicastMgr.GetPort())

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

// IsRunning returns true whether the local managers are running, otherwise false.
func (mgr *MessageManager) IsRunning() bool {
	if !mgr.unicastMgr.IsRunning() {
		return false
	}

	if !mgr.multicastMgr.IsRunning() {
		return false
	}

	return true
}
