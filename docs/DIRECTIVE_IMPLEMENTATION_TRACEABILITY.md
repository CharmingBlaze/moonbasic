# COMPILER_ENGINEER_DIRECTIVE — implementation traceability

This file maps [COMPILER_ENGINEER_DIRECTIVE.md](../COMPILER_ENGINEER_DIRECTIVE.md) requirements to concrete locations in the repo. Update when behavior changes.

Aligned with the directive: **Namespace.Method** as primary style; **`CREATE`** and **`SETPOS`** canonical; **`MAKE`** and **`SETPOSITION`** kept as deprecated aliases during the v0.9 transition; **Easy Mode** only as convenience mappings (see [EASY_MODE.md](EASY_MODE.md)).

| Directive area | Requirement | Primary implementation | Verification |
|----------------|-------------|------------------------|----------------|
| Part 1 — Naming | Canonical `*.CREATE*` | `compiler/builtinmanifest/commands.json` | `api_standardization_test.go` (`TestCreateHasMatchingMakeDeprecatedAlias`); `manifest_test.go` (`TestMakeAliasDeprecationMatchesCanonicalKeys`) |
| Part 1 — Naming | Deprecated `*.MAKE*` aliases | Same manifest + `description` deprecation text | Same test + semantic deprecation notices |
| Part 1 — v0.9 transition | Every manifest `*.CREATE` / `*.MAKE` has a runtime registration (CREATE primary, MAKE alias) | `runtime/**` (`Register`, stub name lists, `regDual` for `PARTICLE`→`PARTICLE3D`) | `runtime_registration_test.go` (`TestManifestCreateKeysAppearInRuntimeSources`, `TestManifestMakeKeysAppearInRuntimeSources`); `PARTICLE3D.*` checked via sibling `PARTICLE.*` literal |
| Part 1 — AUDIO manifest | No duplicate `(key, args)` overload rows | `compiler/builtinmanifest/commands.json` | `TestNoDuplicateAudioManifestArgSignatures` |
| Part 1 — Position | Canonical `SETPOS`, alias `SETPOSITION` | Manifest | `TestSetPosHasSetPositionAlias`; `TestSetPositionDeprecationMatchesSetPosKeys` |
| Part 2 — Universal methods | Handle `.pos`/`.rot`/`.scale`/…; **zero-arg** reads dispatch to `GET*` builtins where they exist (e.g. `CAMERA.GETROT`, `ENTITY.GETPOS` / `GETROT` / `GETSCALE`, `BODY2D.GETPOS` / `GETROT`, `LIGHT.GETPOS` / `LIGHT.GETDIR` / `LIGHT.GETCOLOR`, `SPRITE.GETPOS`, `PARTICLE.GETPOS` / `GETCOLOR` / `GETALPHA` + manifest `PARTICLE3D.GETPOS` / `GETCOLOR` / `GETALPHA`, `LIGHT2D.GETPOS`, `CAMERA2D.GETPOS` / `GETROTATION`, `NAVAGENT.GETPOS`, `INSTANCE.GETPOS` / `GETROT` / `GETSCALE` (instance **index 0**) | `vm/handlecall.go` (`handleCallDispatch` for 0-arg getters; `handleCallBuiltin` for setters); runtime per module; manifest rows for each `GET*` | `vm/handlecall_dispatch_test.go` (includes `BODYREF` kinematic/static/trigger → `BODYREF.GETPOSITION` / `GETROTATION`); `compiler/builtinmanifest` tests; fullruntime samples |
| Part 2 — Sprite collision | **`SPRITE.HIT`** / **`SPRITE.POINTHIT`** match **`SPRITE.DRAW`** (raylib **`DrawTexturePro`**): **scaled** frame size, **scaled origin**, **rotation** | `runtime/sprite/sprite_hit_cgo.go` (SAT + inverse local point test; corners match `rtextures.c`) | `runtime/sprite/sprite_hit_cgo_test.go` |
| Part 4 — ENTITY / WINDOW position | Canonical `ENTITY.SETPOS`, `WINDOW.SETPOS` | `runtime/mbentity/entity_cgo.go`, `runtime/window/window_placement_*.go`; manifest | Deprecation: `SETPOSITION` → `SETPOS` via `builtinmanifest.DeprecationReplacement` |
| Part 3 — Ergonomics | Chainable setters, ≤4 params where refactored | Namespace modules + manifest overloads | Integration `testdata/*.mb` |
| Part 4 — Namespace symmetry | CREATE/FREE, SET/GET pairs | Manifest + `runtime/*/register*.go` | `go run ./tools/apidoc`, `runtime_registration_test.go` (CREATE/MAKE literals), manual namespace docs |
| Part 5 — Easy Mode | Thin wrappers only | `runtime/blitzengine/*`, easy-mode registrations | `docs/EASY_MODE.md`; e.g. `LIGHTPOSITION` → `LIGHT.SETPOS` in `register.go`; `CreateCamera` → `CAMERA.CREATE` in `runtime/mbentity/entity_blitz_cgo.go` |
| Part 6 — Docs | STYLE_GUIDE, conventions, migration | `STYLE_GUIDE.md`, `docs/reference/API_CONVENTIONS.md`, `docs/MIGRATION_CREATE_FROM_MAKE.md` | Doc review |
| Part 6 — Reference tables | Prefer **CREATE** / **SETPOS** in Blitz/camera/image docs | `docs/reference/MODERN_BLITZ_COMMANDS.md`, `CAMERA_LIGHT_RENDER.md`, `IMAGE.md`, `BLITZ2025.md`, `INSTANCE.md`, `RAYCAST.md`; `ARCHITECTURE.md` §3D; teaching tables under `docs/reference/moonbasic-command-set/` (index: `README.md`) | Grep `MAKE` / `SETPOSITION` in `docs/reference/` (migration + directive samples may still show deprecated forms on purpose); command-set **Naming (registry)** in `moonbasic-command-set/README.md` |
| Part 3 — Rays / instancing | **`RAY.CREATE`**, **`MODEL.CREATEINSTANCED`**, **`INSTANCE.CREATE`**, **`INSTANCE.CREATEINSTANCED`** | `runtime/mbcollision/ray_cgo.go`, `runtime/mbmodel3d/model_inst_draw_cgo.go` | Manifest `RAY.CREATE` + deprecated `RAY.MAKE`; `go test ./...` |
| Part 3 — Particles / sky | **`PARTICLE.CREATE`**, **`SKY.CREATE`** | `runtime/mbparticles/particles_cgo.go`, `runtime/sky/register_cgo.go` (+ stubs); **`PARTICLE3D.CREATE`** via `regDual` | `TestMakeAliasDeprecationMatchesCanonicalKeys` (covers `PARTICLE.MAKE` / `SKY.MAKE`) |
| Part 3 — Mem / sprite extras | **`MEM.CREATE`**, **`SPRITELAYER.CREATE`**, **`SPRITEBATCH.CREATE`** | `runtime/mbmem/mem.go`, `runtime/sprite/extras_cgo.go` | Same (`MEM.MAKE`, `SPRITELAYER.MAKE`, `SPRITEBATCH.MAKE`) |
| Part 3 — Sim timers | **`TIMER.CREATE`** / **`TIMER.MAKE`** (deprecated) | `runtime/mbgame/timer_sim.go`, `runtime/mbgame/timers.go` (wall `TIMER.NEW`); manifest | [TIMER.md](../reference/TIMER.md); `TestMakeAliasDeprecationMatchesCanonicalKeys` (`TIMER.MAKE`) |
| Part 7 checklist — Phase 1 | Manifest + apidoc | `tools/apidoc/main.go` (intro: CREATE / deprecated `MAKE`, default-value table); `docs/API_CONSISTENCY.md` | `go run ./tools/apidoc` after manifest edits |
| Part 7 — Doc audit | Namespace → reference table | `tools/cmdaudit/main.go` → `docs/COMMAND_AUDIT.md` | `go run ./tools/cmdaudit` (must exit 0) |
| Part 7 — Phase 2 | Deprecation warnings + optional strict | `compiler/semantic/analyze.go` (`StrictDeprecated`), `compiler/pipeline/compile.go` (`CheckOptions`), `lsp/server.go` | `compiler/semantic/semantic_test.go`; CLI `--strict-deprecated` |
| Part 7 — Tooling | LSP completion order, snippets | `lsp/server.go`, `.vscode/moonbasic-snippets.code-snippets` | Manual editor test |

## Key files (quick index)

- Manifest: `compiler/builtinmanifest/commands.json`
- Manifest parsing + policy tests: `compiler/builtinmanifest/manifest.go`, `manifest_json.go`, `*_test.go` (includes `runtime_registration_test.go`: manifest vs `runtime/` quoted builtin names)
- Semantic analysis + deprecations: `compiler/semantic/analyze.go`, `deprecation.go`
- Compile/check pipeline: `compiler/pipeline/compile.go`
- LSP: `lsp/server.go`
- API doc generator: `tools/apidoc/main.go` → `docs/API_CONSISTENCY.md`
- Command-set index (CREATE / SETPOS): `docs/reference/moonbasic-command-set/README.md`
- Easy Mode: `runtime/blitzengine/`, `docs/EASY_MODE.md`
- Sprite hit tests (aligned with draw): `runtime/sprite/sprite_hit_cgo.go`
- Examples + scripts: `examples/**/*.mb`, `testdata/**/*.mb` (prefer **`CreateCamera()`** / **`*.CREATE`** / **`SetPos`** over deprecated **`*.MAKE`** / **`SetPosition`** where both exist; **`CLOUD`**, **`BIOME`** are in manifest + runtime with **`*.CREATE`** / deprecated **`*.MAKE`**)
