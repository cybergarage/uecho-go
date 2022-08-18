// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"errors"
	"net"
	"strings"
)

const (
	libvirtInterfaceName = "virbr0"
)

// IsIPv6Interface returns true whether the specified interface has a IPv6 address.
func IsIPv6Address(addr string) bool {
	if len(addr) == 0 {
		return false
	}
	if 0 <= strings.Index(addr, ":") {
		return true
	}
	return false
}

// IsIPv4Address returns true whether the specified address is a IPv4 address.
func IsIPv4Address(addr string) bool {
	if len(addr) == 0 {
		return false
	}
	return !IsIPv6Address(addr)
}

// IsIPv6Interface returns true whether the specified address is a IPv6 address.
func IsIPv6Interface(ifi *net.Interface) bool {
	addrs, err := GetInterfaceAddresses(ifi)
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		if IsIPv6Address(addr) {
			return true
		}
	}
	return false
}

// IsIPv4Interface returns true whether the specified address is a IPv4 address.
func IsIPv4Interface(ifi *net.Interface) bool {
	addrs, err := GetInterfaceAddresses(ifi)
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		if IsIPv4Address(addr) {
			return true
		}
	}
	return false
}

// IsLoopbackAddress returns true whether the specified address is a loopback addresses.
func IsLoopbackAddress(addr string) bool {
	localAddrs := []string{
		"127.0.0.1",
		"::1",
	}
	for _, localAddr := range localAddrs {
		if localAddr == addr {
			return true
		}
	}
	return false
}

// IsCommunicableAddress returns true whether the address is a effective address to commnicate with other nodes, othwise false.
func IsCommunicableAddress(addr string) bool {
	if len(addr) == 0 {
		return false
	}
	if IsLoopbackAddress(addr) {
		return false
	}
	return true
}

// IsBridgeInterface returns true when the specified interface is a bridege interface, otherwise false.
func IsBridgeInterface(ifi *net.Interface) bool {
	return ifi.Name == libvirtInterfaceName
}

// IsVirtualInterface returns true when the specified interface is a virtual interface, otherwise false.
func IsVirtualInterface(ifi *net.Interface) bool {
	if strings.HasPrefix(ifi.Name, "utun") { // macOS
		return true
	}
	if strings.HasPrefix(ifi.Name, "llw") { // VirtualBox
		return true
	}
	if strings.HasPrefix(ifi.Name, "awdl") { // AirDrop (macOS)
		return true
	}
	if strings.HasPrefix(ifi.Name, "en6") { // iPhone-USB (macOS)
		return true
	}
	return false
}

// GetInterfaceAddresses returns a IPv4 or IPv6 address of the specivied interface.
func GetInterfaceAddresses(ifi *net.Interface) ([]string, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return nil, err
	}
	ipaddrs := []string{}
	for _, addr := range addrs {
		addrStr := addr.String()
		saddr := strings.Split(addrStr, "/")
		if len(saddr) < 2 {
			continue
		}
		if IsIPv6Address(saddr[0]) {
			continue
		}
		ipaddrs = append(ipaddrs, saddr[0])
	}
	if len(ipaddrs) == 0 {
		return nil, errors.New(errorAvailableAddressNotFound)
	}
	return ipaddrs, nil
}

// GetAvailableInterfaces returns all available interfaces in the node.
func GetAvailableInterfaces() ([]*net.Interface, error) {
	useIfs := make([]*net.Interface, 0)
	localIfs, err := net.Interfaces()
	if err != nil {
		return useIfs, err
	}

	for n := range localIfs {
		localIf := localIfs[n]
		if (localIf.Flags & net.FlagLoopback) != 0 {
			continue
		}
		if (localIf.Flags & net.FlagUp) == 0 {
			continue
		}
		if (localIf.Flags & net.FlagMulticast) == 0 {
			continue
		}
		if IsBridgeInterface(&localIf) {
			continue
		}
		if IsVirtualInterface(&localIf) {
			continue
		}
		_, addrErr := GetInterfaceAddresses(&localIf)
		if addrErr != nil {
			continue
		}

		useIf := localIf
		useIfs = append(useIfs, &useIf)
	}

	if len(useIfs) == 0 {
		return useIfs, errors.New(errorAvailableInterfaceFound)
	}

	return useIfs, err
}

// GetAvailableAddresses returns all available IP addresses in the node.
func GetAvailableAddresses() ([]string, error) {
	addrs := make([]string, 0)
	ifis, err := GetAvailableInterfaces()
	if err != nil {
		return addrs, err
	}
	for _, ifi := range ifis {
		ipaddrs, err := GetInterfaceAddresses(ifi)
		if err != nil {
			continue
		}
		addrs = append(addrs, ipaddrs...)
	}
	return addrs, nil
}
