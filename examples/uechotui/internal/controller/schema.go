package controller

import "github.com/cybergarage/uecho-go/examples/uechotui/internal/mra"

func (d Device) Definitions() map[byte]mra.Property {
	var release byte
	if v := d.Values[0x82]; len(v.Data) == 4 && v.State == "Get success" {
		release = v.Data[2]
	}
	return mra.Properties(uint32(d.Target.EOJ), release)
}
func (d Device) Name() string { return mra.Name(uint32(d.Target.EOJ)) }
