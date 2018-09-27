// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	Port           uint
	messageHandler protocol.MessageHandler
	multicastMgr   *MulticastManager
	unicastMgr     *UnicastManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		Port:           UDPPort,
		messageHandler: nil,
		multicastMgr:   NewMulticastManager(),
		unicastMgr:     NewUnicastManager(),
	}
	return mgr
}

// SetConfig sets all flags.
func (mgr *MessageManager) SetConfig(newConfig *Config) {
	mgr.unicastMgr.SetConfig(newConfig.UnicastConfig)
}

// SetPort sets a listen port.
func (mgr *MessageManager) SetPort(port int) {
	mgr.unicastMgr.SetPort(port)
}

// GetPort returns the listen port.
func (mgr *MessageManager) GetPort() int {
	return mgr.unicastMgr.GetPort()
}

// SetMessageHandler set a listener to all managers.
func (mgr *MessageManager) SetMessageHandler(h protocol.MessageHandler) {
	mgr.multicastMgr.SetHandler(h)
	mgr.unicastMgr.SetHandler(h)
	mgr.messageHandler = h
}

// GetMessageHandler returns the listener of the manager.
func (mgr *MessageManager) GetMessageHandler() protocol.MessageHandler {
	return mgr.messageHandler
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
	return mgr.unicastMgr.SendMessage(addr, port, msg)
}

// AnnounceMessage sends a message to the multicast address.
func (mgr *MessageManager) AnnounceMessage(msg *protocol.Message) error {
	return mgr.unicastMgr.AnnounceMessage(MulticastAddress, UDPPort, msg)
}

// Start starts all transport managers.
func (mgr *MessageManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	for _, ifi := range ifis {
		multicastServer, err := mgr.multicastMgr.Start(ifi)
		if err != nil {
			return err
		}

		unicastServer, err := mgr.unicastMgr.Start(ifi)
		if err != nil {
			return err
		}

		multicastServer.SetUnicastServer(unicastServer)
		mgr.SetPort(mgr.unicastMgr.GetPort())
	}

	return nil
}

// Stop stops all transport managers.
func (mgr *MessageManager) Stop() error {
	var lastErr error
	err := mgr.multicastMgr.Stop()
	if err != nil {
		lastErr = err
	}

	err = mgr.unicastMgr.Stop()
	if err != nil {
		lastErr = err
	}

	return lastErr
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
