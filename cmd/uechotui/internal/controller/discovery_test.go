package controller

import (
	"context"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
	"net"
	"testing"
	"time"
)

func TestDiscoveryCollectsAndValidates(t *testing.T) {
	c := newClient()
	defer c.Close()
	s := NewSession(c)
	s.timeout = 20 * time.Millisecond
	c.send = func(target Target, req *protocol.Message) error {
		if target.IP != "224.0.23.0" || req.DEOJ() != 0x0ef000 {
			t.Error("not multicast node profile")
		}
		go func() {
			base := response(req, "192.0.2.10", []byte{2, 2, 0x90, 1, 1, 0x30, 1})
			base.SetSEOJ(0x0ef001)
			for _, mutate := range []func(*protocol.Message){func(m *protocol.Message) { m.SetTID(req.TID() + 1) }, func(m *protocol.Message) { m.From.Port = 1234 }, func(m *protocol.Message) { m.SetSEOJ(0x029001) }, func(m *protocol.Message) { m.SetDEOJ(0x029001) }, func(m *protocol.Message) { m.SetESV(0x73) }, func(m *protocol.Message) { m.Property(0).SetCode(0xd5) }, func(m *protocol.Message) { m.Property(0).SetData([]byte{2, 2, 0x90, 1}) }} {
				m, _ := protocol.NewMessageWithBytes(base.Bytes())
				m.From.IP = net.ParseIP("192.0.2.10")
				m.From.Port = 3610
				mutate(m)
				c.receive(m, time.Now())
			}
			c.receive(base, time.Now().Add(-time.Minute))
			c.receive(base, time.Now())
			c.receive(base, time.Now())
			other, _ := protocol.NewMessageWithBytes(base.Bytes())
			other.From.IP = net.ParseIP("192.0.2.20")
			other.From.Port = 3610
			c.receive(other, time.Now())
		}()
		return nil
	}
	if err := s.Discover(context.Background(), ""); err != nil {
		t.Fatal(err)
	}
	if n := len(s.Snapshot().Devices); n != 4 {
		t.Fatalf("expected four instances, got %d", n)
	}
}
func TestNotificationAddsInstancesWithoutReplacingGet(t *testing.T) {
	c := newClient()
	defer c.Close()
	s := NewSession(c)
	inf := Request(SourceEOJ, 0x73, 0xd5, []byte{1, 2, 0x90, 1})
	inf.SetSEOJ(0x0ef001)
	inf.From.IP = net.ParseIP("192.0.2.10")
	inf.From.Port = 3610
	c.receive(inf, time.Now())
	if len(s.Snapshot().Devices) != 1 {
		t.Fatal("D5 missing")
	}
	inf.SetSEOJ(0x013001)
	inf.Property(0).SetCode(0x80)
	inf.Property(0).SetData([]byte{0x31})
	c.receive(inf, time.Now())
	if len(s.Snapshot().Devices) != 2 {
		t.Fatal("device INF missing")
	}
	for _, d := range s.Snapshot().Devices {
		if len(d.Values) != 0 {
			t.Fatal("INF became fresh Get")
		}
	}
}
func TestDiscoveryCancelCloseAndNoLateResponse(t *testing.T) {
	for _, closeClient := range []bool{false, true} {
		c := newClient()
		s := NewSession(c)
		sent := make(chan *protocol.Message, 1)
		c.send = func(_ Target, m *protocol.Message) error { sent <- m; return nil }
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		done := make(chan error, 1)
		go func() { done <- s.Discover(ctx, "") }()
		req := <-sent
		if closeClient {
			c.Close()
		} else {
			cancel()
		}
		select {
		case err := <-done:
			if err == nil {
				t.Fatal("expected cancel/close")
			}
		case <-time.After(time.Second):
			t.Fatal("blocked")
		}
		m := response(req, "192.0.2.10", []byte{1, 2, 0x90, 1})
		m.SetSEOJ(0x0ef001)
		c.receive(m, time.Now())
		if len(s.Snapshot().Devices) != 0 {
			t.Fatal("late discovery accepted")
		}
		cancel()
		c.Close()
	}
}
func TestSchemaRejectsWriteWithoutSending(t *testing.T) {
	c := newClient()
	defer c.Close()
	s := NewSession(c)
	target := Target{"192.0.2.10", 0x029001}
	s.Add(target)
	d := s.devices[target]
	d.Maps = true
	d.Get = []byte{0x80, 0xff, 0xb0}
	d.Set = append([]byte(nil), d.Get...)
	c.send = func(Target, *protocol.Message) error { t.Fatal("unsafe write sent"); return nil }
	for _, x := range []struct {
		ep   byte
		data []byte
	}{{0x80, []byte{0x32}}, {0xb0, []byte{101}}, {0xff, []byte{0}}} {
		if s.Set(context.Background(), target, x.ep, x.data) == nil {
			t.Fatal("schema failure accepted")
		}
	}
}

func TestResponseAfterDeadlineIsUnknown(t *testing.T) {
	c := newClient()
	defer c.Close()
	target := Target{"192.0.2.10", 0x029001}
	c.send = func(_ Target, req *protocol.Message) error {
		go c.receive(response(req, target.IP, []byte{0x30}), time.Now().Add(time.Minute))
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, err := c.RoundTrip(ctx, target, 0x62, 0x80, nil); err != context.DeadlineExceeded {
		t.Fatal("expired response accepted", err)
	}
}
