package fileflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportScreenWritesTheContractAndKeepsOtherScreens(t *testing.T) {
	f, ed, dir := setup(t)
	id, err := ed.Add("lipgloss.label", 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := ed.SetBinding(id, "text", "Title"); err != nil {
		t.Fatal(err)
	}
	f.ExportScreen()
	if f.Modal() == nil || !strings.Contains(stripped(strings.Join(f.Modal().Lines(), "\n")), "screens") {
		t.Fatal("a folder prompt suggests a screens folder")
	}
	press(f, true, false)
	out := filepath.Join(dir, "screens")
	for _, name := range []string{"untitled_screen_contract.go", "untitled_screen_view.go", "cuppa_runtime.go", "cuppa_screens.go"} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Fatalf("%s not written: %v", name, err)
		}
	}
	src, _ := os.ReadFile(filepath.Join(out, "untitled_screen_contract.go"))
	if !strings.Contains(string(src), "Title string") {
		t.Errorf("contract:\n%s", src)
	}
	press(f, true, false) // the notice

	// Another design in the same folder leaves the first one alone.
	ed.Load(ed.Document())
	doc := ed.Document()
	doc.Name = "Settings"
	ed.Load(doc)
	f.ExportScreen()
	press(f, true, false)
	if _, err := os.Stat(filepath.Join(out, "untitled_screen_view.go")); err != nil {
		t.Fatal("exporting a second screen removed the first")
	}
	if _, err := os.Stat(filepath.Join(out, "settings_screen_view.go")); err != nil {
		t.Fatal("second screen not written")
	}
}
