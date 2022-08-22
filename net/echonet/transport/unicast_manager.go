// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"
	"net"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A UnicastManager represents a multicast server manager.
type UnicastManager struct {
	*Config
	port    int
	Servers []*UnicastServer
	Handler UnicastHandler
}

// NewUnicastManager returns a new UnicastManager.
func NewUnicastManager() *UnicastManager {
	mgr := &UnicastManager{
		Config:  NewDefaultConfig(),
		port:    UDPPort,
		Servers: make([]*UnicastServer, 0),
		Handler: nil,
	}
	return mgr
}

// SetPort sets a listen port.
func (mgr *UnicastManager) SetPort(port int) {
	mgr.port = port
}

// GetPort returns the listen port.
func (mgr *UnicastManager) Port() int {
	return mgr.port
}

// SetHandler set a listener to all servers.
func (mgr *UnicastManager) SetHandler(l UnicastHandler) {
	mgr.Handler = l
}

// StartWithInterfaceAndPort starts this server on the specified interface and port.
func (mgr *UnicastManager) StartWithInterfaceAndPort(ifi *net.Interface, ifaddr string, port int) (*UnicastServer, error) {
	server := NewUnicastServer()
	server.SetConfig(mgr.Config.UnicastConfig)
	server.Handler = mgr.Handler
	if err := server.Start(ifi, ifaddr, port); err != nil {
		return nil, err
	}
	mgr.Servers = append(mgr.Servers, server)
	return server, nil
}

// Start starts servers on the all avairable interfaces.
func (mgr *UnicastManager) Start() error {
	if err := mgr.Stop(); err != nil {
		return err
	}

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return err
	}

	var lastErr error

	startPort := mgr.Port()
	endPort := startPort
	if mgr.IsAutoPortBindingEnabled() {
		endPort = startPort + UDPPortRange
	}

	for port := startPort; port <= endPort; port++ {
		bindRetryCount := uint(0)
		if !mgr.IsAutoPortBindingEnabled() {
			bindRetryCount = mgr.BindRetryCount()
		}

		for n := uint(0); n <= bindRetryCount; n++ {
			for _, ifi := range ifis {
				ifaddrs, err := GetInterfaceAddresses(ifi)
				if err != nil {
					continue
				}
				for _, ifaddr := range ifaddrs {
					_, lastErr = mgr.StartWithInterfaceAndPort(ifi, ifaddr, port)
					if lastErr != nil {
						break
					}
				}
			}
			if lastErr == nil {
				break
			}
			if n < bindRetryCount {
				time.Sleep(mgr.BindRetryWaitTime())
			}
		}
		if lastErr == nil {
			mgr.SetPort(port)
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

// Stop stops this server.
func (mgr *UnicastManager) getAppropriateServerForInterface(ifi *net.Interface) (*UnicastServer, error) {
	if len(mgr.Servers) == 0 {
		return nil, fmt.Errorf(errorUnicastServerNotRunning)
	}

	for _, server := range mgr.Servers {
		if server == nil {
			continue
		}
		if server.UDPSocket.interfac == ifi {
			return server, nil
		}
	}

	return mgr.Servers[0], nil
}

// IsRunning returns true whether the local servers are running, otherwise false.
func (mgr *UnicastManager) IsRunning() bool {
	return len(mgr.Servers) != 0
}

// SendMessage sends a message to the destination address.
func (mgr *UnicastManager) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	var lastErr error
	for _, server := range mgr.Servers {
		n, err := server.SendMessage(addr, port, msg)
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

// AnnounceMessage sends a message to the multicast address.
func (mgr *UnicastManager) AnnounceMessage(msg *protocol.Message) error {
	var lastErr error
	for _, server := range mgr.Servers {
		err := server.AnnounceMessage(msg)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf(errorUnicastServerNotRunning)
}

// PostMessage posts a message to the destination address and gets the response message.
func (mgr *UnicastManager) PostMessage(addr string, port int, reqMsg *protocol.Message) (*protocol.Message, error) {
	if !mgr.IsTCPEnabled() {
		return nil, fmt.Errorf(errorTCPSocketDisabled)
	}

	var lastErr error
	for _, server := range mgr.Servers {
		resMsg, err := server.TCPSocket.PostMessage(addr, port, reqMsg, mgr.ConnectionTimeout())
		if err == nil {
			return resMsg, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf(errorUnicastServerNotRunning)
}
