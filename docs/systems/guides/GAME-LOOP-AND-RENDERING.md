# Game loop, window, and rendering

> Open a window, run a frame loop, clear the screen, draw the 3D scene, and present — the skeleton every moonBASIC game shares.

**Namespaces:** `APP` · `RENDER` · `SCENE` · **Status:** Shipped · **Platform:** full runtime (`moonrun`)

**Commands:** [COMMAND_REGISTRY.md#core-window-time](../COMMAND_REGISTRY.md#core-window-time) · **Overview:** [01-CORE.md](../01-CORE.md) · [00-START.md](../00-START.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [APP — window and time](#app--window-and-time)
- [RENDER — each frame](#render--each-frame)
- [SCENE — multiple levels](#scene--multiple-levels)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Window, delta time, FPS cap, clear + 3D pass + swap buffers |
| **You need first** | `moonrun` from [Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) |
| **Typical games** | Every game with a visible window |
| **Not for** | Headless compile-only checks — use `moonbasic --check` |

**Why this trio exists:** `APP` owns **time and the OS window**. `RENDER` owns **GPU frame boundaries**. `SCENE` owns **which world** you draw when you have menus vs levels.

---

## When to use this system

**Use when:**

- You are building any real-time game or interactive demo.
- You need stable frame timing (`APP.DELTA()`, `APP.SETFPS`).
- You switch between hub, level, and pause states (`SCENE.*`).

**Skip when:**

- You only validate syntax — `moonbasic --check` never opens `APP`.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Open / close window | `APP.OPEN` / `APP.CLOSE` | Raw GLFW (not exposed) |
| Know if user quit | `APP.SHOULDCLOSE()` | Polling OS events manually |
| Clear + draw 3D | `RENDER.CLEAR` → `BEGIN` → `SCENE.DRAW` → `END` → `FRAME` | Drawing without `BEGIN`/`END` |
| Multiple level layouts | `SCENE.REGISTER` + `SCENE.SWITCH` | One giant `main.mb` forever |
| Place objects in world | `ENTITY.*` ([ENTITY-SYSTEM.md](ENTITY-SYSTEM.md)) | Only `RENDER` — entities live in scene graph |

---

## Core workflow

1. **`APP.OPEN(width, height, title)`** — create window.  
   **Why:** Allocates graphics context; nothing draws until this runs.

2. **`APP.SETFPS(60)`** (optional) — cap frame rate.  
   **Why:** Keeps `APP.DELTA()` stable for movement (`speed * APP.DELTA()`).

3. **Setup once** — camera, entities, lights (outside the loop).

4. **Loop until `APP.SHOULDCLOSE()`:**

   - Update game logic (input, physics, animation).
   - `RENDER.CLEAR(r, g, b)` — reset framebuffer.
   - `RENDER.BEGIN(cam)` — start 3D pass with projection.
   - `SCENE.DRAW()` — draw registered entities / scene content.
   - `RENDER.END()` — end 3D pass.
   - `RENDER.FRAME()` — present to screen.

5. **`APP.CLOSE()`** — release window on exit.

---

## APP — window and time

| Command | Why |
|---------|-----|
| `APP.OPEN(w, h, title)` | Window + graphics context |
| `APP.SHOULDCLOSE()` | Escape / close button → exit loop |
| `APP.SETFPS(n)` | Target frame rate |
| `APP.GETFPS()` | Debug overlay / adaptive quality |
| `APP.WIDTH()` / `HEIGHT()` | UI layout, mouse normalization |
| `APP.TIME()` | Seconds since start (effects, cutscenes) |
| `APP.DELTA()` | Seconds since last frame — **multiply all motion by this** |
| `APP.VERSION()` | Log build in save files or multiplayer handshake |

**Aliases:** `APP.*` maps to `WINDOW.*` / `TIME.*` — use checklist names in tutorials.

---

## RENDER — each frame

| Command | Why |
|---------|-----|
| `RENDER.CLEAR(r, g, b [, a])` | Solid background before draw |
| `RENDER.BEGIN([camera])` | 3D matrices + depth test for world |
| `RENDER.END()` | Close 3D pass before 2D HUD if needed |
| `RENDER.FRAME()` | Swap buffers — **required every loop** |
| `RENDER.SETBACKGROUND(r, g, b)` | Persistent clear color |
| `RENDER.SETWIREFRAME(true)` | Debug collision meshes |
| `RENDER.SCREENSHOT(path)` | Marketing / bug reports |

**Why `BEGIN`/`END`:** Separates world rendering from screen-space HUD (`DRAW.TEXT`, `GUI.*`) so layers do not fight depth settings.

---

## SCENE — multiple levels

| Command | Why |
|---------|-----|
| `SCENE.REGISTER(name)` | Named bucket for entities |
| `SCENE.SWITCH(name)` | Show one world state (level 1 vs menu) |
| `SCENE.DRAW()` | Draw active scene inside `RENDER.BEGIN` |
| `SCENE.SAVESCENE` / `LOADSCENE` | Serialize layout to disk |
| `SCENE.CLEARSCENE()` | Wipe entities when reloading |

**When to switch:** Hub → level → hub without freeing the whole `APP`.

---

## Full example

**Runnable copy:** [examples/guides/game_loop.mb](../../../examples/guides/game_loop.mb)

```basic
; Check: moonbasic --check examples/guides/game_loop.mb
; Run:   moonrun examples/guides/game_loop.mb

APP.OPEN(1280, 720, "Game loop demo")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -8)
CAMERA.LOOKAT(cam, 0, 0, 0)

cube = ENTITY.CREATECUBE(2, 2, 2)
cube.pos(0, 0, 5)

WHILE NOT APP.SHOULDCLOSE()
    cube.turn(0, 60 * APP.DELTA(), 0)
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

ENTITY.FREE(cube)
APP.CLOSE()
```

Also see [examples/foundation/main.mb](../../../examples/foundation/main.mb) (adds a point light).

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Forgot `RENDER.FRAME()` | Black or frozen window — call every loop |
| Motion without `APP.DELTA()` | Speed depends on FPS — multiply movement |
| `SCENE.DRAW` outside `RENDER.BEGIN` | Wrong projection or depth — bracket 3D draw |
| `APP.OPEN` after loop | Open once before `WHILE` |
| Expecting window from `--check` | Use `moonrun` for graphics |

---

## See also

- [ENTITY-SYSTEM.md](ENTITY-SYSTEM.md) — what `SCENE.DRAW` renders
- [CAMERA-AND-INPUT.md](CAMERA-AND-INPUT.md) — camera passed to `RENDER.BEGIN`
- [LIGHTING.md](LIGHTING.md) — lights inside the 3D pass
- [DEBUG-AND-TESTING.md](DEBUG-AND-TESTING.md) — `--check` vs `moonrun`
