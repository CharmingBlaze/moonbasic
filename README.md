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

## Repository layout

The root directory is kept small on purpose: Go sources, `go.mod`, and a handful of top-level guides. Supporting material lives in named folders so the GitHub file tree is easier to scan.

| Path | Purpose |
|------|---------|
| [`cmd/moonbasic`](cmd/moonbasic), [`cmd/moonrun`](cmd/moonrun) | CLI entrypoints (compiler-only vs full runtime). |
| [`compiler/`](compiler/), [`vm/`](vm/) | Language front-end, bytecode, and VM. |
| [`runtime/`](runtime/) | Engine modules (rendering, physics, audio, net, …). |
| [`docs/`](docs/) | Guides and reference; maintainer audits under [`docs/audit/`](docs/audit/). |
| [`testdata/`](testdata/) | `.mb` samples for `--check` and tests (ad hoc copies under [`testdata/dev_samples/`](testdata/dev_samples/)). |
| [`examples/`](examples/) | Runnable projects. |
| [`dist/`](dist/) | What each release flavor contains — see [`dist/README.md`](dist/README.md). |
| [`scripts/`](scripts/), [`tools/`](tools/) | Release packaging and parity / audit helpers. |

---

## Multiplayer through the strengths of Go

Multiplayer functionality is not an afterthought. Because the runtime is hosted in Go, networking can take advantage of goroutines and structured concurrency. The engine handles synchronization, buffering, and timing so developers can think in straightforward multiplayer operations.

Start with **[docs/reference/MULTIPLAYER.md](docs/reference/MULTIPLAYER.md)** (scope + learning path), **[docs/tutorials/FIRST_MULTIPLAYER_GAME.md](docs/tutorials/FIRST_MULTIPLAYER_GAME.md)** (two-process tutorial), then **[docs/reference/NETWORK.md](docs/reference/NETWORK.md)** and the **`NETWORK.*`** commands in **[docs/COMMANDS.md](docs/COMMANDS.md)**.

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

### What to download (prebuilt)

All builds are on **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)** (Windows and Linux **x64**). Each version tag publishes **four** files. Pick one **compiler** archive and/or one **full runtime** archive depending on what you need:

| Your goal | Download (replace `<tag>` with the release, e.g. `v1.2.0`) |
|-----------|-------------------------------------------------------------|
| **Compile** `.mb` → `.mbc`, run **`--check`**, use **`--lsp`** in an editor, no game window | **Compiler only:** `moonbasic-<tag>-compiler-windows-amd64.zip` or `moonbasic-<tag>-compiler-linux-amd64.tar.gz` |
| **Run** games (window, graphics, physics, audio) | **Full runtime:** `moonbasic-<tag>-windows-amd64.zip` or `moonbasic-<tag>-linux-amd64.tar.gz` |

The compiler in **both** flavors uses the **same** builtin command list from the toolchain (`--check` and completions know every command name). Only the **full runtime** download includes **`moonrun`**, which executes those calls on screen.

### What is inside each distribution file

| Archive | After you extract | What it is for |
|---------|-------------------|----------------|
| **Full runtime** (`…-windows-amd64.zip` / `…-linux-amd64.tar.gz`) | **`moonbasic`** / **`moonbasic.exe`**, **`moonrun`** / **`moonrun.exe`**, **`README-RELEASE.txt`** | Full engine: compile, check, **and run** `.mb` / `.mbc` with graphics and physics. May need GPU drivers; on Windows you may need the [VC++ redistributable](https://learn.microsoft.com/en-us/cpp/windows/latest-supported-vc-redist) if a DLL is missing — details are in **`README-RELEASE.txt`** in the zip. |
| **Compiler only** (`…-compiler-…`) | A folder **`MoonBasic-compiler/`** with **`moonbasic`** (or **`moonbasic.exe`**) and a short **`README-COMPILER.txt`** | Toolchain only: **no `moonrun`**, no Raylib next to the compiler. Ideal for authors, CI, and machines where you only compile or lint. |

More detail: **[`dist/README.md`](dist/README.md)**.

### How to use the compiler

1. **Extract** the archive and open a terminal in the folder that contains **`moonbasic`** (on **compiler-only** builds, that is inside **`MoonBasic-compiler/`**).
2. **Verify:** `moonbasic --version` (Windows: `moonbasic.exe --version`).
3. **Check a program** without running it: `moonbasic --check path/to/game.mb` — reports parse/semantic errors.
4. **Compile to bytecode:** `moonbasic path/to/game.mb` — writes **`game.mbc`** next to the source.
5. **Editor support:** run **`moonbasic --lsp`** and point your editor’s MoonBASIC/LSP client at it (stdio).

To **run** a game that opens a window, use a **full runtime** download and run **`moonrun path/to/game.mb`** or **`moonrun path/to/game.mbc`** (same folder as **`moonbasic`** after extract). If you only installed the **compiler-only** bundle, you do not have **`moonrun`** — compile on this machine and run on another that has the full runtime, or download the full archive.

### Running games — no extra compiler to install

**Pre-built `moonrun` from Releases does not call Go, GCC, or Clang.** It embeds the same compile step as `moonbasic`: you can run **`moonrun mygame.mb`** and it compiles in-process, then starts the engine — you do **not** need a separate “moonbasic on PATH” or any programming toolchain. (You still need normal OS pieces: a GPU stack on Linux, and on Windows sometimes the [VC++ x64 runtime](https://learn.microsoft.com/en-us/cpp/windows/latest-supported-vc-redist) if Windows reports a missing DLL — see **`README-RELEASE.txt`** in the zip.)

**Building `moonrun` from source** is different: that requires Go + a C toolchain (see **Build from source** below). End users grabbing a Release zip only extract and run.

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
| [docs/JOLT_WINDOWS_PARITY.md](docs/JOLT_WINDOWS_PARITY.md) | Native Jolt on Windows (CGO), building static `.a` files |
| [examples/README.md](examples/README.md) | Runnable sample programs |

More detail lives in **[docs/reference/](docs/reference/)** and **[ARCHITECTURE.md](ARCHITECTURE.md)**.

---

## Contributing

Contribution guidelines and development notes are in **[CONTRIBUTING.md](CONTRIBUTING.md)** and **[docs/DEVELOPER.md](docs/DEVELOPER.md)**. Continuous integration validates builds, tests, and representative `go run . --check` samples.

On **Windows**, a **`fullruntime`** link that pulls in Jolt requires prebuilt **`libJolt.a`** and **`libjolt_wrapper.a`** in **[third_party/jolt-go/jolt/lib/windows_amd64/](third_party/jolt-go/jolt/lib/windows_amd64/README.md)** (or build them with **`third_party/jolt-go/scripts/build-libs-windows.ps1`**). **`scripts/check-jolt-windows-libs.ps1`** checks that both files are present.

---

## License

**MIT** — see [LICENSE](LICENSE).
