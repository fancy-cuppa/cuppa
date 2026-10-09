package format

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/meta-tui/cuppa/libs/document/design"
)

func sample() design.Document {
	d := design.NewDocument("Login", 80, 24)
	d.Add(design.Node{Component: "bubbles.list", Name: "List", Rect: design.Rect{X: 1, Y: 2, W: 20, H: 8}, Props: map[string]string{"title": "Tea"}})
	d.Add(design.Node{Component: "huh.input", Name: "Input", Rect: design.Rect{X: 30, Y: 4, W: 24, H: 3}})
	return d
}

func TestRoundTrip(t *testing.T) {
	want := sample()
	data, err := Encode(want)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte(Magic)) {
		t.Fatalf("missing magic: %q", data[:8])
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip changed the document:\n got %+v\nwant %+v", got, want)
	}
}

func TestDecodeRejects(t *testing.T) {
	good, _ := Encode(sample())
	newer := append([]byte(nil), good...)
	newer[len(Magic)+1] = byte(CurrentVersion + 1)
	zero := append([]byte(nil), good...)
	zero[len(Magic)], zero[len(Magic)+1] = 0, 0
	cases := map[string]struct {
		data []byte
		want error
	}{
		"empty":      {nil, ErrNotCuppa},
		"json":       {[]byte(`{"document":{}}`), ErrNotCuppa},
		"no version": {[]byte(Magic), ErrCorrupt},
		"version 0":  {zero, ErrCorrupt},
		"newer":      {newer, ErrTooNew},
		"bad gzip":   {append([]byte(Magic+"\x00\x01"), "junk"...), ErrCorrupt},
		"truncated":  {good[:len(good)/2], ErrCorrupt},
	}
	for name, c := range cases {
		if _, err := Decode(c.data); !errors.Is(err, c.want) {
			t.Errorf("%s: got %v, want %v", name, err, c.want)
		}
	}
}

func TestDecodeRepairsImpossibleNodes(t *testing.T) {
	doc := design.Document{Width: 0, Height: -3, Nodes: []design.Node{
		{ID: "a", Rect: design.Rect{W: 5, H: 5}},
		{ID: "a", Rect: design.Rect{W: 5, H: 5}},
		{ID: "", Rect: design.Rect{W: 5, H: 5}},
		{ID: "b", Rect: design.Rect{W: 0, H: 5}},
	}}
	data, _ := Encode(doc)
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.Width != 1 || got.Height != 1 || len(got.Nodes) != 1 {
		t.Fatalf("not repaired: %+v", got)
	}
}

func TestMigrationsRunInOrder(t *testing.T) {
	add := func(field string) migration {
		return func(b rawEnvelope) (rawEnvelope, error) {
			return rawEnvelope(string(b[:len(b)-1]) + `,"` + field + `":1}`), nil
		}
	}
	table := map[uint16]migration{1: add("a"), 2: add("b")}
	got, err := migrateWith(table, rawEnvelope(`{"x":0}`), 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]int
	if err := json.Unmarshal(got, &m); err != nil || m["a"] != 1 || m["b"] != 1 {
		t.Fatalf("got %s (%v)", got, err)
	}
	if _, err := migrateWith(table, rawEnvelope(`{}`), 1, 4); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("missing step: got %v", err)
	}
}

func FuzzDecode(f *testing.F) {
	good, _ := Encode(sample())
	f.Add(good)
	f.Add([]byte(Magic))
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, data []byte) {
		doc, err := Decode(data)
		if err != nil {
			return
		}
		again, err := Encode(doc)
		if err != nil {
			t.Fatalf("decoded document does not re-encode: %v", err)
		}
		back, err := Decode(again)
		if err != nil || !reflect.DeepEqual(back, doc) {
			t.Fatalf("re-encode changed the document: %v", err)
		}
	})
}

func TestLayerFlagsSurviveSaveAndLoad(t *testing.T) {
	d := sample()
	d.Nodes[0].Hidden = true
	d.Nodes[1].Locked = true
	data, err := Encode(d)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Nodes[0].Hidden || got.Nodes[0].Locked || got.Nodes[1].Hidden || !got.Nodes[1].Locked {
		t.Fatalf("flags lost: %+v %+v", got.Nodes[0], got.Nodes[1])
	}
}

func TestDocumentOptionsRoundTripAndBadOnesAreRepaired(t *testing.T) {
	doc := design.NewDocument("opts", 60, 20)
	doc.Background = "#102030"
	doc.HideGrid = true
	doc.Effects = design.Effects{Shadow: true, Scanlines: true, Vignette: true}
	data, err := Encode(doc)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.Background != "#102030" || !got.HideGrid || got.Effects != doc.Effects {
		t.Fatalf("options lost: %+v", got)
	}

	bad := design.NewDocument("bad", 100000, 100000)
	bad.Background = "chartreuse"
	data, _ = Encode(bad)
	got, err = Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.Background != "" || got.Width != design.MaxWidth || got.Height != design.MaxHeight {
		t.Fatalf("a damaged file is repaired: %+v", got)
	}
}
