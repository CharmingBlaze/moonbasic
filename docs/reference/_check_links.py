"""Check for broken internal markdown links in docs/reference."""
import os, re

d = "c:/Users/rain/Documents/GO/moonbasic/docs/reference"
all_files = set(os.listdir(d))
link_pat = re.compile(r"\[([^\]]*)\]\(([^)]+)\)")

broken = []
for fn in sorted(os.listdir(d)):
    if not fn.endswith(".md"):
        continue
    fp = os.path.join(d, fn)
    if not os.path.isfile(fp):
        continue
    text = open(fp, encoding="utf-8").read()
    for m in link_pat.finditer(text):
        target = m.group(2)
        # Skip external URLs and anchors
        if target.startswith("http") or target.startswith("#"):
            continue
        # Strip anchor from target
        target_file = target.split("#")[0]
        # Resolve relative path
        resolved = os.path.normpath(os.path.join(d, target_file))
        if not os.path.exists(resolved):
            line = text[: m.start()].count("\n") + 1
            broken.append(f"{fn}:{line}: -> {target}")

print(f"{len(broken)} broken links:")
for b in broken:
    print(f"  {b}")
