# 5. Keyboard navigation

Status: proposed (issue #43). Parts are built, from the review in [Decisions from the review](#decisions-from-the-review); the rest is still to agree on.

## Context

Cuppa is operated almost entirely with the mouse. The keyboard today does this and little else:

| Where | Keys |
|---|---|
| Anywhere | Ctrl+N, Ctrl+O, Ctrl+S, Ctrl+Q, Ctrl+Z, Ctrl+Y, Ctrl+G (group), Ctrl+U (ungroup), Ctrl+P (preview), Del, Esc |
| Palette search box | typing, Backspace, Enter, Esc |
| Details bar | typing into a field you clicked, Enter, Esc |
| Dialogs | Enter (default button), Esc, typing in text fields |
| Preview | Tab, Shift+Tab, and whatever the focused component reads |

Everything else (choosing a menu, picking a component, selecting, moving or resizing, changing a property, ticking a layer's eye, choosing a colour) needs a pointer. That excludes people who cannot or prefer not to use one, makes the desktop app's screen-reader description (#67) describe things nobody can operate, and is slow for the repeated work a designer does.

The architecture already helps: every edit is an **intent** on the headless editor (`MoveSelectionBy`, `SetRect`, `Reorder`, `SetProp`, `SetHidden`, `SetLocked`, `MoveLayer`, `Group`…), and the panes already know their rows. What is missing is a *focus model* and key handlers that call the same intents.

## Goals

1. **Everything the mouse can do, the keyboard can do**, with the same intents, undo steps and refusals (a locked layer refuses the same way).
2. **One focus model** the user can see and a screen reader can announce, not a different set of rules per pane.
3. **Discoverable**: the status bar and the Shortcuts dialog tell you the keys that work *here*.
4. **No new meaning for existing keys.** What works today keeps working.
5. Keys that work in a plain terminal, in Windows Terminal, and in the desktop window's xterm.js.

Non-goals: Vim-style modal editing, user-remappable keys, chorded sequences (like `g g`). They can come later on top of this.

## Proposal

### Focus model

The shell has one **focus**: the menu bar, the palette, the canvas, or the details bar. A dialog, when open, owns the keyboard as it already does.

- **F6 / Shift+F6** move focus to the next / previous area (the convention of desktop apps). Ctrl+Tab is avoided because terminals cannot send it reliably.
- The focused area shows it: its title is highlighted and a cursor or marker sits on its current row. Clicking an area focuses it, so mouse and keyboard stay in step.
- **Esc** goes up one level: close a menu, cancel an edit, leave a pane for the canvas, and on the canvas clear the selection (as today).
- The screen-reader description already has a `Focused` flag on nodes; the shell sets it from this model, and a change of focus is announced.

### Menu bar

- **F10** (or Alt+F, Alt+E, Alt+V, Alt+X, Alt+H, one per menu) opens the menu; Left/Right change menu, Up/Down move through items, Enter runs one, Esc closes. Items that are not available are skipped, with the same notice as a click.

### Palette

- Up/Down move through the rows, Left/Right fold or unfold a pack, Home/End jump, `/` starts the search (as the search box does now).
- **Enter places the highlighted component on the canvas** at the *canvas cursor* (below) and selects it, so a whole design can be built without a pointer.

### Canvas

- A **canvas cursor** (a visible cell) exists whenever nothing is selected; arrow keys move it, Shift+arrow by 5 cells. It is where Enter in the palette places things.
- **Tab / Shift+Tab** select the next / previous layer; the selection is announced by name. Space adds or removes the highlighted layer from a multiple selection; Ctrl+A selects all.
- With something selected: **arrows move it** by 1 cell (Shift: 5), **Alt+arrows resize** it (width and height, Shift: 5), snapping on or off as the *Snap to guides* setting says. Refusals (locked, edge of the canvas) are announced.
- PageUp / PageDown move the selection one layer forward / back, Home / End to the front / back (the details bar's Front/Fwd/Bwd/Back).
- Del deletes, Ctrl+D duplicates, Ctrl+G / Ctrl+U group and ungroup (already bound), Ctrl+Z / Ctrl+Y undo and redo.
- Enter on a selection moves focus to its settings in the details bar.

### Details bar

- Tab / Shift+Tab (or Up/Down) move between its fields and buttons in the order shown; the focused one is highlighted.
- Enter edits a text or number field (typing replaces, Enter commits, Esc cancels, as now); on a button it presses it; on a colour it opens the picker.
- Left/Right (or `-`/`+`) change a number or a choice by one, Shift by ten. Space toggles a checkbox.
- The **layers list** is a field group of its own: Up/Down move, Enter selects, `h` hides or shows, `l` locks or unlocks, Alt+Up/Alt+Down reorder, Del deletes (the trash).

### Dialogs

- Tab / Shift+Tab move between controls, Enter presses the default one, Esc cancels. Lists (files, packs) use Up/Down, Enter, Space (toggle a pack), and Del or `r` on a removable pack.
- **Colour picker:** Tab moves between the tab bar, the swatch grid or sliders, the value field and the buttons; arrows move through swatches or change the focused slider by 1 (Shift: 10); typing in the value field already works.

### Preview

Unchanged: Tab and Shift+Tab move between running components, the focused one reads the keys, Esc leaves.

### Discoverability

- The status bar shows the keys that work in the focused area, replacing today's single hint.
- *Help → Shortcuts* becomes a table by area, generated from the same key table the shell uses, so it cannot drift.

## Keys to avoid

Terminals cannot tell these apart or deliver them reliably, so none is used: Ctrl+I (it is Tab), Ctrl+M (Enter), Ctrl+[ (Esc), Ctrl+Shift+letter, Ctrl+Tab, Ctrl+Backspace, and bare Alt+Shift chords. Alt+letter and Alt+arrow work in Windows Terminal, iTerm and xterm.js, and are what this proposal relies on; where a terminal swallows one (some send Alt as Esc-prefix) the same action stays reachable through the menus and the details bar.

## Implementation plan

Each step is its own pull request, tested by feeding key sequences to the shell and checking the editor and the description.

1. **Focus model, F6, status hints, menu bar keys.** The foundation; nothing else works without it. Built (issue #135). Dialog Tab order moved to step 4.
2. **Palette and canvas:** the canvas cursor, place with Enter, Tab selection, move and resize with arrows, layer order keys.
3. **Details bar:** field focus, editing, steppers, checkboxes, layers list keys.
4. **Colour picker and the remaining dialogs.**
5. **Screen reader:** `Focused` from the focus model, announcements for selection, moves and refusals; update the user guide.

Existing mouse behaviour and the current shortcuts are not touched in any step.

## Consequences

- Cuppa becomes fully operable without a pointer, which is also what makes the screen-reader description useful.
- A focus model adds state to the shell; it is small (one enum plus a row index per pane) and replaces the ad-hoc "which pane owns the pointer" checks for keys.
- The keys above become a documented contract; changing one later is a breaking change for muscle memory, so the table lives in one file and the Shortcuts dialog is generated from it.

## Decisions from the review

The owner reviewed the proposal with screenshots (2026-10-09) and asked for these, which are **built** (issue #118):

| Key | What | Differs from the proposal |
|---|---|---|
| Ctrl+Shift+S | Save As (listed in the File menu) | New: the proposal avoided Ctrl+Shift+letter |
| Ctrl+D | Duplicate (listed in the Edit menu) | As proposed |
| Arrows | Move the selection 1 cell; **Shift: 10** | The proposal said Shift: 5 |
| Up / Down in a number field | +1 / -1; **Shift: 10** | New |
| Ctrl+F | Focus the component search | New |
| Ctrl+Shift+] / Ctrl+] / Ctrl+[ / Ctrl+Shift+[ | To front / forward one / back one / to back (also in the Edit menu) | The proposal used PageUp/PageDown and Home/End and avoided Ctrl+[ |
| Cmd on a Mac | Menus say Cmd in the desktop app and the web page; Cmd does what Ctrl does | New |

How the keys a terminal cannot express reach the app:

- Bubble Tea reads the enhanced form of a key (CSI u), so Ctrl+Shift+S and Ctrl+[ are distinct keys when the terminal reports them. Kitty-protocol terminals also report Cmd as Super, which the shell treats as Ctrl.
- The desktop app and the web page catch these shortcuts in the page (`shortcuts.ts`) and send the enhanced form themselves, because xterm.js would send Ctrl+Shift+S as Ctrl+S and Ctrl+[ as Esc. On a Mac the page sends Cmd+key as Ctrl+key.
- The menus show Cmd on a Mac in the desktop app (from the operating system) and in the browser (from the User Agent). The terminal program keeps showing Ctrl: the terminal emulator, not Cuppa, owns Cmd there.
- Where a terminal sends only the old codes, the same actions stay in the menus.

### Step 1: focus model and menu keys (issue #135, built)

| Key | What |
|---|---|
| F6 / Shift+F6 | Next / previous area: menu bar, palette, canvas, details bar |
| Alt+1 to Alt+4 | Jump to the menu bar, palette, canvas or details bar |
| F10, Alt+F / E / V / X / H | Open the first menu, or the File / Edit / View / Export / Help menu |
| Left / Right, Up / Down, Home / End, Enter, Esc | In the menu bar: change menu, move through items (separators and disabled items are skipped), run, close |
| Esc | From the menu bar, palette or details bar: back to the canvas without touching the selection |

The focused area's title is highlighted, the status bar lists the keys that work there, and the screen-reader description flags the focused area. Clicking an area, or dropping a component on the canvas, moves the focus there. The arrows still move the selection whichever area has the focus, so existing muscle memory is unchanged; steps 2 and 3 give the palette and the details bar their own arrow keys.

The rest of the proposal (steps 2 to 5) is not built yet.

## Questions (answered)

Answered (2026-10-10, by the owner accepting the recommendations):

1. **F6 plus Alt+1 to Alt+4** to switch areas. Ctrl+Tab is not used.
2. **Space** toggles a checkbox; Enter edits or presses.
3. **No Vim keys** for now; cheap to add later.
4. The **canvas cursor is always drawn** when nothing is selected, so the user can see where things will land.
