# Debug Commands

Commands to help with debugging and profiling your application.

---

### `DEBUG.ENABLE` / `DEBUG.DISABLE`

With **CGO** and Raylib, the watch overlay (`DEBUG.WATCH`) is drawn at frame end only when **`Registry.DebugMode`** is true (for example the host passes **`--debug`** / pipeline **`Options.Debug`**) **or** you have called **`DEBUG.ENABLE`**. **`DEBUG.DISABLE`** turns off that user override; if the process is not in debug mode and you have disabled the override, the overlay is not drawn even when watches are stored.

### `DEBUG.ISENABLED`

Returns **`TRUE`** when the overlay is allowed to draw: **`DEBUG.ENABLE`** was used, or **`Registry.DebugMode`** is on. It does not check whether any **`DEBUG.WATCH`** rows exist.

---

### `ASSERT(condition, message$)`

An assertion is a statement that a condition must be true at a specific point in your program. If the `condition` evaluates to `FALSE`, the program will halt and print the `message$`. This is a powerful tool for catching bugs early.

- `condition`: The expression to check. It should evaluate to `TRUE` if everything is correct.
- `message$`: The error message to display if the assertion fails.

Use assertions to validate assumptions in your code. For example, you can assert that a resource handle is valid after loading it, or that a player's health never drops below zero.

```basic
; Example: Validate a resource handle after loading
player_tex = Texture.Load("player.png")
ASSERT(player_tex <> 0, "Failed to load player texture!")

; Example: Ensure a value is within an expected range
FUNCTION SetHealth(health)
    ASSERT(health >= 0 AND health <= 100, "Health value out of range: " + STR$(health))
    ; ... set health ...
ENDFUNCTION
```

---

### `DUMP(value)`

**[PARTIAL]** Coming soon. Intended to print detailed information about a variable or handle.

---

### `TRACE(value)`

**[PARTIAL]** Coming soon. Intended to enable/disable verbose logging from the runtime.
