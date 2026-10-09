package controller

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

func response(req *protocol.Message, ip string, data []byte) *protocol.Message {
	m := protocol.NewResponseMessageWithMessage(req)
	m.From.IP = net.ParseIP(ip)
	m.From.Port = 3610
	p := protocol.NewPropertyWithCode(req.Property(0).Code())
	p.SetData(data)
	m.AddProperty(p)
	return m
}
func TestCorrelationDuplicatesAndParallel(t *testing.T) {
	c := newClient()
	defer c.Close()
	sent := make(chan *protocol.Message, 32)
	c.send = func(_ Target, m *protocol.Message) error { sent <- m; return nil }
	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.RoundTrip(ctx, Target{"127.0.0.1", 0x029001}, 0x62, byte(0x80+i), nil)
			if err != nil || r.Message == nil || r.Message.Property(0).Data()[0] != byte(i) {
				t.Errorf("roundtrip %d: %v", i, err)
			}
		}()
	}
	for range 16 {
		req := <-sent
		valid := response(req, "127.0.0.1", []byte{byte(req.Property(0).Code()) - 0x80})
		invalid := []func(*protocol.Message){func(m *protocol.Message) { m.From.IP = net.ParseIP("127.0.0.2") }, func(m *protocol.Message) { m.From.Port = 4000 }, func(m *protocol.Message) { m.SetSEOJ(0x013001) }, func(m *protocol.Message) { m.SetDEOJ(0x0ef001) }, func(m *protocol.Message) { m.SetESV(0x73) }, func(m *protocol.Message) { m.SetTID(req.TID() + 100) }, func(m *protocol.Message) { m.Property(0).SetCode(0xff) }}
		for _, mutate := range invalid {
			m := response(req, "127.0.0.1", []byte{0xff})
			mutate(m)
			c.receive(m, time.Now())
		}
		c.receive(valid, time.Now())
		c.receive(valid, time.Now())
	}
	wg.Wait()
}
func TestTimeoutCancelCloseAndTIDExhaustion(t *testing.T) {
	for _, mode := range []string{"timeout", "cancel", "close", "exhaust"} {
		t.Run(mode, func(t *testing.T) {
			c := newClient()
			defer c.Close()
			c.send = func(Target, *protocol.Message) error { return nil }
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			if mode == "cancel" {
				cancel()
			}
			if mode == "close" {
				c.Close()
			}
			if mode == "exhaust" {
				c.next = 65535
			}
			_, err := c.RoundTrip(ctx, Target{"127.0.0.1", 0x029001}, 0x62, 0x80, nil)
			if err == nil {
				t.Fatal("expected error")
			}
			if mode == "timeout" && !errors.Is(err, context.DeadlineExceeded) {
				t.Fatal(err)
			}
		})
	}
	c := newClient()
	sent := make(chan struct{})
	c.send = func(Target, *protocol.Message) error { close(sent); return nil }
	done := make(chan error, 1)
	go func() {
		_, err := c.RoundTrip(context.Background(), Target{"127.0.0.1", 0x029001}, 0x62, 0x80, nil)
		done <- err
	}()
	<-sent
	c.Close()
	if err := <-done; !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
}
func TestMalformedDatagrams(t *testing.T) {
	valid := Request(0x029001, 0x62, 0x80, nil).Bytes()
	for _, b := range [][]byte{nil, valid[:11], valid[:13], append(append([]byte{}, valid...), 0), append(append([]byte{}, valid[:13]...), 3)} {
		if _, err := Decode(b); err == nil {
			t.Fatalf("accepted %X", b)
		}
	}
	if _, err := Decode(valid); err != nil {
		t.Fatal(err)
	}
}
func TestDemoDiscoveryMapsSetReadbackAndNotification(t *testing.T) {
	c := Demo()
	defer c.Close()
	s := NewSession(c)
	ctx := context.Background()
	if err := s.Discover(ctx, "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	target := s.Snapshot().Devices[0].Target
	if err := s.Load(ctx, target); err != nil {
		t.Fatal(err)
	}
	if err := s.Set(ctx, target, 0xb0, []byte{50}); err != nil {
		t.Fatal(err)
	}
	v := s.Snapshot().Devices[0].Values[0xb0]
	if !bytes.Equal(v.Data, []byte{50}) || v.State != "Readback success" || v.Sent.IsZero() || v.Received.Before(v.Sent) {
		t.Fatalf("readback %+v", v)
	}
	events := c.Events()
	var setTID, getTID uint
	for _, e := range events {
		if e.Kind == "TX" {
			if strings.Contains(e.Hex, "6101B0") {
				setTID = e.TID
			}
			if strings.Contains(e.Hex, "6201B0") {
				getTID = e.TID
			}
		}
	}
	if setTID == 0 || getTID == setTID {
		t.Fatal("not distinct TIDs")
	}
	// Inject an old INF carrying different data: Get value must remain unchanged.
	m := Request(SourceEOJ, 0x73, 0xb0, []byte{99})
	m.SetSEOJ(target.EOJ)
	m.From.IP = net.ParseIP(target.IP)
	m.From.Port = 3610
	c.receive(m, time.Now().Add(-time.Minute))
	if s.Snapshot().Devices[0].Values[0xb0].Data[0] != 50 {
		t.Fatal("INF replaced Get")
	}
	if err := s.Set(ctx, target, 0xff, []byte{0}); err == nil {
		t.Fatal("unknown property writable")
	}
}
func TestReadbackMismatchAndUnknown(t *testing.T) {
	for _, mode := range []string{"mismatch", "timeout", "rejected"} {
		t.Run(mode, func(t *testing.T) {
			c := newClient()
			defer c.Close()
			s := NewSession(c)
			s.timeout = 20 * time.Millisecond
			target := Target{"127.0.0.1", 0x029001}
			s.Add(target)
			s.devices[target].Maps = true
			s.devices[target].Set = []byte{0x80}
			c.send = func(_ Target, req *protocol.Message) error {
				if mode == "timeout" && req.ESV() == 0x62 {
					return nil
				}
				m := response(req, target.IP, nil)
				if req.ESV() == 0x62 {
					m.Property(0).SetData([]byte{0x31})
				}
				if mode == "rejected" {
					m.SetESV(0x51)
				}
				go c.receive(m, time.Now())
				return nil
			}
			err := s.Set(context.Background(), target, 0x80, []byte{0x30})
			status := s.Snapshot().Status
			if mode == "mismatch" && (!strings.Contains(status, "mismatch") || err != nil) {
				t.Fatal(status, err)
			}
			if mode == "timeout" && !strings.Contains(status, "readback unknown: timeout") {
				t.Fatal(status)
			}
			if mode == "rejected" && !strings.Contains(status, "SetC rejected") {
				t.Fatal(status)
			}
		})
	}
}
func TestParseHex(t *testing.T) {
	for _, s := range []string{"", "0", "zz", strings.Repeat("00", 256)} {
		if _, err := ParseHex(s); err == nil {
			t.Fatal(s)
		}
	}
	v, err := ParseHex("30 4B\n00")
	if err != nil || !bytes.Equal(v, []byte{0x30, 0x4b, 0}) {
		t.Fatal(v, err)
	}
}

func TestBitmapMapsAndFailedRefreshDisableWrites(t *testing.T) {
	c := newClient()
	defer c.Close()
	s := NewSession(c)
	target := Target{"127.0.0.1", 0x029001}
	invalid := false
	c.send = func(_ Target, req *protocol.Message) error {
		ep := byte(req.Property(0).Code())
		data := []byte{0x30}
		if ep >= 0x9d && ep <= 0x9f {
			data = append([]byte{16}, bytes.Repeat([]byte{1}, 16)...)
			if invalid && ep == 0x9e {
				data = []byte{2, 0x80, 0x80}
			}
		}
		go c.receive(response(req, target.IP, data), time.Now())
		return nil
	}
	if err := s.Load(context.Background(), target); err != nil {
		t.Fatal(err)
	}
	d := s.Snapshot().Devices[0]
	if !d.Maps || len(d.Get) != 16 || len(d.Set) != 16 || d.Get[0] != 0x80 || d.Get[15] != 0x8f {
		t.Fatalf("bitmap %+v", d)
	}
	invalid = true
	if err := s.Load(context.Background(), target); err == nil {
		t.Fatal("duplicate map accepted")
	}
	if s.Snapshot().Devices[0].Maps {
		t.Fatal("old write permission survived failed refresh")
	}
	if err := s.Set(context.Background(), target, 0x80, []byte{0x30}); err == nil {
		t.Fatal("stale write permission")
	}
}
func TestCanceledRequestDoesNotAcceptLateReply(t *testing.T) {
	c := newClient()
	defer c.Close()
	sent := make(chan *protocol.Message, 2)
	c.send = func(_ Target, m *protocol.Message) error { sent <- m; return nil }
	target := Target{"127.0.0.1", 0x029001}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { _, err := c.RoundTrip(ctx, target, 0x62, 0x80, nil); done <- err }()
	old := <-sent
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	go func() { _, err := c.RoundTrip(context.Background(), target, 0x62, 0x80, nil); done <- err }()
	current := <-sent
	c.receive(response(old, target.IP, []byte{0x31}), time.Now())
	select {
	case err := <-done:
		t.Fatalf("late response completed new request: %v", err)
	default:
	}
	c.receive(response(current, target.IP, []byte{0x30}), time.Now())
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestDemoConcurrentOverloadRemainsCancellable(t *testing.T) {
	c := Demo()
	defer c.Close()
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.RoundTrip(ctx, Target{"127.0.0.1", 0x029001}, 0x62, 0x80, nil)
			if err != nil {
				if !strings.Contains(err.Error(), "queue full") {
					t.Error(err)
				}
				return
			}
			if r.Message.Property(0).Data()[0] != 0x30 {
				t.Error("wrong overload response")
			}
		}()
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("demo queue blocked receive/cancel")
	}
}
