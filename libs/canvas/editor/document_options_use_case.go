package editor

import (
	"fmt"

	"github.com/meta-tui/cuppa/libs/color/space"
	"github.com/meta-tui/cuppa/libs/document/design"
)

// Effect names accepted by SetEffect.
const (
	EffectGrid      = "grid"
	EffectShadow    = "shadow"
	EffectScanlines = "scanlines"
	EffectVignette  = "vignette"
)

// SetCanvasSize resizes the canvas, as one undo step. Components outside the
// new size stay where they are; they are just not drawn until it grows again.
func (e *Editor) SetCanvasSize(w, h int) error {
	if w < 1 || h < 1 || w > design.MaxWidth || h > design.MaxHeight {
		return fmt.Errorf("editor: the canvas is 1 to %d wide and 1 to %d high", design.MaxWidth, design.MaxHeight)
	}
	if w == e.doc.Width && h == e.doc.Height {
		return nil
	}
	e.apply(func() bool {
		e.doc.Width, e.doc.Height = w, h
		return true
	})
	return nil
}

// SetBackground sets the canvas colour ("" for the terminal's own), as one undo
// step. The value is validated like a component colour.
func (e *Editor) SetBackground(value string) error {
	value, err := space.Normalise(value)
	if err != nil {
		return fmt.Errorf("editor: %w", err)
	}
	if value == e.doc.Background {
		return nil
	}
	e.apply(func() bool {
		e.doc.Background = value
		return true
	})
	return nil
}

// SetEffect turns one of the canvas options (grid, shadow, scanlines,
// vignette) on or off, as one undo step.
func (e *Editor) SetEffect(name string, on bool) error {
	if !knownEffect(name) {
		return fmt.Errorf("editor: no effect %q", name)
	}
	if e.Effect(name) == on {
		return nil
	}
	e.apply(func() bool {
		switch name {
		case EffectGrid:
			e.doc.HideGrid = !on
		case EffectShadow:
			e.doc.Effects.Shadow = on
		case EffectScanlines:
			e.doc.Effects.Scanlines = on
		case EffectVignette:
			e.doc.Effects.Vignette = on
		}
		return true
	})
	return nil
}

// Effect reports whether a canvas option is on. The grid is on unless hidden.
func (e *Editor) Effect(name string) bool {
	switch name {
	case EffectGrid:
		return !e.doc.HideGrid
	case EffectShadow:
		return e.doc.Effects.Shadow
	case EffectScanlines:
		return e.doc.Effects.Scanlines
	case EffectVignette:
		return e.doc.Effects.Vignette
	}
	return false
}

func knownEffect(name string) bool {
	switch name {
	case EffectGrid, EffectShadow, EffectScanlines, EffectVignette:
		return true
	}
	return false
}
