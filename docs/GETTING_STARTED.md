# Getting Started with MoonBASIC

Welcome to MoonBASIC. Whether you are installing the engine for the first time or writing your first lines of code, this guide will get you up and running in minutes.

> [!TIP]
> **New to game development?**
> Start with **[MoonBASIC: Your First Hour](FIRST_HOUR.md)** for a friendly introduction to the language, modern **Method Chaining**, and rapid prototyping.

> [!TIP]
> **Looking up commands by game system?**
> Browse the **[40-system guide](systems/README.md)** (`APP`, `RENDER`, `ENTITY`, physics, audio, saves, …) — each page follows [DOCUMENTATION_STYLE_GUIDE.md](DOCUMENTATION_STYLE_GUIDE.md).

---

## 1. Installation

Use the **compiled distribution** from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest)** — official Windows/Linux archives with **`moonbasic`** and (for the full runtime) **`moonrun`**. That is the supported way to run games and use the compiler: **no Go, no GCC, no local build** of the engine.

You only need a **clone or ZIP of this repo** if you want example `.mb` sources or documentation; everyday play and compile use **only** the extracted release binaries.

Pick the file that matches what you need (replace `<tag>` with the release version, e.g. `v1.2.20`):

| Your goal | Download |
|-----------|----------|
| **Run games** (window, graphics, physics, audio) | **Full runtime:** `moonbasic-<tag>-windows-amd64.zip`, `moonbasic-<tag>-linux-amd64.tar.gz`, or `moonbasic-<tag>-macos-arm64.tar.gz` |
| **Compile / check / LSP only** (CI, tooling, no `moonrun`) | **Compiler only:** `moonbasic-<tag>-compiler-windows-amd64.zip` or `moonbasic-<tag>-compiler-linux-amd64.tar.gz` |

**Full runtime** includes **`moonbasic`** and **`moonrun`** plus `README-RELEASE.txt`. **Compiler only** ships in a folder such as **`MoonBasic-compiler/`** with **`moonbasic`** (or **`moonbasic.exe`**) and a short readme — there is **no** `moonrun` in that bundle.

Extract the archive somewhere permanent — on **Windows**, keep the **full-runtime** zip contents together (both `.exe` files from the **same** release; do not mix executables from different builds). On **Windows**, use `moonbasic.exe` in the examples below; on **Linux**, use `./moonbasic` if the binary is not on your `PATH`.

