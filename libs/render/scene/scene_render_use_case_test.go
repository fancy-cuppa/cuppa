package scene

import (
	"strings"
	"testing"

	"github.com/meta-tui/cuppa/libs/catalog/standard"
	"github.com/meta-tui/cuppa/libs/document/design"
)

func text(lines []string) string {
	var b strings.Builder
	for _, l := range lines {
		esc := false
		for _, r := range l {
			switch {
			case r == 0x1b:
				esc = true
			case esc && r == 'm':
				esc = false
			case !esc:
				b.WriteRune(r)
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func TestEveryShippedComponentRendersAtDefaultAndMinimumSize(t *testing.T) {
	cat := standard.Default()
	for _, def := range cat.List() {
		for _, sz := range []struct{ w, h int }{{def.DefaultSize.W, def.DefaultSize.H}, {def.MinSize.W, def.MinSize.H}} {
			n := design.Node{Component: def.ID, Name: def.Name, Rect: design.Rect{W: sz.w, H: sz.h}}
			g := RenderNode(n, cat) // must not panic at any legal size
			if g.W != sz.w || g.H != sz.h {
				t.Errorf("%s: grid %dx%d, want %dx%d", def.ID, g.W, g.H, sz.w, sz.h)
			}
		}
		if def.Status == "supported" && !Painted(def.ID) {
			t.Errorf("%s is marked supported but has no painter", def.ID)
		}
	}
}

func TestBoxWithTitle(t *testing.T) {
	cat := standard.Default()
	n := design.Node{Component: "lipgloss.box", Rect: design.Rect{W: 12, H: 3}, Props: map[string]string{"title": "Tea", "border": "normal"}}
	got := text(RenderNode(n, cat).Lines())
	want := "┌─ Tea ────┐\n│          │\n└──────────┘\n"
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

func TestRenderOrdersBackToFront(t *testing.T) {
	cat := standard.Default()
	doc := design.NewDocument("t", 10, 1)
	doc.Add(design.Node{Component: "lipgloss.label", Rect: design.Rect{W: 10, H: 1}, Props: map[string]string{"text": "back"}})
	doc.Add(design.Node{Component: "lipgloss.label", Rect: design.Rect{X: 2, W: 4, H: 1}, Props: map[string]string{"text": "TOP"}})
	got := text(Render(doc, cat).Lines())
	if got != "baTOP     \n" && got != "baTOP\n" && !strings.HasPrefix(got, "baTOP") {
		t.Fatalf("got %q", got)
	}
}

func TestUnknownComponentFallsBackToPlaceholder(t *testing.T) {
	cat := standard.Default()
	n := design.Node{Component: "mystery", Name: "Mystery", Rect: design.Rect{W: 12, H: 3}}
	if got := text(RenderNode(n, cat).Lines()); !strings.Contains(got, "Mystery") {
		t.Fatalf("placeholder missing name: %q", got)
	}
}

func TestNewComponentsDrawRecognisableContent(t *testing.T) {
	cat := standard.Default()
	cases := []struct {
		id    string
		w, h  int
		props map[string]string
		want  string
	}{
		{"bubbles.paginator", 10, 1, map[string]string{"total": "3", "page": "2"}, "• • •"},
		{"bubbles.timer", 20, 1, map[string]string{"value": "01:23", "label": "Steep"}, "Steep 01:23"},
		{"huh.confirm", 30, 3, map[string]string{"title": "Sure?"}, "Sure?"},
		{"ntcharts.sparkline", 8, 1, map[string]string{"values": "1,2,3,4,5,6,7,8"}, "▁"},
		{"glamour.markdown", 30, 3, map[string]string{"markdown": "# Hi|**bold** word"}, "bold word"},
		{"community.statusbar", 40, 1, map[string]string{"left": "main.go"}, "main.go"},
		{"bubbles.help", 40, 1, map[string]string{"bindings": "q:quit"}, "q quit"},
	}
	for _, c := range cases {
		n := design.Node{Component: c.id, Rect: design.Rect{W: c.w, H: c.h}, Props: c.props}
		if got := text(RenderNode(n, cat).Lines()); !strings.Contains(got, c.want) {
			t.Errorf("%s: %q not found in\n%s", c.id, c.want, got)
		}
	}
}
