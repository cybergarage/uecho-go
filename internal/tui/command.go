// Package tui provides the shared terminal controller for both CLI entry points.
package tui

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/uecho-go/internal/tui/controller"
	"github.com/cybergarage/uecho-go/internal/tui/dashboard"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type options struct {
	demo, network     bool
	iface, bind, peer string
}

// NewCommand returns a fresh command; construction and help never open sockets.
func NewCommand() *cobra.Command {
	o := options{network: true}
	c := &cobra.Command{
		Use: "tui", Short: "Discover all ECHONET Lite devices and open the terminal controller",
		Long: "Open the network controller and discover all devices on the selected interface. Use --demo for a socket-free fixture. / filters the listed devices; d/F5 repeats discovery. Select a device for fresh maps/Get; property Enter opens its MRA editor. g reads a property; r refreshes device maps/Get. Writes always require confirmation; write-only outcomes remain unverified.",
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if o.demo && c.Flags().Changed("network") && o.network {
				return fmt.Errorf("--demo and --network cannot be combined")
			}
			return run(c.Context(), o)
		},
	}
	c.Flags().BoolVar(&o.demo, "demo", false, "use the offline demo; open no sockets")
	c.Flags().BoolVar(&o.demo, "offline", false, "alias for --demo")
	c.Flags().BoolVar(&o.network, "network", true, "network mode (default); --network=false selects demo")
	c.Flags().StringVar(&o.iface, "interface", "", "local interface; specify together with --bind")
	c.Flags().StringVar(&o.bind, "bind", "", "local IPv4; specify together with --interface")
	c.Flags().StringVar(&o.peer, "peer", "", "optional literal IPv4 for unicast discovery instead of all-device multicast")
	return c
}

func (o options) validate() error {
	if o.demo || !o.network {
		if o.iface != "" || o.bind != "" || o.peer != "" {
			return fmt.Errorf("network options cannot be used with demo/offline mode")
		}
		return nil
	}
	if (o.iface == "") != (o.bind == "") {
		return fmt.Errorf("specify both --interface and --bind")
	}
	if o.peer != "" {
		ip := net.ParseIP(o.peer)
		if ip == nil || ip.To4() == nil || ip.IsUnspecified() || ip.IsMulticast() || ip.Equal(net.IPv4bcast) {
			return fmt.Errorf("--peer must be literal unicast IPv4")
		}
	}
	return nil
}

func run(ctx context.Context, o options) (result error) {
	if err := o.validate(); err != nil {
		return err
	}
	state, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("interactive terminal required: %w", err)
	}
	defer func() { result = errors.Join(result, term.Restore(int(os.Stdin.Fd()), state)) }()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runSession(ctx, o, dependencies{
		selectInterface: dashboard.SelectInterface,
		open:            controller.Open, openMulticast: controller.OpenMulticast,
		runDashboard: func(ctx context.Context, d *dashboard.Dashboard) error { return d.Run(ctx) },
	})
}

type dependencies struct {
	selectInterface     func(context.Context) (controller.InterfaceOption, error)
	open, openMulticast func(string, string) (*controller.Client, error)
	runDashboard        func(context.Context, *dashboard.Dashboard) error
}

func runSession(ctx context.Context, o options, deps dependencies) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := o.validate(); err != nil {
		return err
	}
	var c *controller.Client
	var err error
	mode, target := "OFFLINE DEMO — no sockets", "127.0.0.1"
	if o.demo || !o.network {
		c = controller.Demo()
	} else {
		if o.iface == "" {
			choice, e := deps.selectInterface(ctx)
			if e != nil {
				return e
			}
			o.iface, o.bind = choice.Name, choice.IP
		}
		target = o.peer
		if target == "" {
			c, err = deps.openMulticast(o.iface, o.bind)
		} else {
			c, err = deps.open(o.iface, o.bind)
		}
		if err != nil {
			return err
		}
		mode = "UDP " + o.iface + " / " + o.bind + ":3610 | all-device multicast discovery"
		if target != "" {
			mode = "UDP " + o.iface + " / " + o.bind + ":3610 -> " + target + ":3610 (unicast)"
		}
	}
	defer c.Close()
	d := dashboard.New(controller.NewSession(c), target, mode)
	d.LoadOnSelection()
	d.DiscoverOnStart()
	return deps.runDashboard(ctx, d)
}
