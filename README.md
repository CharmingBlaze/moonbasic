# moonBASIC

**The most powerful modern BASIC compiler for building real games.**

moonBASIC is not a toy interpreter or a nostalgia dialect frozen in the 1980s. It is a **full compiler**—lexer, parser, AST, and bytecode codegen—paired with a **production-grade game runtime**. You write readable BASIC-style source; the toolchain turns it into compact bytecode and runs it on a fast VM wired straight into **raylib**, **Jolt**, **Box2D**, and **ENet**. That combination—classic syntax, modern compilation, and the same libraries serious indie games use—is the bar we set for a **modern BASIC compiler**: maximum reach from a single language, without leaving the BASIC idiom.

If you grew up on **Blitz BASIC** or **DarkBASIC** and wished that spirit existed with today’s graphics, physics, and networking, this is that idea rebuilt for the current era.

---

## Why call it the most powerful modern BASIC?

- **Real compiler, not a line-by-line runner.** Source is parsed into an AST, lowered to bytecode, and executed by a dedicated VM—predictable performance and room to grow.
- **Scope of the standard library.** One language spans **2D and 3D rendering**, **immediate-mode GUI**, **sprites and atlases**, **shaders and lights**, **tilemaps**, **particles**, **2D and 3D physics**, **character controllers**, **audio**, and **multiplayer networking**—the kind of surface area usually split across many languages and glue code.
- **Industrial backing where it matters.** The runtime sits on **raylib** (window, GPU, input, audio), **Jolt** and **Box2D** for physics, and **ENet** for UDP game traffic—not stripped-down teaching stubs.
- **Fast iteration.** The driver is written in **Go**; compile cycles are short so you stay in flow.

“Most powerful” here means **breadth + seriousness of the stack + a proper compile step**, while keeping the **low-friction, imperative BASIC feel**: no project boilerplate, no headers, no CMake ritual—just your `.mb` file and the `moonbasic` binary.

---

## What does it look like?

Here is a compact 3D example (the [examples/spin_cube](examples/spin_cube/main.mb) demo adds a grid, `Draw.Text`, and proper `Free` calls):

```basic
IF NOT Window.Open(960, 540, "Spinning cube") THEN END
ENDIF
Window.SetFPS(60)

cam = Camera.Make()
cam.SetPos(0, 2, 8)
cam.SetTarget(0, 0, 0)
cam.SetFOV(45)

cube = Mesh.MakeCube(2, 2, 2)
mat  = Material.MakeDefault()
cubeXform  = Transform.Identity()
angle# = 0.0

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    angle# = angle# + 45 * Time.Delta()
    Transform.SetRotation(cubeXform, 0, angle#, 0)
    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cube, mat, cubeXform)
    cam.End()
    Render.Frame()
WEND

Mesh.Free(cube)
Material.Free(mat)
Transform.Free(cubeXform)
Camera.Free(cam)
Window.Close()
```

---

## Features

- **Modern BASIC syntax** — Familiar control flow and types, plus user-defined functions and a clear `Module.Command` API.
- **2D & 3D graphics** — Drawing, text, textures, cameras, models, materials, and more, powered by **raylib**.
- **Physics** — **Jolt** for 3D, **Box2D** for 2D.
- **Networking** — **ENet** for client/server-style multiplayer over UDP.
- **Fast compilation** — Go-hosted toolchain; edit, compile, and run in milliseconds.
- **Cross-platform** — Windows and Linux (with CGO and a C toolchain).

---

## Architecture

moonBASIC is a **compiled, bytecode-driven** language. The `moonbasic` executable bundles the compiler and the VM.

1. **Compile** — `moonbasic mygame.mb` parses `.mb` source into an AST.
2. **Codegen** — The AST becomes a compact bytecode program.
3. **Execute** — The VM runs that bytecode and calls into native libraries (**raylib**, **Jolt**, etc.) via **cgo**.

You get BASIC-level simplicity at the keyboard and a structured pipeline under the hood—not a single giant interpreter loop.

---

## Technology Stack

