package controller

import (
	"errors"
	"net"
	"reflect"
	"testing"
)

type addressText string

func (a addressText) Network() string { return "fixture" }
func (a addressText) String() string  { return string(a) }

func TestInterfaceOptionsOnlyEligibleLANAddresses(t *testing.T) {
	up := net.FlagUp | net.FlagMulticast
	interfaces := []net.Interface{
		{Name: "loop", Flags: up | net.FlagLoopback},
		{Name: "bad", Flags: up},
		{Name: "down", Flags: net.FlagMulticast},
		{Name: "no-multicast", Flags: net.FlagUp},
		{Name: "failed", Flags: up},
		{Name: "lan0", Flags: up},
		{Name: "lan1", Flags: up},
	}
	addresses := map[string][]string{
		"loop": {"127.0.0.1/8", "192.0.2.99/24"},
		"bad":  {"127.0.0.2/8", "0.0.0.0/0", "224.0.23.0/4", "255.255.255.255/32", "169.254.2.3/16", "::1/128", "2001:db8::1/64", "invalid"},
		"down": {"192.0.2.3/24"}, "no-multicast": {"192.0.2.4/24"},
		"lan0": {"192.0.2.20/24", "192.0.2.21/24"}, "lan1": {"198.51.100.20/24"},
	}
	got := interfaceOptions(interfaces, func(i net.Interface) ([]net.Addr, error) {
		if i.Name == "failed" {
			return nil, errors.New("unavailable")
		}
		var out []net.Addr
		for _, a := range addresses[i.Name] {
			out = append(out, addressText(a))
		}
		return out, nil
	})
	want := []InterfaceOption{{"lan0", "192.0.2.20"}, {"lan0", "192.0.2.21"}, {"lan1", "198.51.100.20"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v; want %+v", got, want)
	}
}
func TestInterfaceOptionsNoLocalhostFallback(t *testing.T) {
	got := interfaceOptions([]net.Interface{{Name: "lo", Flags: net.FlagUp | net.FlagMulticast}}, func(net.Interface) ([]net.Addr, error) { return []net.Addr{addressText("127.0.0.1/8")}, nil })
	if len(got) != 0 {
		t.Fatalf("localhost fallback: %+v", got)
	}
}
