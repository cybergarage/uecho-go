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

const (
	errorUnicastServerNotRunning           = "Unicast server is not running"
	errorUnicastServerNoAvailableInterface = "No available interface"
)

// A UnicastManager represents a multicast server manager.
type UnicastManager struct {
	*Config
	Port    int
	Servers []*UnicastServer
	Handler UnicastHandler
}

// NewUnicastManager returns a new UnicastManager.
func NewUnicastManager() *UnicastManager {
	mgr := &UnicastManager{
		Config:  NewDefaultConfig(),
		Port:    UDPPort,
		Servers: make([]*UnicastServer, 0),
		Handler: nil,
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

// SetHandler set a listener to all servers.
func (mgr *UnicastManager) SetHandler(l UnicastHandler) {
	mgr.Handler = l
}

// GetBoundAddresses returns the listen addresses.
func (mgr *UnicastManager) GetBoundAddresses() []string {
	boundAddrs := make([]string, 0)

	if mgr.IsEachInterfaceBindingEnabled() {
		for _, server := range mgr.Servers {
			boundAddrs = append(boundAddrs, server.GetBoundAddresses()...)
		}
	} else {
		addrs, err := GetAvailableAddresses()
		if err == nil {
			boundAddrs = append(boundAddrs, addrs...)
		}
	}

	return boundAddrs
}

// GetBoundInterfaces returns the listen interfaces.
func (mgr *UnicastManager) GetBoundInterfaces() []*net.Interface {
	boundIfs := make([]*net.Interface, 0)

	if mgr.IsEachInterfaceBindingEnabled() {
		for _, server := range mgr.Servers {
			boundIfs = append(boundIfs, server.GetBoundInterface())
		}
	} else {
		ifis, err := GetAvailableInterfaces()
		if err == nil {
			boundIfs = append(boundIfs, ifis...)
		}
	}

	return boundIfs
}

// StartWithInterfaceAndPort starts this server on the specified interface and port.
func (mgr *UnicastManager) StartWithInterfaceAndPort(ifi *net.Interface, port int) (*UnicastServer, error) {
	server := NewUnicastServer()
	server.SetConfig(mgr.Config.UnicastConfig)
	server.Handler = mgr.Handler

	err := server.Start(ifi, port)
	if err != nil {
		return nil, err
	}

	mgr.Servers = append(mgr.Servers, server)

	return server, nil
}

// StartWithInterface starts this server on the specified interface.
func (mgr *UnicastManager) StartWithInterface(ifi *net.Interface) (*UnicastServer, error) {
	startPort := mgr.GetPort()
	endPort := startPort
	if mgr.IsAutoPortBindingEnabled() {
		endPort = startPort + UDPPortRange
	}

	var lastErr error
	for port := startPort; port <= endPort; port++ {
		server, lastErr := mgr.StartWithInterfaceAndPort(ifi, port)
		if lastErr == nil {
			mgr.SetPort(port)
			return server, nil
		}
	}

	return nil, lastErr
}

// Start starts servers on the all avairable interfaces.
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

	startPort := mgr.GetPort()
	endPort := startPort
	if mgr.IsAutoPortBindingEnabled() {
		endPort = startPort + UDPPortRange
	}

	for port := startPort; port <= endPort; port++ {

		bindRetryCount := uint(0)
		if !mgr.IsAutoPortBindingEnabled() {
			bindRetryCount = mgr.GetBindRetryCount()
		}

		for n := uint(0); n <= bindRetryCount; n++ {
			if mgr.IsEachInterfaceBindingEnabled() {
				for _, ifi := range ifis {
					_, lastErr = mgr.StartWithInterfaceAndPort(ifi, port)
					if lastErr != nil {
						break
					}
				}
			} else {
				_, lastErr = mgr.StartWithInterfaceAndPort(nil, port)
			}
			if lastErr == nil {
				break
			}
			if n < bindRetryCount {
				time.Sleep(mgr.GetBindRetryWaitTime())
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
	if len(mgr.Servers) <= 0 {
		return nil, fmt.Errorf(errorUnicastServerNotRunning)
	}

	for _, server := range mgr.Servers {
		if server == nil {
			continue
		}
		if server.Interface == ifi {
			return server, nil
		}
	}

	return mgr.Servers[0], nil
}

// IsRunning returns true whether the local servers are running, otherwise false.
func (mgr *UnicastManager) IsRunning() bool {
	if len(mgr.Servers) <= 0 {
		return false
	}
	return true
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
func (mgr *UnicastManager) AnnounceMessage(addr string, port int, msg *protocol.Message) error {
	var lastErr error
	for _, server := range mgr.Servers {
		err := server.AnnounceMessage(addr, port, msg)
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
		resMsg, err := server.TCPSocket.PostMessage(addr, port, reqMsg, mgr.GetConnectionTimeout())
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
