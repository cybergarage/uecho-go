package controller

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

type discovery struct {
	tid   uint
	sent  time.Time
	until time.Time
	seen  map[string]bool
	count int
}

func validInstances(b []byte) bool {
	if len(b) < 1 || len(b) != 1+3*int(b[0]) {
		return false
	}
	for i := 1; i < len(b); i += 3 {
		if b[i+2] == 0 || b[i] == 0x0e && b[i+1] == 0xf0 {
			return false
		}
	}
	return true
}
func discoveryMatch(p *discovery, m *protocol.Message, now time.Time) bool {
	ip := net.ParseIP(m.SourceAddress())
	return !now.Before(p.sent) && !now.After(p.until) && m.TID() == p.tid && ip != nil && ip.To4() != nil && !ip.IsUnspecified() && !ip.IsMulticast() && m.SourcePort() == 3610 && m.SEOJ() == 0x0ef001 && m.DEOJ() == SourceEOJ && m.ESV() == 0x72 && m.OPC() == 1 && m.Property(0).Code() == 0xd6 && validInstances(m.Property(0).Data())
}
func (c *Client) OnInstances(f func(string, []byte, time.Time)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances = f
}

// Discover uses the selected interface, one fresh TID, a bounded collection
// window, and unicast Get responses. No periodic advertisements or retries.
func (c *Client) Discover(ctx context.Context) (int, error) {
	if _, ok := ctx.Deadline(); !ok {
		return 0, fmt.Errorf("discovery requires a deadline")
	}
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return 0, ErrClosed
	}
	if c.discovery != nil {
		c.mu.Unlock()
		return 0, fmt.Errorf("discovery already active")
	}
	if c.next >= 65535 {
		c.mu.Unlock()
		return 0, fmt.Errorf("TID space exhausted; restart controller")
	}
	c.next++
	m := Request(0x0ef000, 0x62, 0xd6, nil)
	m.SetTID(c.next)
	until, _ := ctx.Deadline()
	p := &discovery{tid: c.next, sent: time.Now(), until: until, seen: map[string]bool{}}
	c.discovery = p
	c.mu.Unlock()
	defer func() { c.mu.Lock(); c.discovery = nil; c.mu.Unlock() }()
	c.Log(Event{Time: p.sent, Kind: "TX", From: fmt.Sprintf("EOJ %06X", uint(SourceEOJ)), To: transport.MulticastIPv4Address + ":3610", TID: p.tid, Hex: fmt.Sprintf("%X", m.Bytes()), Text: "interface-scoped discovery window"})
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return 0, ErrClosed
	}
	if c.unicast != nil {
		deadline, _ := ctx.Deadline()
		c.unicast.Conn.SetWriteDeadline(deadline)
	}
	err := c.send(Target{transport.MulticastIPv4Address, 0x0ef000}, m)
	c.mu.Unlock()
	if err != nil {
		return 0, err
	}
	select {
	case <-c.stop:
		return 0, ErrClosed
	case <-ctx.Done():
	}
	c.mu.Lock()
	n := p.count
	c.mu.Unlock()
	if ctx.Err() == context.Canceled {
		return n, ctx.Err()
	}
	return n, nil
}

// OpenMulticast validates the explicit address/interface without sending, binds a
// reusable standard-port public transport socket, and joins only that interface.
// Opening receives notifications; sending discovery still requires UI confirmation.
func OpenMulticast(ifaceName, bind string) (*Client, error) {
	iface, err := validateBinding(ifaceName, bind)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(bind).To4()
	c := newClient()
	c.local = ip.String()
	c.unicast = transport.NewUnicastUDPSocket()
	config := net.ListenConfig{Control: func(network, address string, raw syscall.RawConn) error {
		var e error
		if err := raw.Control(func(fd uintptr) { e = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1) }); err != nil {
			return err
		}
		return e
	}}
	conn, err := config.ListenPacket(context.Background(), "udp4", net.JoinHostPort(bind, "3610"))
	if err != nil {
		return nil, err
	}
	c.unicast.Conn = conn.(*net.UDPConn)
	c.unicast.SetBoundStatus(iface, bind, 3610)
	raw, err := c.unicast.Conn.SyscallConn()
	if err == nil {
		var e error
		err = raw.Control(func(fd uintptr) {
			var addr [4]byte
			copy(addr[:], ip)
			e = syscall.SetsockoptInet4Addr(int(fd), syscall.IPPROTO_IP, syscall.IP_MULTICAST_IF, addr)
			if e == nil {
				e = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MULTICAST_TTL, 1)
			}
		})
		if err == nil {
			err = e
		}
	}
	if err != nil {
		c.Close()
		return nil, err
	}
	c.multicast = transport.NewMulticastSocket()
	if err = c.multicast.Bind(iface, bind); err != nil {
		if c.multicast.Conn != nil {
			c.multicast.Conn.Close()
		}
		c.multicast = nil
		c.Close()
		return nil, err
	}
	if err = restrictMulticast(c.multicast.Conn); err != nil {
		c.Close()
		return nil, err
	}
	c.send = func(t Target, m *protocol.Message) error { _, e := c.unicast.SendMessage(t.IP, 3610, m); return e }
	c.read(c.unicast.UDPSocket)
	c.read(c.multicast.UDPSocket)
	return c, nil
}

// InterfaceOptions is read-only. It never chooses an interface for the user.
func InterfaceOptions() ([]InterfaceOption, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	out := []InterfaceOption{}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 {
			continue
		}
		addrs, e := iface.Addrs()
		if e != nil {
			continue
		}
		for _, a := range addrs {
			ip, _, e := net.ParseCIDR(a.String())
			if e == nil && ip.To4() != nil && !ip.IsUnspecified() {
				out = append(out, InterfaceOption{iface.Name, ip.String()})
			}
		}
	}
	return out, nil
}

type InterfaceOption struct{ Name, IP string }
