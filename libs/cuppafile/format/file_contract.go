// Package format encodes and decodes the .cuppa file format: a short binary
// header followed by gzip-compressed JSON. See docs/spec/cuppa-format.md.
package format

import "errors"

// Magic opens every .cuppa file.
const Magic = "CUPPA\n"

// CurrentVersion is the format version this build writes and the newest it reads.
const CurrentVersion uint16 = 1

// Extension is the file extension of a saved design, with its dot.
const Extension = ".cuppa"

// MaxDecodedSize bounds the JSON a file may inflate to, so a hostile file
// cannot exhaust memory.
const MaxDecodedSize = 64 << 20

// Errors a decode can return; test with errors.Is.
var (
	ErrNotCuppa = errors.New("not a .cuppa file")
	ErrTooNew   = errors.New("file was written by a newer version of Cuppa")
	ErrCorrupt  = errors.New("file is damaged")
)
