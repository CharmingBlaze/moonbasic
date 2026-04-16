# Debug (`DEBUG.*`, `ASSERT`, `DUMP`, `TRACE`)

Commands to help with debugging and profiling your application.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`** where applicable; Easy Mode names are [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings).

---

### `DEBUG.ENABLE` / `DEBUG.DISABLE`

With **CGO** and Raylib, the watch overlay (`DEBUG.WATCH`) is drawn at frame end only when **`Registry.DebugMode`** is true (for example the host passes **`--debug`** / pipeline **`Options.Debug`**) **or** you have called **`DEBUG.ENABLE`**. **`DEBUG.DISABLE`** turns off that user override; if the process is not in debug mode and you have disabled the override, the overlay is not drawn even when watches are stored.

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
