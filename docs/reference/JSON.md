# JSON documents (`JSON.*`)

The **`JSON`** namespace decodes and mutates JSON values using Go’s [`encoding/json`](https://pkg.go.dev/encoding/json). Values live on the VM heap as **`JSON`** handles (`TagJSON`). The root may be an **object**, **array**, or **scalar**; nested structures use **dot + bracket paths** (see below).

## Why paths instead of flat keys

Earlier builds only supported flat objects. The current implementation keeps **one decoded `interface{}` tree** per handle so you can work with **real-world files** (levels, configs, save data) without flattening.

## Parsing and construction

| Command | Purpose |
|--------|---------|
| `JSON.PARSE(path)` | **Read a UTF-8 file** from disk and decode JSON. |
| `JSON.PARSESTRING(s)` | Decode JSON **from a string** (use this for inline text or loaded buffers). |
| `JSON.MAKE()` | New empty **object** `{}`. |
| `JSON.MAKEARRAY()` | New empty **array** `[]`. |
| `JSON.FREE(handle)` | Release the heap object (`Free` is idempotent). |

**Breaking change (by design):** `JSON.PARSE` is now **file-oriented**. Existing scripts that passed a JSON **string** to `JSON.PARSE` should switch to **`JSON.PARSESTRING`** (see [`examples/rpg/main.mb`](../../examples/rpg/main.mb)).

## Path syntax

Paths are strings composed of:

- **`.`** — separates object keys (`spawn.x`).
- **`[index]`** — array index, 0-based (`items[0].id`).

Examples:

- `spawn.x`
- `items[1].qty`
- `[0]` when the **root** is an array

Empty path `""` refers to the **root** value (for **`LEN`**, **`TOCSV`**, etc.).

## Reading

| Command | Purpose |
|--------|---------|
| `JSON.HAS(handle, path)` | Whether a value exists at the path (map key present or array index in range). |
| `JSON.TYPE(handle, path)` | Rough type string: `object`, `array`, `string`, `number`, `bool`, `null`, or `missing`. |
| `JSON.LEN(handle, path)` | Length of object (key count), array, or string; `0` if not applicable. |
| `JSON.KEYS(handle, path)` | Object keys as a **`StringList`** handle (sorted lexicographically for stable output). |
| `JSON.GETSTRING` / `GETINT` / `GETFLOAT` / `GETBOOL` | Read scalars; optional **3-argument** overloads supply a **default** when missing. |
| `JSON.GETARRAY` / `JSON.GETOBJECT` | Return a **new** `JSON` handle aliasing the nested value (free separately). |

Numbers decode as `float64` internally; **`GETINT`** coerces.

## Writing

| Command | Purpose |
|--------|---------|
| `JSON.SETSTRING` / `SETINT` / `SETFLOAT` / `SETBOOL` | Create intermediate objects/arrays as needed. |
| `JSON.SETNULL` | Store JSON `null`. |
| `JSON.DELETE` | Remove an object key or splice an array index. |
| `JSON.CLEAR` | Empty an object or array at a path. |
| `JSON.APPEND` | Append a value to an array at `path` (creates the array if missing). |

## Serialization and files

| Command | Purpose |
|--------|---------|
| `JSON.TOSTRING` / `JSON.MINIFY` | Compact JSON text. |
| `JSON.PRETTY` | Indented JSON (two spaces). |
| `JSON.TOFILE` / `JSON.TOFILEPRETTY` | Write to disk. |

## `JSON.QUERY` (minimal)

`JSON.QUERY(handle, pattern)` returns a **`StringList`** handle.

- If `pattern` contains **`[*]`**, the prefix before it must resolve to an **array**; for each element, the suffix path after `[*]` is read (if any), and values are collected as strings.
- Without `[*]`, the pattern is a normal path; the result is **one** string in a list (or one empty string if missing).

Example: `items[*].id` collects each `id` field from `items`.

## `JSON.TOCSV`

Exports a JSON **array of objects** to CSV text. Root must be an array, or pass a **path** to an array sub-value. Header row is the union of keys (sorted). Each row is one object; missing fields become empty cells.

## Integration

- Use **`JSON.TOCSV`** with **`CSV.FROMSTRING`** / **`CSV.TOJSON`** for round trips (see [`testdata/data_integration_test.mb`](../../testdata/data_integration_test.mb)).

## Example

```basic
j = JSON.PARSESTRING("{\"player\":{\"hp\":10}}")
PRINT(JSON.GETINT(j, "player.hp"))
JSON.FREE(j)
```

## See also

- [CSV.md](CSV.md), [TABLE.md](TABLE.md), [DATABASE.md](DATABASE.md)
