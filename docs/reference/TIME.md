# Time (`TIME.*`, `TICKCOUNT`, wall-clock)

Commands for elapsed program time, per-frame delta, and wall-clock values.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`** or global builtins as registered; Easy Mode (`Time.Delta`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

## Core Workflow

- **Delta time:** **`TIME.DELTA()`** — seconds since the last frame; use for movement and logic independent of frame rate.
- **Program time:** **`TIME.GET()`** or **`TICKCOUNT()`** — elapsed time since start (**`TICKCOUNT`** is milliseconds; **`TIME.GET`** is seconds).
- **Wall-clock:** **`DATE`**, **`TIME`**, **`DATETIME`**, **`YEAR`**, etc. — real-world time from the system clock (see manifest for the full set). These are **global** builtins; they are not the same as **`TIME.DELTA`** / **`TIME.GET`** (elapsed program time).

---

### `TIME.DELTA()`
Returns seconds since last frame.

- **Returns**: (Float) The frame delta time (DT).
- **Example**:
    ```basic
    pos = pos + speed * TIME.DELTA()
    ```

---

### `WINDOW.SETLOOPMODE(mode$, hz)`
Controls how **`TIME.DELTA()`** reports frame time (see also **`TIME.SETMAXDELTA`** for a simple cap).

| Mode | Behavior |
|------|----------|
| `"variable"` | Default — raw frame time (with optional max-delta cap). |
| `"fixed"` | Constant step **`1/hz`** (default 60 Hz) — use for fixed-timestep physics. |
| `"semi-fixed"` | Cap delta to **`hz`** seconds (default 0.1 s) — avoids huge steps on alt-tab. |

```basic
WINDOW.SETLOOPMODE("fixed", 60)
WINDOW.SETLOOPMODE("semi-fixed", 0.1)
```

Use **`TIME.PHYSICSSTEPS()`** and **`TIME.PHYSICSSTEP()`** with fixed mode to run multiple physics updates per frame when the accumulator has credit:

```basic
WINDOW.SETLOOPMODE("fixed", 60)
WHILE NOT WINDOW.CLOSED()
    steps = TIME.PHYSICSSTEPS()
    FOR i = 1 TO steps
        UpdatePhysics(TIME.PHYSICSSTEP())
    NEXT
    Draw()
WEND
```

---

### `TIME.GET()`
Returns total elapsed seconds since start.

- **Returns**: (Float) Total seconds.

---

### `TICKCOUNT()`
Returns total elapsed milliseconds since start.

- **Returns**: (Integer) Total milliseconds.

---

### `DATE` / `TIME` / `DATETIME` / `TIMESTAMP`
Returns wall-clock strings or epoch values.

- **Returns**: (String/Integer) Formatted date/time or Unix timestamp.

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

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `TIME.GETFPS()` | Returns current measured frames per second. |
| `TIME.SETMAXDELTA(max)` | Cap the maximum delta time returned by `TIME.DELTA` (prevents large spikes). |
| `TIME.UPDATE()` | Manually advance the internal clock (needed outside `RENDER.FRAME` loops). |

## See also

- [WINDOW.md](WINDOW.md) — `WINDOW.SETTARGETFPS`
- [TWEEN.md](TWEEN.md) — `TWEEN.UPDATE(dt)`
