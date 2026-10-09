package inspector

import "testing"

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
