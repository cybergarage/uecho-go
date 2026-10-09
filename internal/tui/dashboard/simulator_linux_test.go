package dashboard

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cybergarage/uecho-go/internal/tui/testutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/internal/tui/controller"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CI runs this only inside an isolated Linux network namespace. The existing
// user's simulator and LAN are never involved; SSE is the GUI's actual input.
func TestSimulatorUIEnumToDisplay(t *testing.T) {
	binary := os.Getenv("UECHOTUI_SIMULATOR_V1")
	if binary == "" || runtime.GOOS != "linux" || os.Getenv("UECHOTUI_ISOLATED_MULTICAST") != "1" {
		t.Skip("isolated Linux simulator opt-in")
	}
	testutil.StartSimulator(t, binary, "http://127.0.0.1:18080/", "--display", "127.0.0.1:18080", "--udp", "127.0.0.2:3610")
	c, err := controller.Open("lo", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	session := controller.NewSession(c)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err = session.Discover(ctx, "127.0.0.2"); err != nil {
		t.Fatal(err)
	}
	for _, device := range session.Snapshot().Devices {
		if err = session.Load(ctx, device.Target); err != nil {
			t.Fatal(err)
		}
	}
	d := New(session, "127.0.0.2", "ISOLATED SIMULATOR")
	defer func() { d.cancel(); d.jobs.Wait() }()
	screen := tcell.NewSimulationScreen("UTF-8")
	d.App.SetScreen(screen)
	defer screen.Fini()
	screen.SetSize(132, 36)
	target := controller.Target{IP: "127.0.0.2", EOJ: 0x029001}
	d.selected = target
	drawScreen(d, screen)
	for _, want := range []byte{0x31, 0x30, 0x31} {
		d.ep = 0x80
		d.App.SetFocus(d.props)
		focusedKey(d, tcell.KeyEnter)
		form := d.modal.(*tview.Form)
		dropdown := form.GetFormItem(0).(*tview.DropDown)
		index, _ := dropdown.GetCurrentOption()
		desired := int(want - 0x30)
		if index != desired {
			if desired == 0 {
				focusedKey(d, tcell.KeyUp)
			} else {
				focusedKey(d, tcell.KeyDown)
			}
		}
		before := len(c.Events())
		focusedKey(d, tcell.KeyEnter)
		if d.modal != nil {
			t.Fatal("send opened an extra confirmation")
		}
		select {
		case <-d.finished:
			d.busy = false
		case <-time.After(3 * time.Second):
			t.Fatal("SET did not complete")
		}
		snapshot := session.Snapshot()
		if !strings.Contains(snapshot.Status, "Readback success") || !strings.Contains(snapshot.Status, fmt.Sprintf("requested EDT %X / received EDT %X", want, want)) {
			t.Fatal(snapshot.Status)
		}
		var txSet, txGet uint
		for _, event := range c.Events()[before:] {
			if event.Kind == "TX" && strings.Contains(event.Hex, fmt.Sprintf("61018001%02X", want)) {
				txSet = event.TID
			}
			if event.Kind == "TX" && strings.Contains(event.Hex, "62018000") {
				txGet = event.TID
			}
		}
		if txSet == 0 || txGet == 0 || txSet == txGet {
			t.Fatalf("missing distinct SET/Get packets: %04X %04X %+v", txSet, txGet, c.Events()[before:])
		}
		response, err := (&http.Client{Timeout: 2 * time.Second}).Get("http://127.0.0.1:18080/events")
		if err != nil {
			t.Fatal(err)
		}
		scanner := bufio.NewScanner(response.Body)
		found := false
		stateEvent := false
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			var event struct {
				Snapshot struct {
					Devices []struct {
						EOJ   uint32 `json:"eoj"`
						Power bool   `json:"power"`
					}
					Events []struct{ Kind, Message string }
				}
			}
			if err = json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &event); err != nil {
				t.Fatal(err)
			}
			for _, update := range event.Snapshot.Events {
				if update.Kind == "STATE" && strings.Contains(update.Message, fmt.Sprintf("029001 / 80 = %X", want)) {
					stateEvent = true
				}
			}
			for _, device := range event.Snapshot.Devices {
				if device.EOJ == 0x029001 {
					found = true
					if device.Power != (want == 0x30) {
						t.Fatalf("GUI/model power=%v EDT=%X", device.Power, want)
					}
				}
			}
			break
		}
		response.Body.Close()
		if !found || !stateEvent {
			t.Fatal("GUI SSE lacks lighting state/update event")
		}
		t.Logf("UI -> SetC %02X TID %04X -> fresh Get TID %04X -> model/SSE power=%v", want, txSet, txGet, want == 0x30)
	}
}
