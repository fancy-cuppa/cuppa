# 6. Responsive layout

Status: accepted (epic #150).

## Context

A design is an interface, not a picture. Today every node has a fixed integer `Rect` and the Go export emits those numbers as constants, so a design built at 120x40 is wrong at 80x24. Groups already scale their children from a base size (`BaseW`, `BaseH`), but nothing else adapts.

The test of the feature: rebuild Cuppa's own shell in Cuppa, with a top bar that is 100% wide and 2 rows tall.

## Decision

### Expressions, not numbers

Each of X, Y, W and H of a node may hold an **expression** instead of a number:

| Form | Meaning |
|---|---|
| `10` | 10 columns (X, W) or rows (Y, H) |
| `50%` | half of the parent's width (X, W) or height (Y, H) |
| `100% - 10` | calc: `+ - * /`, parentheses, unary minus |
| `min(50%, 40)`, `max(…)` | the smaller or larger of its arguments |

The parent is the canvas or the group that holds the node. Percentages are of the parent's size on the same axis. Results are rounded down to whole cells and then clamped so a node is at least 1x1 and never has a negative size.

### A headless layout library

`libs/layout` parses and resolves expressions. It knows nothing about nodes: it takes an expression and a parent size and returns an integer. The document owns where expressions are stored, the editor resolves them, and everything downstream (hit-testing, snapping, rendering) keeps working on resolved rectangles.

### Storage

`Node` gains an optional `Layout` with one expression per axis. A node without one is fixed, as today. `Rect` stays the **resolved** rectangle at the canvas size saved in the file, so renderers and older tools need no change. The field is optional, so by the format policy ([spec](../spec/cuppa-format.md)) the version stays at 1 and old files open unchanged. An older build would drop the expressions when it saves such a file; that is the same trade-off `hidden` and `locked` made.

Expressions apply to top-level components; the children of a group keep scaling with the group.

### Preview size

The canvas width and height become a viewport. Changing it re-resolves every node. Dragging or resizing a node with a percentage writes the result back in the same unit.

### Export

The generated program computes rectangles from the window size at run time (`tea.WindowSizeMsg`) instead of emitting constants, so the app is responsive for real.

### Runtime behaviours

`draggable` and `resizable` are properties of a node in the exported app, separate from the design-time `Locked`. The generated runtime handles the mouse for them.

## Not in this decision

- Flow layouts (rows that wrap, stacks that share space, flex-like containers). They come after constraints prove out.
- Breakpoints (different layout under a width).
- Anchors (left/right/top/bottom as in CSS); `100% - n` covers the common cases.

## Consequences

- Old files open unchanged: every node is fixed.
- The canvas size is the preview size: changing it re-resolves every component that has expressions.
- The editor must keep two things in step for a layout node: the expression and the resolved rectangle.
- The exported runtime grows with the layout and mouse code.

## Order of work

1. Expression parser and resolver (`libs/layout`).
2. Model and editor commands (no version bump), and a resizable preview viewport.
3. Editing expressions and units in the details bar.
4. Responsive Go export.
5. Draggable and resizable runtime flags.
6. Rebuild Cuppa's shell as a `.cuppa`.
