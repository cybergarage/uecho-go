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
	*UnicastConfig
	Port    int
	Servers []*UnicastServer
	Handler UnicastHandler
}

// NewUnicastManager returns a new UnicastManager.
func NewUnicastManager() *UnicastManager {
	mgr := &UnicastManager{
		UnicastConfig: NewDefaultUnicastConfig(),
		Port:          UDPPort,
		Servers:       make([]*UnicastServer, 0),
		Handler:       nil,
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
	server.Handler = mgr.Handler

	var lastErr error
	for port := mgr.GetPort(); (UDPPortMin <= port) && (port <= UDPPortMax); port++ {
		mgr.SetPort(port)
		server.SetConfig(mgr.UnicastConfig)

		err := server.Start(ifi, mgr.Port)
		if err != nil {
			lastErr = err
			continue
		}

		if err == nil {
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

// SendMessage send a message to the destination address.
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
