package gosource

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meta-tui/cuppa/libs/catalog/cupp"
	"github.com/meta-tui/cuppa/libs/catalog/standard"
	"github.com/meta-tui/cuppa/libs/document/design"
)

// richDesign uses every component that has real code, one that does not, a
// group, and a component from a pack.
func richDesign() (design.Document, Catalog) {
	pack := cupp.Pack{ID: "tea-shop", Name: "Tea shop", Components: []design.Composite{{
		ID: "card", Name: "Card", W: 20, H: 5,
		Nodes: []design.Node{
			{ID: "a", Component: "lipgloss.box", Name: "Frame", Rect: design.Rect{W: 20, H: 5}},
			{ID: "b", Component: "lipgloss.label", Name: "Title", Rect: design.Rect{X: 2, Y: 1, W: 10, H: 1}},
		},
		Props: []design.Exposed{{Key: "title", Label: "Title", Kind: "text", Default: "Tea", Target: "b", TargetProp: "text"}},
	}}}
	cat, _ := cupp.Extend(standard.Default(), []cupp.Pack{pack})

	doc := design.NewDocument("My Tea \"Shop\" #1", 120, 40)
	doc.Background = "#102030"
	at := 0
	for _, id := range []string{
		"lipgloss.box", "lipgloss.label", "bubbles.textinput", "bubbles.textarea", "bubbles.list", "bubbles.table",
		"bubbles.viewport", "bubbles.paginator", "bubbles.spinner", "bubbles.progress", "bubbles.stopwatch", "bubbles.timer",
		"huh.input",
	} {
		def, _ := cat.Get(id)
		doc.Add(design.Node{Component: id, Name: def.Name, Rect: design.Rect{X: (at % 4) * 28, Y: (at / 4) * 8, W: def.DefaultSize.W, H: def.DefaultSize.H}})
		at++
	}
	doc.Add(design.Node{Component: "tea-shop.card", Name: "Card", Rect: design.Rect{X: 80, Y: 30, W: 30, H: 6},
		Props: map[string]string{"title": "Matcha \"latte\""}})
	doc.Add(design.Node{Component: design.GroupComponent, Name: "Group", Rect: design.Rect{X: 2, Y: 30, W: 40, H: 8}, BaseW: 20, BaseH: 4,
		Children: []design.Node{
			{ID: "g1", Component: "lipgloss.box", Name: "Frame", Rect: design.Rect{W: 20, H: 4}},
			{ID: "g2", Component: "lipgloss.label", Name: "Text", Rect: design.Rect{X: 1, Y: 1, W: 10, H: 1}},
		}})
	doc.Add(design.Node{Component: "lipgloss.box", Name: "Hidden", Rect: design.Rect{W: 5, H: 3}, Hidden: true})
	return doc, cat
}

func TestTheProjectHasTheFilesAndTheLayoutOfTheDesign(t *testing.T) {
	doc, cat := richDesign()
	p := Generate(doc, cat)
	for _, name := range []string{"go.mod", "main.go", "runtime.go", "layout.go", "README.md"} {
		if p.Files[name] == "" {
			t.Errorf("missing %s", name)
		}
	}
	if p.Module != "my-tea-shop-1" || !strings.Contains(p.Files["go.mod"], "module my-tea-shop-1") {
		t.Fatalf("module = %q", p.Module)
	}
	layout := p.Files["layout.go"]
	for _, want := range []string{
		"canvasW    = 120", `background = "#102030"`, `Kind: "bubbles.textinput"`,
		`"text": "Matcha \"latte\""`, // the pack component's exposed value reaches its part
	} {
		if !strings.Contains(layout, want) {
			t.Errorf("layout.go lacks %s", want)
		}
	}
	if strings.Contains(layout, "tea-shop.card") || strings.Contains(layout, "cuppa.group") {
		t.Error("groups and pack components are expanded into their parts")
	}
	if strings.Contains(layout, `Name: "Hidden"`) {
		t.Error("hidden components are not generated")
	}
	if len(p.Notes) != 1 || !strings.Contains(p.Notes[0], "huh.input") || !strings.Contains(p.Files["README.md"], "huh.input") {
		t.Fatalf("what is not generated is said: %v", p.Notes)
	}
}

func TestGroupPartsAreScaledAndPlacedAtTheGroup(t *testing.T) {
	doc, cat := richDesign()
	layout := Generate(doc, cat).Files["layout.go"]
	// The group is 40 wide for a 20-wide base: its frame spans 40 cells from x=2.
	if !strings.Contains(layout, `Name: "Frame", X: 2, Y: 30, W: 40, H: 8`) {
		t.Fatalf("the grouped frame is scaled to the group:\n%s", layout)
	}
}

func TestTheSameDesignGivesTheSameFiles(t *testing.T) {
	doc, cat := richDesign()
	a, b := Generate(doc, cat), Generate(doc, cat)
	for name := range a.Files {
		if a.Files[name] != b.Files[name] {
			t.Fatalf("%s differs between runs", name)
		}
	}
}

func TestWriteRefusesAFolderThatAlreadyHoldsAProject(t *testing.T) {
	doc, cat := richDesign()
	p := Generate(doc, cat)
	dir := filepath.Join(t.TempDir(), "out")
	if err := Write(dir, p); err != nil {
		t.Fatal(err)
	}
	if err := Write(dir, p); err == nil {
		t.Fatal("a second export into the same folder must not overwrite")
	}
}

func TestSlug(t *testing.T) {
	for in, want := range map[string]string{"Untitled": "untitled", "  Hello, World!": "hello-world", "": "cuppa-design", "???": "cuppa-design"} {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestTheGeneratedProjectCompiles is the point of the exporter: what it writes
// builds and passes go vet. It needs the go tool and the module proxy (or a
// warm module cache), so it is skipped in -short runs.
func TestTheGeneratedProjectCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a project with the go tool")
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("no go tool")
	}
	doc, cat := richDesign()
	dir := filepath.Join(t.TempDir(), "app")
	if err := Write(dir, Generate(doc, cat)); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"mod", "tidy"}, {"vet", "./..."}, {"build", "./..."}} {
		cmd := exec.Command(goBin, args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOWORK=off")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("go %s failed: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
}
