#!/usr/bin/env python3
"""Compare compiler/builtinmanifest/commands.json keys to runtime Register(\"...\") names.

Only string literals inside Register(\"NAME\" are detected; dynamic Register(nameVar, ...)
or shared Register(shortName, ...) tables are not counted, so \"in manifest only (no Register)\"
can include false positives (e.g. math commands registered via a loop).
"""
import json
import os
import re

ROOT = os.path.join(os.path.dirname(__file__), "..")
MAN = os.path.join(ROOT, "compiler", "builtinmanifest", "commands.json")
RT = os.path.join(ROOT, "runtime")

with open(MAN, encoding="utf-8") as f:
    mkeys = {c["key"] for c in json.load(f)["commands"]}

reg = set()
pat = re.compile(r'Register\(\s*"([^"]+)"')
for root, _, files in os.walk(RT):
    for fn in files:
        if not fn.endswith(".go"):
            continue
        p = os.path.join(root, fn)
        try:
            t = open(p, encoding="utf-8").read()
        except OSError:
            continue
        for m in pat.finditer(t):
            reg.add(m.group(1).upper())

only_m = sorted(mkeys - reg)
only_r = sorted(reg - mkeys)
print("manifest unique keys:", len(mkeys))
print("Register() unique keys:", len(reg))
print("in manifest only (no Register):", len(only_m))
for k in only_m[:80]:
    print(" ", k)
if len(only_m) > 80:
    print(" ...")
print()
print("Registered but NOT in manifest:", len(only_r))
for k in only_r:
    print(k)
