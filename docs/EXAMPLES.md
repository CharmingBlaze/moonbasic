# moonBASIC examples (guide)

This page explains what each **repository demo** teaches and shows **short excerpts**. Full sources live under [`examples/`](../examples/README.md).

---

## How to run

From the repo root, **open a window** with the full runtime (not plain `go run .`, which only compiles to `.mbc`):

```bash
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

See [examples/README.md](../examples/README.md) for compile vs run, Windows PowerShell, and the full table.

---

## Concepts used everywhere

| Idea | Typical API |
|------|-------------|
| Open window | `Window.Open(w, h, title$)` — on failure the runtime prints to stderr and exits; use `Window.CanOpen` if you must branch without opening |
| Frame timing | `Window.SetFPS(60)`, then each frame `Time.Delta()` for movement |
| Quit | `Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose()` |
| Clear + present | `Render.Clear(r,g,b)` … `Render.Frame()` |
| 2D shapes | `Draw.Rectangle`, `Draw.Circle`, … |
| HUD text (no font file) | `Draw.Text(msg$, x, y, size, r, g, b, a)` |

Command names are **case-insensitive** at compile time (`Draw.Text` = `DRAW.TEXT`).

For **all commands by namespace**, naming conventions, and `SetPos` / `SetPosition` aliases, see [API_CONSISTENCY.md](API_CONSISTENCY.md) (`go run ./tools/apidoc`). For errors (did-you-mean, runtime line info) see [ERROR_MESSAGES.md](ERROR_MESSAGES.md). **Live HUD:** `DEBUG.WATCH(label$, value)` each frame (on-screen overlay requires CGO; enable with `DEBUG.ENABLE` or host debug mode — see [DEBUG.md](reference/DEBUG.md)). Shortcuts and instant-game helpers (`SCREENW`, `DT`, collision math, timers, …) are documented under [QOL.md](reference/QOL.md).

---

## Modern Blitz-style 3D loop (aliases, CGO)

Full narrative: [GETTING_STARTED.md](GETTING_STARTED.md) (**Modern Blitz-style 3D**). Runnable template: [`examples/high_fidelity/modern_template.mb`](../examples/high_fidelity/modern_template.mb) (any resolution — 1080p shown). **`UpdatePhysics`** = one-call frame tick (`ENTITY.UPDATE` + optional world / 2D / 3D steps). Command map: [MODERN_BLITZ_COMMANDS.md](reference/MODERN_BLITZ_COMMANDS.md).

```basic
; Initialize world
Window.Open(1920, 1080, "Project: High Fidelity")
Window.SetFPS(60)
AppTitle("Project: High Fidelity")
Graphics3D(1920, 1080)
SetMSAA(4)

cam = CreateCamera()
SetSSAO(TRUE)
car = LoadMesh("supercar.gltf")
EntityPBR(car, 0.9, 0.1)

