// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"net"
	"slices"
	"strings"
)

// IsIPv6Address returns true whether the specified interface has a IPv6 address.
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
	return slices.ContainsFunc(addrs, IsIPv6Address)
}

// IsIPv4Interface returns true whether the specified address is a IPv4 address.
func IsIPv4Interface(ifi *net.Interface) bool {
	addrs, err := GetInterfaceAddresses(ifi)
	if err != nil {
		return false
	}
	return slices.ContainsFunc(addrs, IsIPv4Address)
}

// IsLoopbackAddress returns true whether the specified address is a loopback addresses.
func IsLoopbackAddress(addr string) bool {
	localAddrs := []string{
		"127.0.0.1",
		"::1",
	}
	return slices.Contains(localAddrs, addr)
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
	prefixes := []string{
		"virbr0",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(ifi.Name, prefix) {
			return true
		}
	}
	return false
}

// IsVirtualInterface returns true when the specified interface is a virtual interface, otherwise false.
func IsVirtualInterface(ifi *net.Interface) bool {
	prefixes := []string{
		"utun",    // macOS
		"llw",     // VirtualBox
		"awdl",    // AirDrop (macOS)
		"en6",     // iPhone-USB (macOS)
		"docker0", // Docker
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(ifi.Name, prefix) {
			return true
		}
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
		return nil, errAvailableAddressNotFound
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
		return useIfs, errAvailableInterfaceFound
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
