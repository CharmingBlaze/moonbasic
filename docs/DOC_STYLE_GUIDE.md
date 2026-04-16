# Documentation Style Guide

Command and module **reference** pages use a single **page shape** (the **WAVE pattern**). API naming still follows [STYLE_GUIDE.md](../STYLE_GUIDE.md): **registry-first** **`NAMESPACE.ACTION`** in headings and examples unless the page is explicitly about Easy Mode compatibility.

**Live example:** [reference/WAVE.md](reference/WAVE.md) — use it as the structural template when writing or revising docs.

---

## Page shape (WAVE pattern)

1. **Title** — `# [Topic] Commands` (or `# [Namespace] Commands`).
2. **Purpose** — One sentence under the title: what this module is for.
3. **`## Core Workflow`** — Short prose (or a short numbered list if the call order is strict). Explain the typical path (load → edit → export, init → loop → close, etc.).
4. **`---`** — Separator after the workflow section (and before the first command).
5. **Each command** — Repeat this block:
   - `### `signatureWithBackticks`` — Full call as it appears in source (registry style preferred: **`WAVE.LOAD(path)`** not `Wave.Load`).
   - One short paragraph: what it returns or does.
   - If parameters need detail, use a bullet list: `- \`paramName\`: explanation` (plain names; no Blitz **`#` / `$` / `?` / `%`** suffixes).
   - `---` — Separator before the next command (including the last command before **Full Example**).
6. **`## Full Example`** — One sentence describing what the program demonstrates, then a single fenced **`basic`** block with **`;`** comments and realistic cleanup.

Do **not** bury the only runnable sample in the middle of the page unless you also keep a **Full Example** at the end for copy-paste.

---

## Signatures and naming

- **Headings:** Level-3 headings are the **signature** in backticks: `### `MODULE.ACTION(arg1, arg2)``.
- **Parentheses:** Use `()` when there are no arguments.
- **Registry-first:** Prefer **`AUDIO.INIT()`**, **`WAVE.LOAD(path)`**, **`WAVE.FREE(handle)`** in new/edited reference pages — same layout as [WAVE.md](reference/WAVE.md), different spelling. Easy Mode dotted names belong in a compatibility note or alias table, not as the only documented form.

---

## Visual rhythm

- Use **`---`** between **every** command entry (and after **Core Workflow** before the first command) so long pages stay scannable.
- Keep each command’s body short; move edge cases to the narrative reference or a second example.

---

## Platform ordering (project policy)

When a reference page compares **Windows** and **Linux**, list **Windows** first and **Linux** second (tables, columns, and sentences). Exception: a page that is **only** about Linux-only internals may omit Windows or mention it second. Rationale: [DEVELOPER.md](DEVELOPER.md#platform-priority-windows-then-linux).

---

## Consistency check

Verify signatures against `compiler/builtinmanifest/commands.json` before finalizing documentation.