WHILE NOT (KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    CameraSmoothFollow(cam, car, 0.1)
    IF KeyDown(KEY_W) THEN ApplyEntityImpulse(car, 0, 0, 500)
    UpdatePhysics()
    Render.Clear(10, 12, 18)
    RENDER.Begin3D(cam)
        DrawEntities()
    RENDER.End3D()
    Render.Frame()
WEND
Window.Close()
```

---

## 3D spinning cube — `examples/spin_cube/main.mb`

Camera, mesh, material tint, transform matrix rotation, optional **ground grid**, and **cleanup** (`Mesh.Free`, `Material.Free`, `Transform.Free`, `Camera.Free`).

```basic
cam = Camera.Make()
cam.SetPos(0, 2, 8)
cam.SetTarget(0, 0, 0)
cubeMesh = Mesh.MakeCube(2, 2, 2)
cubeMat = Material.MakeDefault()
Material.SetColor(cubeMat, MATERIAL_MAP_ALBEDO, 130, 200, 255, 255)
cubeXform = Transform.Identity()

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    dt# = Time.Delta()
    angle# = angle# + 1.1 * dt#
    Transform.SetRotation(cubeXform, angle#, angle# * 0.65, angle# * 0.35)
    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cubeMesh, cubeMat, cubeXform)
        Draw.Grid(10, 1.0)
    cam.End()
    Render.Frame()
WEND

Mesh.Free(cubeMesh)
Material.Free(cubeMat)
Transform.Free(cubeXform)
Camera.Free(cam)
Window.Close()
```

---

## 3D hop (orbit camera + platforms) — `examples/mario64/`

Third-person hop on a plane and boxes. **`examples/mario64/README.md`** compares sources: **`main.mb`** (implicit typing + **`Draw3D`** only, no entity graph), **`main_entities.mb`** (**`CreateCube`/`CreateSphere`**, **`COLLISIONS`**, **`EntityGrounded`** (coyote), **`EntityMoveCameraRelative`**, **`Camera.OrbitEntity`**, **`CopyEntity`** platform template, **`ENTITY.UPDATE`**, **`DrawEntities`**, child hat), plus Blitz-teaching variants **`main_orbit_simple.mb`**, **`main_v2.mb`**, **`main_v3.mb`**. Older variants use **`Camera.SetOrbit`**, **`Input.Axis`**, **`MOVESTEPX`/`MOVESTEPZ`**, **`BOXTOPLAND`**, **`IIF$`**. See [ENTITY.md](reference/ENTITY.md) (**`MoveEntity`** vs **`TranslateEntity`**, **`EntityHitsType`**), [CAMERA.md](reference/CAMERA.md), [INPUT.md](reference/INPUT.md), [GAMEHELPERS.md](reference/GAMEHELPERS.md), [MATH.md](reference/MATH.md), [LANGUAGE.md](LANGUAGE.md).

---

## 2D + mouse — `docs` snippet (not a separate file)

Use **`Camera2D.Begin()`** / **`Camera2D.End()`** for screen-space 2D (Raylib `BeginMode2D` / `EndMode2D`). There is no `Render.BeginMode2D` in the runtime; see [RENDER](reference/RENDER.md) and [CAMERA](reference/CAMERA.md).

```basic
Window.Open(800, 600, "2D")
Window.SetFPS(60)
WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    mx = Input.MouseX()
    my = Input.MouseY()
    Render.Clear(20, 20, 30)
    Camera2D.Begin()
        IF Input.MouseDown(MOUSE_LEFT_BUTTON) THEN
            Draw.Circle(100, 100, 50, 255, 100, 100, 255)
        ELSE
            Draw.Circle(100, 100, 50, 100, 200, 255, 255)
        ENDIF
        Draw.Rectangle(mx - 25, my - 25, 50, 50, 255, 255, 255, 255)
        Draw.Text("Hello, 2D", 200, 200, 20, 255, 255, 255, 255)
    Camera2D.End()
    Render.Frame()
WEND
Window.Close()
```

---

## GUI (raygui) — `examples/gui_basics/main.mb`, `examples/gui_theme/main.mb`, `examples/gui_form/main.mb`

Requires CGO. See [GUI.md](reference/GUI.md). Use `GUI.THEMEAPPLY("CYBER")` (and other [raygui style](https://github.com/raysan5/raygui/tree/master/styles) names) for bundled themes; see `gui_theme/main.mb`.

```basic
GUI.Enable()
GUI.THEMEAPPLY("DARK")
IF GUI.BUTTON(20, 100, 120, 28, "OK") THEN
    status$ = "OK"
ENDIF
```

---

## File I/O — utilities

Quick text write/read without manual `FILE.*` handles:

```basic
ok = WRITEALLTEXT("out.txt", "hello")
data$ = READALLTEXT$("out.txt")
```

For streaming I/O use `FILE.OPENREAD` / `FILE.OPENWRITE` (see [FILE](reference/FILE.md)).

---

## 3D physics — Linux + Jolt only

`Physics3D` / `Body3D` run on **Linux x64/arm64 with CGO** in this codebase. On Windows you will see stub errors until bindings exist.

```basic
Physics3D.Start()
Physics3D.SetGravity(0, -9.8, 0)
; ... Body3D.Make, Commit, Step, draw with Body3D.GetMatrix ...
Physics3D.Stop()
```

For cross-platform physics, start from [Physics 2D](reference/PHYSICS2D.md) instead.

---

## Other demos (short)

| Demo | File | Focus |
|------|------|--------|
| Pong | `examples/pong/main.mb` | Ball + paddles + score HUD (`Draw.Text`) |
| Platformer | `examples/platformer/main.mb` | Gravity, ground, one platform |
| Arena | `examples/fps/main.mb` | WASD + oscillating targets + `TIME.GET` |
| Racing | `examples/racing/main.mb` | Steer / accelerate + lap line |
| Mini RPG | `examples/rpg/main.mb` | Move in a room, gold counter, `JSON` + `FILE` save on exit |
| 3D hop | `examples/mario64/README.md` | **`main.mb`** / **`main_entities.mb`** + teaching variants; entity sample uses **`EntityGrounded`**, **`EntityMoveCameraRelative`**, **`Camera.OrbitEntity`**, **`CopyEntity`** — see [ENTITY.md](reference/ENTITY.md) |

---

## Next steps

- [Programming guide](PROGRAMMING.md) — structure, types, platforms  
- [Command index](COMMANDS.md) — topic index; [API_CONSISTENCY.md](API_CONSISTENCY.md) lists every registered builtin  
- [Getting started](GETTING_STARTED.md) — install and first window  
- [Camera / culling](reference/CAMERA.md) — CPU frustum and `Cull.*` (section **Culling and visibility**); [`testdata/culling_test.mb`](../testdata/culling_test.mb)  
