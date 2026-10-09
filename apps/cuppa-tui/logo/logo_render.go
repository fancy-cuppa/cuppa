// Package logo draws the Cuppa logo in terminal cells.
//
// It does what an image viewer such as chafa does, with a small symbol set that
// every common font has: each cell shows a 2 by 2 block of the picture, drawn
// with one of the quadrant characters (▘ ▝ ▀ ▖ ▌ ▞ ▛ ▗ ▚ ▐ ▜ ▄ ▙ ▟ █) and two true
// colours, the pair that is closest to the four pixels.
package logo

import (
	"bytes"
	_ "embed" // the logo picture
	"fmt"
	"image"
	"image/png"
	"strings"
	"sync"
)

//go:embed cuppa.png
var pngData []byte

var (
	loadOnce sync.Once
	picture  image.Image
)

func source() image.Image {
	loadOnce.Do(func() {
		img, err := png.Decode(bytes.NewReader(pngData))
		if err != nil {
			img = image.NewRGBA(image.Rect(0, 0, 1, 1))
		}
		picture = img
	})
	return picture
}

// quadrant[mask] is the character that fills the quadrants in mask: 1 is the
// top left, 2 the top right, 4 the bottom left and 8 the bottom right.
var quadrant = [16]rune{' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛', '▗', '▚', '▐', '▜', '▄', '▙', '▟', '█'}

type rgb struct{ r, g, b float64 }

func (c rgb) sub(o rgb) float64 {
	dr, dg, db := c.r-o.r, c.g-o.g, c.b-o.b
	return dr*dr + dg*dg + db*db
}

var (
	cacheMu sync.Mutex
	cache   = map[string][]string{}
)

// Render draws the logo in cols by rows cells and returns one string per row,
// with true-colour escape codes. Cells are about twice as tall as wide, so
// cols = 2*rows keeps the picture square. The background shows through where
// the logo is transparent; pass a dark colour for a dark terminal.
func Render(cols, rows int, background [3]uint8) []string {
	if cols < 1 || rows < 1 {
		return nil
	}
	key := fmt.Sprintf("%dx%d/%v", cols, rows, background)
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if lines, ok := cache[key]; ok {
		return lines
	}
	px := sample(cols*2, rows*2, rgb{float64(background[0]), float64(background[1]), float64(background[2])})
	lines := make([]string, rows)
	for y := 0; y < rows; y++ {
		var b strings.Builder
		for x := 0; x < cols; x++ {
			cell := [4]rgb{
				px[2*y][2*x], px[2*y][2*x+1],
				px[2*y+1][2*x], px[2*y+1][2*x+1],
			}
			mask, fg, bg := bestSplit(cell)
			fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm%c",
				round(fg.r), round(fg.g), round(fg.b), round(bg.r), round(bg.g), round(bg.b), quadrant[mask])
		}
		b.WriteString("\x1b[0m")
		lines[y] = b.String()
	}
	cache[key] = lines
	return lines
}

// bestSplit picks the quadrant pattern whose two average colours are closest
// to the four pixels, and returns the pattern and those colours.
func bestSplit(p [4]rgb) (mask int, fg, bg rgb) {
	// A cell that is (nearly) one colour is a plain space on that colour.
	avg := rgb{(p[0].r + p[1].r + p[2].r + p[3].r) / 4, (p[0].g + p[1].g + p[2].g + p[3].g) / 4, (p[0].b + p[1].b + p[2].b + p[3].b) / 4}
	var spread float64
	for i := 0; i < 4; i++ {
		spread += p[i].sub(avg)
	}
	if spread < 12 {
		return 0, avg, avg
	}
	best := -1.0
	for m := 1; m < 15; m++ {
		var in, out rgb
		var nIn, nOut float64
		for i := 0; i < 4; i++ {
			if m&(1<<i) != 0 {
				in = rgb{in.r + p[i].r, in.g + p[i].g, in.b + p[i].b}
				nIn++
			} else {
				out = rgb{out.r + p[i].r, out.g + p[i].g, out.b + p[i].b}
				nOut++
			}
		}
		in = rgb{in.r / nIn, in.g / nIn, in.b / nIn}
		out = rgb{out.r / nOut, out.g / nOut, out.b / nOut}
		var err float64
		for i := 0; i < 4; i++ {
			if m&(1<<i) != 0 {
				err += p[i].sub(in)
			} else {
				err += p[i].sub(out)
			}
		}
		if best < 0 || err < best {
			best, mask, fg, bg = err, m, in, out
		}
	}
	return mask, fg, bg
}

// sample shrinks the picture to w by h pixels. Each target pixel is the average
// of the source pixels under it, over the background where they are
// transparent, weighted towards the bright ones: the logo is bright blocks with
// thin dark gaps, and a plain average would grey every block at this size.
func sample(w, h int, background rgb) [][]rgb {
	img := source()
	b := img.Bounds()
	out := make([][]rgb, h)
	for ty := 0; ty < h; ty++ {
		out[ty] = make([]rgb, w)
		y0 := b.Min.Y + ty*b.Dy()/h
		y1 := max(b.Min.Y+(ty+1)*b.Dy()/h, y0+1)
		for tx := 0; tx < w; tx++ {
			x0 := b.Min.X + tx*b.Dx()/w
			x1 := max(b.Min.X+(tx+1)*b.Dx()/w, x0+1)
			var r, g, bl, weight float64
			for y := y0; y < y1; y++ {
				for x := x0; x < x1; x++ {
					pr, pg, pb, pa := img.At(x, y).RGBA() // premultiplied, 16 bits
					rest := 1 - float64(pa)/65535
					cr := float64(pr>>8) + background.r*rest
					cg := float64(pg>>8) + background.g*rest
					cb := float64(pb>>8) + background.b*rest
					luma := (0.299*cr + 0.587*cg + 0.114*cb) / 255
					wgt := 0.03 + luma*luma
					r += cr * wgt
					g += cg * wgt
					bl += cb * wgt
					weight += wgt
				}
			}
			out[ty][tx] = rgb{r / weight, g / weight, bl / weight}
		}
	}
	return out
}

func round(v float64) int { return min(max(int(v+0.5), 0), 255) }

// Dialog is the logo at the size the About and Welcome dialogs show it: 20 by
// 10 cells on black.
func Dialog() []string { return Render(20, 10, [3]uint8{0, 0, 0}) }
