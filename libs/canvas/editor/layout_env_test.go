package editor

import (
	"strings"
	"testing"
)

// A layout expression can read an input (the value of the property bound to
// it in the design) and the place of a component before it.
func TestLayoutReadsInputsAndPlacesInTheDesigner(t *testing.T) {
	e := newEditor()
	count := mustAdd(t, e, "bubbles.paginator", 0, 0)
	must(t, e.SetProp(count, "total", "7"))
	must(t, e.SetBinding(count, "total", "PlaylistCount"))
	e.Rename(count, "Counter")

	top := mustAdd(t, e, "lipgloss.box", 0, 2)
	e.Rename(top, "Playlists")
	must(t, e.SetLayout(top, AxisH, "min(max($PlaylistCount + 2, 5), (100% - 4) / 3)"))
	// 7 + 2 is 9, but a third of the 24-row canvas minus 4 is 6.
	if got := rectOf(e, top).H; got != 6 {
		t.Errorf("playlists is %d high, want 6", got)
	}

	under := mustAdd(t, e, "lipgloss.box", 0, 12)
	e.Rename(under, "Entries")
	must(t, e.SetLayout(under, AxisY, `below("Playlists")`))
	if got, want := rectOf(e, under).Y, rectOf(e, top).Y+6; got != want {
		t.Errorf("entries at %d, want %d", got, want)
	}

	// An input that changes moves the expression's result.
	must(t, e.SetProp(count, "total", "1"))
	e.Resolve()
	if got := rectOf(e, top).H; got != 5 {
		t.Errorf("playlists is %d high with one playlist, want the minimum 5", got)
	}

	// A bad expression is refused with a reason.
	if err := e.SetLayout(top, AxisH, "$"); err == nil {
		t.Error("a lone $ was accepted")
	}
}

func TestShowIfTakesAConditionAndTheVariablesListItsInputs(t *testing.T) {
	e := newEditor()
	id := mustAdd(t, e, "lipgloss.box", 0, 0)
	must(t, e.SetShowIf(id, "w >= 100 && $FullLog == 0"))
	n, _ := e.Document().Get(id)
	if n.ShowIf != "w >= 100 && $FullLog == 0" {
		t.Errorf("show-if %q", n.ShowIf)
	}
	if err := e.SetShowIf(id, "w >="); err == nil {
		t.Error("an incomplete condition was accepted")
	}
	found := false
	for _, v := range e.Variables() {
		if v.Name == "FullLog" && v.Kind == VarInput && len(v.Uses) == 1 && v.Uses[0].Role == UseShowIf {
			found = true
		}
	}
	if !found {
		t.Errorf("FullLog is not listed: %+v", e.Variables())
	}

	// Renaming the input rewrites the condition.
	must(t, e.RenameVariable(VarInput, "FullLog", "Log mode"))
	n, _ = e.Document().Get(id)
	if !strings.Contains(n.ShowIf, `$"Log mode"`) {
		t.Errorf("after the rename: %q", n.ShowIf)
	}
}
