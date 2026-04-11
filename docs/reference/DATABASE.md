# SQLite databases (`DB.*`, `ROWS.*`)

The **`DB`** namespace opens **SQLite** databases through [`database/sql`](https://pkg.go.dev/database/sql) and **[`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3)**. This requires **CGO** at build time (same toolchain you use for Raylib). When **`CGO_ENABLED=0`** or no C compiler is available, **`runtime/dbmod`** registers **stubs** that return a clear error for every command.

## Why SQLite

SQLite gives you **durable**, **queryable** storage (saves, inventories, analytics) without a separate server process. Parameters use **`?`** placeholders and map to trailing moonBASIC arguments (see overload rows in **`commands.json`**).

## Handles

| Type | Tag | Notes |
|------|-----|--------|
| Database | `TagDB` | Owns `*sql.DB`, optional active transaction, prepared statement cache. |
| Statement | `TagDBStmt` | **`DB.PREPARE`**; freeing a statement **does not** `Close` the cached `sql.Stmt` — the database owns the cache until **`DB.CLOSE`**. |
| Rows | `TagDBRows` | Result of **`DB.QUERY`**; iterate with **`ROWS.*`**. |
| Transaction | `TagDBTx` | From **`DB.BEGIN`**; **`DB.COMMIT`** / **`DB.ROLLBACK`** frees the tx handle. |

`Free()` on heap objects is **idempotent**. Closing a **`DB`** rolls back an open transaction and closes all cached statements.

### `DB.Open(path)`
Opens a SQLite database file. Use `":memory:"` for an in-memory database. Returns a **database handle**.

### `DB.Close(db)`
Closes the database and releases all associated resources.

---

### `DB.Query(db, sql [, ...params])`
Executes a `SELECT` query and returns a **Rows** handle. Supports bound parameters for safe queries.

### `DB.Exec(db, sql [, ...params])`
Executes a SQL statement (e.g., `CREATE`, `INSERT`, `UPDATE`) with optional bound parameters.

---

### `Rows.Next(rows)`
Advances to the next row in the result set. Returns `FALSE` when no more rows are available.

### `Rows.GetString(rows, column)` / `Rows.GetInt(...)`
Returns the value of the specified column (**1-based** index) for the current row.

### `Rows.Close(rows)`
Closes the result set and frees associated memory.

---

### `DB.Begin(db)`
Starts a SQL transaction and returns a **TX** handle.

### `DB.Commit(tx)` / `DB.Rollback(tx)`
Finalizes the transaction by committing changes or rolling them back, then frees the handle.

Freed **`TX`** handles roll back if not yet committed.

## Prepared statements

| Command | Purpose |
|--------|---------|
| `DB.PREPARE(db, sql)` | Returns a **`STMT`** handle; SQL text is **normalized** (trim + collapse whitespace) so identical queries share one cached `sql.Stmt`. |
| `DB.STMTEXEC(stmt, ...params)` | Executes the statement. Under an active transaction, uses `tx.Stmt`. |
| `DB.STMTCLOSE(stmt)` | Marks the handle freed; underlying `sql.Stmt` stays cached on the DB. |

## Example

```basic
db = DB.OPEN(":memory:")
DB.EXEC(db, "CREATE TABLE t (id INTEGER PRIMARY KEY, v TEXT)")
DB.EXEC(db, "INSERT INTO t (v) VALUES (?)", "ok")
PRINT(DB.LASTINSERTID(db))
DB.CLOSE(db)
```

## See also

- [TABLE.md](TABLE.md) — higher-level tables in memory (future DB bridge hooks may extend this).

Spec note: any legacy typo **`DB.ISOEPEN`** is **not** implemented; use **`DB.ISOPEN`**.
