# Time (`TIME.*`, `TICKCOUNT`, wall-clock)

Commands for elapsed program time, per-frame delta, and wall-clock values.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`** or global builtins as registered; Easy Mode (`Time.Delta`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

## Core concepts

- **Delta time:** **`TIME.DELTA()`** — seconds since the last frame; use for movement and logic independent of frame rate.
- **Program time:** **`TIME.GET()`** or **`TICKCOUNT()`** — elapsed time since start (**`TICKCOUNT`** is milliseconds; **`TIME.GET`** is seconds).
- **Wall-clock:** **`DATE`**, **`TIME`**, **`DATETIME`**, **`YEAR`**, etc. — real-world time from the system clock (see manifest for the full set). These are **global** builtins; they are not the same as **`TIME.DELTA`** / **`TIME.GET`** (elapsed program time).

---

### `TIME.DELTA()`
Returns seconds since last frame.

### `TIME.GET()`
Returns total elapsed seconds since start.

### `TICKCOUNT()`
Returns total elapsed milliseconds since start.

---

### `DATE` / `TIME` / `DATETIME` / `TIMESTAMP`
Returns formatted wall-clock strings or epoch values as documented in the command registry (**`TIMESTAMP`** — Unix seconds).

---

## Full Example: a simple stopwatch

```basic
WINDOW.OPEN(800, 600, "Stopwatch Example")
WINDOW.SETFPS(60)

start_time = 0.0
stop_time = 0.0
running = FALSE

WHILE NOT WINDOW.SHOULDCLOSE()
    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        IF running THEN
            stop_time = TIME.GET()
            running = FALSE
        ELSE
            start_time = TIME.GET()
            stop_time = 0.0
            running = TRUE
        ENDIF
    ENDIF

    elapsed_time = 0.0
    IF running THEN
        elapsed_time = TIME.GET() - start_time
    ELSE
        elapsed_time = stop_time - start_time
    ENDIF

    RENDER.CLEAR(0, 0, 0)
    CAMERA2D.BEGIN()
        DRAW.TEXT("Press SPACE to start/stop", 210, 150, 20, 150, 150, 150, 255)
        DRAW.TEXT(FORMAT(elapsed_time, "%.2f"), 300, 250, 60, 100, 200, 255, 255)
    CAMERA2D.END()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```
