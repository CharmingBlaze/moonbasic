# Utility and filesystem (`UTIL.*`)

Cross-platform **path and file helpers** under the `UTIL.` namespace. Many operations mirror the global names documented in [FILE.md](FILE.md) (`FILEEXISTS`, `READALLTEXT`, …); `UTIL.*` exists for consistent `Module.Command` style and for a few **window-only** helpers.

Implemented in `runtime/mbutil`.

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

### `Util.LoadText(path)` → string

Reads the entire file as UTF-8/text (same idea as `READALLTEXT`).

### `Util.SaveText(path, text)`

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

### `Util.IsFileNameValid(name)` → bool

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
