# Config Commands

Key–value settings store backed by an INI-style file on disk.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load an existing config file with `CONFIG.LOAD`.
2. Read and write typed values with `CONFIG.GETINT` / `CONFIG.SETINT` (and float, string, bool variants).
3. Save changes back to disk with `CONFIG.SAVE`.

For larger or structured data, prefer `JSON.*` or `FILE.*`.

---

### `CONFIG.LOAD(filePath)`
Loads a config file into the module-local store.

- **Arguments**:
    - `filePath`: (String) Path to the `.ini` file.
- **Returns**: (Boolean) `TRUE` if loaded successfully.

---

### `CONFIG.SAVE(filePath)`
Writes the current config store to disk.

- **Returns**: (Boolean) `TRUE` if saved successfully.

---

### `CONFIG.GETINT(key)` / `GETFLOAT` / `GETSTRING` / `GETBOOL`
Returns the typed value for `key`, or a default (0, "", FALSE) if missing.

- **Arguments**:
    - `key`: (String) The setting name.
- **Returns**: (Integer/Float/String/Boolean)

---

### `CONFIG.SETINT(key, value)` / `SETFLOAT` / `SETSTRING` / `SETBOOL`
Sets a typed value in the store.

- **Arguments**:
    - `key`: (String) The setting name.
    - `value`: (Any) The value to store.
- **Returns**: (None)

---

### `CONFIG.HAS(key)`
Returns `TRUE` if the key exists in the store.

---

### `CONFIG.DELETE(key)`
Removes a key from the store.

---

## Full Example

This example loads settings, reads a volume value, changes it, and saves.

```basic
CONFIG.LOAD("settings.ini")

; Read existing volume or default to 80
IF CONFIG.HAS("volume")
    vol = CONFIG.GETINT("volume")
ELSE
    vol = 80
END IF

PRINT "Volume: " + STR(vol)

; Bump volume and save
CONFIG.SETINT("volume", vol + 5)
CONFIG.SETBOOL("fullscreen", TRUE)
CONFIG.SAVE("settings.ini")
PRINT "Settings saved."
```
