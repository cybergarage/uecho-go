// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
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
func (mgr *MulticastManager) GetBoundInterfaces() []*net.Interface {
	boundIfs := make([]*net.Interface, 0)
	for _, server := range mgr.Servers {
		boundIfs = append(boundIfs, server.Interface)
	}
	return boundIfs
}

// StartWithInterface starts this server on the specified interface.
func (mgr *MulticastManager) StartWithInterface(ifi *net.Interface) (*MulticastServer, error) {

	server := NewMulticastServer()
	server.Handler = mgr.Handler
	err := server.Start(ifi)

	if err != nil {
		return nil, err
	}

	mgr.Servers = append(mgr.Servers, server)

	return server, nil
}

// Start starts servers on the all avairable interfaces.
func (mgr *MulticastManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	for _, ifi := range ifis {
		_, err := mgr.StartWithInterface(ifi)
		if err != nil {
			mgr.Stop()
			return err
		}
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

// setUnicastManager sets appropriate unicast servers to all multicast servers to response the multicast messages.
func (mgr *MulticastManager) setUnicastManager(unicastMgr *UnicastManager) error {
	for _, multicastServer := range mgr.Servers {
		unicastServer, err := unicastMgr.getAppropriateServerForInterface(multicastServer.Interface)
		if err != nil {
			mgr.Stop()
			return err
		}
		multicastServer.SetUnicastServer(unicastServer)
	}
	return nil
}
