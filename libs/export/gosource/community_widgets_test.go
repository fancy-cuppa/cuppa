package gosource

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/meta-tui/cuppa/libs/catalog/standard"
	"github.com/meta-tui/cuppa/libs/document/design"
	"github.com/meta-tui/cuppa/libs/render/scene"
)

// TestWidgetsDrawTheSameInTheDesignerAndTheProgram builds a program with every
// variant of the community widgets and checks that each draws, character for
// character, what the designer draws. The two drawings are written twice (the
// painter and community_widgets_go.txt), so this is what keeps them one.
func TestWidgetsDrawTheSameInTheDesignerAndTheProgram(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a project with the go tool")
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("no go tool")
	}
	cat := standard.Default()
	variants := []struct {
		name, id string
		w, h     int
		props    map[string]string
	}{
		{"dropdown open", "community.dropdown", 26, 7, nil},
		{"dropdown closed", "community.dropdown", 26, 7, map[string]string{"open": "false", "selected": "-1"}},
		{"dropdown scrolled", "community.dropdown", 20, 5, map[string]string{"selected": "3"}},
		{"input", "community.promptinput", 36, 2, map[string]string{"value": "Sencha"}},
		{"input placeholder", "community.promptinput", 36, 2, nil},
		{"input hidden", "community.promptinput", 36, 2, map[string]string{"value": "secret", "hidden": "true", "error": "too short"}},
		{"select", "community.promptselect", 30, 6, nil},
		{"select filtered", "community.promptselect", 30, 6, map[string]string{"filter": "e", "selected": "2"}},
		{"select scrolled", "community.promptselect", 30, 4, map[string]string{"selected": "6"}},
		{"tree", "community.datatree", 34, 9, nil},
		{"pdf text", "community.pdfview", 44, 14, nil},
		{"pdf image", "community.pdfview", 44, 14, map[string]string{"mode": "image"}},
		{"3d surface", "ntcharts.chart3d", 44, 14, nil},
		{"3d scatter", "ntcharts.chart3d", 44, 14, map[string]string{"kind": "scatter"}},
		{"3d bar", "ntcharts.chart3d", 44, 14, map[string]string{"kind": "bar"}},
		{"3d line", "ntcharts.chart3d", 44, 14, map[string]string{"kind": "line"}},
		{"3d vector", "ntcharts.chart3d", 44, 14, map[string]string{"kind": "vector", "legend": "false", "title": ""}},
		{"rows", "lipgloss.rows", 44, 5, nil},
		{"rows spaces", "lipgloss.rows", 30, 3, map[string]string{"columns": "A:6,B:6,C", "rows": "  x,,z; y\\,w,,; q,r,s", "styles": "accent,dim,selected"}},
		{"rows colours", "lipgloss.rows", 30, 3, map[string]string{"columns": "Name:8,Swatch:5:colour,Rest", "rows": "Tea,#ff0000,hot;Milk,,cold", "styles": "selected"}},
	}
	doc := design.NewDocument("widgets", 200, 400)
	var nodes []design.Node
	for i, v := range variants {
		n := design.Node{Component: v.id, Name: v.name, Rect: design.Rect{X: 0, Y: i * 15, W: v.w, H: v.h}, Props: v.props}
		doc.Add(n)
		nodes = append(nodes, n)
	}
	dir := filepath.Join(t.TempDir(), "app")
	p := Generate(doc, cat)
	if len(p.Notes) != 0 {
		t.Fatalf("notes: %v", p.Notes)
	}
	if err := Write(dir, p); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dump_test.go"), []byte(dumpTest), 0o644); err != nil {
		t.Fatal(err)
	}
	env := append(os.Environ(), "GOFLAGS=-mod=mod", "GOWORK=off")
	for _, args := range [][]string{{"mod", "tidy"}, {"vet", "."}} {
		cmd := exec.Command(goBin, args...)
		cmd.Dir, cmd.Env = dir, env
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("go %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run := exec.Command(goBin, "test", "-run", "TestDumpViews", "-v", ".")
	run.Dir, run.Env = dir, env
	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("go test: %v\n%s", err, out)
	}
	var views map[string][]string
	for _, line := range strings.Split(string(out), "\n") {
		if rest, ok := strings.CutPrefix(line, "VIEWS="); ok {
			if err := json.Unmarshal([]byte(rest), &views); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, n := range nodes {
		want := trimRows(views[n.Name], n.Rect.H)
		got := trimRows(strings.Split(ansi.Strip(strings.Join(scene.RenderNode(n, cat).Lines(), "\n")), "\n"), n.Rect.H)
		if strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Errorf("%s differs\n-- designer\n%s\n-- program\n%s", n.Name, strings.Join(got, "\n"), strings.Join(want, "\n"))
		}
	}
}

func trimRows(lines []string, h int) []string {
	out := make([]string, h)
	for i := range out {
		if i < len(lines) {
			out[i] = strings.TrimRight(lines[i], " ")
		}
	}
	return out
}
