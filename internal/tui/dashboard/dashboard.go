// Package dashboard owns terminal widgets on the tview event loop.
package dashboard

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cybergarage/uecho-go/internal/tui/controller"
	"github.com/cybergarage/uecho-go/internal/tui/mra"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Dashboard struct {
	discoverOnStart      bool
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
	d.props.SetBorder(true).SetTitle(" MRA settings / Enter: Get / w: SetC ")
	d.logs = tview.NewTextView().SetWrap(false).SetScrollable(true)
	d.logs.SetBorder(true).SetTitle(" Protocol events (256 max) / INF arrival only ")
	d.header = tview.NewTextView().SetText(" UECHO CONTROLLER | " + mode + " | MRA 1.3.0 | fresh Get + raw EDT")
	d.status = tview.NewTextView()
	d.search = tview.NewInputField().SetLabel(" / Filter IP / EOJ / class (empty = all): ")
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
	footer := tview.NewTextView().SetText("Tab/Shift-Tab | arrows | Enter Get | w SetC | d / F5 rediscover all\n/ filter | ? help | Esc cancel | q / Ctrl-C exit")
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
		if d.discoverOnStart {
			d.discoverOnStart = false
			d.start(func(ctx context.Context) error { return d.session.Discover(ctx, d.peer) })
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
		if strings.Contains(strings.ToLower(v.Target.String()+" "+v.Name()), query) {
			d.visible = append(d.visible, v)
		}
	}
	row, _ := d.devices.GetSelection()
	d.devices.Clear()
	d.devices.SetCell(0, 0, tview.NewTableCell("IP / EOJ").SetSelectable(false))
	selection := 0
	for i, v := range d.visible {
		d.devices.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%s / %06X %s", v.Target.IP, uint(v.Target.EOJ), v.Name())))
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
	for i, h := range []string{"EPC / MRA name", "GSI", "Value | raw EDT", "Get RX UTC / TID / state"} {
		d.props.SetCell(0, i, tview.NewTableCell(h).SetSelectable(false).SetTextColor(tcell.ColorAqua))
	}
	v, ok := d.device()
	if !ok {
		d.props.SetTitle(" MRA settings / Enter: Get / w: SetC ")
		return
	}
	d.props.SetTitle(fmt.Sprintf(" MRA settings: %s / %06X | seen %s %s ", v.Name(), uint(v.Target.EOJ), v.LastSeen.UTC().Format("15:04:05Z"), v.Origin))
	definitions := v.Definitions()
	codes := mra.Codes(definitions)
	for ep := range v.Values {
		codes = append(codes, ep)
	}
	for _, list := range [][]byte{v.Get, v.Set, v.Announce} {
		codes = append(codes, list...)
	}
	slices.Sort(codes)
	codes = slices.Compact(codes)
	slices.SortStableFunc(codes, func(a, b byte) int {
		supported := func(ep byte) bool {
			return slices.Contains(v.Get, ep) || slices.Contains(v.Set, ep) || slices.Contains(v.Announce, ep)
		}
		if supported(a) && !supported(b) {
			return -1
		}
		if !supported(a) && supported(b) {
			return 1
		}
		return int(a) - int(b)
	})
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
		definition, known := definitions[ep]
		name := definition.Name
		if !known {
			name = "Unknown property"
		}
		value := v.Values[ep]
		state := "unknown / never read"
		if v.Maps && flags == "---" {
			state = "unsupported by device maps"
		}
		if definition.Reason != "" {
			state += " / " + definition.Reason
		}
		if !value.Received.IsZero() {
			state = fmt.Sprintf("%s / %04X / %s", value.Received.UTC().Format("15:04:05.000Z"), value.TID, value.State)
		}
		for col, text := range []string{fmt.Sprintf("%02X %s", ep, name), flags, fmt.Sprintf("%s | %X", definition.Decode(value.Data), value.Data), state} {
			d.props.SetCell(row, col, tview.NewTableCell(text).SetReference(ep).SetMaxWidth([]int{25, 3, 24, 0}[col]))
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
	if e.Key() == tcell.KeyF5 {
		d.rediscover()
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
		d.confirm("Keys: Tab/Shift-Tab focus; arrows select/scroll; Enter load/Get; w typed SetC; d/F5 rediscover all; / filter listed IP/EOJ/class (empty = all); Esc cancel; q/Ctrl-C exit.\n\nMRA state/number/raw schemas provide names and editors. Unsupported types remain raw and read only. Device maps determine availability; Get RX shows snapshot freshness. INF arrival is not a device timestamp. SetC is followed by a separate Get: success, mismatch or unknown. Cancel never rolls back an applied write.", func() {})
		return nil
	case 'w':
		d.write()
		return nil
	case 'd':
		d.rediscover()
		return nil
	}
	return e
}

// F5 is available even while the list filter has focus. d stays a text character
// in the filter. A running request must finish or be canceled before rediscovery.
func (d *Dashboard) rediscover() {
	if d.busy {
		d.session.Status("Busy: Esc cancels current request; then d / F5 rediscover all")
		return
	}
	destination := d.peer
	if destination == "" {
		destination = "224.0.23.0:3610 on the selected interface (3 second collection window)"
	} else {
		destination += ":3610"
	}
	d.confirm("Rediscover all devices with node-profile Get D6 to "+destination+"?", func() {
		d.search.SetText("")
		d.App.SetFocus(d.devices)
		d.start(func(ctx context.Context) error { return d.session.Discover(ctx, d.peer) })
	})
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
	if !v.Maps || !slices.Contains(v.Get, ep) {
		d.session.Status("Get unavailable: not in fresh device Get map")
		return
	}
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
	definition, known := v.Definitions()[ep]
	if !known || !definition.Editable() || !slices.Contains(v.Get, ep) {
		d.session.Status("Read only: unsupported/ambiguous MRA schema or no Get readback map")
		return
	}
	form := tview.NewForm()
	options := definition.Options()
	field, hasInput := definition.InputField()
	labels := []string{}
	for _, o := range options {
		labels = append(labels, o.Label)
	}
	if hasInput {
		labels = append(labels, "Enter "+field.Kind)
	}
	selected := -1
	for i, o := range options {
		if string(o.Data) == string(v.Values[ep].Data) {
			selected = i
		}
	}
	if hasInput && selected < 0 {
		selected = len(options)
	}
	if len(labels) > 0 {
		form.AddDropDown("Value kind", labels, selected, func(_ string, i int) { selected = i })
	}
	input := tview.NewInputField()
	if hasInput {
		label := "Raw hex"
		current := fmt.Sprintf("%X", v.Values[ep].Data)
		if field.Kind == "number" {
			label = fmt.Sprintf("%.6g..%.6g %s (step %.6g)", field.Min*field.Scale, field.Max*field.Scale, field.Unit, field.Scale)
			decoded := definition.Decode(v.Values[ep].Data)
			current = strings.Fields(decoded)[0]
			if current == "raw" || current == "unread" {
				current = ""
			}
		}
		input.SetLabel(label + ": ").SetText(current)
		form.AddFormItem(input)
	}
	form.AddButton("Cancel", d.closeModal).AddButton("Review", func() {
		var data []byte
		var err error
		if selected >= 0 && selected < len(options) {
			data = append([]byte(nil), options[selected].Data...)
		} else if hasInput {
			data, err = field.Encode(input.GetText())
		} else {
			err = fmt.Errorf("select a value")
		}
		if err == nil {
			err = definition.Validate(data)
		}
		if err != nil {
			form.SetTitle(err.Error())
			return
		}
		d.closeModal()
		d.confirm(fmt.Sprintf("Send SetC EPC %02X %s = %s\nEDT %X (%d bytes)?\n%s\nThen fresh Get readback. Effects may persist after cancel/timeout.", ep, definition.Name, definition.Decode(data), data, len(data), v.Target), func() { d.start(func(ctx context.Context) error { return d.session.Set(ctx, v.Target, ep, data) }) })
	})
	form.SetBorder(true).SetTitle(fmt.Sprintf("SetC %02X / %s / MRA 1.3.0", ep, definition.Name))
	form.SetFocus(form.GetFormItemCount())
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
// DiscoverOnStart schedules one discovery after the first screen draw.
// Call before Run; the empty filter shows every discovered device.
func (d *Dashboard) DiscoverOnStart() { d.discoverOnStart = true }

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
