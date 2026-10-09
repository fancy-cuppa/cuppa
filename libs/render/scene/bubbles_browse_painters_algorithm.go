package scene

import (
	"fmt"
	"strings"

	"github.com/meta-tui/cuppa/libs/render/grid"
)

func paintTextArea(g *grid.Grid, p Props) {
	accent := fg(p.Str("color"))
	lines := pipeLines(p.Str("value"))
	showNumbers := p.Bool("line_numbers")
	gutter := 0
	if showNumbers {
		gutter = 4
	}
	if len(lines) == 0 {
		if showNumbers {
			g.Text(0, 0, "  1 ", dim, g.W)
		}
		g.Set(gutter, 0, grid.Cell{Ch: '█', Style: accent})
		g.Text(gutter+1, 0, p.Str("placeholder"), dim, g.W-gutter-1)
		return
	}
	for i, l := range lines {
		if i >= g.H {
			break
		}
		if showNumbers {
			g.Text(0, i, fmt.Sprintf("%3d ", i+1), dim, gutter)
		}
		n := g.Text(gutter, i, l, grid.Style{}, g.W-gutter)
		if i == len(lines)-1 {
			g.Set(gutter+n, i, grid.Cell{Ch: '█', Style: accent})
		}
	}
}

func paintBubbleList(g *grid.Grid, p Props) {
	items := p.List("items")
	sel := pick(p, "selected", len(items))
	y := 0
	g.Text(1, y, " "+p.Str("title")+" ", grid.Style{Fg: "0", Bg: p.Str("color"), Bold: true}, g.W-1)
	y++
	if p.Bool("show_filter") && g.H > 3 {
		g.Text(2, y, "/ filter", dim, g.W-2)
		y++
	}
	footer := 0
	if p.Bool("show_status") && g.H > 4 {
		footer = 1
	}
	room := g.H - y - footer
	for i, it := range items {
		if i >= room {
			break
		}
		if i == sel {
			g.Text(0, y+i, "┃ "+it, grid.Style{Fg: p.Str("color"), Bold: true}, g.W)
		} else {
			g.Text(2, y+i, it, grid.Style{}, g.W-2)
		}
	}
	if footer == 1 {
		g.Text(2, g.H-1, fmt.Sprintf("%d items", len(items)), dim, g.W-2)
	}
}

func paintBubbleTable(g *grid.Grid, p Props) {
	cols := p.List("columns")
	var rows [][]string
	for _, r := range splitList(p.Str("rows"), ";") {
		rows = append(rows, splitList(r, ","))
	}
	n := len(cols)
	for _, r := range rows {
		n = max(n, len(r))
	}
	if n == 0 {
		paintGeneric(g, "Table")
		return
	}
	widths := make([]int, n)
	measure := func(cells []string) {
		for i, c := range cells {
			widths[i] = max(widths[i], len([]rune(c)))
		}
	}
	measure(cols)
	for _, r := range rows {
		measure(r)
	}
	draw := func(y int, cells []string, s grid.Style) {
		x := 1
		for i, w := range widths {
			if i < len(cells) {
				g.Text(x, y, cells[i], s, min(w, g.W-x))
			}
			x += w + 2
		}
	}
	sel := pick(p, "selected", len(rows))
	y := 0
	if len(cols) > 0 {
		draw(y, cols, bold)
		y++
		g.Text(0, y, strings.Repeat("─", g.W), dim, g.W)
		y++
	}
	for i, r := range rows {
		if y >= g.H {
			break
		}
		if i == sel {
			g.Fill(designRow(g, y), ' ', selected(p.Str("color")))
			draw(y, r, selected(p.Str("color")))
		} else {
			draw(y, r, grid.Style{})
		}
		y++
	}
}

func paintViewport(g *grid.Grid, p Props) {
	lines := pipeLines(p.Str("content"))
	room := g.H
	offset := 0
	if len(lines) > room {
		offset = p.Int("percent", 0) * (len(lines) - room) / 100
	}
	textW := g.W
	if p.Bool("show_scrollbar") && g.W > 3 {
		textW = g.W - 2
	}
	for i := 0; i < room && offset+i < len(lines); i++ {
		g.Text(0, i, lines[offset+i], grid.Style{}, textW)
	}
	if textW < g.W {
		track := fg(p.Str("color"))
		thumbH := max(g.H*room/max(len(lines), room), 1)
		thumbY := 0
		if len(lines) > room {
			thumbY = offset * (g.H - thumbH) / (len(lines) - room)
		}
		for y := 0; y < g.H; y++ {
			ch := '│'
			if y >= thumbY && y < thumbY+thumbH {
				ch = '█'
			}
			g.Set(g.W-1, y, grid.Cell{Ch: ch, Style: track})
		}
	}
}

func paintPaginator(g *grid.Grid, p Props) {
	total := max(p.Int("total", 1), 1)
	page := min(max(p.Int("page", 1), 1), total)
	accent := fg(p.Str("color"))
	if p.Str("style") == "arabic" {
		g.Text(0, 0, fmt.Sprintf("%d/%d", page, total), accent, g.W)
		return
	}
	for i := 1; i <= total && (i-1)*2 < g.W; i++ {
		if i == page {
			g.Set((i-1)*2, 0, grid.Cell{Ch: '•', Style: accent})
		} else {
			g.Set((i-1)*2, 0, grid.Cell{Ch: '•', Style: dim})
		}
	}
}

func paintFilePicker(g *grid.Grid, p Props) {
	accent := grid.Style{Fg: p.Str("color"), Bold: true}
	g.Text(0, 0, p.Str("path"), accent, g.W)
	entries := p.List("entries")
	sel := pick(p, "selected", len(entries))
	for i, e := range entries {
		y := i + 1
		if y >= g.H {
			break
		}
		style := grid.Style{}
		if strings.HasSuffix(e, "/") {
			style = grid.Style{Fg: "39"}
		}
		if i == sel {
			g.Text(0, y, "> ", accent, 2)
			style.Bold = true
		} else {
			g.Text(0, y, "  ", grid.Style{}, 2)
		}
		g.Text(2, y, e, style, g.W-2)
	}
}
