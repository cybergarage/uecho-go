// Package testutil supplies the opt-in isolated integration-test harness.
package testutil

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

// StartSimulator provides a private PTY and waits for its read-only display.
// Callers must use the disposable network namespace supplied by CI.
func StartSimulator(t *testing.T, binary, url string, args ...string) {
	t.Helper()
	// Supply the released terminal mode a private Linux PTY.
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { master.Close() })
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
	t.Cleanup(func() { slave.Close() })
	cmd := exec.Command(binary, args...)
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
	t.Cleanup(func() {
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
	})
	httpClient := http.Client{Timeout: 200 * time.Millisecond}
	ready := false
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		res, e := httpClient.Get(url)
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
}
