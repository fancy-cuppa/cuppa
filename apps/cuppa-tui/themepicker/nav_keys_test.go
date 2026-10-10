package themepicker

import "testing"

func TestArrowsStepAndShiftArrowsJumpTen(t *testing.T) {
	m := New("Theme", "")
	m.Nav("right")
	if m.index != 1 {
		t.Fatalf("index = %d, want 1", m.index)
	}
	m.Nav("shift+right")
	if m.index != 11 {
		t.Fatalf("index = %d, want 11", m.index)
	}
	m.Nav("left")
	if m.index != 10 {
		t.Fatalf("index = %d, want 10", m.index)
	}
}

func TestALetterJumpsToASchemeThatStartsWithIt(t *testing.T) {
	m := New("Theme", "")
	m.Nav("d")
	if got := m.current().Name; len(got) == 0 || got[0] != 'D' && got[0] != 'd' {
		t.Fatalf("scheme = %q, want one starting with D", got)
	}
}
