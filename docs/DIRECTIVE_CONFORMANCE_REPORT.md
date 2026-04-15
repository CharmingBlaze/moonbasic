# Full directive rollout — conformance report

This report closes the **Full Directive Implementation Plan** checklist against [COMPILER_ENGINEER_DIRECTIVE.md](../COMPILER_ENGINEER_DIRECTIVE.md).

**Last validation:** `go test ./...`, `go run ./tools/apidoc` (regenerates [API_CONSISTENCY.md](API_CONSISTENCY.md)), and `go run ./tools/cmdaudit` (regenerates [COMMAND_AUDIT.md](COMMAND_AUDIT.md)).

**Reference doc alignment (ongoing):** Blitz/camera/light/image overview pages list **CREATE** / **SETPOS** as primary and **MAKE** / **SETPOSITION** as deprecated where those files were updated; [MIGRATION_CREATE_FROM_MAKE.md](MIGRATION_CREATE_FROM_MAKE.md) and [COMPILER_ENGINEER_DIRECTIVE.md](../COMPILER_ENGINEER_DIRECTIVE.md) intentionally retain **before** examples.

## Completed deliverables

Directive [COMPILER_ENGINEER_DIRECTIVE.md](../COMPILER_ENGINEER_DIRECTIVE.md) **Part 7 checklist** is marked complete for **v0.9** (in-repo rollout). Recent API work includes **`CAMERA.GETROT`** + zero-arg `cam.rot()` → `CAMERA.GETROT` (`runtime/camera/camera_extras_cgo.go`, `vm/handlecall.go`); **sprite** stored **rot / scale / color / alpha** with universal handle getters; **`SPRITE.HIT`** / **`SPRITE.POINTHIT`** aligned with **`DrawTexturePro`** (`runtime/sprite/sprite_hit_cgo.go`).

| Plan area | Evidence |
|-----------|----------|
| Traceability + baseline | [DIRECTIVE_IMPLEMENTATION_TRACEABILITY.md](DIRECTIVE_IMPLEMENTATION_TRACEABILITY.md), [DIRECTIVE_BASELINE_REPORT.md](DIRECTIVE_BASELINE_REPORT.md) |
| Manifest CREATE/MAKE + SETPOS/SETPOSITION policy | `compiler/builtinmanifest/api_standardization_test.go`; manifest in `compiler/builtinmanifest/commands.json` |
| Compiler + LSP deprecations | `compiler/semantic/analyze.go`, `compiler/pipeline/compile.go` (`CheckOptions`, `CheckFileWithOptions`); `lsp/server.go` diagnostics + completion ordering |
| Strict deprecated mode | `--strict-deprecated` on `moonbasic --check` (`main.go`, `main_fullruntime.go`, `cmd/moonbasic/main.go`) |
| Universal handle documentation | [reference/UNIVERSAL_HANDLE_METHODS.md](reference/UNIVERSAL_HANDLE_METHODS.md) → `vm/handlecall.go` |
| Builder / long-arg guidance | [STYLE_GUIDE.md](../STYLE_GUIDE.md) (long-argument + builder section) |
| Easy Mode → canonical CREATE | `CreateCamera` → `CAMERA.CREATE` in `runtime/blitzengine/register.go` |
| API docs | `go run ./tools/apidoc` → [API_CONSISTENCY.md](API_CONSISTENCY.md) |
| Tooling | `lsp/server.go` completion sort; `.vscode/moonbasic-snippets.code-snippets` |
| Examples | `CAMERA.CREATE` (not `MAKE`); `CAMERA.SETPOS`, `ENTITY.SETPOS`, `WINDOW.SETPOS` in migrated samples (`examples/kcc_*.mb`, `spinning_cube_*.mb`, `testdata/raylib_extras_complete_test.mb`) |
| Validation | `go test ./...` (passed); `go run ./tools/apidoc` (passed) |

## Intentional follow-ups (v1.0 / future)

- Remove deprecated `MAKE` / `SETPOSITION` aliases and strict-mode-by-default when ready for a breaking release (directive Part 8).
- Optional: deduplicate non-`AUDIO` manifest rows where the same `(key, args)` appears twice (e.g. legacy duplicates outside `AUDIO.*`); `AUDIO.*` overloads are guarded by `TestNoDuplicateAudioManifestArgSignatures`.
- Namespace-by-namespace audit of **every** `SET*` without `GET*` remains a documentation/runtime maintenance task as APIs grow.

## Success criteria mapping (directive Part 9)

1. **Naming:** CREATE + SETPOS canonical in manifest/tests; aliases deprecated with warnings or strict errors.
2. **Universal methods:** Documented and dispatched via `handlecall`; per-type gaps tightened incrementally in runtime.
3. **Documentation:** Style guide, conventions, migration, UNIVERSAL handle doc, traceability docs.
4. **Examples:** Migrated to `CREATE` in shipped examples.
5. **Tooling:** LSP completion bias + snippets; optional strict CLI.
6. **Ergonomics:** STYLE_GUIDE documents builder preference; long-arg elimination is iterative.
7. **Symmetry:** Policy encoded in manifest tests where applicable; full GET coverage is ongoing per namespace.
8. **Easy Mode:** Documented as convenience; `CreateCamera` calls canonical `CAMERA.CREATE`.
