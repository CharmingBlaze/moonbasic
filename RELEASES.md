# moonBASIC release artifacts

Every tagged release on **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases)** ships **complete user packages**. You do **not** need Go, GCC, or this repository's source tree.

Replace `<tag>` with the release version (e.g. `v1.3.3`).

---

## Complete package (all platforms)

Each download includes:

- **`moonbasic-ide`** — desktop editor  
- **`moonbasic`** — compiler / `--check` / LSP  
- **`moonrun`** — game runtime  
- **`docs/`** — full documentation offline (also in the IDE sidebar)  
- **`samples/`** — starter scripts  
- **START-IDE.*** launchers  

| Asset | Platform |
|-------|----------|
| `moonbasic-<tag>-ide-windows-amd64.zip` **or** `moonbasic-<tag>-windows-amd64.zip` | Windows x64 (same contents) |
| `moonbasic-<tag>-ide-linux-amd64.tar.gz` **or** `moonbasic-<tag>-linux-amd64.tar.gz` | Linux x64 (same contents) |
| `moonbasic-<tag>-ide-macos-arm64.tar.gz` **or** `moonbasic-<tag>-macos-arm64.tar.gz` | macOS Apple Silicon (same contents) |

**Use when:** you want to write and run moonBASIC games. **No other tools required.**

1. Extract anywhere permanent.  
2. Run **START-IDE.bat** (Windows), **START-IDE.command** (macOS), or **`chmod +x` + `./START-IDE.sh`** (Linux).  
3. Status bar should show **Toolchain ready**.

**Shortcuts:** F5 run · Ctrl+Shift+C check · Ctrl+Shift+B compile · Alt+H help at cursor.

Engine source: [moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler) (`ide/`).

---

## VS Code / Cursor extension (optional)

| Asset | Contents |
|-------|----------|
| `moonbasic-<tag>-vscode.vsix` | Syntax, snippets, LSP client |

```bash
# Or: Extensions → Install from VSIX…
```

---

## After extract

1. Keep every file from the archive together (including any `lib*.dll` on Windows).  
2. Read **`README-IDE-RELEASE.txt`**.  
3. Browse **`docs/`** offline or use the IDE Documentation sidebar.

---

## Platform notes

### Windows

Official builds link Raylib and MinGW runtime pieces into the `.exe` files when possible; some zips also ship `libstdc++` / `libwinpthread` sidecars. Keep them next to `moonrun.exe`.

### Linux

Expect a normal desktop **OpenGL** stack and **glibc**-compatible distro.

### macOS

Apple Silicon (arm64). Make binaries executable if needed: `chmod +x START-IDE.sh moonbasic moonrun`.

---

## Building from source

Not supported from **this** repository. Engine contributors use **[moonbasic-compiler](https://github.com/CharmingBlaze/moonbasic-compiler)**.

---

## Recent releases

### v1.3.3

- **Every OS download is complete:** IDE + compiler + moonrun + `docs/` + samples (short platform names are aliases of the IDE packages).  
- Offline **`docs/`** folder included beside the IDE.  
- Compiler-only / thin runtime-only zips are no longer published to this hub.

### v1.3.2

- Physics stub fix for `BODY3D.COMMIT` / `SETLIFT`; `moonrun -version` reports Jolt backend.  
- See [docs/JOLT_WINDOWS_PARITY.md](docs/JOLT_WINDOWS_PARITY.md).

### v1.3.1

- Windows MinGW DLL / unzip-and-run fix for `moonrun` and IDE F5.
