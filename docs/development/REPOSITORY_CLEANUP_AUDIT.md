# Repository cleanup audit

Inventory for structural cleanup of moonbasic-compiler.  
**Rule:** do not delete unclear/risky files. Prefer moves with `git mv` and path updates.

Date: 2026-07-25  
Branch: `chore/repository-cleanup`  
Baseline: `main` @ post-`8233017`

Legend — **Action:** keep | move | rename | generate | defer | ignore-local  
**Risk:** low | medium | high

---

## Top-level directories

### `compiler/`

| Field | Value |
|-------|--------|
| Current path | `compiler/` |
| Purpose | Lexer, parser, AST, semantic, codegen, include, builtinmanifest, pipeline |
| References | `main.go`, `cmd/moonbasic`, `lsp/`, `vm/`, tests, CI `--check` |
| Recommended destination | Keep |
| Recommended action | keep |
| Risk level | low |
| Reason | Clear architectural boundary; already well named |
| Validation required | `CGO_ENABLED=0 go test ./compiler/...` |

### `vm/`

| Field | Value |
|-------|--------|
| Current path | `vm/` |
| Purpose | Opcode IR, execution, heap, handlecall bridge |
| References | compiler codegen, runtime, `moonrun` |
| Recommended destination | Keep |
| Recommended action | keep |
| Risk level | low |
| Reason | Bytecode contract lives here |
| Validation required | `go test ./vm/...` |

### `runtime/`

| Field | Value |
|-------|--------|
| Current path | `runtime/` (~65 packages, many `mb*` prefixes) |
| Purpose | Built-in command implementations (graphics, physics, net, data, world, …) |
| References | `compiler/pipeline/registry.go` (fullruntime), CI, examples |
| Recommended destination | Keep flat for this cleanup series |
| Recommended action | defer domain regroup (`graphics/`, `physics/`, …) |
| Risk level | high (if moved now) |
| Reason | Build tags, CGO/stub splits, registration order; mass moves break CI |
| Validation required | Fullruntime + stub-only matrices after any future move |

**Deferred domain map (future only):**

- Core: `system`, `time`, `mathmod`, `strmod`, `mbarray`, `mbmem`, `file`, …
- Graphics: `window`, `draw`, `texture`, `mbimage`, `font`, `camera`, `mbmodel3d`, `mblight*`, `shaders`, `mbparticles`, …
- Physics: `physics2d`, `physics3d`, `mbcollision`, `charcontroller`, `player`
- World: `terrain`, `sky`, `water`, `weathermod`, `cloudmod`, `biome`, `scatter`, `worldmgr`
- Gameplay: `mbentity`, `mbscene`, `mbtween`, `mbevent`, `mbcoroutine`, `mbnav`, `mbsave`, …
- Data: `jsonmod`, `csvmod`, `dbmod`, `tablemod`
- Compatibility: `blitzengine`, `checklist_aliases`

### `hal/` + `drivers/`

| Field | Value |
|-------|--------|
| Current path | `hal/`, `drivers/` |
| Purpose | Hardware abstraction; null vs Raylib video drivers |
| References | runtime window/input, docs/architecture/HAL_AND_RENDERING.md |
| Recommended destination | Keep |
| Recommended action | keep |
| Risk level | low |
| Reason | Already the preferred new-code boundary |
| Validation required | Headless registry + fullruntime build |

### `handlecall/`

| Field | Value |
|-------|--------|
| Current path | `handlecall/` |
| Purpose | Shared typed handle method dispatch (semantic + VM) |
| References | `compiler/semantic`, `vm/handlecall.go` |
| Recommended destination | Keep (or later `compiler/handlecall` — defer) |
| Recommended action | keep |
| Risk level | low |
| Reason | New shared package; moving now adds churn |
| Validation required | `go test ./handlecall/... ./compiler/semantic/...` |

### `ide/` → `ide/`

| Field | Value |
|-------|--------|
| Current path | `ide/` (space in name) |
| Purpose | Wails v2 desktop IDE |
| References | `.github/workflows/release.yml`, `scripts/package_ide_bundle.*`, `tools/docsexport`, `tools/ideexport`, internal toolchain paths |
| Recommended destination | `ide/` |
| Recommended action | rename |
| Risk level | medium |
| Reason | Space breaks tooling/scripts; many path strings |
| Validation required | `npm run sync` in ide; `go run ./tools/docsexport`; release packaging scripts |

