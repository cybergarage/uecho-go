package controller

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type Value struct {
	Data           []byte
	Sent, Received time.Time
	TID            uint
	State          string
}
type Device struct {
	Target             Target
	Get, Set, Announce []byte
	Values             map[byte]Value
	Maps               bool
}
type Snapshot struct {
	Devices []Device
	Status  string
}
type Session struct {
	operation chan struct{}
	mu        sync.Mutex
	Client    *Client
	devices   map[Target]*Device
	status    string
	timeout   time.Duration
}

func NewSession(c *Client) *Session {
	s := &Session{Client: c, operation: make(chan struct{}, 1), devices: make(map[Target]*Device), status: "Ready — no requests sent", timeout: 3 * time.Second}
	c.OnNotification(func(m *protocol.Message, at time.Time) {
		// Separate notification cache is deliberate: arrival order cannot prove
		// freshness; INF is never readback or a replacement for a Get snapshot.
		cEvent := Event{Time: at, Kind: "INF", From: fmt.Sprintf("%s:3610", m.SourceAddress()), Text: "arrival time only; freshness unknown", Hex: fmt.Sprintf("%X", m.Bytes()), TID: m.TID()}
		c.Log(cEvent)
	})
	return s
}
func (s *Session) Status(text string) { s.mu.Lock(); s.status = text; s.mu.Unlock() }
func (s *Session) Snapshot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := Snapshot{Status: s.status}
	for _, d := range s.devices {
		copyD := *d
		copyD.Get = append([]byte(nil), d.Get...)
		copyD.Set = append([]byte(nil), d.Set...)
		copyD.Announce = append([]byte(nil), d.Announce...)
		copyD.Values = make(map[byte]Value)
		for ep, v := range d.Values {
			v.Data = append([]byte(nil), v.Data...)
			copyD.Values[ep] = v
		}
		out.Devices = append(out.Devices, copyD)
	}
	slices.SortFunc(out.Devices, func(a, b Device) int { return strings.Compare(a.Target.String(), b.Target.String()) })
	return out
}
func (s *Session) Add(t Target) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.devices[t]; !ok {
		s.devices[t] = &Device{Target: t, Values: make(map[byte]Value)}
	}
}
func (s *Session) request(ctx context.Context, t Target, esv protocol.ESV, ep byte, data []byte) (Reply, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.Client.RoundTrip(ctx, t, esv, ep, data)
}
func (s *Session) Discover(ctx context.Context, ip string) error {
	if err := s.acquire(ctx); err != nil {
		return err
	}
	defer s.release()

	s.Status("Sending discovery Get D6 to " + ip)
	r, err := s.request(ctx, Target{ip, 0x0ef001}, 0x62, 0xd6, nil)
	if err != nil {
		return s.fail("Discovery", err)
	}
	b := r.Message.Property(0).Data()
	if len(b) < 1 || len(b) != 1+3*int(b[0]) {
		return s.fail("Discovery", fmt.Errorf("invalid instance list"))
	}
	for i := 1; i < len(b); i += 3 {
		eoj := protocol.ObjectCode(uint32(b[i])<<16 | uint32(b[i+1])<<8 | uint32(b[i+2]))
		if eoj&0xff != 0 {
			s.Add(Target{ip, eoj})
		}
	}
	s.Status(fmt.Sprintf("Discovery success: %d instances, TID %04X at %s", b[0], r.Message.TID(), stamp(r.Received)))
	return nil
}
func stamp(t time.Time) string { return t.UTC().Format("15:04:05.000Z") }
func (s *Session) fail(operation string, err error) error {
	label := err.Error()
	if errors.Is(err, context.DeadlineExceeded) {
		label = "timeout"
	}
	if errors.Is(err, context.Canceled) {
		label = "canceled"
	}
	s.Status(operation + ": " + label)
	return err
}
func (s *Session) Load(ctx context.Context, t Target) error {
	if err := s.acquire(ctx); err != nil {
		return err
	}
	defer s.release()

	s.Add(t)
	// Clear old permissions immediately; a failed refresh must never leave stale
	// write permission enabled. Publish all maps atomically only when validated.
	s.mu.Lock()
	d := s.devices[t]
	d.Maps = false
	d.Get = nil
	d.Set = nil
	d.Announce = nil
	s.mu.Unlock()
	maps := make(map[byte][]byte)
	for _, ep := range []byte{0x9d, 0x9e, 0x9f} {
		s.Status(fmt.Sprintf("Sending property map Get %02X", ep))
		r, err := s.request(ctx, t, 0x62, ep, nil)
		if err != nil {
			return s.fail("Property map", err)
		}
		data := r.Message.Property(0).Data()
		p := echonet.NewProperty(echonet.WithPropertyCode(echonet.PropertyCode(ep)), echonet.WithPropertyData(data))
		codes, err := p.PropertyMapData()
		if err != nil {
			return s.fail("Property map", err)
		}
		seen := map[byte]bool{}
		for _, code := range codes {
			if code < 0x80 || code > 0xff || seen[byte(code)] {
				return s.fail("Property map", fmt.Errorf("invalid/duplicate EPC"))
			}
			seen[byte(code)] = true
			maps[ep] = append(maps[ep], byte(code))
		}
		if len(codes) != int(data[0]) {
			return s.fail("Property map", fmt.Errorf("property count mismatch"))
		}
		s.save(t, ep, r, "Get success")
	}
	s.mu.Lock()
	d = s.devices[t]
	d.Get = maps[0x9f]
	d.Set = maps[0x9e]
	d.Announce = maps[0x9d]
	d.Maps = true
	s.mu.Unlock()
	for _, ep := range maps[0x9f] {
		if err := s.get(ctx, t, ep); err != nil {
			return err
		}
	}
	s.Status("Property maps and fresh Get values loaded: " + t.String())
	return nil
}
func (s *Session) save(t Target, ep byte, r Reply, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices[t].Values[ep] = Value{Data: append([]byte(nil), r.Message.Property(0).Data()...), Sent: r.Sent, Received: r.Received, TID: r.Message.TID(), State: state}
}
func (s *Session) Get(ctx context.Context, t Target, ep byte) error {
	if err := s.acquire(ctx); err != nil {
		return err
	}
	defer s.release()

	return s.get(ctx, t, ep)
}
func (s *Session) get(ctx context.Context, t Target, ep byte) error {
	s.Add(t)
	s.Status(fmt.Sprintf("Sending Get %02X to %s", ep, t))
	r, err := s.request(ctx, t, 0x62, ep, nil)
	if err != nil {
		s.mu.Lock()
		v := s.devices[t].Values[ep]
		v.State = "Get unknown: " + err.Error()
		s.devices[t].Values[ep] = v
		s.mu.Unlock()
		return s.fail("Get", err)
	}
	s.save(t, ep, r, "Get success")
	s.Status(fmt.Sprintf("Get success EPC %02X TID %04X at %s", ep, r.Message.TID(), stamp(r.Received)))
	return nil
}

