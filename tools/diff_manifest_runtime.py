"""Regenerate or verify manifest vs runtime Register() key lists and MISSING_COMMANDS_AUDIT.md.

Cross-platform equivalent of tools/diff_keys.ps1 + audit_manifest.ps1 + extract_runtime_keys.ps1.
Run from repo root:

  python tools/diff_manifest_runtime.py --write   # refresh docs/audit/*.txt + docs/MISSING_COMMANDS_AUDIT.md
  python tools/diff_manifest_runtime.py --check   # fail if committed files drift
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
MANIFEST_JSON = ROOT / "compiler" / "builtinmanifest" / "commands.json"
RUNTIME_DIR = ROOT / "runtime"
AUDIT_DIR = ROOT / "docs" / "audit"
MANIFEST_TXT = AUDIT_DIR / "manifest_keys.txt"
RUNTIME_TXT = AUDIT_DIR / "runtime_keys.txt"
MISSING_MD = ROOT / "docs" / "MISSING_COMMANDS_AUDIT.md"

# Literal `.Register("KEY"` (same as extract_runtime_keys.ps1).
REGISTER_RE = re.compile(r'\.Register\("([^"]+)"')
# Helpers that register two string keys without repeating `.Register("...")` in source.
REG_FLAT = re.compile(r'regFlat\(\s*"([^"]+)"\s*,\s*"([^"]+)"')
REG_LEGACY2 = re.compile(r'regLegacy2\(\s*"([^"]+)"\s*,\s*"([^"]+)"')
REG_RT0 = re.compile(r'regRT0\(\s*"([^"]+)"\s*,\s*"([^"]+)"')
# bitwise: reg("core.BAND", "BAND", fn)
REG_PAIR = re.compile(r'\breg\(\s*"([^"]+)"\s*,\s*"([^"]+)"\s*,')
# texture atlas: reg("TEXTURE.X", (*Module).method)
REG_SINGLE = re.compile(r'\breg\(\s*"([^"]+)"\s*,\s*\(')


def load_manifest_keys() -> list[str]:
    data = json.loads(MANIFEST_JSON.read_text(encoding="utf-8"))
    return sorted({c["key"].upper() for c in data["commands"]})


def load_runtime_key_sets() -> tuple[set[str], set[str]]:
    """Return (strict_register_literals, wide_all_sources).

    *strict* — only `.Register("KEY"` string literals (true compiler-facing registrations).
    *wide* — strict plus `regFlat` / `regLegacy2` / `regRT0` / two-arg `reg` / single-arg `reg`
    helpers so manifest keys registered only via variables still count as implemented.
    """
    strict: set[str] = set()
    wide: set[str] = set()
    for p in RUNTIME_DIR.rglob("*.go"):
        text = p.read_text(encoding="utf-8", errors="replace")
        for m in REGISTER_RE.finditer(text):
            k = m.group(1).upper()
            strict.add(k)
            wide.add(k)
        for rx in (REG_FLAT, REG_LEGACY2, REG_RT0):
            for m in rx.finditer(text):
                wide.add(m.group(1).upper())
                wide.add(m.group(2).upper())
        for m in REG_PAIR.finditer(text):
            wide.add(m.group(1).upper())
            wide.add(m.group(2).upper())
        for m in REG_SINGLE.finditer(text):
            wide.add(m.group(1).upper())
    return strict, wide


def load_runtime_keys() -> list[str]:
    _, wide = load_runtime_key_sets()
    return sorted(wide)


def write_lines(path: Path, lines: list[str]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("\n".join(lines) + "\n", encoding="utf-8", newline="\n")


def read_key_lines(path: Path) -> list[str]:
    raw = path.read_text(encoding="utf-8-sig")
    return [ln.strip().upper() for ln in raw.splitlines() if ln.strip()]


def build_missing_md(missing_from_manifest: list[str], missing_from_runtime: list[str]) -> str:
    lines: list[str] = [
        "# Missing Commands Audit",
        "",
        "**Generated** by `python tools/diff_manifest_runtime.py --write`. **In Runtime but Missing from Manifest** lists only **string literals** in `.Register(\"…\")` (alias helpers like `regFlat` are not counted here). **In Manifest but Missing from Runtime** uses a **wide** scan: those literals plus keys from `regFlat`, `regLegacy2`, `regRT0`, two-arg `reg`, and single-arg `reg`. Not every builtin is covered. Treat gaps as triage hints.",
        "",
        f"## In Runtime but Missing from Manifest ({len(missing_from_manifest)})",
        "These commands appear as string literals in `.Register(\"…\")` but have no entry in commands.json.",
        "The compiler will reject .mb scripts that try to use them.",
        "",
    ]
    for k in missing_from_manifest:
        lines.append(f"- `{k}`")
    lines.append("")
    lines.append(f"## In Manifest but Missing from Runtime ({len(missing_from_runtime)})")
    lines.append("These commands are declared in commands.json but have no runtime registration.")
    lines.append("Scripts compile but will fail at runtime with 'unknown command'.")
    lines.append("")
    for k in missing_from_runtime:
        lines.append(f"- `{k}`")
    lines.append("")
    return "\n".join(lines)


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument(
        "--write",
        action="store_true",
        help="Write manifest_keys.txt, runtime_keys.txt, and MISSING_COMMANDS_AUDIT.md",
    )
    ap.add_argument(
        "--check",
        action="store_true",
        help="Exit 1 if on-disk files differ from regenerated content",
    )
    args = ap.parse_args()
    if args.write == args.check:
        ap.error("specify exactly one of --write or --check")

    mkeys = load_manifest_keys()
    strict_reg, wide_runtime = load_runtime_key_sets()
    mset = set(mkeys)
    rkeys = sorted(wide_runtime)
    missing_from_manifest = sorted(strict_reg - mset)
    missing_from_runtime = sorted(mset - wide_runtime)

    if args.write:
        write_lines(MANIFEST_TXT, mkeys)
        write_lines(RUNTIME_TXT, rkeys)
        MISSING_MD.write_text(
            build_missing_md(missing_from_manifest, missing_from_runtime),
            encoding="utf-8",
            newline="\n",
        )
        print(f"Wrote {MANIFEST_TXT.relative_to(ROOT)} ({len(mkeys)} keys)")
        print(f"Wrote {RUNTIME_TXT.relative_to(ROOT)} ({len(rkeys)} keys)")
        print(f"Wrote {MISSING_MD.relative_to(ROOT)}")
        print(f"  runtime not in manifest: {len(missing_from_manifest)}")
        print(f"  manifest not in runtime: {len(missing_from_runtime)}")
        return 0

    ok = True
    if read_key_lines(MANIFEST_TXT) != mkeys:
        print(f"MISMATCH: {MANIFEST_TXT} does not match commands.json (run with --write)", file=sys.stderr)
        ok = False
    if read_key_lines(RUNTIME_TXT) != rkeys:
        print(f"MISMATCH: {RUNTIME_TXT} does not match runtime/**/*.go scan (run with --write)", file=sys.stderr)
        ok = False
    expected_md = build_missing_md(missing_from_manifest, missing_from_runtime)
    actual_md = MISSING_MD.read_text(encoding="utf-8-sig").replace("\r\n", "\n")
    if not actual_md.endswith("\n"):
        actual_md += "\n"
    if actual_md != expected_md:
        print(f"MISMATCH: {MISSING_MD} out of date (run with --write)", file=sys.stderr)
        ok = False
    if not ok:
        return 1
    print("OK: manifest_keys.txt, runtime_keys.txt, and MISSING_COMMANDS_AUDIT.md are current.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
