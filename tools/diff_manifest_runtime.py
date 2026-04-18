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

REGISTER_RE = re.compile(r'\.Register\("([^"]+)"')


def load_manifest_keys() -> list[str]:
    data = json.loads(MANIFEST_JSON.read_text(encoding="utf-8"))
    return sorted({c["key"].upper() for c in data["commands"]})


def load_runtime_keys() -> list[str]:
    keys: list[str] = []
    for p in RUNTIME_DIR.rglob("*.go"):
        text = p.read_text(encoding="utf-8", errors="replace")
        for m in REGISTER_RE.finditer(text):
            keys.append(m.group(1).upper())
    return sorted(set(keys))


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
        "**Generated** by `python tools/diff_manifest_runtime.py --write`. Runtime keys follow the same `.Register(` string heuristic as `extract_runtime_keys.ps1` over `runtime/**/*.go`; not every builtin uses that pattern. Treat gaps as triage hints, not a complete defect list.",
        "",
        f"## In Runtime but Missing from Manifest ({len(missing_from_manifest)})",
        "These commands are registered in Go runtime code but have no entry in commands.json.",
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
    rkeys = load_runtime_keys()
    mset, rset = set(mkeys), set(rkeys)
    missing_from_manifest = sorted(rset - mset)
    missing_from_runtime = sorted(mset - rset)

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
