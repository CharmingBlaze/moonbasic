# Getting Started with moonBASIC

**API style:** Prefer **`NAMESPACE.ACTION`** (`WINDOW.OPEN`, `TIME.DELTA`, `INPUT.KEYDOWN`, …) in new programs — see [STYLE_GUIDE.md](../STYLE_GUIDE.md) and [DOC_STYLE_GUIDE.md](DOC_STYLE_GUIDE.md). Some sections below still show **Easy Mode** (`Window.Open`, `CreateCamera`, …) for readability; they compile to the same registry keys.

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
WINDOW.OPEN(960, 540, "Hello Window")
WINDOW.SETFPS(60)

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    RENDER.CLEAR(30, 40, 60)
    ; Default bitmap font — no .ttf file required
    DRAW.TEXT("Press ESC or close the window", 260, 260, 20, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

Run it with the **game runtime** (plain `moonbasic` only compiles to `.mbc`):

```bash
moonrun window.mb
```

If you build from source: `go build -tags fullruntime -o moonrun ./cmd/moonrun`, then `moonrun window.mb`. From the repo without installing binaries: `CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun window.mb`.

You should see a dark blue window with white text. Press the window's close button to exit.

If opening the window fails (invalid size, unavailable display, etc.), the runtime prints a short message to **stderr** and the process exits—you do not need `IF NOT WINDOW.OPEN …` in every program. Use **`WINDOW.CANOPEN`** only when you must choose a fallback resolution or show your own error without opening.

---

## Your First Game Loop

moonBASIC games follow a simple pattern:

```basic
; 1. Setup — runs once
WINDOW.OPEN(800, 600, "My Game")
WINDOW.SETFPS(60)

player_x = 400
player_y = 300

; 2. Loop — runs every frame until the window is closed
WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    dt = TIME.DELTA()         ; seconds since last frame

    ; Update
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN player_x = player_x + 200 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN player_x = player_x - 200 * dt
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN player_y = player_y + 200 * dt
    IF INPUT.KEYDOWN(KEY_UP)    THEN player_y = player_y - 200 * dt

    ; Draw
    RENDER.CLEAR(20, 30, 40)
    DRAW.RECTANGLE(INT(player_x) - 16, INT(player_y) - 16, 32, 32, 100, 200, 255, 255)
    RENDER.FRAME()
WEND

; 3. Cleanup
WINDOW.CLOSE()
```

Key concepts:

| Concept | Explanation |
|---|---|
| **`WINDOW.SHOULDCLOSE()`** | Returns `TRUE` when the player clicks the X or presses Alt+F4. |
| **`INPUT.KEYDOWN(KEY_ESCAPE)`** | Common way to exit demos with the keyboard. |
| **`TIME.DELTA()`** | Seconds since last frame. Multiply speeds by this for frame-rate-independent movement. |
| **`RENDER.CLEAR(r, g, b)`** | Clears the screen to a color. Always call this first in the loop. |
| **`RENDER.FRAME()`** | Shows what was drawn. Always call this last in the loop. |
| **`DRAW.RECTANGLE(x, y, w, h, r, g, b, a)`** | Draws a filled colored rectangle. |
| **`INPUT.KEYDOWN(KEY_*)`** | Returns `TRUE` while a key is held. |

---

## Adding a 3D Camera

For 3D scenes, use **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`** (or **`CAMERA.BEGIN`/`CAMERA.END`**). The snippet below uses **Easy Mode** handle helpers (`cam.Begin`, `Mesh.Draw`, …) as in many examples; registry equivalents are [CAMERA](reference/CAMERA.md), [DRAW3D](reference/DRAW3D.md).

```basic
WINDOW.OPEN(960, 540, "3D Cube")
WINDOW.SETFPS(60)

cam = CreateCamera()
cam.SetPos(0, 3, 8)
cam.SetTarget(0, 0, 0)
cam.SetFOV(45)

cube = Mesh.MakeCube(2, 2, 2)
mat  = Material.MakeDefault()
xform = Transform.Identity()
angle = 0.0

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    angle = angle + 1.5 * TIME.DELTA()
    Transform.SetRotation(xform, angle, angle * 0.7, 0)

    RENDER.CLEAR(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cube, mat, xform)
        Draw3D.Grid(20, 1.0)
    cam.End()
    RENDER.FRAME()
WEND

Mesh.Free(cube)
Material.Free(mat)
Transform.Free(xform)
Camera.Free(cam)
WINDOW.CLOSE()
```

---

## Modern Blitz-style 3D (CGO)

This pattern uses **global aliases** (`Graphics3D`, `CreateCamera`, `LoadMesh`, **`RENDER.BEGIN3D`**, `UpdatePhysics`, …) on top of the same engine as **`WINDOW.*`** / **`ENTITY.*`**.

- **Resolution is your choice** — **`WINDOW.OPEN(w, h, …)`** accepts any reasonable size (720p, 1080p, 1440p, 4K, …). You must **`WINDOW.OPEN`** before **`Graphics3D`**: the latter **only resizes** the client area (omit **`Graphics3D`** if the open size is already what you want).
- **`SetMSAA`**: best-effort; for some drivers MSAA is fixed at **`WINDOW.OPEN`** — see [WINDOW](reference/WINDOW.md) / [BUILDING](BUILDING.md).
- **`SetSSAO` / PBR / lights**: need **CGO Raylib**; stubs return errors on non-graphical builds.
- **`UpdatePhysics`** (same as **`UPDATEPHYSICS`**) runs **`ENTITY.UPDATE(TIME.DELTA())`** and best-effort **`WORLD.UPDATE`**, **`PHYSICS2D.STEP`**, **`PHYSICS3D.STEP`** (inactive worlds no-op or ignored).
- **Frame contract**: **`UpdatePhysics`** → **`RENDER.CLEAR`** → 3D pass → **`RENDER.FRAME`**.

The sample below keeps **Blitz / Easy Mode** scene helpers (`Graphics3D`, `CreateCamera`, `LoadMesh`, …); **window and input** use **`WINDOW.*`** / **`INPUT.*`** so it matches the bullets above.

```basic
; Initialize world
WINDOW.OPEN(1920, 1080, "Project: High Fidelity")
WINDOW.SETFPS(60)
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

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    ; Modern Camera Follow logic
    CameraSmoothFollow(cam, car, 0.1)

    ; Physics Impulse based on modern input
    IF INPUT.KEYDOWN(KEY_W) THEN ApplyEntityImpulse(car, 0, 0, 500)

    UpdatePhysics()

    ; The Render Pass
    RENDER.CLEAR(12, 14, 22)
    RENDER.BEGIN3D(cam)
        DrawEntities()   ; PBR, SSAO, dynamic lights (when CGO + assets load)
    RENDER.END3D()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

Place **`assets/`** paths next to your `.mb` or use paths relative to the process working directory. For manual **`PHYSICS3D.START`** + **`BODY3D`**, you can still call **`PHYSICS3D.STEP`** explicitly; **`UpdatePhysics`** already invokes it with **`TIME.DELTA`** when the world is active. See [PHYSICS3D](reference/PHYSICS3D.md) and [EXAMPLES](EXAMPLES.md).

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

Use **`INPUT.KEYDOWN(key)`** for held, **`INPUT.KEYPRESSED(key)`** for first-press only.

### Mouse

`MOUSE_LEFT_BUTTON`, `MOUSE_RIGHT_BUTTON`, `MOUSE_MIDDLE_BUTTON`

Use **`INPUT.MOUSEDOWN(btn)`**, **`INPUT.MOUSEX()`**, **`INPUT.MOUSEY()`**.

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
