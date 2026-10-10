package screenscmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meta-tui/cuppa/libs/cuppafile/disk"
	"github.com/meta-tui/cuppa/libs/document/design"
)

func saveDesign(t *testing.T, dir, name string) {
	t.Helper()
	doc := design.NewDocument(name, 40, 8)
	doc.Add(design.Node{Component: "lipgloss.label", Name: "Title", Rect: design.Rect{W: 20, H: 1},
		Props: map[string]string{"text": name}, Bind: map[string]string{"text": "Title"}, Event: "Open"})
	doc.Keys = []design.KeyBinding{{Key: "esc", Event: "Back"}}
	if _, err := disk.Save(filepath.Join(dir, strings.ToLower(name)+".cuppa"), doc); err != nil {
		t.Fatal(err)
	}
}

func TestScreensCommandWritesAFolderOfDesigns(t *testing.T) {
	designs, out := t.TempDir(), filepath.Join(t.TempDir(), "screens")
	saveDesign(t, designs, "Colours")
	saveDesign(t, designs, "Settings")

	var stdout, stderr bytes.Buffer
	if code := Run([]string{designs, "-o", out}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit %d: %s", code, stderr.String())
	}
	for _, name := range []string{"colours_screen_contract.go", "settings_screen_view.go", "cuppa_runtime.go"} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Errorf("%s not written", name)
		}
	}
	src, _ := os.ReadFile(filepath.Join(out, "colours_screen_contract.go"))
	if !strings.Contains(string(src), "package screens") || !strings.Contains(string(src), "ColoursOpen") {
		t.Errorf("contract:\n%s", src)
	}

	// A design that is gone takes its files with it; a hand-written file in
	// the way stops the export.
	if err := os.Remove(filepath.Join(designs, "settings.cuppa")); err != nil {
		t.Fatal(err)
	}
	if code := Run([]string{designs, "-o", out}, &stdout, &stderr); code != 0 {
		t.Fatalf("second export: %s", stderr.String())
	}
	if _, err := os.Stat(filepath.Join(out, "settings_screen_view.go")); err == nil {
		t.Error("stale screen kept")
	}
	if err := os.WriteFile(filepath.Join(out, "cuppa_runtime.go"), []byte("package screens\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := Run([]string{designs, "-o", out}, &stdout, &stderr); code != 1 || !strings.Contains(stderr.String(), "did not generate") {
		t.Errorf("hand-written file overwritten: code %d %s", code, stderr.String())
	}
}

func TestScreensCommandUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run(nil, &stdout, &stderr); code != 2 || !strings.Contains(stderr.String(), "usage") {
		t.Errorf("code %d %q", code, stderr.String())
	}
}
