# Stopwatch Commands

High-resolution elapsed-time handles. Useful for measuring performance, timing game events, and profiling code sections.

## Core Workflow

1. `STOPWATCH.NEW()` — start a new stopwatch; it begins counting immediately.
2. `STOPWATCH.ELAPSED(sw)` — read elapsed time in seconds.
3. `STOPWATCH.RESET(sw)` — restart from zero.
4. `STOPWATCH.FREE(sw)` when done.

---

## Creation

### `STOPWATCH.NEW()` 

Creates a new stopwatch handle. The timer starts immediately. Returns a **stopwatch handle**.

---

## Elapsed Time

### `STOPWATCH.ELAPSED(sw)` 

Returns elapsed time in **seconds** since creation or last reset, as a float.

---

## Reset

### `STOPWATCH.RESET(sw)` 

Restarts the stopwatch from zero without freeing it.

---

## Lifetime

### `STOPWATCH.FREE(sw)` 

Frees the stopwatch handle.

---

## Full Example

Measuring frame time and counting a 10-second game event.

```basic
WINDOW.OPEN(800, 450, "Stopwatch Demo")
WINDOW.SETFPS(60)

gameTimer  = STOPWATCH.NEW()
frameTimer = STOPWATCH.NEW()

WHILE NOT WINDOW.SHOULDCLOSE()
    ft = STOPWATCH.ELAPSED(frameTimer)
    STOPWATCH.RESET(frameTimer)

    elapsed = STOPWATCH.ELAPSED(gameTimer)
    remaining = MAX(0, 10.0 - elapsed)

    RENDER.CLEAR(20, 20, 40)
    DRAW.TEXT("Frame ms: " + STR(INT(ft * 1000)), 10, 10, 18, 200, 200, 200, 255)
    DRAW.TEXT("Time left: " + STR(INT(remaining)), 10, 40, 24, 255, 200, 60, 255)

    IF remaining <= 0 THEN
        DRAW.TEXT("DONE!", 340, 200, 32, 80, 255, 80, 255)
    END IF

    RENDER.FRAME()
WEND

STOPWATCH.FREE(gameTimer)
STOPWATCH.FREE(frameTimer)
WINDOW.CLOSE()
```

---

## See also

- [TIMER.md](TIMER.md) — countdown timers
- [TIME.md](TIME.md) — `TIME.DELTA`, `TIME.NOW`, `TICKCOUNT`
- [DEBUG.md](DEBUG.md) — debug assertions and overlays