### `ide/bundled-docs/`

| Field | Value |
|-------|--------|
| Current path | `ide/bundled-docs/` (→ `ide/bundled-docs/`) |
| Purpose | Generated copy of `docs/` for offline IDE |
| References | `tools/docsexport`, IDE `docs.go` |
| Recommended destination | Keep under `ide/` |
| Recommended action | generate (canonical source: `docs/`) |
| Risk level | low |
| Reason | Must stay in sync via tool, not hand-edited |
| Validation required | `go run ./tools/docsexport --check` |

### `editors/vscode-moonbasic/`

| Field | Value |
|-------|--------|
| Current path | `editors/vscode-moonbasic/` |
| Purpose | VS Code / Cursor extension |
| References | `moonbasic install-vscode`, packaging, CI vscode job |
| Recommended destination | Keep |
| Recommended action | keep |
| Risk level | medium |
| Reason | Large tree; do not commit `node_modules` (verify gitignore) |
| Validation required | `npm ci && npm run package` in extension dir |

### `cmd/`

| Field | Value |
|-------|--------|
| Current path | `cmd/moonbasic`, `cmd/moonrun`, `cmd/moondoc`, `cmd/puregohello` |
| Purpose | User-facing tools |
| References | releases, docs, CONTRIBUTING |
| Recommended destination | Keep; IDE stays Wails module under `ide/` (not forced into `cmd/`) |
| Recommended action | keep |
| Risk level | low |
| Reason | Standard Go layout |
| Validation required | `go build ./cmd/...` |

### Root Go entrypoints (`main.go`, `main_*.go`, …)

| Field | Value |
|-------|--------|
| Current path | repo root |
| Purpose | Default `go run .` / `go build .` compiler CLI + fullruntime tags |
| References | CONTRIBUTING, AGENTS, CI `go run . --check` |
| Recommended destination | Keep at root |
| Recommended action | keep (defer move to `cmd/`) |
| Risk level | high if moved |
| Reason | Muscle memory + gopls split-brain docs |
| Validation required | `go build .` and `go build -tags fullruntime .` |

### `docs/`

| Field | Value |
|-------|--------|
| Current path | `docs/` |
| Purpose | Canonical documentation |
| References | IDE export, website, GETTING_STARTED |
| Recommended destination | Keep; add `docs/development/` for contributor meta |
| Recommended action | keep + organize development docs |
| Risk level | low |
| Reason | Already the single source of truth for content |
| Validation required | Link check after moves; docsexport |

### `scripts/`

| Field | Value |
|-------|--------|
| Current path | `scripts/` (flat) |
| Purpose | Build, release, packaging, verification, dev helpers |
| References | CI workflows, Makefile, BUILDING.md, DEVELOPER.md |
| Recommended destination | `scripts/{build,release,packaging,verification,development}/` |
| Recommended action | move |
| Risk level | medium |
| Reason | Many hard-coded paths in YAML/docs |
| Validation required | CI path grep + `check_builds` |

### `tools/`

| Field | Value |
|-------|--------|
| Current path | `tools/` |
| Purpose | apidoc, docsexport, ideexport, grammargen, manifest audits |
| References | CI, npm langdata, docs |
| Recommended destination | Keep; document categories in README |
| Recommended action | keep |
| Risk level | low |
| Reason | Already small and clear |
| Validation required | `go run ./tools/apidoc`, docsexport |

### `packaging/`

| Field | Value |
|-------|--------|
| Current path | `packaging/` |
| Purpose | START-IDE, README-*-RELEASE, samples, ADD-TO-PATH |
| References | release.yml, package_ide_bundle |
| Recommended destination | Keep |
| Recommended action | keep |
| Risk level | low |
| Reason | Release-facing; path stability matters |
| Validation required | stage_ide_extras / package scripts |

### `testdata/`

| Field | Value |
|-------|--------|
| Current path | `testdata/` |
| Purpose | Compiler/runtime fixtures, CI `--check` samples |
| References | CI, unit tests |
| Recommended destination | Keep flat for now |
| Recommended action | defer taxonomy split |
| Risk level | medium if reshuffled |
| Reason | Many relative paths in tests/CI |
| Validation required | Full CI check list |

