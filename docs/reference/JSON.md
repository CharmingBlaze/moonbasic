# JSON documents (`JSON.*`)

The **`JSON`** namespace decodes and mutates JSON values using Go’s [`encoding/json`](https://pkg.go.dev/encoding/json). Values live on the VM heap as **`JSON`** handles (`TagJSON`). The root may be an **object**, **array**, or **scalar**; nested structures use **dot + bracket paths** (see below).

## Why paths instead of flat keys

Earlier builds only supported flat objects. The current implementation keeps **one decoded `interface{}` tree** per handle so you can work with **real-world files** (levels, configs, save data) without flattening.

### `JSON.Parse(path)`
Reads a UTF-8 file from disk and decodes it as JSON. Returns a **handle**.

### `JSON.ParseString(jsonString)`
Decodes JSON from a string. Returns a **handle**.

### `JSON.Free(handle)`
Releases the JSON heap object and unloads it from memory.

---

### `JSON.GetString(handle, path [, default])`
Returns the string value at the specified dot-path. Supports an optional default value if the path is missing.

### `JSON.GetInt(handle, path [, default])`
Returns the integer value at the specified dot-path.

### `JSON.GetBool(handle, path [, default])`
Returns the boolean value at the specified dot-path.

---

### `JSON.SetString(handle, path, value)`
Sets a string value at the specified path, creating intermediate objects or arrays as needed.

### `JSON.ToFile(handle, path)`
Writes the JSON object to a file on disk in a compact format.

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
