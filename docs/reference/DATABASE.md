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

## Core commands

| Command | Purpose |
|--------|---------|
| `DB.OPEN(path)` | `sql.Open("sqlite3", path)`. Use `":memory:"` for RAM-only tests. |
| `DB.CLOSE(handle)` | Close the DB and free resources. |
| `DB.ISOPEN(handle)` | Returns whether the handle is still valid. |
| `DB.EXEC(db, sql, ...params)` | `Exec` with optional bound parameters. |
| `DB.QUERY(db, sql, ...params)` | Returns a **`ROWS`** handle. |
| `DB.QUERYJSON(db, sql, ...params)` | Runs a query and returns a **JSON array of objects** as a string (each column name → value). Handy when you do not need row-by-row iteration. |
| `DB.LASTINSERTID(db)` | `SELECT last_insert_rowid()`. |
| `DB.CHANGES(db)` | `SELECT changes()`. |

## Row cursors (`ROWS.*`)

After **`DB.QUERY`**, call **`ROWS.NEXT(rows)`** until it returns **`FALSE`**. Column values for the current row are read with **`ROWS.GETSTRING`**, **`ROWS.GETINT`**, **`ROWS.GETFLOAT`** using **1-based** column indices. **`ROWS.CLOSE`** frees the result set.

## Transactions

| Command | Purpose |
|--------|---------|
| `DB.BEGIN(db)` | Returns a **`TX`** handle; subsequent **`DB.EXEC`** / **`DB.QUERY`** on the same DB run on the active transaction. |
| `DB.COMMIT(tx)` | Commit and free the handle. |
| `DB.ROLLBACK(tx)` | Roll back and free the handle. |

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
