package shell

import "strings"

// areas is the order F6 walks through: the menu bar, the palette, the canvas
// and the details bar.
var areas = []pane{inMenu, inPalette, inStage, inInspector}

// altAreas maps Alt+1 to Alt+4 to the same areas.
var altAreas = map[string]pane{"alt+1": inMenu, "alt+2": inPalette, "alt+3": inStage, "alt+4": inInspector}

// setFocus puts the keyboard on an area and tells the panes, so each one can
// show it.
func (m *Model) setFocus(p pane) {
	if p == nowhere {
		p = inStage
	}
	m.focus = p
	m.bar.Focus(p == inMenu)
	m.pal.SetFocused(p == inPalette)
	m.ins.SetFocused(p == inInspector)
}

// cycleFocus moves the keyboard to the next area, or the previous one.
func (m *Model) cycleFocus(dir int) {
	at := 0
	for i, a := range areas {
		if a == m.focus {
			at = i
		}
	}
	m.setFocus(areas[((at+dir)%len(areas)+len(areas))%len(areas)])
}

// focusKey handles the keys that move the keyboard between areas and open the
// menus, whichever area has it. It reports whether the key was one of them.
func (m *Model) focusKey(text string) bool {
	switch text {
	case "f6":
		m.cycleFocus(1)
	case "shift+f6":
		m.cycleFocus(-1)
	case "f10":
		m.setFocus(inMenu)
		m.bar.OpenFirst()
	default:
		if a, ok := altAreas[text]; ok {
			m.setFocus(a)
			return true
		}
		if r, ok := strings.CutPrefix(text, "alt+"); ok && len(r) == 1 && m.bar.OpenMnemonic(rune(r[0])) {
			m.setFocus(inMenu)
			return true
		}
		return false
	}
	return true
}

// focusName is the area's name as the status bar and the screen reader say it.
func focusName(p pane) string {
	switch p {
	case inMenu:
		return "Menu bar"
	case inPalette:
		return "Palette"
	case inInspector:
		return "Details"
	}
	return "Canvas"
}

// focusHints are the keys that work in the focused area.
func (m *Model) focusHints() string {
	switch m.focus {
	case inMenu:
		return "Menu bar  ·  ←/→ menu  ·  ↓/↑ item  ·  Enter run  ·  Esc close"
	case inPalette:
		return "Palette  ·  F6 next area  ·  Ctrl+F search  ·  Esc canvas"
	case inInspector:
		return "Details  ·  F6 next area  ·  Esc canvas"
	}
	return "Canvas  ·  F6 next area  ·  F10 menu  ·  arrows move  ·  Del delete"
}
