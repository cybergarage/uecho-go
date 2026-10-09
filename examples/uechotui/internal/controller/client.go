// Package controller implements bounded, correlated developer requests using the
// public uecho-go transport and protocol APIs. It has no terminal dependency.
package controller

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

const SourceEOJ protocol.ObjectCode = 0x05ff01

var ErrClosed = errors.New("controller closed")
var ErrRejected = errors.New("device rejected request")

type Target struct {
	IP  string
	EOJ protocol.ObjectCode
}

func (t Target) String() string { return fmt.Sprintf("%s:3610 / %06X", t.IP, uint(t.EOJ)) }

type Event struct {
	Time                      time.Time
	Kind, From, To, Text, Hex string
	TID                       uint
}
type Reply struct {
	Message        *protocol.Message
	Sent, Received time.Time
}
type pending struct {
	target  Target
	request *protocol.Message
	sent    time.Time
	until   time.Time
	replies chan Reply
}

// Client never reuses a TID during its lifetime. Restart after 65535 requests;
// delayed UDP responses consequently cannot satisfy a newer request.
type Client struct {
	mu           sync.Mutex
	stop         chan struct{}
	next         uint
	closed       bool
	pending      map[uint]*pending
	events       []Event
	unicast      *transport.UnicastUDPSocket
	multicast    *transport.MulticastSocket
	discovery    *discovery
	instances    func(string, []byte, time.Time)
	local        string
	workers      sync.WaitGroup
	send         func(Target, *protocol.Message) error
	notification func(*protocol.Message, time.Time)
}

func newClient() *Client { return &Client{pending: make(map[uint]*pending), stop: make(chan struct{})} }
func (c *Client) Log(e Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e.Time.IsZero() {
		e.Time = time.Now()
	}
	c.events = append(c.events, e)
	if len(c.events) > 256 {
		c.events = append([]Event(nil), c.events[len(c.events)-256:]...)
	}
}
func (c *Client) Events() []Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]Event(nil), c.events...)
}

// Open binds only an explicitly supplied local IPv4 on its named interface.
// No advertisement is sent and no multicast group is joined.
func Open(ifaceName, bind string) (*Client, error) {
	iface, err := validateBinding(ifaceName, bind)
	if err != nil {
		return nil, err
	}
	c := newClient()
	c.unicast = transport.NewUnicastUDPSocket()
	if err = c.unicast.Bind(iface, bind, 3610); err != nil {
		return nil, err
	}
	c.send = func(t Target, m *protocol.Message) error { _, e := c.unicast.SendMessage(t.IP, 3610, m); return e }
	c.read(c.unicast.UDPSocket)
	return c, nil
}

