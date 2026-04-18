# CSV Commands

Load, save, query, and convert tabular CSV data.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load from disk with `CSV.LOAD` or parse a string with `CSV.FROMSTRING`.
2. Read cells with `CSV.GET`, write with `CSV.SET`.
3. Convert to JSON with `CSV.TOJSON` or save with `CSV.SAVE`.
4. Free with `CSV.FREE`.

---

### `CSV.LOAD(path)` / `FROMSTRING`
Loads or parses CSV data into a handle.

- **Returns**: (Handle) The new CSV handle.
- **Example**:
    ```basic
    h = CSV.LOAD("data.csv")
    ```

---

### `CSV.GET(handle, row, col)` / `SET`
Accesses cell values using **1-based** indexing.

- **Arguments**:
    - `handle`: (Handle) The CSV table.
    - `row, col`: (Integer) 1-based coordinates.
- **Returns**: (String) For `GET`.

---

### `CSV.ROWCOUNT(handle)` / `COLCOUNT`
Returns the dimensions of the table.

- **Returns**: (Integer)

---

### `CSV.TOJSON(handle)`
Converts the table to a JSON array of objects (using the first row as headers).

- **Returns**: (Handle) A new JSON handle.

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
