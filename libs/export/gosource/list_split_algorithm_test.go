package gosource

import (
	"reflect"
	"testing"
)

func TestSplitListHonoursEscapes(t *testing.T) {
	bs := string(rune(92))
	for _, c := range []struct {
		in, sep string
		want    []string
	}{
		{"a, b ,, c", ",", []string{"a", "b", "c"}},
		{"Highlight (keys" + bs + ", bars), Dim", ",", []string{"Highlight (keys, bars)", "Dim"}},
		{"a" + bs + bs + ", b", ",", []string{"a" + bs, "b"}},
		{"C:" + bs + "dir, x", ",", []string{"C:" + bs + "dir", "x"}},
		{"a" + bs + ";b;c", ";", []string{"a;b", "c"}},
	} {
		if got := splitList(c.in, c.sep); !reflect.DeepEqual(got, c.want) {
			t.Errorf("splitList(%q, %q) = %q, want %q", c.in, c.sep, got, c.want)
		}
	}
	// A table is cut in two steps: the first keeps the escapes for the second.
	rows := splitRaw("a"+bs+","+bs+";b, c; d", ";")
	if len(rows) != 2 || splitList(rows[0], ",")[0] != "a,;b" {
		t.Errorf("rows = %q", rows)
	}
}
