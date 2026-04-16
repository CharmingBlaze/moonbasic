# moonBASIC examples (guide)

This page explains what each **repository demo** teaches and shows **short excerpts**. Full sources live under [`examples/`](../examples/README.md).

**Style:** Prefer **`NAMESPACE.ACTION`** registry names and **`CREATE`/`SETPOS`/`FREE`** in new code ([STYLE_GUIDE.md](../STYLE_GUIDE.md), [DOC_STYLE_GUIDE.md](DOC_STYLE_GUIDE.md)). Some excerpts below mix **Easy Mode** (`Window.Open`, `CreateCamera`, …) with namespace calls to match legacy demos—see [`examples/high_fidelity/modern_template.mb`](../examples/high_fidelity/modern_template.mb) for a more namespace-first template.

---

## How to run

From the repo root, **open a window** with the full runtime (not plain `go run .`, which only compiles to `.mbc`):

```bash
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

See [examples/README.md](../examples/README.md) for compile vs run, Windows PowerShell, and the full table.

---

## Concepts used everywhere

| Idea | Typical API (registry-first) |
|------|------------------------------|
| Open window | **`WINDOW.OPEN(w, h, title)`** — on failure the runtime may print to stderr and exit; use **`WINDOW.CANOPEN`** if you must branch without opening |
| Frame timing | **`WINDOW.SETFPS(60)`**, then each frame **`TIME.DELTA()`** for movement |
| Quit | **`INPUT.KEYDOWN(KEY_ESCAPE)`** or **`WINDOW.SHOULDCLOSE()`** |
| Clear + present | **`RENDER.CLEAR(r,g,b)`** … **`RENDER.FRAME()`** |
| 2D shapes | **`DRAW.RECTANGLE`**, **`DRAW.CIRCLE`**, … |
| HUD text (no font file) | **`DRAW.TEXT(msg, x, y, size, r, g, b, a)`** |

Command names are **case-insensitive** at compile time (`Draw.Text` and **`DRAW.TEXT`** resolve to the same builtin).

For **all commands by namespace**, naming conventions, and `SetPos` / `SetPosition` aliases, see [API_CONSISTENCY.md](API_CONSISTENCY.md) (`go run ./tools/apidoc`). For errors (did-you-mean, runtime line info) see [ERROR_MESSAGES.md](ERROR_MESSAGES.md). **Live HUD:** `DEBUG.WATCH(label, value)` each frame (on-screen overlay requires CGO; enable with `DEBUG.ENABLE` or host debug mode — see [DEBUG.md](reference/DEBUG.md)). Shortcuts and instant-game helpers (`SCREENW`, `DT`, collision math, timers, …) are documented under [QOL.md](reference/QOL.md).

---

## Modern Blitz-style 3D loop (aliases, CGO)

Full narrative: [GETTING_STARTED.md](GETTING_STARTED.md) (**Modern Blitz-style 3D**). Runnable template: [`examples/high_fidelity/modern_template.mb`](../examples/high_fidelity/modern_template.mb) (any resolution — 1080p shown). **`UPDATEPHYSICS`** = one-call frame tick (**`ENTITY.UPDATE`** + optional world / 2D / 3D steps). Command map: [MODERN_BLITZ_COMMANDS.md](reference/MODERN_BLITZ_COMMANDS.md).

The excerpt uses **Blitz / Easy Mode** scene helpers (`Graphics3D`, `CreateCamera`, …); **window, input, and render pass** use **`WINDOW.*`**, **`INPUT.*`**, **`RENDER.*`** like new projects.

```basic
; Initialize world
WINDOW.OPEN(1920, 1080, "Project: High Fidelity")
WINDOW.SETFPS(60)
AppTitle("Project: High Fidelity")
Graphics3D(1920, 1080)
SetMSAA(4)

cam = CreateCamera()
SetSSAO(TRUE)
car = LoadMesh("supercar.gltf")
EntityPBR(car, 0.9, 0.1)

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    CameraSmoothFollow(cam, car, 0.1)
    IF INPUT.KEYDOWN(KEY_W) THEN ApplyEntityImpulse(car, 0, 0, 500)
    UpdatePhysics()
    RENDER.CLEAR(10, 12, 18)
    RENDER.BEGIN3D(cam)
        DrawEntities()
    RENDER.END3D()
    RENDER.FRAME()
WEND
WINDOW.CLOSE()
```

---

## 3D spinning cube — `examples/spin_cube/main.mb`

Camera, mesh, material tint, transform matrix rotation, optional **ground grid**, and **cleanup** (`Mesh.Free`, `Material.Free`, `Transform.Free`, `Camera.Free`). The loop uses **registry** **`WINDOW.*`**, **`INPUT.*`**, **`TIME.*`**, **`RENDER.*`**; **`cam.Begin`** / **`Mesh.Draw`** match **Easy Mode** in the real file (see [CAMERA](reference/CAMERA.md), [DRAW3D](reference/DRAW3D.md)).

```basic
cam = CreateCamera()
cam.SetPos(0, 2, 8)
cam.SetTarget(0, 0, 0)
cubeMesh = Mesh.MakeCube(2, 2, 2)
cubeMat = Material.MakeDefault()
Material.SetColor(cubeMat, MATERIAL_MAP_ALBEDO, 130, 200, 255, 255)
cubeXform = Transform.Identity()

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    dt = TIME.DELTA()
    angle = angle + 1.1 * dt
    Transform.SetRotation(cubeXform, angle, angle * 0.65, angle * 0.35)
    RENDER.CLEAR(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cubeMesh, cubeMat, cubeXform)
        Draw.Grid(10, 1.0)
    cam.End()
    RENDER.FRAME()
WEND

