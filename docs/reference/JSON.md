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

### `JSON.PARSE(path)` / `PARSESTRING`
Decodes a JSON document from a file or string.

- **Returns**: (Handle) The new JSON handle.
- **Example**:
    ```basic
    j = JSON.PARSESTRING("{\"hero\": {\"hp\": 100}}")
    ```

---

### `JSON.GETSTRING(handle, path [, default])` / `GETINT` / `GETBOOL`
Reads a value at a dot-path (e.g., `"player.inventory[0].id"`).

- **Arguments**:
    - `handle`: (Handle) The JSON document.
    - `path`: (String) The query path.
    - `default`: (Optional) Fallback value if missing.
- **Returns**: (String / Integer / Boolean) The resolved value.

---

### `JSON.SETSTRING(handle, path, value)`
Mutates a value at the specified path.

---

### `JSON.TOFILE(handle, path)` / `TOSTRING`
Serializes the JSON document back to a file or string.

---

### `JSON.FREE(handle)`
Releases the JSON heap object.

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

## Extended Command Reference

### Creation

| Command | Description |
|--------|-------------|
| `JSON.CREATE()` | Create an empty JSON object handle. |
| `JSON.MAKEARRAY()` | Create an empty JSON array handle. |

### Reading

| Command | Description |
|--------|-------------|
| `JSON.GETFLOAT(j, key)` | Read a number value as float. |
| `JSON.GETARRAY(j, key)` | Returns a JSON array handle for `key`. |
| `JSON.GETOBJECT(j, key)` | Returns a JSON object handle for `key`. |
| `JSON.HAS(j, key)` | Returns `TRUE` if `key` exists. |
| `JSON.LEN(j)` | Returns number of keys (object) or elements (array). |
| `JSON.TYPE(j, key)` | Returns type string: `"number"`, `"string"`, `"bool"`, `"array"`, `"object"`, `"null"`. |

### Writing

| Command | Description |
|--------|-------------|
| `JSON.SETINT(j, key, v)` | Set integer value. |
| `JSON.SETFLOAT(j, key, v)` | Set float value. |
| `JSON.SETBOOL(j, key, v)` | Set boolean value. |
| `JSON.SETNULL(j, key)` | Set value to `null`. |
| `JSON.APPEND(j, value)` | Append value to a JSON array. |
| `JSON.DELETE(j, key)` | Delete a key from object. |
| `JSON.CLEAR(j)` | Remove all keys/elements. |

### File I/O

| Command | Description |
|--------|-------------|
| `JSON.LOADFILE(path)` | Load and parse a JSON file, returns handle. |
| `JSON.SAVEFILE(j, path)` | Serialise and write JSON to file (compact). |
| `JSON.TOFILEPRETTY(j, path)` | Serialise and write JSON to file (pretty-printed). |
| `JSON.MINIFY(j)` | Returns compact JSON string. |
| `JSON.PRETTY(j)` | Returns pretty-printed JSON string. |

---

## See also

- [CSV.md](CSV.md), [TABLE.md](TABLE.md), [DATABASE.md](DATABASE.md)
