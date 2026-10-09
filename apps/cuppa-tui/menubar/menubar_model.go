// Package menubar is the File / Edit / Export / Help bar along the top and the
// dropdowns under it. It only reports which action was chosen.
package menubar

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/pointer"
	"github.com/fancy-cuppa/cuppa/apps/cuppa-tui/theme"
)

const (
	brand    = " ☕ Cuppa "
	gap      = 1
	itemPadX = 2
)

// Model is the bar. Row 0 of its own coordinates is the bar, the dropdown
// hangs below from row 1.
type Model struct {
	w int
	// open is the index of the open menu, or -1.
	open int
	// hover is the dropdown item under the pointer, or -1.
	hover int
	// hoverBar is the bar label under the pointer, or -1.
	hoverBar int
	disabled map[Action]bool
	labelX   []int // left column of each label
}

// New returns a closed bar.
func New() *Model {
	return &Model{open: -1, hover: -1, hoverBar: -1, disabled: map[Action]bool{}}
}

// SetWidth sets the screen width.
func (m *Model) SetWidth(w int) { m.w = w }

// SetEnabled greys out (or restores) one action.
func (m *Model) SetEnabled(a Action, enabled bool) { m.disabled[a] = !enabled }

// Open reports whether a dropdown is showing.
func (m *Model) Open() bool { return m.open >= 0 }

// Close hides the dropdown.
func (m *Model) Close() { m.open, m.hover = -1, -1 }

// layout computes where each label sits.
func (m *Model) layout() {
	m.labelX = m.labelX[:0]
	x := ansi.StringWidth(brand) + 2
	for _, mn := range menus {
		m.labelX = append(m.labelX, x)
		x += ansi.StringWidth(mn.label) + 2*gap + 1
	}
}

// Line renders the bar, with right-aligned status text.
func (m *Model) Line(right string) string {
	m.layout()
	var b strings.Builder
	b.WriteString(theme.Title(brand))
	b.WriteString(theme.Faded("│ "))
	for i, mn := range menus {
		label := " " + mn.label + " "
		switch i {
		case m.open:
			label = theme.Selected(label)
		case m.hoverBar:
			label = theme.Hovered(label)
		default:
			label = theme.Dim(label)
		}
		b.WriteString(label + " ")
	}
	left := b.String()
	pad := max(m.w-ansi.StringWidth(left)-ansi.StringWidth(right)-1, 1)
	return theme.Fit(left+strings.Repeat(" ", pad)+right, m.w)
}

// dropdown geometry --------------------------------------------------------

func (m *Model) dropdownWidth() int {
	w := 0
	for _, it := range menus[m.open].items {
		n := ansi.StringWidth(it.label)
		if it.shortcut != "" {
			n += 3 + ansi.StringWidth(it.shortcut)
		}
		w = max(w, n)
	}
	return w + 2*itemPadX
}

// Dropdown returns the open menu's lines and the column to draw them at.
// The first line belongs on screen row 1.
func (m *Model) Dropdown() (x int, lines []string) {
	if m.open < 0 {
		return 0, nil
	}
	m.layout()
	w := m.dropdownWidth()
	x = min(m.labelX[m.open], max(m.w-w-2, 0))
	items := menus[m.open].items
	body := make([]string, 0, len(items))
	for i, it := range items {
		body = append(body, m.itemLine(i, it, w))
	}
	return x, theme.Panel("", body, w+2)
}

func (m *Model) itemLine(i int, it item, w int) string {
	if it.label == separatorLabel {
		return theme.Faded(strings.Repeat("─", w))
	}
	gapW := max(w-2*itemPadX-ansi.StringWidth(it.label)-ansi.StringWidth(it.shortcut), 1)
	text := strings.Repeat(" ", itemPadX) + it.label + strings.Repeat(" ", gapW) + it.shortcut + strings.Repeat(" ", itemPadX)
	switch {
	case m.disabled[it.action]:
		return theme.Faded(text)
	case i == m.hover:
		return theme.Selected(text)
	}
	return text
}

// Handle takes a pointer event in screen coordinates. It returns the chosen
// action (or "") and whether the bar consumed the event. While a dropdown is
// open it consumes every press so a click elsewhere only closes it.
func (m *Model) Handle(e pointer.Event) (Action, bool) {
	m.layout()
	if e.Phase == pointer.Wheel || (e.Phase == pointer.Up && !m.Open()) {
		return nothing, false
	}
	onBar := e.Y == 0
	barIdx := m.labelAt(e.X)
	if onBar {
		m.hoverBar = barIdx
	} else {
		m.hoverBar = -1
	}
	itemIdx := m.itemAt(e.X, e.Y)
	switch e.Phase {
	case pointer.Move:
		if m.Open() {
			if onBar && barIdx >= 0 {
				m.open, m.hover = barIdx, -1
			}
			m.hover = itemIdx
		}
		return nothing, onBar && barIdx >= 0 || m.Open() && itemIdx >= 0
	case pointer.Down:
		return m.press(e, onBar, barIdx, itemIdx)
	case pointer.Up:
		return nothing, true
	}
	return nothing, false
}

func (m *Model) press(e pointer.Event, onBar bool, barIdx, itemIdx int) (Action, bool) {
	if !e.Left {
		return nothing, false
	}
	if onBar && barIdx >= 0 {
		if m.open == barIdx {
			m.Close()
		} else {
			m.open, m.hover = barIdx, -1
		}
		return nothing, true
	}
	if !m.Open() {
		return nothing, false
	}
	if itemIdx >= 0 {
		it := menus[m.open].items[itemIdx]
		if it.action == nothing || m.disabled[it.action] {
			return nothing, true
		}
		m.Close()
		return it.action, true
	}
	m.Close()
	return nothing, true
}

// labelAt is the bar label under column x, or -1.
func (m *Model) labelAt(x int) int {
	for i, mn := range menus {
		if x >= m.labelX[i] && x < m.labelX[i]+ansi.StringWidth(mn.label)+2*gap {
			return i
		}
	}
	return -1
}

// itemAt is the dropdown item under a screen cell, or -1.
func (m *Model) itemAt(x, y int) int {
	if m.open < 0 {
		return -1
	}
	w := m.dropdownWidth()
	left := min(m.labelX[m.open], max(m.w-w-2, 0))
	i := y - 2 // row 0 bar, row 1 box top border
	if x <= left || x >= left+w+1 || i < 0 || i >= len(menus[m.open].items) {
		return -1
	}
	return i
}
