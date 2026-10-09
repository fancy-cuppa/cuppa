// Package image turns ANSI text into a picture (PNG, SVG or WebP) by running
// the Freeze command line tool.
package image

import "errors"

// Format is a picture file type Freeze can write.
type Format string

// Supported formats; the value doubles as the file extension.
const (
	PNG  Format = "png"
	SVG  Format = "svg"
	WebP Format = "webp"
)

// Formats lists every supported format.
var Formats = []Format{PNG, SVG, WebP}

// EnvFreeze names the environment variable that points at the freeze binary.
const EnvFreeze = "CUPPA_FREEZE"

// ErrFreezeMissing is returned when no freeze binary can be found.
var ErrFreezeMissing = errors.New("freeze is not installed: run `go install github.com/charmbracelet/freeze@latest`, or set " + EnvFreeze)
