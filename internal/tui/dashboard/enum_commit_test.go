package dashboard

import (
	"bytes"
	"context"
	"github.com/cybergarage/uecho-go/internal/tui/controller"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"testing"
	"time"
)

// Route keys to the actual focused widget, as Application.Run does.
func focusedKey(d *Dashboard, k tcell.Key) {
	e := d.capture(tcell.NewEventKey(k, 0, tcell.ModNone))
	if e != nil {
		d.App.GetFocus().InputHandler()(e, func(p tview.Primitive) { d.App.SetFocus(p) })
	}
}
func TestEnumArrowCommitsDisplayedValue(t *testing.T) {
	d, s := setup(t)
	drawScreen(d, s)
	d.App.SetFocus(d.props)
	d.ep = 0x80
	focusedKey(d, tcell.KeyEnter)
	form := d.modal.(*tview.Form)
	form.SetFocus(0)
	d.App.SetFocus(form)
	focusedKey(d, tcell.KeyDown)
	option, label := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
	if option != 1 {
		t.Fatalf("displayed OFF candidate was not committed: index=%d label=%s", option, label)
	}
	focusedKey(d, tcell.KeyUp)
	option, _ = form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
	if option != 0 {
		t.Fatal("ON not committed", option)
	}
}

func TestOpenEnumTabCommitsAndCancelSendsNothing(t *testing.T) {
	d, s := setup(t)
	drawScreen(d, s)
	d.App.SetFocus(d.props)
	d.ep = 0x80
	before := len(d.session.Client.Events())
	focusedKey(d, tcell.KeyEnter)
	form := d.modal.(*tview.Form)
	d.App.GetFocus().InputHandler()(tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone), func(p tview.Primitive) { d.App.SetFocus(p) }) // open choices
	focusedKey(d, tcell.KeyDown)
	focusedKey(d, tcell.KeyTab) // commit visible candidate and leave list
	i, label := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
	if i != 1 {
		t.Fatalf("Tab retained old enum %d %s", i, label)
	}
	focusedKey(d, tcell.KeyEscape) // Cancel
	if d.modal != nil || len(d.session.Client.Events()) != before {
		t.Fatal("cancel changed device")
	}
}
func TestFullLogsWrapAndPreserveResult(t *testing.T) {
	d, s := setup(t)
	drawScreen(d, s)
	d.session.Status("Readback mismatch: requested EDT 30 / received EDT 31")
	focusedKey(d, tcell.KeyRune) // no action
	d.showLogs()
	drawScreen(d, s)
	panel, ok := d.modal.(*tview.TextView)
	if !ok || !strings.Contains(panel.GetText(false), "requested EDT 30 / received EDT 31") {
		t.Fatal("missing result detail")
	}
	_, _, w, h := panel.GetRect()
	if w != 132 || h != 36 {
		t.Fatalf("not fullscreen %dx%d", w, h)
	}
	focusedKey(d, tcell.KeyEscape)
	if d.modal != nil {
		t.Fatal("log panel did not close")
	}
}

func waitWrite(t *testing.T, d *Dashboard) {
	t.Helper()
	select {
	case <-d.finished:
		d.busy = false
	case <-time.After(time.Second):
		t.Fatal("write did not finish")
	}
}
func TestEditorEnterDoesNotCascadeAndRepeatedSendIsBusy(t *testing.T) {
	d, s := setup(t)
	drawScreen(d, s)
	d.App.SetFocus(d.props)
	d.ep = 0x80
	before := len(d.session.Client.Events())
	focusedKey(d, tcell.KeyEnter)
	form := d.modal.(*tview.Form)
	if form.GetButtonCount() != 0 || len(d.session.Client.Events()) != before {
		t.Fatal("opening editor sent or retained buttons")
	}
	// Space opens the list; Enter commits only, retaining the editor.
	d.App.GetFocus().InputHandler()(tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone), func(p tview.Primitive) { d.App.SetFocus(p) })
	focusedKey(d, tcell.KeyDown)
	focusedKey(d, tcell.KeyEnter)
	if d.modal != form || len(d.session.Client.Events()) != before {
		t.Fatal("selecting enum cascaded into send")
	}
	focusedKey(d, tcell.KeyEnter)
	for i := 0; i < 5; i++ {
		focusedKey(d, tcell.KeyEnter)
	}
	waitWrite(t, d)
	count := 0
	for _, event := range d.session.Client.Events()[before:] {
		if event.Kind == "TX" && strings.Contains(event.Hex, "61018001") {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("duplicate SET count=%d", count)
	}
	if got := d.session.Snapshot().Devices[0].Values[0x80].Data; !bytes.Equal(got, []byte{0x31}) {
		t.Fatalf("OFF EDT %X", got)
	}
}
func TestMultipleEnumOptionsAndSelectionChange(t *testing.T) {
	d, s := setup(t)
	target := controller.Target{IP: "127.0.0.1", EOJ: 0x013001}
	d.session.Add(target)
	if err := d.session.Load(context.Background(), target); err != nil {
		t.Fatal(err)
	}
	d.selected = target
	drawScreen(d, s)
	d.ep = 0xb0
	d.App.SetFocus(d.props)
	focusedKey(d, tcell.KeyEnter)
	form := d.modal.(*tview.Form)
	dropdown := form.GetFormItem(0).(*tview.DropDown)
	for i := 0; i < 3; i++ {
		focusedKey(d, tcell.KeyDown)
	}
	index, _ := dropdown.GetCurrentOption()
	if index != 2 {
		t.Fatal("multi enum selection", index)
	}
	before := len(d.session.Client.Events())
	d.selected = controller.Target{IP: "127.0.0.1", EOJ: 0x029001}
	focusedKey(d, tcell.KeyEnter)
	if d.modal != nil || len(d.session.Client.Events()) != before {
		t.Fatal("stale selection wrote")
	}
	d.selected = target
	drawScreen(d, s)
	d.ep = 0xb0
	d.App.SetFocus(d.props)
	focusedKey(d, tcell.KeyEnter)
	for i := 0; i < 3; i++ {
		focusedKey(d, tcell.KeyDown)
	}
	var desired []byte
	for _, device := range d.session.Snapshot().Devices {
		if device.Target == target {
			desired = device.Definitions()[0xb0].Options()[2].Data
		}
	}
	focusedKey(d, tcell.KeyEnter)
	waitWrite(t, d)
	for _, device := range d.session.Snapshot().Devices {
		if device.Target == target && !bytes.Equal(device.Values[0xb0].Data, desired) {
			t.Fatalf("multi enum EDT=%X expected=%X", device.Values[0xb0].Data, desired)
		}
	}
}
