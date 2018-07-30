// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
)

const (
	errorServerNotRunning = "Unicast erver is not running"
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

// SetListener set a listener to all servers.
func (mgr *UnicastManager) SetListener(l UnicastListener) {
	mgr.Listener = l
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

	mgr.Servers = make([]*UnicastServer, len(ifis))
	for n, ifi := range ifis {
		server := NewUnicastServer()
		server.Listener = mgr.Listener
		err := server.Start(ifi, mgr.Port)
		if err != nil {
			lastErr = err
		}
		mgr.Servers[n] = server
	}

	return lastErr
}

// Stop stops this server.
func (mgr *UnicastManager) Stop() error {
	var lastErr error = nil

	for _, server := range mgr.Servers {
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
