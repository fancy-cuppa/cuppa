package gosource

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/meta-tui/cuppa/libs/canvas/editor"
	"github.com/meta-tui/cuppa/libs/catalog/standard"
	"github.com/meta-tui/cuppa/libs/cuppafile/disk"
	"github.com/meta-tui/cuppa/libs/document/design"
)

const mvdColoursExample = "../../../examples/mvd-colours.cuppa"

// mvdColours builds, through the editor, the design of MVD's Colours screen
// the way it would be recreated in Cuppa: a frame, the list of colour slots,
// a help text, a status line that comes and goes and a key bar.
func mvdColours(t *testing.T) design.Document {
	t.Helper()
	ed := editor.New(standard.Default(), design.NewDocument("MVD Colours", 80, 24))
	add := func(comp, name string, x, y, w, h int) design.NodeID {
		id, err := ed.Add(comp, x, y)
		if err != nil {
			t.Fatal(err)
		}
		n, _ := ed.Document().Get(id)
		ed.SetRect(id, design.Rect{X: x, Y: y, W: w, H: h}, false)
		ed.Rename(n.ID, name)
		return id
	}
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}

	frame := add("lipgloss.box", "Frame", 0, 0, 80, 24)
	must(ed.SetProp(frame, "title", "MVD · Colours"))
	must(ed.SetLayout(frame, editor.AxisW, "100%"))
	must(ed.SetLayout(frame, editor.AxisH, "100%"))

	slots := add("lipgloss.list", "Slots", 2, 2, 76, 10)
	must(ed.SetProp(slots, "items", "Accent (title), Focus (active border), Highlight (keys, bars), Success, Error, Dim (borders, hints), Text, Selected row (background)"))
	must(ed.SetLayout(slots, editor.AxisW, "100% - 4"))
	must(ed.SetBinding(slots, "items", "Slots"))
	must(ed.SetEvent(slots, "Pick slot"))

	help := add("lipgloss.label", "Help", 2, 13, 76, 1)
	must(ed.SetProp(help, "text", "Colours are #rgb or #rrggbb and apply as soon as you accept them."))
	must(ed.SetLayout(help, editor.AxisW, "100% - 4"))
	must(ed.SetBinding(help, "text", "Help"))

	status := add("lipgloss.label", "Status", 2, 15, 76, 1)
	must(ed.SetProp(status, "text", "not a colour"))
	must(ed.SetLayout(status, editor.AxisW, "100% - 4"))
	must(ed.SetBinding(status, "text", "Status"))
	must(ed.SetShowIf(status, "Show status"))

	keys := add("lipgloss.label", "Key bar", 2, 22, 76, 1)
	must(ed.SetProp(keys, "text", "↑↓ move · enter edit · d default · s save · esc back"))
	must(ed.SetLayout(keys, editor.AxisY, "100% - 2"))
	must(ed.SetLayout(keys, editor.AxisW, "100% - 4"))
	must(ed.SetBinding(keys, "text", "Key bar"))

	must(ed.SetKeys("enter=Edit:edit, d=Default:default, s=Save:save, esc=Back:back"))
	return ed.Document()
}

// TestMVDColoursExample keeps examples/mvd-colours.cuppa in step with the
// design above (set CUPPA_WRITE_EXAMPLES=1 to write it) and checks the screen
// it exports: the contract MVD's Colours screen would be written against.
func TestMVDColoursExample(t *testing.T) {
	want := mvdColours(t)
	path := filepath.FromSlash(mvdColoursExample)
	if os.Getenv("CUPPA_WRITE_EXAMPLES") == "1" {
		if _, err := disk.Save(path, want); err != nil {
			t.Fatal(err)
		}
	}
	got, err := disk.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("examples/mvd-colours.cuppa is out of date; run with CUPPA_WRITE_EXAMPLES=1\n got %+v\nwant %+v", got, want)
	}

	p := GenerateScreens([]design.Document{got}, standard.Default(), "screens")
	if len(p.Notes) != 0 {
		t.Fatalf("notes: %v", p.Notes)
	}
	contract := p.Files["mvd_colours_screen_contract.go"]
	for _, wantSrc := range []string{
		"type MVDColoursProps struct",
		"Slots []string",
		"Help string",
		"Status string",
		"ShowStatus bool",
		"KeyBar string",
		"type MVDColoursPickSlot struct{ X, Y int }",
		"type MVDColoursEdit struct{ X, Y int }",
		"type MVDColoursSave struct{ X, Y int }",
		"type MVDColoursBack struct{ X, Y int }",
		"type MVDColoursDefault struct{ X, Y int }",
	} {
		if !strings.Contains(contract, wantSrc) {
			t.Errorf("contract lacks %q:\n%s", wantSrc, contract)
		}
	}
}
