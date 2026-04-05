# QOL / built-in audit artifacts

Generated for DarkBASIC-style parity tracking. Regenerate when adding runtime commands.

- **`QOL_AUDIT_REGISTERED.txt`** — sorted unique second arguments to `r.Register("KEY", …)` / `reg.Register` across `runtime/**/*.go` (excluding `*_test.go`).
- **`QOL_AUDIT_DUPLICATES.txt`** — keys that appear more than once. **Expected:** stub vs `cgo` pairs (`//go:build !cgo` vs `cgo`) register the same key in mutually exclusive files; only one implementation is linked per build. True conflicts = same key in two files with the **same** build tag.
- **`QOL_AUDIT.txt`** — human matrix: spec / common name → status (DONE / PARTIAL / MISSING / DOCONLY) and primary implementation path.

Canonical **instant game / QOL** implementation lives in **`runtime/mbgame`** (registered from `compiler/pipeline/pipeline.go`). Do **not** add a second package that registers the same command keys.

**Repo root (`python tools/gen_master_audit.py`)** also writes **`REFERENCE_KEY_COVERAGE.txt`** here: manifest keys that appear verbatim in **`docs/reference/*.md`** and **`compiler/errors/MoonBasic.md`**, plus per-file hit counts (same key may count in multiple files).
