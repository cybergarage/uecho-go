package controller

import (
	"net"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// Demo opens no sockets and uses a tiny raw-hex fixture, not an MRA schema.
func Demo() *Client {
	c := newClient()
	values := map[byte][]byte{0x80: {0x30}, 0xb0: {0x4b}, 0x9d: {1, 0x80}, 0x9e: {2, 0x80, 0xb0}, 0x9f: {5, 0x80, 0xb0, 0x9d, 0x9e, 0x9f}}
	packets := make(chan *protocol.Message, 16)
	c.send = func(_ Target, m *protocol.Message) error {
		clone, _ := protocol.NewMessageWithBytes(m.Bytes())
		packets <- clone
		return nil
	}
	c.workers.Add(1)
	go func() {
		defer c.workers.Done()
		for {
			select {
			case <-c.stop:
				return
			case req := <-packets:
				res := protocol.NewResponseMessageWithMessage(req)
				res.From.IP = net.ParseIP("127.0.0.1")
				res.From.Port = 3610
				ep := byte(req.Property(0).Code())
				p := protocol.NewPropertyWithCode(protocol.PropertyCode(ep))
				if req.DEOJ() == 0x0ef001 && ep == 0xd6 {
					p.SetData([]byte{1, 2, 0x90, 1})
				} else if req.ESV() == 0x61 {
					values[ep] = append([]byte(nil), req.Property(0).Data()...)
				} else {
					p.SetData(append([]byte(nil), values[ep]...))
				}
				res.AddProperty(p)
				c.receive(res, time.Now())
				if req.ESV() == 0x61 {
					inf := protocol.NewMessage()
					inf.SetSEOJ(req.DEOJ())
					inf.SetDEOJ(SourceEOJ)
					inf.SetESV(0x73)
					inf.SetTID(req.TID())
					p := protocol.NewPropertyWithCode(protocol.PropertyCode(ep))
					p.SetData(append([]byte(nil), values[ep]...))
					inf.AddProperty(p)
					inf.From.IP = net.ParseIP("127.0.0.1")
					inf.From.Port = 3610
					c.receive(inf, time.Now())
				}
			}
		}
	}()
	c.Log(Event{Kind: "DEMO", Text: "offline fixture; no sockets, raw hex only"})
	return c
}
