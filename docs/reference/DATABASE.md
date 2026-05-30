# Database Commands

SQLite database access: open, query, execute, transactions, and prepared statements.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Open a database with `DB.OPEN` (file path or `":memory:"`).
2. Execute DDL/DML with `DB.EXEC`, query with `DB.QUERY`.
3. Iterate results with `ROWS.NEXT` / `ROWS.GETSTRING` / `ROWS.GETINT`.
4. Close results with `ROWS.CLOSE`, database with `DB.CLOSE`.

Requires **CGO**. See also [CSV_DATABASE.md](CSV_DATABASE.md).

## Handles

| Type | Tag | Notes |
|------|-----|--------|
| Database | `TagDB` | Owns `*sql.DB`, optional active transaction, prepared statement cache. |
| Statement | `TagDBStmt` | **`DB.PREPARE`**; freeing a statement **does not** `Close` the cached `sql.Stmt` — the database owns the cache until **`DB.CLOSE`**. |
| Rows | `TagDBRows` | Result of **`DB.QUERY`**; iterate with **`ROWS.*`**. |
| Transaction | `TagDBTx` | From **`DB.BEGIN`**; **`DB.COMMIT`** / **`DB.ROLLBACK`** frees the tx handle. |

`Free()` on heap objects is **idempotent**. Closing a **`DB`** rolls back an open transaction and closes all cached statements.

### `DB.OPEN(path)`
Opens a SQLite database file (or `":memory:"`).

- **Arguments**:
    - `path`: (String) File path to the database.
- **Returns**: (Handle) The new database handle.
- **Example**:
    ```basic
    db = DB.OPEN("game.db")
    ```

---

### `DB.QUERY(db, sql [, params...])`
Executes a SELECT statement and returns a row set.

- **Arguments**:
    - `db`: (Handle) The database.
    - `sql`: (String) The query with `?` placeholders.
    - `params`: (Optional) Values to bind to placeholders.
- **Returns**: (Handle) A new rows handle.

---

### `DB.EXEC(db, sql [, params...])`
Executes a DML statement (INSERT, UPDATE, DELETE).

- **Returns**: (None)

---

### `ROWS.NEXT(rows)`
Advances to the next row in the result set.

- **Returns**: (Boolean) `TRUE` if a row is available.

---

### `ROWS.GETSTRING(rows, colName)` / `GETINT`
Reads a value from the current row.

- **Returns**: (String / Integer)

---

### `DB.CLOSE(db)` / `ROWS.CLOSE(rows)`
Releases the database or row set resources.

Freed **`TX`** handles roll back if not yet committed.

---

## Prepared statements

| Command | Purpose |
|--------|---------|
| `DB.PREPARE(db, sql)` | Returns a **`STMT`** handle; SQL text is **normalized** (trim + collapse whitespace) so identical queries share one cached `sql.Stmt`. |
| `DB.STMTEXEC(stmt, ...params)` | Executes the statement. Under an active transaction, uses `tx.Stmt`. |
| `DB.STMTCLOSE(stmt)` | Marks the handle freed; underlying `sql.Stmt` stays cached on the DB. |

## Full Example

Create a table, insert a row, query it, and iterate results.

```basic
db = DB.OPEN(":memory:")
DB.EXEC(db, "CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT, score INT)")
DB.EXEC(db, "INSERT INTO items (name, score) VALUES (?, ?)", "sword", 42)
DB.EXEC(db, "INSERT INTO items (name, score) VALUES (?, ?)", "shield", 30)

rows = DB.QUERY(db, "SELECT name, score FROM items ORDER BY score DESC")
WHILE ROWS.NEXT(rows)
    PRINT ROWS.GETSTRING(rows, "name") + " : " + STR(ROWS.GETINT(rows, "score"))
WEND
ROWS.CLOSE(rows)

DB.CLOSE(db)
```

## See also

- [TABLE.md](TABLE.md) — higher-level tables in memory (future DB bridge hooks may extend this).

Spec note: any legacy typo **`DB.ISOEPEN`** is **not** implemented; use **`DB.ISOPEN`**.
