# Game helpers (`runtime/mbgame`)

Fast **screen / time / input** shortcuts and **pure** math helpers for gameplay code. The module is registered in the normal pipeline (`mbgame.NewModule()`).

**CGO:** Raylib-backed shortcuts (`SCREENW`, mouse, `DT`, etc.) need **`CGO_ENABLED=1`**. Without CGO, those calls return an error string that mentions Raylib; **`ENDGAME`** / **`GAME.ENDGAME`** still terminate the VM if wired.

---

## Two spellings: bare names and `Game.*`

Many helpers register **twice** so you can write either style:

| Bare (function-like) | Dotted (namespace) |
|----------------------|--------------------|
| `SCREENW()` | `Game.ScreenW()` → registry key **`GAME.SCREENW`** |
| `DT()` | `Game.DT()` → **`GAME.DT`** |
| `KEYDOWN(KEY_ESCAPE)` | `Game.KeyDown(KEY_ESCAPE)` → **`GAME.KEYDOWN`** |

Handlers are identical; pick one style per project for consistency.

---

## Screen and frame timing

| Key | Returns | Notes |
|-----|---------|--------|
| **`SCREENW`** / **`GAME.SCREENW`** | int | `GetScreenWidth` |
| **`SCREENH`** / **`GAME.SCREENH`** | int | `GetScreenHeight` |
| **`SCREENCX`** / **`GAME.SCREENCX`** | float | Horizontal centre |
| **`SCREENCY`** / **`GAME.SCREENCY`** | float | Vertical centre |
| **`DT`** / **`GAME.DT`** | float | Frame delta (seconds); **`0`** when `GamePaused` is set on the runtime |
| **`FPS`** / **`GAME.FPS`** | int | `GetFPS` |
| **`ENDGAME`** / **`GAME.ENDGAME`** | — | Ends the running script (`TerminateVM`) |

Prefer **`TIME.DELTA()`** when you want the shared time module; **`DT()`** matches Raylib’s frame time and respects pause flags on the registry.

---

## Mouse (aliases included)

| Key | Returns |
|-----|---------|
| **`MX`**, **`MY`**, **`MOUSEX`**, **`MOUSEY`** | int |
| **`MDX`**, **`MDY`** | float (delta) |
| **`MWHEEL`** | float |
| **`MLEFT`**, **`MRIGHT`**, **`MMIDDLE`** | bool (down) |
| **`MLEFTPRESSED`**, **`MRIGHTPRESSED`** | bool (edge) |

Each has a matching **`GAME.*`** registration.

---

## Keyboard shortcuts (vs `INPUT.*`)

| Key | Args | Notes |
|-----|------|--------|
| **`KEYDOWN`**, **`KEYPRESSED`**, **`KEYRELEASED`** | key code | Same codes as **`INPUT.KEY*`** / **`KEY_*`** globals |
| **`KEYCHAR`** | — | `GetCharPressed` (Unicode codepoint) |
| **`ANYKEY`** | — | True if any key had a press this frame |

For full keyboard API and naming, see **[INPUT.md](INPUT.md)**.

---

## Pause and counters

| Key | Role |
|-----|------|
| **`PAUSEGAME`** / **`RESUMEGAME`** / **`GAMEPAUSE`** | Toggle **`Registry.GamePaused`** (affects **`DT`**) |
| **`FRAMECOUNT`** | **`Registry.FrameCount`** (increments on **`RENDER.FRAME`**) |
| **`ELAPSED`** | Seconds since module start |

---

## Pure math and collision (bare names)

Registered in **`register_math.go`**, **`register_collision.go`**, **`register_ease_noise_rand.go`**, **`register_color_format.go`**: movement helpers (**`NEWXVALUE`**, **`POINTDIR2D`**, …), **2D/3D distance**, **AABB/sphere/box** tests, **easing**, **noise**, **RNG**, **RGB** pack/unpack, etc. These are **numeric** only (no Raylib draw).

---

## `GAME.*` engine bridges

Separate dotted commands (single **`GAME.`** prefix in the registry), for example:

- **`GAME.DRAWSCREENFLASH`**, **`GAME.SCREENFLASH`** — full-screen flash (CGO)
- **`GAME.DEBUGRECT`** — debug rectangle overlay (CGO)
- **`GAME.ISGAMEPADAVAILABLE`**, **`GAME.GETGAMEPADNAME$`**
- **`GAME.SETMASTERVOLUME`**, **`GAME.GETMASTERVOLUME`**
- **`GAME.ISCURSORONSCREEN`**
- **`CONFIG.*`** — key/value save-game store (**`CONFIG.LOAD`**, **`CONFIG.GETINT`**, …)
- **`TIMER.*`**, **`STOPWATCH.*`**, simulated **`TIMER.MAKE`** / **`TIMER.UPDATE`** — see **[TIMER.md](TIMER.md)**

---

## See also

- [PROGRAMMING.md](../PROGRAMMING.md) — main loop
- [INPUT.md](INPUT.md) — `INPUT.KEYDOWN`, gestures
- [TIMER.md](TIMER.md) — timers and stopwatches
- [RAYLIB_EXTRAS.md](RAYLIB_EXTRAS.md) — window/render overview
