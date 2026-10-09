package tui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/cybergarage/uecho-go/internal/tui/controller"
	"github.com/cybergarage/uecho-go/internal/tui/dashboard"
)

func TestHelpAndNetworkDefault(t *testing.T) {
	c := NewCommand()
	network, _ := c.Flags().GetBool("network")
	demo, _ := c.Flags().GetBool("demo")
	if !network || demo {
		t.Fatal("network must be the default")
	}
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetArgs([]string{"--help"})
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"--demo", "--interface", "--bind", "/ filters", "d/F5 repeats"} {
		if !strings.Contains(out.String(), text) {
			t.Fatalf("missing help %q", text)
		}
	}
}

func TestSessionModesUseFakeClients(t *testing.T) {
	for _, tc := range []struct {
		name                         string
		opts                         options
		selected, multicast, unicast int
	}{
		{"default all", options{network: true}, 1, 1, 0},
		{"explicit all", options{network: true, iface: "fixture0", bind: "192.0.2.20"}, 0, 1, 0},
		{"peer", options{network: true, iface: "fixture0", bind: "192.0.2.20", peer: "192.0.2.10"}, 0, 0, 1},
		{"demo", options{network: true, demo: true}, 0, 0, 0},
		{"legacy offline", options{}, 0, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			selected, multicast, unicast, displayed := 0, 0, 0, 0
			deps := dependencies{
				selectInterface: func(context.Context) (controller.InterfaceOption, error) {
					selected++
					return controller.InterfaceOption{Name: "fixture0", IP: "192.0.2.20"}, nil
				},
				openMulticast: func(iface, bind string) (*controller.Client, error) {
					multicast++
					if iface != "fixture0" || bind != "192.0.2.20" {
						t.Fatal("wrong local address")
					}
					return controller.Demo(), nil
				},
				open:         func(string, string) (*controller.Client, error) { unicast++; return controller.Demo(), nil },
				runDashboard: func(context.Context, *dashboard.Dashboard) error { displayed++; return nil },
			}
			if err := runSession(context.Background(), tc.opts, deps); err != nil {
				t.Fatal(err)
			}
			if selected != tc.selected || multicast != tc.multicast || unicast != tc.unicast || displayed != 1 {
				t.Fatalf("selection/multicast/unicast/dashboard %d/%d/%d/%d", selected, multicast, unicast, displayed)
			}
		})
	}
}

func TestInvalidOptionsAndCanceledSelectionDoNotOpen(t *testing.T) {
	for _, opts := range []options{
		{network: true, iface: "fixture0"}, {network: true, bind: "192.0.2.20"},
		{network: true, peer: "224.0.23.0"}, {network: true, peer: "255.255.255.255"},
		{demo: true, network: true, peer: "192.0.2.10"},
	} {
		if err := runSession(context.Background(), opts, dependencies{}); err == nil {
			t.Fatal("invalid options accepted")
		}
	}
	err := runSession(context.Background(), options{network: true}, dependencies{
		selectInterface: func(context.Context) (controller.InterfaceOption, error) {
			return controller.InterfaceOption{}, context.Canceled
		},
	})
	if err != context.Canceled {
		t.Fatalf("selection cancellation: %v", err)
	}
}

func TestExplicitLoopbackBypassesNormalCandidateSelection(t *testing.T) {
	opened, displayed := false, false
	err := runSession(context.Background(), options{network: true, iface: "lo", bind: "127.0.0.1", peer: "127.0.0.2"}, dependencies{
		selectInterface: func(context.Context) (controller.InterfaceOption, error) {
			t.Fatal("explicit loopback entered LAN picker")
			return controller.InterfaceOption{}, nil
		},
		openMulticast: func(string, string) (*controller.Client, error) {
			t.Fatal("unicast loopback joined multicast")
			return nil, nil
		},
		open: func(iface, bind string) (*controller.Client, error) {
			if iface != "lo" || bind != "127.0.0.1" {
				t.Fatal("explicit binding changed")
			}
			opened = true
			return controller.Demo(), nil
		},
		runDashboard: func(context.Context, *dashboard.Dashboard) error { displayed = true; return nil },
	})
	if err != nil || !opened || !displayed {
		t.Fatalf("explicit loopback: open=%v displayed=%v err=%v", opened, displayed, err)
	}
}
