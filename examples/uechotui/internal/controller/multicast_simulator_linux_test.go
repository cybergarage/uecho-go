package controller

import (
	"context"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
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
	// Supply the released terminal mode a private Linux PTY.
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	if err = unix.IoctlSetPointerInt(int(master.Fd()), unix.TIOCSPTLCK, 0); err != nil {
		t.Fatal(err)
	}
	number, err := unix.IoctlGetInt(int(master.Fd()), unix.TIOCGPTN)
	if err != nil {
		t.Fatal(err)
	}
	slave, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|syscall.O_NOCTTY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer slave.Close()
	cmd := exec.Command(binary, "--display", "127.0.0.1:18990", "--udp", "192.0.2.10:3610", "--allow-lan", "--multicast-interface", "simtest0")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	cmd.Stdin = slave
	cmd.Stdout = slave
	cmd.Stderr = slave
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true, Ctty: 0}
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	slave.Close()
	drained := make(chan struct{})
	go func() { io.Copy(io.Discard, master); close(drained) }()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		cmd.Process.Signal(os.Interrupt)
		select {
		case err := <-done:
			if err != nil {
				t.Error("simulator exit", err)
			}
		case <-time.After(2 * time.Second):
			cmd.Process.Kill()
			<-done
			t.Error("simulator shutdown timed out")
		}
		master.Close()
		<-drained
	}()
	httpClient := http.Client{Timeout: 200 * time.Millisecond}
	ready := false
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		res, e := httpClient.Get("http://127.0.0.1:18990/")
		if e == nil {
			res.Body.Close()
			ready = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !ready {
		t.Fatal("released simulator failed to start")
	}
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
