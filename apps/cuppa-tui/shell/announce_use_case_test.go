package shell

import (
	"strconv"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func status(m *Model) string { return spoken(m.Describe().Status) }

func TestMovingFocusIsAnnounced(t *testing.T) {
	m := newShell(t)
	send(m, key(tea.KeyF6, 0))
	if got := status(m); got != "Details" {
		t.Fatalf("announced %q, want Details", got)
	}
	send(m, key('2', tea.ModAlt))
	if got := status(m); got != "Palette" {
		t.Fatalf("announced %q, want Palette", got)
	}
}

func TestSelectingByTabAnnouncesTheLayerAndItsPlace(t *testing.T) {
	m := newShell(t)
	place(t, m, 5)
	m.Editor().Clear()
	send(m, key(tea.KeyTab, 0))
	got := status(m)
	if !strings.HasPrefix(got, "Selected ") || !strings.Contains(got, " by ") || !strings.Contains(got, " at ") {
		t.Fatalf("announced %q", got)
	}
}

func TestARefusedMoveSaysWhy(t *testing.T) {
	m := newShell(t)
	n := place(t, m, 5)
	m.Editor().SetLocked(n.ID, true)
	send(m, key(tea.KeyRight, 0))
	if got := status(m); !strings.Contains(got, "Cannot move") || !strings.Contains(got, "locked") {
		t.Fatalf("announced %q, want a refusal that names the lock", got)
	}
}

func TestAMoveAnnouncesWhereItWent(t *testing.T) {
	m := newShell(t)
	n := place(t, m, 5)
	send(m, key(tea.KeyRight, 0))
	want := "Moved to " + strconv.Itoa(n.Rect.X+1) + ", " + strconv.Itoa(n.Rect.Y)
	if got := status(m); got != want {
		t.Fatalf("announced %q, want %q", got, want)
	}
}

func TestTheSameAnnouncementTwiceIsStillAChange(t *testing.T) {
	m := newShell(t)
	send(m, key('2', tea.ModAlt))
	first := m.Describe().Status
	send(m, key('2', tea.ModAlt))
	if second := m.Describe().Status; second == first {
		t.Fatal("a repeated announcement must differ so it is heard again")
	}
}
