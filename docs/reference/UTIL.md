# Util Commands

Cross-platform path helpers, file I/O, directory queries, and drag-and-drop support.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Use `UTIL.*` for namespaced file/path helpers. Global aliases (`FILEEXISTS`, `READALLTEXT`, etc.) are documented in [FILE.md](FILE.md). Both map to the same runtime.

---

## Paths and metadata

### `UTIL.LOADTEXT(path)` / `SAVETEXT`
Reads or writes an entire file as a UTF-8 string.

- **Returns**: (String) For `LOADTEXT`.

---

### `UTIL.FILEEXISTS(path)` / `ISDIR`
Checks for the existence or type of a file system entry.

- **Returns**: (Boolean)

---

### `UTIL.GETFILENAME(path)` / `GETFILEEXT`
Extracts components from a file path string.

- **Returns**: (String)

---

### `UTIL.GETDIRFILES(path)`
Returns a JSON array of all files in a directory.

- **Returns**: (String) JSON-formatted array.

---

### `UTIL.MAKEDIRECTORY(path)` / `CHANGEDIR`
Directory manipulation and process control.

- **Returns**: (Boolean) Success status.

---

## Drag and drop (Raylib / CGO)

When the window layer is available (`CGO` build):

| Command | Role |
|---------|------|
| `Util.IsFileDropped()` | `TRUE` if the user dropped files onto the window this frame. |
| `Util.GetDroppedFiles()` | Returns a **JSON array string** of file paths (Raylib `LoadDroppedFiles`); clears the internal drop list. |
| `Util.ClearDroppedFiles()` | Clears dropped-file state without returning paths. |

Non-CGO builds still register these names but behavior follows the stub implementation.

---

## Full Example

```basic
; Read a text file, modify it, and save
text = UTIL.LOADTEXT("notes.txt")
PRINT "File contents: " + text

UTIL.SAVETEXT("notes_backup.txt", text)
PRINT "Backup saved."

IF UTIL.ISFILENAMEVALID("my_file.txt") THEN
    PRINT "Valid file name."
ENDIF
```

---

## Extended Command Reference

### File & directory operations

| Command | Description |
|--------|-------------|
| `UTIL.COPYFILE(src, dst)` | Copy file from `src` to `dst`. |
| `UTIL.MOVEFILE(src, dst)` | Move/rename a file. |
| `UTIL.RENAMEFILE(src, dst)` | Alias of `UTIL.MOVEFILE`. |
| `UTIL.DELETEFILE(path)` | Delete a file. |
| `UTIL.CREATEDIRECTORY(path)` | Create directory (and parents). |
| `UTIL.DELETEDIR(path)` | Delete an empty directory. |
| `UTIL.GETDIRS(path)` | Returns array of subdirectory names in `path`. |

### File metadata

| Command | Description |
|--------|-------------|
| `UTIL.GETFILESIZE(path)` | Returns file size in bytes. |
| `UTIL.GETFILEMODTIME(path)` | Returns last-modified timestamp (Unix seconds). |
| `UTIL.GETFILEPATH(path)` | Returns directory part of a path. |
| `UTIL.GETFILENAMENOEXT(path)` | Returns filename without extension. |

## See also

- [FILE.md](FILE.md) — low-level read/write streams
- [JSON.md](JSON.md) — `JSON.LOADFILE` / `JSON.SAVEFILE`
