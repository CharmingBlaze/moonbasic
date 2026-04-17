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

Loads a config file into the module-local store. Keys from the file become immediately queryable.

- `filePath`: Path to the config file (e.g. `"settings.ini"`).

---

### `CONFIG.SAVE(filePath)` 

Writes the current config store to disk, overwriting the file if it exists.

- `filePath`: Destination path.

---

### `CONFIG.GETINT(key)` 

Returns the integer value for `key`, or `0` if the key does not exist.

---

### `CONFIG.GETFLOAT(key)` 

Returns the float value for `key`, or `0.0` if the key does not exist.

---

### `CONFIG.GETSTRING(key)` 

Returns the string value for `key`, or `""` if the key does not exist.

---

### `CONFIG.GETBOOL(key)` 

Returns the boolean value for `key`, or `FALSE` if the key does not exist.

---

### `CONFIG.SETINT(key, value)` 

Sets an integer value in the store.

---

### `CONFIG.SETFLOAT(key, value)` 

Sets a float value in the store.

---

### `CONFIG.SETSTRING(key, value)` 

Sets a string value in the store.

---

### `CONFIG.SETBOOL(key, value)` 

Sets a boolean value in the store.

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
