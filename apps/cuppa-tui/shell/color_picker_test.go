package shell

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func typeKeys(m *Model, s string) {
	for _, r := range s {
		send(m, tea.KeyPressMsg{Code: r, Text: string(r)})
	}
}

// boxWithBorderColourOnScreen drops a Box and returns the model, ready to click its colour.
func boxWithBorderColourOnScreen(t *testing.T) *Model {
	t.Helper()
	m := newShell(t)
	send(m, click(5, firstComponentY)) // first palette entry: Box
	send(m, release(m.layout.stage.X+10, 8))
	if _, ok := m.Editor().Primary(); !ok {
		t.Fatal("the dropped box should be selected")
	}
	return m
}

func TestClickingAColourOpensThePickerAndATypedValueIsApplied(t *testing.T) {
	m := boxWithBorderColourOnScreen(t)
	clickText(t, m, "212")
	if m.flow.Modal() == nil {
		t.Fatal("clicking a colour should open the colour dialog")
	}
	screen := strings.Join(screenText(m), "\n")
	for _, want := range []string{"Border color", "RGB", "HSL", "[ OK ]", "[ Cancel ]"} {
		if !strings.Contains(screen, want) {
			t.Fatalf("dialog lacks %q:\n%s", want, screen)
		}
	}
	typeKeys(m, "#ff8800")
	send(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.flow.Modal() != nil {
		t.Fatal("Enter on a valid value should close the dialog")
	}
	n, _ := m.Editor().Primary()
	// A colour set in the details bar is a named colour, named after the
	// property, so it can be reused and a program can change it.
	if n.Props["color"] != "@Border color" {
		t.Fatalf("colour = %q", n.Props["color"])
	}
	if c, _ := m.Editor().Document().Theme.Colour("Border color"); c != "#ff8800" {
		t.Fatalf("named colour = %q", c)
	}
	if !m.Editor().CanUndo() {
		t.Fatal("changing a colour is an undo step")
	}
}

func TestCancellingThePickerChangesNothing(t *testing.T) {
	m := boxWithBorderColourOnScreen(t)
	clickText(t, m, "212")
	typeKeys(m, "#00ff00")
	send(m, tea.KeyPressMsg{Code: tea.KeyEscape})
	if m.flow.Modal() != nil {
		t.Fatal("Esc closes the dialog")
	}
	if n, _ := m.Editor().Primary(); n.Props["color"] != "" {
		t.Fatalf("cancel must leave the colour alone, got %q", n.Props["color"])
	}
}

func TestPickingAPaletteSwatchByMouse(t *testing.T) {
	m := boxWithBorderColourOnScreen(t)
	clickText(t, m, "212")
	dlg := m.flow.Modal()
	if dlg == nil {
		t.Fatal("no dialog")
	}
	// The dialog opens on the 256 tab (the colour is 212). Click inside it
	// somewhere on the swatch area and then OK: some palette colour is chosen.
	screenText(m) // a frame is drawn before the next click, which is when the dialog learns where its swatches are
	r := dlg.Rect()
	send(m, click(r.X+4, r.Y+bodyRowForTest))
	send(m, release(r.X+4, r.Y+bodyRowForTest))
	clickText(t, m, "[ OK ]")
	if m.flow.Modal() != nil {
		t.Fatal("OK closes the dialog")
	}
	n, _ := m.Editor().Primary()
	if n.Props["color"] == "" || n.Props["color"] == "212" {
		t.Fatalf("clicking a swatch should have chosen a different colour, got %q", n.Props["color"])
	}
}

// bodyRowForTest is the first swatch row of the colour dialog, box-relative.
const bodyRowForTest = 3
