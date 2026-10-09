package design

// Rect is a rectangle in terminal cells. X and Y are the top-left cell.
type Rect struct {
	X, Y, W, H int
}

// Right is the first column outside the rectangle.
func (r Rect) Right() int { return r.X + r.W }

// Bottom is the first row outside the rectangle.
func (r Rect) Bottom() int { return r.Y + r.H }

// Contains reports whether the cell (x, y) lies inside the rectangle.
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.Right() && y >= r.Y && y < r.Bottom()
}

// Intersects reports whether the two rectangles share at least one cell.
func (r Rect) Intersects(o Rect) bool {
	return r.X < o.Right() && o.X < r.Right() && r.Y < o.Bottom() && o.Y < r.Bottom()
}

// Translate returns the rectangle moved by (dx, dy).
func (r Rect) Translate(dx, dy int) Rect {
	r.X += dx
	r.Y += dy
	return r
}

// Scale maps r from a box of size fromW x fromH onto one of toW x toH, keeping
// every edge on a whole cell and the result at least 1 x 1.
func (r Rect) Scale(fromW, fromH, toW, toH int) Rect {
	if fromW < 1 || fromH < 1 {
		return r
	}
	x0, x1 := r.X*toW/fromW, (r.X+r.W)*toW/fromW
	y0, y1 := r.Y*toH/fromH, (r.Y+r.H)*toH/fromH
	return Rect{X: x0, Y: y0, W: max(x1-x0, 1), H: max(y1-y0, 1)}
}

// MoveInto returns the rectangle shifted the least needed to lie inside
// bounds. A rectangle larger than bounds is aligned to the top-left of bounds.
func (r Rect) MoveInto(bounds Rect) Rect {
	if r.Right() > bounds.Right() {
		r.X = bounds.Right() - r.W
	}
	if r.Bottom() > bounds.Bottom() {
		r.Y = bounds.Bottom() - r.H
	}
	if r.X < bounds.X {
		r.X = bounds.X
	}
	if r.Y < bounds.Y {
		r.Y = bounds.Y
	}
	return r
}

// Union is the smallest rectangle covering both.
func (r Rect) Union(o Rect) Rect {
	x, y := min(r.X, o.X), min(r.Y, o.Y)
	return Rect{X: x, Y: y, W: max(r.Right(), o.Right()) - x, H: max(r.Bottom(), o.Bottom()) - y}
}
