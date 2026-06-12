# Debug, error, and timer systems

> Logging, on-screen debug draws, compiler diagnostics, and delayed callbacks.

**All commands:** [COMMAND_REGISTRY.md#debug-timer](COMMAND_REGISTRY.md#debug-timer)

**Deep guides:** [guides/DEBUG-AND-TESTING.md](guides/DEBUG-AND-TESTING.md) · [guides/COMPILER-ERRORS.md](guides/COMPILER-ERRORS.md)

**See also:** [reference/DEBUG.md](../reference/DEBUG.md) · [reference/TIMER.md](../reference/TIMER.md) · [ERROR_MESSAGES.md](../ERROR_MESSAGES.md)

---

## Table of contents

- [DEBUG system](#debug-system)
- [ERROR system](#error-system)
- [TIMER system](#timer-system)
- [Full example](#full-example)
- [See also](#see-also)

---

## DEBUG system

Runtime logging and visual debug overlays.

### Core workflow

1. `DEBUG.LOG(msg)` during development.
2. `DEBUG.WATCH(label, value)` each frame for live HUD values (enable with `DEBUG.ENABLE` or host debug mode).
3. `DEBUG.DRAWLINE` / `DRAWBOX` inside the 3D pass for traces.
4. `DEBUG.SHOWFPSGRAPH(true)` for performance graphs.

---

### Logging

| Command | Description |
|---------|-------------|
| `DEBUG.LOG(text)` | Print line to debug output |
| `DEBUG.WARN(text)` | Warning line |
| `DEBUG.LOGFILE(path, text)` | Append to file |

**Example:**

```basic
DEBUG.LOG("Player spawned at " + px)
```

---

### Visual debug

| Command | Description |
|---------|-------------|
| `DEBUG.DRAWLINE(x1,y1,z1,x2,y2,z2)` | World line |
| `DEBUG.DRAWBOX(entity)` | Box around entity |
| `DEBUG.WATCH(label, value)` | Register HUD watch |
| `DEBUG.WATCHCLEAR()` | Clear watches |
| `DEBUG.SHOWFPSGRAPH(enabled)` | FPS graph overlay |
| `DEBUG.ENABLE` / `DEBUG.DISABLE` | Toggle user debug overlay |

On-screen overlay requires **`moonrun`** (full runtime).

**Example:**

```basic
DEBUG.WATCH("fps", APP.GETFPS())
DEBUG.DRAWLINE(0, 0, 0, 5, 0, 0)
```

---

## ERROR system

**Status:** Shipped at compile time — the moonBASIC compiler reports file, line, and helpful messages before run time.

### What you get today

| Feature | How |
|---------|-----|
| File and line | Errors cite `main.mb:42:5` |
| Unknown commands | Clear error + **did you mean** suggestions |
| Bad arity / types | Semantic phase diagnostics |
| Strict deprecations | `moonbasic --check --strict-deprecated` |

**Example compiler output:**

```text
main.mb:42:5
ENTITY.SETPOSITON(player, 0, 1, 5)
      ^^^^^^^^^^
Unknown command: ENTITY.SETPOSITON
Did you mean: ENTITY.SETPOSITION
```

Validate scripts without running:

```bash
moonbasic --check main.mb
```

Runtime crashes surface through `DEBUG` / stderr; full call-stack overlays are evolving — see [ERROR_MESSAGES.md](../ERROR_MESSAGES.md).

---

## TIMER system

Handle timers and frame-scheduled callbacks.

### Handle timers

| Command | Description |
|---------|-------------|
| `TIMER.CREATE(duration)` | Countdown handle |
| `TIMER.START(t)` / `STOP(t)` | Control timer |
| `TIMER.UPDATE(t)` | Tick (or auto from frame hook) |
| `TIMER.DONE(t)` / `FINISHED(t)` | Elapsed? |
| `TIMER.FREE(t)` | Release handle |

**Example:**

```basic
t = TIMER.CREATE(2.0)
TIMER.START(t)
IF TIMER.DONE(t) THEN respawn()
```

---

### Callback timers

| Command | Description |
|---------|-------------|
| `TIMER.AFTER(seconds, functionName)` | Run once after delay |
| `TIMER.EVERY(seconds, functionName)` | Repeat every interval |
| `TIMER.CANCEL(id)` | Cancel scheduled callback |

Callbacks invoke **user functions** by name — define `FUNCTION SpawnEnemy()` etc.

**Example:**

```basic
TIMER.EVERY(1.0, "SpawnWave")
TIMER.AFTER(5.0, "EndRound")

FUNCTION SpawnWave()
    ; spawn logic
ENDFUNCTION
```

`TIMER.AFTER` removes itself after fire; use **`TIMER.CANCEL`** for repeating timers.

---

## Full example

```basic
APP.OPEN(640, 480, "Debug + Timer")
APP.SETFPS(60)
DEBUG.ENABLE()

waves = 0
TIMER.EVERY(2.0, "OnWave")

FUNCTION OnWave()
    waves = waves + 1
    DEBUG.LOG("Wave " + waves)
ENDFUNCTION

WHILE NOT APP.SHOULDCLOSE()
    DEBUG.WATCH("waves", waves)
    DEBUG.WATCH("fps", APP.GETFPS())
    RENDER.CLEAR(0, 0, 0)
    DRAW.TEXT("Waves: " + waves, 10, 10, 18, 255, 255, 255)
    RENDER.FRAME()
WEND

APP.CLOSE()
```

---

## See also

- [11-TOOLING](11-TOOLING.md) — `moonbasic test`
- [MEMORY.md](../MEMORY.md) — `DEBUG.HEAPSTATS`
