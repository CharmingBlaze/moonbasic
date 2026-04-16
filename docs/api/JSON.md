# JSON Commands

Commands for parsing, creating, querying, and serializing JSON data. moonBASIC provides a full in-memory JSON DOM that supports objects, arrays, strings, numbers, booleans, and null values. Useful for config files, save data, level definitions, and web API responses.

## Core Concepts

- **JSON handle** — Every JSON object or array is a heap handle. Must be freed when done.
- **Path keys** — Nested values are accessed with dot-separated keys (e.g., `"player.position.x"`).
- **Type-safe getters** — Use `GetString`, `GetInt`, `GetFloat`, `GetBool` to extract typed values.
- **Mutable** — JSON handles can be modified in-place with `Set*`, `Delete`, `Append`, `Clear`.

---

## Parsing

### `JSON.Parse(filePath)`

Loads and parses a JSON file from disk. Returns a handle to the root object or array.

- `filePath` (string) — Path to a `.json` file.

**Returns:** `handle`

**How it works:** Reads the entire file, parses it into an in-memory tree structure, and allocates a heap handle pointing to the root node.

```basic
config = JSON.Parse("config.json")
```

---

### `JSON.ParseString(jsonString)`

Parses a JSON string directly (not from a file).

- `jsonString` (string) — A valid JSON string.

**Returns:** `handle`

```basic
data = JSON.ParseString("{""name"": ""Player 1"", ""score"": 100}")
```

---

### `JSON.LoadFile(filePath)`

Alias for `JSON.Parse`. Loads and parses a JSON file.

---

## Creating

### `JSON.Create()`

Creates a new empty JSON object `{}`.

**Returns:** `handle`

```basic
saveData = JSON.Create()
JSON.SetString(saveData, "name", "Hero")
JSON.SetInt(saveData, "level", 5)
JSON.SetFloat(saveData, "health", 87.5)
JSON.SetBool(saveData, "hasKey", TRUE)
```

---

### `JSON.MakeArray()`

Creates a new empty JSON array `[]`.

**Returns:** `handle`

```basic
inventory = JSON.MakeArray()
JSON.Append(inventory, "sword")
JSON.Append(inventory, "shield")
JSON.Append(inventory, "potion")
```

---

### `JSON.Free(handle)`

Frees a JSON handle and all its children from memory.

- `handle` (handle) — JSON handle.

```basic
JSON.Free(config)
```

---

## Querying

### `JSON.Has(handle, key)`

Returns `TRUE` if the JSON object contains the given key.

- `handle` (handle) — JSON object handle.
- `key` (string) — Key to check.

**Returns:** `bool`

```basic
IF JSON.Has(config, "debug") THEN
    debugMode = JSON.GetBool(config, "debug")
ENDIF
```

---

### `JSON.Type(handle, key)`

Returns the type of a value as a string: `"string"`, `"number"`, `"object"`, `"array"`, `"bool"`, `"null"`.

- `handle` (handle) — JSON handle.
- `key` (string) — Key to check.

**Returns:** `string`

---

### `JSON.Len(handle)`

Returns the number of elements in a JSON array or keys in an object.

- `handle` (handle) — JSON handle.

**Returns:** `int`

```basic
count = JSON.Len(inventory)
FOR i = 0 TO count - 1
    item = JSON.GetString(inventory, STR(i))
    PRINT item
NEXT
```

---

### `JSON.Keys(handle)`

Returns the keys of a JSON object. Access individual keys by index.

- `handle` (handle) — JSON object handle.

**Returns:** Result accessible by index.

---

## Typed Getters

### `JSON.GetString(handle, key)`

Returns the string value at the given key.

- `handle` (handle) — JSON handle.
- `key` (string) — Key or array index as string.

**Returns:** `string`

```basic
name = JSON.GetString(config, "player.name")
```

---

### `JSON.GetInt(handle, key)`

Returns the integer value at the given key.

**Returns:** `int`

```basic
level = JSON.GetInt(config, "player.level")
```

---

### `JSON.GetFloat(handle, key)`

Returns the float value at the given key.

**Returns:** `float`

```basic
speed = JSON.GetFloat(config, "physics.gravity")
```

---

### `JSON.GetBool(handle, key)`

Returns the boolean value at the given key.

**Returns:** `bool`

---

### `JSON.GetArray(handle, key)`

Returns a JSON array handle nested inside an object.

**Returns:** `handle`

```basic
enemies = JSON.GetArray(levelData, "enemies")
count = JSON.Len(enemies)
```

