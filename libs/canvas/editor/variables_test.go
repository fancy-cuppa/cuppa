package editor

import "testing"

func TestVariablesListEveryNameAndWhereItIsUsed(t *testing.T) {
	e := newEditor()
	a := mustAdd(t, e, "lipgloss.label", 0, 0)
	b := mustAdd(t, e, "lipgloss.label", 0, 2)
	must(t, e.SetColour(a, "color", "#ff0000"))
	must(t, e.UseSwatch(b, "color", "Foreground"))
	must(t, e.SetBinding(a, "text", "Title"))
	must(t, e.SetShowIf(b, "Show b"))
	must(t, e.SetEvent(b, "Open"))
	must(t, e.SetKeys("esc=Open:back"))

	byName := map[string]Variable{}
	for _, v := range e.Variables() {
		byName[string(v.Kind)+":"+v.Name] = v
	}
	if v := byName["colour:Foreground"]; len(v.Uses) != 2 || v.Value != "#ff0000" {
		t.Errorf("colour = %+v", v)
	}
	if v := byName["input:Title"]; len(v.Uses) != 1 || v.Type != "text" || v.Uses[0].Where != "Text" {
		t.Errorf("input = %+v", v)
	}
	if v := byName["input:Show b"]; v.Type != "yes/no" || v.Uses[0].Role != UseShowIf {
		t.Errorf("show-if = %+v", v)
	}
	if v := byName["event:Open"]; len(v.Uses) != 2 {
		t.Errorf("event = %+v", v)
	}
}

func TestRenameVariableChangesEveryUse(t *testing.T) {
	e := newEditor()
	a := mustAdd(t, e, "lipgloss.label", 0, 0)
	b := mustAdd(t, e, "lipgloss.label", 0, 2)
	must(t, e.SetColour(a, "color", "#ff0000"))
	must(t, e.UseSwatch(b, "color", "Foreground"))
	must(t, e.SetEvent(a, "Open"))
	must(t, e.SetKeys("o=Open"))

	must(t, e.RenameVariable(VarColour, "Foreground", "Accent"))
	na, _ := e.Document().Get(a)
	nb, _ := e.Document().Get(b)
	if na.Props["color"] != "@Accent" || nb.Props["color"] != "@Accent" {
		t.Fatalf("uses = %q %q", na.Props["color"], nb.Props["color"])
	}
	if c, _ := e.Document().Theme.Colour("Accent"); c != "#ff0000" {
		t.Fatal("palette not renamed")
	}
	must(t, e.RenameVariable(VarEvent, "Open", "Show"))
	na, _ = e.Document().Get(a)
	if na.Event != "Show" || e.Document().Keys[0].Event != "Show" {
		t.Errorf("event = %q keys = %+v", na.Event, e.Document().Keys)
	}
	must(t, e.SetSwatch("Other", "#00ff00"))
	if err := e.RenameVariable(VarColour, "Accent", "Other"); err == nil {
		t.Error("renamed onto an existing name")
	}
	if err := e.RenameVariable(VarColour, "Nope", "Fresh"); err == nil {
		t.Error("renamed a variable that does not exist")
	}
	e.Undo()
	e.Undo()
	na, _ = e.Document().Get(a)
	if na.Event != "Open" {
		t.Errorf("undo left event %q", na.Event)
	}
}

func TestRelinkMovesOneUse(t *testing.T) {
	e := newEditor()
	a := mustAdd(t, e, "lipgloss.label", 0, 0)
	b := mustAdd(t, e, "lipgloss.label", 0, 2)
	must(t, e.SetColour(a, "color", "#ff0000"))
	must(t, e.UseSwatch(b, "color", "Foreground"))
	must(t, e.SetSwatch("Dim", "#666666"))

	var use Use
	for _, v := range e.Variables() {
		if v.Name == "Foreground" {
			for _, u := range v.Uses {
				if u.Node == b {
					use = u
				}
			}
		}
	}
	must(t, e.Relink(use, VarColour, "Dim"))
	nb, _ := e.Document().Get(b)
	na, _ := e.Document().Get(a)
	if nb.Props["color"] != "@Dim" || na.Props["color"] != "@Foreground" {
		t.Fatalf("a %q b %q", na.Props["color"], nb.Props["color"])
	}
	if err := e.Relink(use, VarColour, "Missing"); err == nil {
		t.Error("linked to a colour that does not exist")
	}
}
