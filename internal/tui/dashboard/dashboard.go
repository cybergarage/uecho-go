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
	submit               func()
	prior                tview.Primitive
	visible              []controller.Device
	selected             controller.Target
	ep                   byte
	last                 string
	deviceNameWidth      int
	autoLoad             bool
	pendingLoad          controller.Target
	readJob              bool
}

func New(s *controller.Session, peer, mode string) *Dashboard {
	d := &Dashboard{App: tview.NewApplication(), session: s, peer: peer, ep: 0x80, finished: make(chan struct{}, 1), screenReady: make(chan tcell.Screen, 1)}
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.devices = tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	d.devices.SetBorder(true).SetTitle(" Devices / select: Get / Enter: refresh ")
	d.props = tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	d.props.SetBorder(true).SetTitle(" MRA settings / Enter: edit / g: Get ")
	d.logs = tview.NewTextView().SetWrap(false).SetScrollable(true)
	d.logs.SetBorder(true).SetTitle(" Protocol events (256 max) / INF arrival only ")
	d.header = tview.NewTextView().SetText(" UECHO CONTROLLER | " + mode + " | MRA 1.3.0 | fresh Get + raw EDT")
	d.status = tview.NewTextView()
	d.search = tview.NewInputField().SetLabel(" / Filter IP / EOJ / class (empty = all): ")
	d.search.SetChangedFunc(func(string) { d.refresh(true) })
	d.search.SetDoneFunc(func(tcell.Key) { d.App.SetFocus(d.devices) })
	d.devices.SetSelectionChangedFunc(func(row, _ int) {
		if row > 0 && row <= len(d.visible) {
			d.selectDevice(d.visible[row-1].Target)
			d.refreshProperties()
		}
	})
	d.devices.SetSelectedFunc(func(int, int) { d.queueLoad(d.selected) })
	d.props.SetSelectionChangedFunc(func(row, _ int) {
		cell := d.props.GetCell(row, 0)
		if ep, ok := cell.GetReference().(byte); ok {
			d.ep = ep
		}
	})
	d.props.SetSelectedFunc(func(int, int) { d.write() })
	d.body = tview.NewFlex()
	footer := tview.NewTextView().SetText("Device select: Get | Enter/w: edit | g: Get | r: refresh\nTab/Shift-Tab | d/F5: rediscover | /: filter | ?: help\nEsc: cancel | q/Ctrl-C: exit | l: full logs")
	d.root = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(d.header, 1, 0, false).AddItem(d.search, 1, 0, false).AddItem(d.body, 0, 1, true).AddItem(d.logs, 9, 0, false).AddItem(footer, 3, 0, false).AddItem(d.status, 2, 0, false)
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
		if w < 120 || h < 28 {
			d.root.ResizeItem(d.logs, max(6, h/3), 0)
			d.setDeviceNameWidth(w)
			d.body.SetDirection(tview.FlexRow).AddItem(d.devices, 6, 0, false).AddItem(d.props, 0, 1, false)
		} else {
			d.root.ResizeItem(d.logs, max(9, h/3), 0)
			deviceWidth := min(max(48, 2*w/5), w-64)
			d.setDeviceNameWidth(deviceWidth)
			d.body.SetDirection(tview.FlexColumn).AddItem(d.devices, deviceWidth, 0, false).AddItem(d.props, 0, 1, false)
		}
		d.refresh(false)
		return false
	})
	d.devices.SetBorderColor(tcell.ColorAqua)
	d.refresh(true)
	return d
}

// Keep the identity columns visible; the name expands only inside its pane.
func (d *Dashboard) setDeviceNameWidth(paneWidth int) {
	width := max(1, paneWidth-25) // borders, 15-column IPv4, 6-column EOJ, two gaps
	if width == d.deviceNameWidth {
		return
	}
	d.deviceNameWidth = width
	for row := 1; row < d.devices.GetRowCount(); row++ {
		d.devices.GetCell(row, 2).SetMaxWidth(width)
	}
}

