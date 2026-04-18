# Rows Commands

Iterate query results from `DB.QUERY` / `DATABASE.QUERY`. `ROWS.*` is the cursor-based row iterator for SQLite result sets.

## Core Workflow

1. `rows = DB.QUERY(db, sql)` — execute a SELECT.
2. `WHILE ROWS.NEXT(rows)` — advance to the next row.
3. `ROWS.GETINT(rows, col)` / `ROWS.GETFLOAT(rows, col)` / `ROWS.GETSTRING(rows, col)` — read columns (0-based index).
4. `ROWS.CLOSE(rows)` — free the result set.

---

## Commands

### `ROWS.NEXT(rows)` 

Advances to the next row. Returns `TRUE` if a row is available, `FALSE` when exhausted.

---

### `ROWS.GETINT(rows, columnIndex)` 

Returns the integer value of column `columnIndex` in the current row (0-based).

---

### `ROWS.GETFLOAT(rows, columnIndex)` 

Returns the float value of the column.

---

### `ROWS.GETSTRING(rows, columnIndex)` 

Returns the string value of the column.

---

### `ROWS.CLOSE(rows)` 

Closes the row iterator and frees resources.

---

## Full Example

```basic
db = DB.OPEN("game.db")
DB.EXEC(db, "CREATE TABLE IF NOT EXISTS scores (name TEXT, score INT)")
DB.EXEC(db, "INSERT INTO scores VALUES ('Alice', 1200)")
DB.EXEC(db, "INSERT INTO scores VALUES ('Bob', 950)")

rows = DB.QUERY(db, "SELECT name, score FROM scores ORDER BY score DESC")
WHILE ROWS.NEXT(rows)
    PRINT ROWS.GETSTRING(rows, 0) + ": " + STR(ROWS.GETINT(rows, 1))
WEND
ROWS.CLOSE(rows)

DB.CLOSE(db)
```

---

## See also

- [DB.md](DB.md) — shorthand database namespace
- [DATABASE.md](DATABASE.md) — full database documentation
