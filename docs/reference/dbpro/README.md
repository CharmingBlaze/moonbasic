# DarkBASIC Professional (DBPro) → moonBASIC

moonBASIC is **not** a DarkBASIC runtime. This folder maps **DBPro-style command names** to **implemented** moonBASIC APIs (usually `NAMESPACE.NAME`). Identifiers are case-insensitive — see [LANGUAGE.md](../../LANGUAGE.md).

**Sources of truth**

- Registry: [`compiler/builtinmanifest/commands.json`](../../../compiler/builtinmanifest/commands.json)
- Human list: [API_CONSISTENCY.md](../../API_CONSISTENCY.md)
- Similar legacy map: [Blitz command index](../BLITZ_COMMAND_INDEX.md)

**Legend (in each table)**

| Mark | Meaning |
|------|---------|
| ✓ | Close equivalent exists under the suggested name(s). |
| ≈ | Same role; different arguments, workflow, or engine feature set. |
| — | No direct command; use the note or build a helper in script. |

---

## Sections

| # | Topic | File |
|---|--------|------|
| 1 | Objects / 3D engine (create, transform, appearance, collision) | [01-objects-3d.md](01-objects-3d.md) |
| 2 | Mesh / limb | [02-mesh-limb.md](02-mesh-limb.md) |
| 3 | Camera | [03-camera.md](03-camera.md) |
| 4 | Lights | [04-lights.md](04-lights.md) |
| 5 | Image | [05-image.md](05-image.md) |
| 6 | Sprite | [06-sprite.md](06-sprite.md) |
| 7 | Sound | [07-sound.md](07-sound.md) |
| 8 | Input | [08-input.md](08-input.md) |
| 9 | File I/O | [09-file.md](09-file.md) |
| 10 | Math | [10-math.md](10-math.md) |

---

## Quick mental model

| DBPro idea | moonBASIC |
|------------|-----------|
| Integer **object id** + **MAKE OBJECT** | Often **`ENTITY.*`** (primitives / file models) or **`MODEL.*`** (loaded model handle) — see [MODEL.md](../MODEL.md), [ENTITY.md](../ENTITY.md). |
| **Limb** hierarchy | Partially: **`MODEL.LIMBCOUNT`**, **`MODEL.CHILDCOUNT`**, parenting — not a full DBPro limb stack. |
| **Camera / Light** handles | **`CAMERA.*`**, **`LIGHT.*`** — [CAMERA.md](../CAMERA.md), [LIGHT.md](../LIGHT.md). |
| **2D sprite** id | **`SPRITE.*`** (Raylib-backed) — [SPRITE.md](../SPRITE.md). |
