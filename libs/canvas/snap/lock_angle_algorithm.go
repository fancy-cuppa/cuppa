package snap

import "math"

// cellAspect is how many times taller a terminal cell is than it is wide, so a
// diagonal that looks like 45° is two columns for every row.
const cellAspect = 2

// LockAngle limits a drag of (dx, dy) cells to the nearest of horizontal,
// vertical and the two diagonals, as they look on screen (a terminal cell is
// about twice as tall as wide), keeping as much of the drag as the direction
// allows.
func LockAngle(dx, dy int) (int, int) {
	if dx == 0 && dy == 0 {
		return 0, 0
	}
	vx, vy := float64(dx), float64(dy*cellAspect)
	angle := math.Atan2(vy, vx)
	step := math.Pi / 4
	locked := math.Round(angle/step) * step
	ux, uy := math.Cos(locked), math.Sin(locked)
	t := vx*ux + vy*uy
	return int(math.Round(t * ux)), int(math.Round(t * uy / cellAspect))
}
