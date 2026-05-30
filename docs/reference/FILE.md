# File I/O Commands

Commands for reading from and writing to files, and for managing the file system.

## Core Workflow

For simple reads/writes use `FILE.READALLTEXT` / `FILE.WRITEALLTEXT`. For streaming, open a handle with `FILE.OPEN`, write with `FILE.WRITE` / `FILE.WRITELINE`, and close with `FILE.CLOSE`. For path helpers see [UTIL.md](UTIL.md).

---

### `FILE.EXISTS(path)`
Returns `TRUE` if the specified file exists on disk.

- **Arguments**:
    - `path`: (String) File path.
- **Returns**: (Boolean)

---

### `FILE.READALLTEXT(path)`
Reads the entire content of a text file into a single string.

- **Returns**: (String) The file contents.

---

### `FILE.WRITEALLTEXT(path, content)`
Writes an entire string to a file.

- **Arguments**:
    - `path`: (String) Destination path.
    - `content`: (String) Text to write.
- **Returns**: (None)

---

## Advanced File Operations

### `FILE.OPEN(path, mode)`
Opens a file and returns a **file handle**.

- **Arguments**:
    - `path`: (String) File path.
    - `mode`: (String) `"r"` (read), `"w"` (write), or `"a"` (append).
- **Returns**: (Handle) The file handle.

---

### `FILE.CLOSE(handle)`
Closes an open file handle.

---

### `FILE.WRITELINE(handle, content)` / `FILE.WRITE`
Writes a string to an open file.

- **Arguments**:
    - `handle`: (Handle) The open file.
    - `content`: (String) Text to write.
- **Returns**: (Handle) The file handle (for chaining).

---

## File System Management

### `FILE.ISDIR(path)`
Returns `TRUE` if the specified path points to a directory.

---

### `FILE.DELETE(path)` / `FILE.COPY` / `FILE.MOVE`
File system mutations.

---

### `FILE.MAKEDIR(path)`
Creates a new directory.

---

### `FILE.GETDIR()` / `FILE.SETDIR(path)`
Working directory management.

---

## Full Example

```basic
log_file = "my_game_log.txt"

; Delete the old log file if it exists
IF FILEEXISTS(log_file) THEN
    DELETEFILE(log_file)
    PRINT "Deleted old log file."
ENDIF

; Write initial messages to the log
WRITEALLTEXT(log_file, "Log file created at: " + DATETIME() + "\n")

; The following would require append mode, which is partial.
; For now, we read all text, append, and write back.

; Simulate adding more log entries
current_log = READALLTEXT(log_file)
new_entry = "Player reached level 2.\n"
WRITEALLTEXT(log_file, current_log + new_entry)

current_log = READALLTEXT(log_file)
new_entry = "Player found a secret item!\n"
WRITEALLTEXT(log_file, current_log + new_entry)


PRINT "--- Final Log Content ---"
PRINT READALLTEXT(log_file)
```

---

## Extended Command Reference

### Stream I/O

| Command | Description |
|--------|-------------|
| `FILE.OPENREAD(path)` | Open file for reading; returns file handle. |
| `FILE.OPENWRITE(path)` | Open file for writing; returns file handle. |
| `FILE.READLINE(f)` | Read and return next line as string. |
| `FILE.WRITELN(f, text)` | Write `text` followed by newline. |
| `FILE.EOF(f)` / `FILE.GETEOF(f)` | Returns `TRUE` if at end-of-file. |
| `FILE.SEEK(f, pos)` | Seek to byte offset `pos`. |
| `FILE.TELL(f)` / `FILE.GETPOS(f)` | Returns current byte offset. |
| `FILE.SIZE(f)` / `FILE.GETSIZE(f)` | Returns file size in bytes. |

## See also

- [UTIL.md](UTIL.md) — `UTIL.LOADTEXT`, `UTIL.SAVETEXT`, directory helpers
- [JSON.md](JSON.md) — structured file I/O
