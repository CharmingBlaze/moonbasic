# In-memory tables (`TABLE.*`)

The **`TABLE`** namespace provides a **column-oriented grid** of cells stored as Go `interface{}` values (strings, numbers, booleans) with **string column names**. It is implemented in pure Go (`runtime/tablemod`, package **`mbtable`**, heap tag **`TagTable`**).

## Why this exists

When you need **structured rows** in RAM (UI tables, rosters, scratch buffers) without SQL overhead, a **`TABLE`** handle is lighter than ad-hoc `ARRAY` juggling. Bridges to **`JSON`** and **`CSV`** let you load/save or ship data to other moonBASIC APIs.

## Creating and filling

| Command | Purpose |
|--------|---------|
| `TABLE.CREATE(cols)` | `cols` is comma-separated names, e.g. `"name,score,hp"`. |
| `TABLE.ADDROW(handle, v1, v2, …)` | One value per column (arity must match column count). |
| `TABLE.FREE(handle)` | Release the table. |

## Access

Indices are **1-based** for rows; columns are addressed by **name** (case-sensitive to the string you passed to **`CREATE`**).

| Command | Purpose |
|--------|---------|
| `TABLE.ROWCOUNT` / `TABLE.COLCOUNT` | Sizes. |
| `TABLE.GET(handle, row, col)` | Read cell; missing/empty cells return sensible defaults for strings. |
| `TABLE.SET(handle, row, col, value)` | Write cell. |

## Bridges

| Command | Purpose |
|--------|---------|
| `TABLE.TOJSON(handle)` | JSON array of objects (one object per row). |
| `TABLE.FROMJSON(handle)` | Expects a **`JSON`** root array of uniform objects; column order is sorted lexicographically by key (deterministic). |
| `TABLE.TOCSV(handle)` | First row = column names; following rows = data (**`CSV`** handle). |
| `TABLE.FROMCSV(handle)` | Inverse of **`TOCSV`** (expects header row). |

## Example

```basic
t = TABLE.CREATE("name,score")
TABLE.ADDROW(t, "ada", 10)
TABLE.ADDROW(t, "bob", 20)
PRINT(TABLE.GET(t, 1, "name"))
TABLE.FREE(t)
```

## See also

- [JSON.md](JSON.md), [CSV.md](CSV.md), [DATABASE.md](DATABASE.md)
