# Programming in moonBASIC

This guide explains how **built-in commands** fit together, how to structure programs, and where to look up APIs. Pair it with [Language Reference](LANGUAGE.md) (syntax), [Command Index](COMMANDS.md) (topic index), and [API_CONSISTENCY.md](API_CONSISTENCY.md) (every registered command name and arity, generated from the manifest).

---

## 1. Commands are `NAMESPACE.NAME`

Built-ins look like method calls: `Window.Open(...)`, `Draw.Rectangle(...)`, `TIME.DELTA()`.

The lexer **folds names to uppercase** for lookup. These are the same call:

- `Window.Open` → `WINDOW.OPEN`
- `draw.rectangle` → `DRAW.RECTANGLE`

Use whatever style reads best; examples in the repo often use **Mixed.Case** for namespaces.

**Consistent verbs across types** (`Load` / `SetPos` / `Free`, when to use `Make` vs `Load`, and how rotation differs for cameras vs models): see [API conventions](reference/API_CONVENTIONS.md).

---

## 2. Arguments and types

Commands are **type-checked** against the manifest (`compiler/builtinmanifest/commands.json`). Typical argument kinds:

| Kind | In source | Example |
|------|-----------|---------|
| Integer | `score`, literal `10` | `ARRAYLEN(arr)` |
| Float | `x#`, `1.5` | `MATH.SIN(angle#)` |
| String | `msg$`, `"hi"` | `FILE.OPEN(path$, "r")` |
| Boolean | `ok?`, `TRUE` / `FALSE` | `Input.KeyDown(KEY_SPACE)` |
| Handle | value from `Load`, `Make`, etc. | `Mesh.Draw(mesh, mat, transform)` |

Numeric **widening** is allowed where the manifest marks alternatives (many APIs accept int or float for coordinates).

Variable types follow **suffix rules** (`$` string, `#` float, `?` bool, none = int) — see [LANGUAGE.md](LANGUAGE.md).

---

## 3. The usual game / app loop

Almost all graphical programs follow this shape:

```basic
Window.Open(960, 540, "Title")
Window.SetFPS(60)

; setup handles, load assets, set variables

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    dt# = Time.Delta()

    ; --- update (physics, input, AI) ---

    ; --- draw ---
    Render.Clear(r, g, b)
    ; optional: Render.BeginMode2D() for screen-space 2D
    ; optional: cam.Begin() for 3D
    Draw.Rectangle(...)
    ; optional: cam.End() / Render.EndMode2D()
    Render.Frame()
WEND

; free heap handles (fonts, meshes, textures) if you loaded any
Window.Close()
```

Rules of thumb:

- **`Render.Clear`** — first drawing call each frame (or after `BeginMode2D` / camera, depending on your pipeline).
- **`Render.Frame`** — last call each frame; swaps / presents the buffer.
- **`Time.Delta()`** — seconds since last frame; multiply speeds by `dt#` for **frame-rate-independent** motion.
- **`Window.ShouldClose()`** — true when the user closes the window.
- **`Input.KeyDown(KEY_ESCAPE)`** — common explicit quit.

moonBASIC does **not** provide a hidden **`Game.Loop()`** / **`Game.Begin()`** / **`Game.End()`** wrapper: the **`WHILE`** + **`dt#`** pattern stays visible so you control ordering, pausing, and multi-pass rendering. Helpers like **`Input.Orbit`**, **`LANDBOXES`**, and **`MOVESTEPX`** shorten the *body*, not the loop shell.

---

## 4. 2D vs 3D drawing

- **Screen-space 2D** (pixels): wrap drawing in `Render.BeginMode2D()` … `Render.EndMode2D()` when you want coordinates in window pixels (see [RENDER](reference/RENDER.md)).
- **3D**: create `cam = Camera.Make()`, configure position/target/FOV, then `cam.Begin()` … `cam.End()` around `Mesh.Draw`, `Draw.Grid`, etc. (see [CAMERA](reference/CAMERA.md), [MODEL](reference/MODEL.md)).

Some 3D helpers are also registered under `DRAW.*` (e.g. `Draw.Grid` inside a camera block).

---

## 5. Text without shipping a font file

`Draw.Text(text$, x, y, size, r, g, b, a)` uses Raylib’s **default font** — no `.ttf` path required. Use this in small demos and HUD.

For a **custom** font, `Font.Load(path$)` returns a handle; draw with `Draw.TextEx` / `Draw.TextFont` style APIs (see [FONT](reference/FONT.md)). The repo **does not** ship `.ttf` files under `assets/`; add your own or rely on `Draw.Text`.

---

## 6. GUI (`GUI.*`)

`GUI.*` wraps **raygui** when **CGO** is enabled. On **Windows** with **`CGO_ENABLED=0`**, a **minimal** Raylib-drawn `GUI.*` subset runs instead (not full raygui); see [GUI.md](reference/GUI.md).

