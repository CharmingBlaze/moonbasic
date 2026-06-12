# Save and progress — persist player state

> Store level, score, and inventory in memory, then write a **JSON save file** players can reload.

**Namespaces:** `SAVE` · `FILE` · `JSON` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../COMMAND_REGISTRY.md#data) · [09-DATA.md](../09-DATA.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to save](#when-to-save)
- [Core workflow](#core-workflow)
- [Key commands](#key-commands)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Layer | API | Why |
|-------|-----|-----|
| **Runtime table** | `SAVE.SET` / `SAVE.GET` | Fast key/value while playing |
| **Disk file** | `SAVE.WRITE` / `SAVE.READ` | Persist across sessions |
| **Raw files** | `FILE.READTEXT` / `WRITETEXT` | Config, mods, custom formats |
| **Structured** | `JSON.PARSE` / `STRINGIFY` | Complex nested data |

**Case:** Save **keys** are **case-insensitive** in scripts.

---

## When to save

| Trigger | Pattern |
|---------|---------|
| Level complete | `SAVE.WRITE` once |
| Autosave timer | `TIMER.EVERY` → write |
| On exit | `SAVE.WRITE` before `APP.CLOSE` |
| Death | Usually **don’t** write — reload last save |

**Why separate SET and WRITE:** `SET` is cheap RAM; `WRITE` is disk IO — batch before writing.

---

## Core workflow

1. **Load** (if file exists): `SAVE.READ("save1.json")`.
2. **Play** — update `SAVE.SET("level", n)` when values change.
3. **Save** — `SAVE.WRITE("save1.json")` at milestones.
4. **Ship** — include save path in player folder; document filename.

```basic
IF FILE.EXISTS("save1.json") THEN SAVE.READ("save1.json")
level = SAVE.GET("level")
IF level = 0 THEN level = 1
```

---

## Key commands

| Command | Why |
|---------|-----|
| `SAVE.SET(key, value)` / `SAVE.DATA` | Store in memory |
| `SAVE.GET(key)` | Read back |
| `SAVE.WRITE(path)` | Flush table to JSON file |
| `SAVE.READ(path)` | Load file into table |
| `SAVE.WRITEFILE` / `READFILE` | Raw blob helpers |

For custom JSON not using SAVE table:

```basic
text = FILE.READTEXT("config.json")
doc = JSON.PARSE(text)
hp = JSON.GET(doc, "player.health")
```

---

## Full example

```basic
APP.OPEN(400, 300, "Save demo")
score = 0
IF FILE.EXISTS("hiscore.json") THEN
    SAVE.READ("hiscore.json")
    score = SAVE.GET("best")
ENDIF

WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYHIT(KEY_SPACE) THEN score = score + 1
    IF score > SAVE.GET("best") THEN
        SAVE.SET("best", score)
        SAVE.WRITE("hiscore.json")
    ENDIF
    DRAW.TEXT("Score " + score, 10, 10, 16, 255, 255, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
```

Example with quit save: [examples/rpg](../examples/rpg/main.mb).

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| `GET` before `READ` | Load file first or default values |
| Save every frame | Write on events only |
| Absolute paths | Use relative paths next to `moonrun` cwd |
| Secrets in save | Players can edit JSON — don’t trust client |

---

## See also

- [09-DATA.md](../09-DATA.md)
- [11-TOOLING.md](../11-TOOLING.md) — ship save folder with game
