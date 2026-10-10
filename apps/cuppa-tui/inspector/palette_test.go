package inspector

import (
	"strings"
	"testing"
)

func TestSettingAColourNamesItAfterTheProperty(t *testing.T) {
	m, ed, id := setup(t, "lipgloss.label")
	// Typing a colour into the Foreground property names it.
	clickText(t, m, "255", "255")
	for range "255" {
		m.Key("", true, false, false)
	}
	for _, r := range "#ff0000" {
		m.Key(string(r), false, false, false)
	}
	m.Key("", false, true, false)
	n, _ := ed.Document().Get(id)
	if n.Props["color"] != "@Foreground" {
		t.Fatalf("color = %q (message %q)", n.Props["color"], m.message)
	}
	if c, _ := ed.Document().Theme.Colour("Foreground"); c != "#ff0000" {
		t.Fatalf("palette = %+v", ed.Document().Theme.Palette)
	}

	// The name shows under the property, and nothing selected lists it.
	lines := strings.Join(m.Lines(), "\n")
	if !strings.Contains(stripANSI(lines), "Foreground") {
		t.Errorf("name not shown:\n%s", lines)
	}
	ed.Select()
	if got := stripANSI(strings.Join(m.Lines(), "\n")); !strings.Contains(got, "Named colours") || !strings.Contains(got, "#ff0000") {
		t.Errorf("palette not listed:\n%s", got)
	}
}

func TestANamedColourIsAddedByTypingItsName(t *testing.T) {
	m, ed, _ := setup(t, "lipgloss.label")
	ed.Select()
	m.Lines()
	clickText(t, m, "[+ colour]", "[+ colour]")
	for _, r := range "Accent" {
		m.Key(string(r), false, false, false)
	}
	m.Key("", false, true, false)
	if _, ok := ed.Document().Theme.Colour("Accent"); !ok {
		t.Fatalf("palette = %+v (message %q)", ed.Document().Theme.Palette, m.message)
	}
}
