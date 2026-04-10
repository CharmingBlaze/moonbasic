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

## Build tags: `fullruntime` vs default

- **Default** (`go build .`, `go run .`): builds [`main.go`](../main.go) — **compiler only** (no linked game runtime). Running `go run . game.mb` **writes `game.mbc`** next to the source; it does **not** open a window.
- **`-tags fullruntime`**: includes [`main_fullruntime.go`](../main_fullruntime.go) instead, links the full runtime. Use **`moonrun`**, or **`go run -tags fullruntime . --run file.mb`**, to execute graphical programs.

Details: [BUILDING.md](BUILDING.md).

## Developer environment: VS Code, gopls, and “split brain”

The repo uses **mutually exclusive** `//go:build` lines at the roots of the main binaries:

| Tags | Root | `cmd/moonbasic` | `cmd/moonrun` |
|------|------|-----------------|---------------|
| **Default** (no `fullruntime`) | [`main.go`](../main.go) | yes (`!fullruntime`) | excluded |
| **`fullruntime`** | [`main_fullruntime.go`](../main_fullruntime.go) | excluded | yes |

**gopls** runs `go list` with a **single** set of build tags. It cannot load both sides of that split at once, which is why you may see **“No packages found”** or files **greyed out** when the active tags do not match the file you opened.

### Default IDE setup (fullruntime)

For day-to-day work on **physics, rendering, VM + runtime modules, and `cmd/moonrun`**, the workspace [`.vscode/settings.json`](../.vscode/settings.json) sets:

```json
"gopls": { "buildFlags": ["-tags=fullruntime"], ... }
```

That enables analysis for [`main_fullruntime.go`](../main_fullruntime.go), [`cmd/moonrun/`](../cmd/moonrun/), and the full CGO / Jolt / Raylib graph.

### Switching to “compiler CLI” mode

To edit [`main.go`](../main.go) or [`cmd/moonbasic/`](../cmd/moonbasic/) with full IntelliSense:

1. Open [`.vscode/settings.json`](../.vscode/settings.json).
2. Remove or comment out the `"buildFlags": ["-tags=fullruntime"]` entry inside `gopls`.
3. Run **Go: Restart Language Server** (Command Palette) or reload the window.

Switch back when you return to runtime-heavy code.

### Why this exists

- The **default** toolchain stays **small** (compiler, LSP, `--check`): suitable for “zero extra DLL” compiler builds and fast iteration.
- **`fullruntime`** pulls in the **heavy** game stack (Raylib, optional Jolt on Linux, etc.) for **`moonrun`** and `--run`.

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
| Type-check | `go run . --check path/to/script.mb` |
| Compile to `.mbc` | `go run . path/to/script.mb` |
| Run game (source) | `CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun path/to/script.mb` |
| Run game (alternate) | `CGO_ENABLED=1 go run -tags fullruntime . --run path/to/script.mb` |
| Disassemble bytecode | `go run . --disasm path/to/script.mbc` |
| All Go tests | `go test ./...` |
| Regenerate API consistency doc | `go run ./tools/apidoc` |

Shortcuts: see [Makefile](../Makefile) (Unix/Git Bash) and [scripts/dev.ps1](../scripts/dev.ps1) / [scripts/dev.sh](../scripts/dev.sh).

## CI parity

Linux CI runs `go test ./...` and a set of `go run . --check …` commands on representative samples. See [.github/workflows/ci.yml](../.github/workflows/ci.yml). Running the same `--check` lines locally before pushing catches most compile-time regressions.

## Editing the command manifest

[`compiler/builtinmanifest/commands.json`](../compiler/builtinmanifest/commands.json) is the source of truth for **names, arity, and types** exposed to MoonBASIC. After edits, validate with `--check` on real scripts and refresh [API_CONSISTENCY.md](API_CONSISTENCY.md) via `go run ./tools/apidoc` when user-visible API changes.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).
