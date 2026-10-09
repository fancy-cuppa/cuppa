package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// wailsEvents offers Wails' event bus to tty-go, so the terminal app talks to
// the page through the window itself: no server, no port, no origin check.
type wailsEvents struct{ ctx context.Context }

// On calls handler with the string each event called name carries.
func (e wailsEvents) On(name string, handler func(string)) func() {
	return runtime.EventsOn(e.ctx, name, func(data ...interface{}) {
		if len(data) == 0 {
			return
		}
		if s, ok := data[0].(string); ok {
			handler(s)
		}
	})
}

// Emit sends an event to the page.
func (e wailsEvents) Emit(name, data string) { runtime.EventsEmit(e.ctx, name, data) }