### `examples/`

| Field | Value |
|-------|--------|
| Current path | `examples/` |
| Purpose | Showcase games/scripts |
| References | CI, README, packaging |
| Recommended destination | Keep |
| Recommended action | defer beginner/2d/3d regroup |
| Risk level | medium if reshuffled |
| Reason | Docs and CI cite concrete paths |
| Validation required | `--check` on cited examples |

### `third_party/`

| Field | Value |
|-------|--------|
| Current path | `third_party/jolt-go`, `go-enet`, `raylib-go-raylib` |
| Purpose | Vendored native / bindings |
| References | go.mod replace, CGO builds, PHYSICS.md |
| Recommended destination | Keep |
| Recommended action | keep |
| Risk level | high if altered |
| Reason | Licences + static libs + amalgamation |
| Validation required | fullruntime builds |

### `benchmarks/`, `assets/`, `web/`, `dist/`, `dap/`, `lsp/`, `internal/`, `registry/`, `tests/`, `lineprof/`

| Path | Purpose | Action | Risk |
|------|---------|--------|------|
| `benchmarks/` | Perf `.mb` samples | keep | low |
| `assets/` | Shared sample assets | keep | low |
| `web/` | Web-related stubs/site helpers | keep | low |
| `dist/` | Local/release staging + README | keep (gitignore archives as needed) | low |
| `dap/` | Debug adapter | keep | low |
| `lsp/` | Language server | keep | low |
| `internal/` | Internal shared Go | keep | low |
| `registry/` | Small registry helper | keep | low |
| `tests/` | Sparse integration | keep | low |
| `lineprof/` | Profiling helper | keep | low |

### `.github/`

| Field | Value |
|-------|--------|
| Current path | `.github/workflows/` |
| Purpose | CI / release |
| Recommended action | update paths when scripts/IDE move |
| Risk level | medium |
| Validation required | PR CI green |

---

## Root markdown files

| Current path | Purpose | Recommended destination | Action | Risk |
|--------------|---------|-------------------------|--------|------|
| `README.md` | Project front door | root | keep | low |
| `LICENSE` | Licence | root | keep | low |
| `CONTRIBUTING.md` | Contributor entry | root | keep | low |
| `AGENTS.md` | Agent/contributor notes | root | keep | low |
| `ARCHITECTURE.md` | High-level architecture | root | keep | low |
| `STYLE_GUIDE.md` | Language style | root | keep | low |
| `COMPILER_ENGINEER_DIRECTIVE.md` | Compiler engineering directive | `docs/development/` | move | low |
| `ENGINE_IR_V2.md` | IR v2 notes | `docs/architecture/` | move | low |
| `ENGINE_IR_V3.md` | IR v3 notes | `docs/architecture/` | move | low |
| `Masterplan.md` | Product/engineering masterplan | `docs/development/` | move | low |
| `final polish and docs.md` | Scratch polish notes | `docs/audit/archive/` | move | low |

---

## Local scratch (do not commit)

| Path | Notes |
|------|--------|
| `*.exe`, `moonbasic.exe`, `moonrun.exe` | `.gitignore` already covers `*.exe` |
| `save.json`, `spin_test.txt`, `_mattest/` | Already ignored |

---

## Duplication

| Item | Canonical | Copy | Action |
|------|-----------|------|--------|
| User docs | `docs/` | `ide/bundled-docs/` (after rename) | generate via `tools/docsexport` |
| API surface | `compiler/builtinmanifest/commands.json` | `docs/api/*`, `docs/API_CONSISTENCY.md` | generate via `tools/apidoc` |
| Command lists | manifest | various audits | prefer generators |

---

## Explicitly deferred (do not do in this series)

1. Regroup `runtime/` into domain folders  
2. Mass `testdata/` / `examples/` taxonomy  
3. Rename `mb*` Go packages  
4. Move root `main.go` into `cmd/`  
5. Delete historical docs without link+history proof  

---

## Planned safe actions (this branch)

1. Add this audit, cleanup log, codebase map  
2. Move root planning docs into `docs/`  
3. Reorganize `scripts/` into subfolders + update references  
4. Rename `ide` → `ide`  
5. Add `docsexport --check` + generated marker  
6. Package READMEs + CONTRIBUTING link to codebase map  
