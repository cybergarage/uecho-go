// Package dashboard owns terminal widgets on the tview event loop.
package dashboard

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cybergarage/uecho-go/examples/uechotui/internal/controller"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Dashboard struct {
	screenReady          chan tcell.Screen
	App                  *tview.Application
	session              *controller.Session
	peer                 string
	ctx                  context.Context
	cancel               context.CancelFunc
	jobs                 sync.WaitGroup
	jobCancel            context.CancelFunc
	busy                 bool
	finished             chan struct{}
	pages                *tview.Pages
	root, body           *tview.Flex
	devices, props       *tview.Table
	logs, header, status *tview.TextView
	search               *tview.InputField
	modal                tview.Primitive
	prior                tview.Primitive
	visible              []controller.Device
	selected             controller.Target
	ep                   byte
	last                 string
}

func New(s *controller.Session, peer, mode string) *Dashboard {
	d := &Dashboard{App: tview.NewApplication(), session: s, peer: peer, ep: 0x80, finished: make(chan struct{}, 1), screenReady: make(chan tcell.Screen, 1)}
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.devices = tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	d.devices.SetBorder(true).SetTitle(" Devices / Enter: load maps + Get ")
	d.props = tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	d.props.SetBorder(true).SetTitle(" Raw properties / Enter: Get / w: SetC ")
	d.logs = tview.NewTextView().SetWrap(false).SetScrollable(true)
	d.logs.SetBorder(true).SetTitle(" Protocol events (256 max) / INF arrival only ")
	d.header = tview.NewTextView().SetText(" UECHO CONTROLLER | " + mode + " | raw hex — schema not inferred")
	d.status = tview.NewTextView()
	d.search = tview.NewInputField().SetLabel(" / Search IP / EOJ: ")
	d.search.SetChangedFunc(func(string) { d.refresh(true) })
	d.search.SetDoneFunc(func(tcell.Key) { d.App.SetFocus(d.devices) })
	d.devices.SetSelectionChangedFunc(func(row, _ int) {
		if row > 0 && row <= len(d.visible) {
			d.selected = d.visible[row-1].Target
			d.refreshProperties()
		}
	})
	d.devices.SetSelectedFunc(func(int, int) {
		if d.selected.EOJ != 0 {
			target := d.selected
			d.confirm("Send property map Gets (9D/9E/9F), then fresh Get for each readable EPC?\n"+target.String(), func() { d.start(func(ctx context.Context) error { return s.Load(ctx, target) }) })
		}
	})
	d.props.SetSelectionChangedFunc(func(row, _ int) {
		cell := d.props.GetCell(row, 0)
		if ep, ok := cell.GetReference().(byte); ok {
			d.ep = ep
		}
	})
	d.props.SetSelectedFunc(func(int, int) { d.get() })
	d.body = tview.NewFlex()
	footer := tview.NewTextView().SetText("Tab/Shift-Tab | arrows | Enter Get | w SetC | d discover\n/ search | ? help | Esc cancel | q / Ctrl-C exit")
	d.root = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(d.header, 1, 0, false).AddItem(d.search, 1, 0, false).AddItem(d.body, 0, 1, true).AddItem(d.logs, 9, 0, false).AddItem(footer, 2, 0, false).AddItem(d.status, 2, 0, false)
	d.pages = tview.NewPages().AddPage("main", d.root, true, true)
	d.App.SetRoot(d.pages, true).EnableMouse(false).EnablePaste(true).SetFocus(d.devices).SetInputCapture(d.capture)
	for _, b := range []*tview.Box{d.devices.Box, d.props.Box, d.logs.Box} {
		b.SetFocusFunc(func() { b.SetBorderColor(tcell.ColorAqua) })
		b.SetBlurFunc(func() { b.SetBorderColor(tcell.ColorGray) })
	}
	d.App.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		select {
		case d.screenReady <- screen:
		default:
		}
		w, h := screen.Size()
		d.body.Clear()
		if w < 100 || h < 28 {
			d.root.ResizeItem(d.logs, 6, 0)
			d.body.SetDirection(tview.FlexRow).AddItem(d.devices, 6, 0, false).AddItem(d.props, 0, 1, false)
		} else {
			d.root.ResizeItem(d.logs, 9, 0)
			d.body.SetDirection(tview.FlexColumn).AddItem(d.devices, 38, 0, false).AddItem(d.props, 0, 1, false)
		}
		d.refresh(false)
		return false
	})
	d.devices.SetBorderColor(tcell.ColorAqua)
	d.refresh(true)
	return d
}
func (d *Dashboard) refresh(force bool) {
	select {
	case <-d.finished:
		d.busy = false
		d.jobCancel = nil
	default:
	}
	snap := d.session.Snapshot()
	events := d.session.Client.Events()
	key := fmt.Sprintf("%v/%v/%s", snap, events, d.search.GetText())
	if !force && key == d.last {
		return
	}
	d.last = key
	d.status.SetText(snap.Status)
	d.visible = nil
	query := strings.ToLower(d.search.GetText())
	for _, v := range snap.Devices {
		if strings.Contains(strings.ToLower(v.Target.String()), query) {
			d.visible = append(d.visible, v)
		}
	}
	row, _ := d.devices.GetSelection()
	d.devices.Clear()
	d.devices.SetCell(0, 0, tview.NewTableCell("IP / EOJ").SetSelectable(false))
	selection := 0
	for i, v := range d.visible {
		d.devices.SetCell(i+1, 0, tview.NewTableCell(v.Target.String()))
		if v.Target == d.selected {
			selection = i + 1
		}
	}
	if selection == 0 && len(d.visible) > 0 {
		selection = min(max(row, 1), len(d.visible))
		d.selected = d.visible[selection-1].Target
	}
	if len(d.visible) == 0 {
		d.selected = controller.Target{}
	}
	d.devices.Select(selection, 0)
	d.refreshProperties()
	var b strings.Builder
	for _, e := range events {
		fmt.Fprintf(&b, "%s %-7s TID %04X %s -> %s %s\n", e.Time.UTC().Format("15:04:05.000Z"), e.Kind, e.TID, e.From, e.To, e.Text)
		if e.Hex != "" {
			fmt.Fprintf(&b, "  %s\n", e.Hex)
		}
	}
	d.logs.SetText(b.String())
	if !d.logs.HasFocus() {
		d.logs.ScrollToEnd()
	}
}
func (d *Dashboard) device() (controller.Device, bool) {
	for _, v := range d.session.Snapshot().Devices {
		if v.Target == d.selected {
			return v, true
		}
	}
	return controller.Device{}, false
}
func (d *Dashboard) refreshProperties() {
	d.props.Clear()
	for i, h := range []string{"EPC", "Get/Set/INF", "EDT (raw hex)", "Get RX UTC / TID / state"} {
		d.props.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetTextColor(tcell.ColorAqua))
	}
	v, ok := d.device()
	if !ok {
		return
	}
	codes := []byte{}
	for ep := range v.Values {
		codes = append(codes, ep)
	}
	for _, list := range [][]byte{v.Get, v.Set, v.Announce} {
		codes = append(codes, list...)
	}
	slices.Sort(codes)
	codes = slices.Compact(codes)
	selection := 1
	for i, ep := range codes {
		row := i + 1
		if ep == d.ep {
			selection = row
		}
		flags := "---"
		if v.Maps {
			f := []byte(flags)
			if slices.Contains(v.Get, ep) {
				f[0] = 'G'
			}
			if slices.Contains(v.Set, ep) {
				f[1] = 'S'
			}
			if slices.Contains(v.Announce, ep) {
				f[2] = 'I'
			}
			flags = string(f)
		}
		value := v.Values[ep]
		state := "unknown / never read"
		if !value.Received.IsZero() {
			state = fmt.Sprintf("%s / %04X / %s", value.Received.UTC().Format("15:04:05.000Z"), value.TID, value.State)
		}
		for col, text := range []string{fmt.Sprintf("%02X", ep), flags, fmt.Sprintf("%X", value.Data), state} {
			d.props.SetCell(row, col, tview.NewTableCell(text).SetReference(ep))
		}
	}
	if len(codes) > 0 {
		d.props.Select(selection, 0)
	}
}
func (d *Dashboard) capture(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyF24 {
		return nil
	}
	if e.Key() == tcell.KeyCtrlC {
		d.cancel()
		d.App.Stop()
		return nil
	}
	if d.modal != nil {
		if e.Key() == tcell.KeyEscape {
			d.closeModal()
			return nil
		}
		return e
	}
	if e.Key() == tcell.KeyEscape {
		if d.jobCancel != nil {
			d.jobCancel()
			d.session.Status("Canceled: in-flight SetC may have applied; use fresh Get")
		}
		d.search.SetText("")
		d.App.SetFocus(d.devices)
		return nil
	}
	if d.App.GetFocus() == d.search {
		return e
	}
	if e.Key() == tcell.KeyTab || e.Key() == tcell.KeyBacktab {
		list := []tview.Primitive{d.devices, d.props, d.logs}
		index := 0
		for i, p := range list {
			if p == d.App.GetFocus() {
				index = i
			}
		}
		step := 1
		if e.Key() == tcell.KeyBacktab {
			step = 2
		}
		d.App.SetFocus(list[(index+step)%3])
		return nil
	}
	switch e.Rune() {
	case 'q':
		d.cancel()
		d.App.Stop()
		return nil
	case '/':
		d.App.SetFocus(d.search)
		return nil
	case '?':
		d.confirm("Keys: Tab/Shift-Tab focus; arrows select/scroll; Enter load/Get; w raw SetC; d discovery; / search; Esc cancel; q/Ctrl-C exit.\n\nAll EDT is raw hex. INF arrival is not a device timestamp. SetC is followed by a separate Get: success, mismatch or unknown. Cancel never rolls back an applied write.", func() {})
		return nil
	case 'w':
		d.write()
		return nil
	case 'd':
		d.confirm("Send node-profile Get D6 discovery to "+d.peer+":3610?", func() { d.start(func(ctx context.Context) error { return d.session.Discover(ctx, d.peer) }) })
		return nil
	}
	return e
}
func (d *Dashboard) closeModal() {
	d.pages.RemovePage("dialog")
	d.modal = nil
	d.App.SetFocus(d.prior)
}
func (d *Dashboard) show(p tview.Primitive) {
	d.prior = d.App.GetFocus()
	d.modal = p
	centered := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(nil, 0, 1, false).AddItem(tview.NewFlex().AddItem(nil, 0, 1, false).AddItem(p, 0, 8, true).AddItem(nil, 0, 1, false), 0, 6, true).AddItem(nil, 0, 1, false)
	d.pages.AddPage("dialog", centered, true, true)
	d.App.SetFocus(p)
}
func (d *Dashboard) confirm(text string, apply func()) {
	m := tview.NewModal().SetText(text).AddButtons([]string{"Cancel", "Confirm"}).SetDoneFunc(func(i int, _ string) {
		d.closeModal()
		if i == 1 {
			apply()
		}
	})
	d.show(m)
}
func (d *Dashboard) get() {
	v, ok := d.device()
	if !ok {
		return
	}
	ep := d.ep
	d.confirm(fmt.Sprintf("Send fresh Get EPC %02X?\n%s", ep, v.Target), func() { d.start(func(ctx context.Context) error { return d.session.Get(ctx, v.Target, ep) }) })
}
func (d *Dashboard) write() {
	v, ok := d.device()
	if !ok {
		return
	}
	ep := d.ep
	if !v.Maps || !slices.Contains(v.Set, ep) {
		d.session.Status("Read only / unknown permission: load property maps first")
		return
	}
	form := tview.NewForm()
	input := tview.NewInputField().SetLabel("EDT raw hex: ").SetText(fmt.Sprintf("%X", v.Values[ep].Data))
	form.AddFormItem(input).AddButton("Cancel", d.closeModal).AddButton("Review", func() {
		data, err := controller.ParseHex(input.GetText())
		if err != nil {
			form.SetTitle(err.Error())
			return
		}
		d.closeModal()
		d.confirm(fmt.Sprintf("Send SetC EPC %02X EDT %X (%d bytes)?\n%s\nThen send fresh Get to verify. Device effect may persist after cancel/timeout.", ep, data, len(data), v.Target), func() { d.start(func(ctx context.Context) error { return d.session.Set(ctx, v.Target, ep, data) }) })
	})
	form.SetBorder(true).SetTitle(fmt.Sprintf("SetC %02X — schema unknown; raw hex", ep))
	form.SetFocus(1)
	d.show(form)
}
func (d *Dashboard) start(f func(context.Context) error) {
	if d.busy {
		d.session.Status("Busy: Esc cancels current request")
		return
	}
	d.busy = true
	ctx, cancel := context.WithCancel(d.ctx)
	d.jobCancel = cancel
	d.jobs.Add(1)
	go func() { defer d.jobs.Done(); defer cancel(); _ = f(ctx); d.finished <- struct{}{} }()
}

// Run polls detached snapshots; workers never mutate widgets or block waiting on
// QueueUpdateDraw after Stop. Fini remains owned by tview.
func (d *Dashboard) Run(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	// Let tview create/init native screens: SetScreen discards Init errors.
	done := make(chan struct{})
	go func() {
		var screen tcell.Screen
		select {
		case screen = <-d.screenReady:
		case <-done:
			return
		}
		lastState := ""
		tick := time.NewTicker(100 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				d.cancel()
				d.App.Stop()
				return
			case <-done:
				return
			case <-tick.C:
				state := fmt.Sprintf("%v/%v", d.session.Snapshot(), d.session.Client.Events())
				if state != lastState {
					lastState = state
					screen.PostEvent(tcell.NewEventKey(tcell.KeyF24, 0, tcell.ModNone))
				}
			}
		}
	}()
	err := d.App.Run()
	close(done)
	d.cancel()
	d.jobs.Wait()
	return err
}
