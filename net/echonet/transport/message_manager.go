// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"

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

// GetMulticastManager returns the multicast manager.
func (mgr *MessageManager) GetMulticastManager() *MulticastManager {
	return mgr.multicastMgr
}

// GetMulticastManager returns the unicast manager.
func (mgr *MessageManager) GeUnicastManager() *UnicastManager {
	return mgr.unicastMgr
}

// GetBoundAddresses return the bounded interface addresses.
func (mgr *MessageManager) GetBoundAddresses() []string {
	ifaddrs := []string{}
	for _, server := range mgr.GetMulticastManager().Servers {
		ifaddrs = append(ifaddrs, server.BoundAddress)
	}
	return ifaddrs
}

// SetConfig sets all configuration flags.
func (mgr *MessageManager) SetConfig(newConfig *Config) {
	mgr.unicastMgr.SetConfig(newConfig)
}

// GetConfig returns all current configurations.
func (mgr *MessageManager) GetConfig() *Config {
	return mgr.unicastMgr.Config
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

// GetBoundPort returns the listen port.
func (mgr *MessageManager) GetBoundPort() (int, error) {
	if !mgr.IsRunning() {
		return 0, fmt.Errorf(errorMessageManagerNotRunning)
	}
	return mgr.unicastMgr.GetPort(), nil
}

// SendMessage send a message to the destination address.
func (mgr *MessageManager) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	return mgr.unicastMgr.SendMessage(addr, port, msg)
}

// AnnounceMessage sends a message to the multicast address.
func (mgr *MessageManager) AnnounceMessage(msg *protocol.Message) error {
	return mgr.unicastMgr.AnnounceMessage(msg)
}

// PostMessage posts a message to the destination address and gets the response message.
func (mgr *MessageManager) PostMessage(addr string, port int, msg *protocol.Message) (*protocol.Message, error) {
	return mgr.unicastMgr.PostMessage(addr, port, msg)
}

// Start starts all transport managers.
func (mgr *MessageManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}
	err = mgr.unicastMgr.Start()
	if err != nil {
		return err
	}

	err = mgr.multicastMgr.Start()
	if err != nil {
		mgr.Stop()
		return err
	}

	// Set appropriate unicast servers to all multicast servers to response the multicast messages
	err = mgr.multicastMgr.setUnicastManager(mgr.unicastMgr)
	if err != nil {
		mgr.Stop()
		return err
	}

	mgr.SetPort(mgr.unicastMgr.GetPort())

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