func validateBinding(ifaceName, bind string) (*net.Interface, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(bind)
	if ip == nil || ip.To4() == nil || ip.IsUnspecified() || ip.IsMulticast() {
		return nil, fmt.Errorf("bind must be a literal local IPv4")
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	found := false
	for _, a := range addrs {
		local, _, e := net.ParseCIDR(a.String())
		if e == nil && local.Equal(ip) {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("bind IP does not belong to interface %s", ifaceName)
	}
	return iface, nil
}

// Read raw datagrams on the connection owned by the public transport socket so
// we can reject truncated/trailing payloads before the permissive library parser.
func (c *Client) read(socket *transport.UDPSocket) {
	c.workers.Add(1)
	go func() {
		defer c.workers.Done()
		b := make([]byte, 65536)
		for {
			n, from, err := socket.Conn.ReadFromUDP(b)
			if err != nil {
				c.mu.Lock()
				closed := c.closed
				c.mu.Unlock()
				if !closed {
					c.Log(Event{Kind: "ERROR", Text: err.Error()})
				}
				return
			}
			now := time.Now()
			m, err := Decode(b[:n])
			if err != nil {
				c.Log(Event{Time: now, Kind: "INVALID", From: from.String(), Hex: fmt.Sprintf("%X", b[:min(n, 512)]), Text: err.Error()})
				continue
			}
			m.From.IP = from.IP
			m.From.Port = from.Port
			if m.SourceAddress() != c.local {
				c.receive(m, now)
			}
		}
	}()
}
func Decode(b []byte) (*protocol.Message, error) {
	if len(b) < 12 || len(b) > 4096 || b[0] != 0x10 || b[1] != 0x81 || b[11] == 0 {
		return nil, fmt.Errorf("invalid Format 1 frame")
	}
	pos := 12
	for range int(b[11]) {
		if pos+2 > len(b) {
			return nil, fmt.Errorf("truncated property")
		}
		pos += 2 + int(b[pos+1])
		if pos > len(b) {
			return nil, fmt.Errorf("truncated EDT")
		}
	}
	if pos != len(b) {
		return nil, fmt.Errorf("trailing bytes / unsupported service format")
	}
	return protocol.NewMessageWithBytes(b)
}
func (c *Client) OnNotification(f func(*protocol.Message, time.Time)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.notification = f
}
func match(p *pending, m *protocol.Message) bool {
	r := p.request
	if m.TID() != r.TID() || m.SourceAddress() != p.target.IP || m.SourcePort() != 3610 || m.DEOJ() != r.SEOJ() || m.SEOJ() != r.DEOJ() {
		return false
	}
	success := protocol.ESV(byte(r.ESV()) + 0x10)
	failure := protocol.ESV(byte(r.ESV()) - 0x10)
	if m.ESV() != success && m.ESV() != failure || m.OPC() != r.OPC() {
		return false
	}
	for i, p := range r.Properties() {
		if m.Property(i).Code() != p.Code() {
			return false
		}
		if m.ESV() == 0x71 && m.Property(i).Size() != 0 {
			return false
		}
		if m.ESV() == 0x72 && m.Property(i).Size() == 0 {
			return false
		}
	}
	return true
}
func (c *Client) receive(m *protocol.Message, now time.Time) {
	c.Log(Event{Time: now, Kind: "RX", From: fmt.Sprintf("%s:%d", m.SourceAddress(), m.SourcePort()), To: fmt.Sprintf("EOJ %06X", uint(m.DEOJ())), Hex: fmt.Sprintf("%X", m.Bytes()), TID: m.TID(), Text: fmt.Sprintf("SEOJ %06X ESV %02X", uint(m.SEOJ()), byte(m.ESV()))})
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	if p := c.pending[m.TID()]; p != nil && !now.Before(p.sent) && (p.until.IsZero() || !now.After(p.until)) && match(p, m) {
		select {
		case p.replies <- Reply{m, p.sent, now}:
		default:
		}
		c.mu.Unlock()
		return
	}
	if c.discovery != nil && discoveryMatch(c.discovery, m, now) {
		p := c.discovery
		key := m.SourceAddress()
		if !p.seen[key] && len(p.seen) < 256 {
			p.seen[key] = true
			p.count++
			instances := c.instances
			c.mu.Unlock()
			if instances != nil {
				instances(key, m.Property(0).Data(), now)
			}
			return
		}
	}
	notification := c.notification
	c.mu.Unlock()
	// INF arrival never verifies a write or replaces a fresh Get value.
	ip := net.ParseIP(m.SourceAddress())
	if ip != nil && ip.To4() != nil && !ip.IsUnspecified() && !ip.IsMulticast() && m.ESV() == 0x73 && m.SourcePort() == 3610 && (m.DEOJ() == SourceEOJ || m.DEOJ() == 0x0ef001 || m.DEOJ() == 0x0ef000) && notification != nil {
		notification(m, now)
	}
}

func Request(eoj protocol.ObjectCode, esv protocol.ESV, epc byte, data []byte) *protocol.Message {
	m := protocol.NewMessage()
	m.SetSEOJ(SourceEOJ)
	m.SetDEOJ(eoj)
	m.SetESV(esv)
	p := protocol.NewPropertyWithCode(protocol.PropertyCode(epc))
	p.SetData(append([]byte(nil), data...))
	m.AddProperty(p)
	return m
}
func (c *Client) RoundTrip(ctx context.Context, t Target, esv protocol.ESV, epc byte, data []byte) (Reply, error) {
	ip := net.ParseIP(t.IP)
	if ip == nil || ip.To4() == nil || ip.IsUnspecified() || ip.IsMulticast() || t.EOJ == 0 {
		return Reply{}, fmt.Errorf("invalid target")
	}
	if esv != 0x61 && esv != 0x62 {
		return Reply{}, fmt.Errorf("only Get and SetC supported")
	}
	if len(data) > 255 {
		return Reply{}, fmt.Errorf("EDT exceeds 255 bytes")
	}
	if err := ctx.Err(); err != nil {
		return Reply{}, err
	}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return Reply{}, ErrClosed
	}
	if c.next >= 65535 {
		c.mu.Unlock()
		return Reply{}, fmt.Errorf("TID space exhausted; restart controller")
	}
	c.next++
	m := Request(t.EOJ, esv, epc, data)
	m.SetTID(c.next)
	deadline, _ := ctx.Deadline()
	p := &pending{target: t, request: m, sent: time.Now(), until: deadline, replies: make(chan Reply, 1)}
	c.pending[c.next] = p
	c.mu.Unlock()
	defer func() { c.mu.Lock(); delete(c.pending, m.TID()); c.mu.Unlock() }()
	c.Log(Event{Time: p.sent, Kind: "TX", From: fmt.Sprintf("EOJ %06X", uint(SourceEOJ)), To: t.String(), Hex: fmt.Sprintf("%X", m.Bytes()), TID: m.TID()})
	// Serialize sending with Close; the receive path never holds this lock while reading.
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return Reply{}, ErrClosed
	}
	if c.unicast != nil {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(3 * time.Second)
		}
		c.unicast.Conn.SetWriteDeadline(deadline)
	}
	err := c.send(t, m)
	c.mu.Unlock()
	if err != nil {
		return Reply{}, err
	}
	select {
	case res := <-p.replies:
		if res.Message == nil {
			return Reply{}, ErrClosed
		}
		if res.Message.ESV() == protocol.ESV(byte(esv)-0x10) {
			return res, fmt.Errorf("%w: ESV %02X", ErrRejected, byte(res.Message.ESV()))
		}
		return res, nil
	case <-ctx.Done():
		return Reply{}, ctx.Err()
	}
}
func (c *Client) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	close(c.stop)
	for _, p := range c.pending {
		select {
		case p.replies <- Reply{}:
		default:
		}
	}
	// Close the connections, retaining socket.Conn until readers finish (the
	// upstream Close mutates it). net.UDPConn permits concurrent read/close.
	if c.unicast != nil {
		c.unicast.Conn.Close()
	}
	if c.multicast != nil {
		c.multicast.Conn.Close()
	}
	c.mu.Unlock()
	c.workers.Wait()
	return nil
}
