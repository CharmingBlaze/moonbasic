# Directive baseline / conformance snapshot

Generated as part of the full directive rollout. Use this as a reference point for regressions.

**Normative roadmap:** [API_STANDARDIZATION_DIRECTIVE.md](API_STANDARDIZATION_DIRECTIVE.md)

## Canonical API policy (compiler-enforced manifest tests)

- Every `*.CREATE*` command entry has a matching deprecated `*.MAKE*` alias (same arity) with deprecation text.
- Every `*.SETPOS` overload has a matching `*.SETPOSITION` alias where applicable (includes **ENTITY** and **WINDOW** position APIs).

See `compiler/builtinmanifest/api_standardization_test.go`.

## Deprecation surfacing

- CLI: `go run . --check` and compile print `[moonBASIC] Warning:` lines for deprecated `MAKE` / `SETPOSITION`-style aliases (see `compiler/pipeline/compile.go`).
- **Strict mode:** `moonbasic --check --strict-deprecated file.mb` fails on deprecated aliases instead of warning (`compiler/semantic/analyze.go`, `pipeline.CheckOptions`).
- LSP: `textDocument/publishDiagnostics` includes warning-level diagnostics for deprecated builtins (`lsp/server.go`).

## Universal handle methods

See [reference/UNIVERSAL_HANDLE_METHODS.md](reference/UNIVERSAL_HANDLE_METHODS.md) and `vm/handlecall.go`.

## Documentation entry points

- Normative style: `STYLE_GUIDE.md`
- API naming: `docs/reference/API_CONVENTIONS.md`
- Easy Mode scope: `docs/EASY_MODE.md`
- Migration: `docs/MIGRATION_CREATE_FROM_MAKE.md`
- Language (case insensitivity): `docs/LANGUAGE.md`

## Regenerating API consistency

```bash
go run ./tools/apidoc
```

Outputs/updates `docs/API_CONSISTENCY.md`.

## Full validation (local)

```bash
go test ./...
go run ./tools/apidoc
go run . --check testdata/camera_complete_test.mb
```

Adjust tags per `docs/DEVELOPER.md` for fullruntime-only tests.
