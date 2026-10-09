package main

import (
	"context"
	"sync/atomic"

	tea "charm.land/bubbletea/v2"
	"github.com/meta-tui/cuppa/apps/cuppa-desktop/terminal"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/shell"
	"github.com/meta-tui/cuppa/libs/catalog/standard"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App ties the window to the terminal app. Its exported methods are what the
// page can call.
type App struct {
	ctx      context.Context
	server   *terminal.Server
	startErr error
	// closing is set once the app itself has decided to quit, so closing the
	// window then goes ahead instead of asking the app again.
	closing atomic.Bool
}

func newApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.server, a.startErr = terminal.Start(a.newModel)
}

// newModel builds the Cuppa terminal app. The window closes when the user
// quits from inside it.
func (a *App) newModel() tea.Model {
	app := shell.New(standard.Default())
	app.Welcome()
	return terminal.OnQuit(app, func() {
		a.closing.Store(true)
		runtime.Quit(a.ctx)
	})
}

// TerminalURL is the address the page connects its terminal to.
func (a *App) TerminalURL() (string, error) {
	if a.startErr != nil {
		return "", a.startErr
	}
	return a.server.URL(), nil
}

// beforeClose turns "the user closed the window" into Ctrl+Q inside the app,
// so its own unsaved-changes prompt runs. It returns true to keep the window
// open until the app quits itself.
func (a *App) beforeClose(context.Context) bool {
	if a.closing.Load() || a.server == nil {
		return false
	}
	return a.server.Send(tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl})
}

func (a *App) shutdown(context.Context) {
	if a.server != nil {
		_ = a.server.Close()
	}
}
