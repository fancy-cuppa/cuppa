// Package inspector is the right bar: details of the selection (position,
// size, layer order, properties) and the layer list.
package inspector

import (
	"strconv"
	"strings"

	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/pointer"
	"github.com/fancy-cuppa/cuppa/libs/canvas/editor"
	"github.com/fancy-cuppa/cuppa/libs/catalog/definition"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
)

// Catalog is what the inspector needs to know about components.
type Catalog interface {
	Get(id string) (definition.Definition, bool)
}

// region is a clickable span of one content line.
type region struct {
	y, x0, x1 int
	act       func()
}

// Model is the inspector pane.
type Model struct {
	ed      *editor.Editor
	cat     Catalog
	w, h    int
	scroll  int
	regions []region
	total   int // content height of the last render

	editing string // field being typed into: "x", "y", "w", "h", "name" or "prop:<key>"
	buf     string
	message string
}

// New returns an inspector editing through ed.
func New(ed *editor.Editor, cat Catalog) *Model { return &Model{ed: ed, cat: cat} }

// SetSize sets the pane size in cells.
func (m *Model) SetSize(w, h int) { m.w, m.h = w, h }

// Editing reports whether a field is taking keyboard input.
func (m *Model) Editing() bool { return m.editing != "" }

// Handle processes a pointer event in pane coordinates.
func (m *Model) Handle(e pointer.Event) {
	switch e.Phase {
	case pointer.Wheel:
		m.scroll = min(max(m.scroll+e.WheelY*2, 0), max(m.total-m.h, 0))
	case pointer.Down:
		if !e.Left {
			return
		}
		y := e.Y + m.scroll
		for _, r := range m.regions {
			if r.y == y && e.X >= r.x0 && e.X < r.x1 {
				m.cancelEdit()
				r.act()
				return
			}
		}
		m.cancelEdit()
	}
}

// Key handles typing into the field being edited.
func (m *Model) Key(text string, backspace, enter, escape bool) {
	switch {
	case escape:
		m.cancelEdit()
	case enter:
		m.commit()
	case backspace:
		if r := []rune(m.buf); len(r) > 0 {
			m.buf = string(r[:len(r)-1])
		}
	default:
		m.buf += text
	}
}

func (m *Model) startEdit(field, initial string) {
	m.editing, m.buf, m.message = field, initial, ""
}

func (m *Model) cancelEdit() { m.editing, m.buf = "", "" }

func (m *Model) commit() {
	field, value := m.editing, m.buf
	m.cancelEdit()
	n, ok := m.ed.Primary()
	if !ok {
		return
	}
	switch {
	case field == "name":
		m.ed.Rename(n.ID, value)
	case strings.HasPrefix(field, "prop:"):
		m.report(m.ed.SetProp(n.ID, strings.TrimPrefix(field, "prop:"), value))
	default:
		v, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			m.message = "enter a whole number"
			return
		}
		m.setGeometry(n, field, v)
	}
}

func (m *Model) report(err error) {
	m.message = ""
	if err != nil {
		m.message = strings.TrimPrefix(err.Error(), "editor: ")
	}
}

func (m *Model) setGeometry(n design.Node, field string, v int) {
	r := n.Rect
	switch field {
	case "x":
		r.X = v
	case "y":
		r.Y = v
	case "w":
		r.W = v
	case "h":
		r.H = v
	}
	m.ed.SetRect(n.ID, r, true)
}

func (m *Model) nudge(field string, delta int) {
	n, ok := m.ed.Primary()
	if !ok {
		return
	}
	switch field {
	case "x":
		m.setGeometry(n, field, n.Rect.X+delta)
	case "y":
		m.setGeometry(n, field, n.Rect.Y+delta)
	case "w":
		m.setGeometry(n, field, n.Rect.W+delta)
	case "h":
		m.setGeometry(n, field, n.Rect.H+delta)
	}
}
