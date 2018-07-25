// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// A MulticastManager represents a multicast server manager.
type MulticastManager struct {
	Listener MulticastListener
	Servers  []*MulticastServer
}

// NewMulticastManager returns a new MulticastManager.
func NewMulticastManager() *MulticastManager {
	server := &MulticastManager{}
	server.Servers = make([]*MulticastServer, 0)
	server.Listener = nil
	return server
}

// SetListener set a listener to all servers.
func (mgr *MulticastManager) SetListener(l MulticastListener) {
	mgr.Listener = l
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
func (mgr *MulticastManager) Stop() error {
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
