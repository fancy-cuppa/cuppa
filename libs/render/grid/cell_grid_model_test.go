package grid

import (
	"strings"
	"testing"

	"github.com/meta-tui/cuppa/libs/document/design"
)

func plain(lines []string) []string {
	// Strip ANSI escape sequences so tests compare characters only.
	out := make([]string, len(lines))
	for i, l := range lines {
		var b strings.Builder
		esc := false
		for _, r := range l {
			switch {
			case r == 0x1b:
				esc = true
			case esc && r == 'm':
				esc = false
			case !esc:
				b.WriteRune(r)
			}
		}
		out[i] = b.String()
	}
	return out
}

func TestBoxAndText(t *testing.T) {
	g := New(6, 3)
	g.Box(design.Rect{W: 6, H: 3}, BorderNamed("normal"), Style{Fg: "212"})
	g.Text(1, 1, "hi", Style{}, 0)
	got := plain(g.Lines())
	want := []string{"┌────┐", "│hi  │", "└────┘"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("row %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestBlitSkipsUnpaintedAndClips(t *testing.T) {
	dst := New(4, 1)
	dst.Text(0, 0, "abcd", Style{}, 0)
	src := New(3, 1)
	src.Set(1, 0, Cell{Ch: 'X'})
	dst.Blit(src, 3, 0)
	if got := plain(dst.Lines())[0]; got != "abcd" {
		t.Fatalf("blit clipped wrongly: %q", got)
	}
	dst.Blit(src, 1, 0)
	if got := plain(dst.Lines())[0]; got != "abXd" {
		t.Fatalf("blit = %q", got)
	}
}

func TestTextRespectsMaxWidthAndBounds(t *testing.T) {
	g := New(3, 1)
	if n := g.Text(0, 0, "abcdef", Style{}, 2); n != 2 {
		t.Fatalf("wrote %d", n)
	}
	g.Text(2, 0, "xyz", Style{}, 0) // clipped at the grid edge, must not panic
	if got := plain(g.Lines())[0]; got != "abx" {
		t.Fatalf("row = %q", got)
	}
}

func TestStyledOutputContainsEscapes(t *testing.T) {
	g := New(2, 1)
	g.Text(0, 0, "ab", Style{Fg: "212", Bold: true}, 0)
	if !strings.Contains(g.String(), "\x1b[") {
		t.Fatal("styled cell rendered without escape codes")
	}
}
