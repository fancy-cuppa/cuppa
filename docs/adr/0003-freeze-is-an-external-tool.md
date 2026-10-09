# 3. Freeze stays an external tool

Status: accepted (M3, issue #34)

## Context

Cuppa exports pictures through [Freeze](https://github.com/charmbracelet/freeze)
(`github.com/charmbracelet/freeze`, MIT). We could

1. call the `freeze` binary,
2. import Freeze as a Go library, or
3. bundle the binary in every Cuppa release.

Checked against Freeze v0.2.2:

- `main` is a CLI built with `kong`; the packages that draw the SVG and
  rasterise it are not a documented public API, and PNG/WebP output goes
  through extra native-free but large dependencies (SVG rasteriser).
- `freeze --language ansi --output out.png < colour.txt` accepts exactly what
  Cuppa produces (ANSI text) and chooses PNG, SVG or WebP from the extension.
- It installs with `go install`, `brew`, `scoop` and from release archives.

## Decision

Cuppa **runs `freeze` as a separate program** and pipes ANSI text to it
(`libs/export/image`).

- Found via `$CUPPA_FREEZE` first, then the `PATH`.
- If it is missing the export fails with a message that says how to install
  it, and nothing else in Cuppa is affected.
- Export runs in the background (a `tea.Cmd`), with a 60 s limit.
- ANSI and plain-text export (`libs/export/text`) need no external tool.

## Consequences

- The release stays small and has no second license to ship.
- Freeze upgrades reach users without a Cuppa release.
- Users need one extra install to export pictures; the error says exactly
  what to run.
- Revisit bundling if users report the extra step as a blocker, or if Freeze
  publishes a stable Go API.
