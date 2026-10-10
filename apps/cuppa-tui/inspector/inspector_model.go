// Package inspector is the right bar: details of the selection (position,
// size, layer order, properties) and the layer list.
package inspector

import (
	"strconv"
	"strings"

	"github.com/meta-tui/cuppa/apps/cuppa-tui/pointer"
	"github.com/meta-tui/cuppa/libs/canvas/editor"
	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/document/design"
)

// Catalog is what the inspector needs to know about components.
type Catalog interface {
	Get(id string) (definition.Definition, bool)
}

// region is a clickable span of one content line.
type region struct {
	y, x0, x1 int
	act       func()
	// step changes the value by one in a direction (-1 or +1), for Left and
	// Right on the keyboard stop; mouseOnly spans are not keyboard stops.
	step      func(delta int)
	mouseOnly bool
	// label is the stop as it is read out: the text before it on its line and
	// its own text.
	label string
	// layer is set on a layer name: pressing it can start a drag of that layer.
	layer design.NodeID
}

// Model is the inspector pane.
type Model struct {
	ed      *editor.Editor
	cat     Catalog
	w, h    int
	focused bool
	scroll  int
	regions []region
	total   int // content height of the last render
	// stop is the keyboard stop that has the focus and stops how many there
	// were at the last render.
	stop, stops int

	editing string // field being typed into: "x", "y", "w", "h", "name" or "prop:<key>"
	buf     string
	message string

	// snap reads and changes the stage setting this pane exposes as a checkbox.
	snapGet func() bool
	snapSet func(bool)

	// pickColor opens the colour dialog; nil means colours are edited as text.
	pickColor ColorPicker
	// pickTheme opens the dialog that fills the theme from a colour scheme.
	pickTheme ThemePicker

	// The layer list: a drag in progress, where the list starts and how long it
	// is (content rows, set at render), and where the trash button is.
	drag                     *layerDrag
	layerTop, layerCount     int
	trashY, trashX0, trashX1 int
}

// New returns an inspector editing through ed.
func New(ed *editor.Editor, cat Catalog) *Model { return &Model{ed: ed, cat: cat} }

// BindSnap connects the "snap to guides" checkbox to the stage setting.
func (m *Model) BindSnap(get func() bool, set func(bool)) { m.snapGet, m.snapSet = get, set }

// SetSize sets the pane size in cells.
func (m *Model) SetSize(w, h int) { m.w, m.h = w, h }

// ColorPicker opens a colour dialog for a property: it shows title with the
// current colour selected, and calls apply with the chosen one ("" for none).
// The shell provides it.
type ColorPicker func(title, current string, apply func(color string))

// ThemePicker opens the dialog that fills the theme from a colour scheme and
// calls apply with the background and the theme colours it chose. The shell
// provides it.
type ThemePicker func(apply func(background string, t design.Theme))

// BindThemePicker sets how the Theme section's scheme button opens the dialog.
// Without it the button is not shown.
func (m *Model) BindThemePicker(p ThemePicker) { m.pickTheme = p }

// BindColorPicker sets how clicking a colour opens the picker. Without it a
// click edits the colour as text.
func (m *Model) BindColorPicker(p ColorPicker) { m.pickColor = p }

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
				if r.layer != "" {
					// Selecting changes what the bar shows above the list, which
					// would slide the list out from under a drag; so a layer is
					// selected on release, when it was not dragged.
					m.pressLayer(r.layer, y, r.act)
					return
				}
				r.act()
				return
			}
		}
		m.cancelEdit()
	case pointer.Move:
		if m.drag != nil && e.Held {
			m.dragTo(e)
		}
	case pointer.Up:
		if m.drag != nil {
			m.drop(e)
		}
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
	m.apply(field, value)
}

// Step adds delta to the number being edited and applies it at once, so the
// design follows the arrow keys; Enter then ends the edit. It reports whether
// the field holds a number.
func (m *Model) Step(delta int) bool {
	if m.editing == "" || m.editing == "background" || m.editing == "name" {
		return false
	}
	n, err := strconv.Atoi(strings.TrimSpace(m.buf))
	if err != nil {
		return false
	}
	m.buf = strconv.Itoa(n + delta)
	m.apply(m.editing, m.buf)
	return true
}

// apply sets the field to value. The field is "background", "canvas-…",
// "name", "prop:key" or one of the geometry fields.
func (m *Model) apply(field, value string) {
	switch {
	case field == "background":
		m.report(m.ed.SetBackground(value))
		return
	case strings.HasPrefix(field, "theme-"):
		m.report(m.ed.SetThemeColor(strings.TrimPrefix(field, "theme-"), value))
		return
	case strings.HasPrefix(field, "canvas-"):
		m.commitCanvas(field, value)
		return
	}
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
