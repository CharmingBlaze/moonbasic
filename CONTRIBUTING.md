# Contributing to moonBASIC

Thank you for helping improve the compiler, VM, runtime, or docs.

## Prerequisites

- **Go** — version in [`go.mod`](go.mod) (currently 1.25.3+).
- **C toolchain** — required for the default **CGO** build (raylib, physics, networking). See [docs/BUILDING.md](docs/BUILDING.md) for Windows (MinGW) and Linux packages.

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
| **Compiler only** (default) | `go build -o moonbasic .` | `.mb` → `.mbc`, `--check`, `--lsp`, `--disasm`. No game window. |
| **Full runtime** | `go build -tags fullruntime -o moonrun ./cmd/moonrun` | Run graphical programs from `.mb` / `.mbc`. |

Alternatively: `go build -tags fullruntime -o moonbasic .` gives a single binary that can **`--run`** (see [docs/DEVELOPER.md](docs/DEVELOPER.md)).

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
2. Implement registration and behavior under [`runtime/`](runtime/) (and related packages).
3. Run `go run . --check` on a sample that exercises the change.
4. Regenerate API docs when the public surface changes: `go run ./tools/apidoc` (updates [`docs/API_CONSISTENCY.md`](docs/API_CONSISTENCY.md)).

**Ease-of-use helpers:** New pattern commands (movement, snapping, camera-relative input, etc.) should complement—not replace—existing `MATH.*` / vector primitives. Naming, tuples, and documentation expectations are summarized in [`docs/EASY_LANGUAGE.md`](docs/EASY_LANGUAGE.md).

## Architecture

High-level pipeline and layout: [ARCHITECTURE.md](ARCHITECTURE.md). Deeper contributor map: [docs/DEVELOPER.md](docs/DEVELOPER.md).

## Optional: command coverage

[`COMMAND_AUDIT.txt`](COMMAND_AUDIT.txt) tracks implementation status (`DONE`, `PARTIAL`, `MISSING`) for builtins—useful for larger features, not required for every small fix.

## First-Time Contributor's Checklist
When contributing to MoonBASIC, remember our Static-First philosophy to ensure single-binary Zero-DLL purity across releases!
- [ ] No `C` headers bridging shared dll calls natively (unless encapsulated via CGO `#cgo LDFLAGS` specifying static archives).
- [ ] Only utilize Pure-Go parsers (e.g. `qmuntal/gltf`) targeting GPU buffers statically.
- [ ] Make sure resources utilize `//go:embed` targeting payload bundles rather than enforcing loose paths avoiding runtime disk crashes natively.
- [ ] Remember to update the `commands.json` API manifest exactly aligning new methods directly with handle models.
