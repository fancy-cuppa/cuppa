package snap

import "testing"

func TestLockAngle(t *testing.T) {
	cases := []struct{ dx, dy, wx, wy int }{
		{10, 1, 10, 0},   // nearly horizontal
		{-10, 2, -10, 0}, // still horizontal
		{1, 9, 0, 9},     // nearly vertical
		{10, 5, 10, 5},   // a diagonal as it looks: two columns per row
		{-10, -5, -10, -5},
		{10, -5, 10, -5},
		{0, 0, 0, 0},
	}
	for _, c := range cases {
		if x, y := LockAngle(c.dx, c.dy); x != c.wx || y != c.wy {
			t.Errorf("LockAngle(%d, %d) = (%d, %d), want (%d, %d)", c.dx, c.dy, x, y, c.wx, c.wy)
		}
	}
}
