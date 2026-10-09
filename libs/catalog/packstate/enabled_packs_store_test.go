package packstate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPacksAreOnUntilSwitchedOffAndTheChoiceIsRemembered(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cuppa", "packs.json")
	s := Open(path)
	if !s.Enabled("bubbles") {
		t.Fatal("packs start on")
	}
	s.SetEnabled("bubbles", false)
	s.SetEnabled("community", false)
	s.SetEnabled("community", true)

	again := Open(path)
	if again.Enabled("bubbles") || !again.Enabled("community") {
		t.Fatal("the next run sees the same choices")
	}
}

func TestMissingOrDamagedFilesMeanEverythingIsOn(t *testing.T) {
	dir := t.TempDir()
	if !Open(filepath.Join(dir, "none.json")).Enabled("huh") {
		t.Fatal("missing file")
	}
	bad := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(bad, []byte("{nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Open(bad).Enabled("huh") {
		t.Fatal("damaged file")
	}
}

func TestWithoutAPathNothingIsWritten(t *testing.T) {
	s := Open("")
	s.SetEnabled("huh", false)
	if s.Enabled("huh") {
		t.Fatal("the choice still holds in memory")
	}
}
