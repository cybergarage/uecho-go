// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
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

// AnnounceMessage announces the message to the bound multicast address.
func (mgr *MulticastManager) AnnounceMessage(msg *protocol.Message) error {
	var lastErr error
	for _, server := range mgr.Servers {
		err := server.AnnounceMessage(msg)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// StartWithInterface starts this server on the specified interface.
func (mgr *MulticastManager) StartWithInterface(ifi *net.Interface, ifaddr string) (*MulticastServer, error) {
	server := NewMulticastServer()
	server.Handler = mgr.Handler
	if err := server.Start(ifi, ifaddr); err != nil {
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
		ifaddrs, err := GetInterfaceAddresses(ifi)
		if err != nil {
			continue
		}
		for _, ifaddr := range ifaddrs {
			_, err := mgr.StartWithInterface(ifi, ifaddr)
			if err != nil {
				mgr.Stop()
				return err
			}
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
	return len(mgr.Servers) != 0
}

// setUnicastManager sets appropriate unicast servers to all multicast servers to response the multicast messages.
func (mgr *MulticastManager) setUnicastManager(unicastMgr *UnicastManager) error {
	for _, multicastServer := range mgr.Servers {
		unicastServer, err := unicastMgr.getAppropriateServerForInterface(multicastServer.Socket.interfac)
		if err != nil {
			mgr.Stop()
			return err
		}
		multicastServer.SetUnicastServer(unicastServer)
	}
	return nil
}
