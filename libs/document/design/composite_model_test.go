package design

import "testing"

func card() Composite {
	return Composite{
		ID: "card", Name: "Card", W: 20, H: 5,
		Nodes: []Node{
			{ID: "a", Component: "lipgloss.box", Name: "Frame", Rect: Rect{W: 20, H: 5}},
			{ID: "b", Component: "lipgloss.label", Name: "Title", Rect: Rect{X: 2, Y: 1, W: 10, H: 1}},
		},
		Props: []Exposed{{Key: "title", Label: "Title", Kind: "text", Default: "Hi", Target: "b", TargetProp: "text"}},
	}
}

func TestAGoodCompositePasses(t *testing.T) {
	if err := card().Check(); err != nil {
		t.Fatal(err)
	}
}

func TestCompositeProblemsAreNamed(t *testing.T) {
	cases := map[string]func(*Composite){
		"bad id":        func(c *Composite) { c.ID = "Not Valid" },
		"no name":       func(c *Composite) { c.Name = "" },
		"no size":       func(c *Composite) { c.W = 0 },
		"repeated part": func(c *Composite) { c.Nodes[1].ID = "a" },
		"tiny part":     func(c *Composite) { c.Nodes[0].Rect.W = 0 },
		"unknown kind":  func(c *Composite) { c.Props[0].Kind = "sparkle" },
		"dangling prop": func(c *Composite) { c.Props[0].Target = "zzz" },
		"repeated key":  func(c *Composite) { c.Props = append(c.Props, c.Props[0]) },
	}
	for name, mutate := range cases {
		c := card().Clone()
		mutate(&c)
		if c.Check() == nil {
			t.Errorf("%s should be refused", name)
		}
	}
}

func TestCloneIsIndependent(t *testing.T) {
	a := card()
	b := a.Clone()
	b.Nodes[0].Name = "changed"
	b.Props[0].Label = "changed"
	if a.Nodes[0].Name == "changed" || a.Props[0].Label == "changed" {
		t.Fatal("clone shares state")
	}
}
