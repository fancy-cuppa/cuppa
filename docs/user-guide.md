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
| Shift+click | Add it to, or remove it from, the selection (every selected component shows hollow corner marks and the group gets a dashed frame; a single selection shows solid resize handles) |
| Drag on empty canvas | Draw a box; everything it touches is selected |
| Click empty canvas, or Esc | Deselect |

**Move and resize.** Drag a selected component to move it (a whole selection moves together). Drag a corner handle to resize. Components never go below their smallest sensible size or off the canvas. While you drag, edges and centres snap to neighbours and a guide line appears; switch this off with *Snap to guides* in the details bar.

**Group components.** Select two or more components (Shift-click or draw a box) and press Ctrl+G, choose *Edit → Group*, or click *Group* in the details bar. They become one component that moves, resizes (its parts scale) and layers as a single row in Layers, and takes the place of its front-most member. Groups can contain groups. Select a group and press Ctrl+U (*Ungroup*) to get the parts back at the group's current place and size. Locked components cannot be grouped, and a locked group cannot be ungrouped.

**Component packs.** The palette is made of packs: Lip Gloss, Bubbles, Huh forms, Glamour, ntcharts and Community. *Edit → Component packs…* lists them with a checkbox each; click one to switch it off or on, and the palette follows at once. Switching a pack off only hides it from the palette: designs that already use its components keep drawing them. Your choice is remembered in `packs.json` in the Cuppa folder of your user config directory.

**Light and dark terminals, colour profiles.** With nothing selected, *Terminal* in the details bar has two choices. *Preview on* **Dark** or **Light** shows the canvas the way a light or dark terminal would: text and backgrounds that use the terminal's own colours turn dark-on-pale, while colours you chose stay, so light text on a light terminal looks as bad as it will be. It is a preview only; exports keep the terminal's own colours. *Colours* steps through **true colour**, **256 colours**, **16 colours** and **no colour**: every colour is reduced to the nearest one the profile has, in the editor and in every export, so you can see (and ship) how the design looks on a limited terminal. With *no colour* only bold and faint text remain.

**Preview.** *View → Preview design* (Ctrl+P) runs the design: the Bubbles components on the canvas become the real Bubbles models and respond to you. Click a component to give it the keyboard (Tab and Shift+Tab move between them), then type in a **text input** or **text area**, move through a **list** or **table** with the arrow keys, page a **paginator** with left and right, scroll a **viewport** with the wheel, nudge a **progress** bar with the arrows, and press Space on a **stopwatch** or **timer**; **spinners** just spin. Everything else (Lip Gloss, Huh, Glamour, charts, your own components) stays as drawn. Nothing you do in the preview changes the design. Esc, Ctrl+P again, or any editing command goes back to editing.

