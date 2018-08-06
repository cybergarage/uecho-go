// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
)

const (
	errorServerNotRunning = "Unicast server is not running"
)

// A UnicastManager represents a multicast server manager.
type UnicastManager struct {
	Port     int
	Servers  []*UnicastServer
	Listener UnicastListener
}

// NewUnicastManager returns a new UnicastManager.
func NewUnicastManager() *UnicastManager {
	mgr := &UnicastManager{
		Port:     UDPPort,
		Servers:  make([]*UnicastServer, 0),
		Listener: nil,
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

// SetListener set a listener to all servers.
func (mgr *UnicastManager) SetListener(l UnicastListener) {
	mgr.Listener = l
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
func (mgr *UnicastManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	var lastErr error

	for port := mgr.GetPort(); (UDPPortMin <= port) && (port <= UDPPortMax); port++ {
		mgr.Servers = make([]*UnicastServer, len(ifis))
		mgr.SetPort(port)

		for n, ifi := range ifis {
			server := NewUnicastServer()
			server.Listener = mgr.Listener
			lastErr = server.Start(ifi, mgr.Port)

			if lastErr == nil {
				mgr.Servers[n] = server
			} else {
				mgr.Stop()
				break
			}
		}

		if lastErr == nil {
			break
		}
	}

	return lastErr
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

func (mgr *UnicastManager) Write(addr string, port int, b []byte) (int, error) {
	if 0 < len(mgr.Servers) {
		server := mgr.Servers[0]
		return server.Socket.Write(addr, port, b)
	}
	return 0, fmt.Errorf(errorServerNotRunning)
}
