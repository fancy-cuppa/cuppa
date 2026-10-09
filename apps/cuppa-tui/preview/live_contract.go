// Package preview runs a design: components that have a real Bubbles model
// behind them (text inputs, lists, spinners…) are started as that model and
// respond to the keyboard, the mouse and time, while the rest stay as the
// designer draws them. It is a mode of the app, not part of the headless
// engine, because the models are Bubble Tea models.
package preview

import (
	tea "charm.land/bubbletea/v2"
)

// live is one running component.
type live struct {
	// update gives the component a message and returns what it wants run.
	update func(tea.Msg) tea.Cmd
	// view is the component's current screen, ANSI styled.
	view func() string
	// focus gives or takes the keyboard; nil when the component has no focus state.
	focus func(on bool) tea.Cmd
	// start is what to run when the preview begins (a blink, a tick).
	start tea.Cmd
	// keys is true when the component reads the keyboard; clicking it focuses it.
	keys bool
}

// props are a node's property values with the catalog defaults underneath.
type props map[string]string

func (p props) str(key string) string { return p[key] }

func (p props) flag(key string) bool { return p[key] == "true" }

func (p props) integer(key string, fallback int) int {
	var n int
	if _, err := fmtSscan(p[key], &n); err != nil {
		return fallback
	}
	return n
}
