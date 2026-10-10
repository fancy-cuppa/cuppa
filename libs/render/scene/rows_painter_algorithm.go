package scene

import (
	"strconv"
	"strings"

	"github.com/meta-tui/cuppa/libs/document/design"
	"github.com/meta-tui/cuppa/libs/render/grid"
)

// rowsColumn is one column of the Rows component: its name, its width (0 for
// the rest of the row) and whether its cells are colours.
type rowsColumn struct {
	name   string
	width  int
	colour bool
}

// rowsColumns reads "Name[:width][:colour]" items separated by commas. A
// column without a width takes twelve cells, except the last, which takes what
// is left of the row.
func rowsColumns(spec string) []rowsColumn {
	var out []rowsColumn
	for _, item := range rowsSplit(spec, ",", true) {
		parts := strings.Split(strings.TrimSpace(item), ":")
		col := rowsColumn{name: strings.TrimSpace(parts[0])}
		for _, extra := range parts[1:] {
			extra = strings.ToLower(strings.TrimSpace(extra))
			if n, err := strconv.Atoi(extra); err == nil {
				col.width = max(n, 0)
			} else if extra == "colour" || extra == "color" {
				col.colour = true
			}
		}
		out = append(out, col)
	}
	for i := range out {
		if out[i].width == 0 && i < len(out)-1 {
			out[i].width = 12
		}
	}
	return out
}

// rowsSplit cuts s at sep. Unlike the list splitting it keeps empty pieces and
// whitespace, because a cell may be empty or start with spaces. A backslash
// before a comma, a semicolon or another backslash is that character itself
// (unescape false keeps the backslash, for the first cut of the rows).
func rowsSplit(s, sep string, unescape bool) []string {
	var out []string
	var cur strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) && (s[i+1] == ',' || s[i+1] == ';' || s[i+1] == '\\') {
			if !unescape {
				cur.WriteByte(c)
			}
			cur.WriteByte(s[i+1])
			i++
			continue
		}
		if strings.HasPrefix(s[i:], sep) {
			out = append(out, cur.String())
			cur.Reset()
			i += len(sep) - 1
			continue
		}
		cur.WriteByte(c)
	}
	return append(out, cur.String())
}

// rowsCells reads the rows: cells separated by commas, rows by semicolons.
func rowsCells(s string) [][]string {
	var rows [][]string
	for _, raw := range rowsSplit(s, ";", false) {
		if raw == "" {
			continue
		}
		rows = append(rows, rowsSplit(raw, ",", true))
	}
	return rows
}

// paintRows draws the Rows component: one line per row, the cells laid out in
// the columns the design names, and each row in the style it was given:
// "selected" inverts the row in the accent colour, "dim" mutes it and "accent"
// colours its text.
func paintRows(g *grid.Grid, p Props) {
	accent := p.Str("color")
	columns := rowsColumns(p.Str("columns"))
	styles := p.List("styles")
	for y, row := range rowsCells(p.Str("rows")) {
		if y >= g.H {
			break
		}
		style := grid.Style{}
		if y < len(styles) {
			switch styles[y] {
			case "selected":
				style = selected(accent)
				g.Fill(design.Rect{X: 0, Y: y, W: g.W, H: 1}, ' ', style)
			case "dim":
				style = p.Dim()
			case "accent":
				style = grid.Style{Fg: accent}
			}
		}
		x := 0
		for i, col := range columns {
			w := col.width
			if w == 0 {
				w = max(g.W-x, 0)
			}
			w = min(w, max(g.W-x, 0))
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			switch {
			case w <= 0:
			case col.colour && cell != "":
				g.Fill(design.Rect{X: x, Y: y, W: w, H: 1}, ' ', grid.Style{Bg: cell})
			case col.colour:
				g.Text(x, y, "·", p.Dim(), w)
			default:
				g.Text(x, y, cell, style, w)
			}
			x += w
		}
	}
}
