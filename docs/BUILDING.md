# Building moonBASIC from Source

This guide provides detailed instructions for compiling the `moonBASIC` interpreter from its source code.

---

## Dependencies

Before you can build, you need the following software installed on your system.

### All Systems
- **Go**: Version **1.25.3** or later (see [`go.mod`](../go.mod)). Download from the [official Go website](https://go.dev/dl/).
- **Git**: For cloning the repository.

### raylib-go and “Raylib 5.5”

In **[gen2brain/raylib-go](https://github.com/gen2brain/raylib-go)**, Git tags such as **`raylib/v0.55.0`** / **`v0.55.1`** correspond to bindings for **Raylib C 5.5**. This repository does **not** pin those tags today: it uses a **newer** `raylib` + `raygui` module version (**`v0.56.0-dev`**-style pseudo-version) because **`GUI.*`** and other code target the **current raygui** Go API (`ControlID`, `PropertyID`, `SetAlpha`, color helpers, etc.). Downgrading only the module to **`v0.55.x`** without a large port breaks the build.

**OpenGL profile (Windows CGO):** Unless you pass alternate **`raylib-go` build tags** (e.g. **`opengl43`** on the **`raylib`** module), the upstream **`cgo_windows*.go`** files default to **`-DGRAPHICS_API_OPENGL_33`**. moonBASIC does **not** enable **`opengl43`** by default, so linked Raylib is aimed at **OpenGL 3.3**, which matches most integrated GPUs from the last decade.

For the **native Raylib library** (`raylib.dll`, `libraylib.so`, …), install a **C Raylib** build whose **ABI matches** the **raylib-go** revision you compile against (check upstream release notes for that commit). If you specifically need **Raylib C 5.5** artifacts, pair them with **`raylib-go` `v0.55.x`** only after adapting `runtime/mbgui` and any other API-drift call sites.

### Raylib 5.5 and “Go only” (no CGO / no C compiler)

**What you can get:** On **Windows**, you can build with **`CGO_ENABLED=0`** so the **Go toolchain never invokes a C compiler** and **core `github.com/gen2brain/raylib-go/raylib`** uses the **purego** backend: it loads a prebuilt **`raylib.dll`** at runtime via [`purego`](https://github.com/ebitengine/purego). You still **ship that DLL** (or put it on `PATH`); it is the normal Raylib **native** library, not a second Go implementation of Raylib.

**Raylib 5.5 pairing:** Upstream tags **`raylib/v0.55.x`** are the **Go bindings** aimed at **Raylib C 5.5**. Your **`raylib.dll`** should be a **5.5** build from the same family so symbols match. This repository currently pins a **newer** `raylib-go` revision; for a strict **5.5** stack you would use **`v0.55.x`** bindings **and** a **5.5** DLL once the code is ported (see above).

**What is not “Go only” here:** Upstream **`raygui-go`** is **CGO + C**. On **Windows** with **`CGO_ENABLED=0`**, moonbasic still provides a **minimal** `GUI.*` layer drawn with Raylib (not full raygui). Advanced widgets (text entry, list views, `.rgs` themes, etc.) still need **CGO**. **`DB.*`** defaults to **`mattn/go-sqlite3`** (CGO); for **pure Go** SQLite with **`CGO_ENABLED=0`**, build with **`-tags modernc_sqlite`** ([`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite)). **ENet** still needs **CGO** for the linked **libenet** path.

**Linux / macOS:** **gen2brain/raylib-go** does **not** ship a non-CGO desktop Raylib for non-Windows; you link Raylib with **CGO** there.

### Windows
- **Optional: Zig CC wrapper** — For experimenting with static-friendly Windows builds, see [`scripts/build_static.ps1`](../scripts/build_static.ps1) (sets `CGO_ENABLED=1` and `CC="zig cc -target x86_64-windows-gnu"` before `go build -tags fullruntime ./cmd/moonrun`). You still need Raylib and GLFW/GL headers/libs on the compiler search path (e.g. MSYS2 MinGW + `mingw-w64-x86_64-raylib`).

- **C toolchain (recommended full build)**  
  For the default **CGO** build (linked **raylib**, **raygui**, ENet, SQLite, etc.), install a C compiler. We recommend **MinGW-w64** via **MSYS2**:
  1. Install MSYS2 from [https://www.msys2.org/](https://www.msys2.org/).
  2. In the MSYS2 **MINGW64** shell, install GCC (and optionally `mingw-w64-x86_64-raylib` if you link against the system library):
     ```bash
     pacman -S mingw-w64-x86_64-gcc
     ```
- **Pure Go on Windows (no CGO)**  
  You can build with **`CGO_ENABLED=0`** so **core Raylib** comes from **raylib-go’s purego** path (loads **`raylib.dll`** at runtime for **shipping** interactive binaries). **`GUI.*`** uses a **built-in minimal** widget set (see [GUI.md](reference/GUI.md)); full **raygui** still needs **CGO**. **`DB.*`** can use **`-tags modernc_sqlite`** (no CGO); **ENet** still requires **CGO**.  
  Place **`raylib.dll`** (matching your Raylib 5.x ABI) next to the executable or on **`PATH`** when you **run** a windowed game built this way. **`go test`** on Windows: vendored purego **`init()`** **defers** loading the DLL for **`*.test`** binaries (and when **`MOONBASIC_SKIP_RAYLIB_DLL=1`**) so packages that import **`raylib-go`** do not panic just because the DLL is absent—see [HAL_AND_RENDERING.md](architecture/HAL_AND_RENDERING.md). By default, **`go build .`** at the repo root and **`go build ./cmd/moonbasic`** produce a **compiler-only** binary (no game runtime): **`--check`**, **`.mbc`**, **`--lsp`**, and **`--disasm`** need no Raylib. For **`moonbasic --run`**, **`moonrun`**, or **`pipeline.RunProgram`**, build with **`-tags fullruntime`**. For the full test matrix (all packages, all tags), follow [DEVELOPER.md](DEVELOPER.md) and CI; some integration tests still expect CGO or a display stack on Linux.

### Linux (Debian / Ubuntu)
- **A C Compiler and Libraries**: You'll need `gcc` and the development headers for the libraries `raylib` depends on.
  ```bash
  sudo apt-get update
  sudo apt-get install -y gcc libgl1-mesa-dev libxi-dev \
    libxcursor-dev libxrandr-dev libxinerama-dev \
    libwayland-dev libxkbcommon-dev
  ```

---

## Build Steps

### 1. Clone the Repository

First, get the source code from GitHub:
```bash
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic
```

### 2. Build on Windows

Open a standard Command Prompt (`cmd.exe`) or PowerShell. You must tell Go where to find the C compiler you installed.

```bat
REM Set the CGO_ENABLED flag to allow Go to call C code
set CGO_ENABLED=1

REM Point to the MinGW GCC compiler (adjust path if you installed MSYS2 elsewhere)
set CC=C:\msys64\mingw64\bin\gcc.exe

REM Build the executable
go build -o moonbasic.exe .
```

**Optional — Windows without CGO (Raylib purego):** from `cmd` or PowerShell, no MinGW required for the Go link step:

```bat
set CGO_ENABLED=0
go build -o moonbasic.exe .
```

Ensure **`raylib.dll`** is available at runtime. For **full raygui** (`GUI.*`), **`DB.*`**, or ENet, use **`CGO_ENABLED=1`** and a C toolchain as above.

**Smoke test (purego only):** [`cmd/puregohello/`](../cmd/puregohello/) loads Raylib via [`internal/raylibpurego`](../internal/raylibpurego/) and moves a textured quad with the keyboard. Build with **`CGO_ENABLED=0`** and run with the same sidecar Raylib shared library as the main binary.

### 3. Build on Linux

Open a terminal and run the following commands:

```bash
# Set the CGO_ENABLED flag
export CGO_ENABLED=1

# Build the executable
go build -o moonbasic .
```

After a successful build, you can run the interpreter directly or add it to your system's PATH to run it from any directory.

### Distribution-style builds (full runtime)

Release archives ship **`moonbasic`** and **`moonrun`** built with **`-tags fullruntime`** so windowed programs, **`moonbasic --run`**, and the full builtin surface match what contributors test with **`go test -tags fullruntime ./...`**.

```bash
# Linux example
export CGO_ENABLED=1
go build -tags fullruntime -o moonbasic .
go build -tags fullruntime -o moonrun ./cmd/moonrun
```

On **Windows**, set **`CGO_ENABLED=1`** and point **`CC`** at MinGW **`gcc.exe`** as in [Build on Windows](#2-build-on-windows) above, then use the same **`-tags fullruntime`** lines (outputs **`moonbasic.exe`** / **`moonrun.exe`**).

**3D physics:** native **Jolt** (`PHYSICS3D.*` / `BODY3D.*`) is available on **Linux and Windows x64** when **`CGO_ENABLED=1`** and the Jolt static libraries are present (see [JOLT_WINDOWS_PARITY.md](JOLT_WINDOWS_PARITY.md)). Other builds get a **full graphics** runtime with physics builtins **stubbed** with a clear error—see [PHYSICS3D.md](reference/PHYSICS3D.md).

**Jolt on Windows (LTO / GCC mismatch):** If the link step fails with **LTO** or **GCC version** errors (e.g. vendored **`libJolt.a`** built with a different toolchain than your MinGW **`gcc`**), rebuild the static archives with the **same** compiler you use for CGO — see [`third_party/jolt-go/jolt/lib/windows_amd64/README.md`](../third_party/jolt-go/jolt/lib/windows_amd64/README.md). As a **temporary** experiment, you can try disabling LTO on the **final** link only, for example: `go build -ldflags="-extldflags=-fno-lto" -tags fullruntime ./cmd/moonrun`. That may not fix every mismatch; rebuilding the **`.a`** files to match your environment remains the reliable fix.

---

## Windows static-linked `moonrun` (no `raylib.dll` / `jolt.dll`)

For a **standalone `.exe`** where native code is linked statically (game content may still load from disk next to the binary):

1. Build Jolt static archives for Windows amd64 (see [`third_party/jolt-go/jolt/lib/windows_amd64/README.md`](../third_party/jolt-go/jolt/lib/windows_amd64/README.md)).
2. From the repo root, run:

```powershell
powershell -File scripts/check-jolt-windows-libs.ps1   # optional preflight
powershell -File scripts/build_static.ps1
```

Default output: **`moonrun_static.exe`**. Optional: `$env:OUTPUT="moonrun.exe"` before the script.

3. **Verify** the linker did not leave runtime DLLs for Raylib or Jolt:

```powershell
dumpbin /dependents moonrun_static.exe
```

You should **not** see **`raylib.dll`** or **`jolt.dll`**. Non-system DLLs (e.g. `VCRUNTIME140.dll`) may still appear depending on toolchain. If `dumpbin` is not on `PATH`, run the same command from a **Visual Studio Developer** shell.

**Purego note:** Builds with **`CGO_ENABLED=0`** on Windows load **`raylib.dll`** at runtime; they are **not** covered by the static script above.

---

## Beta zip distribution (exe + loose folders)

Shipping **scripts, shaders, and assets** as files next to the binary keeps rebuilds fast and matches common engine layouts. A helper script builds the static runner and packs a standard tree:

```powershell
powershell -File scripts/package_beta_zip.ps1
```

Default archive: **`dist/MoonBasic-beta-windows-amd64.zip`**, containing a **`MoonBasic/`** root with **`moonrun.exe`**, **`shaders/shd/`** (mirror of [`runtime/shaders/shd`](../runtime/shaders/shd)), **`assets/`**, **`examples/`**, and **`README-BETA.txt`**.

| Issue | What to check |
|-------|----------------|
| **File not found** | Restore the full zip layout; paths are often relative to the bundle folder or the `.exe` directory (see **`RES.PATH`** in scripts). |
| **Wrong working directory** | Run commands from the unzipped **`MoonBasic`** folder (or the folder that contains **`moonrun.exe`**) so relative paths in samples resolve. |
| **Missing DLL error** | You are likely running a **non-static** build (e.g. purego). Rebuild with **`scripts/build_static.ps1`** or use the packaged **`moonrun.exe`**. |

**Clean-room check:** On a PC **without** Go, Zig, or **`raylib.dll`** on `PATH`, unzip the archive, open a terminal in **`MoonBasic`**, run **`.\moonrun.exe examples\sphere_drop\main.mb`**. The window should open if GPU drivers are available.

**Future:** A single-file **`embed.FS`** bundle (atomic ship) is optional and deferred until physics/rendering are stable across the full matrix; when added, prefer **raw** embedded bytes (no startup decompression) for fast boot.

---

**See also:** [DEVELOPER.md](DEVELOPER.md) (command cheat sheet, compile vs run), [CONTRIBUTING.md](../CONTRIBUTING.md).
