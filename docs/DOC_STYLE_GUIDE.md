# Documentation Style Guide

Command and module **reference** pages use a single **page shape** (the **WAVE pattern**). API naming still follows [STYLE_GUIDE.md](../STYLE_GUIDE.md): **registry-first** **`NAMESPACE.ACTION`** in headings and examples unless the page is explicitly about Easy Mode compatibility.

**Live example:** [reference/WAVE.md](reference/WAVE.md) — use it as the structural template when writing or revising docs.

---

## Page shape (WAVE pattern)

1. **Title** — `# [Topic] Commands` (or `# [Namespace] Commands`).
2. **Purpose** — One or two sentences under the title: what this module is for and any platform/build constraints.
3. **`## Core Workflow`** — Short prose or a numbered list. Explain the typical path (init → loop → cleanup). If the module supports handle chaining, add a `### Method chaining` sub-section here.
4. **`---`** — Separator after the workflow section and before the first command.
5. **Section headers (`## Section Name`)** — Group related commands under `##` headers (e.g. `## World Management`, `## Body Creation`, `## Collision Queries`). Each section may be as short as one command.
6. **Each command** — Repeat this exact block:
   ```
   ### `NAMESPACE.COMMAND(arg1, arg2)` 

   One short paragraph: what it does or returns.

   - *Handle shortcut*: `handle.method(arg1, arg2)`  ← include only if a handle dispatch exists
   - **Arguments**:
     - `arg1` (type): What this argument represents.
     - `arg2` (type): What this argument represents.
   - **Returns**: (type) Description of the return value (e.g., a new handle or `self` for chaining).

   - **Example**:
     ```basic
     ; Short, focused snippet showing this command in action
     res = NAMESPACE.COMMAND(a, b)
     ```

   ---
   ```
   **Rules:**
   - The `###` heading signature ends with a **trailing space** before the newline. This is the canonical style.
   - One blank line between the heading and the body paragraph.
   - If the command takes no arguments use `()`.
   - For aliased commands: `Alias: NAMESPACE.OTHERNAME` in the body — do **not** create a duplicate heading.
   - Handle shortcut bullet uses `*Handle shortcut*:` (italic, no bold).
   - **Arguments** and **Returns** sections are mandatory for clarity.
   - **Example** snippets are highly recommended for complex commands.
   - `---` separator comes **after** the last block of the command entry.
7. **`## Full Example`** — One sentence intro, then a single fenced ` ```basic ` block with `;` comments and realistic cleanup (`FREE`, `STOP`, `WINDOW.CLOSE`).
8. **`## See also`** — Short bullet list of related pages with one-phrase descriptions.

Do **not** bury the only runnable sample in the middle of the page — always keep `## Full Example` at the end.

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
