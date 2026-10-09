package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/cmd/uechotui/internal/controller"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func setup(t *testing.T) (*Dashboard, tcell.SimulationScreen) {
	t.Helper()
	c := controller.Demo()
	t.Cleanup(func() { c.Close() })
	s := controller.NewSession(c)
	if err := s.Discover(context.Background(), "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	target := s.Snapshot().Devices[0].Target
	if err := s.Load(context.Background(), target); err != nil {
		t.Fatal(err)
	}
	d := New(s, "127.0.0.1", "OFFLINE DEMO - no sockets")
	screen := tcell.NewSimulationScreen("UTF-8")
	d.App.SetScreen(screen)
	screen.SetSize(132, 36)
	t.Cleanup(func() {
		d.cancel()
		d.jobs.Wait()
		if _, w, h := screen.GetContents(); w != 0 || h != 0 {
			screen.Fini()
		}
	})
	return d, screen
}
func drawScreen(d *Dashboard, s tcell.Screen) {
	w, h := s.Size()
	s.Clear()
	d.App.GetBeforeDrawFunc()(s)
	d.pages.SetRect(0, 0, w, h)
	d.pages.Draw(s)
	s.Show()
}
func key(d *Dashboard, k tcell.Key, r rune) {
	e := tcell.NewEventKey(k, r, tcell.ModNone)
	if e = d.capture(e); e != nil {
		d.pages.InputHandler()(e, func(p tview.Primitive) { d.App.SetFocus(p) })
	}
}
func TestNavigationSearchConfirmCancelResize(t *testing.T) {
	d, s := setup(t)
	drawScreen(d, s)
	key(d, tcell.KeyTab, 0)
	if d.App.GetFocus() != d.props {
		t.Fatal("Tab")
	}
	key(d, tcell.KeyBacktab, 0)
	if d.App.GetFocus() != d.devices {
		t.Fatal("Shift Tab")
	}
	key(d, tcell.KeyRune, '/')
	key(d, tcell.KeyRune, '0')
	key(d, tcell.KeyRune, '2')
	key(d, tcell.KeyRune, '9')
	drawScreen(d, s)
	if len(d.visible) != 1 {
		t.Fatalf("search text=%q visible=%d focus=%T", d.search.GetText(), len(d.visible), d.App.GetFocus())
	}
	key(d, tcell.KeyEscape, 0)
	key(d, tcell.KeyEnter, 0)
	if d.modal == nil {
		t.Fatal("load dialog")
	}
	before := len(d.session.Client.Events())
	key(d, tcell.KeyEnter, 0)
	if d.modal != nil || len(d.session.Client.Events()) != before {
		t.Fatal("cancel must send nothing")
	}
	d.App.SetFocus(d.props)
	d.ep = 0x80
	key(d, tcell.KeyRune, 'w')
	if d.modal == nil {
		t.Fatal("write form")
	}
	key(d, tcell.KeyEscape, 0)
	if len(d.session.Client.Events()) != before {
		t.Fatal("Esc wrote")
	}
	key(d, tcell.KeyRune, '?')
	s.SetSize(55, 22)
	drawScreen(d, s)
	key(d, tcell.KeyEscape, 0)
	if d.modal != nil {
		t.Fatal("resize Esc")
	}
	d.search.SetText("no match")
	drawScreen(d, s)
	key(d, tcell.KeyRune, 'w')
	if d.modal != nil {
		t.Fatal("search operated hidden selection")
	}
}
func TestWriteReviewAndApply(t *testing.T) {
	d, s := setup(t)
	drawScreen(d, s)
	d.App.SetFocus(d.props)
	d.ep = 0x80
	d.write()
	form := d.modal.(*tview.Form)
	input := form.GetFormItem(0).(*tview.DropDown)
	input.SetCurrentOption(1)
	form.SetFocus(2)
	key(d, tcell.KeyEnter, 0)
	if _, ok := d.modal.(*tview.Modal); !ok {
		t.Fatal("missing review confirmation")
	}
	before := len(d.session.Client.Events())
	key(d, tcell.KeyRight, 0)
	key(d, tcell.KeyEnter, 0)
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		select {
		case <-d.finished:
			d.busy = false
			goto complete
		default:
			time.Sleep(time.Millisecond)
		}
	}
	t.Fatal("write did not finish")
