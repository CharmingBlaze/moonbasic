# Distribution artifacts

## Two kinds of downloads

| Artifact | Contents | End-user needs |
|----------|----------|----------------|
| **Full runtime** (`moonbasic-<tag>-windows-amd64.zip` / `linux-amd64.tar.gz`) | `moonbasic` + `moonrun` + README | OS + GPU/OpenGL stack; may need VC++ redist on Windows (see `packaging/README-RELEASE.txt`). |
| **Compiler only** (`moonbasic-<tag>-compiler-windows-amd64.zip` / `linux-amd64.tar.gz`) | `moonbasic` only, **CGO off** — no `raylib.dll` for the compiler | **Nothing extra** (no Python, no C compiler) — compile `.mb` → `.mbc`, `--check`, `--lsp`, `--disasm` only. |

Build the compiler bundle locally: `scripts/release_compiler_windows.ps1` / `scripts/release_compiler_linux.sh` (see `docs/BUILDING.md`).

---

## Other paths in this folder

- **windows/** — NSIS script `moonbasic.nsi` builds an installer. Run on Windows with NSIS 3.x after placing `moonbasic.exe` and required MinGW DLLs next to the script (see CI release job).
- **linux/** — `build-appimage.sh` and `build-deb.sh` expect a staged tree under `dist/stage/` with `bin/moonbasic`, `share/moonbasic/{examples,assets}`.

Release CI (`.github/workflows/release.yml`) uploads **both** full-runtime and compiler-only archives on version tags. NSIS/AppImage integration may require local paths adjusted for your Raylib/MinGW layout.
