package shell

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// keyboardPlace puts the first component of the palette on the canvas without
// the pointer.
func keyboardPlace(t *testing.T, m *Model) {
	t.Helper()
	send(m, key('2', tea.ModAlt))
	send(m, key(tea.KeyDown, 0)) // first component under the first pack
	send(m, key(tea.KeyEnter, 0))
}

func TestEnterInThePaletteLeavesTheComponentAtTheCanvasCursor(t *testing.T) {
	m := newShell(t)
	m.curX, m.curY = 7, 3
	keyboardPlace(t, m)
	n, ok := m.Editor().Primary()
	if !ok {
		t.Fatal("nothing was placed")
	}
	if n.Rect.X != 7 || n.Rect.Y != 3 {
		t.Fatalf("placed at (%d,%d), want (7,3)", n.Rect.X, n.Rect.Y)
	}
	if m.curY <= n.Rect.Bottom()-1 {
		t.Fatalf("cursor stayed on the component: y = %d", m.curY)
	}
}

func TestEnterOnAPackHeaderFoldsItAndPlacesNothing(t *testing.T) {
	m := newShell(t)
	send(m, key('2', tea.ModAlt))
	send(m, key(tea.KeyEnter, 0)) // the cursor starts on the first pack header
	if n := len(m.Editor().Document().Nodes); n != 0 {
		t.Fatalf("Enter on a pack header placed %d components", n)
	}
	send(m, key(tea.KeyDown, 0))
	send(m, key(tea.KeyEnter, 0)) // the next row is the second pack's header now
	if n := len(m.Editor().Document().Nodes); n != 0 {
		t.Fatalf("a folded pack left a component to place: %d", n)
	}
}

func TestSlashInThePaletteStartsTheSearch(t *testing.T) {
	m := newShell(t)
	send(m, key('2', tea.ModAlt))
	send(m, typed('/'))
	if !m.pal.Searching() {
		t.Fatal("/ should start the search")
	}
}

func TestArrowsMoveTheCanvasCursorWhenNothingIsSelected(t *testing.T) {
	m := newShell(t)
	x, y := m.curX, m.curY
	send(m, key(tea.KeyRight, 0))
	send(m, key(tea.KeyDown, tea.ModShift))
	if m.curX != x+1 || m.curY != y+10 {
		t.Fatalf("cursor (%d,%d), want (%d,%d)", m.curX, m.curY, x+1, y+10)
	}
	send(m, key(tea.KeyLeft, tea.ModShift))
	if m.curX != 0 {
		t.Fatalf("cursor left the canvas: x = %d", m.curX)
	}
}

func TestAltArrowsResizeTheSelection(t *testing.T) {
	m := newShell(t)
	n := place(t, m, 5)
	send(m, key(tea.KeyRight, tea.ModAlt))
	send(m, key(tea.KeyDown, tea.ModAlt|tea.ModShift))
	got, _ := m.Editor().Primary()
	if got.Rect.W != n.Rect.W+1 || got.Rect.H != n.Rect.H+10 {
		t.Fatalf("size %dx%d, want %dx%d", got.Rect.W, got.Rect.H, n.Rect.W+1, n.Rect.H+10)
	}
	if got.Rect.X != n.Rect.X || got.Rect.Y != n.Rect.Y {
		t.Fatal("resizing moved the component")
	}
	send(m, ctrl('z'))
	send(m, ctrl('z'))
	got, _ = m.Editor().Primary()
	if got.Rect != n.Rect {
		t.Fatalf("two undos give %+v, want %+v", got.Rect, n.Rect)
	}
}

func TestTabWalksTheLayersAndSpaceExtendsTheSelection(t *testing.T) {
	m := newShell(t)
	a := place(t, m, 5)
	b := place(t, m, 40)
	place(t, m, 70)
	m.Editor().Clear()
	send(m, key(tea.KeyTab, 0))
	if p, _ := m.Editor().Primary(); p.ID != a.ID {
		t.Fatalf("first Tab selected %s, want the back layer", p.ID)
	}
	send(m, key(tea.KeyTab, 0))
	if p, _ := m.Editor().Primary(); p.ID != b.ID {
		t.Fatalf("second Tab selected %s", p.ID)
	}
	send(m, key(tea.KeyTab, tea.ModShift))
	if p, _ := m.Editor().Primary(); p.ID != a.ID {
		t.Fatalf("Shift+Tab went to %s, want the back layer", p.ID)
	}
	send(m, typed(' '))
	send(m, typed(' '))
	if n := len(m.Editor().Selected()); n != 3 {
		t.Fatalf("two Space presses give a selection of %d, want 3", n)
	}
}

func TestCtrlASelectsEverythingAndEnterFocusesTheDetails(t *testing.T) {
	m := newShell(t)
	place(t, m, 5)
	place(t, m, 40)
	send(m, ctrl('a'))
	if n := len(m.Editor().Selected()); n != 2 {
		t.Fatalf("Ctrl+A selected %d, want 2", n)
	}
	send(m, key(tea.KeyEnter, 0))
	if m.focus != inInspector {
		t.Fatalf("Enter focused %v, want the details bar", m.focus)
	}
}
