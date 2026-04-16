# MoonBASIC Style Guide (v1.0 Standard)

**Roadmap:** Full standardization goals, phased checklist, and Easy Mode placement are in [docs/API_STANDARDIZATION_DIRECTIVE.md](docs/API_STANDARDIZATION_DIRECTIVE.md).

## Official API Style

MoonBASIC uses `Namespace.Method` as the canonical API style. All new code should prioritize this pattern to ensure consistency across the ecosystem.

### Canonical Pattern
Use `CREATE` for constructors and `SETPOS` for position setters. Methods should return `self` to support chaining.

```basic
' Recommended: Method Chaining
camera = CAMERA.CREATE()
              .pos(0, 10, 20)
              .look(0, 0, 0)
              .fov(60)

' Alternative: Explicit Namespace Calls
model = MODEL.CREATE()
MODEL.SETPOS(model, 10, 0, 10)
MODEL.SETROT(model, 0, 90, 0)

' Symmetry: CREATE/FREE, SET/GET
pos = model.pos()
MODEL.FREE(model)
```

## Naming Conventions

- **Variables:** `camelCase` (e.g., `heroHandle`, `mainLight`)
- **Constants:** `SCREAMING_SNAKE_CASE` (e.g., `GRAVITY_VAL`)
- **Types:** `PascalCase` (e.g., `EnemyMetadata`)
- **Namespaces:** `SCREAMING_SNAKE_CASE` (e.g., `CAMERA`, `VEHICLE`)
- **Methods:** `SCREAMING_SNAKE_CASE` in reference, but command lookup is case-insensitive.

### No Blitz-style suffix characters

moonBASIC **does not** use **`#`**, **`$`**, **`?`**, or **`%`** as type or string-function suffixes on names (no `speed#`, `msg$`, `ok?`, or `x%` style). Use **plain identifiers**, implicit typing from assignment, **`DIM` / `AS`**, and **`Namespace.Method`** calls. Some **registry keys** in `commands.json` still contain `$` for legacy compatibility; new examples and docs should use the **documented API names** (e.g. `String.*`, `STR`, `FORMAT` patterns in [API_CONSISTENCY.md](docs/API_CONSISTENCY.md)), not Blitz punctuation.

## Mandatory API Principles

1. **Standard Verbs:** Always use `CREATE` for construction. `MAKE` is deprecated.
2. **Standard Properties:** Always use `SETPOS` for position. `SETPOSITION` is deprecated.
3. **Chainable Setters:** All setter methods (`.pos()`, `.rot()`, `.scale()`, `.col()`, `.alpha()`) must return the handle to allow chaining.
4. **Universal Methods:** Every spatial handle type (Model, Camera, Light, Body, Sprite) MUST implement:
   - `.pos(x, y, z)` / `.pos()` (getter)
   - `.rot(p, y, r)` / `.rot()` (getter)
   - `.scale(sx, sy, sz)` / `.scale()` (getter)
   - `.free()` (canonical cleanup)
5. **Short Parameter Lists:** Prefer chainable setters over constructors with more than 4 positional arguments.

## Easy Mode (Compatibility Layer)

Easy Mode (`CreateCamera()`, `PositionEntity()`, etc.) is provided as a thin wrapper for compatibility with legacy Blitz-style code. It is NOT the primary API and is not recommended for new high-performance games.

```basic
' Legacy Style (Discouraged)
cam = CreateCamera()
PositionEntity(cam, 0, 5, -10)
```

## Advanced API Patterns

### Builders
For complex objects like `PARTICLE` or `TILEMAP`, use the "Create-Configure-Finalize" pattern:

```basic
fire = PARTICLE.CREATE()
              .texture("fire.png")
              .rate(100)
              .velocity(0, 1, 0)
              .play()
```

## Documentation

- **Docs tree entry point:** [docs/STYLE_GUIDE.md](docs/STYLE_GUIDE.md) points here so links under `docs/` stay stable—**do not** duplicate long examples there.
- **Reference page layout (canonical shape):** Command reference pages should follow the **WAVE pattern** — see [docs/DOC_STYLE_GUIDE.md](docs/DOC_STYLE_GUIDE.md) and the live example [docs/reference/WAVE.md](docs/reference/WAVE.md): `# … Commands`, one-line purpose, **`## Core Workflow`** (short narrative), then each command as **`### \`signature\``** with a short paragraph and optional parameter bullets (\`- \`param\`: …\`), **`---`** between every block, and a closing **`## Full Example`** (one-sentence intro + fenced `basic` with comments). This structure is fixed; naming inside headings/snippets follows registry-first rules below.
- **Registry-first in reference pages:** Use uppercase **`NAMESPACE.ACTION`** in headings and snippets (see [docs/reference/API_CONVENTIONS.md](docs/reference/API_CONVENTIONS.md)); match [docs/API_CONSISTENCY.md](docs/API_CONSISTENCY.md) when in doubt.
- **Canonical verbs:** Prefer **`CREATE`**, **`SETPOS`**, **`FREE`**; note deprecated aliases (`MAKE`, `SETPOSITION`) only when documenting migration.
- **Easy Mode:** Mention dotted facades (`Input.KeyDown`, `CreateCamera`, …) as **compatibility**, not the primary path for new examples (see **Easy Mode** above).
- **Identifiers in prose:** Do not use Blitz-style **`#` / `$` / `?` / `%`** suffixes in new docs—use plain parameter names ([Naming conventions](#naming-conventions)).
- **Cross-links:** Point to [API Standardization Directive](docs/API_STANDARDIZATION_DIRECTIVE.md) for phased rollout and manifest workflow.

