# moonBASIC release artifacts

Every tagged release on **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases)** ships pre-built binaries. You do **not** need Go, GCC, or this repository's source tree to use them.

Replace `<tag>` with the release version (e.g. `v1.2.28`).

---

## moonBASIC IDE (recommended)

| Asset | Platform | Contents |
|-------|----------|----------|
| `moonbasic-<tag>-ide-windows-amd64.zip` | Windows x64 | `moonbasic-ide.exe`, `moonbasic.exe`, `moonrun.exe`, `README-IDE-RELEASE.txt`, `START-IDE.bat` |
| `moonbasic-<tag>-ide-linux-amd64.tar.gz` | Linux x64 | `moonbasic-ide`, `moonbasic`, `moonrun`, `README-IDE-RELEASE.txt`, `START-IDE.sh` |
| `moonbasic-<tag>-ide-macos-arm64.tar.gz` | macOS Apple Silicon | same as Linux layout |

**Use when:** you want the easiest setup — desktop editor, full documentation inside the IDE, check, compile, and run games. **No other tools required.**

1. Extract anywhere permanent.
2. Run **START-IDE.bat** (Windows) or **`chmod +x` + `./START-IDE.sh`** (Linux/macOS).
3. Status bar should show **Toolchain ready** — `moonbasic` and `moonrun` sit in the same folder.

**Shortcuts:** F5 run · Ctrl+Shift+C check · Ctrl+Shift+B compile · Alt+H help at cursor · gear menu for themes/fonts.

Engine source: [moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) (`moonbasic ide/`).

---

## Full runtime (games)

| Asset | Platform | Contents |
|-------|----------|----------|
| `moonbasic-<tag>-windows-amd64.zip` | Windows x64 | `moonbasic.exe`, `moonrun.exe`, `README-RELEASE.txt`, `moonbasic-<tag>-vscode.vsix`, `INSTALL-VSCODE.bat` |
| `moonbasic-<tag>-linux-amd64.tar.gz` | Linux x64 | `moonbasic`, `moonrun`, `README-RELEASE.txt`, `.vsix`, `INSTALL-VSCODE.sh` |
| `moonbasic-<tag>-macos-arm64.tar.gz` | macOS Apple Silicon | `moonbasic`, `moonrun`, `README-RELEASE.txt`, `.vsix`, `INSTALL-VSCODE.sh` |

**Use when:** you want to open a window, run graphics, physics, audio, and play or develop games.

**Primary command:** `moonrun yourgame.mb` — compiles `.mb` in-process when needed, then runs the engine.

**Also includes `moonbasic`** for `--check`, `--lsp`, `moonbasic new`, and standalone `.mbc` output.

---

## Compiler only (tooling)

| Asset | Platform | Contents |
|-------|----------|----------|
| `moonbasic-<tag>-compiler-windows-amd64.zip` | Windows x64 | `moonbasic.exe` (+ short readme in `MoonBasic-compiler/`) |
| `moonbasic-<tag>-compiler-linux-amd64.tar.gz` | Linux x64 | `moonbasic` binary only |
| `moonbasic-<tag>-compiler-macos-arm64.tar.gz` | macOS Apple Silicon | `moonbasic` binary only |

**Use when:** CI linting, `--check`, `--lsp`, or compiling `.mb` → `.mbc` **without** linking the game engine.

**Does not include `moonrun`** — cannot open a game window.

The compiler uses the **same built-in command catalog** as the full runtime for `--check` and diagnostics.

---

## VS Code / Cursor extension

| Asset | Contents |
|-------|----------|
| `moonbasic-<tag>-vscode.vsix` | Syntax, snippets, LSP, check/compile/run commands, debugger |
| (in full-runtime zip) | Same `.vsix` + **`INSTALL-VSCODE.bat`** / **`INSTALL-VSCODE.sh`** |

**Easiest install (full runtime folder):**

```bash
moonbasic install-vscode
```

Or double-click **`INSTALL-VSCODE.bat`** (Windows) / **`./INSTALL-VSCODE.sh`**.

**Manual:** VS Code → Extensions → **⋯** → **Install from VSIX…**

The installer sets `moonbasic.languageServerPath` and `moonbasic.moonrunPath` when binaries sit beside the archive.

See [docs/GETTING_STARTED.md](docs/GETTING_STARTED.md#vs-code-syntax-and-lsp) · source: [editors/vscode-moonbasic/](editors/vscode-moonbasic/).

---

## After extract

1. Keep **full-runtime** executables from the **same** zip together — do not mix `moonrun` from one release with `moonbasic` from another.
2. Read **`README-RELEASE.txt`** in the archive for platform-specific notes.
3. Clone or download **this repository** for **`examples/`** and documentation — release zips do not bundle the example tree.

---

## Platform notes

### Windows

Official builds link Raylib and MinGW runtime pieces into the `.exe` files — you should **not** need companion DLLs beside the binaries. If Windows reports a missing DLL, re-extract the **entire** zip from one release; avoid dropping stray MinGW DLLs next to the exes.

Some systems may need the [Microsoft VC++ x64 redistributable](https://learn.microsoft.com/en-us/cpp/windows/latest-supported-vc-redist).

### Linux

Full-runtime binaries expect a normal desktop **OpenGL** stack and **glibc**-compatible distro (built on a current Ubuntu-style CI image). Install GPU drivers; no compiler `-dev` packages are required to **run** releases.

### macOS

Apple Silicon (arm64) full-runtime tarball. Make binaries executable: `chmod +x moonbasic moonrun`.

---

## Building from source

Not supported from **this** repository. Engine contributors use **[moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)**.
