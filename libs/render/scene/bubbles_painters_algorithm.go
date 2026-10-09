package scene

import (
	"strings"

	"github.com/meta-tui/cuppa/libs/render/grid"
)

func paintTextInput(g *grid.Grid, p Props) {
	accent := fg(p.Str("color"))
	prompt := p.Str("prompt")
	n := g.Text(0, 0, prompt, accent, g.W)
	value := p.Str("value")
	if value == "" {
		g.Text(n+1, 0, p.Str("placeholder"), dim, g.W-n-1)
		g.Set(n, 0, grid.Cell{Ch: '█', Style: accent})
		return
	}
	m := g.Text(n, 0, value, grid.Style{}, g.W-n)
	g.Set(n+m, 0, grid.Cell{Ch: '█', Style: accent})
}

var spinnerFrames = map[string]string{
	"line": "|", "dot": "⣾", "minidot": "⠋", "jump": "⢄", "pulse": "█",
	"points": "∙", "globe": "◐", "moon": "◑", "monkey": "o",
}

func paintSpinner(g *grid.Grid, p Props) {
	frame := spinnerFrames[p.Str("style")]
	if frame == "" {
		frame = "⣾"
	}
	accent := fg(p.Str("color"))
	g.Text(0, 0, frame, accent, 1)
	g.Text(2, 0, p.Str("label"), grid.Style{}, g.W-2)
}

func paintProgress(g *grid.Grid, p Props) {
	pct := min(max(p.Int("percent", 0), 0), 100)
	label := ""
	if p.Bool("show_percentage") {
		label = " " + pad3(pct) + "%"
	}
	barW := max(g.W-len(label), 1)
	filled := barW * pct / 100
	accent := fg(p.Str("color"))
	g.Text(0, 0, strings.Repeat("█", filled), accent, barW)
	g.Text(filled, 0, strings.Repeat("░", barW-filled), dim, barW-filled)
	g.Text(barW, 0, label, grid.Style{}, 0)
}

func pad3(n int) string {
	s := []byte("   ")
	for i := 2; i >= 0 && n > 0; i-- {
		s[i] = byte('0' + n%10)
		n /= 10
	}
	if s[2] == ' ' {
		s[2] = '0'
	}
	return string(s)
}
