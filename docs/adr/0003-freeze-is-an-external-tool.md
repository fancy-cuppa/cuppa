# 3. Freeze stays an external tool

Status: accepted (M3, issue #34)

## Context

Cuppa exports pictures through [Freeze](https://github.com/charmbracelet/freeze)
(`github.com/charmbracelet/freeze`, MIT). We could

1. call the `freeze` binary,
2. import Freeze as a Go library, or
3. bundle the binary in every Cuppa release.

Checked against Freeze v0.2.2 (source read in the module cache):

- The code that turns ANSI into SVG and PNG is in the root `package main`
  (`ansi.go`, `png.go`, `main.go`), which Go cannot import. Only helper
  packages (`svg`, `font`, `input`) are importable.
- PNG and WebP output goes through a wasm SVG renderer (`resvg-go`). Copying
  that into Cuppa would mean forking and maintaining it.
- `freeze --language ansi --output out.png < colour.txt` accepts exactly what
  Cuppa produces (ANSI text) and chooses PNG, SVG or WebP from the extension.
- It installs with `go install`, `brew`, and release archives for Linux,
  macOS and Windows.

## Decision

Cuppa **runs `freeze` as a separate program** and pipes ANSI text to it
(`libs/export/image`).

- Found via `$CUPPA_FREEZE` first, then the `PATH`.
- If it is missing, the first launch shows a one-time notice (marker file in
  the user config directory), the PNG/SVG/WebP menu items are greyed out, and
  clicking one opens a popup with install instructions. Nothing else in Cuppa
  is affected.
- Export runs in the background (a `tea.Cmd`), with a 60 s limit.
- ANSI and plain-text export (`libs/export/text`) need no external tool.

## Consequences

- The release stays small and has no second license to ship.
- Freeze upgrades reach users without a Cuppa release.
- Users need one extra install to export pictures; the error says exactly
  what to run.
- Revisit bundling if users report the extra step as a blocker, or if Freeze
  publishes a stable Go API.
