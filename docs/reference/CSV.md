# CSV Commands

Load, save, query, and convert tabular CSV data.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load from disk with `CSV.LOAD` or parse a string with `CSV.FROMSTRING`.
2. Read cells with `CSV.GET`, write with `CSV.SET`.
3. Convert to JSON with `CSV.TOJSON` or save with `CSV.SAVE`.
4. Free with `CSV.FREE`.

---

### `CSV.LOAD(path)` 
Reads a CSV file from disk and returns a new **handle**.

---

### `CSV.SAVE(handle, path)` 
Writes the CSV table to disk as UTF-8 text.

---

### `CSV.FREE(handle)` 
Releases the heap object and frees memory.

---

### `CSV.ROWCOUNT(handle)` / `CSV.COLCOUNT(handle)` 
Returns the number of rows or columns. If the table is empty, column count is 0; otherwise it follows the first row's width. Note that indexing is **1-based**.

---

### `CSV.GET(handle, row, col)` 
Returns the cell string at the specified **1-based** index (matching typical BASIC grids).

---

### `CSV.SET(handle, row, col, value)` 
Sets the cell string at the specified **1-based** index.

---

### `CSV.TOJSON(handle)` 
Converts the CSV table to a JSON array handle. The first row is treated as header names; each following row becomes one object (`header -> cell string`).

Use this with **`JSON.*`** for structured data (`CSV.TOJSON` → `JSON.TOCSV` for uniform arrays of objects).

## Full Example

```basic
nl = CHR(10)
h = CSV.FROMSTRING("name,hp" + nl + "hero,10")
PRINT(CSV.GET(h, 2, 1))
j = CSV.TOJSON(h)
JSON.FREE(j)
CSV.FREE(h)
```

## See also

- [JSON.md](JSON.md) — nested documents and **`JSON.TOCSV`**
- [TABLE.md](TABLE.md) — typed in-memory tables with **`TABLE.FROMCSV`** / **`TABLE.TOCSV`**
