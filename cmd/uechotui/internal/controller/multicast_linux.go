//go:build linux

package controller

import (
	"net"
	"syscall"
)

// Linux otherwise delivers groups joined by other sockets/interfaces too.
func restrictMulticast(conn *net.UDPConn) error {
	raw, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	var e error
	err = raw.Control(func(fd uintptr) { e = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, 49, 0) })
	if err != nil {
		return err
	}
	return e
}
