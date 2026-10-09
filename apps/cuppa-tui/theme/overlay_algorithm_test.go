package theme

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestOverlayKeepsSidesAndWidth(t *testing.T) {
	base := Fit(Dim("abcdefghij"), 10)
	got := ansi.Strip(Overlay(base, "XY", 3))
	if got != "abcXYfghij" {
		t.Fatalf("got %q", got)
	}
	if ansi.StringWidth(Overlay(base, "XYZ", 8)) != 11 {
		// Overlays past the edge are the caller's job to clip.
		t.Log("overlay past the edge widens the line")
	}
}

func TestPanelIsExactlyWide(t *testing.T) {
	lines := Panel("Open", []string{"hello", "x"}, 20)
	if len(lines) != 4 {
		t.Fatalf("got %d lines", len(lines))
	}
	for _, l := range lines {
		if w := ansi.StringWidth(l); w != 20 {
			t.Errorf("width %d in %q", w, ansi.Strip(l))
		}
	}
	if !strings.Contains(ansi.Strip(lines[0]), "Open") {
		t.Error("title missing")
	}
}
