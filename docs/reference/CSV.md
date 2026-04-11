# CSV tables (`CSV.*`)

The **`CSV`** namespace loads, saves, and inspects **tabular text** using Go’s [`encoding/csv`](https://pkg.go.dev/encoding/csv). Data lives on the VM heap as a **`CSV`** handle (`TagCSV`): each row is a slice of string fields.

## Why this exists

Games and tools often exchange **spreadsheets** (loot tables, localisation, balance). CSV is a simple, diff-friendly interchange format. moonBASIC keeps **everything as strings** at the cell level so round-trips do not silently change types; convert to numbers in script when needed.

### `CSV.Load(path)`
Reads a CSV file from disk and returns a new **handle**.

### `CSV.Save(handle, path)`
Writes the CSV table to disk as UTF-8 text.

### `CSV.Free(handle)`
Releases the heap object and frees memory.

---

### `CSV.RowCount(handle)` / `CSV.ColCount(handle)`
Returns the number of rows or columns. If the table is empty, column count is 0; otherwise it follows the first row's width. Note that indexing is **1-based**.

### `CSV.Get(handle, row, col)`
Returns the cell string at the specified **1-based** index (matching typical BASIC grids).

### `CSV.Set(handle, row, col, value)`
Sets the cell string at the specified **1-based** index.

---

### `CSV.ToJSON(handle)`
Converts the CSV table to a JSON array handle. The first row is treated as header names; each following row becomes one object (`header -> cell string`).

Use this with **`JSON.*`** for structured data (`CSV.TOJSON` → `JSON.TOCSV` for uniform arrays of objects).

## Example

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
