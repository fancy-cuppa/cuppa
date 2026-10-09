// Package packsdialog is the Packs box: every component pack with a checkbox
// that switches it on or off in the palette.
package packsdialog

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/modal"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/pointer"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/theme"
	"github.com/meta-tui/cuppa/libs/catalog/definition"
	"github.com/meta-tui/cuppa/libs/document/design"
)

const width = 56

// Entry is one pack with how many components it holds.
type Entry struct {
	Pack  definition.Pack
	Count int
}

// Model is the open dialog. Switching a pack calls toggle at once, so the
// palette behind it follows; the box only ends with Done.
type Model struct {
	entries []Entry
	enabled func(definition.Family) bool
	toggle  func(definition.Family)

	rect    design.Rect
	hoverY  int
	doneHot bool
	outcome modal.Outcome
	done    bool
}

// New returns the dialog. enabled reads the current state of a pack.
func New(entries []Entry, enabled func(definition.Family) bool, toggle func(definition.Family)) *Model {
	return &Model{entries: entries, enabled: enabled, toggle: toggle, hoverY: -1}
}

// Place implements modal.Modal.
func (m *Model) Place(sw, sh int) {
	w := min(width, sw)
	h := len(m.entries)*3 + 5
	m.rect = design.Rect{X: max((sw-w)/2, 0), Y: max((sh-h)/3, 0), W: w, H: h}
}

// Rect implements modal.Modal.
func (m *Model) Rect() design.Rect { return m.rect }

// rowTop is the y of an entry's first line, inside the box.
func rowTop(i int) int { return 2 + i*3 }

// Lines implements modal.Modal.
func (m *Model) Lines() []string {
	inner := m.rect.W - 4
	content := []string{""}
	for i, e := range m.entries {
		mark := "[ ]"
		if m.enabled(e.Pack.ID) {
			mark = "[x]"
		}
		head := fmt.Sprintf(" %s %s", mark, e.Pack.Name)
		count := theme.Faded(fmt.Sprintf("%d components", e.Count))
		head = theme.Fit(head, max(inner-ansi.StringWidth(count)-1, 1))
		head += " " + count
		if m.hoverY >= rowTop(i) && m.hoverY < rowTop(i)+2 {
			head = theme.Selected(ansi.Strip(head))
		}
		content = append(content, " "+head, "      "+theme.Faded(theme.Fit(e.Pack.Description, inner-6)), "")
	}
	button := "[ Done ]"
	if m.doneHot {
		button = theme.Selected(button)
	} else {
		button = theme.Button(button, true)
	}
	pad := max(m.rect.W-2-ansi.StringWidth(button)-1, 0)
	content = append(content, spaces(pad)+button)
	return theme.Panel("Component packs", content, m.rect.W)
}

func spaces(n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = ' '
	}
	return string(out)
}

// Handle implements modal.Modal.
func (m *Model) Handle(e pointer.Event) {
	if m.done {
		return
	}
	x, y := e.X-m.rect.X, e.Y-m.rect.Y
	m.hoverY, m.doneHot = -1, false
	if i, ok := m.entryAt(y); ok && x > 0 && x < m.rect.W-1 {
		m.hoverY = rowTop(i)
		if e.Phase == pointer.Down && e.Left {
			m.toggle(m.entries[i].Pack.ID)
		}
		return
	}
	if y == m.rect.H-2 && x >= m.rect.W-12 && x < m.rect.W-1 {
		m.doneHot = true
		if e.Phase == pointer.Down && e.Left {
			m.finish(false)
		}
	}
}

func (m *Model) entryAt(y int) (int, bool) {
	for i := range m.entries {
		if y >= rowTop(i) && y < rowTop(i)+2 {
			return i, true
		}
	}
	return 0, false
}

// Key implements modal.Modal.
func (m *Model) Key(_ string, _, enter, esc bool) {
	if enter || esc {
		m.finish(esc)
	}
}

func (m *Model) finish(canceled bool) {
	m.done = true
	m.outcome = modal.Outcome{Button: "Done", Canceled: canceled}
}

// Outcome implements modal.Modal.
func (m *Model) Outcome() (modal.Outcome, bool) { return m.outcome, m.done }
