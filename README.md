# moonBASIC

**A modern BASIC for building 2D and 3D games** — write `.mb` source, download pre-built binaries, and run. **No Go, no C compiler, no build tools** on your machine. Download from here: (https://github.com/CharmingBlaze/moonbasic/releases/tag/v1.2.27)

One download gives you **`moonrun`** (play games), **`moonbasic`** (check, compile to `.mbc`, language server), and optionally the **moonBASIC IDE** — a desktop editor with the full documentation built in. The engine bundles **Raylib**, **Box2D**, and **Jolt** — graphics, audio, 2D physics, 3D physics, networking, terrain, UI, particles, and more behind **4,200+ built-in commands** across 40+ namespaces (`WINDOW.*`, `ENTITY.*`, `PHYSICS3D.*`, `GUI.*`, `NET.*`, …).

This repository is the **official home** for documentation, examples, and **pre-compiled downloads** for **Windows**, **Linux**, and **macOS (Apple Silicon)**.

---

## Download (nothing else required)

**Latest release:** **[github.com/CharmingBlaze/moonbasic/releases/latest](https://github.com/CharmingBlaze/moonbasic/releases/latest)**

**Direct links by platform:** **[charmingblaze.github.io/moonbasic/](https://charmingblaze.github.io/moonbasic/)**

You do **not** need to install Go, GCC, Node.js, or Raylib separately.

| Your goal | Download |
|-----------|----------|
| **moonBASIC IDE** (editor + compiler + runtime + docs — **easiest**) | **IDE bundle** — `moonbasic-<tag>-ide-windows-amd64.zip`, `linux-amd64.tar.gz`, or `macos-arm64.tar.gz` |
| **Play / make games** (terminal only) | **Full runtime** — `moonbasic-<tag>-windows-amd64.zip`, `linux-amd64.tar.gz`, or `macos-arm64.tar.gz` |
| **Lint / compile / LSP only** (no game window) | **Compiler only** — `moonbasic-<tag>-compiler-…` |
| **VS Code / Cursor** | Included in full-runtime zip + **`moonbasic-<tag>-vscode.vsix`** on Releases |

**IDE bundle** — `moonbasic-ide`, `moonbasic`, `moonrun`, `README-IDE-RELEASE.txt`, `START-IDE`. Documentation is inside the IDE.

**Full runtime** — `moonbasic`, `moonrun`, `README-RELEASE.txt`, VS Code `.vsix`, `INSTALL-VSCODE.bat` / `INSTALL-VSCODE.sh`.

All artifact names and platform notes: **[RELEASES.md](RELEASES.md)**

---

## moonBASIC IDE (recommended)

The fastest way to start — one folder, no extra setup.

![moonBASIC IDE — editor, documentation, check, compile, and run in one app](docs/images/moonbasic-ide.png)

1. Download **`moonbasic-<tag>-ide-…`** for your OS from [Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest).
2. Extract anywhere permanent.
3. Run **START-IDE.bat** (Windows) or **`chmod +x START-IDE.sh moonbasic-ide moonbasic moonrun && ./START-IDE.sh`** (Linux/macOS).
4. Write or open a `.mb` file in the editor.

| Shortcut | Action |
|----------|--------|
| **F5** | Run game (`moonrun`) |
| **Ctrl+Shift+C** | Check syntax |
| **Ctrl+Shift+B** | Compile to `.mbc` |
| **Alt+H** | Help at cursor |
| **Gear menu** (title bar) | Themes, fonts, colors, settings |

The IDE auto-detects `moonbasic` and `moonrun` in the same folder. Open **Documentation** in the sidebar for **BEGIN_HERE.md**, guides, and the full command reference.

Engine source: [moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) (`moonbasic ide/`).

---

## Quick start (terminal)

If you downloaded the **full runtime** instead of the IDE:

```bash
moonrun --version          # Windows: moonrun.exe --version
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

**Try an example** from this repo (clone or Download ZIP — examples ship here, not inside release zips):

```bash
moonrun examples/spin_cube/main.mb
moonrun examples/guides/game_loop.mb
```

---

## VS Code / Cursor (one command)

After extracting the **full runtime**:

```bash
moonbasic install-vscode
```

Or double-click **`INSTALL-VSCODE.bat`** (Windows) / run **`./INSTALL-VSCODE.sh`**.

That installs the extension and sets **`moonbasic.languageServerPath`** and **`moonbasic.moonrunPath`** automatically. Then open any `.mb` file — completions, hover help, **Ctrl+F5** run, **Ctrl+Shift+C** check, **Alt+H** help at cursor.

Extension source: **[editors/vscode-moonbasic/](editors/vscode-moonbasic/)** · Guide: **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md#vs-code-syntax-and-lsp)**

---

## Documentation

| Start here | |
|------------|---|
| **[docs/BEGIN_HERE.md](docs/BEGIN_HERE.md)** | Install, first 10 minutes, learning path |
| **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)** | Ship your game, VS Code, player installs |
| **[docs/systems/GUIDES.md](docs/systems/GUIDES.md)** | Topic guides (physics, lighting, multiplayer, math, …) |
| **[docs/COMMANDS.md](docs/COMMANDS.md)** | Command index |
| **[docs/reference/](docs/reference/)** | Full API by namespace |
| **[web/command-browser.html](web/command-browser.html)** | Searchable offline command browser |
| **[examples/](examples/)** | Runnable demos + **`examples/guides/`** copies of doc examples |

---

## What's built in

| Layer | Technology |
|-------|------------|
| Graphics & audio | **Raylib** |
| 2D physics | **Box2D** |
| 3D physics & KCC | **Jolt** |
| Multiplayer | **ENet** (where enabled) |

---

## Example

```moonbasic
APP.OPEN(960, 540, "Hello moonBASIC")
WHILE NOT APP.SHOULDCLOSE()
    RENDER.CLEAR(20, 24, 32)
    RENDER.FRAME()
WEND
APP.CLOSE()
```

(`APP.*` aliases `WINDOW.*` — same engine, clearer tutorials.)

---

## Tooling

| Tool / command | What it does |
|----------------|--------------|
| **moonBASIC IDE** (`moonbasic-ide`) | Desktop editor — syntax, docs, check, compile, run, themes |
| **`moonrun game.mb`** | Compile (if needed) and run — **primary terminal entry point** |
| **`moonbasic --check game.mb`** | Parse and type-check without running |
| **`moonbasic game.mb`** | Compile to **`game.mbc`** bytecode |
| **`moonbasic --lsp`** | Language server (stdio) for other editors |
| **`moonbasic new Name`** | New project: `main.mb`, `assets/`, `.vscode/` |
| **`moonbasic install-vscode`** | Install VS Code / Cursor extension + configure paths |

---

## License

**MIT** — see [LICENSE](LICENSE).

**Engine source (contributors):** [github.com/CharmingBlaze/moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)
