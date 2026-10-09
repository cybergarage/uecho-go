package dashboard

import (
	"context"
	"fmt"
	"github.com/cybergarage/uecho-go/examples/uechotui/internal/controller"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SelectInterface lists actual local addresses without binding or sending. The
// user must explicitly confirm a selection; Cancel is the initial focus.
func SelectInterface(ctx context.Context) (controller.InterfaceOption, error) {
	options, err := controller.InterfaceOptions()
	if err != nil {
		return controller.InterfaceOption{}, err
	}
	if len(options) == 0 {
		return controller.InterfaceOption{}, fmt.Errorf("no up multicast IPv4 interface; use explicit isolated --interface/--bind/--peer")
	}
	return selectInterface(ctx, options, nil)
}
func selectInterface(ctx context.Context, options []controller.InterfaceOption, screen tcell.Screen) (controller.InterfaceOption, error) {
	app := tview.NewApplication()
	if screen != nil {
		app.SetScreen(screen)
	}
	form := tview.NewForm()
	labels := []string{}
	for _, v := range options {
		labels = append(labels, v.Name+" / "+v.IP)
	}
	selected := 0
	confirmed := false
	form.AddDropDown("Interface / IPv4", labels, 0, func(_ string, i int) { selected = i }).AddButton("Cancel", app.Stop).AddButton("Confirm receive only", func() { confirmed = true; app.Stop() })
	form.SetBorder(true).SetTitle("UDP opt-in: choose interface; d confirms discovery later")
	form.SetFocus(1)
	app.SetRoot(form, true).SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEscape || e.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}
		return e
	})
	ready := make(chan struct{}, 1)
	app.SetBeforeDrawFunc(func(tcell.Screen) bool {
		select {
		case ready <- struct{}{}:
		default:
		}
		return false
	})
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-ready:
		}
		select {
		case <-done:
		case <-ctx.Done():
			app.Stop()
		}
	}()
	err := app.Run()
	close(done)
	if err != nil {
		return controller.InterfaceOption{}, err
	}
	if !confirmed {
		return controller.InterfaceOption{}, context.Canceled
	}
	return options[selected], nil
}
