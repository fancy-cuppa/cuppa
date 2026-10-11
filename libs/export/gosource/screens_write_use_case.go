package gosource

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// WriteScreens puts the generated screens in dir, creating it if needed. A
// file Cuppa generated earlier is replaced; a file it did not write is never
// touched, and the call refuses the whole folder when one is in the way. Other
// screens already in the folder stay.
func WriteScreens(dir string, p Project) error {
	for name := range p.Files {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		if !strings.HasPrefix(string(data), GeneratedHeader) {
			return fmt.Errorf("%s has a %s that Cuppa did not generate; screens need a folder of their own", dir, name)
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, content := range p.Files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// CheckScreens says how the folder differs from what GenerateScreens made: a
// file that is missing or has other content, and a generated file of a screen
// that is no longer designed. It writes nothing; an empty result means the
// folder is current. Line endings are not compared.
func CheckScreens(dir string, p Project) []string {
	var out []string
	same := func(a, b string) bool {
		return strings.ReplaceAll(a, "\r\n", "\n") == strings.ReplaceAll(b, "\r\n", "\n")
	}
	names := make([]string, 0, len(p.Files))
	for name := range p.Files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		switch {
		case err != nil:
			out = append(out, name+": missing")
		case !same(string(data), p.Files[name]):
			out = append(out, name+": differs from what the designs export")
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}
	for _, e := range entries {
		if _, kept := p.Files[e.Name()]; kept || e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		if data, err := os.ReadFile(filepath.Join(dir, e.Name())); err == nil && strings.HasPrefix(string(data), GeneratedHeader) {
			out = append(out, e.Name()+": generated for a screen that is no longer designed")
		}
	}
	return out
}

// ReplaceScreens is WriteScreens for a project that holds every screen of the
// folder: the generated files of a screen that is no longer designed are
// removed. Files Cuppa did not generate are never removed.
func ReplaceScreens(dir string, p Project) error {
	if err := WriteScreens(dir, p); err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if _, kept := p.Files[e.Name()]; kept || e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err == nil && strings.HasPrefix(string(data), GeneratedHeader) {
			if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
