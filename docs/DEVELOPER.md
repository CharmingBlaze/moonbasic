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

## Build tags: `fullruntime` vs default

- **Default** (`go build .`, `go run .`): builds [`main.go`](../main.go) — **compiler only** (no linked game runtime). Running `go run . game.mb` **writes `game.mbc`** next to the source; it does **not** open a window.
- **`-tags fullruntime`**: includes [`main_fullruntime.go`](../main_fullruntime.go) instead, links the full runtime. Use **`moonrun`**, or **`go run -tags fullruntime . --run file.mb`**, to execute graphical programs.

Details: [BUILDING.md](BUILDING.md).

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
