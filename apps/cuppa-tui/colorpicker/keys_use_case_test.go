package colorpicker

import "testing"

func TestTabChangesTheTabAndArrowsMoveThroughSwatches(t *testing.T) {
	m := New("Colour", "")
	m.Nav("right") // nothing chosen yet: starts at the first swatch
	if m.Value() != "0" {
		t.Fatalf("value = %q, want 0", m.Value())
	}
	m.Nav("right")
	m.Nav("down")
	if m.Value() != "9" {
		t.Fatalf("value = %q, want 9 (one right, one row down)", m.Value())
	}
	m.Nav("tab")
	if m.tab != tab256 {
		t.Fatalf("tab = %v, want the 256 tab", m.tab)
	}
	m.Nav("down")
	if m.Value() != "15" {
		t.Fatalf("value = %q, want 15 (a row of the cube is six)", m.Value())
	}
}

func TestSlidersAreChosenWithUpDownAndChangedWithLeftRight(t *testing.T) {
	m := New("Colour", "#000000")
	m.Nav("shift+right") // R +10
	m.Nav("down")        // G
	m.Nav("right")
	if got := m.Value(); got != "#0a0100" {
		t.Fatalf("value = %q, want #0a0100", got)
	}
	m.Nav("left")
	m.Nav("left") // clamped at 0
	if got := m.Value(); got != "#0a0000" {
		t.Fatalf("value = %q, want #0a0000", got)
	}
}
