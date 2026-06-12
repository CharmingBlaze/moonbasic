# moonBASIC documentation

Welcome to the moonBASIC documentation. Start here whether you browse on [GitHub](https://github.com/CharmingBlaze/moonbasic-compiler) or offline in the repo.

---

## Quick start

| I want to… | Read |
|------------|------|
| **Start from zero (install + why each command)** | [BEGIN_HERE.md](BEGIN_HERE.md) |
| **Deep guides (all 40 systems + net)** | [systems/GUIDES.md](systems/GUIDES.md) — 24 topic guides |
| Install and run a game | [GETTING_STARTED.md](GETTING_STARTED.md) |
| **List every beginner-system command** | [systems/COMMAND_REGISTRY.md](systems/COMMAND_REGISTRY.md) |
| Learn the language | [LANGUAGE.md](LANGUAGE.md) · [FIRST_HOUR.md](FIRST_HOUR.md) |
| Build a game loop | [PROGRAMMING.md](PROGRAMMING.md) |
| Browse commands by **system** (APP, ENTITY, …) | [systems/README.md](systems/README.md) |
| Look up a **namespace** (deep reference) | [COMMANDS.md](COMMANDS.md) → [reference/](reference/) |
| See every registered command | [API_CONSISTENCY.md](API_CONSISTENCY.md) |

---

## Documentation map

### Game systems (beginner API)

The **[systems guide](systems/README.md)** documents the 40 systems from the moonBASIC foundation checklist: APP, RENDER, ENTITY, physics, audio, saves, tooling, and more. Each page follows [DOCUMENTATION_STYLE_GUIDE.md](DOCUMENTATION_STYLE_GUIDE.md).

| Part | File | Systems covered |
|------|------|-----------------|
| 1 | [systems/01-CORE.md](systems/01-CORE.md) | APP, RENDER, SCENE, ENTITY |
| 2 | [systems/02-CAMERA-LIGHT.md](systems/02-CAMERA-LIGHT.md) | CAMERA, LIGHT |
| 3 | [systems/03-ASSETS.md](systems/03-ASSETS.md) | MESH, MODEL, MATERIAL, TEXTURE, ASSET |
| 4 | [systems/04-INPUT.md](systems/04-INPUT.md) | INPUT, ACTION |
| 5 | [systems/05-PHYSICS.md](systems/05-PHYSICS.md) | PHYSICS, BODY, COLLISION, PICK |
| 6 | [systems/06-AUDIO.md](systems/06-AUDIO.md) | AUDIO, AUDIO3D |
| 7 | [systems/07-2D-WORLD.md](systems/07-2D-WORLD.md) | SPRITE, TILEMAP, TERRAIN, PARTICLE, ANIMATION |
| 8 | [systems/08-UI-TEXT.md](systems/08-UI-TEXT.md) | UI, FONT, TEXT |
| 9 | [systems/09-DATA.md](systems/09-DATA.md) | SAVE, FILE, JSON, MATH, VEC3 |
| 10 | [systems/10-DEBUG-TIMER.md](systems/10-DEBUG-TIMER.md) | DEBUG, ERROR, TIMER |
| 11 | [systems/11-TOOLING.md](systems/11-TOOLING.md) | PROJECT, PACKAGE, MODULE, HELP, TEST, TEMPLATE |

All eleven parts follow [DOCUMENTATION_STYLE_GUIDE.md](DOCUMENTATION_STYLE_GUIDE.md).

**Foundation example:** [../examples/foundation/main.mb](../examples/foundation/main.mb)

### Language and workflow

- [LANGUAGE.md](LANGUAGE.md) — syntax, types, case insensitivity
- [PROGRAMMING.md](PROGRAMMING.md) — game loop, 2D/3D
- [MEMORY.md](MEMORY.md) — handles, `FREE`, `ERASE ALL`
- [ROADMAP.md](ROADMAP.md) — shipped vs planned

### Engine reference (namespaces)

- [COMMANDS.md](COMMANDS.md) — topic index
- [reference/](reference/) — per-namespace pages (`RENDER.md`, `ENTITY.md`, …)
- [COMMAND_AUDIT.md](COMMAND_AUDIT.md) — namespace → doc file map

### Contributors

- [DEVELOPER.md](DEVELOPER.md) · [BUILDING.md](BUILDING.md) · [CONTRIBUTING.md](../CONTRIBUTING.md)
- [DOCUMENTATION_STYLE_GUIDE.md](DOCUMENTATION_STYLE_GUIDE.md) — **how to write docs**
- [API_CONSISTENCY.md](API_CONSISTENCY.md) — every registered command (manifest listing)

---

## Conventions

- Commands are **case-insensitive**: `app.open`, `APP.OPEN`, and `App.Open` are the same.
- Canonical registry names use `NAMESPACE.COMMAND` (uppercase).
- Checklist names like `APP.*` are **aliases** documented in systems pages; internals may use `WINDOW.*`, `TIME.*`, etc.

---

## See also

- [Repository README](../README.md)
- [Examples](../examples/README.md)
- [FINAL_POLISH_SYSTEMS.md](FINAL_POLISH_SYSTEMS.md) — checklist status summary
