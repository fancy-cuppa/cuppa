package space

import (
	"math"
	"testing"
)

func TestANSIPalette(t *testing.T) {
	cases := map[int]RGB{
		0: {0, 0, 0}, 9: {255, 0, 0}, 15: {255, 255, 255},
		16: {0, 0, 0}, 21: {0, 0, 255}, 196: {255, 0, 0}, 212: {255, 135, 215}, 231: {255, 255, 255},
		232: {8, 8, 8}, 255: {238, 238, 238},
		-5: {0, 0, 0}, 999: {238, 238, 238},
	}
	for i, want := range cases {
		if got := ANSI(i); got != want {
			t.Errorf("ANSI(%d) = %v, want %v", i, got, want)
		}
	}
}

func TestNearestFindsTheCubeColour(t *testing.T) {
	if got := Nearest(RGB{255, 135, 215}); got != 212 {
		t.Fatalf("got %d, want 212", got)
	}
	if got := Nearest(RGB{10, 10, 10}); got != 232 {
		t.Fatalf("near-black should be the first grey, got %d", got)
	}
}

func TestHSLRoundTrip(t *testing.T) {
	for _, c := range []RGB{{0, 0, 0}, {255, 255, 255}, {255, 0, 0}, {18, 200, 77}, {99, 99, 99}, {255, 135, 215}, {1, 2, 254}} {
		back := ToHSL(c).RGB()
		if math.Abs(float64(back.R)-float64(c.R)) > 1 || math.Abs(float64(back.G)-float64(c.G)) > 1 || math.Abs(float64(back.B)-float64(c.B)) > 1 {
			t.Errorf("%v -> %+v -> %v", c, ToHSL(c), back)
		}
	}
}

func TestHSLKnownValues(t *testing.T) {
	if got := (HSL{H: 120, S: 1, L: 0.5}).RGB(); got != (RGB{0, 255, 0}) {
		t.Fatalf("pure green: %v", got)
	}
	if h := ToHSL(RGB{0, 0, 255}); h.H != 240 || h.S != 1 || h.L != 0.5 {
		t.Fatalf("pure blue: %+v", h)
	}
	if got := (HSL{H: -60, S: 1, L: 0.5}).RGB(); got != (RGB{255, 0, 255}) {
		t.Fatalf("negative hue wraps: %v", got)
	}
	if got := (HSL{H: 10, S: 5, L: 2}).RGB(); got != (RGB{255, 255, 255}) {
		t.Fatalf("out of range values clamp: %v", got)
	}
}

func TestNormalise(t *testing.T) {
	ok := map[string]string{
		"": "", "  ": "", "212": "212", " 7 ": "7", "007": "7",
		"#FF5FD7": "#ff5fd7", "#f5d": "#ff55dd", "#000": "#000000",
	}
	for in, want := range ok {
		if got, err := Normalise(in); err != nil || got != want {
			t.Errorf("Normalise(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"256", "-1", "red", "#12", "#12345", "#gggggg", "0x10", "1.5"} {
		if got, err := Normalise(in); err == nil {
			t.Errorf("Normalise(%q) = %q, want an error", in, got)
		}
	}
}

func TestResolve(t *testing.T) {
	if c, ok := Resolve("#0000ff"); !ok || c != (RGB{0, 0, 255}) {
		t.Fatalf("hex: %v %v", c, ok)
	}
	if c, ok := Resolve("21"); !ok || c != (RGB{0, 0, 255}) {
		t.Fatalf("index: %v %v", c, ok)
	}
	if _, ok := Resolve(""); ok {
		t.Fatal("empty is no colour")
	}
	if _, ok := Resolve("nope"); ok {
		t.Fatal("garbage is no colour")
	}
}

func TestLuma(t *testing.T) {
	if (RGB{255, 255, 255}).Luma() < 250 || (RGB{0, 0, 0}).Luma() != 0 || (RGB{0, 0, 255}).Luma() > 40 {
		t.Fatal("luma should rank white above blue above black")
	}
}
