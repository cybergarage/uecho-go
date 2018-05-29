// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// A MulticastServerManager represents a multicast server manager.
type MulticastServerManager struct {
	Listener MulticastListener
	Servers  []*MulticastServer
}

// NewMulticastServerManager returns a new MulticastServerManager.
func NewMulticastServerManager() *MulticastServerManager {
	server := &MulticastServerManager{}
	server.Servers = make([]*MulticastServer, 0)
	server.Listener = nil
	return server
}

// SetListener set a listener to all servers.
func (mgr *MulticastServerManager) SetListener(l UnicastListener) {
	mgr.Listener = l
}

// Start starts this server.
func (mgr *MulticastServerManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	var lastErr error = nil

	mgr.Servers = make([]*MulticastServer, len(ifis))
	for n, ifi := range ifis {
		server := NewMulticastServer()
		server.Listener = mgr.Listener
		err := server.Start(ifi)
		if err != nil {
			lastErr = err
		}
		mgr.Servers[n] = server
	}

	return lastErr
}

// Stop stops this server.
func (mgr *MulticastServerManager) Stop() error {
	var lastErr error = nil

	for _, server := range mgr.Servers {
		err := server.Stop()
		if err != nil {
			lastErr = err
		}
	}

	mgr.Servers = make([]*MulticastServer, 0)

	return lastErr
}
