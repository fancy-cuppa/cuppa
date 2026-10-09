package registry

import (
	"testing"

	"github.com/fancy-cuppa/cuppa/libs/catalog/definition"
)

func sample() []definition.Definition {
	return []definition.Definition{
		{ID: "bubbles.spinner", Name: "Spinner", Family: definition.FamilyBubbles},
		{ID: "lipgloss.label", Name: "Label", Family: definition.FamilyLipgloss},
		{ID: "lipgloss.box", Name: "Box", Family: definition.FamilyLipgloss, Description: "bordered"},
	}
}

func TestNewRejectsDuplicatesAndEmptyIDs(t *testing.T) {
	if _, err := New(definition.Definition{Name: "x"}); err == nil {
		t.Fatal("empty id accepted")
	}
	d := definition.Definition{ID: "a"}
	if _, err := New(d, d); err == nil {
		t.Fatal("duplicate id accepted")
	}
}

func TestListOrdersByFamilyThenName(t *testing.T) {
	r, err := New(sample()...)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, d := range r.List() {
		got = append(got, d.ID)
	}
	want := []string{"lipgloss.box", "lipgloss.label", "bubbles.spinner"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v, want %v", got, want)
		}
	}
	if fs := r.Families(); len(fs) != 2 || fs[0] != definition.FamilyLipgloss {
		t.Fatalf("families = %v", fs)
	}
}

func TestSearchAndGet(t *testing.T) {
	r, _ := New(sample()...)
	if got := r.Search("BORDER"); len(got) != 1 || got[0].ID != "lipgloss.box" {
		t.Fatalf("search = %v", got)
	}
	if got := r.Search(""); len(got) != 3 {
		t.Fatalf("empty search = %d", len(got))
	}
	if _, ok := r.Get("nope"); ok {
		t.Fatal("unknown id found")
	}
}
