package controller

import (
	"context"
	"github.com/cybergarage/uecho-go/internal/tui/testutil"
	"net"
	"os"
	"runtime"
	"testing"
	"time"
)

// Never runs against a host interface: the external test harness creates a
// namespace containing only loopback + simtest0 with documentation-only IPs.
func TestSimulatorMulticast(t *testing.T) {
	if os.Getenv("UECHOTUI_ISOLATED_MULTICAST") != "1" {
		t.Skip("isolated namespace opt-in")
	}
	if runtime.GOOS != "linux" {
		t.Skip("Linux isolated namespace")
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range ifaces {
		if i.Name != "lo" && i.Name != "simtest0" {
			addresses, e := i.Addrs()
			if e != nil || i.Flags&net.FlagUp != 0 || len(addresses) > 0 {
				t.Fatalf("refuse non-isolated interface %s", i.Name)
			}
		}
	}
	binary := os.Getenv("UECHOTUI_SIMULATOR_V1")
	if binary == "" {
		t.Fatal("released v1.0.0 binary required")
	}
	testutil.StartSimulator(t, binary, "http://127.0.0.1:18990/", "--display", "127.0.0.1:18990", "--udp", "192.0.2.10:3610", "--allow-lan", "--multicast-interface", "simtest0")
	c, err := OpenMulticast("simtest0", "192.0.2.20")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	s := NewSession(c)
	s.timeout = 300 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = s.Discover(ctx, ""); err != nil {
		t.Fatal(err)
	}
	if n := len(s.Snapshot().Devices); n != 3 {
		t.Fatalf("multicast discovery instances=%d", n)
	}
	for _, d := range s.Snapshot().Devices {
		if err = s.Load(ctx, d.Target); err != nil {
			t.Fatal(err)
		}
	}
	target := Target{"192.0.2.10", 0x029001}
	if err = s.Set(ctx, target, 0xb0, []byte{75}); err != nil {
		t.Fatal(err)
	}
	for _, d := range s.Snapshot().Devices {
		if d.Target == target {
			v := d.Values[0xb0]
			if v.State != "Readback success" || d.Definitions()[0xb0].Decode(v.Data) != "75 %" {
				t.Fatalf("typed readback %+v", v)
			}
			t.Logf("v1.0.0 multicast discovery=3; lighting setting=75%% raw=%X Get TID=%04X TX=%s RX=%s", v.Data, v.TID, v.Sent.Format(time.RFC3339Nano), v.Received.Format(time.RFC3339Nano))
		}
	}
	found := false
	for _, e := range c.Events() {
		if e.Kind == "INF" && e.From == "192.0.2.10:3610" {
			found = true
		}
	}
	if !found {
		t.Fatal("multicast INF missing")
	}
}
