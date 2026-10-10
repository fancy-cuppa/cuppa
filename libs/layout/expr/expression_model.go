// Package expr is a size expression: a number of cells, a percentage of the
// parent, or a calculation over them. "10", "50%", "100% - 10" and
// "min(50%, 40)" are expressions.
package expr

import (
	"math"
	"strconv"
	"strings"
)

// Expr is a parsed expression. The zero value is the constant 0.
type Expr struct {
	root node
	src  string
}

// Resolve returns the expression in whole cells for a parent that is parent
// cells long on the same axis. The result is rounded down; it is not clamped,
// so a calculation can come out negative. A division by zero gives 0.
func (e Expr) Resolve(parent int) int {
	if e.root == nil {
		return 0
	}
	v := e.root.eval(float64(parent))
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return int(math.Floor(v + 1e-9))
}

// IsFixed reports whether the expression does not depend on the parent size.
func (e Expr) IsFixed() bool { return e.root == nil || !e.root.usesParent() }

// GoSource writes the expression as a Go float64 expression for generated
// code. parent is the name of an int variable that holds the parent's size on
// the axis; the code uses a function div(a, b float64) float64 that gives 0
// for a division by zero, and the built-in min and max. Wrap it in a rounding
// down to get whole cells, as Resolve does.
func (e Expr) GoSource(parent string) string {
	if e.root == nil {
		return "0.0"
	}
	return e.root.goSource(parent)
}

// String is the expression as written, trimmed.
func (e Expr) String() string {
	if e.root == nil {
		return "0"
	}
	return e.src
}

// Cells is the expression for a fixed number of cells.
func Cells(n int) Expr {
	return Expr{root: number{v: float64(n)}, src: strconv.Itoa(n)}
}

// Percent is the expression for p percent of the parent.
func Percent(p float64) Expr {
	s := strconv.FormatFloat(p, 'f', -1, 64) + "%"
	return Expr{root: percent{v: p}, src: s}
}

// node is one part of the expression tree.
type node interface {
	eval(parent float64) float64
	usesParent() bool
	// goSource is the node as a Go float64 expression; parent names the
	// variable that holds the parent's size.
	goSource(parent string) string
}

type number struct{ v float64 }

func (n number) eval(float64) float64 { return n.v }
func (number) usesParent() bool       { return false }
func (n number) goSource(string) string {
	s := strconv.FormatFloat(n.v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

type percent struct{ v float64 }

func (p percent) eval(parent float64) float64 { return parent * p.v / 100 }
func (percent) usesParent() bool              { return true }
func (p percent) goSource(parent string) string {
	return "(float64(" + parent + ") * " + number(p).goSource("") + " / 100.0)"
}

type negation struct{ x node }

func (n negation) eval(parent float64) float64 { return -n.x.eval(parent) }
func (n negation) usesParent() bool            { return n.x.usesParent() }
func (n negation) goSource(parent string) string { return "(-" + n.x.goSource(parent) + ")" }

type binary struct {
	op   byte
	l, r node
}

func (b binary) eval(parent float64) float64 {
	l, r := b.l.eval(parent), b.r.eval(parent)
	switch b.op {
	case '+':
		return l + r
	case '-':
		return l - r
	case '*':
		return l * r
	}
	if r == 0 {
		return 0
	}
	return l / r
}
func (b binary) usesParent() bool { return b.l.usesParent() || b.r.usesParent() }

func (b binary) goSource(parent string) string {
	l, r := b.l.goSource(parent), b.r.goSource(parent)
	if b.op == '/' {
		return "div(" + l + ", " + r + ")"
	}
	return "(" + l + " " + string(b.op) + " " + r + ")"
}

type call struct {
	fn   string
	args []node
}

func (c call) eval(parent float64) float64 {
	out := c.args[0].eval(parent)
	for _, a := range c.args[1:] {
		v := a.eval(parent)
		if c.fn == "min" {
			out = math.Min(out, v)
		} else {
			out = math.Max(out, v)
		}
	}
	return out
}

func (c call) goSource(parent string) string {
	args := make([]string, len(c.args))
	for i, a := range c.args {
		args[i] = a.goSource(parent)
	}
	return c.fn + "(" + strings.Join(args, ", ") + ")"
}

func (c call) usesParent() bool {
	for _, a := range c.args {
		if a.usesParent() {
			return true
		}
	}
	return false
}

func clean(s string) string { return strings.Join(strings.Fields(s), " ") }
