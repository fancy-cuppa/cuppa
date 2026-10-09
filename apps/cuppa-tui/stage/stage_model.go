// Package stage is the centre pane: it draws the document and turns pointer
// gestures into editor commands (select, move, resize).
package stage

import (
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/pointer"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/theme"
	"github.com/fancy-cuppa/cuppa/libs/canvas/editor"
	"github.com/fancy-cuppa/cuppa/libs/canvas/hittest"
	"github.com/fancy-cuppa/cuppa/libs/catalog/definition"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
)

// Catalog is what the stage needs to know about components.
type Catalog interface {
	Get(id string) (definition.Definition, bool)
}

type gesture int

const (
	idle gesture = iota
	moving
	resizing
	marquee
)

// Model is the stage pane.
type Model struct {
	ed   *editor.Editor
	cat  Catalog
	w, h int
	// offX, offY is the canvas cell shown at the pane's top-left corner.
	offX, offY int

	mode gesture
	// moving: where inside the primary node the pointer grabbed it.
	grabDX, grabDY int
	// resizing: the node, its starting rectangle, the handle and the press point.
	target         design.NodeID
	origin         design.Rect
	handle         hittest.Handle
	startX, startY int
	// marquee: the press point and current point, in canvas cells.
	marqueeFrom, marqueeTo [2]int
	marqueeAdditive        bool

	ghost *design.Rect
}

// New returns a stage editing through ed.
func New(ed *editor.Editor, cat Catalog) *Model {
	return &Model{ed: ed, cat: cat}
}

// SetSize sets the pane size in cells.
func (m *Model) SetSize(w, h int) {
	m.w, m.h = w, h
	m.clampOffset()
}

// Busy reports whether a pointer gesture started on the stage is in progress.
func (m *Model) Busy() bool { return m.mode != idle }

// Canvas converts pane coordinates to canvas cells.
func (m *Model) Canvas(x, y int) (int, int) { return x + m.offX, y + m.offY }

// Contains reports whether pane coordinates lie inside the pane.
func (m *Model) Contains(x, y int) bool { return x >= 0 && y >= 0 && x < m.w && y < m.h }

// SetGhost shows (or, with nil, hides) the outline of a component being
// dragged in from the palette. The rectangle is in canvas cells.
func (m *Model) SetGhost(r *design.Rect) { m.ghost = r }

func (m *Model) clampOffset() {
	doc := m.ed.Document()
	m.offX = min(max(m.offX, 0), max(doc.Width-m.w, 0))
	m.offY = min(max(m.offY, 0), max(doc.Height-m.h, 0))
}

// Handle processes a pointer event in pane coordinates.
func (m *Model) Handle(e pointer.Event) {
	cx, cy := m.Canvas(e.X, e.Y)
	switch e.Phase {
	case pointer.Wheel:
		if e.Shift {
			m.offX += e.WheelY * 3
		} else {
			m.offY += e.WheelY * 2
		}
		m.offX += e.WheelX * 3
		m.clampOffset()
	case pointer.Down:
		if e.Left {
			m.press(cx, cy, e.Shift)
		}
	case pointer.Move:
		if e.Held {
			m.drag(cx, cy)
		}
	case pointer.Up:
		m.release()
	}
}

func (m *Model) press(cx, cy int, additive bool) {
	doc := m.ed.Document()
	if sel := m.ed.Selected(); len(sel) == 1 {
		if n, ok := m.ed.Primary(); ok {
			if h := hittest.HandleAt(n.Rect, cx, cy); h != hittest.HandleNone {
				m.ed.Checkpoint()
				m.mode, m.target, m.origin, m.handle = resizing, n.ID, n.Rect, h
				m.startX, m.startY = cx, cy
				return
			}
		}
	}
	if id, ok := hittest.Node(doc, cx, cy); ok {
		switch {
		case additive:
			m.ed.Toggle(id)
		case !m.ed.IsSelected(id):
			m.ed.Select(id)
		}
		if !m.ed.IsSelected(id) {
			return
		}
		n, _ := m.ed.Primary()
		m.ed.Checkpoint()
		m.mode = moving
		m.grabDX, m.grabDY = cx-n.Rect.X, cy-n.Rect.Y
		return
	}
	if !additive {
		m.ed.Clear()
	}
	m.mode = marquee
	m.marqueeFrom, m.marqueeTo = [2]int{cx, cy}, [2]int{cx, cy}
	m.marqueeAdditive = additive
}

func (m *Model) drag(cx, cy int) {
	switch m.mode {
	case moving:
		n, ok := m.ed.Primary()
		if !ok {
			return
		}
		m.ed.MoveSelectionBy(cx-m.grabDX-n.Rect.X, cy-m.grabDY-n.Rect.Y)
	case resizing:
		minW, minH := 1, 1
		if n, ok := m.ed.Document().Get(m.target); ok {
			if def, known := m.cat.Get(n.Component); known {
				minW, minH = max(def.MinSize.W, 1), max(def.MinSize.H, 1)
			}
		}
		r := hittest.Resize(m.origin, m.handle, cx-m.startX, cy-m.startY, minW, minH)
		m.ed.SetRect(m.target, r, false)
	case marquee:
		m.marqueeTo = [2]int{cx, cy}
	}
}

func (m *Model) release() {
	switch m.mode {
	case moving, resizing:
		m.ed.EndGesture()
	case marquee:
		r := m.marqueeRect()
		ids := hittest.NodesIn(m.ed.Document(), r)
		if m.marqueeAdditive {
			ids = append(m.ed.Selected(), ids...)
		}
		m.ed.Select(ids...)
	}
	m.mode = idle
}

// marqueeRect is the rubber-band rectangle, normalised and at least one cell.
func (m *Model) marqueeRect() design.Rect {
	a, b := m.marqueeFrom, m.marqueeTo
	x, y := min(a[0], b[0]), min(a[1], b[1])
	return design.Rect{X: x, Y: y, W: max(a[0], b[0]) - x + 1, H: max(a[1], b[1]) - y + 1}
}

// Lines renders the pane as exactly h lines of w cells.
func (m *Model) Lines() []string {
	return theme.Block(m.render().Lines(), m.w, m.h)
}
