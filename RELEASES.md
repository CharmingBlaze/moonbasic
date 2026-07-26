# moonBASIC release artifacts

Every tagged release on **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases)** ships pre-built binaries. You do **not** need Go, GCC, or this repository's source tree to use them.

Replace `<tag>` with the release version (e.g. `v1.3.2`).

---

## moonBASIC IDE (recommended)

Includes the desktop editor **and built-in documentation** (Begin Here, guides, full command reference), plus `moonbasic` and `moonrun`.

| Asset | Platform | Contents |
|-------|----------|----------|
| `moonbasic-<tag>-ide-windows-amd64.zip` | Windows x64 | `moonbasic-ide.exe`, `moonbasic.exe`, `moonrun.exe`, docs inside IDE, `samples/`, `START-IDE.bat`, `README-IDE-RELEASE.txt` |
| `moonbasic-<tag>-ide-linux-amd64.tar.gz` | Linux x64 | `moonbasic-ide`, `moonbasic`, `moonrun`, docs inside IDE, `samples/`, `START-IDE.sh` |
| `moonbasic-<tag>-ide-macos-arm64.tar.gz` | macOS Apple Silicon | same layout as Linux |

**Use when:** easiest setup — edit, read docs, check, compile, and run. **No other tools required.**

1. Extract anywhere permanent.
2. Run **START-IDE.bat** (Windows), **START-IDE.command** (macOS), or **`chmod +x` + `./START-IDE.sh`** (Linux).
3. Status bar should show **Toolchain ready**.

**Shortcuts:** F5 run · Ctrl+Shift+C check · Ctrl+Shift+B compile · Alt+H help at cursor.

Engine source: [moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) (`ide/`).

---

## Full runtime (games)

| Asset | Platform | Contents |
|-------|----------|----------|
| `moonbasic-<tag>-windows-amd64.zip` | Windows x64 | `moonbasic.exe`, `moonrun.exe`, `README-RELEASE.txt`, `moonbasic-<tag>-vscode.vsix`, `INSTALL-VSCODE.bat` |
| `moonbasic-<tag>-linux-amd64.tar.gz` | Linux x64 | `moonbasic`, `moonrun`, `README-RELEASE.txt`, `.vsix`, `INSTALL-VSCODE.sh` |
| `moonbasic-<tag>-macos-arm64.tar.gz` | macOS Apple Silicon | `moonbasic`, `moonrun`, `README-RELEASE.txt`, `.vsix`, `INSTALL-VSCODE.sh` |

**Primary command:** `moonrun yourgame.mb`

Browse documentation in this repo’s [`docs/`](docs/) tree (same content the IDE ships).

---

## Compiler only (tooling)

| Asset | Platform | Contents |
|-------|----------|----------|
| `moonbasic-<tag>-compiler-windows-amd64.zip` | Windows x64 | `moonbasic.exe` |
| `moonbasic-<tag>-compiler-linux-amd64.tar.gz` | Linux x64 | `moonbasic` |

**Does not include `moonrun`** — cannot open a game window.

---

## VS Code / Cursor extension

| Asset | Contents |
|-------|----------|
| `moonbasic-<tag>-vscode.vsix` | Syntax, snippets, LSP client |
| (in full-runtime zip) | Same `.vsix` + **`INSTALL-VSCODE.bat`** / **`INSTALL-VSCODE.sh`** |

```bash
moonbasic install-vscode
```

---

## After extract

1. Keep executables from the **same** release together.
2. Read **`README-IDE-RELEASE.txt`** or **`README-RELEASE.txt`** in the archive.
3. Clone or browse **this repository** for **`docs/`** and **`examples/`** online — IDE bundles also embed the docs.

---

## Platform notes

### Windows

Official builds link Raylib and MinGW runtime pieces into the `.exe` files — you should **not** need companion DLLs. Some systems may need the [Microsoft VC++ x64 redistributable](https://learn.microsoft.com/en-us/cpp/windows/latest-supported-vc-redist).

### Linux

Full-runtime / IDE binaries expect a normal desktop **OpenGL** stack and **glibc**-compatible distro.

### macOS

Apple Silicon (arm64). Make binaries executable: `chmod +x moonbasic-ide moonbasic moonrun`.

---

## Building from source

Not supported from **this** repository. Engine contributors use **[moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)**.

---

## Recent releases

### v1.3.2

- **Physics stub fix:** `BODY3D.COMMIT` no longer returns `NULL` on non-CGO / soft-stub builds (fixes `SETLIFT` crash on `physics_power_test`).
- **`moonrun -version`** reports `Jolt backend: native` vs `stub` so you can tell if real Jolt is linked.
- Docs clarify **Windows or Linux + CGO** for native Jolt (official IDE / fullruntime zips ship native). See [docs/JOLT_WINDOWS_PARITY.md](docs/JOLT_WINDOWS_PARITY.md).

### v1.3.1

- **Windows fix:** `moonrun` / IDE **F5** no longer exits immediately from missing MinGW DLLs. Release zips embed or ship `libstdc++` / `libwinpthread` so unzip-and-run works on clean PCs.
- Prefer this tag over **v1.3.0** on Windows.

### v1.3.0

- Ship IDE bundles for **Windows**, **Linux**, and **macOS** (editor + compiler + runtime + **built-in documentation** + samples).
- Full runtime, compiler-only, and VS Code `.vsix` for the same tag.
- Repository cleanup: public download hub docs/examples synced from the engine; install guides point here.
- Prior hardening: shared handle-call dispatch, multi-error LSP diagnostics, BOM/INCLUDE path fixes, offline IDE CodeMirror, CI smoke on Win/Linux/macOS.

### v1.2.29

- Fix chained entity/camera method calls (`ENTITY.CREATECUBE(…).scale(…).pos(…).col(…)`, `CAMERA.CREATE().fov(…)`).
- Harden PBR material initialization on Windows purego builds so 3D examples (e.g. `examples/spin_cube`) render correctly.
