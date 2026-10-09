//go:build !linux

package controller

import "net"

func restrictMulticast(conn *net.UDPConn) error { return nil }