Mesh.Free(cubeMesh)
Material.Free(cubeMat)
Transform.Free(cubeXform)
Camera.Free(cam)
WINDOW.CLOSE()
```

---

## 3D hop (orbit camera + platforms) — `examples/mario64/`

Third-person hop on a plane and boxes. **`examples/mario64/README.md`** compares sources: **`main.mb`** (implicit typing + **`Draw3D`** only, no entity graph), **`main_entities.mb`** (**`CreateCube`/`CreateSphere`**, **`COLLISIONS`**, **`EntityGrounded`** (coyote), **`EntityMoveCameraRelative`**, **`Camera.OrbitEntity`**, **`CopyEntity`** platform template, **`ENTITY.UPDATE`**, **`DrawEntities`**, child hat), plus Blitz-teaching variants **`main_orbit_simple.mb`**, **`main_v2.mb`**, **`main_v3.mb`**. Older variants use **`CAMERA.SETORBIT`**, **`INPUT.AXIS`**, **`MOVESTEPX`/`MOVESTEPZ`**, **`BOXTOPLAND`**, **`IIF`**. See [ENTITY.md](reference/ENTITY.md) (**`MoveEntity`** vs **`TranslateEntity`**, **`EntityHitsType`**), [CAMERA.md](reference/CAMERA.md), [INPUT.md](reference/INPUT.md), [GAMEHELPERS.md](reference/GAMEHELPERS.md), [MATH.md](reference/MATH.md), [LANGUAGE.md](LANGUAGE.md).

---

## 2D + mouse — `docs` snippet (not a separate file)

Use **`CAMERA2D.BEGIN()`** / **`CAMERA2D.END()`** for screen-space 2D (Raylib `BeginMode2D` / `EndMode2D`). **`RENDER.BEGINMODE2D`** / **`RENDER.ENDMODE2D`** are also registered for the Raylib mode stack when you need that path; see [RENDER](reference/RENDER.md), [CAMERA](reference/CAMERA.md), [SPRITE](reference/SPRITE.md).

```basic
WINDOW.OPEN(800, 600, "2D")
WINDOW.SETFPS(60)
WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    mx = INPUT.MOUSEX()
    my = INPUT.MOUSEY()
    RENDER.CLEAR(20, 20, 30)
    CAMERA2D.BEGIN()
        IF INPUT.MOUSEDOWN(MOUSE_LEFT_BUTTON) THEN
            DRAW.CIRCLE(100, 100, 50, 255, 100, 100, 255)
        ELSE
            DRAW.CIRCLE(100, 100, 50, 100, 200, 255, 255)
        ENDIF
        DRAW.RECTANGLE(mx - 25, my - 25, 50, 50, 255, 255, 255, 255)
        DRAW.TEXT("Hello, 2D", 200, 200, 20, 255, 255, 255, 255)
    CAMERA2D.END()
    RENDER.FRAME()
WEND
WINDOW.CLOSE()
```

---

## GUI (raygui) — `examples/gui_basics/main.mb`, `examples/gui_theme/main.mb`, `examples/gui_form/main.mb`

Requires CGO. See [GUI.md](reference/GUI.md). Use `GUI.THEMEAPPLY("CYBER")` (and other [raygui style](https://github.com/raysan5/raygui/tree/master/styles) names) for bundled themes; see `gui_theme/main.mb`.

```basic
GUI.Enable()
GUI.THEMEAPPLY("DARK")
IF GUI.BUTTON(20, 100, 120, 28, "OK") THEN
    status = "OK"
ENDIF
```

---

## File I/O — utilities

Quick text write/read without manual `FILE.*` handles:

```basic
ok = WRITEALLTEXT("out.txt", "hello")
data = READALLTEXT("out.txt")
```

For streaming I/O use `FILE.OPENREAD` / `FILE.OPENWRITE` (see [FILE](reference/FILE.md)).

---

## 3D physics — Linux + Jolt only

**`PHYSICS3D`** / **`BODY3D`** run on **Linux x64/arm64 with CGO** in this codebase. On Windows you will see stub errors until bindings exist.

```basic
PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -9.8, 0)
; ... BODY3D builder, COMMIT, STEP, draw with body matrix helpers ...
PHYSICS3D.STOP()
```

For cross-platform physics, start from [Physics 2D](reference/PHYSICS2D.md) instead.

---

## Other demos (short)

| Demo | File | Focus |
|------|------|--------|
| Pong | `examples/pong/main.mb` | Ball + paddles + score HUD (`DRAW.TEXT`) |
| Platformer | `examples/platformer/main.mb` | Gravity, ground, one platform |
| Arena | `examples/fps/main.mb` | WASD + oscillating targets + `TIME.GET` |
| Racing | `examples/racing/main.mb` | Steer / accelerate + lap line |
| Mini RPG | `examples/rpg/main.mb` | Move in a room, gold counter, `JSON` + `FILE` save on exit |
| 3D hop | `examples/mario64/README.md` | **`main.mb`** / **`main_entities.mb`** + teaching variants; entity sample uses **`EntityGrounded`**, **`EntityMoveCameraRelative`**, **`CAMERA.ORBITENTITY`**, **`CopyEntity`** — see [ENTITY.md](reference/ENTITY.md) |

---

## Next steps

- [Programming guide](PROGRAMMING.md) — structure, types, platforms  
- [Command index](COMMANDS.md) — topic index; [API_CONSISTENCY.md](API_CONSISTENCY.md) lists every registered builtin  
- [Getting started](GETTING_STARTED.md) — install and first window  
- [Camera / culling](reference/CAMERA.md) — CPU frustum and **`CULL.*`** (section **Culling and visibility**); [`testdata/culling_test.mb`](../testdata/culling_test.mb)  
