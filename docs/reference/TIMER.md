# Timer Commands

Wall-clock and simulation timers for countdowns, cooldowns, and timed events.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**Wall-clock timers** — Create with `TIMER.NEW`, check `TIMER.FINISHED`, free with `TIMER.FREE`.

**Simulation timers** — Create with `TIMER.CREATE`, start with `TIMER.START`, advance with `TIMER.UPDATE(timer, dt)` each frame, check `TIMER.DONE` for edge-triggered completion.

`TIMER.REMAINING` works for both types.

---

### `TIMER.NEW(duration)` 

Creates a wall-clock timer that expires `duration` seconds from now. Returns a handle.

---

### `TIMER.RESET(timerHandle, duration)` 

Reschedules a wall-clock timer's deadline.

---

### `TIMER.FINISHED(timerHandle)` 

Returns `TRUE` after a wall-clock timer's deadline has passed.

---

### `TIMER.CREATE(duration)` 

Creates a stopped simulation timer with the given duration in seconds. Returns a handle.

---

### `TIMER.MAKE(duration)` 

Deprecated alias for `TIMER.CREATE`.

---

### `TIMER.START(timerHandle)` 

Starts or restarts a simulation timer from zero.

---

### `TIMER.STOP(timerHandle)` 

Pauses a simulation timer.

---

### `TIMER.REWIND(timerHandle)` 

Resets elapsed time to zero.

---

### `TIMER.UPDATE(timerHandle, dt)` 

Advances a simulation timer by `dt` seconds (non-negative).

---

### `TIMER.DONE(timerHandle)` 

Returns `TRUE` for one call when a simulation timer cycle completes (edge-triggered).

---

### `TIMER.FRACTION(timerHandle)` 

Returns `elapsed / duration`, clamped to 0.0–1.0.

---

### `TIMER.SETLOOP(timerHandle, loop)` 

If `TRUE`, the timer wraps around and `TIMER.DONE` pulses each cycle.

---

### `TIMER.REMAINING(timerHandle)` 

Returns seconds remaining. Works for both wall-clock and simulation timers.

---

### `TIMER.FREE(timerHandle)` 

Frees the timer handle.

---

## Full Example

This example uses a simulation timer for a 3-second cooldown.

```basic
cd = TIMER.CREATE(3.0)
TIMER.START(cd)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = DELTATIME()
    TIMER.UPDATE(cd, dt)

    IF TIMER.DONE(cd)
        PRINT "Cooldown ready!"
        TIMER.REWIND(cd)
        TIMER.START(cd)
    END IF

    pct = TIMER.FRACTION(cd)
    RENDER.BEGINFRAME()
    DRAW.RECT(50, 50, INT(200 * pct), 20, 0, 255, 0, 255)
    RENDER.ENDFRAME()
WEND

TIMER.FREE(cd)
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `TIMER.GETLOOP(timer)` | Returns the loop count setting (-1 = infinite). |

## See also

- [TIME.md](TIME.md) — `TIME.DELTA`, `TIME.GETFPS`
- [TWEEN.md](TWEEN.md) — value animation with callbacks
