# File I/O Commands

Commands for reading from and writing to files, and for managing the file system.

## Core Concepts

-   **Simple I/O**: For basic needs, `READALLTEXT$` and `WRITEALLTEXT` are the easiest way to read or write entire files at once.
-   **Advanced I/O**: For more control (e.g., reading a large file line-by-line), you must `OPENFILE`, perform your read/write operations, and then `CLOSEFILE`.

---

## Simple File Operations

### `READALLTEXT$(filePath$)`

Reads the entire content of a text file into a single string.

### `WRITEALLTEXT(filePath$, content$)`

Writes an entire string to a file, creating it if it doesn't exist and overwriting it if it does.

---

## Advanced File Operations

### `OPENFILE(filePath$, mode$)`

**[PARTIAL]** Opens a file and returns a handle. `mode$` can be `"r"` (read), `"w"` (write), or `"a"` (append).

### `CLOSEFILE(fileHandle)`

**[PARTIAL]** Closes a file handle that was opened with `OPENFILE`.

### `FILE.WRITELN(fileHandle, content$)`

Writes a string to an open file, followed by a newline character.

---

## File System Management

### `FILEEXISTS(filePath$)` / `DIREXISTS(dirPath$)`

Returns `TRUE` if the specified file or directory exists.

### `DELETEFILE(filePath$)` / `DELETEDIR(dirPath$)`

Deletes a file or an empty directory.

### `COPYFILE(source$, dest$)` / `MOVEFILE(source$, dest$)`

Copies or moves a file.

### `MAKEDIR(dirPath$)`

Creates a new directory.

### `GETDIR$()` / `SETDIR(dirPath$)`

Gets or sets the current working directory.

---

## Full Example: Creating and Managing a Log File

```basic
log_file$ = "my_game_log.txt"

; Delete the old log file if it exists
IF FILEEXISTS(log_file$) THEN
    DELETEFILE(log_file$)
    PRINT "Deleted old log file."
ENDIF

; Write initial messages to the log
WRITEALLTEXT(log_file$, "Log file created at: " + DATETIME$() + "\n")

; The following would require append mode, which is partial.
; For now, we read all text, append, and write back.

; Simulate adding more log entries
current_log$ = READALLTEXT$(log_file$)
new_entry$ = "Player reached level 2.\n"
WRITEALLTEXT(log_file$, current_log$ + new_entry$)

current_log$ = READALLTEXT$(log_file$)
new_entry$ = "Player found a secret item!\n"
WRITEALLTEXT(log_file$, current_log$ + new_entry$)


PRINT "--- Final Log Content ---"
PRINT READALLTEXT$(log_file$)
```
