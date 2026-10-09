//go:build js && wasm

// Command cuppa-web is the Cuppa designer compiled to WebAssembly. The page
// shows it in the same TReactUI terminal the desktop app uses; there is no
// server, the terminal app runs inside the browser tab.
package main

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/shell"
	"github.com/meta-tui/cuppa/libs/catalog/standard"
	ttygo "github.com/meta-tui/treactui/packages/tty-go"
)

func newModel() tea.Model {
	app := shell.New(standard.Default())
	app.Welcome()
	return app
}

func main() {
	shared := ttygo.NewSharedProgram(newModel)
	ttygo.BindShared(context.Background(), shared, newPageEvents(), ttygo.BindOptions{})
	select {} // keep the module alive; everything happens in callbacks
}
