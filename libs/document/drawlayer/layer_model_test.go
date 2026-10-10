package drawlayer

import "testing"

func TestEncodeDecodeRoundTrip(t *testing.T) {
	l := New()
	l.Paint(Cell{X: 3, Y: 1, Ch: '╭', Fg: "#ff87d7"}, Cell{X: 0, Y: 0, Ch: ';', Fg: "212"}, Cell{X: 1, Y: 0, Ch: ' ', Bg: "#101010"})
	back := Decode(l.Encode())
	if back.Len() != 3 {
		t.Fatalf("len = %d, want 3", back.Len())
	}
	if c, _ := back.Get(0, 0); c.Ch != ';' || c.Fg != "212" {
		t.Fatalf("cell (0,0) = %+v: a semicolon must survive", c)
	}
	if c, _ := back.Get(1, 0); c.Bg != "#101010" {
		t.Fatalf("cell (1,0) = %+v", c)
	}
}

func TestEraseRemovesOnlyThosePositions(t *testing.T) {
	l := New()
	l.Paint(Cell{X: 0, Y: 0, Ch: 'a'}, Cell{X: 1, Y: 0, Ch: 'b'})
	l.Erase(Cell{X: 0, Y: 0})
	if _, ok := l.Get(0, 0); ok {
		t.Fatal("erased cell is still there")
	}
	if _, ok := l.Get(1, 0); !ok {
		t.Fatal("erase took a neighbour")
	}
}

func TestPaintingNothingIsIgnoredAndBadEntriesAreSkipped(t *testing.T) {
	l := New()
	l.Paint(Cell{X: 0, Y: 0})
	if l.Len() != 0 {
		t.Fatal("an empty cell was painted")
	}
	if Decode("junk;1,2,zz,,;3,4,41,,").Len() != 1 {
		t.Fatal("only the valid entry should survive")
	}
}
