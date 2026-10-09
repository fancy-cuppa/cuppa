# The `.cupp` format

A `.cupp` file (a "cuppa component" pack) holds components that a user made out of
other components. It is a sibling of the [`.cuppa` design format](cuppa-format.md)
and is built the same way.

## Layout

```
offset  size  content
0       5     magic: the ASCII bytes "CUPP" followed by a line feed (0x0A)
5       2     format version, unsigned 16-bit big endian (currently 1)
7       …     a gzip stream whose content is UTF-8 JSON
```

The JSON body of version 1:

```json
{
  "pack": {
    "id": "tea-shop",
    "name": "Tea shop",
    "version": "1.0.0",
    "description": "Cards for a tea shop",
    "components": [
      {
        "id": "card",
        "name": "Card",
        "w": 20, "h": 5,
        "nodes": [
          { "id": "a", "component": "lipgloss.box",   "name": "Frame", "rect": { "x": 0, "y": 0, "w": 20, "h": 5 } },
          { "id": "b", "component": "lipgloss.label", "name": "Title", "rect": { "x": 2, "y": 1, "w": 10, "h": 1 } }
        ],
        "props": [
          { "key": "title", "label": "Title", "kind": "text", "default": "Tea", "target": "b", "targetProp": "text" }
        ]
      }
    ]
  }
}
```

| Field | Meaning |
|---|---|
| `pack.id` | Lowercase letters, digits and single dashes, at most 64 characters. Prefixes every component id: the card above is `tea-shop.card`. Must differ from the built-in packs (`lipgloss`, `bubbles`, `huh`, `glamour`, `ntcharts`, `community`) and from other installed packs. |
| `pack.name`, `version`, `description` | What the Packs dialog shows. `name` is required. |
| `component.id` | Same character rules as the pack id; unique in the pack. |
| `component.w`, `h` | Default size when placed. At least 1. Inner rectangles are relative to a box of this size and scale with the placed component. |
| `component.nodes` | Inner components, back to front, in the same shape as the nodes of a `.cuppa` file. An inner component may itself be a pack component; a loop stops at a fixed depth. |
| `component.props` | Properties the component shows in the details bar. `kind` is `text`, `int`, `bool`, `color` or `choice` (`choice` needs `choices`). `target` is the inner node and `targetProp` the property of it that the value sets. |

## Reading rules

A reader:

1. rejects a file that does not start with the magic (`not a .cupp file`);
2. rejects a version greater than it knows (`written by a newer version`);
3. rejects version 0, a missing version, a bad gzip stream, or a missing or invalid pack id or name (`damaged`);
4. refuses bodies that inflate beyond 16 MiB;
5. drops a component that is invalid (bad id, no size, repeated or damaged part, property that drives nothing) or repeats an earlier id, and reports it, without losing the rest of the pack.

A writer refuses a pack the reader would not accept, so Cuppa never saves a damaged file.

## Where packs live

Cuppa loads every `*.cupp` file in the `packs` folder of its user config directory
(`%AppData%\cuppa\packs` on Windows, `~/Library/Application Support/cuppa/packs` on macOS,
`~/.config/cuppa/packs` on Linux), in file name order. A file that cannot be read, or
whose pack id is already taken, is skipped and listed in a notice at start.

A design that uses a component of a pack that is not installed still opens: the node is
drawn as a labelled placeholder.

## Versioning

Same policy as `.cuppa`: optional fields need no version bump; anything else bumps
`CurrentVersion` and adds a migration step.
