# DB Commands

`DB.*` is the shorthand alias namespace for SQLite database commands. All commands are identical to those in [DATABASE.md](DATABASE.md) — use whichever prefix you prefer.

## Quick Reference

| Command | Description |
|---|---|
| `DB.OPEN(path)` | Open or create a SQLite file |
| `DB.CLOSE(db)` | Close the database |
| `DB.EXEC(db, sql)` | Execute a statement (no results) |
| `DB.QUERY(db, sql)` | Run a SELECT; returns a rows handle |
| `DB.QUERYJSON(db, sql)` | Run a SELECT; returns JSON string handle |
| `DB.PREPARE(db, sql)` | Prepare a statement |
| `DB.STMTEXEC(stmt, args...)` | Execute a prepared statement |
| `DB.STMTCLOSE(stmt)` | Close a prepared statement |
| `DB.LASTINSERTID(db)` | Last auto-increment row id |
| `DB.CHANGES(db)` | Rows affected by last write |
| `DB.ISOPEN(db)` | `TRUE` if database is open |
| `DB.BEGIN(db)` | Begin a transaction |
| `DB.COMMIT(db)` | Commit a transaction |
| `DB.ROLLBACK(db)` | Rollback a transaction |

Read rows from `DB.QUERY` with [ROWS.md](ROWS.md).

## See also

- [DATABASE.md](DATABASE.md) — full SQL database documentation
- [ROWS.md](ROWS.md) — iterating query results
