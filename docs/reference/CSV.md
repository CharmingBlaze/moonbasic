# CSV tables (`CSV.*`)

The **`CSV`** namespace loads, saves, and inspects **tabular text** using Go’s [`encoding/csv`](https://pkg.go.dev/encoding/csv). Data lives on the VM heap as a **`CSV`** handle (`TagCSV`): each row is a slice of string fields.

## Why this exists

Games and tools often exchange **spreadsheets** (loot tables, localisation, balance). CSV is a simple, diff-friendly interchange format. moonBASIC keeps **everything as strings** at the cell level so round-trips do not silently change types; convert to numbers in script when needed.

## Lifecycle

| Command | Purpose |
|--------|---------|
| `CSV.LOAD(path)` | Read a file from disk; returns a new **`CSV`** handle. |
| `CSV.SAVE(handle, path)` | Write the table back to disk (UTF-8 text). |
| `CSV.FROMSTRING(s)` | Parse CSV text already in memory (no file). |
| `CSV.TOSTRING(handle)` | Serialize the table to a single string (newline-terminated rows). |
| `CSV.FREE(handle)` | Release the heap object. |

## Shape and indexing

- **`CSV.ROWCOUNT`** / **`CSV.COLCOUNT`**: row and column counts. If the table is empty, column count is **0**; otherwise column count follows the **first** row’s width.
- **`CSV.GET(handle, row, col)`** / **`CSV.SET(handle, row, col, val)`**: **1-based** row and column indices (first row is **1**, first column is **1**), matching typical BASIC-style grids.

## JSON bridge

| Command | Purpose |
|--------|---------|
| `CSV.TOJSON(handle)` | Builds a **`JSON`** handle whose root is an **array of objects**: row **1** is treated as **header names**; each following row becomes one object (`header → cell string`). |

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
