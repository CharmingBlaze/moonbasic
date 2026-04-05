# moonBASIC

**The most powerful BASIC compiler ever built for real games—and the most fun way to ship one.**

moonBASIC is a **full compiler**: lexer, parser, semantic analyzer, AST, and bytecode codegen, plus a **fast VM** that talks directly to the same native libraries the industry uses for shipping titles. You write imperative BASIC with namespaces like `Window.Open`, `Camera.Begin`, `Physics3D.Step`, and `Net.Update`; the toolchain compiles your `.mb` sources to compact bytecode and runs them with **Go**-hosted tooling and **CGO** bridges into **raylib**, **Jolt**, **Box2D**, and **ENet**.

It is **not** a toy line-by-line interpreter. It is **not** a game engine you wrestle with for days before drawing a triangle. It is a **language and runtime** designed so you **program the game**, not the engine—and so you move fast without giving up serious graphics, physics, or networking.

---

## A word from the lab

*Legend has it moonBASIC was forged by **one hundred computer scientists** sealed inside a **secret underground moon base**, where the only entertainment was proving that BASIC could drive a modern GPU and a physics engine without apologizing. Whether or not you believe the base exists, the outcome is real: a compiler and VM engineered with **memory safety and lifecycle discipline** in mind—handles, generations, idempotent `Free`, `Heap.FreeAll` on shutdown—so the fun stays in your game loop, not in chasing leaks.*

*That engineering is also why moonBASIC stays **dynamically typed** at the language level: the compiler is built to work **with** BASIC’s suffix conventions (`name$`, `count#`, `flag?`) and runtime values, and to **catch mistakes** through semantic analysis and manifest-driven builtins—not by forcing you into a second job as a type bureaucrat. You think in game terms; the toolchain keeps the contract between your script and the native stack honest.*

---

## Why program in moonBASIC instead of “using a game engine”?

| You get | What it means |
|--------|----------------|
| **A real compile step** | Source → AST → bytecode → VM. Predictable structure, room for optimization, and a clear story when you debug. |
| **Less ceremony** | No scene-graph religion, no asset pipeline tax for a first prototype—open a window, clear the screen, draw, step physics, send a packet. |
| **Industry-grade natives** | **raylib** for windowing, drawing, input, and audio. **Jolt** (3D) and **Box2D** (2D) for physics families trusted in production. **ENet** for UDP-style multiplayer—the same *kind* of stack studios reach for when they need performance without writing everything from scratch. |
| **A huge builtin surface** | Hundreds of dotted commands: rendering, GUI, sprites, tilemaps, particles, lights, shaders, navigation helpers, terrain/world streaming, data (JSON, CSV, tables, SQLite when enabled), noise, networking, and more—see [docs/COMMANDS.md](docs/COMMANDS.md). The checked-in **`MASTER_AUDIT.txt`** compares the manifest to registered natives (see [ARCHITECTURE.md](ARCHITECTURE.md)). |
| **Rewarding workflow** | You author **your** loop, **your** systems, **your** feel—closer to classic **Blitz** / **DarkBASIC** energy than to configuring someone else’s editor for a week. |

“Faster” here means **faster to a playable build** for many projects: short compile cycles (Go-hosted toolchain), one binary, one language from menu code to netcode—without abandoning the **readable, imperative BASIC** style.

---

## What does code look like?

From [examples/spin_cube](examples/spin_cube/main.mb)—a real window, a camera, a mesh, a frame loop:

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

That is **compiled** moonBASIC—not interpreted line noise. Same language scales to terrain, props, physics, and multiplayer as you grow.

---

## How the toolchain fits together

1. **Edit** `.mb` source (`.mbc` MOON bytecode is optional for shipping or tooling).
2. **Compile** — parse, analyze, generate IR v2 bytecode ([ARCHITECTURE.md](ARCHITECTURE.md)).
3. **Run** — VM executes opcodes and dispatches **natives** registered from `compiler/builtinmanifest/commands.json` into modular `runtime/*` packages.
4. **Native layer** — Raylib on the main OS thread where required; physics and network code behind the same manifest keys with CGO/stub splits where platforms differ.

Memory and resource rules are documented for contributors in [docs/MEMORY.md](docs/MEMORY.md): handle tables, idempotent teardown, and how shutdown ties to `Window.Close` and registry shutdown.

---

## Technology stack (what actually ships)

