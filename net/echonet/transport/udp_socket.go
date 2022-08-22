// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"encoding/hex"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A UDPSocket represents a socket for UDP.
type UDPSocket struct {
	*Socket
	Conn           *net.UDPConn
	ReadBufferSize int
	ReadBuffer     []byte
}

// NewUDPSocket returns a new UDPSocket.
func NewUDPSocket() *UDPSocket {
	sock := &UDPSocket{
		Socket:         NewSocket(),
		Conn:           nil,
		ReadBufferSize: MaxPacketSize,
		ReadBuffer:     make([]byte, 0),
	}
	sock.SetReadBufferSize(MaxPacketSize)
	return sock
}

// SetReadBufferSize sets the read buffer size.
func (sock *UDPSocket) SetReadBufferSize(n int) {
	sock.ReadBufferSize = n
	sock.ReadBuffer = make([]byte, n)
}

// GetReadBufferSize returns the read buffer size.
func (sock *UDPSocket) GetReadBufferSize() int {
	return sock.ReadBufferSize
}

// Close closes the current opened socket.
func (sock *UDPSocket) Close() error {
	if sock.Conn == nil {
		return nil
	}

	// FIXME : sock.Conn.Close() hung up on darwin
	/*
		err := sock.Conn.Close()
		if err != nil {
			return err
		}
	*/
	go sock.Conn.Close()
	time.Sleep(time.Millisecond * 100)

	sock.Conn = nil

	return nil
}

func (sock *UDPSocket) outputReadLog(logLevel log.Level, logType string, msgFrom string, msg string, msgSize int) {
	msgTo, _ := sock.IPAddr()
	outputSocketLog(logLevel, logType, logSocketDirectionRead, msgFrom, msgTo, msg, msgSize)
}

func (sock *UDPSocket) outputWriteLog(logLevel log.Level, msgTo string, msg string, msgSize int) {
	msgFrom, _ := sock.IPAddr()
	outputSocketLog(logLevel, logSocketTypeUDPUnicast, logSocketDirectionWrite, msgFrom, msgTo, msg, msgSize)
}

// SendBytes sends the specified bytes.
func (sock *UDPSocket) SendBytes(addr string, port int, b []byte) (int, error) {
	toAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return 0, err
	}

	// Send from binding port

	if sock.Conn != nil {
		n, err := sock.Conn.WriteToUDP(b, toAddr)
		sock.outputWriteLog(log.LevelTrace, toAddr.String(), hex.EncodeToString(b), n)
		if err != nil {
			log.Error(err.Error())
		}
		return n, err
	}

	// Send from no binding port

	conn, err := net.Dial("udp", toAddr.String())
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	n, err := conn.Write(b)
	sock.outputWriteLog(log.LevelTrace, toAddr.String(), hex.EncodeToString(b), n)
	if err != nil {
		log.Error(err.Error())
	}
	conn.Close()

	return n, err
}

// SendMessage send a message to the destination address.
func (sock *UDPSocket) SendMessage(addr string, port int, msg *protocol.Message) (int, error) {
	return sock.SendBytes(addr, port, msg.Bytes())
}

// AnnounceMessage announces the message to the bound multicast address.
func (sock *UDPSocket) AnnounceMessage(msg *protocol.Message) error {
	ifaddr, err := sock.Address()
	if err != nil {
		return err
	}
	maddr := MulticastIPv4Address
	if IsIPv6Address(ifaddr) {
		maddr = MulticastIPv6Address
	}
	_, err = sock.SendMessage(maddr, Port, msg)
	return err
}

// ReadMessage reads a message from the current opened socket.
func (sock *UDPSocket) ReadMessage() (*protocol.Message, error) {
	if sock.Conn == nil {
		return nil, errors.New(errorSocketClosed)
	}

	n, from, err := sock.Conn.ReadFromUDP(sock.ReadBuffer)
	if err != nil {
		return nil, err
	}

	msg, err := protocol.NewMessageWithBytes(sock.ReadBuffer[:n])
	if err != nil {
		sock.outputReadLog(log.LevelError, logSocketTypeUDPUnicast, (*from).String(), hex.EncodeToString(sock.ReadBuffer[:n]), n)
		log.Error(err.Error())
		return nil, err
	}

	msg.From.IP = from.IP
	msg.From.Port = from.Port
	msg.Interface = sock.Socket.interfac

	return msg, nil
}
