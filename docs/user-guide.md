# User guide

Cuppa has three panes and a menu bar:

```
 ☕ Cuppa │ File  Edit  Export  Help                 design.cuppa •
 COMPONENTS      │   canvas                     │  DETAILS
 search          │                              │  X Y W H, layer order
 ▼ Lip Gloss     │      (your design)           │  properties
   Box           │                              │  Layers
 ▼ Bubbles       │                              │
 …               │                              │
 status line: what is selected, or what just happened
```

The dot after the name means there are unsaved changes.

## Building a design

**Add a component.** Press on a name in the left bar, drag onto the canvas and release. A ghost outline shows where it will land; releasing anywhere else cancels, and so does Esc. Click a family heading (Lip Gloss, Bubbles…) to fold or unfold it. Scroll with the wheel, and type in the search box at the top to filter.

**Select.**

| Do | Result |
|---|---|
| Click a component | Select it |
| Shift+click | Add it to, or remove it from, the selection |
| Drag on empty canvas | Draw a box; everything it touches is selected |
| Click empty canvas, or Esc | Deselect |

**Move and resize.** Drag a selected component to move it (a whole selection moves together). Drag a corner handle to resize. Components never go below their smallest sensible size or off the canvas. While you drag, edges and centres snap to neighbours and a guide line appears; switch this off with *Snap to guides* in the details bar.

**Edit details.** The right bar shows the selected component:

- **X, Y, W, H**: click `[-]` / `[+]` to nudge by one cell, or click the number and type.
- **Layer**: `[Front]` `[Fwd]` `[Bwd]` `[Back]` change what is drawn on top of what. The *Layers* list at the bottom shows every component, front first; click one to select it.
- **Duplicate** and **Delete**.
- **Properties**: what you can change depends on the component. Text fields are click-then-type (Enter to accept, Esc to cancel). Numbers have `[-]`/`[+]`. Choices (border style, alignment…) are arrows. Yes/no options are checkboxes. Colours are ANSI 256 numbers with a swatch.

**Undo and redo.** `[Undo]` and `[Redo]` at the top of the details bar, the Edit menu, or Ctrl+Z / Ctrl+Y. A drag or a typed edit is one step.

A `~` after a component in the palette means its preview is an approximation of the real thing.

## Files

| Menu item | Shortcut | What it does |
|---|---|---|
| File ▸ New | Ctrl+N | Start an empty design |
| File ▸ Open… | Ctrl+O | Pick a `.cuppa` file |
| File ▸ Save | Ctrl+S | Save; asks for a name the first time |
| File ▸ Save As… | | Save under a new name |
| File ▸ Quit | Ctrl+Q | Leave |

If there are unsaved changes, New, Open and Quit ask whether to save, discard or cancel. Saving over an existing file asks first.

**The file dialog.** Click a folder to go into it (`..` goes up). Click a file once to pick it and again to open it, or type a name and press Enter. Typing a folder path and pressing Enter goes there. On Windows type `D:\` to switch drive. Esc or *Cancel* closes it.

## Exporting

| Export item | Result |
|---|---|
| Image (PNG / SVG / WebP) | A picture of the canvas in a window frame, made by Freeze |
| Colour text (ANSI) | A `.ans` file you can `cat` in a terminal |
| Plain text | A `.txt` file, no colours |

If Freeze is not installed, the image items are greyed out; clicking one explains how to install it. After installing, restart Cuppa, or just click the item again.

## Shortcuts

| Key | Action |
|---|---|
| Ctrl+N / Ctrl+O / Ctrl+S | New / Open / Save |
| Ctrl+Z / Ctrl+Y | Undo / Redo |
| Del | Delete the selection |
| Esc | Deselect, cancel a drag, close a menu or dialog |
| Ctrl+Q or Ctrl+C | Quit |

## Troubleshooting

- **Clicks do nothing**: your terminal needs mouse reporting. Try Windows Terminal, iTerm2, kitty or WezTerm.
- **Odd squares instead of ☕ or ⌕**: your font lacks those symbols; nothing is wrong.
- **Colours look off**: Cuppa uses ANSI 256 colours; use a terminal set to 256 colours or true colour.
- **"not a Cuppa design"** when opening: the file is not a `.cuppa` file or is damaged. See the [format](spec/cuppa-format.md).
