package inspector

import (
	"github.com/meta-tui/cuppa/apps/cuppa-tui/pointer"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/theme"
	"github.com/meta-tui/cuppa/libs/document/design"
)

// layerDrag is a layer being dragged in the list. Pressing a layer does
// nothing yet: moving to another row starts the drag, and releasing without
// moving selects it.
type layerDrag struct {
	id     design.NodeID
	startY int
	// selectIt is what a plain click on the layer does; it runs on release.
	selectIt func()
	moved    bool
	// row is the list row the layer would land on (0 is the front).
	row       int
	overTrash bool
}

// Dragging reports whether a layer is being dragged.
func (m *Model) Dragging() bool { return m.drag != nil && m.drag.moved }

// CancelDrag drops the drag without changing anything.
func (m *Model) CancelDrag() { m.drag = nil }

// layers lists the layers front to back. Each row has an eye (show or hide) and
// a padlock (lock or unlock); the name selects it and can be dragged to a new
// place in the list, or onto the trash button.
func (m *Model) layers(b *builder) {
	doc := m.ed.Document()
	b.blank()
	b.text(theme.Bold(" Layers")).text("  ")
	m.trashButton(b)
	b.end()
	n := len(doc.Nodes)
	if n == 0 {
		b.text(theme.Faded(" (none)")).end()
		m.layerCount = 0
		return
	}
	m.layerTop, m.layerCount = len(b.lines), n
	for row := 0; row < n; row++ {
		node := doc.Nodes[n-1-row]
		m.layerRow(b, node, row)
		b.end()
	}
}

// trashButton draws the delete button: click it to delete the selected layers,
// or drop a dragged layer on it.
func (m *Model) trashButton(b *builder) {
	label := "[Trash]"
	var styled string
	switch {
	case m.drag != nil && m.drag.overTrash:
		styled = theme.Selected(label)
	case m.drag != nil && m.drag.moved:
		styled = theme.Title(label)
	default:
		styled = theme.Button(label, len(m.ed.Selected()) > 0)
	}
	m.trashY, m.trashX0 = len(b.lines), b.x
	b.add(styled, func() { m.ed.Delete() })
	m.trashX1 = b.x
}

func (m *Model) layerRow(b *builder, node design.Node, row int) {
	id := node.ID
	marker := " "
	if m.drag != nil && m.drag.moved && !m.drag.overTrash && m.drag.row == row {
		marker = theme.Title("▸")
	}
	b.text(marker)
	eye := "◉"
	if node.Hidden {
		eye = theme.Faded("○")
	}
	b.add(eye, func() { m.ed.SetHidden(id, !node.Hidden) }).text(" ")
	lock := theme.Faded("▢")
	if node.Locked {
		lock = theme.Title("▣")
	}
	b.add(lock, func() { m.ed.SetLocked(id, !node.Locked) }).text(" ")

	name := node.Name
	rest := max(m.w-b.x, 1)
	var styled string
	switch {
	case m.drag != nil && m.drag.moved && m.drag.id == id:
		styled = theme.Faded(theme.Fit(name, rest))
	case m.ed.IsSelected(id):
		styled = theme.Selected(theme.Fit(name, rest))
	case node.Hidden:
		styled = theme.Faded(name)
	default:
		styled = name
	}
	b.add(styled, func() { m.ed.Select(id) })
	b.regions[len(b.regions)-1].layer = id
}

// pressLayer starts watching for a drag of the layer under the pointer.
func (m *Model) pressLayer(id design.NodeID, y int, selectIt func()) {
	m.drag = &layerDrag{id: id, startY: y, selectIt: selectIt}
}

// dragTo follows the pointer while a layer is held.
func (m *Model) dragTo(e pointer.Event) {
	d := m.drag
	y := e.Y + m.scroll
	if y != d.startY {
		d.moved = true
	}
	d.overTrash = y == m.trashY && e.X >= m.trashX0 && e.X < m.trashX1
	if m.layerCount > 0 {
		d.row = min(max(y-m.layerTop, 0), m.layerCount-1)
	}
}

// drop ends the drag: onto the trash deletes, elsewhere moves the layer to the
// row it was dropped on. A press that never moved was just a selection.
func (m *Model) drop(e pointer.Event) {
	d := m.drag
	if d == nil {
		return
	}
	if d.moved {
		m.dragTo(e) // where the pointer was let go
	}
	m.drag = nil
	if !d.moved {
		if d.selectIt != nil {
			d.selectIt()
		}
		return
	}
	if d.overTrash {
		if m.ed.IsSelected(d.id) {
			m.ed.Delete()
		} else {
			m.ed.DeleteLayer(d.id)
		}
		return
	}
	m.ed.MoveLayer(d.id, m.layerCount-1-d.row)
}