More detail on what each archive contains: **[`dist/README.md`](../dist/README.md)** (in the source tree) or the **[main README](https://github.com/CharmingBlaze/moonbasic-compiler#download-and-use-recommended)** on GitHub.

### VS Code: syntax and LSP

After you extract the **full runtime** zip, install the editor extension in **one step** (no Node.js, no manual VSIX):

| Platform | Easiest |
|----------|---------|
| **Windows** | Double-click **`INSTALL-VSCODE.bat`** in the extracted folder |
| **Linux / macOS** | Run **`./INSTALL-VSCODE.sh`** in the extracted folder |
| **Any** | **`moonbasic install-vscode`** |

The release zip already includes **`moonbasic-*-vscode.vsix`**. The installer finds **VS Code** or **Cursor**, installs the extension, and sets **`moonbasic.languageServerPath`** / **`moonbasic.moonrunPath`** to the **`moonbasic`** / **`moonrun`** next to the installer.

If no local VSIX is found, **`moonbasic install-vscode`** downloads the latest **`moonbasic-*-vscode.vsix`** from GitHub automatically.

**Manual fallback:** Extensions → **⋯** → **Install from VSIX…** → pick the `.vsix` in the zip.

You get **syntax highlighting**, **snippets**, **LSP** (completions, hover help, diagnostics), and **moonBASIC** commands. **Ctrl+F5** run, **Ctrl+Shift+C** check, **Alt+H** help at cursor. Projects from **`moonbasic new`** include **`.vscode/`** for check, compile, run, and debug.

**From source:** `powershell -File scripts/install-vscode-extension.ps1` or `./scripts/install-vscode-extension.sh`

### VS Code: debugging (full runtime)

The **`moonbasic-<tag>-vscode.vsix`** extension can **debug** `.mb` games when you have **`moonrun`** from a **full runtime** archive (not compiler-only):

1. Run **`moonbasic install-vscode`** (or **`INSTALL-VSCODE.bat`** / **`./INSTALL-VSCODE.sh`**) from the extracted full-runtime folder — paths are configured automatically.
2. Open a **`.mb`** file, set breakpoints, then **Run and Debug** → **Debug moonBASIC** (or use **`.vscode/launch.json`** from **`moonbasic new`**).
3. Only if needed: **Settings** → **`moonbasic.moonrunPath`** → path to **`moonrun`** when it is not beside **`moonbasic`**.

Debugging uses **`moonrun --dap`** under the hood. Breakpoints pause the game; the **Globals** scope shows live variables.

### Start a new project (`moonbasic new`)

From any folder where **`moonbasic`** is on **`PATH`** (full runtime or compiler-only):

```bash
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

This creates **`main.mb`**, **`assets/`**, **`.vscode/`** (launch, tasks, extensions), and a short **`README.md`**. Edit **`main.mb`**, then run with **`moonrun`**.

---

## 2. Ship your game (for authors)

You can share games in two straightforward ways:

**A — Minimal install for players (recommended)**  
Ship your **`.mb`** source and/or **`.mbc`** bytecode, plus any **assets** (images, sounds, data files) using the **paths your scripts expect** (working directory when they run `moonrun`, or paths you set with **`RES.PATH`** and similar APIs). Tell players to install the **same [full runtime](#1-installation) archive** for their OS from [Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) — **not** the **compiler-only** download (that bundle has no `moonrun` and cannot open a game window). Prefer the **same moonBASIC release tag** you used to build and test: bytecode and engine behavior stay aligned across patch versions.

**B — Folder bundle (one zip per game)**  
Ship a folder that contains **`moonrun`** (and optionally **`moonbasic`**) next to your game and assets so players extract and run from that folder. Copy **`moonrun.exe`** / **`moonrun`** from the **same full-runtime release** you used to test. On **Linux**, use the official tarball layout or **`moonbasic package linux`** from your project folder.

---

## 3. Run your first game (`moonrun`)

After extracting the **full runtime**, use **`moonrun`** as your main entry point — it compiles `.mb` in-process when needed and starts the engine (same as BlitzBASIC's single-exe workflow).

```bash
moonrun --version
moonrun path/to/game.mb
```

On Windows:

```bat
moonrun.exe --version
moonrun.exe path\to\game.mb
```

You do **not** need a separate compile step for everyday development. Use **`moonbasic`** when you want lint/CI (`--check`), editor LSP (`--lsp`), or standalone `.mbc` files — see the next section.

---

## 4. Using the moonbasic compiler (optional tooling)

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

## 5. Your First Program

Create a file named `hello.mb`:

```basic
PRINT "Hello, MoonBASIC!"
```

Run it using the runtime (full runtime install):

```bash
moonrun hello.mb
```

---

## 6. Opening a Window

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

## 7. Modern 3D with Method Chaining

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

## 8. Modern Blitz-Style (High Fidelity)

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

## Language features (2026)

Recent syntax additions (full detail in [LANGUAGE.md](LANGUAGE.md)):

| Feature | Example |
|---------|---------|
| String interpolation | `PRINT($"Score: {score}")` |
| Multi-return | `RETURN x, y, z` then `a, b, c = GetPos()` |
| Enums | `ENUM State … ENDENUM`, `IF s = State.IDLE` |
| Array loops | `FOR EACH e IN enemies … NEXT` |

Runnable demos: [examples/tilemap](../examples/tilemap/README.md), [examples/gamepad](../examples/gamepad/README.md). Roadmap: [ROADMAP.md](ROADMAP.md).

---

## Next Steps

Explore the specialized documentation to master every aspect of the engine:

| Topic | Reference |
|-------|-----------|
| **Start here (install + why)** | [BEGIN_HERE.md](BEGIN_HERE.md) · [systems/00-START.md](systems/00-START.md) |
| **Every beginner command listed** | [systems/COMMAND_REGISTRY.md](systems/COMMAND_REGISTRY.md) |
| **40 game systems** | [systems/README.md](systems/README.md) |
| **Core Workflow** | [Programming Guide](PROGRAMMING.md) |
| **Language Syntax** | [Language Reference](LANGUAGE.md) — includes **`$"..."`**, **`ENUM`**, multi-return |
| **Roadmap** | [ROADMAP.md](ROADMAP.md) — shipped vs planned language work |
| **Examples** | [examples/README.md](../examples/README.md) — tilemap, gamepad, platformer, … |
| **3D Entities** | [Entity Reference](reference/ENTITY.md) |
| **Physics** | [Physics 3D Reference](reference/PHYSICS3D.md) |
| **Atmosphere** | [Camera & Render Reference](reference/CAMERA_LIGHT_RENDER.md) |
| **Gameplay Helpers** | [Beginner Full Stack](reference/BEGINNER_FULL_STACK.md) |

**Happy Coding!**
