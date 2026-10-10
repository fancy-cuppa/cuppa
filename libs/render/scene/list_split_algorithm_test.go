package scene

import (
	"reflect"
	"testing"
)

func TestListItemsMayHoldAnEscapedComma(t *testing.T) {
	bs := string(rune(92))
	p := Props{"items": "Highlight (keys" + bs + ", bars), Dim", "rows": "a" + bs + ",b, c; d"}
	if got, want := p.List("items"), []string{"Highlight (keys, bars)", "Dim"}; !reflect.DeepEqual(got, want) {
		t.Errorf("items = %q, want %q", got, want)
	}
	rows := splitRaw(p["rows"], ";")
	if len(rows) != 2 || !reflect.DeepEqual(splitList(rows[0], ","), []string{"a,b", "c"}) {
		t.Errorf("rows = %q", rows)
	}
}

func TestListEnumeratorNoneHasNoMarker(t *testing.T) {
	if marker("none", 0) != "" || marker("bullet", 0) == "" {
		t.Error("marker")
	}
}
