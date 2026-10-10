package filedialog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpDownMoveThroughTheListAndEnterOpensAFolder(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.cuppa"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	m := New(Spec{Title: "Open", Dir: dir, Ext: ".cuppa"})
	m.Place(100, 30)
	m.Nav("down") // ".." is first when the folder has a parent
	m.Nav("down") // sub
	if got := m.entries[m.picked].name; got != "sub" {
		t.Fatalf("highlight on %q, want sub", got)
	}
	m.Key("", false, true, false) // Enter opens it
	if filepath.Base(m.Dir()) != "sub" {
		t.Fatalf("dir = %s, want sub", m.Dir())
	}
}

func TestHighlightingAFileFillsTheNameAndEnterAcceptsIt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.cuppa"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	m := New(Spec{Title: "Open", Dir: dir, Ext: ".cuppa"})
	m.Place(100, 30)
	m.Nav("end")
	m.Key("", false, true, false)
	out, done := m.Outcome()
	if !done || filepath.Base(out.Path) != "a.cuppa" {
		t.Fatalf("outcome %+v done=%v", out, done)
	}
}