// Set performs one SetC then a distinct fresh Get. Cancellation/timeout after
// transmission leaves outcome unknown; it does not retry a possibly applied Set.
func (s *Session) Set(ctx context.Context, t Target, ep byte, data []byte) error {
	if err := s.acquire(ctx); err != nil {
		return err
	}
	defer s.release()

	s.mu.Lock()
	d := s.devices[t]
	allowed := d != nil && d.Maps && slices.Contains(d.Set, ep)
	s.mu.Unlock()
	if !allowed {
		return s.fail("SetC", fmt.Errorf("not writable in a freshly loaded Set map"))
	}
	if len(data) == 0 || len(data) > 255 {
		return s.fail("SetC", fmt.Errorf("EDT must contain 1..255 bytes"))
	}
	s.Status(fmt.Sprintf("Sending SetC %02X=%X to %s", ep, data, t))
	ack, err := s.request(ctx, t, 0x61, ep, data)
	if err != nil {
		s.mark(t, ep, "SetC unknown / possibly applied")
		if errors.Is(err, ErrRejected) {
			s.mark(t, ep, "SetC rejected")
			return s.fail("SetC rejected", err)
		}
		return s.fail("SetC outcome unknown / possibly applied", err)
	}
	s.Status(fmt.Sprintf("SetC acknowledged TID %04X; sending fresh Get readback", ack.Message.TID()))
	r, err := s.request(ctx, t, 0x62, ep, nil)
	if err != nil {
		s.mark(t, ep, "SetC acknowledged; readback unknown")
		return s.fail("SetC acknowledged; readback unknown", err)
	}
	state := "Readback success"
	if !bytes.Equal(data, r.Message.Property(0).Data()) {
		state = "Readback mismatch"
	}
	s.save(t, ep, r, state)
	s.Status(fmt.Sprintf("%s: SetC TID %04X / Get TID %04X RX %s EDT %X", state, ack.Message.TID(), r.Message.TID(), stamp(r.Received), r.Message.Property(0).Data()))
	return nil
}
func ParseHex(text string) ([]byte, error) {
	compact := strings.Join(strings.Fields(text), "")
	if compact == "" || len(compact)%2 != 0 || len(compact) > 510 {
		return nil, fmt.Errorf("enter 1..255 bytes as even-length raw hex")
	}
	b, err := hex.DecodeString(compact)
	if err != nil {
		return nil, fmt.Errorf("invalid raw hex")
	}
	return b, nil
}

func (s *Session) mark(t Target, ep byte, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d := s.devices[t]; d != nil {
		v := d.Values[ep]
		v.State = state
		d.Values[ep] = v
	}
}

func (s *Session) acquire(ctx context.Context) error {
	select {
	case s.operation <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (s *Session) release() { <-s.operation }
