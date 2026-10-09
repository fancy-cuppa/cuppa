package packcmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const goodJSON = `{"pack":{"id":"tea-shop","name":"Tea shop","version":"1.0.0","description":"Cards","components":[
 {"id":"card","name":"Card","w":20,"h":5,"nodes":[{"id":"a","component":"lipgloss.box","name":"Frame","rect":{"x":0,"y":0,"w":20,"h":5}}]}]}}`

func run(args ...string) (int, string, string) {
	var out, errOut bytes.Buffer
	code := Run(args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestBuildThenJSONThenCheckRoundTrip(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tea.json")
	if err := os.WriteFile(src, []byte(goodJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out, errs := run("build", src)
	if code != 0 || !strings.Contains(out, "tea-shop") {
		t.Fatalf("build: %d %q %q", code, out, errs)
	}
	built := filepath.Join(dir, "tea.cupp")
	code, out, _ = run("check", built)
	if code != 0 || !strings.Contains(out, "tea-shop.card") {
		t.Fatalf("check .cupp: %d %q", code, out)
	}
	code, out, _ = run("json", built)
	if code != 0 || !strings.Contains(out, `"id": "tea-shop"`) {
		t.Fatalf("json: %d %q", code, out)
	}
	if code, _, _ = run("check", src); code != 0 {
		t.Fatal("check also reads the JSON form")
	}
}

func TestProblemsAreReportedWithANonZeroExit(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(bad, []byte(strings.Replace(goodJSON, `"tea-shop"`, `"Tea Shop"`, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, _, errs := run("build", bad); code != 1 || !strings.Contains(errs, "pack id") {
		t.Fatalf("a bad id is explained: %d %q", code, errs)
	}
	if code, _, _ := run("build", filepath.Join(dir, "missing.json")); code != 1 {
		t.Fatal("a missing file fails")
	}
	if code, _, _ := run("build", bad, filepath.Join(dir, "x.txt")); code != 1 {
		t.Fatal("the output must be .cupp")
	}
	if code, _, errs := run(); code != 2 || !strings.Contains(errs, "usage") {
		t.Fatal("no arguments shows the usage")
	}
	if code, _, _ := run("frobnicate"); code != 2 {
		t.Fatal("unknown subcommand")
	}
	if code, _, _ := run("json"); code != 1 {
		t.Fatal("missing file argument")
	}
}

func TestCatalogListsComponentsWithPropertyKeys(t *testing.T) {
	code, out, _ := run("catalog", "label")
	if code != 0 || !strings.Contains(out, "lipgloss.label") || !strings.Contains(out, "text") || !strings.Contains(out, "align") {
		t.Fatalf("catalog: %d\n%s", code, out)
	}
	if strings.Contains(out, "bubbles.spinner") {
		t.Fatal("the filter narrows the list")
	}
	if code, out, _ = run("catalog"); code != 0 || !strings.Contains(out, "bubbles.spinner") {
		t.Fatal("no filter lists everything")
	}
}
