package scene

import (
	"fmt"
	"strings"

	"github.com/fancy-cuppa/cuppa/libs/document/design"
	"github.com/fancy-cuppa/cuppa/libs/render/grid"
)

func paintBubbleTableCommunity(g *grid.Grid, p Props) {
	paintTable(g, Props{"headers": p.Str("columns"), "rows": p.Str("rows"), "border": "rounded", "color": p.Str("color")})
	if f := p.Str("footer"); f != "" {
		rows := len(splitList(p.Str("rows"), ";"))
		if g.H >= rows+6 {
			g.Text(g.W-len([]rune(f))-1, g.H-1, f, dim, g.W)
		}
	}
}

// cells splits length into n slots separated by one-cell gaps and returns the
// start and size of slot i.
func slot(length, n, i int) (start, size int) {
	size = max((length-(n-1))/n, 1)
	return i * (size + 1), size
}

func paintFlexBox(g *grid.Grid, p Props) {
	rows, cols := max(p.Int("rows", 1), 1), max(p.Int("columns", 1), 1)
	style := fg(p.Str("color"))
	for r := 0; r < rows; r++ {
		y, h := slot(g.H, rows, r)
		for c := 0; c < cols; c++ {
			x, w := slot(g.W, cols, c)
			g.Box(design.Rect{X: x, Y: y, W: w, H: h}, grid.BorderNamed("rounded"), style)
			if h >= 3 && w >= 5 {
				g.Text(x+2, y+h/2, fmt.Sprintf("%d,%d", r+1, c+1), dim, w-3)
			}
		}
	}
}

func paintBoxer(g *grid.Grid, p Props) {
	n := max(p.Int("children", 2), 1)
	style := fg(p.Str("color"))
	vertical := p.Str("orientation") == "vertical"
	for i := 0; i < n; i++ {
		var r design.Rect
		if vertical {
			y, h := slot(g.H, n, i)
			r = design.Rect{Y: y, W: g.W, H: h}
		} else {
			x, w := slot(g.W, n, i)
			r = design.Rect{X: x, W: w, H: g.H}
		}
		g.Box(r, grid.BorderNamed("normal"), style)
		if r.W >= 9 && r.H >= 3 {
			g.Text(r.X+2, r.Y+r.H/2, fmt.Sprintf("model %d", i+1), dim, r.W-3)
		}
	}
}

func paintDatePicker(g *grid.Grid, p Props) {
	accent := grid.Style{Fg: p.Str("color"), Bold: true}
	month := p.Str("month")
	g.Text(0, 0, "◂", accent, 1)
	g.Text((g.W-len([]rune(month)))/2, 0, month, bold, g.W)
	g.Text(g.W-1, 0, "▸", accent, 1)
	g.Text(0, 1, "Su Mo Tu We Th Fr Sa", dim, g.W)
	const firstWeekday = 4 // 1 October 2026 is a Thursday.
	day := min(max(p.Int("day", 1), 1), 31)
	for d := 1; d <= 31; d++ {
		slotN := firstWeekday + d - 1
		x, y := (slotN%7)*3, 2+slotN/7
		s := grid.Style{}
		if d == day {
			s = selected(p.Str("color"))
		}
		g.Text(x, y, fmt.Sprintf("%2d", d), s, 2)
	}
}

func paintOverlay(g *grid.Grid, p Props) {
	color := p.Str("color")
	g.Box(full(g), grid.BorderNamed("thick"), fg(color))
	if g.W > 4 {
		g.Text(2, 0, " "+p.Str("title")+" ", grid.Style{Fg: color, Bold: true}, g.W-4)
	}
	for i, l := range wrap(p.Str("body"), g.W-4) {
		if 2+i >= g.H-2 {
			break
		}
		g.Text(2, 1+i, l, grid.Style{}, g.W-4)
	}
	if g.H >= 5 {
		n := g.Text(2, g.H-2, " OK ", selected(color), g.W-4)
		g.Text(3+n, g.H-2, " Cancel ", grid.Style{Fg: "250", Bg: "238"}, g.W-5-n)
	}
}

func paintStatusBar(g *grid.Grid, p Props) {
	g.Fill(full(g), ' ', grid.Style{Bg: "236"})
	x := g.Text(0, 0, " "+p.Str("left")+" ", selected(p.Str("color")), g.W)
	x += g.Text(x, 0, " "+p.Str("middle")+" ", grid.Style{Fg: "252", Bg: "240"}, g.W-x)
	end := " " + p.Str("end") + " "
	right := " " + p.Str("right") + " "
	ew, rw := len([]rune(end)), len([]rune(right))
	if g.W-ew >= x {
		g.Text(g.W-ew, 0, end, grid.Style{Fg: "0", Bg: "250", Bold: true}, ew)
	}
	if g.W-ew-rw >= x {
		g.Text(g.W-ew-rw, 0, right, grid.Style{Fg: "252", Bg: "238"}, rw)
	}
}

func paintFileTree(g *grid.Grid, p Props) {
	entries := p.List("entries")
	sel := pick(p, "selected", len(entries))
	for i, e := range entries {
		if i >= g.H {
			break
		}
		name, style := "  "+e, grid.Style{}
		if strings.HasSuffix(e, "/") {
			name, style = "▸ "+e, grid.Style{Fg: "39"}
		}
		if i == sel {
			style = selected(p.Str("color"))
			g.Fill(designRow(g, i), ' ', style)
		}
		g.Text(0, i, name, style, g.W)
	}
}
