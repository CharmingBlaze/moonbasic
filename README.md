# moonBASIC

**A modern game programming language for people who want to make games, not fight a compiler.**

moonBASIC is a simple, fast, and fun language for making 2D and 3D games. It is inspired by the classic simplicity of **BlitzBasic** and **DarkBASIC**, but built for the modern era with a powerful set of integrated libraries.

Our philosophy is to remove as much friction as possible between an idea and a playable result. This means no complex project setup, no header files, no build system configuration, and a focus on a clean, imperative API. Just write your code and run it.

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
rot  = Mat4.Identity()
angle# = 0.0

WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    angle# = angle# + 45 * Time.Delta()
    Mat4.SetRotation(rot, 0, angle#, 0)
    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cube, mat, rot)
    cam.End()
    Render.Frame()
WEND

Mesh.Free(cube)
Material.Free(mat)
Mat4.Free(rot)
Camera.Free(cam)
Window.Close()
```

---

## Features

-   **Simple Syntax**: A clean, readable syntax inspired by classic BASIC, with modern features like user-defined functions and a dot-notation command system.
-   **2D & 3D Graphics**: A comprehensive set of commands for 2D drawing and 3D model rendering, powered by **raylib**.
-   **Physics**: Integrated 3D physics with **Jolt Physics** and 2D physics with **Box2D**.
-   **Networking**: Built-in support for multiplayer games using **ENet**.
-   **Fast Compilation**: The compiler, written in Go, is extremely fast. Go from code to running program in milliseconds.
-   **Cross-Platform**: Build and run on Windows and Linux.

---

## Architecture

moonBASIC is a compiled, bytecode-interpreted language. The `moonbasic` executable is a self-contained compiler and virtual machine (VM).

1.  **Compile**: When you run `moonbasic mygame.mb`, the compiler first parses your `.mb` source code into an Abstract Syntax Tree (AST).
2.  **Codegen**: It then walks the AST to generate a compact bytecode program.
3.  **Execute**: The bytecode is immediately executed by the built-in VM. The VM is written in Go and communicates with the underlying C-based libraries (like raylib and Jolt) via CGo.

This architecture provides a balance of high-level simplicity for the user and high performance for the game.

---

## Technology Stack

moonBASIC stands on the shoulders of giants. It integrates several industry-standard, open-source libraries.

| Library | Version | Purpose |
|---|---|---|
| [Go](https://go.dev) | 1.25+ (see `go.mod`) | The core language for the compiler, VM, and runtime.
| [raylib](https://raylib.com) | 5.5 | The foundation for windowing, rendering, input, and audio.
| [Jolt Physics](https://github.com/jrouwe/JoltPhysics) | 5.1 | High-performance 3D physics engine.
| [Box2D](https://box2d.org) | 3.0 | Robust 2D physics engine.
| [ENet](http://enet.bespin.org) | 1.3 | Reliable UDP networking library.

---

## Getting Started

### 1. Install Dependencies

-   **Go 1.25+** (see `go.mod`)
-   **A C Compiler**: MinGW-w64 on Windows, or GCC on Linux.

See the [Building from Source](docs/BUILDING.md) guide for detailed setup instructions.

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

From the repo root, run a bundled demo (requires **CGO**):

```bash
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

Or build the driver and pass any `.mb` file:

```bash
./moonbasic mygame.mb
```

More samples: [examples/README.md](examples/README.md).

---

## Documentation

| Document | What it covers |
|---|---|
| [Getting Started](docs/GETTING_STARTED.md) | Installation, your first window, and key concepts. |
| [Programming Guide](docs/PROGRAMMING.md) | Game loop, command style, 2D/3D, CGO and platforms. |
| [Language Reference](docs/LANGUAGE.md) | Variables, types, loops, and functions. |
| [Command Index](docs/COMMANDS.md) | Every command in one place — built-ins and module commands. |
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
| [Font](docs/reference/FONT.md) | Load custom TTF/OTF fonts. |
| [Sprite & Atlas](docs/reference/SPRITE.md) | Aseprite animation, texture atlases. |
| [Model, Mesh & Material](docs/reference/MODEL.md) | 3D models, procedural meshes, PBR materials. |
| [Mat4](docs/reference/MAT4.md) | 4×4 matrix math for transforms. |
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

Contributions are welcome! This project is in active development. The best way to get started is to check the `COMMAND_AUDIT.txt` file, which lists all commands and their implementation status (`DONE`, `PARTIAL`, or `MISSING`). Implementing a `PARTIAL` or `MISSING` command is a great first contribution.

---

## License

moonBASIC is licensed under the MIT License. See [LICENSE](LICENSE) for details.
