# Getting Started with MoonBASIC

Welcome to MoonBASIC. Whether you are installing the engine for the first time or writing your first lines of code, this guide will get you up and running in minutes.

> [!TIP]
> **New to game development?**
> Start with **[MoonBASIC: Your First Hour](FIRST_HOUR.md)** for a friendly introduction to the language, modern **Method Chaining**, and rapid prototyping.

---

## 1. Installation

Use the **compiled distribution** from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)** — official Windows/Linux archives with **`moonbasic`** and (for the full runtime) **`moonrun`**. That is the supported way to run games and use the compiler: **no Go, no GCC, no local build** of the engine.

You only need a **clone or ZIP of this repo** if you want example `.mb` sources or documentation; everyday play and compile use **only** the extracted release binaries.

Pick the file that matches what you need (replace `<tag>` with the release version, e.g. `v1.2.20`):

| Your goal | Download |
|-----------|----------|
| **Run games** (window, graphics, physics, audio) | **Full runtime:** `moonbasic-<tag>-windows-amd64.zip` or `moonbasic-<tag>-linux-amd64.tar.gz` |
| **Compile / check / LSP only** (CI, tooling, no `moonrun`) | **Compiler only:** `moonbasic-<tag>-compiler-windows-amd64.zip` or `moonbasic-<tag>-compiler-linux-amd64.tar.gz` |

**Full runtime** includes **`moonbasic`** and **`moonrun`** plus `README-RELEASE.txt`. **Compiler only** ships in a folder such as **`MoonBasic-compiler/`** with **`moonbasic`** (or **`moonbasic.exe`**) and a short readme — there is **no** `moonrun` in that bundle.

Extract the archive somewhere permanent — on **Windows**, keep **every file** from the full-runtime zip in one folder (runtime DLLs sit next to `moonrun.exe`; copying only the `.exe` elsewhere can trigger **“Entry Point Not Found”** / `nanosleep64` errors). On **Windows**, use `moonbasic.exe` in the examples below; on **Linux**, use `./moonbasic` if the binary is not on your `PATH`.

More detail on what each archive contains: **[`dist/README.md`](../dist/README.md)** (in the source tree) or the **[main README](https://github.com/CharmingBlaze/moonbasic#download-and-use-recommended)** on GitHub.

To **build moonbasic from source** (contributors), see **[BUILDING.md](BUILDING.md)**.

---

## 2. Using the moonbasic compiler

Open a terminal in the directory that contains **`moonbasic`** (on **compiler-only** installs, that is usually inside **`MoonBasic-compiler/`**).

### Check the binary

```bash
moonbasic --version
```

On Windows:

```bat
moonbasic.exe --version
```

### Lint without running (`--check`)

Parses and type-checks a program and reports errors. Does not require `moonrun` or a GPU.

```bash
moonbasic --check path/to/game.mb
```

Use this in editors, pre-commit hooks, and CI.

### Compile to bytecode (`.mbc`)

```bash
moonbasic path/to/game.mb
```

This writes **`game.mbc`** next to **`game.mb`** (same base name). The compiler does not run the game — it only produces bytecode.

### Language server (`--lsp`)

For editor integration, run:

```bash
moonbasic --lsp
```

Configure your editor’s MoonBASIC/LSP client to use **stdio**. The same **`moonbasic`** binary serves **`--check`**, compile, and **`--lsp`**; the full builtin list is always available in the compiler.

### Running games

Running **`.mb`** or **`.mbc`** with a window requires **`moonrun`** from a **full runtime** download:

```bash
moonrun path/to/game.mb
```

Release **`moonrun`** compiles `.mb` in-process when needed, then starts the engine — you still do **not** need Go or a separate compiler install on the player machine for pre-built releases.

If you only installed the **compiler-only** archive, use **`moonbasic`** to produce **`.mbc`** here and run **`moonrun`** on another machine that has the **full runtime**, or install the full runtime on the same machine.

---

## 3. Your First Program

Create a file named `hello.mb`:

```basic
PRINT "Hello, MoonBASIC!"
```

Run it using the runtime (full runtime install):

```bash
moonrun hello.mb
```

---

## 4. Opening a Window

MoonBASIC makes window management effortless. Create `display.mb`:

```basic
WINDOW.OPEN(1280, 720, "MoonBASIC Window")
WINDOW.SETFPS(60)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(30, 40, 60)
    DRAW.TEXT("Press ESC to exit", 540, 350, 20, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## 5. Modern 3D with Method Chaining

MoonBASIC supports **Method Chaining** (Fluent API), allowing you to configure objects in a single, readable line.

```basic
WINDOW.OPEN(1280, 720, "3D Cube Demo")
cam = CAMERA.CREATE().SETPOS(0, 5, 10).SETTARGET(0, 0, 0)
cube = ENTITY.CREATECUBE(2.0).SETCOLOR(100, 200, 255, 255)

WHILE NOT WINDOW.SHOULDCLOSE()
    ; Update rotation using a fluent method
    cube.SETROT(0, TIME.GET() * 50, 0)

    RENDER.CLEAR(10, 10, 20)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(50, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

---

## 6. Modern Blitz-Style (High Fidelity)

For advanced users, MoonBASIC provides a "High Fidelity" path with PBR materials, dynamic lighting, and SSAO.

```basic
WINDOW.OPEN(1920, 1080, "Project: High Fidelity")
cam = CAMERA.CREATE().SETPOS(0, 5, 10)
sun = LIGHT.CREATEDIRECTIONAL(0, -1, 0, 255, 255, 200, 2.0)

; Load a high-poly model with modern effects
car = ENTITY.LOADMESH("supercar.glb").SETPBR(0.9, 0.1)
RENDER.SETSSAO(TRUE)
RENDER.SETBLOOM(0.8)

WHILE NOT WINDOW.SHOULDCLOSE()
    CAMERA.FOLLOWENTITY(cam, car, 10.0, 3.0, 5.0)
    
    ENTITY.UPDATE(TIME.DELTA())

    RENDER.CLEAR(12, 14, 22)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

---

## Next Steps

Explore the specialized documentation to master every aspect of the engine:

| Topic | Reference |
|-------|-----------|
| **Core Workflow** | [Programming Guide](PROGRAMMING.md) |
| **Language Syntax** | [Language Reference](LANGUAGE.md) |
| **3D Entities** | [Entity Reference](reference/ENTITY.md) |
| **Physics** | [Physics 3D Reference](reference/PHYSICS3D.md) |
| **Atmosphere** | [Camera & Render Reference](reference/CAMERA_LIGHT_RENDER.md) |
| **Gameplay Helpers** | [Beginner Full Stack](reference/BEGINNER_FULL_STACK.md) |

**Happy Coding!**
