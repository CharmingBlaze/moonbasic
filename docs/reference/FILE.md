# File I/O Commands

Commands for reading from and writing to files, and for managing the file system.

## Core Workflow

For simple reads/writes use `FILE.READALLTEXT` / `FILE.WRITEALLTEXT`. For streaming, open a handle with `FILE.OPEN`, write with `FILE.WRITE` / `FILE.WRITELINE`, and close with `FILE.CLOSE`. For path helpers see [UTIL.md](UTIL.md).

---

### `FILE.EXISTS(path)` 
Returns `TRUE` if the specified file exists on disk.

---

### `FILE.READALLTEXT(path)` 
Reads the entire content of a text file into a single string.

---

### `FILE.WRITEALLTEXT(path, content)` 
Writes an entire string to a file, overwriting it if it already exists.

---

## Advanced File Operations

### `FILE.OPEN(path, mode)` 
Opens a file and returns a **file handle**. `mode` can be `"r"` (read), `"w"` (write), or `"a"` (append).

---

### `FILE.CLOSE(handle)` 
Closes an open file handle and releases its resources.

---

### `FILE.WRITELINE(handle, content)` 
Writes a string to an open file followed by a newline character.

---

### `FILE.WRITE(handle, content)` 
Writes a string to an open file without a newline.

---

## File System Management

### `FILE.ISDIR(path)` 
Returns `TRUE` if the specified path points to a directory.

---

### `FILE.DELETE(path)` 
Deletes a file or an empty directory from the system.

---

### `FILE.COPY(source, dest)` 
Copies a file from the source path to the destination path.

---

### `FILE.MOVE(source, dest)` 
Moves or renames a file or directory.

---

### `FILE.MAKEDIR(path)` 
Creates a new directory.

---

### `FILE.GETDIR()` 
Returns the current working directory.

---

### `FILE.SETDIR(path)` 
Sets the current working directory.

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
