# moonBASIC

**A modern BASIC for building 2D and 3D games** — write `.mb` source, run with pre-built binaries. No Go, no compiler toolchain, no engine source required.

This repository is the **official distribution**: documentation, examples, and **GitHub Releases** with compiled **`moonbasic`** and **`moonrun`** for **Windows**, **Linux**, and **macOS (Apple Silicon)**.

> **Engine developers:** compiler and runtime source live at [github.com/CharmingBlaze/moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) — not here.

---

## Download

**[GitHub Releases → latest](https://github.com/CharmingBlaze/moonbasic/releases/latest)**

| What you need | Download |
|---------------|----------|
| **Play and build games** (window, graphics, physics, audio) | **Full runtime** — `moonbasic-<tag>-windows-amd64.zip`, `moonbasic-<tag>-linux-amd64.tar.gz`, or `moonbasic-<tag>-macos-arm64.tar.gz` |
| **Compile / lint / LSP only** (CI, editors — no game window) | **Compiler only** — `moonbasic-<tag>-compiler-windows-amd64.zip` or `moonbasic-<tag>-compiler-linux-amd64.tar.gz` |
| **VS Code** (syntax, LSP, debugger) | **`moonbasic-<tag>-vscode.vsix`** on the same release page |

**Browse downloads:** [charmingblaze.github.io/moonbasic](https://charmingblaze.github.io/moonbasic/) (static landing page with direct links)

See **[RELEASES.md](RELEASES.md)** for what each archive contains and platform notes.

---

## Quick start (full runtime)

1. Download and extract the **full runtime** zip/tar for your OS from [Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest).
2. Open a terminal in that folder.
3. Check it works:

   ```bash
   moonrun --version
   ```

   Windows: `moonrun.exe --version`

4. Start a new game:

   ```bash
   moonbasic new MyGame
   cd MyGame
   moonrun main.mb
   ```

5. Or run an example from this repo (clone or download ZIP):

   ```bash
   moonrun examples/spin_cube/main.mb
   ```

Read **`README-RELEASE.txt`** inside the archive for troubleshooting (Windows DLL notes, Linux GPU stack, etc.).

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

More samples: **[examples/README.md](examples/README.md)** · full 3D cube: **`examples/spin_cube/`**

---

## Tooling

| Command | What it does |
|---------|----------------|
| **`moonrun game.mb`** | Compile (if needed) and run with the full engine |
| **`moonbasic --check game.mb`** | Parse and type-check without running |
| **`moonbasic game.mb`** | Compile to **`game.mbc`** bytecode |
| **`moonbasic --lsp`** | Language server on stdio (for editors) |
| **`moonbasic new Name`** | Scaffold a new project with `main.mb` and assets folder |

---

## License

**MIT** — see [LICENSE](LICENSE).
