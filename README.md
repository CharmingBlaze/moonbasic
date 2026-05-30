# moonBASIC

**A modern BASIC for building 2D and 3D games** — write `.mb` source, download pre-built binaries, and run. No Go, no C compiler, no build step required.

The engine ships with **Raylib** (graphics, audio, input), **Box2D** (2D physics), and **Jolt** (3D physics) built in — one language, one toolchain, ready to make games.

This repository is the **official home** for documentation, examples, and **pre-compiled downloads** for **Windows**, **Linux**, and **macOS (Apple Silicon)**.

---

## Where to download

**Latest release (recommended):**  
**https://github.com/CharmingBlaze/moonbasic/releases/latest**

**Download landing page (direct links by platform):**  
**https://charmingblaze.github.io/moonbasic/**

Every release is a tagged bundle on the **Releases** tab of this repository. You do **not** clone or compile anything to play — download the zip/tar for your OS, extract, and run.

---

## How to download

1. Open **[Releases → Latest](https://github.com/CharmingBlaze/moonbasic/releases/latest)** (or use the [download page](https://charmingblaze.github.io/moonbasic/) for one-click links).
2. Find the asset row for your platform (see table below).
3. Click the file name to download.
4. Extract the archive to a folder you keep (Desktop, `Projects`, etc.).
5. Open a terminal **in that extracted folder**.
6. Run `moonrun --version` (Windows: `moonrun.exe --version`) to confirm it works.

Inside each **full runtime** archive you get **`moonbasic`**, **`moonrun`**, and **`README-RELEASE.txt`** with platform-specific notes.

Full details on every file name: **[RELEASES.md](RELEASES.md)**

---

## Which file do I need?

Replace `<tag>` with the release version shown on the page (e.g. `v1.2.26`).

| Your goal | Download this asset |
|-----------|---------------------|
| **Play or make games** (window, 2D/3D graphics, physics, audio) | **Full runtime** — `moonbasic-<tag>-windows-amd64.zip`, `moonbasic-<tag>-linux-amd64.tar.gz`, or `moonbasic-<tag>-macos-arm64.tar.gz` |
| **Compile / lint / LSP only** (CI, editors — no game window) | **Compiler only** — `moonbasic-<tag>-compiler-windows-amd64.zip`, `moonbasic-<tag>-compiler-linux-amd64.tar.gz`, or `moonbasic-<tag>-compiler-macos-arm64.tar.gz` |
| **VS Code** (syntax highlighting, LSP, debugger) | **`moonbasic-<tag>-vscode.vsix`** on the same release page |

- **Full runtime** includes **`moonrun`** (the game engine) and **`moonbasic`** (compiler, `--check`, `--lsp`).
- **Compiler only** includes **`moonbasic`** only — no window, no `moonrun`.

---

## Quick start (full runtime)

After you extract the **full runtime** archive:

```bash
moonrun --version
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

Windows:

```bat
moonrun.exe --version
moonbasic.exe new MyGame
cd MyGame
moonrun.exe main.mb
```

**Run an example** from this repo (clone or **Code → Download ZIP** — examples are not inside the release zips):

```bash
moonrun examples/spin_cube/main.mb
```

---

## What's built in

moonBASIC is a complete game stack in one download:

| System | Role |
|--------|------|
| **Raylib** | Window, rendering, textures, sprites, audio, input |
| **Box2D** | 2D rigid-body physics, joints, collisions |
| **Jolt Physics** | 3D physics, character controllers, collision layers |

You call them through moonBASIC commands (`WINDOW.*`, `RENDER.*`, `PHYSICS2D.*`, `PHYSICS3D.*`, etc.) — no separate installs.

---

## What's in this repository

| Path | Purpose |
|------|---------|
| [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md) | Install, first window, VS Code, ship your game |
| [docs/LANGUAGE.md](docs/LANGUAGE.md) | Syntax, variables, control flow, functions |
| [docs/PROGRAMMING.md](docs/PROGRAMMING.md) | Game loop, 2D/3D patterns |
| [docs/COMMANDS.md](docs/COMMANDS.md) | Built-in command index |
| [docs/reference/](docs/reference/) | Full API reference by namespace |
| [examples/](examples/) | Runnable sample games and demos |
| [web/command-browser.html](web/command-browser.html) | Searchable offline command browser (open in a browser) |
| [RELEASES.md](RELEASES.md) | Release artifact guide |

---

## Example

```moonbasic
WINDOW.OPEN(960, 540, "Hello moonBASIC")
WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR WINDOW.SHOULDCLOSE())
    RENDER.CLEAR(20, 24, 32)
    RENDER.FRAME()
WEND
WINDOW.CLOSE()
```

More samples: **[examples/README.md](examples/README.md)** · 3D cube: **`examples/spin_cube/`**

---

## Tooling

| Command | What it does |
|---------|----------------|
| **`moonrun game.mb`** | Compile (if needed) and run with the full engine |
| **`moonbasic --check game.mb`** | Parse and type-check without running |
| **`moonbasic game.mb`** | Compile to **`game.mbc`** bytecode |
| **`moonbasic --lsp`** | Language server on stdio (for editors) |
| **`moonbasic new Name`** | Scaffold a new project with `main.mb` and an assets folder |

---

## License

**MIT** — see [LICENSE](LICENSE).
