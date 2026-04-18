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

### `TABLE.CREATE(columnNames)`
Allocates a new in-memory table with the specified columns.

- **Arguments**:
    - `columnNames`: (String) Comma-separated names (e.g., "id,name,score").
- **Returns**: (Handle) The new table handle.
- **Example**:
    ```basic
    t = TABLE.CREATE("id,name")
    ```

---

### `TABLE.ADDROW(handle, v1, v2, ...)`
Appends a new row of values to the table.

- **Returns**: (Handle) The table handle.

---

### `TABLE.GET(handle, row, columnName)` / `SET`
Accesses a cell by its **1-based** row index and column name.

- **Returns**: (String) For `GET`.

---

### `TABLE.TOJSON(handle)` / `FROMJSON`
Bridges the table to a JSON array of objects.

- **Returns**: (Handle) A new JSON handle.

## Full Example

```basic
t = TABLE.CREATE("name,score")
TABLE.ADDROW(t, "ada", 10)
TABLE.ADDROW(t, "bob", 20)
PRINT(TABLE.GET(t, 1, "name"))
TABLE.FREE(t)
```

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `TABLE.MAKE(...)` | Deprecated alias of `TABLE.CREATE`. |
| `TABLE.ROWCOUNT(t)` | Returns number of rows. |
| `TABLE.COLCOUNT(t)` | Returns number of columns. |

---

## See also

- [JSON.md](JSON.md), [CSV.md](CSV.md), [DATABASE.md](DATABASE.md)