complete:
	if len(d.session.Client.Events()) <= before || !strings.Contains(d.session.Snapshot().Status, "Readback success") {
		t.Fatal(d.session.Snapshot().Status)
	}
	if got := d.session.Snapshot().Devices[0].Values[0x80].Data[0]; got != 0x31 {
		t.Fatalf("EDT %X", got)
	}
}
func TestTerminalFinalization(t *testing.T) {
	for _, quit := range []tcell.Key{tcell.KeyCtrlC, tcell.KeyRune} {
		t.Run(fmt.Sprint(quit), func(t *testing.T) {
			d, s := setup(t)
			drawn := make(chan struct{}, 1)
			d.App.SetAfterDrawFunc(func(tcell.Screen) {
				select {
				case drawn <- struct{}{}:
				default:
				}
			})
			done := make(chan error, 1)
			go func() { done <- d.App.Run() }()
			select {
			case <-drawn:
			case <-time.After(time.Second):
				t.Fatal("no draw")
			}
			r := rune(0)
			if quit == tcell.KeyRune {
				r = 'q'
			}
			s.InjectKey(quit, r, tcell.ModNone)
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(time.Second):
				t.Fatal("exit blocked")
			}
			_, w, h := s.GetContents()
			if w != 0 || h != 0 {
				t.Fatalf("screen not finalized: %dx%d", w, h)
			}
		})
	}
}

// Export pixels of actual tcell widget cells. The optional path is used to
// regenerate documentation images, not to create an artist's mockup.
func TestScreenshot(t *testing.T) {
	d, s := setup(t)
	drawScreen(d, s)
	cells, w, h := s.GetContents()
	var text strings.Builder
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := cells[y*w+x].Runes
			if len(r) > 0 {
				text.WriteRune(r[0])
			} else {
				text.WriteByte(' ')
			}
		}
		text.WriteByte('\n')
	}
	if !strings.Contains(text.String(), "MRA settings") || !strings.Contains(text.String(), "Protocol events") {
		t.Fatal(text.String())
	}
	dir := os.Getenv("UECHOTUI_SCREENSHOT_DIR")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, shot := range []struct {
		name   string
		w, h   int
		dialog bool
		ep     byte
	}{{"tui", 132, 36, false, 0}, {"tui-compact", 68, 26, false, 0}, {"tui-set", 132, 36, true, 0x80}, {"tui-number", 132, 36, true, 0xb0}} {
		s.SetSize(shot.w, shot.h)
		if shot.dialog {
			if d.modal != nil {
				d.closeModal()
			}
			d.ep = shot.ep
			d.write()
		}
		drawScreen(d, s)
		cells, w, h = s.GetContents()
		img := image.NewRGBA(image.Rect(0, 0, w*7, h*14))
		ttf := gomono.TTF
		if path := os.Getenv("UECHOTUI_SCREENSHOT_FONT"); path != "" {
			var err error
			ttf, err = os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
		}
		parsed, err := opentype.Parse(ttf)
		if err != nil {
			collection, e := opentype.ParseCollection(ttf)
			if e != nil {
				t.Fatal(e)
			}
			parsed, err = collection.Font(0)
			if err != nil {
				t.Fatal(err)
			}
		}
		face, err := opentype.NewFace(parsed, &opentype.FaceOptions{Size: 12, DPI: 72, Hinting: font.HintingFull})
		if err != nil {
			t.Fatal(err)
		}
		drawer := font.Drawer{Dst: img, Face: face}
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				cell := cells[y*w+x]
				fg, bg, _ := cell.Style.Decompose()
				fr, fgreen, fb := fg.RGB()
				br, bg2, bb := bg.RGB()
				if fg == tcell.ColorDefault {
					fr, fgreen, fb = 220, 220, 220
				}
				if bg == tcell.ColorDefault {
					br, bg2, bb = 15, 20, 28
				}
				draw.Draw(img, image.Rect(x*7, y*14, (x+1)*7, (y+1)*14), &image.Uniform{color.RGBA{uint8(br), uint8(bg2), uint8(bb), 255}}, image.Point{}, draw.Src)
				if len(cell.Runes) > 0 && cell.Runes[0] >= 32 {
					drawer.Src = &image.Uniform{color.RGBA{uint8(fr), uint8(fgreen), uint8(fb), 255}}
					drawer.Dot = fixed.P(x*7, y*14+12)
					drawer.DrawString(string(cell.Runes))
				}
			}
		}
		f, err := os.Create(filepath.Join(dir, shot.name+".png"))
		if err != nil {
			t.Fatal(err)
		}
		if err = png.Encode(f, img); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
	data, _ := json.MarshalIndent(d.session.Snapshot(), "", "  ")
	os.WriteFile(filepath.Join(dir, "screenshot-state.json"), data, 0644)
}

