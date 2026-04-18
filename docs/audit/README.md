# Audit artifacts

This directory groups **generated inventories**, **manual implementation logs**, and **historical / baseline dumps** so the repository root stays readable on GitHub.

## Layout

| Subpath | Contents |
|---------|----------|
| *(this folder)* | **`MASTER_AUDIT*.txt`**, **`manifest_keys.txt`**, **`runtime_keys.txt`**, **`COMMAND_AUDIT.txt`** (hand-maintained implementation log), **`REFERENCE_KEY_COVERAGE.txt`**, QOL audit `QOL_AUDIT*.txt`, **`raylib_symbol_gap.txt`**. |
| **`baselines/`** | Optional regression captures (race detector, gccheckmark, escape analysis, Valgrind/Dr. Memory placeholders, etc.) — see [`docs/MEMORY.md`](../MEMORY.md). |
| **`archives/`** | Large one-off exports and scratch audits (raw doc dumps, grep captures, legacy notes). Not required to build or run the project. |

## Regeneration

- **`python tools/gen_master_audit.py`** (from repo root) — updates **`MASTER_AUDIT*.txt`** and **`REFERENCE_KEY_COVERAGE.txt`**.
- **`tools/audit_manifest.ps1`** → **`manifest_keys.txt`**; **`tools/extract_runtime_keys.ps1`** → **`runtime_keys.txt`**; **`tools/diff_keys.ps1`** writes **`docs/MISSING_COMMANDS_AUDIT.md`**.
- **`go run ./tools/cmdaudit`** — updates **`docs/COMMAND_AUDIT.md`** (separate from **`COMMAND_AUDIT.txt`** here).

## QOL / built-in audit

Generated for DarkBASIC-style parity tracking. Regenerate when adding runtime commands.

- **`QOL_AUDIT_REGISTERED.txt`** — sorted unique second arguments to `r.Register("KEY", …)` / `reg.Register` across `runtime/**/*.go` (excluding `*_test.go`).
- **`QOL_AUDIT_DUPLICATES.txt`** — keys that appear more than once. **Expected:** stub vs `cgo` pairs (`//go:build !cgo` vs `cgo`) register the same key in mutually exclusive files; only one implementation is linked per build. True conflicts = same key in two files with the **same** build tag.
- **`QOL_AUDIT.txt`** — human matrix: spec / common name → status (DONE / PARTIAL / MISSING / DOCONLY) and primary implementation path.

Canonical **instant game / QOL** implementation lives in **`runtime/mbgame`** (registered from `compiler/pipeline/pipeline.go`). Do **not** add a second package that registers the same command keys.
