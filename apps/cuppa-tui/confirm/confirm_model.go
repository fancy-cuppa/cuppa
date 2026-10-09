// Package confirm is a small dialog with a message and a row of buttons:
// "Save changes?", "Replace file?", error notices.
package confirm

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/modal"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/pointer"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/theme"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
)

const (
	minW = 34
	maxW = 64
)

// Model is one open dialog. The first button is the default (Enter).
type Model struct {
	title   string
	body    []string
	buttons []string

	rect    design.Rect
	hover   int
	spans   []span
	outcome modal.Outcome
	done    bool
}

type span struct{ x0, x1 int }

// New returns a dialog. With no buttons it gets a single "OK".
func New(title, body string, buttons ...string) *Model {
	if len(buttons) == 0 {
		buttons = []string{"OK"}
	}
	return &Model{title: title, body: strings.Split(body, "\n"), buttons: buttons, hover: -1}
}

// Place implements modal.Modal.
func (m *Model) Place(sw, sh int) {
	w := minW
	for _, l := range m.body {
		w = max(w, ansi.StringWidth(l)+6)
	}
	w = min(w, min(maxW, sw))
	lines := m.wrapped(w - 4)
	h := len(lines) + 5
	m.rect = design.Rect{X: max((sw-w)/2, 0), Y: max((sh-h)/3, 0), W: w, H: h}
}

// Rect implements modal.Modal.
func (m *Model) Rect() design.Rect { return m.rect }

func (m *Model) wrapped(w int) []string {
	var out []string
	for _, l := range m.body {
		out = append(out, strings.Split(ansi.Wrap(l, w, ""), "\n")...)
	}
	return out
}

// Lines implements modal.Modal.
func (m *Model) Lines() []string {
	body := m.wrapped(m.rect.W - 4)
	content := []string{""}
	for _, l := range body {
		content = append(content, "  "+l)
	}
	row, spans := m.buttonRow()
	m.spans = spans
	content = append(content, "", row)
	return theme.Panel(m.title, content, m.rect.W)
}

// buttonRow lays the buttons out right-aligned inside the box and records
// their columns (relative to the box).
func (m *Model) buttonRow() (string, []span) {
	var parts []string
	var spans []span
	total := 0
	for _, b := range m.buttons {
		total += ansi.StringWidth(b) + 4 + 1
	}
	x := max(m.rect.W-2-total, 1)
	line := strings.Repeat(" ", x)
	col := 1 + x
	for i, b := range m.buttons {
		label := " " + b + " "
		style := theme.Button
		text := "[" + label + "]"
		if i == m.hover {
			text = theme.Selected(text)
		} else {
			text = style(text, true)
		}
		parts = append(parts, text)
		spans = append(spans, span{col, col + ansi.StringWidth(text)})
		col += ansi.StringWidth(text) + 1
	}
	return line + strings.Join(parts, " "), spans
}

// Handle implements modal.Modal.
func (m *Model) Handle(e pointer.Event) {
	if m.done {
		return
	}
	x, y := e.X-m.rect.X, e.Y-m.rect.Y
	idx := -1
	if y == m.rect.H-2 {
		for i, s := range m.spans {
			if x >= s.x0 && x < s.x1 {
				idx = i
			}
		}
	}
	switch e.Phase {
	case pointer.Move:
		m.hover = idx
	case pointer.Down:
		if e.Left && idx >= 0 {
			m.finish(m.buttons[idx], false)
		}
	}
}

// Key implements modal.Modal.
func (m *Model) Key(_ string, _, enter, esc bool) {
	switch {
	case esc:
		m.finish("", true)
	case enter:
		m.finish(m.buttons[0], false)
	}
}

func (m *Model) finish(button string, canceled bool) {
	m.done = true
	m.outcome = modal.Outcome{Button: button, Canceled: canceled}
}

// Outcome implements modal.Modal.
func (m *Model) Outcome() (modal.Outcome, bool) { return m.outcome, m.done }
