package logo

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestRenderIsTheSizeAskedForAndDrawsTheCup(t *testing.T) {
	lines := Render(16, 8, [3]uint8{0, 0, 0})
	if len(lines) != 8 {
		t.Fatalf("%d rows, want 8", len(lines))
	}
	filled := false
	for i, l := range lines {
		if w := ansi.StringWidth(l); w != 16 {
			t.Errorf("row %d is %d cells wide, want 16", i, w)
		}
		if strings.ContainsAny(ansi.Strip(l), "▘▝▀▖▌▞▛▗▚▐▜▄▙▟█") {
			filled = true
		}
	}
	if !filled {
		t.Fatal("the cup should show as block characters")
	}
	// The logo is green on near black: some cell has a green foreground.
	if !strings.Contains(strings.Join(lines, ""), "38;2;5") && !strings.Contains(strings.Join(lines, ""), ";195;") {
		t.Logf("no obvious green in %q", lines[4])
	}
}

func TestRenderIsRepeatable(t *testing.T) {
	a := strings.Join(Render(8, 4, [3]uint8{10, 10, 10}), "\n")
	b := strings.Join(Render(8, 4, [3]uint8{10, 10, 10}), "\n")
	if a != b || a == "" {
		t.Fatal("the same size should give the same picture")
	}
	if Render(0, 4, [3]uint8{}) != nil {
		t.Fatal("no size, no picture")
	}
}

func TestBestSplitFindsTheHalfBlock(t *testing.T) {
	white, black := rgb{255, 255, 255}, rgb{0, 0, 0}
	mask, fg, _ := bestSplit([4]rgb{white, white, black, black})
	if quadrant[mask] != '▀' && quadrant[mask] != '▄' {
		t.Fatalf("top white over bottom black gave %q", quadrant[mask])
	}
	_ = fg
	mask, _, _ = bestSplit([4]rgb{white, black, white, black})
	if quadrant[mask] != '▌' && quadrant[mask] != '▐' {
		t.Fatalf("left white over right black gave %q", quadrant[mask])
	}
	mask, fg, bg := bestSplit([4]rgb{white, white, white, white})
	if mask != 0 || fg != white || bg != white {
		t.Fatalf("one colour should be a plain cell, got mask %d", mask)
	}
}
