# moonBASIC

**A modern BASIC for developers who want to build without unnecessary friction**

Many engines impose their own complexity before you can begin creating. Navigating Godot’s node hierarchy can feel cumbersome when your goal is simply to make a game. Unity has grown increasingly bloated and demands a significant amount of setup before you can begin meaningful work. Unreal often requires hours of preparation before anything appears on screen.

**moonBASIC** takes a different approach. You open a window, create a camera, and immediately step into a hardware-accelerated, Jolt-powered three-dimensional environment. The intention is to let you begin building while your coffee is still warm.

moonBASIC is a modern BASIC language designed for developers who appreciate the clarity of earlier eras but expect contemporary rendering, physics, and networking. It is not an interpreter. It is a real compiler that produces bytecode executed by a production-grade runtime built on the same technologies used in professional game development today.

---

## Modern capability with a direct and purposeful design philosophy

moonBASIC is not an attempt to recreate the past. It is a tool for shipping real projects with a workflow that values clarity and speed.

- **Vertical integration** replaces scattered glue code. You work within one toolchain and one mental model, and your game runs in the same environment you develop in.
- **Two-dimensional and three-dimensional development** share the same workflow. You can create a PlayStation One–style horror prototype or a Super Nintendo–style platformer without switching paradigms.
- **The aesthetic direction** favors a clean, cyber-retro sensibility reminiscent of late-nineties software. The interface is intentionally straightforward, leaving room for your own artistic direction rather than imposing a default look.

