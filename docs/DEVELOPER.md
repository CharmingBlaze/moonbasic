# Developer guide

Quick orientation for people and tools working on the **moonBASIC** repository (Go compiler + VM + game runtime).

## Repository layout

| Area | Role |
|------|------|
| [`compiler/`](../compiler/) | Lexer, parser, AST, semantic analysis, optimizer, codegen → bytecode |
| [`vm/`](../vm/) | Opcode definitions, MOON container (`.mbc`), VM execution |
| [`runtime/`](../runtime/) | Native modules: raylib, physics, audio, entities, etc. |
| [`lsp/`](../lsp/) | Language server for editors |
| [`cli_compile.go`](../cli_compile.go) | Shared “compile `.mb` → `.mbc`” helper used by CLI entrypoints |

Design overview: [ARCHITECTURE.md](../ARCHITECTURE.md).
Entity system refactor: [ARCHITECTURE_MODULAR_ENTITIES.md](ARCHITECTURE_MODULAR_ENTITIES.md).

## Platform priority (Windows, then Linux)

**Policy:** Treat **Windows** as the **first** platform for **day-to-day development**, **default tooling**, and **how reference docs are ordered**. Treat **Linux** as the **second** platform for **full Jolt** behavior (native **`PHYSICS3D`** / **`CharacterVirtual`**, rigid bodies, picks) and for **Unix-style CI** (`bash scripts/check_builds.sh`).

- **Docs:** In tables and bullet lists that compare OSes, put **Windows** before **Linux** unless the page is explicitly Linux-only (e.g. Jolt implementation notes).
- **Code:** Still maintain **both** paths (`*_cgo.go` / `*_stub.go`, same manifest keys); see [CONTRIBUTING.md](../CONTRIBUTING.md) and [AGENTS.md](../AGENTS.md).

## Build tags: `fullruntime` vs default (Headless)

- **Default** (`go build .`, `go run .`): builds [`main.go`](../main.go) — **headless compiler** with a **Null** hardware driver. Running `go run . game.mb` **writes `game.mbc`** and validates semantics without needing `raylib.dll` or a GPU. Suitable for CI/CD and servers.
- **`-tags fullruntime`**: includes [`main_fullruntime.go`](../main_fullruntime.go), links the **Raylib** hardware driver. Use **`moonrun`**, or **`go run -tags fullruntime . --run file.mb`**, to execute graphical programs.

Details: [BUILDING.md](BUILDING.md). **HAL / drivers / Windows purego vs CGO:** [architecture/HAL_AND_RENDERING.md](architecture/HAL_AND_RENDERING.md).

## Developer environment: VS Code, gopls, and “split brain”

The repo uses **mutually exclusive** `//go:build` lines at the roots of the main binaries:

| Tags | Root | `cmd/moonbasic` | `cmd/moonrun` |
|------|------|-----------------|---------------|
| **Default** (no `fullruntime`) | [`main.go`](../main.go) | yes (`!fullruntime`) | excluded |
| **`fullruntime`** | [`main_fullruntime.go`](../main_fullruntime.go) | excluded | yes |

**gopls** runs `go list` with a **single** set of build tags. It cannot load both sides of that split at once, which is why you may see **“No packages found”** or files **greyed out** when the active tags do not match the file you opened.

### Default IDE setup (fullruntime)

For day-to-day work on **physics, rendering, VM + runtime modules, and `cmd/moonrun`**, the workspace [`.vscode/settings.json`](../.vscode/settings.json) sets **`go.buildTags`** / **`gopls.buildFlags`** to **`fullruntime,gopls_stub`**, plus **`CGO_ENABLED=1`** (**`gopls.build.env`** + **`go.toolsEnvVars`**). The **`gopls_stub`** tag is **for gopls only** on **Windows**: it includes **`runtime/terrain/*_stub.go`** (e.g. **`heap_objects_stub.go`**) in the analysis build so you do not get **“No packages found”** when opening them. **`go build -tags fullruntime`** (no **`gopls_stub`**) is unchanged and still uses the real CGO terrain sources.

