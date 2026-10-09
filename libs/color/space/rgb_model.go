// Package space is colour in the forms the designer offers: the 16 and 256
// ANSI palettes, hex and RGB, and HSL. It converts between them and parses and
// normalises the text a design stores.
package space

import "fmt"

// RGB is a colour as red, green and blue, 0 to 255 each.
type RGB struct{ R, G, B uint8 }

// Hex is the canonical "#rrggbb" form.
func (c RGB) Hex() string { return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B) }

// Luma is the perceived brightness, 0 (black) to 255 (white); it decides
// whether dark or light text reads better on top of the colour.
func (c RGB) Luma() float64 {
	return 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
}
