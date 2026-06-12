# Begin here — using moonBASIC

> Install once, write `.mb` files, run with **`moonrun`**. This page explains what to download, which tool to use, and where every command is documented.

---

## Table of contents

- [What is moonBASIC?](#what-is-moonbasic)
- [Install (5 minutes)](#install-5-minutes)
- [moonbasic vs moonrun — which do I use?](#moonbasic-vs-moonrun--which-do-i-use)
- [Your first 10 minutes](#your-first-10-minutes)
- [The game loop (and why each step exists)](#the-game-loop-and-why-each-step-exists)
- [Where every command is documented](#where-every-command-is-documented)
- [Learning path](#learning-path)
- [See also](#see-also)

---

## What is moonBASIC?

moonBASIC is a **game programming language** and runtime:

- You write **`.mb`** source files (BASIC-style syntax).
- **`moonrun`** opens a window, runs physics, draws graphics, and plays audio.
- **`moonbasic`** checks syntax, compiles to **`.mbc`** bytecode, and powers the editor (LSP).

You do **not** need Go, C++, or a compiler toolchain on your machine — use the **pre-built release** from GitHub.

**Why commands instead of a giant engine class?**  
Built-ins are grouped by **namespace** (`WINDOW.OPEN`, `ENTITY.SETPOS`, …) so you can search docs by *what you are trying to do* (open window, move entity, play sound). Aliases like **`APP.OPEN`** map to the same behavior as **`WINDOW.OPEN`** for readability in tutorials.

---

## Install (5 minutes)

1. Open [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest).
2. Download the **full runtime** for your OS (`windows-amd64.zip` or `linux-amd64.tar.gz`).
3. Extract somewhere permanent. You should see **`moonbasic`** and **`moonrun`** (or `.exe` on Windows).
4. **VS Code / Cursor (one command):** run **`moonbasic install-vscode`** — or double-click **`INSTALL-VSCODE.bat`** (Windows) / **`./INSTALL-VSCODE.sh`** (Linux/macOS) inside the extracted folder. No manual VSIX install.

Details: [GETTING_STARTED.md](GETTING_STARTED.md).

---

## moonbasic vs moonrun — which do I use?

| Tool | Purpose | Opens a window? |
|------|---------|-----------------|
| **`moonrun game.mb`** | **Play** your game — compile if needed, then run the engine | Yes (full runtime) |
| **`moonbasic --check game.mb`** | **Validate** code without running (CI, editor, quick errors) | No |
| **`moonbasic game.mb`** | Write **`game.mbc`** bytecode next to the source | No |
| **`moonbasic --lsp`** | Language server for VS Code | No |
| **`moonbasic new MyGame`** | Create `main.mb`, `assets/`, project folder | No (then use `moonrun`) |

**Rule:** Day-to-day development = edit `.mb` → **`moonrun main.mb`**. Use **`moonbasic --check`** before you commit or when the editor shows diagnostics.

---

## Your first 10 minutes

```bash
moonbasic new MyFirstGame
cd MyFirstGame
moonrun main.mb
```

You should see a window. Open **`main.mb`** in an editor and change the title or colors.

**Check without running:**

```bash
moonbasic --check main.mb
```

**Ask the runtime for help on one command:**

```basic
HELP("WINDOW.OPEN")
```

---

## The game loop (and why each step exists)

Every interactive moonBASIC program follows the same rhythm:

```
setup once → WHILE running → update logic → draw → present frame → WEND → cleanup
```

| Step | Typical commands | Why you need it |
|------|------------------|-----------------|
| **Open display** | `APP.OPEN` / `WINDOW.OPEN` | Creates the OS window and graphics context. Nothing draws until this runs. |
| **Frame rate** | `APP.SETFPS(60)` | Keeps simulation and input stable across machines. |
| **Create world** | `ENTITY.*`, `CAMERA.*`, `LIGHT.*` | Entities are your game objects; camera defines the view; lights make 3D readable. |
| **Loop test** | `WHILE NOT APP.SHOULDCLOSE()` | Runs until the user closes the window or you quit. |
| **Simulation** | `INPUT.*`, `PHYSICS.STEP`, `ENTITY.MOVE` | Read controls, advance physics, move things — use **`APP.DELTA()`** so speed is frame-independent. |
| **Clear** | `RENDER.CLEAR(r,g,b)` | Wipes the previous frame’s pixels and depth buffer. |
| **3D pass** | `RENDER.BEGIN(cam)` … `SCENE.DRAW()` … `RENDER.END()` | Sets camera matrices and draws 3D entities in depth order. |
| **2D / HUD** | `DRAW.TEXT`, `GUI.*` | Overlay text and menus after or before 3D depending on your design. |
| **Present** | `RENDER.FRAME()` | Shows the finished frame on screen. **Without this, you see a frozen or blank window.** |
| **Cleanup** | `APP.CLOSE`, `ENTITY.FREE`, `TEXTURE.FREE` | Releases GPU and OS resources cleanly. |

Annotated line-by-line walkthrough: [systems/00-START.md](systems/00-START.md).

Minimal 3D sample: [examples/foundation/main.mb](../examples/foundation/main.mb).

---

## Where every command is documented

moonBASIC registers **thousands** of command overloads. They are all listed in generated registries — you do not have to memorize them.

| Resource | What it contains |
|----------|------------------|
| **[systems/00-START.md](systems/00-START.md)** | Why beginners use specific commands; first loop explained |
| **[systems/GUIDES.md](systems/GUIDES.md)** | **Deep guides:** entity, 2D/3D collision, UI, multiplayer |
| **[systems/README.md](systems/README.md)** | 40 game systems in teaching order |
| **[systems/01-CORE.md](systems/01-CORE.md)** … **[11-TOOLING.md](systems/11-TOOLING.md)** | Narrative guides with examples per system |
| **[systems/COMMAND_REGISTRY.md](systems/COMMAND_REGISTRY.md)** | **Complete list** of commands for beginner namespaces (arity + returns) |
| **[API_CONSISTENCY.md](API_CONSISTENCY.md)** | **Every** registered command in the entire engine |
| **[COMMAND_AUDIT.md](COMMAND_AUDIT.md)** | Every namespace → deep reference page |
| **[COMMANDS.md](COMMANDS.md)** | Topic index (globals, math, strings, …) |
| **`HELP("COMMAND")`** | Quick console help while coding |
| **`docs/reference/*.md`** | Deep dives per namespace (`ENTITY.md`, `CAMERA.md`, …) |

**Case:** `window.open`, `WINDOW.OPEN`, and `Window.Open` are the same command.

---

## Learning path

1. [BEGIN_HERE.md](BEGIN_HERE.md) (this page) + [systems/00-START.md](systems/00-START.md)
2. [FIRST_HOUR.md](FIRST_HOUR.md) — language primer
3. [systems/GUIDES.md](systems/GUIDES.md) — **24 topic guides** covering all 40 beginner systems (+ multiplayer)
4. [systems/01-CORE.md](systems/01-CORE.md) → follow build order in [systems/README.md](systems/README.md)
4. [PROGRAMMING.md](PROGRAMMING.md) — 2D/3D patterns
5. [examples/README.md](../examples/README.md) — runnable demos
6. Look up any command in [COMMAND_REGISTRY.md](systems/COMMAND_REGISTRY.md) or [API_CONSISTENCY.md](API_CONSISTENCY.md)

---

## See also

- [GETTING_STARTED.md](GETTING_STARTED.md) — install, ship, VS Code
- [LANGUAGE.md](LANGUAGE.md) — syntax, types, `IMPORT`
- [DOCUMENTATION_STYLE_GUIDE.md](DOCUMENTATION_STYLE_GUIDE.md) — how docs are organized
