# moonBASIC

**A modern BASIC for building 2D and 3D games** — download pre-built binaries and start coding. **No Go, no C compiler, no build tools** on your machine.

This repository is the **official download and documentation hub** for **Windows**, **Linux**, and **macOS (Apple Silicon)**. Engine source for contributors lives in **[moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)**.

---

## Download (nothing else required)

**Latest release:** **[v1.3.3 — Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)**

Every OS package is **complete** for users: **IDE** + **compiler (`moonbasic`)** + **runtime (`moonrun`)** + **`docs/`** + samples.

| Your OS | Download (either name — same contents) |
|---------|----------------------------------------|
| **Windows x64** | `moonbasic-v1.3.3-ide-windows-amd64.zip` or `moonbasic-v1.3.3-windows-amd64.zip` |
| **Linux x64** | `moonbasic-v1.3.3-ide-linux-amd64.tar.gz` or `moonbasic-v1.3.3-linux-amd64.tar.gz` |
| **macOS Apple Silicon** | `moonbasic-v1.3.3-ide-macos-arm64.tar.gz` or `moonbasic-v1.3.3-macos-arm64.tar.gz` |
| **VS Code / Cursor** (optional) | `moonbasic-v1.3.3-vscode.vsix` |
| **The Book** (optional) | `moonBASIC-The-Book.zip` — Markdown + Word + PDF |

Artifact naming: **[RELEASES.md](RELEASES.md)** · Book online: **[docs/THE_MOONBASIC_BOOK.md](docs/THE_MOONBASIC_BOOK.md)**

---

## Start in 30 seconds

1. Download the package for your OS from [Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest).
2. Extract anywhere permanent.
3. Run **START-IDE.bat** (Windows), **START-IDE.command** (macOS), or **`./START-IDE.sh`** (Linux).
4. Open **Documentation** in the sidebar (or browse the `docs/` folder), or open `samples/hello.mb` and press **F5**.

| Shortcut | Action |
|----------|--------|
| **F5** | Run game (`moonrun`) |
| **Ctrl+Shift+C** | Check syntax |
| **Ctrl+Shift+B** | Compile to `.mbc` |
| **Alt+H** | Help at cursor |

---

## What's in the folder

| File | Role |
|------|------|
| `moonbasic-ide` / `.exe` / `.app` | Desktop IDE |
| `moonbasic` / `.exe` | Compiler + LSP |
| `moonrun` / `.exe` | Game runtime |
| `docs/` | Full documentation offline |
| `samples/` | Starter scripts |

Engine / IDE source: [moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) (`ide/`).

---

## Documentation

Browse online in this repo, inside the IDE sidebar, or the `docs/` folder in your download:

| Start here | |
|------------|---|
| **[docs/THE_MOONBASIC_BOOK.md](docs/THE_MOONBASIC_BOOK.md)** | **The Book** — funny full guide ([Word/PDF](docs/book/)) |
| **[docs/BEGIN_HERE.md](docs/BEGIN_HERE.md)** | Install, first 10 minutes, learning path |
| **[docs/GETTING_STARTED.md](docs/GETTING_STARTED.md)** | IDE, shipping games, VS Code |
| **[docs/FIRST_HOUR.md](docs/FIRST_HOUR.md)** | Friendly first hour |
| **[docs/systems/GUIDES.md](docs/systems/GUIDES.md)** | Topic guides |
| **[docs/COMMANDS.md](docs/COMMANDS.md)** | Command index |
| **[docs/reference/](docs/reference/)** | API by namespace |
| **[examples/](examples/)** | Runnable demos |

---

## Quick start (terminal)

`moonbasic` and `moonrun` sit beside the IDE. Optional: run **ADD-TO-PATH** so any terminal can find them.

```bash
moonrun --version
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

---

## VS Code / Cursor

Install **`moonbasic-v1.3.3-vscode.vsix`** from Releases (Extensions → Install from VSIX…), or use the IDE if you prefer an all-in-one editor.

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
