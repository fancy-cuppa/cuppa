package space

import "math"

// HSL is a colour as hue (degrees, 0 up to 360), saturation and lightness
// (each 0 to 1).
type HSL struct{ H, S, L float64 }

// ToHSL converts from RGB.
func ToHSL(c RGB) HSL {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	hi, lo := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	l := (hi + lo) / 2
	if hi == lo {
		return HSL{L: l}
	}
	d := hi - lo
	s := d / (1 - math.Abs(2*l-1))
	var h float64
	switch hi {
	case r:
		h = math.Mod((g-b)/d, 6)
	case g:
		h = (b-r)/d + 2
	default:
		h = (r-g)/d + 4
	}
	h *= 60
	if h < 0 {
		h += 360
	}
	return HSL{H: h, S: s, L: l}
}

// RGB converts to RGB, clamping out-of-range values.
func (c HSL) RGB() RGB {
	h := math.Mod(c.H, 360)
	if h < 0 {
		h += 360
	}
	s, l := clamp01(c.S), clamp01(c.L)
	chroma := (1 - math.Abs(2*l-1)) * s
	x := chroma * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - chroma/2
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = chroma, x, 0
	case h < 120:
		r, g, b = x, chroma, 0
	case h < 180:
		r, g, b = 0, chroma, x
	case h < 240:
		r, g, b = 0, x, chroma
	case h < 300:
		r, g, b = x, 0, chroma
	default:
		r, g, b = chroma, 0, x
	}
	return RGB{to255(r + m), to255(g + m), to255(b + m)}
}

func clamp01(v float64) float64 { return math.Min(math.Max(v, 0), 1) }

func to255(v float64) uint8 { return uint8(math.Round(clamp01(v) * 255)) }
