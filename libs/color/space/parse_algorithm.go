package space

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse reads a colour as a design stores it: a palette index ("212"), or hex
// ("#ff5fd7" or "#f5d"). It reports the colour and, for an index, which one.
// The empty string is valid and means "no colour" (ok is true, isIndex false,
// and the RGB is the zero value); use Normalise to tell them apart.
func Parse(s string) (c RGB, index int, isIndex, ok bool) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "#") {
		c, ok = parseHex(s[1:])
		return c, -1, false, ok
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n > 255 {
		return RGB{}, -1, false, false
	}
	return ANSI(n), n, true, true
}

func parseHex(h string) (RGB, bool) {
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) != 6 {
		return RGB{}, false
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return RGB{}, false
	}
	return RGB{uint8(v >> 16), uint8(v >> 8), uint8(v)}, true
}

// Normalise returns the canonical text of a colour: "" stays "", an index is
// kept as its number, and hex becomes lowercase "#rrggbb". It rejects anything
// else, so a design can only ever hold a colour a terminal understands.
func Normalise(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}
	c, index, isIndex, ok := Parse(s)
	switch {
	case !ok:
		return "", fmt.Errorf("%q is not a colour: use a number from 0 to 255 or hex like #ff5fd7", s)
	case isIndex:
		return strconv.Itoa(index), nil
	}
	return c.Hex(), nil
}

// Resolve is the RGB of a stored colour, or false when it is empty or invalid.
func Resolve(s string) (RGB, bool) {
	c, _, _, ok := Parse(s)
	if !ok || strings.TrimSpace(s) == "" {
		return RGB{}, false
	}
	return c, true
}
