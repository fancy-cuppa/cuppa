package scene

import (
	"strings"

	"github.com/meta-tui/cuppa/libs/render/grid"
)

// huhFrame draws the focus bar, the title and the description shared by every
// form field, and returns the first free row for the field's own content.
func huhFrame(g *grid.Grid, p Props) int {
	accent := grid.Style{Fg: p.Str("color")}
	for y := 0; y < g.H; y++ {
		g.Set(0, y, grid.Cell{Ch: '┃', Style: accent})
	}
	y := 0
	if t := p.Str("title"); t != "" {
		g.Text(2, y, t, grid.Style{Fg: p.Str("color"), Bold: true}, g.W-2)
		y++
	}
	if d := p.Str("description"); d != "" && y < g.H {
		g.Text(2, y, d, dim, g.W-2)
		y++
	}
	return y
}

func paintHuhInput(g *grid.Grid, p Props) {
	y := huhFrame(g, p)
	accent := fg(p.Str("color"))
	n := g.Text(2, y, "> ", accent, g.W-2)
	if v := p.Str("value"); v != "" {
		m := g.Text(2+n, y, v, grid.Style{}, g.W-2-n)
		g.Set(2+n+m, y, grid.Cell{Ch: '█', Style: accent})
		return
	}
	g.Set(2+n, y, grid.Cell{Ch: '█', Style: accent})
	g.Text(3+n, y, p.Str("placeholder"), dim, g.W-3-n)
}

func paintHuhText(g *grid.Grid, p Props) {
	y := huhFrame(g, p)
	accent := fg(p.Str("color"))
	lines := pipeLines(p.Str("value"))
	if len(lines) == 0 {
		g.Set(2, y, grid.Cell{Ch: '█', Style: accent})
		g.Text(3, y, p.Str("placeholder"), dim, g.W-3)
		return
	}
	for i, l := range lines {
		if y+i >= g.H {
			break
		}
		n := g.Text(2, y+i, l, grid.Style{}, g.W-2)
		if i == len(lines)-1 {
			g.Set(2+n, y+i, grid.Cell{Ch: '█', Style: accent})
		}
	}
}

func paintHuhSelect(g *grid.Grid, p Props) {
	y := huhFrame(g, p)
	opts := p.List("options")
	sel := pick(p, "selected", len(opts))
	for i, o := range opts {
		if y+i >= g.H {
			break
		}
		if i == sel {
			g.Text(2, y+i, "> "+o, grid.Style{Fg: p.Str("color"), Bold: true}, g.W-2)
		} else {
			g.Text(2, y+i, "  "+o, grid.Style{}, g.W-2)
		}
	}
}

func paintHuhMultiSelect(g *grid.Grid, p Props) {
	y := huhFrame(g, p)
	opts := p.List("options")
	checked := map[string]bool{}
	for _, c := range p.List("checked") {
		checked[c] = true
	}
	cursor := pick(p, "cursor", len(opts))
	for i, o := range opts {
		if y+i >= g.H {
			break
		}
		box := "[ ] "
		if checked[strings.TrimSpace(itoa(i))] {
			box = "[•] "
		}
		prefix := "  "
		style := grid.Style{}
		if i == cursor {
			prefix, style = "> ", grid.Style{Fg: p.Str("color"), Bold: true}
		}
		g.Text(2, y+i, prefix+box+o, style, g.W-2)
	}
}

func paintHuhConfirm(g *grid.Grid, p Props) {
	y := huhFrame(g, p)
	yes, no := " "+p.Str("affirmative")+" ", " "+p.Str("negative")+" "
	on, off := selected(p.Str("color")), grid.Style{Fg: "250", Bg: "238"}
	ys, ns := off, on
	if p.Bool("value") {
		ys, ns = on, off
	}
	n := g.Text(2, y, yes, ys, g.W-2)
	g.Text(2+n+1, y, no, ns, g.W-3-n)
}

func paintHuhNote(g *grid.Grid, p Props) {
	y := huhFrame(g, p)
	for i, l := range wrap(p.Str("body"), g.W-2) {
		if y+i >= g.H {
			break
		}
		g.Text(2, y+i, l, grid.Style{}, g.W-2)
	}
}

func paintHuhFilePicker(g *grid.Grid, p Props) {
	y := huhFrame(g, p)
	entries := p.List("entries")
	sel := pick(p, "selected", len(entries))
	for i, e := range entries {
		if y+i >= g.H {
			break
		}
		style := grid.Style{}
		if strings.HasSuffix(e, "/") {
			style = grid.Style{Fg: "39"}
		}
		prefix := "  "
		if i == sel {
			prefix, style.Bold = "> ", true
		}
		g.Text(2, y+i, prefix+e, style, g.W-2)
	}
}

func paintHuhForm(g *grid.Grid, p Props) {
	color := p.Str("color")
	g.Box(full(g), grid.BorderNamed("rounded"), fg(color))
	if t := p.Str("title"); t != "" && g.W > 4 {
		g.Text(2, 0, " "+t+" ", grid.Style{Fg: color, Bold: true}, g.W-4)
	}
	rows := (g.H - 3) / 2
	for i := 0; i < rows; i++ {
		y := 1 + i*2
		g.Set(2, y, grid.Cell{Ch: '┃', Style: fg(color)})
		g.Text(4, y, strings.Repeat("▁", max(g.W-8, 0)), dim, g.W-6)
	}
	if f := p.Str("footer"); f != "" && g.H > 2 {
		g.Text(2, g.H-2, f, dim, g.W-4)
	}
}
