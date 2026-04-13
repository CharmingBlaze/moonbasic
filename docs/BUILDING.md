# Building moonBASIC from Source

This guide provides detailed instructions for compiling the `moonBASIC` interpreter from its source code.

---

## Dependencies

Before you can build, you need the following software installed on your system.

### All Systems
- **Go**: Version 1.22 or later. You can download it from the [official Go website](https://go.dev/dl/).
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

**3D physics:** native **Jolt** (`PHYSICS3D.*` / `BODY3D.*`) is **Linux + CGO** in this repository. Other platforms still get a **full graphics** runtime, but those builtins **fail with a clear stub error**—see [PHYSICS3D.md](reference/PHYSICS3D.md).

---

**See also:** [DEVELOPER.md](DEVELOPER.md) (command cheat sheet, compile vs run), [CONTRIBUTING.md](../CONTRIBUTING.md).
