package shell

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/menubar"
)

func TestAboutShowsTheLogoAboveTheText(t *testing.T) {
	m := newShell(t)
	m.perform(menubar.HelpAbout)
	dlg := m.flow.Modal()
	if dlg == nil {
		t.Fatal("About should open a dialog")
	}
	lines := dlg.Lines()
	if len(lines) != dlg.Rect().H {
		t.Fatalf("%d lines in a rect %d tall", len(lines), dlg.Rect().H)
	}
	plain := ansi.Strip(strings.Join(lines, "\n"))
	if !strings.ContainsAny(plain, "▘▝▀▖▌▞▛▗▚▐▜▄▙▟█") {
		t.Fatalf("the logo's blocks are missing:\n%s", plain)
	}
	if !strings.Contains(plain, "A designer for Bubble Tea interfaces") {
		t.Fatal("the text should still be there")
	}
	for i, l := range lines {
		if w := ansi.StringWidth(l); w != dlg.Rect().W {
			t.Errorf("line %d is %d wide, the box is %d", i, w, dlg.Rect().W)
		}
	}
}
