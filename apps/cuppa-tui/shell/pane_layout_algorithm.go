package shell

import "github.com/meta-tui/cuppa/libs/document/design"

// layout places the three panes below the title bar and above the status bar.
// A one-cell separator column sits between neighbouring panes.
type layout struct {
	palette, stage, inspector design.Rect
}

const (
	preferredPaletteW   = 26
	preferredInspectorW = 32
	minStageW           = 20
	chromeRows          = 2 // title bar + status bar
)

func computeLayout(w, h int) layout {
	pal, ins := preferredPaletteW, preferredInspectorW
	// Give the stage room first: shrink the side bars on narrow terminals.
	for w-pal-ins-2 < minStageW && (pal > 16 || ins > 20) {
		if pal > 16 {
			pal--
		}
		if ins > 20 {
			ins--
		}
	}
	stageW := max(w-pal-ins-2, 1)
	body := max(h-chromeRows, 1)
	return layout{
		palette:   design.Rect{X: 0, Y: 1, W: pal, H: body},
		stage:     design.Rect{X: pal + 1, Y: 1, W: stageW, H: body},
		inspector: design.Rect{X: pal + 1 + stageW + 1, Y: 1, W: ins, H: body},
	}
}

// paneAt returns the pane under a screen cell.
func (l layout) paneAt(x, y int) pane {
	switch {
	case l.palette.Contains(x, y):
		return inPalette
	case l.stage.Contains(x, y):
		return inStage
	case l.inspector.Contains(x, y):
		return inInspector
	}
	return nowhere
}
