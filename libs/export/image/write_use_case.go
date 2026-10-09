package image

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Runner runs freeze with ANSI text on stdin; tests replace it.
type Runner func(ctx context.Context, bin string, args []string, stdin string) (output string, err error)

// Options style the picture. The zero value gives Freeze's defaults without
// window decoration.
type Options struct {
	// Window draws macOS-style window controls.
	Window bool
	// Background fills behind the text, e.g. "#171717"; empty keeps Freeze's default.
	Background string
}

// timeout bounds one freeze run so a stuck process cannot hang the app.
const timeout = 60 * time.Second

// FormatOf picks the format from a file name's extension.
func FormatOf(path string) (Format, bool) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	for _, f := range Formats {
		if string(f) == ext {
			return f, true
		}
	}
	return "", false
}

// Write renders ansi to a picture at path. The format comes from the extension.
func Write(path, ansi string, opts Options) error {
	return write(path, ansi, opts, Locate, run)
}

func write(path, ansi string, opts Options, locate func() (string, error), runner Runner) error {
	if _, ok := FormatOf(path); !ok {
		return fmt.Errorf("unsupported picture type %q: use .png, .svg or .webp", filepath.Ext(path))
	}
	bin, err := locate()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if out, err := runner(ctx, bin, args(path, opts), ansi); err != nil {
		return fmt.Errorf("freeze failed: %w: %s", err, strings.TrimSpace(out))
	}
	return nil
}

func args(path string, opts Options) []string {
	a := []string{"--language", "ansi", "--output", path}
	if opts.Window {
		a = append(a, "--window")
	}
	if opts.Background != "" {
		a = append(a, "--background", opts.Background)
	}
	return a
}

func run(ctx context.Context, bin string, args []string, stdin string) (string, error) {
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdin = strings.NewReader(stdin)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	return out.String(), err
}
