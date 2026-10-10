package inspector

import (
	"strings"
	"testing"
)

// The Change row offers components that can take the place of the selected one
// without losing its variables, and pressing a name changes the component.
func TestChangeRowChangesTheComponentAndKeepsTheVariables(t *testing.T) {
	m, ed, id := setup(t, "bubbles.textinput")
	if err := ed.SetBinding(id, "value", "Colour"); err != nil {
		t.Fatal(err)
	}
	m.Lines()
	plain := stripANSI(strings.Join(m.Lines(), "\n"))
	if !strings.Contains(plain, "Change") || !strings.Contains(plain, "keeps its variables") {
		t.Fatalf("no Change row:\n%s", plain)
	}

	// Step through the candidates until the colour picker is the one shown.
	shown := func() string {
		for _, l := range m.Lines() {
			if p := stripANSI(l); strings.Contains(p, "Change ") && strings.Contains(p, "◂") {
				return p
			}
		}
		return ""
	}
	for i := 0; i < 30 && !strings.Contains(shown(), "Colour picker"); i++ {
		clickText(t, m, "Change ", "▸")
	}
	if !strings.Contains(shown(), "Colour picker") {
		t.Fatalf("the colour picker is not offered for a text input:\n%s", shown())
	}
	clickText(t, m, "Change ", "Colour picker")

	n, _ := ed.Document().Get(id)
	if n.Component != "lipgloss.colourpicker" || n.Bind["value"] != "Colour" {
		t.Fatalf("component %q, bind %v (message %q)", n.Component, n.Bind, m.message)
	}
	if !ed.CanUndo() {
		t.Error("the change should be undoable")
	}
	ed.Undo()
	if n, _ = ed.Document().Get(id); n.Component != "bubbles.textinput" || n.Bind["value"] != "Colour" {
		t.Errorf("undo: %q %v", n.Component, n.Bind)
	}
}

func TestChangeRowIsAbsentWhenNothingCanTakeThePlace(t *testing.T) {
	m, _, _ := setup(t, "lipgloss.box")
	if plain := stripANSI(strings.Join(m.Lines(), "\n")); strings.Contains(plain, "keeps its variables") {
		t.Errorf("a box has no compatible component, but:\n%s", plain)
	}
}