| Piece | Role |
|-------|------|
| **[Go](https://go.dev)** | Compiler, VM, runtime orchestration (`go.mod` pins 1.25+). |
| **[raylib](https://www.raylib.com)** (via [raylib-go](https://github.com/gen2brain/raylib-go)) | Window, rendering, input, audio, images, GPU textures. |
| **[Jolt Physics](https://github.com/jrouwe/JoltPhysics)** (via [jolt-go](https://github.com/bbitechnologies/jolt-go)) | 3D rigid-body style simulation where the binding is available (see [ARCHITECTURE.md](ARCHITECTURE.md) §12 for platform notes). |
| **[Box2D](https://box2d.org)** (via [ByteArena/box2d](https://github.com/ByteArena/box2d)) | 2D physics API surface in the manifest. |
| **[ENet](http://enet.bespin.org)** (via [go-enet](https://github.com/codecat/go-enet)) | Reliable UDP networking for client/server style games. |

These are the **same families** of libraries you will find behind many shipped indie and AA games—wired here so your BASIC calls land on real simulation and real sockets, not toy stubs (stubs exist only where a platform cannot link the native library yet).

---

## Build & run

**Requirements:** Go **1.25+**, a **C toolchain** (MinGW-w64 on Windows, GCC on Linux), **CGO enabled** for full graphics/physics/net.

```bash
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic

# Windows (example: adjust CC to your MinGW gcc)
set CGO_ENABLED=1
set CC=C:\msys64\mingw64\bin\gcc.exe
go build -o moonbasic.exe .

# Linux
CGO_ENABLED=1 go build -o moonbasic .
```

Run a sample from the repo root:

```bash
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

Deeper setup: [docs/BUILDING.md](docs/BUILDING.md). More runnable projects: [examples/README.md](examples/README.md).

---

## Documentation map

| Start here | Why |
|------------|-----|
| [Getting Started](docs/GETTING_STARTED.md) | Install, first window, mental model. |
| [Programming Guide](docs/PROGRAMMING.md) | Game loop, commands, 2D/3D, platforms. |
| [Language Reference](docs/LANGUAGE.md) | Variables, control flow, functions, suffix types. |
| [Command Index](docs/COMMANDS.md) | Built-ins in one place. |
| [Examples](docs/EXAMPLES.md) | Narrated snippets. |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Ground truth for pipeline, registry, phases. |
| [MEMORY.md](docs/MEMORY.md) | Handles, GC vs native, `FreeAll`. |

### API reference by module

| Module | Topics |
|--------|--------|
| [Window](docs/reference/WINDOW.md) | Open/close, monitors, flags |
| [Render](docs/reference/RENDER.md) | Clear, frame, blend, screenshots |
| [Camera](docs/reference/CAMERA.md) | 3D and 2D cameras |
| [Draw 2D](docs/reference/DRAW2D.md) | Primitives, text, textures |
| [GUI](docs/reference/GUI.md) | Immediate-mode `GUI.*` (raygui) |
| [Texture](docs/reference/TEXTURE.md) | Load, upload, generate |
| [Image](docs/reference/IMAGE.md) | CPU images, `Texture.FromImage` |
| [Font](docs/reference/FONT.md) | TTF/OTF |
| [Sprite & Atlas](docs/reference/SPRITE.md) | Animation, atlases |
| [Model / Mesh / Material](docs/reference/MODEL.md) | 3D assets, procedural meshes |
| [Transform](docs/reference/TRANSFORM.md) | Matrices (`Mat4.*` legacy in [MAT4.md](docs/reference/MAT4.md)) |
| [Shader](docs/reference/SHADER.md) | GLSL |
| [Light](docs/reference/LIGHT.md) | Scene lights |
| [Physics 3D (Jolt)](docs/reference/PHYSICS3D.md) | Rigid bodies, shapes |
| [Physics 2D (Box2D)](docs/reference/PHYSICS2D.md) | 2D colliders |
| [Character Controller](docs/reference/CHARCONTROLLER.md) | Kinematic player |
| [Tilemap](docs/reference/TILEMAP.md) | Tiled `.tmx` |
| [Particles](docs/reference/PARTICLES.md) | Emitters |
| [Input](docs/reference/INPUT.md) | Keyboard, mouse, gamepad |
| [Audio](docs/reference/AUDIO.md) | Sounds, music |
| [Network (ENet)](docs/reference/NETWORK.md) | Multiplayer |
| [Math](docs/reference/MATH.md) | Math helpers |
| [String](docs/reference/STRING.md) | Text |
| [Array](docs/reference/ARRAY.md) | `DIM`, array ops |
| [File](docs/reference/FILE.md) | Filesystem |
| [Time](docs/reference/TIME.md) | Delta time |
| [Debug](docs/reference/DEBUG.md) | Debug helpers |
| [Bitwise](docs/reference/BITWISE.md) | Flags |
| [Console](docs/reference/CONSOLE.md) | `PRINT`, `INPUT` |
| [System](docs/reference/SYSTEM.md) | Env, argv, clipboard |

More topics (terrain, world, weather, JSON, CSV, noise, …) live under [docs/reference/](docs/reference/) and [docs/COMMANDS.md](docs/COMMANDS.md).

---

## Contributing

Issues and PRs welcome. For implementation status of builtins, see `COMMAND_AUDIT.txt` (and related audit files). Closing gaps there is a great first contribution.

---

## License

MIT — see [LICENSE](LICENSE).
