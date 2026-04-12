# Getting Started with moonBASIC

## What you need

- **Windows x64** or **Linux x64** — project docs and contributor workflows assume **Windows first**, **Linux** second for parity checks ([DEVELOPER.md](DEVELOPER.md#platform-priority-windows-then-linux)).
- **Either** a pre-built binary from [GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest) **or** Go + a C toolchain to build from source (below).

---

## Install (pre-built binaries)

Download the latest **`moonbasic-v*-windows-amd64.zip`** (Windows) or **`moonbasic-v*-linux-amd64.tar.gz`** (Linux) from **[Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)**. Extract anywhere; each archive includes **`README-RELEASE.txt`** with paths, `chmod` on Linux, and example commands.

| In the archive | Purpose |
|----------------|---------|
| `moonbasic` / `moonbasic.exe` | Compiler: `.mb` → `.mbc`, `--check`, `--lsp` |
| `moonrun` / `moonrun.exe` | Full game runtime (run `.mb` / `.mbc`) |

No Go installation required for this path.

---

## What you need (build from source)

- Go **1.25.3** or later (see `go.mod` in the repo) — https://go.dev/dl/
- A C compiler:
  - Windows: MinGW-w64 — https://www.mingw-w64.org/
  - Linux: GCC — usually already installed

---

## Install (build from source)

### Windows

1. Install Go from https://go.dev/dl/
2. Install MSYS2 from https://www.msys2.org/
3. Open MSYS2 MINGW64 terminal and run:
```bash
pacman -S mingw-w64-x86_64-gcc
```
4. Clone and build moonBASIC:
```bat
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic
set CGO_ENABLED=1
set CC=C:\msys64\mingw64\bin\gcc.exe
go build -o moonbasic.exe .
```
5. Add `moonbasic.exe` to your PATH (or run it from the project folder).

### Linux (Ubuntu / Debian)
```bash
sudo apt-get update
sudo apt-get install -y gcc libgl1-mesa-dev libxi-dev \
  libxcursor-dev libxrandr-dev libxinerama-dev \
  libwayland-dev libxkbcommon-dev

git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic
CGO_ENABLED=1 go build -o moonbasic .
sudo cp moonbasic /usr/local/bin/
```

### Verify installation
```bash
moonbasic --version
```

Expected output:
```
moonBASIC vX.X.X
```

---

## Your First Program

Create a file called `hello.mb`:

```basic
PRINT("Hello, moonBASIC!")
```

Run it (needs the **runtime** — `moonbasic` alone only compiles to `.mbc`):

```bash
moonrun hello.mb
```

Output:

```
Hello, moonBASIC!
```

From source without installing binaries: `go run -tags fullruntime ./cmd/moonrun hello.mb`. To only compile: `moonbasic hello.mb` → writes `hello.mbc`.

---

## Your First Window

Create `window.mb`:

```basic
Window.Open(960, 540, "Hello Window")
Window.SetFPS(60)

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    Render.Clear(30, 40, 60)
    ; Default bitmap font — no .ttf file required
    Draw.Text("Press ESC or close the window", 260, 260, 20, 255, 255, 255, 255)
    Render.Frame()
WEND

Window.Close()
```

Run it with the **game runtime** (plain `moonbasic` only compiles to `.mbc`):

```bash
moonrun window.mb
```

If you build from source: `go build -tags fullruntime -o moonrun ./cmd/moonrun`, then `moonrun window.mb`. From the repo without installing binaries: `CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun window.mb`.

You should see a dark blue window with white text. Press the window's close button to exit.

If opening the window fails (invalid size, unavailable display, etc.), the runtime prints a short message to **stderr** and the process exits—you do not need `IF NOT Window.Open …` in every program. Use **`Window.CanOpen`** only when you must choose a fallback resolution or show your own error without opening.

---

## Your First Game Loop

moonBASIC games follow a simple pattern:

```basic
; 1. Setup — runs once
Window.Open(800, 600, "My Game")
Window.SetFPS(60)

player_x = 400
player_y = 300

; 2. Loop — runs every frame until the window is closed
WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    dt = Time.Delta()         ; seconds since last frame

    ; Update
    IF Input.KeyDown(KEY_RIGHT) THEN player_x = player_x + 200 * dt
    IF Input.KeyDown(KEY_LEFT)  THEN player_x = player_x - 200 * dt
    IF Input.KeyDown(KEY_DOWN)  THEN player_y = player_y + 200 * dt
    IF Input.KeyDown(KEY_UP)    THEN player_y = player_y - 200 * dt

    ; Draw
    Render.Clear(20, 30, 40)
    Draw.Rectangle(INT(player_x) - 16, INT(player_y) - 16, 32, 32, 100, 200, 255, 255)
    Render.Frame()
WEND

; 3. Cleanup
Window.Close()
```

Key concepts:

| Concept | Explanation |
|---|---|
| `Window.ShouldClose()` | Returns `TRUE` when the player clicks the X or presses Alt+F4. |
| `Input.KeyDown(KEY_ESCAPE)` | Common way to exit demos with the keyboard. |
| `Time.Delta()` | Seconds since last frame. Multiply speeds by this for frame-rate-independent movement. |
| `Render.Clear(r, g, b)` | Clears the screen to a color. Always call this first in the loop. |
| `Render.Frame()` | Shows what was drawn. Always call this last in the loop. |
| `Draw.Rectangle(x, y, w, h, r, g, b, a)` | Draws a filled colored rectangle. |
| `Input.KeyDown(KEY_*)` | Returns `TRUE` while a key is held. |

---

## Adding a 3D Camera

For 3D scenes, wrap your drawing commands in `cam.Begin()` / `cam.End()`:

```basic
Window.Open(960, 540, "3D Cube")
Window.SetFPS(60)

cam = Camera.Make()
cam.SetPos(0, 3, 8)
cam.SetTarget(0, 0, 0)
cam.SetFOV(45)

cube = Mesh.MakeCube(2, 2, 2)
mat  = Material.MakeDefault()
xform = Transform.Identity()
angle = 0.0

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    angle = angle + 1.5 * Time.Delta()
    Transform.SetRotation(xform, angle, angle * 0.7, 0)

    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cube, mat, xform)
        Draw3D.Grid(20, 1.0)
    cam.End()
    Render.Frame()
WEND

Mesh.Free(cube)
Material.Free(mat)
Transform.Free(xform)
Camera.Free(cam)
Window.Close()
```

---

## Modern Blitz-style 3D (CGO)

This pattern uses **global aliases** (`Graphics3D`, `CreateCamera`, `LoadMesh`, `RENDER.Begin3D`, `UpdatePhysics`, …) on top of the same engine as `Window.*` / `ENTITY.*`.

- **Resolution is your choice** — `Window.Open(w, h, …)` accepts any reasonable size (720p, 1080p, 1440p, 4K, …). You must **`Window.Open`** before **`Graphics3D`**: the latter **only resizes** the client area (omit **`Graphics3D`** if the open size is already what you want).
- **`SetMSAA`**: best-effort; for some drivers MSAA is fixed at **`Window.Open`** — see [WINDOW](reference/WINDOW.md) / [BUILDING](BUILDING.md).
- **`SetSSAO` / PBR / lights**: need **CGO Raylib**; stubs return errors on non-graphical builds.
- **`UpdatePhysics`** (same as **`UPDATEPHYSICS`**) runs **`ENTITY.UPDATE(Time.Delta)`** and best-effort **`WORLD.UPDATE`**, **`PHYSICS2D.STEP`**, **`PHYSICS3D.STEP`** (inactive worlds no-op or ignored).
- **Frame contract**: **`UpdatePhysics`** → **`Render.Clear`** → 3D pass → **`Render.Frame`**.

```basic
; Initialize world
Window.Open(1920, 1080, "Project: High Fidelity")
Window.SetFPS(60)
AppTitle("Project: High Fidelity")
Graphics3D(1920, 1080)   ; optional resize (omit if Open already matched size)
SetMSAA(4)               ; Clean edges

; Setup Scene
cam = CreateCamera()
SetSSAO(TRUE)            ; Modern shadows

; Load a high-poly PBR model (paths relative to working directory)
car = LoadMesh("supercar.gltf")
EntityPBR(car, 0.9, 0.1)                    ; Shiny Chrome look
EntityNormalMap(car, LoadTexture("car_normals.png"))

; Attach dynamic headlights
L_Light = CreatePointLight(car, 255, 255, 200)
TranslateEntity(L_Light, -1, 0, 2)

WHILE NOT (KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    ; Modern Camera Follow logic
    CameraSmoothFollow(cam, car, 0.1)

    ; Physics Impulse based on modern input
    IF KeyDown(KEY_W) THEN ApplyEntityImpulse(car, 0, 0, 500)

    UpdatePhysics()

    ; The Render Pass
    Render.Clear(12, 14, 22)
    RENDER.Begin3D(cam)
        DrawEntities()   ; PBR, SSAO, dynamic lights (when CGO + assets load)
    RENDER.End3D()
    Render.Frame()
WEND

Window.Close()
```

Place **`assets/`** paths next to your `.mb` or use paths relative to the process working directory. For manual **`PHYSICS3D.START`** + **`BODY3D`**, you can still call **`PHYSICS3D.STEP`** explicitly; **`UpdatePhysics`** already invokes it with **`Time.Delta`** when the world is active. See [PHYSICS3D](reference/PHYSICS3D.md) and [EXAMPLES](EXAMPLES.md).

---

## Variable Types

Variables types are determined implicitly by the value assigned on first use.

| Type | Example |
|---|---|
| String | `name = "Player"` |
| Float | `speed = 5.5` |
| Boolean | `alive = TRUE` |
| Integer | `score = 100` |

There is no separate declaration step — variables are created on first assignment.
The language is dynamically typed (implicit `Any`).

---

## Key Constants

### Keyboard

Common key constants: `KEY_W`, `KEY_A`, `KEY_S`, `KEY_D`, `KEY_UP`, `KEY_DOWN`,
`KEY_LEFT`, `KEY_RIGHT`, `KEY_SPACE`, `KEY_ESCAPE`, `KEY_ENTER`, `KEY_LSHIFT`.

Use `Input.KeyDown(key)` for held, `Input.KeyPressed(key)` for first-press only.

### Mouse

`MOUSE_LEFT_BUTTON`, `MOUSE_RIGHT_BUTTON`, `MOUSE_MIDDLE_BUTTON`

Use `Input.MouseDown(btn)`, `Input.MouseX()`, `Input.MouseY()`.

### Material Map Slots

`MATERIAL_MAP_ALBEDO` (also `MAP_DIFFUSE`) — base color texture/tint  
`MATERIAL_MAP_METALNESS` — metalness  
`MATERIAL_MAP_ROUGHNESS` — roughness

---

## Project Layout

There is no mandatory project layout. A single `.mb` file is a complete program.
For larger projects, a common convention is:

```
mygame/
  main.mb          ← entry point
  assets/
    fonts/
    textures/
    sounds/
    maps/
```

Run from the project root:
```bash
moonbasic mygame/main.mb
```

---

## Next Steps

| Topic | Where to go |
|---|---|
| How commands fit together (loop, 2D/3D, CGO) | [Programming Guide](PROGRAMMING.md) |
| Language syntax (IF, FOR, FUNCTION…) | [Language Reference](LANGUAGE.md) |
| Command index by topic | [Command Index](COMMANDS.md) |
| Full manifest (every name + arity) | [API_CONSISTENCY.md](API_CONSISTENCY.md) |
| Copy-paste snippets + narrative | [Examples](EXAMPLES.md) |
| Runnable programs under `examples/` | [examples/README.md](../examples/README.md) |
| 2D drawing | [Draw Reference](reference/DRAW2D.md) |
| 3D models & meshes | [Model Reference](reference/MODEL.md) |
| 2D physics (Box2D) | [Physics 2D](reference/PHYSICS2D.md) |
| 3D physics (Jolt) | [Physics 3D](reference/PHYSICS3D.md) |
| Multiplayer (ENet) | [Network Reference](reference/NETWORK.md) |
| Tilemaps | [Tilemap Reference](reference/TILEMAP.md) |
| Particles | [Particles Reference](reference/PARTICLES.md) |
