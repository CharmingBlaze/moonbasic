# CSV & Database Commands

Tabular data import/export (CSV) and embedded SQLite database access (DB/ROWS).

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**CSV** — Load a CSV file or parse a string, query rows and columns, modify cells, and export back to string or file.

**Database** — Open an SQLite file, execute statements, iterate result rows, and use transactions for batch writes.

---

### `CSV.LOAD(filePath)` 

Loads a CSV file and returns a handle to the parsed table.

---

### `CSV.FROMSTRING(csvText)` 

Parses a CSV-formatted string and returns a table handle.

---

### `CSV.SAVE(csvHandle, filePath)` 

Writes the CSV table to a file on disk.

---

### `CSV.FREE(csvHandle)` 

Frees the CSV table from memory.

---

### `CSV.ROWCOUNT(csvHandle)` 

Returns the number of rows in the table.

---

### `CSV.COLCOUNT(csvHandle)` 

Returns the number of columns in the table.

---

### `CSV.GET(csvHandle, row, col)` 

Returns the string value at the given row and column (0-based).

---

### `CSV.SET(csvHandle, row, col, value)` 

Sets the cell at the given row and column to a string value.

---

### `CSV.TOSTRING(csvHandle)` 

Returns the entire table as a CSV-formatted string.

---

### `CSV.TOJSON(csvHandle)` 

Converts the CSV table to a JSON array-of-objects handle (column headers become keys).

---

### `DB.OPEN(filePath)` 

Opens or creates an SQLite database file. Returns a database handle.

---

### `DB.CLOSE(dbHandle)` 

Closes an open database connection.

---

### `DB.EXEC(dbHandle, sql)` 

Executes a SQL statement that returns no rows (INSERT, UPDATE, DELETE, CREATE TABLE, etc.).

---

### `DB.QUERY(dbHandle, sql)` 

Executes a SQL SELECT and returns a rows handle for iteration.

---

### `DB.QUERYJSON(dbHandle, sql)` 

Executes a SQL SELECT and returns the entire result set as a JSON array handle.

---

### `DB.PREPARE(dbHandle, sql)` 

Prepares a parameterised SQL statement. Returns a statement handle.

---

### `DB.STMTEXEC(stmtHandle, param)` 

Executes a prepared statement with the given parameter value.

---

### `DB.STMTCLOSE(stmtHandle)` 

Closes a prepared statement and frees its resources.

---

### `DB.LASTINSERTID(dbHandle)` 

Returns the row ID of the most recent INSERT.

---

### `DB.CHANGES(dbHandle)` 

Returns the number of rows affected by the last INSERT, UPDATE, or DELETE.

---

### `DB.ISOPEN(dbHandle)` 

Returns `TRUE` if the database connection is still open.

---

### `DB.BEGIN(dbHandle)` 

Begins a transaction.

---

### `DB.COMMIT(dbHandle)` 

Commits the current transaction.

---

### `DB.ROLLBACK(dbHandle)` 

Rolls back the current transaction.

---

### `ROWS.NEXT(rowsHandle)` 

Advances to the next row in a query result. Returns `TRUE` if a row is available.

---

### `ROWS.GETINT(rowsHandle, colIndex)` 

Returns the integer value of the given column in the current row (0-based).

---

### `ROWS.GETFLOAT(rowsHandle, colIndex)` 

Returns the float value of the given column in the current row.

---

### `ROWS.GETSTRING(rowsHandle, colIndex)` 

Returns the string value of the given column in the current row.

---

### `ROWS.CLOSE(rowsHandle)` 

Closes the rows handle and frees query resources.

---

## Full Example

This example loads a CSV, prints it, then stores it in a database and reads it back.

```basic
; --- CSV ---
csv = CSV.LOAD("scores.csv")
PRINT "Rows: " + STR(CSV.ROWCOUNT(csv))

FOR i = 0 TO CSV.ROWCOUNT(csv) - 1
    name  = CSV.GET(csv, i, 0)
    score = CSV.GET(csv, i, 1)
    PRINT name + " => " + score
NEXT
CSV.FREE(csv)

; --- Database ---
db = DB.OPEN("game.db")
DB.EXEC(db, "CREATE TABLE IF NOT EXISTS scores (name TEXT, score INT)")

DB.BEGIN(db)
DB.EXEC(db, "INSERT INTO scores VALUES ('Alice', 100)")
DB.EXEC(db, "INSERT INTO scores VALUES ('Bob', 85)")
DB.COMMIT(db)

rows = DB.QUERY(db, "SELECT name, score FROM scores ORDER BY score DESC")
WHILE ROWS.NEXT(rows)
    PRINT ROWS.GETSTRING(rows, 0) + ": " + STR(ROWS.GETINT(rows, 1))
WEND
ROWS.CLOSE(rows)

DB.CLOSE(db)
```