// Draw callbacks execute under Application's mutex. Exercise a changed snapshot
// through the real application loop so accidental reentrant GetFocus calls are
// caught (direct widget rendering alone cannot expose this deadlock).
func TestRunningSnapshotUpdateAndExit(t *testing.T) {
	d, s := setup(t)
	drawn := make(chan string, 16)
	d.App.SetAfterDrawFunc(func(tcell.Screen) {
		select {
		case drawn <- d.session.Snapshot().Status:
		default:
		}
	})
	done := make(chan error, 1)
	go func() { done <- d.Run(context.Background()) }()
	select {
	case <-drawn:
	case <-time.After(time.Second):
		t.Fatal("first draw")
	}
	d.session.Status("changed snapshot under draw lock")
	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()
	for {
		select {
		case status := <-drawn:
			if status == "changed snapshot under draw lock" {
				goto exit
			}
		case <-deadline.C:
			t.Fatal("snapshot redraw blocked")
		}
	}
exit:
	s.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("quit blocked after changed snapshot")
	}
}

func TestContextShutdownCancelsWorker(t *testing.T) {
	d, s := setup(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	worker := make(chan struct{})
	d.start(func(ctx context.Context) error { <-ctx.Done(); close(worker); return ctx.Err() })
	ready := make(chan struct{}, 1)
	d.App.SetAfterDrawFunc(func(tcell.Screen) {
		select {
		case ready <- struct{}{}:
		default:
		}
	})
	done := make(chan error, 1)
	go func() { done <- d.Run(ctx) }()
	select {
	case <-ready:
	case <-time.After(time.Second):
		t.Fatal("not drawn")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown with worker blocked")
	}
	select {
	case <-worker:
	default:
		t.Fatal("worker not canceled")
	}
	_, w, h := s.GetContents()
	if w != 0 || h != 0 {
		t.Fatal("screen not finalized")
	}
}

func TestCompactKeysRemainVisible(t *testing.T) {
	d, s := setup(t)
	s.SetSize(68, 26)
	drawScreen(d, s)
	cells, w, h := s.GetContents()
	var text strings.Builder
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := cells[y*w+x].Runes
			if len(r) > 0 {
				text.WriteRune(r[0])
			}
		}
		text.WriteByte('\n')
	}
	for _, hint := range []string{"/ search | ? help", "Esc cancel | q / Ctrl-C exit", "Value | raw EDT"} {
		if !strings.Contains(text.String(), hint) {
			t.Fatalf("compact hint %q hidden:\n%s", hint, text.String())
		}
	}
}

func TestTypedNumberAndUnsupportedControls(t *testing.T) {
	d, screen := setup(t)
	drawScreen(d, screen)
	d.ep = 0xb0
	d.write()
	form, ok := d.modal.(*tview.Form)
	if !ok {
		t.Fatal("number editor unavailable")
	}
	input := form.GetFormItem(1).(*tview.InputField)
	before := len(d.session.Client.Events())
	input.SetText("101")
	form.GetButton(1).InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), func(tview.Primitive) {})
	if d.modal != form || len(d.session.Client.Events()) != before {
		t.Fatal("invalid value advanced or sent")
	}
	input.SetText("75")
	form.GetButton(1).InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), func(tview.Primitive) {})
	if _, ok := d.modal.(*tview.Modal); !ok {
		t.Fatal("typed value lacks confirmation")
	}
	d.closeModal()
	if len(d.session.Client.Events()) != before {
		t.Fatal("cancel transmitted")
	}
	d.ep = 0xff
	d.write()
	if d.modal != nil {
		t.Fatal("unknown editor")
	}
	d.ep = 0x81
	d.get()
	if d.modal != nil {
		t.Fatal("unsupported Get confirmation")
	}
}

func TestInterfacePickerRequiresConfirmation(t *testing.T) {
	for _, confirm := range []bool{false, true} {
		screen := tcell.NewSimulationScreen("UTF-8")
		screen.SetSize(100, 30)
		done := make(chan error, 1)
		go func() {
			choice, err := selectInterface(context.Background(), []controller.InterfaceOption{{Name: "fixture0", IP: "192.0.2.20"}}, screen)
			if confirm && err == nil && choice.IP != "192.0.2.20" {
				err = fmt.Errorf("wrong selection")
			}
			done <- err
		}()
		time.Sleep(20 * time.Millisecond)
		if confirm {
			screen.PostEventWait(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
		}
		screen.PostEventWait(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
		select {
		case err := <-done:
			if confirm && err != nil {
				t.Fatal(err)
			}
			if !confirm && err != context.Canceled {
				t.Fatal("default must cancel", err)
			}
		case <-time.After(time.Second):
			t.Fatal("picker blocked")
		}
	}
}
