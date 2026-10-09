// uechotui is an explicitly opt-in developer controller. Default is offline.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/term"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/uecho-go/examples/uechotui/internal/controller"
	"github.com/cybergarage/uecho-go/examples/uechotui/internal/dashboard"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() (result error) {
	network := flag.Bool("network", false, "explicitly enable UDP; requires --interface, --bind, --peer")
	iface := flag.String("interface", "", "confirmed local interface name")
	bind := flag.String("bind", "", "confirmed local IPv4; listen port is 3610")
	peer := flag.String("peer", "", "confirmed device/simulator literal IPv4; discovery is unicast")
	flag.Parse()
	state, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("interactive terminal required: %w", err)
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), state); err != nil {
			result = errors.Join(result, fmt.Errorf("restore terminal: %w", err))
		}
	}()
	var c *controller.Client
	mode := "OFFLINE DEMO — no sockets"
	target := "127.0.0.1"
	if *network {
		if *iface == "" || *bind == "" || *peer == "" {
			return fmt.Errorf("--network requires explicit --interface, --bind and --peer")
		}
		ip := net.ParseIP(*peer)
		if ip == nil || ip.To4() == nil || ip.IsUnspecified() || ip.IsMulticast() {
			return fmt.Errorf("--peer must be literal unicast IPv4")
		}
		c, err = controller.Open(*iface, *bind)
		if err != nil {
			return err
		}
		target = ip.String()
		mode = "UDP " + *bind + ":3610 -> " + target + ":3610"
	} else {
		if *iface != "" || *bind != "" || *peer != "" {
			return fmt.Errorf("network options require --network")
		}
		c = controller.Demo()
	}
	defer c.Close()
	s := controller.NewSession(c)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return dashboard.New(s, target, mode).Run(ctx)
}