| Library | Version | Purpose |
|---|---|---|
| [Go](https://go.dev) | 1.25+ (see `go.mod`) | Compiler, VM, and runtime host. |
| [raylib](https://raylib.com) | via raylib-go | Windowing, rendering, input, audio. |
| [Jolt Physics](https://github.com/jrouwe/JoltPhysics) | via jolt-go | 3D rigid-body simulation. |
| [Box2D](https://box2d.org) | 3.x | 2D physics. |
| [ENet](http://enet.bespin.org) | 1.3 | Reliable UDP networking. |

---

## Getting Started

### 1. Install Dependencies

- **Go 1.25+** (see `go.mod`)
- **A C compiler**: MinGW-w64 on Windows, or GCC on Linux.

See [Building from Source](docs/BUILDING.md) for detailed setup.

### 2. Build moonBASIC

```bash
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic

# On Windows (in a regular command prompt)
set CGO_ENABLED=1
set CC=C:\msys64\mingw64\bin\gcc.exe
go build -o moonbasic.exe .

# On Linux
CGO_ENABLED=1 go build -o moonbasic .
```

### 3. Run Your First Program

From the repo root (requires **CGO**):

```bash
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

Or use a built binary:

```bash
./moonbasic mygame.mb
```

More samples: [examples/README.md](examples/README.md).

---

## Documentation

| Document | What it covers |
|---|---|
| [Getting Started](docs/GETTING_STARTED.md) | Installation, first window, core ideas. |
| [Programming Guide](docs/PROGRAMMING.md) | Game loop, command style, 2D/3D, CGO and platforms. |
| [Language Reference](docs/LANGUAGE.md) | Variables, types, loops, functions. |
| [Command Index](docs/COMMANDS.md) | Built-ins and module commands in one place. |
| [Examples](docs/EXAMPLES.md) | Narrated snippets (see also `examples/`). |
| [Examples (repo)](examples/README.md) | Runnable `main.mb` programs and how to launch them. |

### API Reference (by module)

| Module | Description |
|---|---|
| [Window](docs/reference/WINDOW.md) | Open/close window, monitors, flags. |
| [Render](docs/reference/RENDER.md) | Clear, frame, blend modes, screenshots. |
| [Camera](docs/reference/CAMERA.md) | 3D camera and 2D scrolling camera. |
| [Draw 2D](docs/reference/DRAW2D.md) | Rectangles, circles, text, textures. |
| [GUI (raygui)](docs/reference/GUI.md) | Immediate-mode `GUI.*` widgets. |
| [Texture](docs/reference/TEXTURE.md) | Load, free, and generate textures. |
| [Image (CPU)](docs/reference/IMAGE.md) | Pixel buffers, software draw, `Texture.FromImage`. |
| [Font](docs/reference/FONT.md) | Load custom TTF/OTF fonts. |
| [Sprite & Atlas](docs/reference/SPRITE.md) | Texture strips, **`Anim.*`** FSM, atlases ([ATLAS.md](docs/reference/ATLAS.md)). |
| [Model, Mesh & Material](docs/reference/MODEL.md) | 3D models, procedural meshes, PBR materials. |
| [Transform](docs/reference/TRANSFORM.md) | 4×4 object/world transform matrices (legacy `Mat4.*` in [MAT4.md](docs/reference/MAT4.md)). |
| [Shader](docs/reference/SHADER.md) | Load and apply GLSL shaders. |
| [Light](docs/reference/LIGHT.md) | Directional, point, and spot lights. |
| [Physics 3D (Jolt)](docs/reference/PHYSICS3D.md) | Rigid bodies, shapes, forces. |
| [Physics 2D (Box2D)](docs/reference/PHYSICS2D.md) | 2D bodies and colliders. |
| [Character Controller](docs/reference/CHARCONTROLLER.md) | Kinematic player controller. |
| [Tilemap (Tiled)](docs/reference/TILEMAP.md) | Load .tmx maps, collision, layers. |
| [Particles](docs/reference/PARTICLES.md) | Emitter-based particle effects. |
| [Input](docs/reference/INPUT.md) | Keyboard, mouse, gamepad, action maps. |
| [Audio](docs/reference/AUDIO.md) | Sounds, music streaming, audio streams. |
| [Network (ENet)](docs/reference/NETWORK.md) | Multiplayer server/client over UDP. |
| [Math](docs/reference/MATH.md) | Trig, interpolation, randomization. |
| [String](docs/reference/STRING.md) | String manipulation and parsing. |
| [Array](docs/reference/ARRAY.md) | DIM arrays and array operations. |
| [File I/O](docs/reference/FILE.md) | Read/write files and directories. |
| [Time](docs/reference/TIME.md) | Delta time, wall-clock time. |
| [Debug](docs/reference/DEBUG.md) | ASSERT and debug utilities. |
| [Bitwise](docs/reference/BITWISE.md) | Bit flags and manipulation. |
| [Console](docs/reference/CONSOLE.md) | PRINT, INPUT, CLS. |
| [System](docs/reference/SYSTEM.md) | Environment, clipboard, command-line args. |

---

## Contributing

Contributions are welcome. This project is under active development. Start with `COMMAND_AUDIT.txt` for command implementation status (`DONE`, `PARTIAL`, `MISSING`). Closing gaps there is a great first contribution.

---

## License

moonBASIC is licensed under the MIT License. See [LICENSE](LICENSE) for details.
