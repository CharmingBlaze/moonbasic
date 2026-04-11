# File I/O Commands

Commands for reading from and writing to files, and for managing the file system.

## Core Concepts

-   **Simple I/O**: For basic needs, `ReadAllText` and `WriteAllText` are the easiest way to read or write entire files at once.

---

### `File.Exists(path)`
Returns `TRUE` if the specified file exists on disk.

### `File.ReadAllText(path)`
Reads the entire content of a text file into a single string.

### `File.WriteAllText(path, content)`
Writes an entire string to a file, overwriting it if it already exists.

---

## Advanced File Operations

### `File.Open(path, mode)`
Opens a file and returns a **file handle**. `mode` can be `"r"` (read), `"w"` (write), or `"a"` (append).

### `File.Close(handle)`
Closes an open file handle and releases its resources.

### `File.WriteLine(handle, content)`
Writes a string to an open file followed by a newline character.

### `File.Write(handle, content)`
Writes a string to an open file without a newline.

---

## File System Management

### `File.IsDir(path)`
Returns `TRUE` if the specified path points to a directory.

### `File.Delete(path)`
Deletes a file or an empty directory from the system.

### `File.Copy(source, dest)`
Copies a file from the source path to the destination path.

### `File.Move(source, dest)`
Moves or renames a file or directory.

### `File.MakeDir(path)`
Creates a new directory.

### `File.GetDir()`
Returns the current working directory.

### `File.SetDir(path)`
Sets the current working directory.

---

## Full Example: Creating and Managing a Log File

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
