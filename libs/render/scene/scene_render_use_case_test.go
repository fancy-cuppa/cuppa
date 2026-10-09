package scene

import (
	"strings"
	"testing"

	"github.com/fancy-cuppa/cuppa/libs/catalog/standard"
	"github.com/fancy-cuppa/cuppa/libs/document/design"
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
