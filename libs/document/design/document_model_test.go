package design

import "testing"

func ids(d Document) []NodeID {
	out := make([]NodeID, len(d.Nodes))
	for i, n := range d.Nodes {
		out[i] = n.ID
	}
	return out
}

func newThree() Document {
	d := NewDocument("t", 80, 24)
	d.Add(Node{Component: "a"})
	d.Add(Node{Component: "b"})
	d.Add(Node{Component: "c"})
	return d
}

func TestAddAssignsUniqueIDsOnTop(t *testing.T) {
	d := newThree()
	got := ids(d)
	want := []NodeID{"n1", "n2", "n3"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ids = %v, want %v", got, want)
		}
	}
	d.Remove("n3")
	if n := d.Add(Node{}); n.ID != "n4" {
		t.Fatalf("id reused after remove: %s", n.ID)
	}
}

func TestZOrder(t *testing.T) {
	d := newThree()
	if !d.ToFront("n1") || ids(d)[2] != "n1" {
		t.Fatalf("ToFront: %v", ids(d))
	}
	if !d.ToBack("n1") || ids(d)[0] != "n1" {
		t.Fatalf("ToBack: %v", ids(d))
	}
	if !d.Raise("n1") || ids(d)[1] != "n1" {
		t.Fatalf("Raise: %v", ids(d))
	}
	if !d.Lower("n1") || ids(d)[0] != "n1" {
		t.Fatalf("Lower: %v", ids(d))
	}
	if d.Lower("n1") {
		t.Fatal("Lower at the back should report no change")
	}
	if d.Raise("missing") {
		t.Fatal("unknown id should report no change")
	}
}

func TestCloneIsDeep(t *testing.T) {
	d := newThree()
	d.Update("n1", func(n *Node) { n.Props = map[string]string{"k": "v"} })
	c := d.Clone()
	c.Update("n1", func(n *Node) { n.Props["k"] = "changed" })
	c.Remove("n2")
	if n, _ := d.Get("n1"); n.Props["k"] != "v" {
		t.Fatal("clone shares props with the original")
	}
	if len(d.Nodes) != 3 {
		t.Fatal("clone shares the node slice")
	}
}

func TestRectGeometry(t *testing.T) {
	r := Rect{X: 2, Y: 2, W: 4, H: 3}
	if !r.Contains(2, 2) || !r.Contains(5, 4) || r.Contains(6, 4) || r.Contains(5, 5) {
		t.Fatal("Contains edges wrong")
	}
	if !r.Intersects(Rect{X: 5, Y: 4, W: 3, H: 3}) || r.Intersects(Rect{X: 6, Y: 2, W: 2, H: 2}) {
		t.Fatal("Intersects wrong")
	}
	got := Rect{X: 78, Y: -3, W: 10, H: 4}.MoveInto(Rect{W: 80, H: 24})
	if got != (Rect{X: 70, Y: 0, W: 10, H: 4}) {
		t.Fatalf("MoveInto = %+v", got)
	}
	if u := r.Union(Rect{X: 0, Y: 0, W: 1, H: 1}); u != (Rect{X: 0, Y: 0, W: 6, H: 5}) {
		t.Fatalf("Union = %+v", u)
	}
}

func TestMoveToIndexPlacesANodeExactly(t *testing.T) {
	d := NewDocument("t", 10, 10)
	a := d.Add(Node{Name: "a", Rect: Rect{W: 1, H: 1}})
	b := d.Add(Node{Name: "b", Rect: Rect{W: 1, H: 1}})
	c := d.Add(Node{Name: "c", Rect: Rect{W: 1, H: 1}})
	if !d.MoveToIndex(a.ID, 2) || d.Nodes[0].ID != b.ID || d.Nodes[1].ID != c.ID || d.Nodes[2].ID != a.ID {
		t.Fatalf("order %v", d.Nodes)
	}
	if d.MoveToIndex(a.ID, 2) || d.MoveToIndex(a.ID, 9) || d.MoveToIndex("zz", 0) {
		t.Fatal("no-ops must report false")
	}
}
