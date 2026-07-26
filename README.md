# moonBASIC

**A modern BASIC for building 2D and 3D games** — download pre-built binaries and start coding. **No Go, no C compiler, no build tools** on your machine.

This repository is the **official download and documentation hub** for **Windows**, **Linux**, and **macOS (Apple Silicon)**. Engine source for contributors lives in **[moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)**.

---

## Download (nothing else required)

**Latest release:** **[v1.3.2 — Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)**

You do **not** need to install Go, GCC, Node.js, or Raylib separately.

| Your goal | Download |
|-----------|----------|
| **moonBASIC IDE** (editor + compiler + runtime + **full docs** — **easiest**) | `moonbasic-v1.3.2-ide-windows-amd64.zip`, `moonbasic-v1.3.2-ide-linux-amd64.tar.gz`, or `moonbasic-v1.3.2-ide-macos-arm64.tar.gz` |
| **Play / make games** (terminal only) | `moonbasic-v1.3.2-windows-amd64.zip`, `…-linux-amd64.tar.gz`, or `…-macos-arm64.tar.gz` |
| **Lint / compile / LSP only** (no game window) | `moonbasic-v1.3.2-compiler-…` |
| **VS Code / Cursor** | `moonbasic-v1.3.2-vscode.vsix` (also inside full-runtime zips) |

Artifact naming for any tag: **[RELEASES.md](RELEASES.md)**

---

## moonBASIC IDE (recommended)

One folder: desktop editor, **`moonbasic`**, **`moonrun`**, samples, and the **complete documentation** built into the IDE (Begin Here, guides, command reference).

1. Download the **IDE bundle** for your OS from [Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest).
2. Extract anywhere permanent.
3. Run **START-IDE.bat** (Windows), **START-IDE.command** (macOS), or **`./START-IDE.sh`** (Linux).
4. Open **Documentation** in the sidebar, or open `samples/hello.mb` and press **F5**.

| Shortcut | Action |
|----------|--------|
| **F5** | Run game (`moonrun`) |
| **Ctrl+Shift+C** | Check syntax |
| **Ctrl+Shift+B** | Compile to `.mbc` |
| **Alt+H** | Help at cursor |

Engine / IDE source: [moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) (`ide/`).

---

## Documentation

Browse online in this repo, or inside the IDE sidebar:

| Start here | |
|------------|---|
| **[docs/BEGIN_HERE.md](docs/BEGIN_HERE.md)** | Install, first 10 minutes, learning path |
| **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)** | IDE, full runtime, VS Code, shipping games |
| **[docs/FIRST_HOUR.md](docs/FIRST_HOUR.md)** | Friendly first hour |
| **[docs/systems/GUIDES.md](docs/systems/GUIDES.md)** | Topic guides |
| **[docs/COMMANDS.md](docs/COMMANDS.md)** | Command index |
| **[docs/reference/](docs/reference/)** | API by namespace |
| **[examples/](examples/)** | Runnable demos |

---

## Quick start (terminal)

If you downloaded the **full runtime** instead of the IDE:

```bash
moonrun --version          # Windows: moonrun.exe --version
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

```bash
moonrun examples/spin_cube/main.mb
moonrun examples/guides/game_loop.mb
```

---

## VS Code / Cursor

After extracting the **full runtime**:

```bash
moonbasic install-vscode
```

Or install **`moonbasic-v1.3.2-vscode.vsix`** from Releases (Extensions → Install from VSIX…).

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

---

## License

**MIT** — see [LICENSE](LICENSE).

**Engine source (contributors):** [github.com/CharmingBlaze/moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)
