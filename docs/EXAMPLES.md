# moonBASIC examples (guide)

This page explains what each **repository demo** teaches and shows **short excerpts**. Full sources live under [`examples/`](../examples/README.md) — run them with `CGO_ENABLED=1` from the repo root.

---

## How to run

```bash
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

See [examples/README.md](../examples/README.md) for the full table and Windows notes.

---

## Concepts used everywhere

| Idea | Typical API |
|------|-------------|
| Open window | `Window.Open(w, h, title$)` — check return or use `IF NOT Window.Open(...) THEN ... ENDIF` |
| Frame timing | `Window.SetFPS(60)`, then each frame `Time.Delta()` for movement |
| Quit | `Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose()` |
| Clear + present | `Render.Clear(r,g,b)` … `Render.Frame()` |
| 2D shapes | `Draw.Rectangle`, `Draw.Circle`, … |
| HUD text (no font file) | `Draw.Text(msg$, x, y, size, r, g, b, a)` |

Command names are **case-insensitive** at compile time (`Draw.Text` = `DRAW.TEXT`).

---

## 3D spinning cube — `examples/spin_cube/main.mb`

Camera, mesh, material tint, matrix rotation, optional **ground grid**, and **cleanup** (`Mesh.Free`, `Material.Free`, `Mat4.Free`, `Camera.Free`).

```basic
cam = Camera.Make()
cam.SetPos(0, 2, 8)
cam.SetTarget(0, 0, 0)
cubeMesh = Mesh.MakeCube(2, 2, 2)
cubeMat = Material.MakeDefault()
Material.SetColor(cubeMat, MATERIAL_MAP_ALBEDO, 130, 200, 255, 255)
rot = Mat4.Identity()

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    dt# = Time.Delta()
    angle# = angle# + 1.1 * dt#
    Mat4.SetRotation(rot, angle#, angle# * 0.65, angle# * 0.35)
    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cubeMesh, cubeMat, rot)
        Draw.Grid(10, 1.0)
    cam.End()
    Render.Frame()
WEND

Mesh.Free(cubeMesh)
Material.Free(cubeMat)
Mat4.Free(rot)
Camera.Free(cam)
Window.Close()
```

---

## 2D + mouse — `docs` snippet (not a separate file)

Use `Render.BeginMode2D()` when you want pixel coordinates for UI-style 2D (see [RENDER](reference/RENDER.md)).

```basic
Window.Open(800, 600, "2D")
Window.SetFPS(60)
WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    mx = Input.MouseX()
    my = Input.MouseY()
    Render.Clear(20, 20, 30)
    Render.BeginMode2D()
        IF Input.MouseDown(MOUSE_LEFT_BUTTON) THEN
            Draw.Circle(100, 100, 50, 255, 100, 100, 255)
        ELSE
            Draw.Circle(100, 100, 50, 100, 200, 255, 255)
        ENDIF
        Draw.Rectangle(mx - 25, my - 25, 50, 50, 255, 255, 255, 255)
        Draw.Text("Hello, 2D", 200, 200, 20, 255, 255, 255, 255)
    Render.EndMode2D()
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

---

## Next steps

- [Programming guide](PROGRAMMING.md) — structure, types, platforms  
- [Command index](COMMANDS.md) — look up any builtin  
- [Getting started](GETTING_STARTED.md) — install and first window  
