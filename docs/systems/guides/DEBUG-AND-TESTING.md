# Debug, testing, and validation

> Log problems, watch live values, validate scripts before run, and run smoke tests.

**Namespaces:** `DEBUG` · CLI · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#debug-timer](../COMMAND_REGISTRY.md#debug-timer) · [10-DEBUG-TIMER.md](../10-DEBUG-TIMER.md) · [11-TOOLING.md](../11-TOOLING.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [Validate before run](#validate-before-run)
- [Runtime logging](#runtime-logging)
- [On-screen watches](#on-screen-watches)
- [Visual debug draws](#visual-debug-draws)
- [Project tests](#project-tests)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Tool | When | Why |
|------|------|-----|
| `moonbasic --check` | Before every run / in CI | Catch typos and arity without GPU |
| `DEBUG.LOG` | Development | Trace flow in console |
| `DEBUG.WATCH` | Tuning gameplay | Live HUD values |
| `DEBUG.DRAWLINE` | 3D debug | See rays, paths |
| `moonbasic test` | Maintainer smoke | Language/runtime regression |

---

## Validate before run

```bash
moonbasic --check main.mb
```

**Why:** Compiler reports unknown commands (`ENTITY.SETPOSITON`) with **did you mean** hints — faster than a black screen.

Use in VS Code via LSP diagnostics (install VSIX from Releases).

---

## Runtime logging

```basic
DEBUG.LOG("Spawn wave " + wave)
DEBUG.WARN("Low FPS: " + APP.GETFPS())
```

**Why `LOG` vs `PRINT`:** `DEBUG.*` can integrate with host debug mode and log files (`DEBUG.LOGFILE`).

---

## On-screen watches

```basic
DEBUG.ENABLE()
WHILE NOT APP.SHOULDCLOSE()
    DEBUG.WATCH("fps", APP.GETFPS())
    DEBUG.WATCH("hp", hp)
    ; ... game ...
WEND
```

**Why:** Overlay at frame end when debug enabled — no manual `DRAW.TEXT` for every variable.

Clear when done: `DEBUG.WATCHCLEAR()`.

---

## Visual debug draws

Inside 3D pass:

```basic
DEBUG.DRAWLINE(0, 0, 0, 5, 0, 0)
DEBUG.DRAWBOX(targetEntity)
```

**Why:** See collision volumes and AI paths in the world.

---

## Project tests

```bash
moonbasic test
```

Runs compiler/runtime smoke checks when developing from a full source tree. For your game:

```bash
moonbasic --check main.mb
moonbasic --check levels/level1.mb
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Only `PRINT` in loop | Floods console — use `WATCH` |
| Skip `--check` in CI | Add to pipeline |
| Debug draw after `RENDER.FRAME` | Draw before frame present |

---

## See also

- [ERROR_MESSAGES.md](../../ERROR_MESSAGES.md)
- [10-DEBUG-TIMER.md](../10-DEBUG-TIMER.md)
