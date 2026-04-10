#!/usr/bin/env python3
"""Remove legacy BASIC type-suffix clutter from Markdown *prose* only.

- Skips fenced code blocks (``` ... ```).
- Skips generated docs/API_CONSISTENCY.md.
- Does not strip '#' from Markdown headings (lines that start with optional whitespace then #).

Safe rules:
- word$  -> word   (string suffix in prose; not inside URLs)
- word#  -> word   only when # is not followed by ( or digit (avoids result#(0) and grid#1)
- Optional markers like (name?) -> (name, optional) in prose
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
# LANGUAGE.md defines $ # ? suffix rules; keep authoritative wording.
EXCLUDE_NAMES = frozenset({"API_CONSISTENCY.md", "LANGUAGE.md"})


def transform_prose(chunk: str) -> str:
    lines = chunk.splitlines(keepends=True)
    out: list[str] = []
    for line in lines:
        if re.match(r"^\s*#+\s", line):
            out.append(line)
            continue
        s = line
        # URLs: do not touch query strings
        s = re.sub(r"\b([A-Za-z_][A-Za-z0-9_]*)\$(?![A-Za-z0-9_])", r"\1", s)
        s = re.sub(
            r"(?<![A-Za-z0-9_#])([A-Za-z_][A-Za-z0-9_]*)\#(?![#\d\(])",
            r"\1",
            s,
        )
        s = re.sub(r"\(([A-Za-z][A-Za-z0-9_]*)\?\)", r"(\1, optional)", s)
        out.append(s)
    return "".join(out)


def process_file(path: Path) -> bool:
    if path.name in EXCLUDE_NAMES:
        return False
    text = path.read_text(encoding="utf-8")
    parts = re.split(r"(```[\s\S]*?```)", text)
    new_parts: list[str] = []
    for part in parts:
        if part.startswith("```"):
            new_parts.append(part)
        else:
            new_parts.append(transform_prose(part))
    new_text = "".join(new_parts)
    if new_text != text:
        nl = "\r\n" if "\r\n" in text else "\n"
        path.write_text(new_text.replace("\n", nl), encoding="utf-8")
        return True
    return False


def main() -> int:
    changed: list[str] = []
    singles = [ROOT / "README.md", ROOT / "ARCHITECTURE.md", ROOT / "CONTRIBUTING.md", ROOT / "AGENTS.md"]
    for f in singles:
        if f.is_file() and process_file(f):
            changed.append(str(f.relative_to(ROOT)))
    doc_root = ROOT / "docs"
    if doc_root.is_dir():
        for p in sorted(doc_root.rglob("*.md")):
            if process_file(p):
                changed.append(str(p.relative_to(ROOT)))
    for c in changed:
        print("updated:", c)
    print(f"strip_doc_typography: {len(changed)} file(s) changed", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
