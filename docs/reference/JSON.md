# JSON Commands

Decode, query, mutate, and serialize JSON documents from files or strings.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Parse from a file with `JSON.PARSE` or from a string with `JSON.PARSESTRING`.
2. Read values with `JSON.GETSTRING`, `JSON.GETINT`, `JSON.GETBOOL` using dot-bracket paths.
3. Mutate with `JSON.SETSTRING`.
4. Write back with `JSON.TOFILE`.
5. Free with `JSON.FREE`.

Nested structures use dot + bracket paths (e.g. `"player.inventory[0].name"`).

### `JSON.PARSE(path)` 
Reads a UTF-8 file from disk and decodes it as JSON. Returns a **handle**.

---

### `JSON.PARSESTRING(jsonString)` 
Decodes JSON from a string. Returns a **handle**.

---

### `JSON.FREE(handle)` 
Releases the JSON heap object.

---

### `JSON.GETSTRING(handle, path, default)` 
Returns the string value at the specified dot-path. Optional `default` if missing.

---

### `JSON.GETINT(handle, path, default)` 
Returns the integer value at the specified dot-path.

---

### `JSON.GETBOOL(handle, path, default)` 
Returns the boolean value at the specified dot-path.

---

### `JSON.SETSTRING(handle, path, value)` 
Sets a string value at the specified path, creating intermediates as needed.

---

### `JSON.TOFILE(handle, path)` 
Writes the JSON object to a file on disk.

---

## `JSON.QUERY` (minimal)

`JSON.QUERY(handle, pattern)` returns a **`StringList`** handle.

- If `pattern` contains **`[*]`**, the prefix before it must resolve to an **array**; for each element, the suffix path after `[*]` is read (if any), and values are collected as strings.
- Without `[*]`, the pattern is a normal path; the result is **one** string in a list (or one empty string if missing).

Example: `items[*].id` collects each `id` field from `items`.

## `JSON.TOCSV`

Exports a JSON **array of objects** to CSV text. Root must be an array, or pass a **path** to an array sub-value. Header row is the union of keys (sorted). Each row is one object; missing fields become empty cells.

## Integration

- Use **`JSON.TOCSV`** with **`CSV.FROMSTRING`** / **`CSV.TOJSON`** for round trips (see [`testdata/data_integration_test.mb`](../../testdata/data_integration_test.mb)).

## Full Example

```basic
j = JSON.PARSESTRING("{\"player\":{\"hp\":10}}")
PRINT(JSON.GETINT(j, "player.hp"))
JSON.FREE(j)
```

## See also

- [CSV.md](CSV.md), [TABLE.md](TABLE.md), [DATABASE.md](DATABASE.md)
