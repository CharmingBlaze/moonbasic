# Data systems: SAVE, FILE, JSON, MATH, VEC3

> Persistence, file I/O, JSON configs, game math, and 3D vectors.

**All commands:** [COMMAND_REGISTRY.md#data](COMMAND_REGISTRY.md#data)

**See also:** [reference/SAVE.md](../reference/SAVE.md) · [reference/JSON.md](../reference/JSON.md) · [MEMORY.md](../MEMORY.md)

**Case:** **SAVE** keys and JSON path segments are **case-insensitive** in scripts.

**Deep guides:** [guides/SAVE-AND-PROGRESS.md](guides/SAVE-AND-PROGRESS.md) · [guides/FILES-AND-JSON.md](guides/FILES-AND-JSON.md) · [guides/MATH-AND-VECTORS.md](guides/MATH-AND-VECTORS.md) · [guides/math/README.md](guides/math/README.md) (2D/3D math & vectors)

---

## Table of contents

- [SAVE system](#save-system)
- [FILE system](#file-system)
- [JSON system](#json-system)
- [MATH system](#math-system)
- [VEC3 system](#vec3-system)
- [Full example](#full-example)
- [See also](#see-also)

---

## SAVE system

Key/value game state with JSON file export.

### Core workflow

1. `SAVE.SET(key, value)` or `SAVE.DATA(key, value)` — store in memory.
2. `SAVE.GET(key)` — read back (typed coercion).
3. `SAVE.WRITE(path)` / `SAVE.READ(path)` — persist to disk.
4. `SAVE.WRITEFILE` / `SAVE.READFILE` — raw file helpers.

---

### `SAVE.SET(key, value)` / `SAVE.DATA(key, value)`

Stores a value under a string key.

| Argument | Type | Description |
|----------|------|-------------|
| key | string | Save slot key |
| value | any | Number, string, bool, … |

**Returns:** nothing

**Example:**

```basic
SAVE.SET("level", 3)
SAVE.DATA("health", 80)
```

---

### `SAVE.GET(key)`

Reads a stored value.

**Returns:** stored value or default behavior per type — see [reference/SAVE.md](../reference/SAVE.md)

**Example:**

```basic
level = SAVE.GET("level")
```

---

### `SAVE.WRITE(path)` / `SAVE.READ(path)`

Writes or loads the in-memory save table as JSON.

**Example:**

```basic
SAVE.WRITE("save1.json")
SAVE.READ("save1.json")
hp = SAVE.GET("health")
```

---

## FILE system

Low-level text file helpers.

| Command | Description |
|---------|-------------|
| `FILE.EXISTS(path)` | `bool` — file present |
| `FILE.READTEXT(path)` | `string` — entire file |
| `FILE.WRITETEXT(path, text)` | Write string to file |
| `FILE.DELETE(path)` | Remove file |

**Aliases:** checklist `READTEXT` / `WRITETEXT` map to these names.

**Example:**

```basic
IF FILE.EXISTS("config.json") THEN
    cfg = FILE.READTEXT("config.json")
ENDIF
```

Paths follow OS rules; use forward slashes in scripts when possible.

---

## JSON system

Parse configs, saves, and network payloads.

### `JSON.PARSE(text)`

Parses a JSON string into a handle/object.

**Returns:** `handle` or JSON value per runtime

**Example:**

```basic
text = FILE.READTEXT("config.json")
doc = JSON.PARSE(text)
```

---

### `JSON.GET(json, path)` / `JSON.SET(json, path, value)`

Path-based access (e.g. `"player.health"`).

**Example:**

```basic
hp = JSON.GET(doc, "player.health")
JSON.SET(doc, "player.health", 100)
```

---

### `JSON.STRINGIFY(json)` / `JSON.TOSTRING(json)`

Serialize back to string.

**Aliases:** `JSON.STRINGIFY` ↔ stringify helpers

**Example:**

```basic
out = JSON.STRINGIFY(doc)
FILE.WRITETEXT("out.json", out)
```

---

## MATH system

Common game math helpers.

| Command | Description |
|---------|-------------|
| `MATH.RAND(min, max)` | Random int in range |
| `MATH.RANDF(min, max)` | Random float |
| `MATH.CLAMP(v, lo, hi)` | Clamp value |
| `MATH.LERP(a, b, t)` | Linear interpolation |
| `MATH.DISTANCE(x1,y1,z1,x2,y2,z2)` | 3D distance |
| `VEC3.DISTANCE(a, b)` | Distance between vec handles |

**Example:**

```basic
roll = MATH.RAND(1, 6)
t = MATH.CLAMP(hp / 100, 0, 1)
x = MATH.LERP(0, 100, t)
```

See [reference/MATH.md](../reference/MATH.md) for trig, smoothstep, remaps.

---

## VEC3 system

3D vector handles for cleaner math.

| Command | Description |
|---------|-------------|
| `VEC3.CREATE(x, y, z)` | New vector |
| `VEC3.ADD(a, b)` | Add vectors |
| `VEC3.NORMALIZE(v)` | Unit vector |
| `VEC3.LENGTH(v)` | Magnitude |
| `VEC3.DOT` / `CROSS` | Products |

**Example:**

```basic
a = VEC3.CREATE(0, 1, 0)
b = VEC3.CREATE(1, 0, 0)
c = VEC3.ADD(a, b)
```

---

## Full example

```basic
; Save / load high score
APP.OPEN(400, 300, "Save demo")
APP.SETFPS(60)

IF FILE.EXISTS("hiscore.json") THEN
    SAVE.READ("hiscore.json")
ENDIF
best = SAVE.GET("best")
IF best = 0 THEN best = 0

score = 0
WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYHIT(KEY_SPACE) THEN score = score + MATH.RAND(1, 5)
    IF score > best THEN
        best = score
        SAVE.SET("best", best)
        SAVE.WRITE("hiscore.json")
    ENDIF
    DRAW.TEXT("Score " + score + "  Best " + best, 10, 10, 16, 255, 255, 255)
    RENDER.FRAME()
WEND

APP.CLOSE()
```

Check with **`moonbasic --check`**; run with **`moonrun`**.

---

## See also

- [examples/rpg](../examples/rpg/main.mb) — JSON save on exit
- [11-TOOLING](11-TOOLING.md) — ship save paths with your game folder
