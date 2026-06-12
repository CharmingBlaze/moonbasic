# Files and JSON — configs, mods, and raw data

> Read and write text files, parse JSON configs, and serialize data — separate from the in-game `SAVE` table.

**Namespaces:** `FILE` · `JSON` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../COMMAND_REGISTRY.md#data) · **Overview:** [09-DATA.md](../09-DATA.md)

**Related:** Player progress slots → [SAVE-AND-PROGRESS.md](SAVE-AND-PROGRESS.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow — config file](#core-workflow--config-file)
- [FILE commands](#file-commands)
- [JSON commands](#json-commands)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Exists/read/write/delete text files; parse and stringify JSON |
| **You need first** | Paths relative to game folder (ship with `moonbasic package`) |
| **Typical games** | Options menu, level manifests, mod-friendly data |
| **Not for** | Simple high-score only — `SAVE.WRITE` may be enough |

**Why FILE + JSON:** `SAVE.*` is a **key/value game state** API. `FILE` + `JSON` handle **arbitrary documents** (nested objects, arrays, configs).

---

## When to use this system

**Use when:**

- Loading `config.json` for difficulty, key bindings export.
- Reading level list or quest tables from data files.
- Writing generated content (replay, editor export).
- Parsing network or tool output (with validation).

**Skip when:**

- Only `level` and `health` keys — [SAVE-AND-PROGRESS.md](SAVE-AND-PROGRESS.md).
- Binary assets — `TEXTURE.LOAD` / `ASSET.*`.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| High score / flags | `SAVE.SET` + `SAVE.WRITE` | Manual JSON for one int |
| Nested config | `FILE.READTEXT` + `JSON.PARSE` | Many `SAVE` keys mimicking tree |
| Raw log append | `FILE.WRITETEXT` or `DEBUG.LOGFILE` | `JSON` for plain logs |
| Ship default config | `FILE.EXISTS` → else write defaults | Crash if missing |
| Mod path | `FILE.READTEXT` from user folder | Hard-coded only in code |

---

## Core workflow — config file

1. **Check** — `IF FILE.EXISTS("config.json")`.  
   **Why:** First run may have no file yet.

2. **Read** — `text = FILE.READTEXT("config.json")`.  
   **Why:** Whole file as string.

3. **Parse** — `doc = JSON.PARSE(text)`.  
   **Why:** Structured access.

4. **Use** — `JSON.GET(doc, "audio.volume")`.  
   **Why:** Path syntax for nested keys.

5. **Save changes** — `JSON.SET` → `JSON.STRINGIFY` → `FILE.WRITETEXT`.

Paths in scripts: prefer forward slashes; OS resolves on disk.

---

## FILE commands

| Command | Why |
|---------|-----|
| `FILE.EXISTS(path)` | Avoid read errors on first run |
| `FILE.READTEXT(path)` | Load entire text file |
| `FILE.WRITETEXT(path, text)` | Save string (overwrite) |
| `FILE.DELETE(path)` | Remove save slot or temp file |

**Aliases:** Checklist `READTEXT` / `WRITETEXT` map here.

---

## JSON commands

| Command | Why |
|---------|-----|
| `JSON.PARSE(text)` | String → document handle |
| `JSON.GET(doc, path)` | Read `"player.health"` style paths |
| `JSON.SET(doc, path, value)` | Write nested field |
| `JSON.STRINGIFY(doc)` / `TOSTRING` | Document → string for `WRITETEXT` |

**Case:** Path segments are **case-insensitive** in scripts ([LANGUAGE.md](../../LANGUAGE.md)).

---

## Full example

**Runnable:** [examples/guides/files_json.mb](../../../examples/guides/files_json.mb)

```basic
; Check: moonbasic --check examples/guides/files_json.mb
; Run:   moonrun examples/guides/files_json.mb

volume = 0.8
IF FILE.EXISTS("options.json") THEN
    rawJson = FILE.READTEXT("options.json")
    doc = JSON.PARSESTRING(rawJson)
    volume = JSON.GETFLOAT(doc, "audio.volume")
    JSON.FREE(doc)
ENDIF

APP.OPEN(400, 200, "Files + JSON")
APP.SETFPS(60)

WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYHIT(KEY_UP) THEN volume = volume + 0.1
    IF INPUT.KEYHIT(KEY_DOWN) THEN volume = volume - 0.1
    volume = MATH.CLAMP(volume, 0, 1)

    IF INPUT.KEYHIT(KEY_S) THEN
        doc = JSON.PARSESTRING("{\"audio\":{\"volume\":0}}")
        JSON.SETFLOAT(doc, "audio.volume", volume)
        FILE.WRITETEXT("options.json", JSON.STRINGIFY(doc))
        JSON.FREE(doc)
    ENDIF

    DRAW.TEXT("Volume " + volume + "  (S save)", 10, 10, 16, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Crash if no file | `FILE.EXISTS` before `READTEXT` |
| Invalid JSON | Validate after parse; ship default `{}` |
| Wrong working directory | Run from project root; package copies data folder |
| Huge file every frame | Read once at boot; write on save action |
| Mixing SAVE and JSON paths | Use SAVE for game state; FILE for raw docs |

---

## See also

- [SAVE-AND-PROGRESS.md](SAVE-AND-PROGRESS.md) — `SAVE.WRITE` / `READ`
- [MATH-AND-VECTORS.md](MATH-AND-VECTORS.md) — clamp volume with `MATH.CLAMP`
- [PROJECT-WORKFLOW.md](PROJECT-WORKFLOW.md) — ship JSON with `package`
- [examples/rpg](../../../examples/rpg/main.mb) — save on exit
