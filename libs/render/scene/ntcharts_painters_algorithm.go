package scene

import (
	"strings"

	"github.com/fancy-cuppa/cuppa/libs/render/grid"
)

func paintBarChart(g *grid.Grid, p Props) {
	vs := numbers(p, "values")
	if len(vs) == 0 {
		paintGeneric(g, "Bar chart")
		return
	}
	labels := p.List("labels")
	style := fg(p.Str("color"))
	_, hi := extent(append([]float64{0}, vs...))
	plotH := max(g.H-1, 1)
	slot := max(g.W/len(vs), 2)
	barW := max(slot-1, 1)
	for i, v := range vs {
		x0 := i * slot
		if x0 >= g.W {
			break
		}
		eighths := int(v / hi * float64(plotH*8))
		for dx := 0; dx < barW; dx++ {
			fillColumn(g, x0+dx, plotH-1, eighths, style)
		}
		if i < len(labels) {
			l := []rune(labels[i])
			if len(l) > barW {
				l = l[:barW]
			}
			g.Text(x0+(barW-len(l))/2, g.H-1, string(l), dim, barW)
		}
	}
}

// axes draws a left and bottom axis and returns the plot rectangle origin and size.
func axes(g *grid.Grid, show bool) (x0, y0, w, h int) {
	if !show || g.W < 4 || g.H < 3 {
		return 0, 0, g.W, g.H
	}
	for y := 0; y < g.H-1; y++ {
		g.Set(0, y, grid.Cell{Ch: '│', Style: dim})
	}
	g.Set(0, g.H-1, grid.Cell{Ch: '└', Style: dim})
	for x := 1; x < g.W; x++ {
		g.Set(x, g.H-1, grid.Cell{Ch: '─', Style: dim})
	}
	return 1, 0, g.W - 1, g.H - 1
}

func paintLineChart(g *grid.Grid, p Props) {
	vs := numbers(p, "values")
	if len(vs) == 0 {
		paintGeneric(g, "Line chart")
		return
	}
	x0, y0, w, h := axes(g, p.Bool("axes"))
	lo, hi := extent(vs)
	style := fg(p.Str("color"))
	yOf := func(v float64) int { return y0 + h - 1 - int((v-lo)/(hi-lo)*float64(h-1)+0.5) }
	prevX, prevY := -1, -1
	for i, v := range vs {
		x := x0
		if len(vs) > 1 {
			x = x0 + i*(w-1)/(len(vs)-1)
		}
		y := yOf(v)
		if prevX >= 0 {
			for cx := prevX + 1; cx < x; cx++ {
				t := float64(cx-prevX) / float64(x-prevX)
				g.Set(cx, prevY+int(t*float64(y-prevY)+0.5), grid.Cell{Ch: '·', Style: style})
			}
		}
		g.Set(x, y, grid.Cell{Ch: '•', Style: style})
		prevX, prevY = x, y
	}
}

func paintSparkline(g *grid.Grid, p Props) {
	vs := resample(numbers(p, "values"), g.W)
	lo, hi := extent(vs)
	style := fg(p.Str("color"))
	for x, v := range vs {
		level := int((v - lo) / (hi - lo) * 7)
		g.Set(x, 0, grid.Cell{Ch: sparkLevels[min(max(level, 0), 7)], Style: style})
	}
}

func paintStreamline(g *grid.Grid, p Props) {
	vs := resample(numbers(p, "values"), g.W)
	_, hi := extent(append([]float64{0}, vs...))
	style := fg(p.Str("color"))
	for x, v := range vs {
		fillColumn(g, x, g.H-1, int(v/hi*float64(g.H*8)), style)
	}
}

func paintTimeSeries(g *grid.Grid, p Props) {
	paintLineChart(g, p)
	if g.H >= 4 && g.W >= 12 && p.Bool("axes") {
		g.Text(1, g.H-1, "─09:00", dim, 6)
		g.Text(g.W-6, g.H-1, "17:00", dim, 6)
	}
}

// heatColors run from cold to hot.
var heatColors = []string{"17", "19", "21", "33", "45", "190", "214", "196"}

func paintHeatMap(g *grid.Grid, p Props) {
	seed := p.Int("seed", 0)
	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			v := (x*7 + y*13 + seed*31 + (x*y)%5*3) % len(heatColors)
			g.Set(x, y, grid.Cell{Ch: '█', Style: grid.Style{Fg: heatColors[v]}})
		}
	}
}

func paintChartCanvas(g *grid.Grid, p Props) {
	x0, _, _, _ := axes(g, true)
	for y := 1; y < g.H-1; y += 2 {
		for x := x0 + 1; x < g.W; x += 4 {
			g.Set(x, y, grid.Cell{Ch: '·', Style: dim})
		}
	}
	if t := p.Str("title"); t != "" && g.W > 6 {
		g.Text(g.W-len([]rune(t))-1, 0, strings.TrimSpace(t), fg(p.Str("color")), g.W)
	}
}
