package scene

import (
	"strconv"
	"strings"

	"github.com/meta-tui/cuppa/libs/render/grid"
)

func itoa(n int) string { return strconv.Itoa(n) }

// markdownPalette is how one Glamour style colours the parts of a document.
type markdownPalette struct{ text, heading, code, quote, bg string }

var markdownStyles = map[string]markdownPalette{
	"dark":  {text: "252", heading: "212", code: "215", quote: "245"},
	"light": {text: "235", heading: "129", code: "130", quote: "243", bg: "255"},
	"pink":  {text: "218", heading: "205", code: "213", quote: "175"},
	"notty": {},
}

func paintMarkdown(g *grid.Grid, p Props) {
	pal, ok := markdownStyles[p.Str("style")]
	if !ok {
		pal = markdownStyles["dark"]
	}
	base := grid.Style{Fg: pal.text, Bg: pal.bg}
	g.Fill(full(g), ' ', base)
	for y, line := range pipeLines(p.Str("markdown")) {
		if y >= g.H {
			break
		}
		switch {
		case strings.HasPrefix(line, "# "):
			g.Text(0, y, line[2:], grid.Style{Fg: pal.heading, Bg: pal.bg, Bold: true}, g.W)
		case strings.HasPrefix(line, "## "), strings.HasPrefix(line, "### "):
			g.Text(0, y, strings.TrimLeft(line, "# "), grid.Style{Fg: pal.heading, Bg: pal.bg, Bold: true}, g.W)
		case strings.HasPrefix(line, "- "), strings.HasPrefix(line, "* "):
			g.Text(0, y, "• ", grid.Style{Fg: pal.heading, Bg: pal.bg}, 2)
			inline(g, 2, y, line[2:], base, pal)
		case strings.HasPrefix(line, "> "):
			g.Text(0, y, "│ ", grid.Style{Fg: pal.quote, Bg: pal.bg}, 2)
			inline(g, 2, y, line[2:], grid.Style{Fg: pal.quote, Bg: pal.bg}, pal)
		default:
			inline(g, 0, y, line, base, pal)
		}
	}
}

// inline writes a line, styling **bold** and `code` spans.
func inline(g *grid.Grid, x, y int, line string, base grid.Style, pal markdownPalette) {
	rs := []rune(line)
	style := base
	for i := 0; i < len(rs) && x < g.W; i++ {
		switch {
		case rs[i] == '`':
			if style.Fg == pal.code && pal.code != "" {
				style = base
			} else {
				style = grid.Style{Fg: pal.code, Bg: pal.bg}
			}
		case rs[i] == '*' && i+1 < len(rs) && rs[i+1] == '*':
			style.Bold = !style.Bold
			i++
		default:
			g.Set(x, y, grid.Cell{Ch: rs[i], Style: style})
			x++
		}
	}
}