- The [GUI reference](reference/GUI.md) is the full catalog: **every `GUI.*` command**, **how to theme and restyle** (`GUI.THEMEAPPLY`, `SETCOLOR`, `SETSTYLE`, `GCTL_*` / `GPROP_*`), and **stateful array handles** (`SCROLLPANEL`, `LISTVIEW`, `DROPDOWNBOX`, …). Use **`GUI.THEMENAMES$`** for the list of built-in / bundled theme names.
- Runnable demos: `examples/gui_basics/main.mb`, `examples/gui_theme/main.mb`, `examples/gui_form/main.mb`.

---

## 7. Platform and build flags

| Area | Notes |
|------|--------|
| **Graphics, audio, window** | **Linux / macOS:** **CGO** + C toolchain (linked Raylib). **Windows:** either **CGO** + MinGW, or **`CGO_ENABLED=0`** + **`raylib.dll`** (purego; see [BUILDING.md](BUILDING.md)). **`GUI.*`**: full **raygui** needs **CGO**; **Windows + no CGO** uses the minimal GUI layer. |
| **Physics 3D** (`Physics3D`, `Body3D`) | Implemented on **Linux x64/arm64** with Jolt; other OS builds use stubs until bindings exist. |
| **Physics 2D** | Box2D path — see [PHYSICS2D](reference/PHYSICS2D.md). |
| **gopls / IDE** | If the editor analyzes with `CGO_ENABLED=0`, Raylib symbols may look “missing”; set `buildFlags`: `["-tags=cgo"]` and enable CGO where possible. |

---

## 8. Arrays, `DIM`, and handles

- **`DIM a(10)`** — numeric array; indices `0` … `9`.
- **`DIM plat AS Platform(4)`** — array of a **record type** defined with **`TYPE` … `ENDTYPE`** (see [LANGUAGE.md](LANGUAGE.md)). Use **`plat(i) = Platform(...)`** and **`plat(i).field`**.
- Some builtins return **handles** to heap arrays (e.g. `MEASURETEXTEX`, `GUI.GETCOLOR`). Index with the same `arr(i)` syntax as `DIM` arrays.
- **`ERASE(name)`** — frees a `DIM` array or typed array and clears the variable when you no longer need it.
- **`ERASE ALL`** / **`FREE.ALL`** — frees every VM heap object and nulls handle variables; see [MEMORY.md](MEMORY.md).
- **`ARRAYFREE(handle)`** when you are done with a heap array you no longer need.

---

## 9. Where to look things up

| Need | Document |
|------|----------|
| Syntax (`IF`, `FUNCTION`, …) | [LANGUAGE.md](LANGUAGE.md) |
| Topic command index | [COMMANDS.md](COMMANDS.md) |
| Every manifest name (arity, types) | [API_CONSISTENCY.md](API_CONSISTENCY.md) (`go run ./tools/apidoc`) |
| Namespace → reference map (counts, blurbs) | [COMMAND_AUDIT.md](COMMAND_AUDIT.md) (`go run ./tools/cmdaudit`) |
| Consistent verbs (`LOAD`, `SETPOS`, …) | [reference/API_CONVENTIONS.md](reference/API_CONVENTIONS.md) |
| Copy-paste samples | [EXAMPLES.md](EXAMPLES.md) |
| Install & first run | [GETTING_STARTED.md](GETTING_STARTED.md) |
| Deep dive per topic | [reference/](reference/WINDOW.md) (module pages) |
| Handles, leaks, `FreeAll` | [MEMORY.md](MEMORY.md) |
| 2D physics tuning (`SetStep`, `SetIterations`) | [reference/PHYSICS2D.md](reference/PHYSICS2D.md) |
| Purego `GUI.*` (stable rects, internal caps) | [reference/GUI.md](reference/GUI.md) |

---

## 10. Performance checklist

Use this alongside the loop in **§3**:

- **Motion and animation** — Multiply speeds by **`Time.Delta()`** so gameplay stays consistent when FPS changes.
- **2D physics** — Call **`Physics2D.Step()`** once per frame in the common case; set **`Physics2D.SetStep(dt#)`** to match that step (e.g. `1/60` with **`Window.SetFPS(60)`**). Tune cost vs stability with **`Physics2D.SetIterations`** — see [PHYSICS2D.md](reference/PHYSICS2D.md).
- **Heap handles** — Call **`*.Free`** for textures, fonts, sounds, and other handles when you are done, especially in long sessions or when reloading assets. **`Window.Close`** and process shutdown still run **`Heap.FreeAll`** as a safety net — see [MEMORY.md](MEMORY.md).
- **Churn** — Avoid creating many new handles or large temporary work every frame when you can reuse values or keep allocations outside the inner loop.
- **Assets** — Prefer texture sizes and counts appropriate for the target resolution; fewer draw state changes usually help.
- **Platform** — On **Windows**, **`CGO_ENABLED=0`** builds need **`raylib.dll`** on the DLL search path. Full **`Physics3D`** (Jolt) is only on **Linux** with CGO today — see **§7**.

---

## 11. Running repository demos

From the **repository root** (so relative paths behave as documented):

```bash
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

On Windows (PowerShell):

```powershell
$env:CGO_ENABLED="1"
go run . examples\spin_cube\main.mb
```

See [examples/README.md](../examples/README.md) for the full list.
