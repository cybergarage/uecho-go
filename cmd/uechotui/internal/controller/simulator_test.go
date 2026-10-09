package controller

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

// Deliberately opt-in: default tests never bind sockets. Linux's 127/8 loopback
// permits two distinct addresses on standard port 3610 without aliases or LAN.
func TestSimulatorLoopback(t *testing.T) {
	binary := os.Getenv("UECHOTUI_SIMULATOR_V1")
	if binary == "" {
		t.Skip("set UECHOTUI_SIMULATOR_V1 to a v1.0.0 binary; loopback sockets only")
	}
	if runtime.GOOS != "linux" {
		t.Skip("standard-port two-address test requires Linux 127/8; no host aliases added")
	}
	cmd := exec.Command(binary, "--plain", "--udp", "127.0.0.2:3610")
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { in.Close(); cmd.Process.Kill() }()
	ready := make(chan struct{})
	go func() {
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "Explicit loopback UDP:") {
				close(ready)
				return
			}
		}
	}()
	select {
	case <-ready:
	case <-time.After(5 * time.Second):
		t.Fatal("simulator failed to bind isolated loopback")
	}
	c, err := Open("lo", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	s := NewSession(c)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = s.Discover(ctx, "127.0.0.2"); err != nil {
		t.Fatal(err)
	}
	devices := s.Snapshot().Devices
	if len(devices) != 3 {
		t.Fatalf("discovery %+v", devices)
	}
	for _, d := range devices {
		if err = s.Load(ctx, d.Target); err != nil {
			t.Fatal(err)
		}
	}
	target := Target{"127.0.0.2", 0x029001}
	if err = s.Set(ctx, target, 0x80, []byte{0x30}); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	hasINF := false
	for time.Now().Before(deadline) {
		for _, e := range c.Events() {
			if e.Kind == "INF" && e.From == "127.0.0.2:3610" {
				hasINF = true
			}
		}
		if hasINF {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !hasINF {
		t.Fatal("v1 status-change INF not received")
	}
	for _, d := range s.Snapshot().Devices {
		if d.Target == target {
			v := d.Values[0x80]
			if len(v.Data) != 1 || v.Data[0] != 0x30 || v.State != "Readback success" {
				t.Fatalf("readback %+v", v)
			}
			t.Logf("v1.0.0 discovery=3 devices; maps/Get all; SetC 80=30; readback TID=%04X TX=%s RX=%s; INF received", v.TID, v.Sent.UTC().Format(time.RFC3339Nano), v.Received.UTC().Format(time.RFC3339Nano))
		}
	}
	fmt.Fprintln(in, "quit")
	in.Close()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err = <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("simulator shutdown blocked")
	}
}
