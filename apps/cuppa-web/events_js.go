//go:build js && wasm

package main

import (
	"sync"
	"syscall/js"
)

// pageEvents is the event bus between this program and the page, for tty-go's
// Bind. The page reaches Go through globalThis.cuppaBridge.up(name, data) and
// Go reaches the page through cuppaBridge.down(name, data); web_runtime.ts
// turns those into the EventsOn/EventsEmit that createWailsSocket expects.
type pageEvents struct {
	mu       sync.Mutex
	handlers map[string][]func(string)
}

func newPageEvents() *pageEvents {
	e := &pageEvents{handlers: map[string][]func(string){}}
	bridge := js.Global().Get("cuppaBridge")
	if bridge.IsUndefined() {
		bridge = js.Global().Get("Object").New()
		js.Global().Set("cuppaBridge", bridge)
	}
	bridge.Set("up", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) == 2 {
			// A callback runs on the event loop: hand the work to a goroutine so a
			// handler that waits cannot freeze the page.
			go e.dispatch(args[0].String(), args[1].String())
		}
		return nil
	}))
	return e
}

func (e *pageEvents) dispatch(name, data string) {
	e.mu.Lock()
	list := append([]func(string){}, e.handlers[name]...)
	e.mu.Unlock()
	for _, h := range list {
		h(data)
	}
}

// On calls handler with the string each event called name carries.
func (e *pageEvents) On(name string, handler func(string)) func() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[name] = append(e.handlers[name], handler)
	index := len(e.handlers[name]) - 1
	return func() {
		e.mu.Lock()
		defer e.mu.Unlock()
		if index < len(e.handlers[name]) {
			e.handlers[name][index] = func(string) {}
		}
	}
}

// Emit sends an event to the page.
func (e *pageEvents) Emit(name, data string) {
	if down := js.Global().Get("cuppaBridge").Get("down"); down.Type() == js.TypeFunction {
		down.Invoke(name, data)
	}
}
