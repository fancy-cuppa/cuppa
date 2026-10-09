package confirm

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/meta-tui/cuppa/apps/cuppa-tui/pointer"
)

func clickButton(t *testing.T, m *Model, label string) {
	t.Helper()
	lines := m.Lines()
	for y, l := range lines {
		plain := ansi.Strip(l)
		if i := strings.Index(plain, "[ "+label+" ]"); i >= 0 {
			m.Handle(pointer.Event{X: m.rect.X + len([]rune(plain[:i])) + 2, Y: m.rect.Y + y, Phase: pointer.Down, Left: true})
			return
		}
	}
	t.Fatalf("no %q button in:\n%s", label, strings.Join(lines, "\n"))
}

func TestClickingAButtonEndsTheDialog(t *testing.T) {
	m := New("Unsaved changes", "Save changes to demo.cuppa?", "Save", "Discard", "Cancel")
	m.Place(100, 30)
	clickButton(t, m, "Discard")
	out, done := m.Outcome()
	if !done || out.Button != "Discard" || out.Canceled {
		t.Fatalf("got %+v done=%v", out, done)
	}
}

func TestEnterPicksDefaultAndEscCancels(t *testing.T) {
	m := New("t", "b", "Save", "Cancel")
	m.Place(80, 24)
	m.Key("", false, true, false)
	if out, _ := m.Outcome(); out.Button != "Save" {
		t.Fatalf("enter: %+v", out)
	}
	m = New("t", "b", "Save", "Cancel")
	m.Place(80, 24)
	m.Key("", false, false, true)
	if out, done := m.Outcome(); !done || !out.Canceled {
		t.Fatalf("esc: %+v", out)
	}
}

func TestRenderFitsItsRect(t *testing.T) {
	m := New("Error", strings.Repeat("long message ", 12), "OK")
	m.Place(70, 20)
	lines := m.Lines()
	if len(lines) != m.Rect().H {
		t.Fatalf("lines %d, rect %d", len(lines), m.Rect().H)
	}
	for _, l := range lines {
		if ansi.StringWidth(l) != m.Rect().W {
			t.Fatalf("width %d != %d: %q", ansi.StringWidth(l), m.Rect().W, ansi.Strip(l))
		}
	}
}

func TestClicksOutsideButtonsDoNothing(t *testing.T) {
	m := New("t", "body", "OK")
	m.Place(80, 24)
	m.Lines()
	m.Handle(pointer.Event{X: 0, Y: 0, Phase: pointer.Down, Left: true})
	if _, done := m.Outcome(); done {
		t.Fatal("ended by a stray click")
	}
}
