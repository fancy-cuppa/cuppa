package image

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatOf(t *testing.T) {
	for path, want := range map[string]Format{"a.png": PNG, "b/C.SVG": SVG, "d.webp": WebP} {
		if got, ok := FormatOf(path); !ok || got != want {
			t.Errorf("%s: %v %v", path, got, ok)
		}
	}
	if _, ok := FormatOf("x.gif"); ok {
		t.Error("gif should be unsupported")
	}
}

func TestLocate(t *testing.T) {
	missing := func(string) (string, error) { return "", errors.New("nope") }
	if _, err := locate(func(string) string { return "" }, missing); !errors.Is(err, ErrFreezeMissing) {
		t.Fatalf("path miss: %v", err)
	}
	if _, err := locate(func(string) string { return filepath.Join(t.TempDir(), "x") }, missing); !errors.Is(err, ErrFreezeMissing) {
		t.Fatalf("env miss: %v", err)
	}
	f := filepath.Join(t.TempDir(), "freeze")
	_ = os.WriteFile(f, nil, 0o755)
	if got, err := locate(func(string) string { return f }, missing); err != nil || got != f {
		t.Fatalf("env hit: %v %v", got, err)
	}
	found := func(string) (string, error) { return "/bin/freeze", nil }
	if got, _ := locate(func(string) string { return "" }, found); got != "/bin/freeze" {
		t.Fatalf("path hit: %v", got)
	}
}

func TestWritePassesArgsAndInput(t *testing.T) {
	var gotArgs []string
	var gotIn string
	runner := func(_ context.Context, _ string, a []string, in string) (string, error) {
		gotArgs, gotIn = a, in
		return "", nil
	}
	err := write("out.png", "\x1b[31mhi", Options{Window: true, Background: "#000"}, func() (string, error) { return "freeze", nil }, runner)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(gotArgs, " ")
	for _, want := range []string{"--language ansi", "--output out.png", "--window", "--background #000"} {
		if !strings.Contains(joined, want) {
			t.Errorf("args %q lack %q", joined, want)
		}
	}
	if gotIn != "\x1b[31mhi" {
		t.Errorf("stdin = %q", gotIn)
	}
}

func TestWriteErrors(t *testing.T) {
	ok := func() (string, error) { return "freeze", nil }
	if err := write("a.gif", "", Options{}, ok, nil); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("gif: %v", err)
	}
	if err := write("a.png", "", Options{}, func() (string, error) { return "", ErrFreezeMissing }, nil); !errors.Is(err, ErrFreezeMissing) {
		t.Errorf("missing: %v", err)
	}
	failing := func(context.Context, string, []string, string) (string, error) { return "boom", errors.New("exit 1") }
	if err := write("a.png", "", Options{}, ok, failing); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Errorf("failure should carry freeze output: %v", err)
	}
}

// TestRealFreeze runs the real binary when it is installed.
func TestRealFreeze(t *testing.T) {
	if _, err := Locate(); err != nil {
		t.Skip("freeze not installed")
	}
	out := filepath.Join(t.TempDir(), "out.svg")
	if err := Write(out, "\x1b[31mhello\x1b[0m\n", Options{}); err != nil {
		t.Fatal(err)
	}
	if fi, err := os.Stat(out); err != nil || fi.Size() == 0 {
		t.Fatalf("no picture written: %v", err)
	}
}
