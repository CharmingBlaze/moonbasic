# Distribution artifacts

## “All commands” — what that means

1. **Full builtin catalog (language + tooling)**  
   Every builtin name the product knows about is defined in **`compiler/builtinmanifest/commands.json`**. Any **`moonbasic`** binary (compiler-only **or** full bundle) uses that catalog for **`--check`**, **`--lsp`**, and compiling to **`.mbc`**. So for **authors and CI**, you already have the **complete command set** at compile/check time — no extra download for “more command names” in the compiler.

2. **Running games (engine at runtime)**  
   Calling **`WINDOW.*`**, **`PHYSICS3D.*`**, etc. requires the **engine**, which ships as **`moonrun`** in the **full runtime** archives below — not in the compiler-only zip. Use the **full runtime** download when you need to **execute** the full surface area on a machine (graphics, physics, net, …).

3. **Release layout**  
   A version tag produces **four** assets on GitHub Releases: **two full-runtime** (Linux + Windows) and **two compiler-only**. Together they cover “all commands” for tooling + “all commands that can run” when you pair the right binary with the right workload.

## Two kinds of downloads

**Release binaries do not require Go, GCC, or Clang on the player’s machine** — `moonrun` compiles `.mb` internally. (Building *from source* still needs those tools; see `docs/BUILDING.md`.)

| Artifact | Contents | End-user needs |
|----------|----------|----------------|
| **Full runtime** (`moonbasic-<tag>-windows-amd64.zip` / `linux-amd64.tar.gz`) | `moonbasic` + `moonrun` + README | OS + GPU/OpenGL stack; may need VC++ redist on Windows (see `packaging/README-RELEASE.txt`). **Use this for “all commands” at run time** (games). |
| **Compiler only** (`moonbasic-<tag>-compiler-windows-amd64.zip` / `linux-amd64.tar.gz`) | `moonbasic` only, **CGO off** — no `raylib.dll` for the compiler | **Nothing extra** — full **`--check`** / compile / LSP against the same manifest; no `moonrun`. |

Build the compiler bundle locally: `scripts/release_compiler_windows.ps1` / `scripts/release_compiler_linux.sh` (see `docs/BUILDING.md`).

---

## Other paths in this folder

- **windows/** — NSIS script `moonbasic.nsi` builds an installer. Run on Windows with NSIS 3.x after placing `moonbasic.exe` and required MinGW DLLs next to the script (see CI release job).
- **linux/** — `build-appimage.sh` and `build-deb.sh` expect a staged tree under `dist/stage/` with `bin/moonbasic`, `share/moonbasic/{examples,assets}`.

Release CI (`.github/workflows/release.yml`) uploads **both** full-runtime and compiler-only archives on version tags. NSIS/AppImage integration may require local paths adjusted for your Raylib/MinGW layout.
