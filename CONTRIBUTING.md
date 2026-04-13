# Contributing to moonBASIC

Thank you for helping improve the compiler, VM, runtime, or docs.

## Prerequisites

- **Go** — version in [`go.mod`](go.mod) (currently 1.25.3+).
- **C toolchain** — required for the default **CGO** build (raylib, physics, networking). See [docs/BUILDING.md](docs/BUILDING.md) for Windows (MinGW) and Linux packages.

## Platform priority (Windows, then Linux)

**Windows first, Linux second:** The project assumes most contributors run **Windows** for the default **fullruntime** + **CGO** + **Raylib** loop (`moonrun`, `--check`). **Linux** is the follow-on target for **full Jolt** (KCC, rigid-body **`PHYSICS3D`**) and for running **`bash scripts/check_builds.sh`** in CI. When you document OS-specific behavior, list **Windows** before **Linux** in tables and prose. Details: [docs/DEVELOPER.md](docs/DEVELOPER.md#platform-priority-windows-then-linux).

## Clone and verify

From the repository root:

```bash
go test ./...
go run . --check examples/mario64/main_entities.mb
```

The [`--check`](.github/workflows/ci.yml) samples mirror a subset of CI; fixing failures before opening a PR saves round-trips.

## Two build modes

| Build | Command | What you get |
|--------|---------|----------------|
| **Headless Compiler** (default) | `go build -o moonbasic .` | `.mb` → `.mbc`, `--check`, `--lsp`, `--disasm`. Uses **Null Driver** (dependency-free). |
| **Full Interactive Runtime** | `go build -tags fullruntime -o moonrun ./cmd/moonrun` | Run graphical programs from `.mb` / `.mbc`. Uses **Raylib Driver**. |

Alternatively: `go build -tags fullruntime -o moonbasic .` gives a single binary that can **`--run`** locally. To produce a **standalone static exe** on Windows, use [`scripts/build_static.ps1`](scripts/build_static.ps1).

### IDE: gopls and build tags (“split brain”)

VS Code / Cursor is configured (see [`.vscode/settings.json`](.vscode/settings.json)) with **`gopls` `buildFlags`: `-tags=fullruntime,gopls_stub`** so IntelliSense covers the game runtime, [`main_fullruntime.go`](main_fullruntime.go), and [`cmd/moonrun/`](cmd/moonrun/), and **`runtime/terrain`** stub files analyze on **Windows** (see **`docs/DEVELOPER.md`**). That **excludes** the default compiler entrypoints [`main.go`](main.go) and [`cmd/moonbasic/`](cmd/moonbasic/); you may see **“No packages found”** for those until you adjust tags and **restart the Go language server**. Full rationale and switching steps: **[docs/DEVELOPER.md](docs/DEVELOPER.md#developer-environment-vs-code-gopls-and-split-brain)**.

Before pushing Go changes, run **`bash scripts/check_builds.sh`** (or **`make check-builds`**, or **`powershell -File scripts/check_builds.ps1`** on Windows) to compile **both** the default compiler path and **`-tags fullruntime`** (`moonrun` + full root). On Windows, if PowerShell hits **`runtime/cgo`** errors on the fullruntime steps, use **Git Bash / MSYS2** and `bash scripts/check_builds.sh` (see [docs/DEVELOPER.md](docs/DEVELOPER.md#pre-push-validate-both-build-paths)).

**`ENTITY` spatial macros:** literal entity indices are range-checked in **semantic analysis** (visible to **`--check`**) and **codegen**; see [docs/COMPILER_SPEC.md](docs/COMPILER_SPEC.md). Regression: `go run . --check testdata/entity_spatial_id_oob.mb` must **fail** with a type error; other scripts in **`testdata/`** used by CI should still pass **`--check`** as before.

**Memory note (`MaxEntitySpatialIndex` = 2²⁴):** This value is an **upper bound on numeric ids** the compiler and VM accept for **`ENTITY.X` / `Y` / `Z` / …** macros—not a preallocated heap. The host **SoA** (`runtime.SpatialBuffer`: six **`float32`** columns) starts at a **small capacity** and **grows on demand** when entities are created (see **`runtime/mbentity/entity_cgo.go`**). A *theoretical* full 2²⁴ rows would be on the order of **hundreds of MiB** for those six slices alone (plus separate entity structs and engine data); normal games stay far below that. There is **no** separate user-facing **`MAX_ENTITIES`** config yet—tightening caps for mobile/embedded would be a future engine option.

## Typical workflows

| Goal | Command |
|------|---------|
| Type-check a script | `go run . --check path/to/game.mb` |
| Compile to bytecode | `go run . path/to/game.mb` → writes `path/to/game.mbc` |
| Run a game (window) | `go run -tags fullruntime ./cmd/moonrun path/to/game.mb` or use a built `moonrun` binary |
| Language server | `go run . --lsp` (stdio) |

**Important:** Plain `go run . file.mb` (without `-tags fullruntime`) only **compiles** to `.mbc`; it does not open a window. See [docs/DEVELOPER.md](docs/DEVELOPER.md).

## Changing builtins / commands

1. Add or update the declaration in [`compiler/builtinmanifest/commands.json`](compiler/builtinmanifest/commands.json).
2. Implement behavior in a `runtime/` package using **`hal` types** (`hal.V3`, etc.).
3. If the feature interacts with hardware, call via `rt.Driver`.
4. Ensure the feature satisfies the **Null Driver** (even as a no-op) so compiler tests remain headless.
5. Run `go run . --check` on a sample that exercises the change.
6. Regenerate API docs: `go run ./tools/apidoc` (**[`docs/API_CONSISTENCY.md`](docs/API_CONSISTENCY.md)**).

**Ease-of-use helpers:** New pattern commands (movement, snapping, camera-relative input, etc.) should complement—not replace—existing `MATH.*` / vector primitives. Naming, tuples, and documentation expectations are summarized in [`docs/EASY_LANGUAGE.md`](docs/EASY_LANGUAGE.md).

## Architecture

High-level pipeline and layout: [ARCHITECTURE.md](ARCHITECTURE.md). Deeper contributor map: [docs/DEVELOPER.md](docs/DEVELOPER.md).

## Optional: command coverage

[`COMMAND_AUDIT.txt`](COMMAND_AUDIT.txt) tracks implementation status (`DONE`, `PARTIAL`, `MISSING`) for builtins—useful for larger features, not required for every small fix.

## First-Time Contributor's Checklist
When contributing to MoonBASIC, remember our Static-First philosophy to ensure single-binary Zero-DLL purity across releases!
- [ ] **HAL Compliance**: All hardware access must go through `rt.Driver`. Direct imports of Raylib in `runtime/` are forbidden.
- [ ] **Static Linking**: Specify static archives in CGO LDFLAGS (already configured in `drivers/video/raylib`).
- [ ] **Headless Parity**: Ensure the `Null` driver is updated so `--check` and tests stay dependency-free.
- [ ] **Embed Resources**: Utilize `//go:embed` targeting payload bundles rather than enforcing loose paths.
- [ ] **Manifest Alignment**: Update `commands.json` API manifest exactly aligning new methods directly with handle types.
