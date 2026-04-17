# tools/

Developer utilities for auditing and maintaining the MoonBASIC compiler manifest and documentation.

## Manifest Patching

| Script | Purpose |
|--------|---------|
| `patch_manifest_missing.py` | Appends runtime-registered commands missing from `compiler/builtinmanifest/commands.json`. Entry data is split across the `manifest_entries/` subpackage (one module per namespace group). Idempotent — re-running adds only genuinely new entries. |
| `manifest_entries/` | Subpackage containing `(key, args, returns, phase)` tuples organised by namespace: `audio`, `camera`, `data`, `draw`, `draw_objects`, `entity`, `input`, `misc`, `model`, `physics`, `ray`, `rendering`, `world`. Each file has a docstring describing its scope. |

## Auditing

| Script | Purpose |
|--------|---------|
| `audit_manifest.ps1` | Extracts unique keys from `commands.json` → `manifest_keys.txt`. |
| `extract_runtime_keys.ps1` | Greps `r.Register("KEY"` patterns from `runtime/` Go source → `runtime_keys.txt`. |
| `diff_keys.ps1` | Compares the two key sets and writes `docs/MISSING_COMMANDS_AUDIT.md`. |
| `audit_manifest_runtime.py` | Python alternative for the same audit. |
| `gen_master_audit.py` | Generates a broader master-audit report. |

## Documentation

| Script | Purpose |
|--------|---------|
| `cmdaudit/` | Go tool (`go run ./tools/cmdaudit`) that audits manifest vs. doc-file coverage by namespace. |
| `apidoc/` | Go tool (`go run ./tools/apidoc`) that regenerates `docs/API_CONSISTENCY.md` from the manifest. |
| `strip_doc_typography.py` | Normalises unicode characters in Markdown doc files. |

## Obsolete / Intermediate Artefacts

These files were generated during initial audit work and can be safely removed:

- `_audit_out.txt` — raw diff output (superseded by `docs/MISSING_COMMANDS_AUDIT.md`).
- `new_manifest_entries.json` — early generic-arg entries (superseded by `manifest_entries/`).
- `gen_manifest_entries.ps1` — early PowerShell generator (superseded by the Python subpackage).
