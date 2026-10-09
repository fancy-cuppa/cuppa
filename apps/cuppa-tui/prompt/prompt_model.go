// Package prompt is a dialog that asks for one line of text, such as the name
// of a new component.
package prompt

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/a11y"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/modal"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/pointer"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/theme"
	"github.com/meta-tui/cuppa/libs/document/design"
)

const (
	boxW   = 50
	okText = "[ OK ]"
	noText = "[ Cancel ]"
)

// Model is one open prompt. Enter confirms a non-empty text; Esc cancels.
type Model struct {
	title, label string
	text         []rune
	hint         string

	rect    design.Rect
	hot     string
	outcome modal.Outcome
	done    bool
}

// New returns a prompt with the text pre-filled.
func New(title, label, initial string) *Model {
	return &Model{title: title, label: label, text: []rune(initial)}
}

// Place implements modal.Modal.
func (m *Model) Place(sw, sh int) {
	w := min(boxW, sw)
	m.rect = design.Rect{X: max((sw-w)/2, 0), Y: max((sh-7)/3, 0), W: w, H: 7}
}

// Rect implements modal.Modal.
func (m *Model) Rect() design.Rect { return m.rect }

// Lines implements modal.Modal.
func (m *Model) Lines() []string {
	inner := m.rect.W - 6
	field := theme.Fit(string(m.text)+"█", inner)
	hint := ""
	if m.hint != "" {
		hint = theme.Faded(m.hint)
	}
	row := m.button("ok", okText) + " " + m.button("no", noText)
	pad := max(m.rect.W-2-ansi.StringWidth(row)-1, 0)
	content := []string{" " + m.label, "  " + field, "  " + hint, "", strings.Repeat(" ", pad) + row}
	return theme.Panel(m.title, content, m.rect.W)
}

func (m *Model) button(id, label string) string {
	if m.hot == id {
		return theme.Selected(label)
	}
	return theme.Button(label, true)
}

// buttonAt says which button the box column x is over.
func (m *Model) buttonAt(x int) string {
	noW, okW := ansi.StringWidth(noText), ansi.StringWidth(okText)
	noX := m.rect.W - 2 - noW
	okX := noX - 1 - okW
	switch {
	case x >= noX && x < noX+noW:
		return "no"
	case x >= okX && x < okX+okW:
		return "ok"
	}
	return ""
}

// Handle implements modal.Modal.
func (m *Model) Handle(e pointer.Event) {
	if m.done {
		return
	}
	x, y := e.X-m.rect.X, e.Y-m.rect.Y
	m.hot = ""
	if y != m.rect.H-2 {
		return
	}
	m.hot = m.buttonAt(x)
	if e.Phase == pointer.Down && e.Left {
		switch m.hot {
		case "ok":
			m.confirm()
		case "no":
			m.finish(true)
		}
	}
}

// Key implements modal.Modal.
func (m *Model) Key(text string, back, enter, esc bool) {
	switch {
	case esc:
		m.finish(true)
	case enter:
		m.confirm()
	case back:
		if len(m.text) > 0 {
			m.text = m.text[:len(m.text)-1]
		}
	case text != "":
		m.text = append(m.text, []rune(text)...)
	}
}

func (m *Model) confirm() {
	if strings.TrimSpace(string(m.text)) == "" {
		m.hint = "Type a name first"
		return
	}
	m.finish(false)
}

func (m *Model) finish(canceled bool) {
	m.done = true
	m.outcome = modal.Outcome{Button: "OK", Value: strings.TrimSpace(string(m.text)), Canceled: canceled}
}

// Outcome implements modal.Modal.
func (m *Model) Outcome() (modal.Outcome, bool) { return m.outcome, m.done }

// Describe reads the dialog out: its question, the text so far and the buttons.
func (m *Model) Describe() []a11y.Node {
	field := a11y.Field(m.label, string(m.text))
	field.Focused = true
	nodes := []a11y.Node{a11y.Heading(m.title), field}
	if m.hint != "" {
		nodes = append(nodes, a11y.Text(m.hint))
	}
	return append(nodes, a11y.Button("OK"), a11y.Button("Cancel"))
}