**`third_party/raylib-go-raylib` purego files** (e.g. **`raylib_purego.go`**, **`frustum_cull_purego_windows.go`**) use **`//go:build !cgo && windows`**. With the default **`CGO_ENABLED=1`**, gopls builds the CGO variant of **`raylib`** and those files are **out of the build**, so the editor may show **“No packages found”** when they are focused. To get IntelliSense while editing them, temporarily set **`gopls.build.env.CGO_ENABLED`** (and **`go.toolsEnvVars`**) to **`0`**, run **Go: Restart Language Server**, then switch back to **`1`** when returning to CGO-heavy work.

### Switching to “compiler CLI” mode

To edit [`main.go`](../main.go) or [`cmd/moonbasic/`](../cmd/moonbasic/) with full IntelliSense:

1. Open [`.vscode/settings.json`](../.vscode/settings.json).
2. Remove or comment out the `"buildFlags": ["-tags=fullruntime"]` entry inside `gopls`.
3. Run **Go: Restart Language Server** (Command Palette) or reload the window.

Switch back when you return to runtime-heavy code.

### Why this exists

- **Default (Headless)**: The toolchain stays **small** and **dependency-free** (compiler, LSP, `--check`). It uses a `Null` hardware backend, making it suitable for servers and fast unit testing.
- **`fullruntime`**: Pulls in the **heavy** hardware stack (Raylib, Jolt) via the `hal` package for **interactive** use.

### Pre-push: validate both build paths

After touching shared packages, confirm **both** tag axes still compile (avoids leaking imports across the boundary):

```bash
# Unix / Git Bash / WSL
bash scripts/check_builds.sh
# or
make check-builds
```

On Windows PowerShell:

```powershell
powershell -File scripts/check_builds.ps1
```

If **fullruntime** steps fail with **`runtime/cgo`** in plain PowerShell, run the same script from **Git Bash** or **MSYS2 MINGW64** so **MinGW `gcc`** is on `PATH` (same as [scripts/release-windows.sh](../scripts/release-windows.sh)):

```bash
bash scripts/check_builds.sh
```

The full-runtime steps expect **`CGO_ENABLED=1`** and a C toolchain (see [BUILDING.md](BUILDING.md)), same as a normal `moonrun` build.

## Command cheat sheet (repo root)

Replace paths as needed. On Windows, set `CGO_ENABLED=1` and `CC` per BUILDING.md when building the full runtime.

| Action | Command |
|--------|---------|
| Type-check (Headless) | `go run . --check path/to/script.mb` |
| Compile to `.mbc` (Headless) | `go run . path/to/script.mb` |
| Run game (source) | `CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun path/to/script.mb` |
| Run game (alternate) | `CGO_ENABLED=1 go run -tags fullruntime . --run path/to/script.mb` |
| Static Build (Windows) | `powershell -File scripts/build_static.ps1` |
| Disassemble bytecode | `go run . --disasm path/to/script.mbc` |
| All Go tests | `go test ./...` |
| Regenerate API consistency doc | `go run ./tools/apidoc` |

Shortcuts: see [Makefile](../Makefile) (Unix/Git Bash) and [scripts/dev.ps1](../scripts/dev.ps1) / [scripts/dev.sh](../scripts/dev.sh).

## CI parity

Linux CI runs `go test ./...` and a set of `go run . --check …` commands on representative samples. See [.github/workflows/ci.yml](../.github/workflows/ci.yml). Running the same `--check` lines locally before pushing catches most compile-time regressions.

## Editing the command manifest

[`compiler/builtinmanifest/commands.json`](../compiler/builtinmanifest/commands.json) is the source of truth for **names, arity, and types** exposed to MoonBASIC. After edits, validate with `--check` on real scripts and refresh [API_CONSISTENCY.md](API_CONSISTENCY.md) via `go run ./tools/apidoc` when user-visible API changes.

## Rendering stability and defaults

Stability is a priority. To prevent viewport masking (e.g. the "Black Screen" issue), global draw hooks and default lighting states must adhere to the stability guidelines.

Details: [RENDERING_STABILITY_AND_DEFAULTS.md](architecture/RENDERING_STABILITY_AND_DEFAULTS.md).

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).
