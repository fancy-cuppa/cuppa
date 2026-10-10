package editor

import "github.com/meta-tui/cuppa/libs/document/design"

// MoveSelectionBy shifts every selected node by (dx, dy), as a group, stopping
// at the canvas edge. It does not checkpoint: drags call Checkpoint once first.
func (e *Editor) MoveSelectionBy(dx, dy int) bool {
	ids := e.unlockedSelection()
	if len(ids) == 0 {
		return false
	}
	var union design.Rect
	for i, id := range ids {
		n, _ := e.doc.Get(id)
		if i == 0 {
			union = n.Rect
		} else {
			union = union.Union(n.Rect)
		}
	}
	moved := union.Translate(dx, dy).MoveInto(e.doc.Bounds())
	dx, dy = moved.X-union.X, moved.Y-union.Y
	if dx == 0 && dy == 0 {
		return false
	}
	for _, id := range ids {
		before, _ := e.doc.Get(id)
		e.doc.Update(id, func(n *design.Node) { n.Rect = n.Rect.Translate(dx, dy) })
		e.followRect(id, before.Rect)
	}
	return true
}

// Nudge moves the selection by (dx, dy) as one undo step, for keyboard moves.
// It leaves no step when nothing moved (locked, or already at the canvas edge).
func (e *Editor) Nudge(dx, dy int) bool {
	return e.apply(func() bool { return e.MoveSelectionBy(dx, dy) })
}

// MoveTo puts the node's top-left at (x, y), clamped to the canvas, as one undo step.
func (e *Editor) MoveTo(id design.NodeID, x, y int) bool {
	n, ok := e.doc.Get(id)
	if !ok {
		return false
	}
	return e.SetRect(id, design.Rect{X: x, Y: y, W: n.Rect.W, H: n.Rect.H}, true)
}

// SetRect gives the node a new rectangle, enforcing the component's minimum
// size and the canvas bounds. With checkpoint false it joins the gesture
// already in progress.
func (e *Editor) SetRect(id design.NodeID, r design.Rect, checkpoint bool) bool {
	n, ok := e.doc.Get(id)
	if !ok || n.Locked {
		return false
	}
	minW, minH := e.minSize(n)
	r.W = min(max(r.W, minW), e.doc.Width)
	r.H = min(max(r.H, minH), e.doc.Height)
	r = r.MoveInto(e.doc.Bounds())
	if r == n.Rect {
		return false
	}
	set := func() bool {
		e.doc.Update(id, func(n *design.Node) { n.Rect = r })
		e.followRect(id, n.Rect)
		return true
	}
	if checkpoint {
		return e.apply(set)
	}
	return set()
}
