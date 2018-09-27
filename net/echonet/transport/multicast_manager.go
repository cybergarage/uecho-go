// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
)

const (
	errorMulticastServerNoAvailableInterface = "No available interface"
)

// A MulticastManager represents a multicast server manager.
type MulticastManager struct {
	Servers []*MulticastServer
	Handler MulticastHandler
}

// NewMulticastManager returns a new MulticastManager.
func NewMulticastManager() *MulticastManager {
	mgr := &MulticastManager{
		Servers: make([]*MulticastServer, 0),
		Handler: nil,
	}
	return mgr
}

// SetHandler set a listener to all servers.
func (mgr *MulticastManager) SetHandler(l MulticastHandler) {
	mgr.Handler = l
}

// GetBoundAddresses returns the listen addresses.
func (mgr *MulticastManager) GetBoundAddresses() []string {
	boundAddrs := make([]string, 0)
	for _, server := range mgr.Servers {
		boundAddrs = append(boundAddrs, server.GetBoundAddresses()...)
	}
	return boundAddrs
}

// GetBoundInterfaces returns the listen interfaces.
func (mgr *MulticastManager) GetBoundInterfaces() []net.Interface {
	boundIfs := make([]net.Interface, 0)
	for _, server := range mgr.Servers {
		boundIfs = append(boundIfs, server.Interface)
	}
	return boundIfs
}

// Start starts this server.
func (mgr *MulticastManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	mgr.Servers = make([]*MulticastServer, 0)

	var lastErr error
	for _, ifi := range ifis {
		server := NewMulticastServer()
		server.Handler = mgr.Handler
		err := server.Start(ifi)
		if err != nil {
			lastErr = err
		}
		mgr.Servers = append(mgr.Servers, server)
	}

	if lastErr != nil {
		return lastErr
	}

	if len(mgr.Servers) <= 0 {
		return fmt.Errorf(errorMulticastServerNoAvailableInterface)
	}

	return nil
}

// Stop stops this server.
func (mgr *MulticastManager) Stop() error {
	var lastErr error

	for _, server := range mgr.Servers {
		err := server.Stop()
		if err != nil {
			lastErr = err
		}
	}

	mgr.Servers = make([]*MulticastServer, 0)

	return lastErr
}

// IsRunning returns true whether the local servers are running, otherwise false.
func (mgr *MulticastManager) IsRunning() bool {
	if len(mgr.Servers) <= 0 {
		return false
	}
	return true
}