The compiler is stable, the two-dimensional and three-dimensional pipelines are active, and the standard library already includes support for Tiled maps, materials, sprites, atlases, particles, audio, lighting, shaders, and more. You can explore these systems through the [documentation index](#documentation) and [ARCHITECTURE.md](ARCHITECTURE.md).

---

## A development environment focused on shipping real projects

| You have experienced | moonBASIC provides |
|----------------------|-------------------|
| Complex editors that slow iteration | A single executable, a single `.mb` file, and rapid feedback |
| Scripting languages that limit performance | A compiled pipeline with bytecode and a virtual machine |
| Engines that force one physics model | **Jolt** for three-dimensional physics and **Box2D** for two-dimensional precision |
| Multiplayer systems that feel bolted on | Networking that integrates naturally with the language |

moonBASIC prioritizes developer speed and conceptual clarity. It is intentionally opinionated and designed to reduce friction.

---

## Technical foundation

These components form the core of the runtime and are chosen for their reliability and performance.

| Component | Purpose |
|-----------|---------|
| **Go** | Fast compilation, simple cross-platform builds, and a predictable runtime host |
| **Raylib** (via **raylib-go**) | Efficient hardware-accelerated two-dimensional and three-dimensional rendering, input, and audio |
| **Jolt Physics** | High-performance three-dimensional rigid-body simulation |
| **Box2D** | Precise and proven two-dimensional physics |
| **ENet** | Reliable UDP networking suitable for real-time games |

moonBASIC uses the right physics engine for the right dimensionality. Three-dimensional games benefit from Jolt, while two-dimensional games rely on Box2D. This avoids the compromises common in engines that force a single physics solution for every project.

---

## Multiplayer through the strengths of Go

Multiplayer functionality is not an afterthought. Because the runtime is hosted in Go, networking can take advantage of goroutines and structured concurrency. The engine handles synchronization, buffering, and timing so developers can think in straightforward multiplayer operations.

Start with **[docs/reference/NETWORK.md](docs/reference/NETWORK.md)** and the **`NETWORK.*`** commands in **[docs/COMMANDS.md](docs/COMMANDS.md)**.

---

## Example

Below is a compact three-dimensional example with **no** `#`, `$`, or `?` type suffixes—implicit typing only.

```moonbasic
WINDOW.OPEN(960, 540, "Spinning cube")
WINDOW.SETFPS(60)

cam = CAMERA.MAKE()
CAMERA.SETPOSITION(cam, 0, 2, 8)
CAMERA.SETTARGET(cam, 0, 0, 0)
CAMERA.SETFOV(cam, 45)

cube = MESH.MAKECUBE(2, 2, 2)
mat  = MATERIAL.MAKEDEFAULT()
cubeXform = TRANSFORM.IDENTITY()
angle = 0

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    dt = TIME.DELTA()
    angle = angle + 45 * dt
    TRANSFORM.SETROTATION(cubeXform, 0, angle, 0)
    RENDER.CLEAR(12, 14, 22)
    CAMERA.BEGIN(cam)
        MESH.DRAW(cube, mat, cubeXform)
    CAMERA.END(cam)
    RENDER.FRAME()
WEND

MESH.FREE(cube)
MATERIAL.FREE(mat)
TRANSFORM.FREE(cubeXform)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

The full [examples/spin_cube](examples/spin_cube/main.mb) demo adds a grid and on-screen text.

---

## Architecture in brief

1. **Compilation** produces bytecode with the **`.mbc`** extension.
2. **Execution** is handled by the virtual machine, which interfaces with Raylib, Jolt, Box2D, ENet, and other systems through **CGO** where enabled.

This provides the clarity of BASIC at the language level and the performance of a real engine beneath it.

---

## Getting started

### Prebuilt binaries

**Windows x64** and **Linux x64** builds are available on **[Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)**. Each archive includes the compiler, the runtime, and a release guide (`README-RELEASE.txt`). After extraction, verify with `moonbasic --version` or run any example under [`examples/`](examples/).

### Build from source

Building from source requires **Go** and a **C toolchain**. Full graphical programs need the **`fullruntime`** build tag and **`moonrun`** (or **`moonbasic --run`** from a full-runtime build). See **[AGENTS.md](AGENTS.md)**, **[CONTRIBUTING.md](CONTRIBUTING.md)**, and **[docs/BUILDING.md](docs/BUILDING.md)**.

```bash
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic
# Windows (example): set CGO_ENABLED=1 and a working gcc, then:
go build -o moonbasic.exe .

# Run a 3D sample (full runtime + CGO):
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

---

## Documentation

The repository includes extensive documentation covering installation, language features, engine commands, examples, and internal architecture. Reference modules include Window, Render, Camera, Draw2D, GUI, Texture, Image, Sprite, Atlas, Model, Mesh, Material, Shader, Light, Physics3D, Physics2D, Character Controller, Tilemap, Network, Input, and Audio.

| Document | What it covers |
|----------|----------------|
| [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) | Install, first window, mental model |
| [docs/PROGRAMMING.md](docs/PROGRAMMING.md) | Game loop, modules, 2D/3D |
| [docs/LANGUAGE.md](docs/LANGUAGE.md) | Variables, control flow, functions |
| [docs/COMMANDS.md](docs/COMMANDS.md) | Built-in command index |
| [examples/README.md](examples/README.md) | Runnable sample programs |

More detail lives in **[docs/reference/](docs/reference/)** and **[ARCHITECTURE.md](ARCHITECTURE.md)**.

---

## Contributing

Contribution guidelines and development notes are in **[CONTRIBUTING.md](CONTRIBUTING.md)** and **[docs/DEVELOPER.md](docs/DEVELOPER.md)**. Continuous integration validates builds, tests, and representative `go run . --check` samples.

On **Windows**, a **`fullruntime`** link that pulls in Jolt requires prebuilt **`libJolt.a`** and **`libjolt_wrapper.a`** in **[third_party/jolt-go/jolt/lib/windows_amd64/](third_party/jolt-go/jolt/lib/windows_amd64/README.md)** (or build them with **`third_party/jolt-go/scripts/build-libs-windows.ps1`**). **`scripts/check-jolt-windows-libs.ps1`** checks that both files are present.

---

## License

**MIT** — see [LICENSE](LICENSE).
