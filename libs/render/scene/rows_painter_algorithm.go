package scene

import (
	"strconv"
	"strings"

	"github.com/meta-tui/cuppa/libs/document/design"
	"github.com/meta-tui/cuppa/libs/render/grid"
)

// rowsColumn is one column of the Rows component: its name, its width (0 for
// the rest of the row) and how its cells are drawn.
type rowsColumn struct {
	name  string
	width int
	// colour makes the cell a block of that colour (a background); glyph makes
	// it full blocks drawn in that colour.
	colour, glyph bool
	// ellipsis cuts a cell that is too long with "…".
	ellipsis bool
	// fg and bg are the colour of the column's text and its background; selBg
	// is the background the column takes on a selected row. They are colours or
	// @Name of the palette.
	fg, bg, selBg string
}

// rowsColumns reads "Name[:width][:token]..." items separated by commas. A
// token is colour, glyph, ellipsis, fg=<colour>, bg=<colour> or
// selbg=<colour>. A column without a width takes twelve cells, except the
// last, which takes what is left of the row.
func rowsColumns(spec string) []rowsColumn {
	var out []rowsColumn
	for _, item := range rowsSplit(spec, ",", true) {
		parts := strings.Split(strings.TrimSpace(item), ":")
		col := rowsColumn{name: strings.TrimSpace(parts[0])}
		for _, extra := range parts[1:] {
			extra = strings.TrimSpace(extra)
			lower := strings.ToLower(extra)
			if n, err := strconv.Atoi(extra); err == nil {
				col.width = max(n, 0)
				continue
			}
			switch {
			case lower == "colour" || lower == "color":
				col.colour = true
			case lower == "glyph":
				col.glyph = true
			case lower == "ellipsis":
				col.ellipsis = true
			case strings.HasPrefix(lower, "fg="):
				col.fg = strings.TrimSpace(extra[3:])
			case strings.HasPrefix(lower, "bg="):
				col.bg = strings.TrimSpace(extra[3:])
			case strings.HasPrefix(lower, "selbg="):
				col.selBg = strings.TrimSpace(extra[6:])
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

// paletteColour resolves a colour of a column or a cell: a name of the
// design's palette (@Name) through the palette the screen hands its parts, any
// other text as it is.
func paletteColour(p Props, c string) string {
	c = strings.TrimSpace(c)
	name, isName := strings.CutPrefix(c, "@")
	if !isName {
		return c
	}
	for _, entry := range strings.Split(p.Str("theme.palette"), ";") {
		if key, value, ok := strings.Cut(entry, "="); ok && key == name {
			return value
		}
	}
	return ""
}

// cellStyleOf reads the style a program gave one cell: "fg|bg|flags" with the
// flags b (bold), d (dim) and r (reverse).
func cellStyleOf(p Props, code string, base grid.Style) grid.Style {
	if code == "" || code == "-" {
		return base
	}
	parts := strings.Split(code, "|")
	if len(parts) > 0 && parts[0] != "" {
		base.Fg = paletteColour(p, parts[0])
	}
	if len(parts) > 1 && parts[1] != "" {
		base.Bg = paletteColour(p, parts[1])
	}
	if len(parts) > 2 {
		base.Bold = base.Bold || strings.Contains(parts[2], "b")
		base.Dim = base.Dim || strings.Contains(parts[2], "d")
		base.Reverse = base.Reverse || strings.Contains(parts[2], "r")
	}
	return base
}

// paintRows draws the Rows component: one line per row, the cells laid out in
// the columns the design names, and each row in the style it was given:
// "selected" inverts the row in the accent colour (or, when a column has a
// selbg, only gives those columns their background), "dim" mutes it and
// "accent" colours its text. A column can set its own colours, cut with an
// ellipsis, or draw a colour as a block or as glyphs; a program can restyle a
// cell, and the text of a cell may carry SGR colours.
func paintRows(g *grid.Grid, p Props) {
	accent := p.Str("color")
	columns := rowsColumns(p.Str("columns"))
	styles := p.List("styles")
	cellStyles := rowsCells(p.Str("cellstyles"))
	anySelBg := false
	for _, c := range columns {
		anySelBg = anySelBg || c.selBg != ""
	}
	for y, row := range rowsCells(p.Str("rows")) {
		if y >= g.H {
			break
		}
		base := grid.Style{}
		isSelected := false
		if y < len(styles) {
			switch styles[y] {
			case "selected":
				isSelected = true
				if !anySelBg {
					base = selected(accent)
					g.Fill(design.Rect{X: 0, Y: y, W: g.W, H: 1}, ' ', base)
				}
			case "dim":
				base = p.Dim()
			case "accent":
				base = grid.Style{Fg: accent}
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
			style := base
			if col.fg != "" {
				style.Fg = paletteColour(p, col.fg)
			}
			if col.bg != "" {
				style.Bg = paletteColour(p, col.bg)
			}
			if isSelected && col.selBg != "" {
				style.Bg = paletteColour(p, col.selBg)
			}
			if y < len(cellStyles) && i < len(cellStyles[y]) {
				style = cellStyleOf(p, cellStyles[y][i], style)
			}
			switch {
			case w <= 0:
			case col.colour && cell != "":
				g.Fill(design.Rect{X: x, Y: y, W: w, H: 1}, ' ', grid.Style{Bg: cell})
			case col.glyph && cell != "":
				g.Text(x, y, strings.Repeat("█", w), grid.Style{Fg: cell}, w)
			case col.colour || col.glyph:
				g.Text(x, y, "·", p.Dim(), w)
			default:
				if style.Bg != "" {
					g.Fill(design.Rect{X: x, Y: y, W: w, H: 1}, ' ', grid.Style{Bg: style.Bg})
				}
				at := x
				for _, run := range cutRuns(parseSGR(cell, style), w, col.ellipsis) {
					at += g.Text(at, y, run.text, run.style, 0)
				}
			}
			x += w
		}
	}
}
