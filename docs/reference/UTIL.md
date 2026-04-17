# Util Commands

Cross-platform path helpers, file I/O, directory queries, and drag-and-drop support.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Use `UTIL.*` for namespaced file/path helpers. Global aliases (`FILEEXISTS`, `READALLTEXT`, etc.) are documented in [FILE.md](FILE.md). Both map to the same runtime.

---

## Paths and metadata

| Command | Returns | Notes |
|---------|---------|------|
| `Util.FileExists(path)` | bool | |
| `Util.IsDir(path)` | bool | |
| `Util.GetFileExt(path)` | string | Includes the dot (e.g. `.png`). |
| `Util.GetFileName(path)` | string | Final path segment. |
| `Util.GetFileNameNoExt(path)` | string | Base name without extension. |
| `Util.GetFilePath(path)` | string | Directory portion. |
| `Util.GetFileSize(path)` | int | `0` if stat fails. |
| `Util.GetFileModTime(path)` | int | **Unix seconds** since epoch; `0` if stat fails. |

---

## Read / write text

### `UTIL.LOADTEXT(path)` → string 

Reads the entire file as UTF-8/text (same idea as `READALLTEXT`).

---

### `UTIL.SAVETEXT(path, text)` 

Writes a file, replacing contents.

---

## Directories and working directory

| Command | Returns / behavior |
|---------|---------------------|
| `Util.GetDirFiles(dir)` | **JSON array string** of **all** entry names in the directory (files and subdirs), from `os.ReadDir`. |
| `Util.ChangeDir(path)` | **Bool** — `TRUE` if `chdir` succeeded. |
| `Util.MakeDirectory(path)` | **Bool** — `TRUE` if `MkdirAll` succeeded. |
Current working directory and subdirectory listing use the global names **`GETDIR`** and **`GETDIRS`** (same implementation in `mbutil`; see [FILE.md](FILE.md)), not `UTIL.*` prefixes in the manifest.

---

## File operations (globals)

`DELETEFILE`, `COPYFILE`, `RENAMEFILE`, `MOVEFILE`, `DELETEDIR` are registered on the same module as **`UTIL.*`** helpers but use **global** names in the manifest. See [FILE.md](FILE.md) for semantics; implementations live in `mbutil`.

---

## Validation

### `UTIL.ISFILENAMEVALID(name)` → bool 

Checks whether a file name is acceptable on the current OS (invalid characters, reserved names, etc., per Go `path/filepath` usage in the runtime).

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
