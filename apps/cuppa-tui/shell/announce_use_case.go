package shell

import (
	"fmt"
	"strings"
)

// say sets what a screen reader announces next. The same words twice in a row
// would not be heard again, so a zero-width mark alternates at the end.
func (m *Model) say(text string) {
	m.announceN++
	if m.announceN%2 == 0 {
		text += "​"
	}
	m.announce = text
}

// settleAnnouncement lets the file flow's own feedback ("Saved x") take over
// when it changes.
func (m *Model) settleAnnouncement() {
	if st := m.flow.Status(); st != m.lastFlowStatus {
		m.lastFlowStatus = st
		if st != "" {
			m.announce = ""
		}
	}
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

// sayFocus announces where the keyboard went.
func (m *Model) sayFocus() { m.say(focusName(m.focus)) }

// saySelection announces what is selected: the layer's name and where it is,
// or how many layers.
func (m *Model) saySelection() {
	sel := m.ed.Selected()
	switch len(sel) {
	case 0:
		m.say("Nothing selected")
	case 1:
		if n, ok := m.ed.Primary(); ok {
			state := ""
			switch {
			case n.Locked:
				state = ", locked"
			case n.Hidden:
				state = ", hidden"
			}
			m.say(fmt.Sprintf("Selected %s, %d by %d at %d, %d%s", n.Name, n.Rect.W, n.Rect.H, n.Rect.X, n.Rect.Y, state))
		}
	default:
		m.say(fmt.Sprintf("%d layers selected", len(sel)))
	}
}

// sayRefusal announces why a move or a resize did nothing.
func (m *Model) sayRefusal(verb string) {
	doc := m.ed.Document()
	for _, id := range m.ed.Selected() {
		if n, ok := doc.Get(id); ok && n.Locked {
			m.say("Cannot " + verb + ", " + n.Name + " is locked")
			return
		}
	}
	m.say("Cannot " + verb + ", at the edge or the smallest size")
}

// sayCursor announces the canvas cursor.
func (m *Model) sayCursor() { m.say(fmt.Sprintf("Cursor at %d, %d", m.curX, m.curY)) }

// spoken trims an announcement for tests.
func spoken(s string) string { return strings.TrimSuffix(s, "​") }
