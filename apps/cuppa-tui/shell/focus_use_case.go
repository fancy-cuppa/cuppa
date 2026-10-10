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
	m.sayFocus()
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
		m.say(m.bar.Current())
	default:
		if a, ok := altAreas[text]; ok {
			m.setFocus(a)
			m.sayFocus()
			return true
		}
		if r, ok := strings.CutPrefix(text, "alt+"); ok && len(r) == 1 && m.bar.OpenMnemonic(rune(r[0])) {
			m.setFocus(inMenu)
			m.say(m.bar.Current())
			return true
		}
		return false
	}
	return true
}

// focusHints are the keys that work in the focused area.
func (m *Model) focusHints() string {
	switch m.focus {
	case inMenu:
		return "Menu bar  ·  ←/→ menu  ·  ↓/↑ item  ·  Enter run  ·  Esc close"
	case inPalette:
		return "Palette  ·  ↑/↓ choose  ·  Enter place  ·  ←/→ fold  ·  / search  ·  F6 next area"
	case inInspector:
		return "Details  ·  Tab/↑↓ move  ·  Enter edit  ·  ←/→ change  ·  h hide  ·  l lock  ·  Esc canvas"
	}
	return "Canvas  ·  arrows move  ·  Alt+arrows resize  ·  Tab layer  ·  F6 area"
}
