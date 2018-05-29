// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// A UnicastServerManager represents a multicast server manager.
type UnicastServerManager struct {
	Listener UnicastListener
	Servers  []*UnicastServer
}

// NewUnicastServerManager returns a new UnicastServerManager.
func NewUnicastServerManager() *UnicastServerManager {
	server := &UnicastServerManager{}
	server.Servers = make([]*UnicastServer, 0)
	server.Listener = nil
	return server
}

// Start starts this server.
func (mgr *UnicastServerManager) Start() error {
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
func (mgr *UnicastServerManager) Stop() error {
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
