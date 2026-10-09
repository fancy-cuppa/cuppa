package scheme

import (
	"sort"
	"strings"
	"testing"
)

func TestAllListsHundredsOfSchemesSortedByName(t *testing.T) {
	list := All()
	if len(list) < 300 {
		t.Fatalf("%d schemes, want hundreds", len(list))
	}
	if !sort.SliceIsSorted(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	}) {
		t.Error("schemes are not sorted by name")
	}
	for _, s := range list {
		if s.Name == "" {
			t.Fatal("a scheme has no name")
		}
	}
}

func TestFindIgnoresCaseAndReturnsTheSixteenColours(t *testing.T) {
	s, ok := Find("dracula")
	if !ok {
		t.Fatal("Dracula should be there")
	}
	if !s.Dark {
		t.Error("Dracula is a dark scheme")
	}
	if s.Background == s.Foreground {
		t.Error("foreground and background should differ")
	}
	distinct := map[string]bool{}
	for _, c := range s.ANSI {
		distinct[c.Hex()] = true
	}
	if len(distinct) < 8 {
		t.Errorf("only %d distinct ANSI colours", len(distinct))
	}
	if _, ok := Find("no such scheme"); ok {
		t.Error("an unknown name should not be found")
	}
}

func TestThemeTakesPageTextAndThreeAccentsFromAScheme(t *testing.T) {
	s, ok := Find("dracula")
	if !ok {
		t.Fatal("Dracula should be there")
	}
	th := s.Theme()
	if th.Background != s.Background || th.Text != s.Foreground {
		t.Error("page and text come straight from the scheme")
	}
	if th.Border != s.ANSI[4] || th.Secondary != s.ANSI[5] {
		t.Error("border is blue and secondary is purple")
	}
	if th.Muted == th.Background || th.Muted == th.Text {
		t.Errorf("muted must differ from both: %v", th.Muted)
	}
	for _, sc := range All() {
		m := sc.Theme()
		if m.Muted == m.Background && m.Background == m.Text {
			t.Fatalf("%s has no usable muted colour", sc.Name)
		}
	}
}
