# Debug Commands

Runtime assertions, debug overlays, and diagnostic logging.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Use `ASSERT` to guard invariants. Enable the debug overlay with `DEBUG.ENABLE`, add watches with `DEBUG.WATCH`, and toggle visibility with `DEBUG.ENABLE` / `DEBUG.DISABLE`.

---

### `DEBUG.ENABLE` / `DEBUG.DISABLE` 

With **CGO** and Raylib, the watch overlay (`DEBUG.WATCH`) is drawn at frame end only when **`Registry.DebugMode`** is true (for example the host passes **`--debug`** / pipeline **`Options.Debug`**) **or** you have called **`DEBUG.ENABLE`**. **`DEBUG.DISABLE`** turns off that user override; if the process is not in debug mode and you have disabled the override, the overlay is not drawn even when watches are stored.

---

### `DEBUG.ISENABLED` 

Returns **`TRUE`** when the overlay is allowed to draw: **`DEBUG.ENABLE`** was used, or **`Registry.DebugMode`** is on. It does not check whether any **`DEBUG.WATCH`** rows exist.

---

### `ASSERT(condition, message)` / `DEBUG.ASSERT(condition, message)` 

If **`condition`** is **`FALSE`**, the program halts and prints **`message`**. **`ASSERT`** and **`DEBUG.ASSERT`** invoke the same implementation.

- `condition`: Should be **`TRUE`** when the program state is valid.
- `message`: Error text when the assertion fails.

```basic
player_tex = TEXTURE.LOAD("player.png")
ASSERT(player_tex <> 0, "Failed to load player texture!")

FUNCTION SetHealth(health)
    ASSERT(health >= 0 AND health <= 100, "Health value out of range: " + STR(health))
ENDFUNCTION
```

---

### `DUMP(value)` 

**[PARTIAL]** Coming soon. Intended to print detailed information about a variable or handle.

---

### `TRACE(value)` 

**[PARTIAL]** Coming soon. Intended to enable/disable verbose logging from the runtime.

---

## Full Example

Assertions and a debug watch overlay showing player health and position.

```basic
WINDOW.OPEN(800, 450, "Debug Demo")
WINDOW.SETFPS(60)

DEBUG.ENABLE()
health = 100
player = ENTITY.CREATECUBE(1.0)
ENTITY.SETPOS(player, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    ASSERT(health >= 0, "Health went negative!")

    DEBUG.WATCH("health", health)
    px, py, pz = ENTITY.GETPOS(player)
    DEBUG.WATCH("pos", STR(px) + " " + STR(py) + " " + STR(pz))

    ENTITY.UPDATE(dt)
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN3D(CAMERA.CREATE())
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [SYSTEM.md](SYSTEM.md) — `PRINT`, `PRINTERROR`, runtime info
- [CONSOLE.md](CONSOLE.md) — in-game console overlay