---

### `JSON.GetObject(handle, key)`

Returns a JSON object handle nested inside another object.

**Returns:** `handle`

```basic
playerData = JSON.GetObject(saveFile, "player")
health = JSON.GetFloat(playerData, "health")
```

---

## Typed Setters

### `JSON.SetString(handle, key, value)`

Sets a string value at the given key. Creates the key if it doesn't exist.

- `handle` (handle) — JSON handle.
- `key` (string) — Key name.
- `value` (string) — Value to set.

---

### `JSON.SetInt(handle, key, value)`

Sets an integer value.

---

### `JSON.SetFloat(handle, key, value)`

Sets a float value.

---

### `JSON.SetBool(handle, key, value)`

Sets a boolean value.

---

### `JSON.SetNull(handle, key)`

Sets a value to `null`.

---

### `JSON.Delete(handle, key)`

Removes a key and its value from a JSON object.

---

### `JSON.Clear(handle)`

Removes all keys from an object or all elements from an array.

---

### `JSON.Append(handle, value)`

Appends a value to a JSON array.

- `handle` (handle) — JSON array handle.
- `value` (any) — Value to append.

```basic
scores = JSON.MakeArray()
JSON.Append(scores, 100)
JSON.Append(scores, 250)
JSON.Append(scores, 50)
```

---

## Querying (Advanced)

### `JSON.Query(handle, path)`

Performs a JSONPath-style query on the data and returns the result.

- `handle` (handle) — JSON handle.
- `path` (string) — Query path.

**Returns:** Result handle or value.

---

## Serialization

### `JSON.ToString(handle)`

Converts a JSON handle to a compact JSON string.

- `handle` (handle) — JSON handle.

**Returns:** `string`

```basic
PRINT JSON.ToString(saveData)
; Output: {"name":"Hero","level":5,"health":87.5,"hasKey":true}
```

---

### `JSON.Pretty(handle)`

Converts a JSON handle to a pretty-printed JSON string with indentation.

**Returns:** `string`

---

### `JSON.Minify(handle)`

Returns a minified (no whitespace) JSON string.

**Returns:** `string`

---

### `JSON.ToFile(handle, filePath)` / `JSON.SaveFile(handle, filePath)`

Saves a JSON handle to a file as compact JSON.

- `handle` (handle) — JSON handle.
- `filePath` (string) — Output file path.

```basic
JSON.SaveFile(saveData, "save.json")
```

---

### `JSON.ToFilePretty(handle, filePath)`

Saves a JSON handle to a file as pretty-printed JSON.

---

### `JSON.ToCSV(handle)`

Converts a JSON array of objects to CSV format.

**Returns:** `string`

---

## Full Example

A complete save/load system using JSON.

```basic
; === Save Game ===
FUNCTION SaveGame(fileName, playerName, level, health, x, y, z)
    save = JSON.Create()
    JSON.SetString(save, "name", playerName)
    JSON.SetInt(save, "level", level)
    JSON.SetFloat(save, "health", health)

    ; Save position as nested object
    pos = JSON.Create()
    JSON.SetFloat(pos, "x", x)
    JSON.SetFloat(pos, "y", y)
    JSON.SetFloat(pos, "z", z)
    JSON.SetString(save, "position", JSON.ToString(pos))
    JSON.Free(pos)

    ; Save inventory as array
    inv = JSON.MakeArray()
    JSON.Append(inv, "sword")
    JSON.Append(inv, "shield")
    JSON.SetString(save, "inventory", JSON.ToString(inv))
    JSON.Free(inv)

    JSON.ToFilePretty(save, fileName)
    JSON.Free(save)
    PRINT "Game saved to " + fileName
END FUNCTION

; === Load Game ===
FUNCTION LoadGame(fileName)
    save = JSON.Parse(fileName)
    IF save = 0 THEN
        PRINT "No save file found"
        RETURN
    ENDIF

    playerName = JSON.GetString(save, "name")
    level = JSON.GetInt(save, "level")
    health = JSON.GetFloat(save, "health")

    PRINT "Loaded: " + playerName + " (Level " + STR(level) + ")"
    PRINT "Health: " + STR(health)

    JSON.Free(save)
END FUNCTION

; Usage
SaveGame("save.json", "Hero", 5, 87.5, 10.0, 0.0, -5.0)
LoadGame("save.json")
```

---

## See Also

- [FILE](FILE.md) — Raw file I/O
- [NET](NET.md) — Parse JSON from network responses
