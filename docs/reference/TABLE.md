# Table Commands

Column-oriented in-memory grids for structured rows, with JSON and CSV bridges.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a table with `TABLE.CREATE`, listing column names.
2. Add rows with `TABLE.ADDROW`.
3. Read/write cells with `TABLE.GET` / `TABLE.SET`.
4. Bridge to JSON or CSV with `TABLE.TOJSON` / `TABLE.FROMJSON` / `TABLE.TOCSV` / `TABLE.FROMCSV`.
5. Free with `TABLE.FREE`.

---

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

## Full Example

```basic
t = TABLE.CREATE("name,score")
TABLE.ADDROW(t, "ada", 10)
TABLE.ADDROW(t, "bob", 20)
PRINT(TABLE.GET(t, 1, "name"))
TABLE.FREE(t)
```

## See also

- [JSON.md](JSON.md), [CSV.md](CSV.md), [DATABASE.md](DATABASE.md)
