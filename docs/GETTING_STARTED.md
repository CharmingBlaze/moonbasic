# Getting Started with moonBASIC

Welcome to moonBASIC. This guide gets you from download to a running game in minutes.

> **New to game dev?** Start with **[Your First Hour](FIRST_HOUR.md)** for a gentle intro to the language and game loop.

---

## 1. Installation

Download pre-built binaries from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)** — no Go, GCC, or build tools required.

| Your goal | Download |
|-----------|----------|
| **Run games** (window, graphics, physics, audio) | **Full runtime:** `moonbasic-<tag>-windows-amd64.zip`, `moonbasic-<tag>-linux-amd64.tar.gz`, or `moonbasic-<tag>-macos-arm64.tar.gz` |
| **Compile / check / LSP only** (CI, editors — no window) | **Compiler only:** `moonbasic-<tag>-compiler-windows-amd64.zip` or `moonbasic-<tag>-compiler-linux-amd64.tar.gz` |

**Full runtime** includes **`moonbasic`**, **`moonrun`**, and **`README-RELEASE.txt`**. **Compiler only** ships **`moonbasic`** without **`moonrun`**.

See **[RELEASES.md](../RELEASES.md)** for a full breakdown of each archive.

Extract the archive. On Windows, keep both `.exe` files from the **same** release zip together.

**Browse downloads:** [charmingblaze.github.io/moonbasic](https://charmingblaze.github.io/moonbasic/)

### VS Code: syntax and LSP

1. On the **same release page**, download **`moonbasic-<tag>-vscode.vsix`**.
2. VS Code → Extensions → **⋯** → **Install from VSIX…**
3. If **`moonbasic`** is not on **PATH**, set **`moonbasic.languageServerPath`** to the full path of the binary.

### VS Code: debugging

Requires **full runtime** (`moonrun`):

1. Set **`moonbasic.moonrunPath`** if needed.
2. Open a `.mb` file, set breakpoints, **Run and Debug** → **Debug moonBASIC**.

Projects from **`moonbasic new`** include a `.vscode/launch.json` template.

### Start a new project

```bash
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

Creates **`main.mb`**, **`assets/`**, and editor debug config.

---

## 2. Ship your game

**Minimal (recommended):** Share your `.mb` / `.mbc` and assets. Tell players to install the **full runtime** for their OS from [Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest) (same tag you tested with).

**Bundle:** Ship a folder with **`moonrun`**, your game, and assets in one zip so players extract and run.

---

## 3. Run your first game

```bash
moonrun --version
moonrun path/to/game.mb
```

Windows: `moonrun.exe --version`

**`moonrun`** compiles `.mb` in-process when needed — no separate compile step for everyday dev.

---

## 4. Compiler tooling (`moonbasic`)

| Command | Purpose |
|---------|---------|
| `moonbasic --check game.mb` | Lint / type-check (no window) |
| `moonbasic game.mb` | Compile to **`game.mbc`** |
| `moonbasic --lsp` | Language server on stdio |
| `moonbasic new Name` | Scaffold a project |

Running games with a window always requires **`moonrun`** from a **full runtime** download.

---

## 5. Your first program

Create `hello.mb`:

```moonbasic
PRINT("Hello, moonBASIC!")
```

Run:

```bash
moonrun hello.mb
```

---

## 6. Opening a window

```moonbasic
WINDOW.OPEN(1280, 720, "moonBASIC")
WINDOW.SETFPS(60)

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    RENDER.CLEAR(30, 40, 60)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## 7. Examples in this repository

Clone or download this repo for the **`examples/`** folder (not included in release zips):

```bash
moonrun examples/spin_cube/main.mb
moonbasic --check examples/spin_cube/main.mb
```

See **[examples/README.md](../examples/README.md)**.

---

## Language features

| Feature | Example |
|---------|---------|
| String interpolation | `PRINT($"Score: {score}")` |
| Multi-return | `RETURN x, y, z` then `a, b, c = GetPos()` |
| Enums | `ENUM State … ENDENUM` |
| Array loops | `FOR EACH e IN arr … NEXT` |
| Typed functions | `FUNCTION f(x AS FLOAT) AS INT` |

Full syntax: **[LANGUAGE.md](LANGUAGE.md)**

---

## Next steps

| Topic | Document |
|-------|----------|
| Game loop & modules | [PROGRAMMING.md](PROGRAMMING.md) |
| Language syntax | [LANGUAGE.md](LANGUAGE.md) |
| Command index | [COMMANDS.md](COMMANDS.md) |
| API reference | [reference/](reference/) |
| Examples | [examples/README.md](../examples/README.md) |
| Blitz porting | [BLITZ3D_PORTING.md](BLITZ3D_PORTING.md) |

**Engine contributors:** source and build instructions are at [github.com/CharmingBlaze/moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) — not in this repository.
