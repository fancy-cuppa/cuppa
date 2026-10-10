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

The focused area's title is highlighted, the status bar lists the keys that work there, and the screen-reader description flags the focused area. Clicking an area, or dropping a component on the canvas, moves the focus there.

### Step 2: palette and canvas (issue #136, built)

| Key | What | Differs from the proposal |
|---|---|---|
| Up / Down, Home / End, Left / Right, `/` (palette) | Move, jump, fold or unfold a pack, start the search | As proposed |
| Enter (palette) | Place the component under the cursor at the canvas cursor and select it; on a pack header, fold it | As proposed; the cursor then moves below the new component |
| Arrows (canvas, nothing selected) | Move the canvas cursor, Shift: 10; drawn as `┼` | Shift is 10, as for moving |
| Alt+arrows (canvas, one selected) | Resize by 1 cell, Shift: 10 | Not snapped: the step is the size |
| Tab / Shift+Tab (canvas) | Select the next / previous layer, back to front, hidden layers skipped | As proposed |
| Space (canvas) | Add the next layer to the selection (repeat for a run) | The proposal toggled the "highlighted" layer; there is no separate highlight, so Space extends the selection |
| Ctrl+A | Select every visible layer | As proposed |
| Enter (canvas, with a selection) | Move the keyboard to the details bar | As proposed |

With the palette focused the arrows move through the palette, not the selection; the canvas keeps the nudge. PageUp/PageDown and Home/End for layer order were replaced by the Ctrl+[ family, as in the review decisions.

### Step 3: details bar (issue #137, built)

Every control of the details bar that a click can press is a **stop**; the stops are in reading order and the focused one is highlighted. The `[-]`/`[+]` and `◂`/`▸` buttons are not stops of their own: the value between them is, and Left/Right (or `-`/`+`) change it.

| Key | What |
|---|---|
| Tab / Down, Shift+Tab / Up | Next / previous stop (wraps) |
| Enter, Space | Press the stop: edit a text or number, press a button, toggle a checkbox, step a choice, open the colour picker, select a layer |
| Left / Right, `-` / `+` | On a number or a choice: change by 1, with Shift by 10. On other stops: previous / next |
| `h`, `l` (on a layer) | Hide or show, lock or unlock |
| Alt+Up / Alt+Down (on a layer) | Move it toward the front / back |
| Delete (on a layer) | Delete it (a locked layer refuses, as for the mouse) |

While a field is being edited its own keys apply, as before: typing, Up/Down to step a number, Enter to commit, Esc to cancel. Enter on the canvas, with a selection, puts the keyboard here.

### Step 4: dialogs (issue #138, built)

A dialog that can take navigation keys implements `modal.Navigator`; the shell offers it every key by name first, and what it does not use reaches it as typed text, as before. Enter and Esc keep their meaning everywhere.

| Dialog | Keys |
|---|---|
| Message boxes (Save changes?, notices) | Tab / Shift+Tab or arrows choose the button Enter presses |
| Name prompt | Tab / Shift+Tab move between the field, OK and Cancel; Enter on Cancel cancels |
| Open / Save | Up / Down move through the list (a file's name fills the name field), PageUp / PageDown, Home / End, Alt+Up for the parent folder; Enter opens a highlighted folder, otherwise accepts the name |
| Component packs | Up / Down choose a pack, Space switches it, Delete or `r` removes an installed pack, `a` adds one |
| Colour scheme for the theme | Left / Right step (Shift: 10), Up / Down and PageUp / PageDown by 10, Home / End, a letter or digit jumps |
| Colour picker | Tab / Shift+Tab change tab; arrows move through the swatches (16: 8 per row; 256: 6 per row, the cube's rows; Themes: the scheme's 16, with PageUp / PageDown for the scheme); on RGB and HSL Up / Down choose a slider and Left / Right change it (Shift: 10); typing still edits the value field |

The Tab order the proposal listed for the colour picker (tab bar, swatches, value, buttons) became: Tab changes the tab, the value field always takes typing, Enter accepts. A separate focus ring would have made typing a hex value need an extra key.

### Step 5: screen reader and the key table (issue #139, built)

- The shell announces what the keyboard did through the description's status line: focus moving to an area, the layer selected (name, size, place, locked or hidden), a move or resize and where it ended, a refusal with its reason (locked, edge, smallest size), the menu item or palette row or details control reached. A repeat of the same words alternates a zero-width mark so it is heard again; feedback from the file flow (Saved…) takes over when it changes.
- The description flags the focused node of the menu bar, palette, canvas and details bar, the cursor row of the palette, and the focused control of the name prompt and message boxes.
- Help → Shortcuts is built from one table (`shortcuts_table_model.go`) and a test fails when the user guide does not mention a key in it, so the dialog and the guide cannot drift.

### Modifier keys on mouse gestures (owner request)

As in an image editor: Shift while moving keeps the move to 0°, 45° or 90° (as it looks, two columns per row); Shift while resizing keeps the proportions; Alt while resizing grows from the centre; both together scale around the centre. Snapping to guides is off during a locked move. The keyboard resize (Alt+arrows) is unchanged: Alt is its own modifier there and Shift already means ten.

The proposal is built.

## Questions (answered)

Answered (2026-10-10, by the owner accepting the recommendations):

1. **F6 plus Alt+1 to Alt+4** to switch areas. Ctrl+Tab is not used.
2. **Space** toggles a checkbox; Enter edits or presses.
3. **No Vim keys** for now; cheap to add later.
4. The **canvas cursor is always drawn** when nothing is selected, so the user can see where things will land.
