package shell

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestEnterOnTheCanvasFocusesTheDetailsAndRightStepsANumber(t *testing.T) {
	m := newShell(t)
	n := place(t, m, 5)
	send(m, key(tea.KeyEnter, 0))
	if m.focus != inInspector {
		t.Fatalf("focus = %v", m.focus)
	}
	m.View() // lay the bar out so its stops are known
	// Stops: Undo, Redo, (name is plain text), X, Y, W, H ... walk to X.
	for i := 0; i < 8; i++ {
		m.render()
		got, _ := m.Editor().Primary()
		before := got.Rect.X
		send(m, key(tea.KeyRight, 0))
		m.render()
		after, _ := m.Editor().Primary()
		if after.Rect.X == before+1 {
			send(m, key(tea.KeyLeft, tea.ModShift))
			back, _ := m.Editor().Primary()
			if back.Rect.X != n.Rect.X+1-10 && back.Rect.X != 0 {
				t.Fatalf("Shift+Left moved x to %d", back.Rect.X)
			}
			return
		}
	}
	t.Fatal("no stop changed the x position with Right")
}

func TestTabWalksTheDetailsAndEnterPressesAButton(t *testing.T) {
	m := newShell(t)
	place(t, m, 5)
	send(m, key('4', tea.ModAlt))
	m.render()
	// First stop is [Undo]: Enter undoes the placement.
	send(m, key(tea.KeyEnter, 0))
	if n := len(m.Editor().Document().Nodes); n != 0 {
		t.Fatalf("Enter on [Undo] left %d components", n)
	}
}

func TestLayerKeysHideLockAndDelete(t *testing.T) {
	m := newShell(t)
	n := place(t, m, 5)
	send(m, key('4', tea.ModAlt))
	m.render()
	send(m, key(tea.KeyTab, tea.ModShift)) // wraps to the last stop: the layer
	m.render()
	send(m, typed('h'))
	got, _ := m.Editor().Document().Get(n.ID)
	if !got.Hidden {
		t.Fatal("h should hide the layer")
	}
	send(m, typed('l'))
	got, _ = m.Editor().Document().Get(n.ID)
	if !got.Locked {
		t.Fatal("l should lock the layer")
	}
	send(m, typed('l')) // a locked layer refuses deletion, as it does for the mouse
	send(m, key(tea.KeyDelete, 0))
	if _, ok := m.Editor().Document().Get(n.ID); ok {
		t.Fatal("Delete should remove the layer")
	}
}