// LoadOnSelection enables automatic fresh maps/Get in the interactive CLI.
func (d *Dashboard) LoadOnSelection() { d.autoLoad = true; d.queueLoad(d.selected) }
func (d *Dashboard) selectDevice(target controller.Target) {
	if target == d.selected {
		return
	}
	d.selected = target
	if d.autoLoad {
		d.queueLoad(target)
	}
}
func (d *Dashboard) queueLoad(target controller.Target) {
	d.pendingLoad = target
	if d.busy && d.readJob && d.jobCancel != nil {
		d.jobCancel()
	}
	// Pending reads run on the event loop after the previous worker releases.
}
func (d *Dashboard) startRead(f func(context.Context) error) {
	if d.busy {
		return
	}
	d.start(f)
	d.readJob = true
}

func (d *Dashboard) refresh(force bool) {
	select {
	case <-d.finished:
		d.busy = false
		d.jobCancel = nil
	default:
	}
	if !d.busy && d.pendingLoad.EOJ != 0 {
		target := d.pendingLoad
		d.pendingLoad = controller.Target{}
		d.startRead(func(ctx context.Context) error { return d.session.Load(ctx, target) })
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
	for col, label := range []string{"IP", "EOJ", "Device"} {
		d.devices.SetCell(0, col, tview.NewTableCell(label).SetSelectable(false).SetTextColor(tcell.ColorAqua))
	}
	selection := 0
	for i, v := range d.visible {
		d.devices.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%-15s", v.Target.IP)).SetMaxWidth(15))
		d.devices.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("%06X", uint(v.Target.EOJ))).SetMaxWidth(6))
		d.devices.SetCell(i+1, 2, tview.NewTableCell(v.Name()).SetExpansion(1).SetMaxWidth(max(1, d.deviceNameWidth)))
		if v.Target == d.selected {
			selection = i + 1
		}
	}
	if selection == 0 && len(d.visible) > 0 {
		selection = min(max(row, 1), len(d.visible))
		d.selectDevice(d.visible[selection-1].Target)
	}
	if len(d.visible) == 0 {
		d.selectDevice(controller.Target{})
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
	if panel, ok := d.modal.(*tview.TextView); ok {
		panel.SetText(b.String() + "\nRESULT: " + snap.Status)
	}
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
		d.props.SetTitle(" MRA settings / Enter: edit / g: Get ")
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
		if value.State != "" {
			state = value.State
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
		if e.Key() == tcell.KeyEnter && d.submit != nil {
			if _, choosing := d.App.GetFocus().(*tview.List); !choosing {
				d.submit()
				return nil
			}
		}
		// Leaving an open enum list commits its visible candidate. Otherwise
		// tview's Tab closes/moves focus while retaining the previous EDT.
		if _, ok := d.modal.(*tview.Form); ok && (e.Key() == tcell.KeyTab || e.Key() == tcell.KeyBacktab) {
			if list, ok := d.App.GetFocus().(*tview.List); ok {
				list.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), func(p tview.Primitive) { d.App.SetFocus(p) })
			}
		}
		if e.Key() == tcell.KeyEscape {
			d.closeModal()
			return nil
		}
		return e
	}
	if e.Key() == tcell.KeyEscape {
		d.pendingLoad = controller.Target{}
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
		d.confirm("Keys: Tab/Shift-Tab focus; arrows select/scroll; select device loads maps/Get; Enter or w typed SetC; g fresh property Get; r refresh device maps/Get; l full wrapped logs; d/F5 rediscover all; / filter listed IP/EOJ/class (empty = all); Esc cancel; q/Ctrl-C exit.\n\nMRA state/number/raw schemas provide names and editors. Unsupported types remain raw and read only. Device maps determine availability; Get RX shows snapshot freshness. INF arrival is not a device timestamp. Readable SetC has fresh Get readback; write-only acknowledgment remains unverified. Cancel never rolls back an applied write.", func() {})
		return nil
	case 'l':
		d.showLogs()
		return nil
	case 'g':
		d.get()
		return nil
	case 'r':
		d.queueLoad(d.selected)
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

// Full-screen, wrapped copy preserves complete TX/RX frames and operation results.
func (d *Dashboard) showLogs() {
	d.prior = d.App.GetFocus()
	detail := tview.NewTextView().SetWrap(true).SetScrollable(true)
	detail.SetBorder(true).SetTitle(" Full protocol log / arrows, PgUp/PgDn / Esc close ")
	detail.SetText(d.logs.GetText(false) + "\nRESULT: " + d.session.Snapshot().Status)
	detail.ScrollToEnd()
	d.modal = detail
	d.pages.AddPage("dialog", detail, true, true)
	d.App.SetFocus(detail)
}

func (d *Dashboard) closeModal() {
	d.pages.RemovePage("dialog")
	d.modal = nil
	d.submit = nil
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
	if d.busy {
		d.session.Status("Busy: wait or Esc cancels current request")
		return
	}
	d.startRead(func(ctx context.Context) error { return d.session.Get(ctx, v.Target, ep) })
}
func (d *Dashboard) write() {
	v, ok := d.device()
	if !ok {
		return
	}
	ep := d.ep
	if d.busy {
		d.session.Status("Busy: wait or Esc cancels current request")
		return
	}
	if !v.Maps || !slices.Contains(v.Set, ep) {
		d.session.Status(fmt.Sprintf("Not editable EPC %02X: read only or maps not loaded; select device / r refresh; g reads Get-permitted values", ep))
		return
	}
	definition, known := v.Definitions()[ep]
	if !known || !definition.Editable() {
		d.session.Status("Not editable: unsupported/ambiguous MRA schema")
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
	updatePreview := func() {}
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
		form.AddDropDown("Value (arrows apply; Space opens list)", labels, selected, func(_ string, i int) { selected = i; updatePreview() })
		dropdown := form.GetFormItem(0).(*tview.DropDown)
		dropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			i, _ := dropdown.GetCurrentOption()
			switch event.Key() {
			case tcell.KeyDown:
				i = min(i+1, len(labels)-1)
			case tcell.KeyUp:
				i = max(i-1, 0)
			default:
				return event
			}
			dropdown.SetCurrentOption(i)
			return nil
		})
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
		input.SetLabel(label + ": ").SetText(current).SetChangedFunc(func(string) { updatePreview() })
		form.AddFormItem(input)
	}
	submit := func() {
		if d.selected != v.Target {
			d.closeModal()
			d.session.Status("Canceled: device selection changed; reopen the editor")
			return
		}
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
		d.start(func(ctx context.Context) error { return d.session.Set(ctx, v.Target, ep, data) })
	}
	form.AddTextView("Send to", "", 0, 5, false, false)
	preview := form.GetFormItem(form.GetFormItemCount() - 1).(*tview.TextView)
	updatePreview = func() {
		value := "select a value"
		if selected >= 0 && selected < len(options) {
			value = fmt.Sprintf("%s / EDT %X", definition.Decode(options[selected].Data), options[selected].Data)
		} else if hasInput {
			data, err := field.Encode(input.GetText())
			if err != nil {
				value = "Invalid: " + err.Error()
			} else {
				value = fmt.Sprintf("%s / EDT %X", definition.Decode(data), data)
			}
		}
		readback := "Fresh Get verifies the result"
		if !slices.Contains(v.Get, ep) {
			readback = "Write-only: acknowledged result remains unverified"
		}
		preview.SetText(fmt.Sprintf("%s\nEPC %02X %s\n%s\n%s\nReturn: send current value / Esc: cancel", v.Target, ep, definition.Name, value, readback))
	}
	updatePreview()
	form.SetBorder(true).SetTitle(fmt.Sprintf("SetC %02X / %s / MRA 1.3.0", ep, definition.Name))
	form.SetFocus(0)
	if hasInput && selected == len(options) {
		form.SetFocus(1)
	}
	d.show(form)
	d.submit = submit
}
func (d *Dashboard) start(f func(context.Context) error) {
	if d.busy {
		d.session.Status("Busy: Esc cancels current request")
		return
	}
	d.busy = true
	d.readJob = false
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
