# Start here — why these commands exist

> Before memorizing APIs, understand the **minimum game loop** and **why** each builtin is on the critical path.

**Next:** [01-CORE.md](01-CORE.md) · **All commands:** [COMMAND_REGISTRY.md](COMMAND_REGISTRY.md) · **Hub:** [../BEGIN_HERE.md](../BEGIN_HERE.md)

**Case:** All command names are **case-insensitive**.

---

## Table of contents

- [What you need installed](#what-you-need-installed)
- [Create a project](#create-a-project)
- [Foundation program explained](#foundation-program-explained)
- [Command groups at a glance](#command-groups-at-a-glance)
- [When to add more systems](#when-to-add-more-systems)
- [Verify your script](#verify-your-script)
- [See also](#see-also)

---

## What you need installed

| Piece | Why |
|-------|-----|
| **Full runtime** zip from [Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) | Includes **`moonrun`** to open windows and run games |
| **`moonbasic`** (in the same zip) | Check, compile, LSP, `moonbasic new` |
| Optional **VSIX** | Syntax + completions in VS Code |

You do **not** need to build the compiler from source to make games.

---

## Create a project

```bash
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

**Why `moonbasic new`?** It creates a standard layout (`main.mb`, `assets/`, debug config) so paths and tooling work the same for every beginner.

Templates: `moonbasic new --template 3d|platformer|ui|physics MyGame` — see [11-TOOLING.md](11-TOOLING.md).

---

## Foundation program explained

This is the smallest **3D** loop most games grow from ([`examples/foundation/main.mb`](../../examples/foundation/main.mb)):

```basic
APP.OPEN(1280, 720, "Test")
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

APP.CLOSE()
```

### Line by line — what and why

| Code | What it does | Why you need it |
|------|----------------|-----------------|
| `APP.OPEN(1280, 720, "Test")` | Opens window 1280×720 | Without a window, Raylib cannot create a GL context — **no graphics**. Alias of `WINDOW.OPEN`. |
| `APP.SETFPS(60)` | Targets 60 frames per second | Keeps `DELTA()` predictable so movement feels the same on fast and slow PCs. |
| `CAMERA.CREATE()` | Creates a view camera handle | 3D drawing needs eye position, orientation, and projection — the camera stores that. |
| `CAMERA.SETACTIVE(cam)` | Default camera for `RENDER.BEGIN()` | So you can call `RENDER.BEGIN()` without repeating the handle every frame. |
| `CAMERA.SETPOS` / `LOOKAT` | Aim the camera | You must place the eye and look-at point or you may stare at empty space or inside geometry. |
| `ENTITY.CREATECUBE(2,2,2)` | Mesh + entity in one step | **Entities** are the object model — everything visible is an entity (or drawn through one). |
| `cube.pos(0, 0, 5)` | Move cube along +Z | Handle shortcut for `ENTITY.SETPOS` — positions matter for visibility and collision. |
| `WHILE NOT APP.SHOULDCLOSE()` | Main loop | Games run continuously until quit; this tests the close button / OS exit. |
| `cube.turn(0, 60 * APP.DELTA(), 0)` | Rotate 60°/sec on Y | **`APP.DELTA()`** = seconds since last frame — multiply speed so rotation is smooth at any FPS. |
| `RENDER.CLEAR(20, 20, 30)` | Fill background color | Clears color + depth buffers; old pixels would smear without this. |
| `RENDER.BEGIN(cam)` | Start 3D pass with this camera | Sets view/projection matrices for depth-correct 3D. |
| `SCENE.DRAW()` | Draw registered entities | Batches entity drawing for the active scene graph. |
| `RENDER.END()` | End 3D pass | Restores 2D state if you draw HUD after 3D. |
| `RENDER.FRAME()` | Swap buffers / present | **Mandatory** each frame — otherwise the window never updates. |
| `APP.CLOSE()` | Shutdown window | Releases display after the loop; good practice on desktop. |

### 2D-only games

Skip `CAMERA` / `RENDER.BEGIN` and use:

```basic
WINDOW.OPEN(800, 600, "2D")
WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(40, 44, 52)
    DRAW.RECTANGLE(100, 100, 200, 80, 255, 100, 50, 255)
    DRAW.TEXT("Hello", 20, 20, 18, 255, 255, 255)
    RENDER.FRAME()
WEND
WINDOW.CLOSE()
```

**Why `DRAW.*` instead of `ENTITY`?** For simple 2D, immediate draw calls are enough; entities help when you have many objects, hierarchy, and 3D.

---

## Command groups at a glance

| You want to… | Start with namespace | Guide |
| **Full walkthrough** | — | [GUIDES.md](GUIDES.md) (entity, 2D/3D collision, UI, multiplayer) |
|--------------|----------------------|-------|
| Window, time, quit | `APP.*` / `WINDOW.*`, `TIME.*` | [01-CORE](01-CORE.md) |
| Clear screen, 3D pass, show frame | `RENDER.*` | [01-CORE](01-CORE.md) |
| Game objects | `ENTITY.*` | [01-CORE](01-CORE.md) |
| Where the player looks | `CAMERA.*` | [02-CAMERA-LIGHT](02-CAMERA-LIGHT.md) |
| Lit 3D scenes | `LIGHT.*` | [02-CAMERA-LIGHT](02-CAMERA-LIGHT.md) |
| Images, models, packs | `TEXTURE.*`, `MODEL.*`, `ASSET.*` | [03-ASSETS](03-ASSETS.md) |
| Keyboard / gamepad | `INPUT.*`, `ACTION.*` | [04-INPUT](04-INPUT.md) |
| Gravity, collisions, clicks | `PHYSICS.*`, `BODY.*`, `PICK.*` | [05-PHYSICS](05-PHYSICS.md) |
| Sound | `AUDIO.*`, `SOUND.*` | [06-AUDIO](06-AUDIO.md) |
| Tiles, terrain, particles | `SPRITE.*`, `TILEMAP.*`, … | [07-2D-WORLD](07-2D-WORLD.md) |
| Menus, HUD text | `GUI.*`, `DRAW.TEXT` | [08-UI-TEXT](08-UI-TEXT.md) |
| Saves, JSON, math | `SAVE.*`, `FILE.*`, `JSON.*` | [09-DATA](09-DATA.md) |
| Logs, timers | `DEBUG.*`, `TIMER.*` | [10-DEBUG-TIMER](10-DEBUG-TIMER.md) |
| New project, package | CLI | [11-TOOLING](11-TOOLING.md) |

**Full command list per group:** [COMMAND_REGISTRY.md](COMMAND_REGISTRY.md)  
**Entire engine:** [API_CONSISTENCY.md](../API_CONSISTENCY.md)

---

## When to add more systems

Recommended order (from the foundation checklist):

1. Core loop works (window + clear + frame).
2. **DEBUG** — `DEBUG.LOG` / `DEBUG.WATCH` while tuning.
3. **INPUT** — move something with keys.
4. **ASSET** pack — stop hard-coding file paths.
5. **PHYSICS** — when position should respect collisions.
6. **AUDIO** — feedback and music.
7. **SAVE** — persist progress.

Do **not** block on networking, visual editors, or advanced shaders until the loop above is solid.

---

## Verify your script

```bash
moonbasic --check main.mb
```

**Why check?** Catches typos (`ENTITY.SETPOSITON`), wrong argument counts, and bad types **before** you run — faster than hunting a black screen.

Run the game:

```bash
moonrun main.mb
```

---

## See also

- [COMMAND_REGISTRY.md](COMMAND_REGISTRY.md) — every beginner-namespace command
- [README.md](README.md) — 40-system index
- [../FIRST_HOUR.md](../FIRST_HOUR.md) — language basics
