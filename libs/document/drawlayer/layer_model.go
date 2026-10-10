// Package drawlayer is the drawing of a design: one layer of styled cells that
// the drawing tools paint into and erase from. It is stored as the "cells"
// property of a single node, so it saves, undoes and exports like any
// component.
package drawlayer

import (
	"sort"
	"strconv"
	"strings"
)

const (
	// Component is the id of the drawing's component.
	Component = "draw.layer"
	// PropCells is the property that holds the encoded cells.
	PropCells = "cells"
)

// Cell is one painted cell of the drawing.
type Cell struct {
	X, Y int
	Ch   rune
	// Fg and Bg are colours as designs store them ("212", "#ff87d7") or "".
	Fg, Bg string
}

// Layer is the set of painted cells, at most one to a position.
type Layer struct {
	cells map[[2]int]Cell
}

// New returns an empty layer.
func New() *Layer { return &Layer{cells: map[[2]int]Cell{}} }

// Len is the number of painted cells.
func (l *Layer) Len() int { return len(l.cells) }

// Get returns the cell at (x, y), if painted.
func (l *Layer) Get(x, y int) (Cell, bool) {
	c, ok := l.cells[[2]int{x, y}]
	return c, ok
}

// Paint puts cells on the layer, replacing what was there. A cell with no
// character and no background paints nothing.
func (l *Layer) Paint(cells ...Cell) {
	for _, c := range cells {
		if c.Ch == 0 && c.Bg == "" {
			continue
		}
		if c.Ch == 0 {
			c.Ch = ' '
		}
		l.cells[[2]int{c.X, c.Y}] = c
	}
}

// Erase removes the cells at the positions of cells, so the components under
// the drawing show again.
func (l *Layer) Erase(cells ...Cell) {
	for _, c := range cells {
		delete(l.cells, [2]int{c.X, c.Y})
	}
}

// Cells lists the painted cells in reading order.
func (l *Layer) Cells() []Cell {
	out := make([]Cell, 0, len(l.cells))
	for _, c := range l.cells {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Y != out[j].Y {
			return out[i].Y < out[j].Y
		}
		return out[i].X < out[j].X
	})
	return out
}

// Encode writes the layer as text: "x,y,codepoint,fg,bg" per cell, joined by
// semicolons, in reading order.
func (l *Layer) Encode() string {
	var b strings.Builder
	for i, c := range l.Cells() {
		if i > 0 {
			b.WriteByte(';')
		}
		b.WriteString(strconv.Itoa(c.X) + "," + strconv.Itoa(c.Y) + "," + strconv.FormatInt(int64(c.Ch), 16) + "," + c.Fg + "," + c.Bg)
	}
	return b.String()
}

// Decode reads what Encode wrote. Entries it cannot read are skipped.
func Decode(s string) *Layer {
	l := New()
	for _, entry := range strings.Split(s, ";") {
		f := strings.Split(entry, ",")
		if len(f) != 5 {
			continue
		}
		x, errX := strconv.Atoi(f[0])
		y, errY := strconv.Atoi(f[1])
		r, errR := strconv.ParseInt(f[2], 16, 32)
		if errX != nil || errY != nil || errR != nil {
			continue
		}
		l.Paint(Cell{X: x, Y: y, Ch: rune(r), Fg: f[3], Bg: f[4]})
	}
	return l
}
