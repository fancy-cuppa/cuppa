package editor

import (
	"testing"

	"github.com/meta-tui/cuppa/libs/catalog/standard"
	"github.com/meta-tui/cuppa/libs/document/design"
)

func blank() *Editor { return New(standard.Default(), design.NewDocument("t", 100, 40)) }

func TestCanvasSizeIsCheckedAndUndoable(t *testing.T) {
	ed := blank()
	if err := ed.SetCanvasSize(0, 10); err == nil {
		t.Error("zero width must be refused")
	}
	if err := ed.SetCanvasSize(design.MaxWidth+1, 10); err == nil {
		t.Error("too wide must be refused")
	}
	if err := ed.SetCanvasSize(80, 24); err != nil {
		t.Fatal(err)
	}
	if d := ed.Document(); d.Width != 80 || d.Height != 24 {
		t.Fatalf("size = %dx%d", d.Width, d.Height)
	}
	ed.Undo()
	if d := ed.Document(); d.Width != 100 || d.Height != 40 {
		t.Fatal("undo restores the size")
	}
}

func TestBackgroundIsValidatedLikeAColour(t *testing.T) {
	ed := blank()
	if err := ed.SetBackground("not a colour"); err == nil {
		t.Fatal("bad colour must be refused")
	}
	if err := ed.SetBackground("  #A0B0C0 "); err != nil {
		t.Fatal(err)
	}
	if got := ed.Document().Background; got != "#a0b0c0" {
		t.Fatalf("background = %q", got)
	}
	if err := ed.SetBackground("#a0b0c0"); err != nil || len(ed.undo) != 1 {
		t.Errorf("setting the same colour is not a step: %v, %d steps", err, len(ed.undo))
	}
	if err := ed.SetBackground(""); err != nil {
		t.Fatal(err)
	}
	if ed.Document().Background != "" {
		t.Fatal("empty clears it")
	}
}

func TestEffectsToggleAndUndo(t *testing.T) {
	ed := blank()
	if !ed.Effect(EffectGrid) || ed.Effect(EffectShadow) {
		t.Fatal("grid starts on, shadow off")
	}
	for _, name := range []string{EffectGrid, EffectShadow, EffectScanlines, EffectVignette} {
		want := !ed.Effect(name)
		if err := ed.SetEffect(name, want); err != nil || ed.Effect(name) != want {
			t.Fatalf("%s: toggled = %v, err %v", name, ed.Effect(name), err)
		}
	}
	if !ed.Document().HideGrid || !ed.Document().Effects.Vignette {
		t.Fatalf("document = %+v", ed.Document())
	}
	ed.Undo()
	if ed.Effect(EffectVignette) {
		t.Fatal("undo turns the last one off")
	}
	if err := ed.SetEffect("sparkle", true); err == nil {
		t.Fatal("unknown effect")
	}
}
