# moonBASIC

## Stop fighting your engine.

**Godot** wants you to learn its nodes. **Unity** wants your email and a project wizard. **Unreal** wants three hours before you see a cube.  

**moonBASIC** wants you to type **`Window.Open`**, drop in a **`Camera`**, and feel a **hardware-accelerated**, **Jolt-powered** 3D world come alive while your coffee is still hot.

This is a **language that looks like C/C++ but feels like BASIC**—built for **indie devs** who miss the **honest simplicity of the 90s** but refuse to give up **modern rendering, physics, and networking**. It is not a toy interpreter: it is a **real compiler** (lexer → AST → bytecode) plus a **production runtime** wired into the same libraries serious games use today.

---

## Modern power, retro soul

We are not here to cosplay the past. We are here to **ship**.

- **Vertical integration beats glue code.** One toolchain, one mental model, one place where your game actually runs.
- **2D and 3D parity.** Build a **PS1-style 3D horror** prototype or a **SNES-style 2D platformer** without changing your workflow—same commands, same loop, same clarity.
- **Aesthetic intent.** moonBASIC leans into a **Windows 98 / PS1 “cyber-retro”** energy: flat UI honesty, chunky clarity, and room for **your** art direction—not a generic PBR showroom unless **you** want that.

**Bold claim, human proof:** the **compiler is stable**, the **2D and 3D pipelines are live**, and the standard library already spans **Tiled maps**, **PBR-style materials**, **sprites and atlases**, **particles**, **audio**, **lights**, **shaders**, and more—see the [documentation index](#documentation) below.

---

## The engine for people who actually want to ship

| You have been… | moonBASIC offers… |
|----------------|-------------------|
| Drowning in editor complexity | **One executable**, **one `.mb` file**, fast iteration |
| Fighting a scripting ceiling | A **compiled** pipeline—bytecode + VM—not a slow REPL |
| Stuck with “good enough” physics | **Dual engines**: **Jolt** for 3D, **Box2D** for 2D—**no** forced compromise |
| Treating multiplayer as a weekend project | **Networking that feels native**—see [Multiplayer](#multiplayer-the-go-edge) |

We pick **developer speed** and **clarity** over the “big guys” and their ceremony. Opinionated? **Yes.** Apologetic? **No.**

---

## What’s under the hood

These are the pillars—not marketing stickers. The runtime is actually built on them.

| | Role |
|---|------|
| **Go** | **Insane compile speed**, simple cross-platform builds, and a runtime host that stays out of your way. |
| **Raylib (raylib-go)** | **Fast**, **clean**, **hardware-accelerated** 2D/3D, input, audio—ideal for indie scope. |
| **Jolt Physics** | The industry’s current **darling** for **high-performance 3D** rigid simulation. |
| **Box2D** | The **gold standard** for **pixel-perfect 2D**—not a toy substitute. |
| **ENet** | **Reliable UDP** for game-style client/server traffic. |

**Physics choice matters.** Most engines lock you into **one** physics story. moonBASIC gives you **Jolt** when the problem is **three-dimensional**, and **Box2D** when the problem is **flat and precise**—because **indie games are not all the same shape**.

---

## Multiplayer: the Go edge

**Multiplayer commands** are not a bolt-on script layer. Because the **host runtime is Go**, networking can lean on **goroutines** and **sane concurrency** so the **heavy lifting**—connections, buffers, timing—lives in the **engine**, not in your face.

You think in **simple multiplayer verbs**; the **Go backend** carries synchronization and state so it feels like **part of the language**, not a research project. Start with **[NETWORK.md](docs/reference/NETWORK.md)** and the **`NETWORK.*`** surface in **[COMMANDS.md](docs/COMMANDS.md)**.

---

## What does it look like?

A compact **3D** sample—the full [examples/spin_cube](examples/spin_cube/main.mb) demo adds a grid and on-screen text. **No** project boilerplate: just your script.

```basic
WINDOW.OPEN(960, 540, "Spinning cube")
WINDOW.SETFPS(60)

cam = CAMERA.MAKE()
CAMERA.SETPOSITION(cam, 0, 2, 8)
CAMERA.SETTARGET(cam, 0, 0, 0)
CAMERA.SETFOV(cam, 45)

cube = MESH.MAKECUBE(2, 2, 2)
mat  = MATERIAL.MAKEDEFAULT()
cubeXform = TRANSFORM.IDENTITY()
angle# = 0

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    dt# = TIME.DELTA()
    angle# = angle# + 45 * dt#
    TRANSFORM.SETROTATION(cubeXform, 0, angle#, 0)
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

---

## Architecture (ten seconds)

1. **Compile** — `moonbasic mygame.mb` → bytecode **`.mbc`**.
2. **Run** — the **VM** executes bytecode and calls into **Raylib**, **Jolt**, **Box2D**, **ENet**, etc., via **CGO** where enabled.

You get **BASIC-level clarity** at the keyboard and a **real pipeline** under the hood—**not** a giant interpreter loop pretending to be an engine.

---

## Getting started

### Pre-built binaries (fastest)

**Windows x64** and **Linux x64** builds ship on **[Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)**.

| Platform | Archive | Inside |
|----------|---------|--------|
| **Windows** | `moonbasic-v*-windows-amd64.zip` | `moonbasic.exe`, `moonrun.exe`, **README-RELEASE.txt** |
| **Linux** | `moonbasic-v*-linux-amd64.tar.gz` | `moonbasic`, `moonrun`, **README-RELEASE.txt** |

Extract anywhere, read **README-RELEASE.txt**, then run **`moonbasic --version`** or **`moonrun`** on a sample under [`examples/`](examples/).

### Build from source

- **Go** — see **`go.mod`** (toolchain **1.25.3+**).
- **C toolchain** — **MinGW-w64** on Windows, **GCC** on Linux (for full **CGO** graphics and physics).

```bash
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic

# Windows (example: MSYS2 MinGW64)
set CGO_ENABLED=1
set CC=C:\msys64\mingw64\bin\gcc.exe
go build -o moonbasic.exe .

# Linux
CGO_ENABLED=1 go build -o moonbasic .
```

**Graphical programs** need **`-tags fullruntime`** and **`moonrun`** (or **`moonbasic --run`** from a full-runtime build). A plain **`go run . file.mb`** compile-only path emits **`.mbc`** without opening a window—see **[AGENTS.md](AGENTS.md)**, **[CONTRIBUTING.md](CONTRIBUTING.md)**, and **[docs/BUILDING.md](docs/BUILDING.md)**.

```bash
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

---

## Documentation

| Document | What it covers |
|----------|----------------|
| [CONTRIBUTING.md](CONTRIBUTING.md) | PR workflow, builtins, **`go run . --check`**. |
| [docs/DEVELOPER.md](docs/DEVELOPER.md) | Repo layout, build tags, commands. |
| [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) | Install, first window, mental model. |
| [docs/PROGRAMMING.md](docs/PROGRAMMING.md) | Game loop, modules, 2D/3D, platforms. |
| [docs/LANGUAGE.md](docs/LANGUAGE.md) | Variables, control flow, functions. |
| [docs/COMMANDS.md](docs/COMMANDS.md) | Built-in index. |
| [docs/EXAMPLES.md](docs/EXAMPLES.md) | Narrated snippets. |
| [examples/README.md](examples/README.md) | Runnable **`main.mb`** programs. |

### API reference (modules)

| Module | Topics |
|--------|--------|
| [Window](docs/reference/WINDOW.md) | Open/close, monitors |
| [Render](docs/reference/RENDER.md) | Clear, frame, screenshots |
| [Camera](docs/reference/CAMERA.md) | 3D and 2D cameras |
| [Draw 2D](docs/reference/DRAW2D.md) | Primitives, text, textures |
| [GUI](docs/reference/GUI.md) | Immediate-mode UI |
| [Texture](docs/reference/TEXTURE.md) / [Image](docs/reference/IMAGE.md) | GPU textures, CPU images |
| [Sprite](docs/reference/SPRITE.md) / [Atlas](docs/reference/ATLAS.md) | Strips, animation, atlases |
| [Model / Mesh / Material](docs/reference/MODEL.md) | 3D assets, PBR-style materials |
| [Shader](docs/reference/SHADER.md) / [Light](docs/reference/LIGHT.md) | GLSL, lighting |
| [Physics 3D (Jolt)](docs/reference/PHYSICS3D.md) | Rigid bodies, scenes |
| [Physics 2D (Box2D)](docs/reference/PHYSICS2D.md) | 2D colliders |
| [Character controller](docs/reference/CHARCONTROLLER.md) | Kinematic player |
| [Tilemap (Tiled)](docs/reference/TILEMAP.md) | **`.tmx`** maps |
| [Network](docs/reference/NETWORK.md) | **ENet**, multiplayer |
| [Input](docs/reference/INPUT.md) / [Audio](docs/reference/AUDIO.md) | Input and sound |

More in **[docs/reference/](docs/reference/)** and **[ARCHITECTURE.md](ARCHITECTURE.md)**.

---

## Contributing

Read **[CONTRIBUTING.md](CONTRIBUTING.md)** and **[docs/DEVELOPER.md](docs/DEVELOPER.md)**. CI exercises **`go test ./...`**, **`go build`**, and **`go run . --check`** on representative samples. For AI/editor context: **[AGENTS.md](AGENTS.md)**.

---

## License

**MIT** — see [LICENSE](LICENSE).

---

*A manifesto for a new era of indie gamedev: **fast tools**, **honest syntax**, **real libraries**, and **your** game on screen—**today**.*
