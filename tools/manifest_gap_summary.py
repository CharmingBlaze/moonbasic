#!/usr/bin/env python3
"""Print manifest keys with no wide-runtime match, grouped by namespace prefix.

Uses the same discovery logic as tools/diff_manifest_runtime.py (runpy-loaded).

Usage (repo root)::

    python tools/manifest_gap_summary.py
    python tools/manifest_gap_summary.py --list JOLT
"""
from __future__ import annotations

import argparse
import runpy
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def load_dmr():
    ns = runpy.run_path(str(ROOT / "tools" / "diff_manifest_runtime.py"))
    return ns["load_manifest_keys"], ns["load_runtime_key_sets"]


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument(
        "--list",
        metavar="PREFIX",
        help='Print gap keys: use a namespace like JOLT or ENTITY; use "global" for flat globals (no dot).',
    )
    args = ap.parse_args()

    load_manifest_keys, load_runtime_key_sets = load_dmr()
    mset = set(load_manifest_keys())
    _, wide = load_runtime_key_sets()
    gap = sorted(mset - wide)

    if args.list:
        pfx = args.list.strip()
        pfx_u = pfx.upper()
        is_global = pfx_u in ("(GLOBAL)", "GLOBAL")
        for k in gap:
            if "." not in k:
                if is_global:
                    print(k)
                continue
            ns, _rest = k.split(".", 1)
            if ns.upper() == pfx_u:
                print(k)
        return 0

    # Count by first namespace segment
    from collections import Counter

    counts: Counter[str] = Counter()
    for k in gap:
        if "." in k:
            counts[k.split(".", 1)[0]] += 1
        else:
            counts["(global)"] += 1

    print(f"Manifest keys with no wide-runtime match: {len(gap)} (see tools/diff_manifest_runtime.py)\n")
    print(f"{'count':>6}  namespace / bucket")
    print(f"{'------':>6}  ------------------")
    for ns, n in sorted(counts.items(), key=lambda x: (-x[1], x[0])):
        print(f"{n:6d}  {ns}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
