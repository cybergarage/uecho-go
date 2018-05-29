// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// A UnicastManager represents a multicast server manager.
type UnicastManager struct {
	Listener UnicastListener
	Servers  []*UnicastServer
}

// NewUnicastManager returns a new UnicastManager.
func NewUnicastManager() *UnicastManager {
	server := &UnicastManager{}
	server.Servers = make([]*UnicastServer, 0)
	server.Listener = nil
	return server
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

	var lastErr error = nil

	mgr.Servers = make([]*UnicastServer, len(ifis))
	for n, ifi := range ifis {
		server := NewUnicastServer()
		server.Listener = mgr.Listener
		err := server.Start(ifi, UDP_PORT)
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
