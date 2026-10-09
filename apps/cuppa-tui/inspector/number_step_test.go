package inspector

import (
	"strings"
	"testing"
)

func TestUpAndDownStepANumberFieldThatIsBeingEdited(t *testing.T) {
	m, ed, id := setup(t, "lipgloss.box")
	n, _ := ed.Document().Get(id)
	m.startEdit("w", "24")
	if !m.Step(1) {
		t.Fatal("a number field should step")
	}
	got, _ := ed.Document().Get(id)
	if got.Rect.W != 25 || m.buf != "25" {
		t.Fatalf("after +1: width %d, buffer %q", got.Rect.W, m.buf)
	}
	m.Step(10)
	m.Step(-1)
	got, _ = ed.Document().Get(id)
	if got.Rect.W != 34 {
		t.Fatalf("after +10 and -1: width %d, want 34", got.Rect.W)
	}
	if !m.Editing() {
		t.Fatal("stepping should keep the field open until Enter")
	}
	m.Key("", false, true, false)
	if m.Editing() {
		t.Fatal("Enter should end the edit")
	}
	_ = n
}

func TestTextAndColourFieldsDoNotStep(t *testing.T) {
	m, _, _ := setup(t, "lipgloss.box")
	m.startEdit("name", "123")
	if m.Step(1) {
		t.Fatal("a name is text, even when it looks like a number")
	}
	m.startEdit("background", "212")
	if m.Step(1) {
		t.Fatal("a colour is not stepped")
	}
	m.startEdit("w", "abc")
	if m.Step(1) {
		t.Fatal("text in a number field does not step")
	}
	m.cancelEdit()
	if m.Step(1) {
		t.Fatal("nothing is being edited")
	}
}

func TestCanvasSizeSteps(t *testing.T) {
	m, ed, _ := setup(t, "lipgloss.box")
	m.startEdit("canvas-w", "100")
	m.Step(10)
	if w := ed.Document().Width; w != 110 {
		t.Fatalf("canvas width %d, want 110", w)
	}
}

func TestThemeSectionShowsTheFiveColoursAndAnOverrideCanGoBackToTheTheme(t *testing.T) {
	m, ed, id := setup(t, "lipgloss.box")
	m.SetSize(32, 70)
	ed.Clear()
	text := stripANSI(strings.Join(m.Lines(), "\n"))
	for _, want := range []string{"Theme", "Background", "Text", "Muted", "Border", "Secondary"} {
		if !strings.Contains(text, want) {
			t.Errorf("the theme section lacks %q:\n%s", want, text)
		}
	}
	// A border colour in the theme reaches the component, which follows it.
	if err := ed.SetThemeColor("border", "#336699"); err != nil {
		t.Fatal(err)
	}
	ed.Select(id)
	text = stripANSI(strings.Join(m.Lines(), "\n"))
	if !strings.Contains(text, "#336699") || !strings.Contains(text, "theme") {
		t.Fatalf("the colour should show the theme's value and say it follows it:\n%s", text)
	}
	// Its own colour is marked, and [theme] takes it back.
	if err := ed.SetProp(id, "color", "212"); err != nil {
		t.Fatal(err)
	}
	m.Lines()
	clickText(t, m, "212", "[theme]")
	n, _ := ed.Document().Get(id)
	if _, set := n.Props["color"]; set {
		t.Fatalf("[theme] should remove the override: %v", n.Props)
	}
}
