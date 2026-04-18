# Save Commands

Persist and retrieve JSON data to/from a local save slot. Minimal key/value store backed by a JSON file on disk.

For full SQLite database persistence see [DATABASE.md](DATABASE.md).

## Core Workflow

1. `SAVE.DATA(key, jsonString)` — write a JSON string to `key`.
2. `SAVE.GET(key)` — read it back as a string.
3. Parse/encode with `JSON.*` commands.

---

## Commands

### `SAVE.DATA(key, jsonString)` 

Writes `jsonString` to persistent storage under `key`. Overwrites existing data for the same key.

---

### `SAVE.GET(key)` 

Reads the JSON string stored under `key`. Returns an empty string if the key does not exist.

---

## Full Example

Saving and loading player progress.

```basic
WINDOW.OPEN(800, 450, "Save Demo")
WINDOW.SETFPS(60)

score = 0
; load saved score
saved = SAVE.GET("playerData")
IF saved <> ""
    obj   = JSON.PARSE(saved)
    score = JSON.GETINT(obj, "score")
    PRINT "Loaded score: " + STR(score)
END IF

WHILE NOT WINDOW.SHOULDCLOSE()
    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        score = score + 10
        ; save immediately
        data = JSON.MAKE()
        JSON.SETINT(data, "score", score)
        SAVE.DATA("playerData", JSON.STRINGIFY(data))
        JSON.FREE(data)
    END IF

    RENDER.CLEAR(20, 20, 40)
    DRAW.TEXT("Score: " + STR(score), 10, 10, 24, 255, 255, 255, 255)
    DRAW.TEXT("SPACE = +10 and save", 10, 45, 18, 180, 180, 180, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [DATABASE.md](DATABASE.md) — full SQLite database
- [JSON.md](JSON.md) — JSON encode/decode
- [FILE.md](FILE.md) — raw file I/O
- [CONFIG.md](CONFIG.md) — INI-style config files
