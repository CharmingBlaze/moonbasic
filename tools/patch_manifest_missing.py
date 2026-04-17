#!/usr/bin/env python3
"""Append runtime-registered commands missing from compiler/builtinmanifest/commands.json.

Entry data lives in the ``manifest_entries`` subpackage — one module per
namespace group — so no single file exceeds a few hundred lines.
See STYLE_GUIDE.md §Documentation for project conventions.

Usage::

    python tools/patch_manifest_missing.py
"""
import json
import os
import sys

# Ensure the tools/ directory is importable when run as a script.
sys.path.insert(0, os.path.dirname(__file__))

from manifest_entries.audio import ENTRIES as AUDIO
from manifest_entries.camera import ENTRIES as CAMERA
from manifest_entries.data import ENTRIES as DATA
from manifest_entries.draw import ENTRIES as DRAW
from manifest_entries.draw_objects import ENTRIES as DRAW_OBJECTS
from manifest_entries.entity import ENTRIES as ENTITY
from manifest_entries.input import ENTRIES as INPUT
from manifest_entries.misc import ENTRIES as MISC
from manifest_entries.model import ENTRIES as MODEL
from manifest_entries.physics import ENTRIES as PHYSICS
from manifest_entries.ray import ENTRIES as RAY
from manifest_entries.rendering import ENTRIES as RENDERING
from manifest_entries.world import ENTRIES as WORLD

ROOT = os.path.join(os.path.dirname(__file__), "..")
MAN = os.path.join(ROOT, "compiler", "builtinmanifest", "commands.json")

# Aggregate all namespace modules into one flat list.
NEW: list[tuple] = (
    AUDIO + CAMERA + DRAW + DRAW_OBJECTS + ENTITY + INPUT + DATA
    + PHYSICS + RENDERING + WORLD + RAY + MODEL + MISC
)


def main():
    with open(MAN, encoding="utf-8") as f:
        root = json.load(f)
    have = {c["key"] for c in root["commands"]}
    added = 0
    for spec in NEW:
        key = spec[0]
        if key in have:
            continue
        args = spec[1]
        rec = {"key": key, "args": args, "phase": spec[3]}
        if len(spec) > 2 and spec[2]:
            rec["returns"] = spec[2]
        root["commands"].append(rec)
        added += 1
        have.add(key)
    with open(MAN, "w", encoding="utf-8", newline="\n") as f:
        json.dump(root, f, indent=2)
        f.write("\n")
    print("added", added, "commands")


if __name__ == "__main__":
    main()
