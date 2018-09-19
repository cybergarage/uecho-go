// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
)

const (
	errorUnicastServerNotRunning           = "Unicast server is not running"
	errorUnicastServerNoAvailableInterface = "No available interface"
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
func (mgr *UnicastManager) Start(conf *Config) error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	mgr.Servers = make([]*UnicastServer, 0)

	var lastErr error
	for port := mgr.GetPort(); (UDPPortMin <= port) && (port <= UDPPortMax); port++ {
		mgr.SetPort(port)

		for _, ifi := range ifis {
			server := NewUnicastServer()
			server.Listener = mgr.Listener
			lastErr = server.Start(conf, ifi, mgr.Port)

			if lastErr == nil {
				mgr.Servers = append(mgr.Servers, server)
			} else {
				mgr.Stop()
				break
			}
		}

		if lastErr == nil {
			break
		}
	}

	if lastErr != nil {
		return lastErr
	}

	if len(mgr.Servers) <= 0 {
		return fmt.Errorf(errorUnicastServerNoAvailableInterface)
	}

	return nil
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

func (mgr *UnicastManager) Write(addr string, port int, b []byte) (int, error) {
	var lastErr error
	for _, server := range mgr.Servers {
		n, err := server.UDPSocket.Write(addr, port, b)
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
