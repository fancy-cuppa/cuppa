# The `.cuppa` file format

Version **1**. Code: `libs/cuppafile/format`.

A `.cuppa` file is a short binary header followed by gzip-compressed JSON.
The header lets a reader reject foreign or too-new files before it touches
any JSON; the JSON stays easy to inspect and diff after `gunzip`.

```
offset  size  field
0       6     magic      the bytes "CUPPA\n"  (43 55 50 50 41 0A)
6       2     version    unsigned, big endian (currently 0x0001)
8       …     body       gzip stream; inflates to one JSON object
```

## Body (version 1)

```json
{
  "document": {
    "name": "Login screen",
    "width": 120,
    "height": 40,
    "seq": 3,
    "nodes": [
      {
        "id": "n1",
        "component": "huh.input",
        "name": "Input",
        "rect": { "x": 4, "y": 2, "w": 30, "h": 3 },
        "props": { "title": "Name", "placeholder": "Ada" }
      }
    ]
  }
}
```

| Field | Meaning |
|---|---|
| `document.name` | Display name of the design. |
| `document.width`, `height` | Canvas size in terminal cells. At least 1. |
| `document.seq` | Counter behind generated node ids; only ever grows, so ids are never reused. |
| `nodes` | Back to front: the **last** node is drawn on top. |
| `node.id` | Unique within the document. |
| `node.component` | Catalog id (`bubbles.list`, `huh.select`, `community.flexbox`…). Unknown ids load and render as a labelled placeholder, so a file from a newer catalog still opens. |
| `node.rect` | Position and size in cells. `w` and `h` are at least 1. |
| `node.props` | String values the designer changed. Anything not listed uses the component's default. |

## Reading rules

A reader:

1. rejects a file that does not start with the magic (`not a .cuppa file`);
2. rejects a version greater than it knows (`written by a newer version`);
3. rejects version 0, a missing version or a bad gzip stream (`damaged`);
4. refuses bodies that inflate beyond 64 MiB;
5. migrates older versions one step at a time to the current one;
6. drops nodes with an empty or repeated id or a rect under 1×1, and raises a
   canvas under 1×1 to 1×1, so a damaged or hand-edited file never gives the
   editor an impossible document.

## Writing rules

A writer always writes the current version, and saves atomically: it writes a
temporary file next to the target and renames it into place, so a crash never
leaves half a file where the old one was.

## Changing the format

Adding an optional field needs no version bump. Anything that changes meaning
or removes a field:

1. raise `CurrentVersion` in `file_contract.go`;
2. add `migrations[oldVersion]` in `migrations_algorithm.go` that turns the old
   JSON body into the new one;
3. keep a fixture of the old version in the tests.

`FuzzDecode` checks that no input crashes the reader and that anything it does
accept survives an encode/decode round trip unchanged.
