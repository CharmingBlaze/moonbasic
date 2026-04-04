# Getting Started with moonBASIC

## What you need

- Windows x64 or Linux x64
- Go **1.25.3** or later (see `go.mod` in the repo) — https://go.dev/dl/
- A C compiler:
  - Windows: MinGW-w64 — https://www.mingw-w64.org/
  - Linux: GCC — usually already installed

---

## Install

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

Run it:
```bash
moonbasic hello.mb
```

Output:
```
Hello, moonBASIC!
```

---

## Your First Window

Create `window.mb`:

```basic
IF NOT Window.Open(960, 540, "Hello Window") THEN
    PRINT("Could not open window")
    END
ENDIF
Window.SetFPS(60)

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    Render.Clear(30, 40, 60)
    ; Default bitmap font — no .ttf file required
    Draw.Text("Press ESC or close the window", 260, 260, 20, 255, 255, 255, 255)
    Render.Frame()
WEND

Window.Close()
```

Run it:
```bash
moonbasic window.mb
```

You should see a dark blue window with white text. Press the window's close button to exit.

---

## Your First Game Loop

moonBASIC games follow a simple pattern:

```basic
; 1. Setup — runs once
IF NOT Window.Open(800, 600, "My Game") THEN
    PRINT("Could not open window")
    END
ENDIF
Window.SetFPS(60)

player_x# = 400
player_y# = 300

; 2. Loop — runs every frame until the window is closed
WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    dt# = Time.Delta()         ; seconds since last frame

    ; Update
    IF Input.KeyDown(KEY_RIGHT) THEN player_x# = player_x# + 200 * dt#
    IF Input.KeyDown(KEY_LEFT)  THEN player_x# = player_x# - 200 * dt#
    IF Input.KeyDown(KEY_DOWN)  THEN player_y# = player_y# + 200 * dt#
    IF Input.KeyDown(KEY_UP)    THEN player_y# = player_y# - 200 * dt#

    ; Draw
    Render.Clear(20, 30, 40)
    Draw.Rectangle(INT(player_x#) - 16, INT(player_y#) - 16, 32, 32, 100, 200, 255, 255)
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
IF NOT Window.Open(960, 540, "3D Cube") THEN
    PRINT("Could not open window")
    END
ENDIF
Window.SetFPS(60)

cam = Camera.Make()
cam.SetPos(0, 3, 8)
cam.SetTarget(0, 0, 0)
cam.SetFOV(45)

cube = Mesh.MakeCube(2, 2, 2)
mat  = Material.MakeDefault()
xform = Mat4.Identity()
angle# = 0.0

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    angle# = angle# + 1.5 * Time.Delta()
    Mat4.SetRotation(xform, angle#, angle# * 0.7, 0)

    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cube, mat, xform)
        Draw.Grid(20, 1.0)
    cam.End()
    Render.Frame()
WEND

Mesh.Free(cube)
Material.Free(mat)
Mat4.Free(xform)
Camera.Free(cam)
Window.Close()
```

---

## Variable Types

Variables are typed by their **name suffix**:

| Suffix | Type | Example |
|---|---|---|
| `$` | String | `name$ = "Player"` |
| `#` | Float | `speed# = 5.5` |
| `?` | Boolean | `alive? = TRUE` |
| (none) | Integer | `score = 100` |

There is no separate declaration step — variables are created on first assignment.

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
| Every command in one place | [Command Index](COMMANDS.md) |
| Copy-paste snippets + narrative | [Examples](EXAMPLES.md) |
| Runnable programs under `examples/` | [examples/README.md](../examples/README.md) |
| 2D drawing | [Draw Reference](reference/DRAW2D.md) |
| 3D models & meshes | [Model Reference](reference/MODEL.md) |
| 2D physics (Box2D) | [Physics 2D](reference/PHYSICS2D.md) |
| 3D physics (Jolt) | [Physics 3D](reference/PHYSICS3D.md) |
| Multiplayer (ENet) | [Network Reference](reference/NETWORK.md) |
| Tilemaps | [Tilemap Reference](reference/TILEMAP.md) |
| Particles | [Particles Reference](reference/PARTICLES.md) |
