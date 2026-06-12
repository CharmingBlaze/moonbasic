# moonBASIC documentation style guide

**Follow this guide** when writing or editing moonBASIC documentation (`docs/`, README files, in-repo reference).

Older pages may use the [WAVE pattern](DOC_STYLE_GUIDE.md); new **topic guides** use the **beginner template** below.

---

## Audience

**User documentation** explains how to **use** moonBASIC:

- Install from [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest)
- Run games with **`moonrun`**
- Check, compile, and edit with **`moonbasic`** (`--check`, `--lsp`, `.mbc`)
- Scaffold with **`moonbasic new`**, ship with **`package`**

**Do not** put compiler-build instructions in user-facing pages (`go run`, `go build`, `-tags fullruntime`, `CGO_ENABLED`, `tools/apidoc`). Those live in [DEVELOPER.md](DEVELOPER.md), [BUILDING.md](BUILDING.md), [CONTRIBUTING.md](../CONTRIBUTING.md).

---

## How a beginner should read any system doc

Every system page should let someone answer five questions:

| Question | Section to include |
|----------|-------------------|
| **What problem does this solve?** | Opening quote + **At a glance** |
| **Do I need this in my game?** | **When to use** / **When to skip** |
| **What do I call first, second, third?** | **Core workflow** (numbered, with **why**) |
| **What if I picked the wrong API?** | **Choose the right tool** comparison table |
| **What goes wrong in practice?** | **Common mistakes** |

Then: commands (or link to [COMMAND_REGISTRY.md](systems/COMMAND_REGISTRY.md)), **full example**, **`moonbasic --check`** / **`moonrun`** note.

**Reading order for new authors:** [BEGIN_HERE.md](BEGIN_HERE.md) → [systems/00-START.md](systems/00-START.md) → [systems/GUIDES.md](systems/GUIDES.md) → numbered systems `01`–`11` → deep [reference/](reference/).

---

## Goals

1. **Work on GitHub** — headings, tables, fenced code, relative links, TOC.
2. **Work off GitHub** — VS Code, static sites, PDF.
3. **Explain how and why** — not only signatures.
4. **Align with** `compiler/builtinmanifest/commands.json` (arity).
5. **Case insensitivity** — document `NAMESPACE.COMMAND`; note `app.open` = `APP.OPEN` ([LANGUAGE.md](LANGUAGE.md)).

---

## File types

| Type | Location | Purpose |
|------|----------|---------|
| **Onboarding** | `docs/BEGIN_HERE.md`, `systems/00-START.md` | Install, loop, why each step |
| **Hub** | `docs/README.md`, `systems/README.md`, `systems/GUIDES.md` | Maps |
| **Systems (overview)** | `docs/systems/01-*.md` … `11-*.md` | 40-system summaries |
| **Topic guides (deep)** | `docs/systems/guides/*.md` | Entity, 2D/3D collision, UI, multiplayer, … |
| **Reference** | `docs/reference/*.md` | Namespace deep dives |
| **Generated** | `API_CONSISTENCY.md`, `systems/COMMAND_REGISTRY.md` | Every overload (maintainers regenerate via [DEVELOPER.md](DEVELOPER.md)) |

**Rule:** Beginner **how/why** narratives → `systems/` or `systems/guides/`. Exhaustive overload lists → generated registries + `reference/`.

---

## Beginner topic guide template (`systems/guides/*.md`)

Use this for entity, collision, UI, multiplayer, and similar **deep** pages.

```markdown
# [Topic title]

> One sentence: the player-visible problem this solves.

**Namespaces:** `ENTITY`, `…` · **Status:** Shipped · **Platform:** full runtime (Windows / Linux)

**Commands (complete list):** [COMMAND_REGISTRY.md#anchor](COMMAND_REGISTRY.md#anchor)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Key commands](#key-commands)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | … |
| **You need first** | Window + loop ([00-START.md](../00-START.md)) |
| **Typical games** | Platformer, FPS, menu, … |

---

## When to use this system

**Use when:** …

**Skip when:** … (link to simpler alternative)

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| … | `NAMESPACE.*` | … |

---

## Core workflow

1. **Step** — command — **Why:** reason
2. …

---

## Key commands

(Top 8–15 commands with args table + example each, OR tables grouped by task)

---

## Full example

; Runnable sketch — moonrun required if graphics

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| … | … |

---

## Memory notes

---

## See also
```

