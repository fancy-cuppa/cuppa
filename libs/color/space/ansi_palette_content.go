package space

// system16 are the 16 standard terminal colours (xterm values). Terminals let
// users change them, so these are the usual look, not a promise.
var system16 = [16]RGB{
	{0, 0, 0}, {205, 0, 0}, {0, 205, 0}, {205, 205, 0},
	{0, 0, 238}, {205, 0, 205}, {0, 205, 205}, {229, 229, 229},
	{127, 127, 127}, {255, 0, 0}, {0, 255, 0}, {255, 255, 0},
	{92, 92, 255}, {255, 0, 255}, {0, 255, 255}, {255, 255, 255},
}

// cubeLevels are the six steps of each channel in the 6x6x6 colour cube.
var cubeLevels = [6]uint8{0, 95, 135, 175, 215, 255}

// ANSI is the colour of an index in the 256-colour palette: 0 to 15 the system
// colours, 16 to 231 the colour cube, 232 to 255 the grey ramp. Out-of-range
// indexes are clamped.
func ANSI(index int) RGB {
	index = min(max(index, 0), 255)
	switch {
	case index < 16:
		return system16[index]
	case index < 232:
		i := index - 16
		return RGB{cubeLevels[i/36], cubeLevels[(i/6)%6], cubeLevels[i%6]}
	}
	v := uint8(8 + 10*(index-232))
	return RGB{v, v, v}
}

// Nearest is the index in the 256-colour palette closest to c (ignoring the
// 16 system colours, which terminals may redefine), for showing a hex colour
// on a terminal that cannot display it.
func Nearest(c RGB) int {
	best, bestDist := 16, int(^uint(0)>>1)
	for i := 16; i < 256; i++ {
		p := ANSI(i)
		dr, dg, db := int(p.R)-int(c.R), int(p.G)-int(c.G), int(p.B)-int(c.B)
		if d := dr*dr + dg*dg + db*db; d < bestDist {
			best, bestDist = i, d
		}
	}
	return best
}
