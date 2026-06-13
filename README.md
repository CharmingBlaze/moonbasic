# moonBASIC

**A modern BASIC for building 2D and 3D games** — write `.mb` source, download pre-built binaries, and run. **No Go, no C compiler, no build tools** on your machine.

One download gives you **`moonrun`** (play games) and **`moonbasic`** (check, compile to `.mbc`, language server). The engine bundles **Raylib**, **Box2D**, and **Jolt** — graphics, audio, 2D physics, 3D physics, networking, terrain, UI, particles, and more behind **4,200+ built-in commands** across 40+ namespaces (`WINDOW.*`, `ENTITY.*`, `PHYSICS3D.*`, `GUI.*`, `NET.*`, …).

This repository is the **official home** for documentation, examples, and **pre-compiled downloads** for **Windows**, **Linux**, and **macOS (Apple Silicon)**.

---

## Download (nothing else required)

**Latest release:** **[github.com/CharmingBlaze/moonbasic/releases/latest](https://github.com/CharmingBlaze/moonbasic/releases/latest)**

**Direct links by platform:** **[charmingblaze.github.io/moonbasic/](https://charmingblaze.github.io/moonbasic/)**

1. Download the **full runtime** for your OS (table below).
2. Extract anywhere permanent.
3. `moonrun --version` (Windows: `moonrun.exe --version`).
4. `moonrun main.mb` or `moonbasic new MyGame` → `moonrun main.mb`.

You do **not** need to install Go, GCC, Node.js, or Raylib separately.

| Your goal | Download |
|-----------|----------|
| **moonBASIC IDE** (editor + compiler + runtime + docs — **easiest**) | **IDE bundle** — `moonbasic-<tag>-ide-windows-amd64.zip`, `linux-amd64.tar.gz`, or `macos-arm64.tar.gz` |
| **Play / make games** | **Full runtime** — `moonbasic-<tag>-windows-amd64.zip`, `linux-amd64.tar.gz`, or `macos-arm64.tar.gz` |
| **Lint / compile / LSP only** (no game window) | **Compiler only** — `moonbasic-<tag>-compiler-…` |
| **VS Code / Cursor** | Included in full-runtime zip + **`moonbasic-<tag>-vscode.vsix`** on Releases |

IDE archives include **`moonbasic-ide`**, **`moonbasic`**, **`moonrun`**, **`README-IDE-RELEASE.txt`**, and **`START-IDE`**. Documentation is built into the IDE — extract and run **START-IDE** to begin.

Full-runtime archives include **`moonbasic`**, **`moonrun`**, **`README-RELEASE.txt`**, the **VS Code `.vsix`**, and **`INSTALL-VSCODE.bat`** / **`INSTALL-VSCODE.sh`**. Details: **[RELEASES.md](RELEASES.md)**

---

## VS Code / Cursor (one command)

After extracting the full runtime:

```bash
moonbasic install-vscode
```

Or double-click **`INSTALL-VSCODE.bat`** (Windows) / run **`./INSTALL-VSCODE.sh`**.

That installs the extension and sets **`moonbasic.languageServerPath`** and **`moonbasic.moonrunPath`** automatically. Then open any `.mb` file — completions, hover help, **Ctrl+F5** run, **Ctrl+Shift+C** check, **Alt+H** help at cursor.

Extension source: **[editors/vscode-moonbasic/](editors/vscode-moonbasic/)** · Guide: **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md#vs-code-syntax-and-lsp)**

---

## Quick start

```bash
moonrun --version
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

| Command | What it does |
|---------|----------------|
| **`moonrun game.mb`** | Compile (if needed) and run — **primary entry point** |
| **`moonbasic --check game.mb`** | Parse and type-check without running |
| **`moonbasic game.mb`** | Compile to **`game.mbc`** bytecode |
| **`moonbasic --lsp`** | Language server (stdio) for editors |
| **`moonbasic new Name`** | New project: `main.mb`, `assets/`, `.vscode/` |
| **`moonbasic install-vscode`** | Install VS Code / Cursor extension + configure paths |

---

## License

**MIT** — see [LICENSE](LICENSE).

**Engine source (contributors):** [github.com/CharmingBlaze/moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)