---

## Overview systems template (`systems/01-*.md`)

Shorter than topic guides. Link to **guides/** for depth.

Required sections: **Why this page**, TOC, per-system blocks, **Full example**, **See also**, link to **COMMAND_REGISTRY** anchor.

---

## Command entry rules

1. **Heading** — `### `COMMAND.NAME(args)``
2. **Summary** — what + **when to call** (one sentence).
3. **Why** — optional one line if non-obvious (e.g. “Call after `PHYSICS2D.STEP`”).
4. **Arguments table** — when args exist.
5. **Returns** — always.
6. **Example** — `basic` fence, `;` comments.
7. **Aliases** — in **Aliases**, not duplicate headings.
8. **`---`** after each command block in detailed guides.

For **hundreds** of overloads in one namespace, use **task tables** + link to [COMMAND_REGISTRY.md](systems/COMMAND_REGISTRY.md) instead of duplicating every row.

---

## Code examples

- Fence: `basic`
- Comments: `;`
- Prefer handle chaining: `cube.pos(0,1,0).turn(0,45,0)`
- End with check/run line: `moonbasic --check` / `moonrun`
- No legacy `$` / `#` variable suffixes in new docs

---

## “Choose the right tool” patterns

Document these comparisons wherever beginners confuse APIs:

| Topic | Options to compare |
|-------|-------------------|
| **2D overlap** | Manual rects, `COLLISION.*` math, `PHYSICS2D` + `BODY2D` |
| **3D overlap** | `COLLISION.*`, `PHYSICS3D` + Jolt, `PICK.*` rays |
| **Movement** | `ENTITY.MOVE`, kinematic body, dynamic `BODY3D`, `CHAR.*` KCC |
| **UI** | `DRAW.TEXT` HUD, `GUI.*` widgets, `UI.*` helpers |
| **Networking** | `SERVER.*`/`CLIENT.*`, `NET.*`+`PEER.*`, external HTTP API |
| **Drawing 3D** | `ENTITY.DRAWALL`, `SCENE.DRAW`, raw `DRAW3D.*` |

---

## Navigation

```
docs/BEGIN_HERE.md
  → systems/00-START.md
  → systems/GUIDES.md (entity, collision, UI, multiplayer, …)
  → systems/01–11 (overview)
  → systems/COMMAND_REGISTRY.md (all overloads)
  → reference/*.md (namespace internals)
```

Every hub: `#` title, TOC, **See also**.

---

## Terminology

| Term | Meaning |
|------|---------|
| **Command** | `NAMESPACE.NAME` or global `PRINT` |
| **Handle** | Runtime reference (camera, body, texture) |
| **System** | Beginner group (APP, ENTITY, …) |
| **Topic guide** | Deep `guides/*.md` page |
| **Namespace** | First segment of dotted name (`CAMERA`) |

---

## Platform notes

- **Windows first**, then Linux.
- **Full runtime** vs **compiler-only** — graphics, audio, physics, net need **`moonrun`** from full runtime zip.

---

## Checklist before merging doc changes

- [ ] **At a glance** + **when to use** + **choose the right tool** on topic guides
- [ ] Core workflow steps include **why**
- [ ] **Common mistakes** table on topic guides
- [ ] Signatures match `commands.json`
- [ ] Link to COMMAND_REGISTRY for full overload list
- [ ] TOC links work on GitHub
- [ ] `moonbasic --check` / `moonrun` noted for examples
- [ ] No compiler-build instructions on user pages

---

## Related

- [STYLE_GUIDE.md](../STYLE_GUIDE.md) — API naming in source
- [DOC_STYLE_GUIDE.md](DOC_STYLE_GUIDE.md) — legacy WAVE
- [systems/README.md](systems/README.md) — 40-system index
- [systems/GUIDES.md](systems/GUIDES.md) — deep topic index
