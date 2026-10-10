# Examples

- `cuppa-shell.cuppa`: Cuppa's own screen (menu bar, left bar, canvas, details bar, status bar) built with
  [layout expressions](../docs/user-guide.md#building-a-design) such as `100%` and `100% - 24 - 34`. Open it in
  Cuppa and change the canvas size (or click `[80×24]`, `[120×40]`, `[160×50]`) to see it follow, or export it
  to Go source to get a program that does the same when the terminal is resized. The left and details bars are
  also marked *Resizable*.
- `mvd-colours.cuppa`: the Colours screen of MVD, written for [screen contracts](../docs/adr/0007-screen-contracts.md):
  the colour slots, the help text, the key bar and the status line are bound inputs, the status line has a *Show if*,
  the list raises `Pick slot` when clicked, and `enter`, `d`, `s` and `esc` raise events. Export it with
  `cuppa screens examples/mvd-colours.cuppa -o screens` to get `MVDColours(props, w, h)`, `MVDColoursHandle` and the
  contract types.

- `mvd-colours-rows.cuppa`: the same screen with one *Rows* component for the eight slots instead of a swatch and a
  cursor marker each ([ADR 0008](../docs/adr/0008-row-templates.md)). The program fills `Slots []MVDColoursRowsSlotsRow`
  and marks the cursor with `Style: RowSelected`; a click raises `Pick` with `Y` as the row.
- `community-widgets.cuppa`: one of each component of the third catalog round (dropdown, prompt input and select,
  data tree, PDF viewer, 3D chart), as a screenshot reference.