**Screen readers (desktop app).** The desktop window describes the screen to assistive technology: the menus and their items, the palette, the canvas size, the selected component with its position, size and properties, each layer with whether it is hidden or locked, and every dialog as the whole screen while it is open. Feedback such as "Saved design.cuppa" is announced. The description follows what is on screen; operating Cuppa still needs the mouse until keyboard navigation is designed (issue #43).

**Packs in depth.** [Component packs](packs.md) covers making, editing, installing, sharing and contributing packs, step by step.

**Make your own component.** Select a group (or two or more components) and choose *Edit → Save as component…*, then type a name. The group is saved in the pack *My components* (`my-components.cupp`) and appears in the palette at once. Placed copies scale like a group, and the text properties of its parts (titles, labels, items…) become the component's own properties, so each copy can say something different. Saving again with the same name makes `name-2`.

**Add and remove packs.** *Edit → Component packs…* also has *Add pack…*, which copies a `.cupp` file you picked into the `packs` folder of Cuppa's config directory (see the [format](spec/cupp-format.md)), and a *Remove* button on each installed pack, which asks first and deletes the file. Built-in packs can only be switched off. A pack that cannot be read is skipped, and a notice at start says why.

**Sharing designs.** A `.cuppa` file carries a copy of every pack component it uses, so it opens and draws correctly on a computer that does not have your packs; those components are listed under *From this design*. If the pack is installed, the installed one is used.

**Document options.** With nothing selected the details bar shows the canvas: its width and height (click a number to type one; up to 400 × 200), the background colour (click it to open the colour picker; *None* keeps the terminal's own) and effects applied over the whole design, in the editor and in every export: *Shadows* (a drop shadow under each component), *Scanlines* (every other row dimmed) and *Vignette* (dimmed edges). *Grid dots* only hides the editor's dotted grid; it is never exported. Each change is one undo step and is saved in the `.cuppa` file. Image exports use the background colour as the picture's backdrop.

**Resize the side bars.** Drag the thin line between the left bar and the canvas, or between the canvas and the
details bar, to make a bar wider or narrower. The line lights up when you hover it. The canvas keeps at least 20
columns, and Cuppa remembers the widths for next time.

**Edit details.** The right bar shows the selected component:

- **X, Y, W, H**: click `[-]` / `[+]` to nudge by one cell, or click the number and type.
- **Layer**: `[Front]` `[Fwd]` `[Bwd]` `[Back]` change what is drawn on top of what. The *Layers* list at the bottom shows every component, front first; click one to select it.
- **Duplicate** and **Delete**.
- **Theme**: with nothing selected, the details bar shows the design's theme: **Background**, **Text**, **Muted**, **Border** and **Secondary**. Click one to choose its colour (the colour dialog's Themes tab lists terminal colour schemes to pick from). Components that have no colour of their own follow the theme: text for labels, lists and trees, muted for placeholders and hints, border for boxes, tabs, tables and frames, secondary for prompts, bars and highlights. A new component comes with no colours of its own, so it follows the theme; the theme is saved in the `.cuppa` file. When you pick a colour for a component's property, that colour is its own and stays when the theme changes: the property then shows `[theme]`, and clicking it removes the override so the component follows the theme again.
- **Properties**: what you can change depends on the component. Text fields are click-then-type (Enter to accept, Esc to cancel). Numbers have `[-]`/`[+]`. Choices (border style, alignment…) are arrows. Yes/no options are checkboxes. Click a colour to open the colour dialog: the **16** and **256** tabs are swatches of the terminal palettes
(a click picks one), **RGB** and **HSL** have sliders for any colour (click or drag the bar, `[-]`/`[+]` for single
steps, true colour on terminals that support it). The **Themes** tab lists 342 terminal colour schemes (Dracula, Nord, Solarized and more, from [bubbletint](https://github.com/lrstanley/bubbletint)): step with `«` `◂` `▸` `»` or the wheel, jump by letter, and click one of its 16 colours, or its text, page, cursor or selection colour; the design stores that colour as plain hex. Type a palette number (`0` to `255`) or hex (`#ff5fd7`) in the
value field and press Enter, or choose **None** for no colour. Long text values wrap inside the bar.

**Undo and redo.** `[Undo]` and `[Redo]` at the top of the details bar, the Edit menu, or Ctrl+Z / Ctrl+Y. A drag or a typed edit is one step.

A `~` after a component in the palette means its preview is an approximation of the real thing.

## Files

| Menu item | Shortcut | What it does |
|---|---|---|
| File ▸ New | Ctrl+N | Start an empty design |
| File ▸ Open… | Ctrl+O | Pick a `.cuppa` file |
| File ▸ Save | Ctrl+S | Save; asks for a name the first time |
| File ▸ Save As… | Ctrl+Shift+S | Save under a new name |
| File ▸ Quit | Ctrl+Q | Leave |

If there are unsaved changes, New, Open and Quit ask whether to save, discard or cancel. Saving over an existing file asks first.

**The file dialog.** Click a folder to go into it (`..` goes up). Click a file once to pick it and again to open it, or type a name and press Enter. Typing a folder path and pressing Enter goes there. On Windows type `D:\` to switch drive. Esc or *Cancel* closes it.

## Exporting

| Export item | Result |
|---|---|
| Image (PNG / SVG / WebP) | A picture of the canvas in a window frame, made by Freeze |
| Colour text (ANSI) | A `.ans` file you can `cat` in a terminal |
| Plain text | A `.txt` file, no colours |
| Go source (Bubble Tea) | A folder with a Go program that runs the design (see below) |

**Go source.** *Export → Go source (Bubble Tea)…* asks for a folder (it suggests `<design name>-app` next to the design; it must not already hold a project) and writes a Bubble Tea v2 program there: `main.go` (the model), `layout.go` (your design: where each component sits and its properties), `runtime.go` (the code that runs them) and `go.mod`. Then:

```sh
cd <folder>
go mod tidy
go run .
```

The program places the **real components**: the Bubbles models (text input, text area, list, table, tree, file picker, help, viewport, paginator, spinner, progress, stopwatch, timer), the Huh fields, Glamour markdown, the ntcharts charts, the Evertras bubble table, and Lip Gloss boxes (with their border gradient), labels, lists, tables, trees and tabs. The Huh spinner, the layout sketches (join, place) and the community widgets without a Bubble Tea v2 library (frame, dialog, status message, toast, flex box, boxer, date picker, overlay, status bar, file tree) are drawn to look like the library they stand for, using Lip Gloss. Components that need a library of their own (Huh, Glamour, ntcharts, the bubble table, big text, the QR code) add it to `go.mod` only when the design uses them, and the image reads the file named in its File property from where the program runs. Tab and Shift+Tab move between components that take input, a click focuses one, Esc quits. Groups and pack components are expanded into their parts, hidden components are left out, and the canvas background colour fills behind the components. A component the generator does not know appears as an empty frame; the export notice and the project's README list it. Effects, colour profiles and the light-terminal preview are not part of the generated program.

If Freeze is not installed, the image items are greyed out; clicking one explains how to install it. After installing, restart Cuppa, or just click the item again.

## Shortcuts

| Key | Action |
|---|---|
| Ctrl+N / Ctrl+O / Ctrl+S | New / Open / Save |
| Ctrl+Shift+S | Save As |
| Ctrl+Z / Ctrl+Y (or Ctrl+Shift+Z) | Undo / Redo |
| Ctrl+C / Ctrl+V | Copy / Paste the selection (each paste lands a little further on; works between designs) |
| Ctrl+D | Duplicate the selection |
| Del | Delete the selection |
| Ctrl+G / Ctrl+U | Group / Ungroup |
| Ctrl+Shift+] / Ctrl+] | Move the selection to the front / one layer forward |
| Ctrl+[ / Ctrl+Shift+[ | Move the selection one layer back / to the back |
| Arrow keys | Move the selection 1 cell (with Shift, 10 cells); each press is one undo step |
| Up / Down in a number field | Add or subtract 1 (with Shift, 10); Enter ends the edit |
| Ctrl+F | Go to the component search |
| Esc | Deselect, cancel a drag, close a menu or dialog |
| Ctrl+Q | Quit (Ctrl+C is Copy) |

**On a Mac**, the desktop app and the web page show **Cmd** where this page says Ctrl (and accept it); the Cmd key does what Ctrl does. In a terminal on a Mac, the terminal program keeps Cmd for itself, so use Ctrl there, unless your terminal reports Cmd to programs (kitty, WezTerm and Ghostty can), in which case Cmd works too.

**Keys a plain terminal cannot tell apart.** Ctrl+Shift+S looks like Ctrl+S, and Ctrl+[ looks like Esc, to a terminal that only sends the old key codes. Terminals that report full key presses (those that implement the Kitty keyboard protocol or CSI u, such as kitty, WezTerm and Ghostty) send them properly, and so do the desktop app and the web page. Check yours: if Ctrl+Shift+S saves without asking for a name, it does not. Elsewhere, use the menus: *File ▸ Save As…* and the *Edit* menu's layer items do the same.

## Troubleshooting

- **Clicks do nothing**: your terminal needs mouse reporting. Try Windows Terminal, iTerm2, kitty or WezTerm.
- **Odd squares instead of ☕ or ⌕**: your font lacks those symbols; nothing is wrong.
- **Colours look off**: Cuppa uses ANSI 256 colours; use a terminal set to 256 colours or true colour.
- **"not a Cuppa design"** when opening: the file is not a `.cuppa` file or is damaged. See the [format](spec/cuppa-format.md).
