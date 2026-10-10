package expr

import "testing"

func TestShift(t *testing.T) {
	cases := []struct {
		src    string
		delta  int
		parent int
		want   string
	}{
		{"10", 5, 100, "15"},
		{"2 * 5", 1, 100, "11"},
		{"50%", 10, 100, "60%"},
		{"50%", 1, 120, "50.83%"},
		{"100% - 10", -3, 120, "100% - 13"},
		{"100% - 10", 12, 120, "100% + 2"},
		{"100% - 10", 10, 120, "100%"},
		{"50% + 2", 3, 120, "50% + 5"},
		{"(100% - 4) / 2", 1, 120, "(100% - 4) / 2 + 1"},
		{"min(50%, 40)", -2, 120, "min(50%, 40) - 2"},
		{"50%", 0, 100, "50%"},
	}
	for _, c := range cases {
		e, err := Parse(c.src)
		if err != nil {
			t.Fatal(err)
		}
		got := e.Shift(c.delta, c.parent)
		if got.String() != c.want {
			t.Errorf("%q shifted by %d = %q, want %q", c.src, c.delta, got.String(), c.want)
		}
		if _, err := Parse(got.String()); err != nil {
			t.Errorf("%q does not parse back: %v", got.String(), err)
		}
	}
}

func TestShiftKeepsTheResolvedValueMoving(t *testing.T) {
	for _, src := range []string{"100% - 10", "50%", "min(50%, 40)", "10"} {
		e, _ := Parse(src)
		before := e.Resolve(120)
		after := e.Shift(7, 120).Resolve(120)
		if after-before < 6 || after-before > 7 {
			t.Errorf("%q: %d -> %d", src, before, after)
		}
	}
}

func TestPercentOfResolvesBackToTheSameCells(t *testing.T) {
	for parent := 1; parent <= 200; parent += 7 {
		for cells := 0; cells <= parent; cells++ {
			if got := PercentOf(cells, parent).Resolve(parent); got != cells {
				t.Fatalf("%d of %d: %q resolves to %d", cells, parent, PercentOf(cells, parent).String(), got)
			}
		}
	}
	if PercentOf(24, 80).String() != "30%" {
		t.Errorf("24 of 80 = %q", PercentOf(24, 80).String())
	}
}
