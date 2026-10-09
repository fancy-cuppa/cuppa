// Package pointer is the mouse input every pane understands, independent of
// the terminal library that produced it.
package pointer

// Phase is what the pointer did.
type Phase int

// Pointer phases.
const (
	// Down is a button press.
	Down Phase = iota
	// Move is motion, with or without a button held (see Event.Held).
	Move
	// Up is a button release.
	Up
	// Wheel is a scroll step; see Event.WheelX and Event.WheelY.
	Wheel
)

// Event is one pointer event in the coordinates of whoever receives it.
type Event struct {
	X, Y  int
	Phase Phase
	// Left is true for the primary button on Down and Up.
	Left bool
	// Held is true while the primary button is down (drags).
	Held bool
	// Shift is true while the shift key is held.
	Shift bool
	// WheelX and WheelY are -1, 0 or 1: the scroll direction on a Wheel event.
	WheelX, WheelY int
}

// Translate returns the event with its position moved by (-dx, -dy), turning
// screen coordinates into the coordinates of a pane whose top-left is (dx, dy).
func (e Event) Translate(dx, dy int) Event {
	e.X -= dx
	e.Y -= dy
	return e
}
