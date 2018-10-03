// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	errorUnicastServerNotRunning           = "Unicast server is not running"
	errorUnicastServerNoAvailableInterface = "No available interface"
)

// A UnicastManager represents a multicast server manager.
type UnicastManager struct {
	*Config
	Port    int
	Servers []*UnicastServer
	Handler UnicastHandler
}

// NewUnicastManager returns a new UnicastManager.
func NewUnicastManager() *UnicastManager {
	mgr := &UnicastManager{
		Config:  NewDefaultConfig(),
		Port:    UDPPort,
		Servers: make([]*UnicastServer, 0),
		Handler: nil,
	}
	return mgr
}

// SetPort sets a listen port.
func (mgr *UnicastManager) SetPort(port int) {
	mgr.Port = port
}

// GetPort returns the listen port.
func (mgr *UnicastManager) GetPort() int {
	return mgr.Port
}

// SetHandler set a listener to all servers.
func (mgr *UnicastManager) SetHandler(l UnicastHandler) {
	mgr.Handler = l
}

// GetBoundAddresses returns the listen addresses.
func (mgr *UnicastManager) GetBoundAddresses() []string {
	boundAddrs := make([]string, 0)
	for _, server := range mgr.Servers {
		boundAddrs = append(boundAddrs, server.GetBoundAddresses()...)
	}
	return boundAddrs
}

// GetBoundInterfaces returns the listen interfaces.
func (mgr *UnicastManager) GetBoundInterfaces() []net.Interface {
	boundIfs := make([]net.Interface, 0)
	for _, server := range mgr.Servers {
		boundIfs = append(boundIfs, server.GetBoundInterface())
	}
	return boundIfs
}

// Start starts this server.
func (mgr *UnicastManager) Start(ifi net.Interface) (*UnicastServer, error) {
	server := NewUnicastServer()
	server.SetConfig(mgr.Config.UnicastConfig)
	server.Handler = mgr.Handler

	startPort := mgr.GetPort()
	endPort := startPort
	if !mgr.IsAutoBindingEnabled() {
		endPort = UDPPortMax
	}

	var lastErr error
	for port := startPort; port <= endPort; port++ {
		err := server.Start(ifi, port)
		if err != nil {
			lastErr = err
			continue
		}

		if err == nil {
			mgr.SetPort(port)
			lastErr = nil
			break
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}

	mgr.Servers = append(mgr.Servers, server)

	return server, nil
}

// Stop stops this server.
func (mgr *UnicastManager) Stop() error {
	var lastErr error
	for _, server := range mgr.Servers {
		if server == nil {
			continue
		}
		err := server.Stop()
		if err != nil {
			lastErr = err
		}
	}
	mgr.Servers = make([]*UnicastServer, 0)
	return lastErr
}

// IsRunning returns true whether the local servers are running, otherwise false.
func (mgr *UnicastManager) IsRunning() bool {
	if len(mgr.Servers) <= 0 {
		return false
	}
	return true
}

// SendMessage sends a message to the destination address.
func (mgr *UnicastManager) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	var lastErr error
	for _, server := range mgr.Servers {
		n, err := server.SendMessage(addr, port, msg)
		if err == nil {
			return n, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return 0, lastErr
	}
	return 0, fmt.Errorf(errorUnicastServerNotRunning)
}

// AnnounceMessage sends a message to the multicast address.
func (mgr *UnicastManager) AnnounceMessage(addr string, port int, msg *protocol.Message) error {
	var lastErr error
	for _, server := range mgr.Servers {
		err := server.AnnounceMessage(addr, port, msg)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf(errorUnicastServerNotRunning)
}

// PostMessage posts a message to the destination address and gets the response message.
func (mgr *UnicastManager) PostMessage(addr string, port int, reqMsg *protocol.Message) (*protocol.Message, error) {
	if !mgr.IsTCPEnabled() {
		return nil, fmt.Errorf(errorTCPSocketDisabled)
	}

	var lastErr error
	for _, server := range mgr.Servers {
		resMsg, err := server.TCPSocket.PostMessage(addr, port, reqMsg, mgr.GetConnectionTimeout())
		if err == nil {
			return resMsg, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf(errorUnicastServerNotRunning)
}
