package confirm

import "testing"

func TestTabChoosesTheButtonEnterPresses(t *testing.T) {
	m := New("Save?", "Save changes?", "Save", "Don't save", "Cancel")
	m.Place(80, 24)
	m.Nav("tab")
	m.Nav("tab")
	m.Key("", false, true, false)
	out, done := m.Outcome()
	if !done || out.Button != "Cancel" {
		t.Fatalf("outcome %+v done=%v, want Cancel", out, done)
	}
}

func TestShiftTabWrapsToTheLastButton(t *testing.T) {
	m := New("Save?", "x", "Yes", "No")
	m.Place(80, 24)
	m.Nav("shift+tab")
	m.Key("", false, true, false)
	if out, _ := m.Outcome(); out.Button != "No" {
		t.Fatalf("button = %q, want No", out.Button)
	}
}
