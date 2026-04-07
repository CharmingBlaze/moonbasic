# MoonBASIC command set (engine design)

**Blitz3D spirit · DBPro power · MoonBASIC simplicity**

This folder is the **canonical, human-oriented** view of what moonBASIC provides: short names in these tables are **teaching names**; the runtime usually exposes **`NAMESPACE.NAME`** (see [LANGUAGE.md](../../LANGUAGE.md)). The **Implementation** column gives the real registry keys to search in [API_CONSISTENCY.md](../../API_CONSISTENCY.md).

---

## Memory model (read this first)

| Kind | What it is | How to free |
|------|------------|-------------|
| **VM heap handle** | `KindHandle` — cameras, textures, arrays, wrappers, … | **`FREE`** / **`ERASE`** per handle, or **`FREE.ALL`** / **`ERASE ALL`** for a full reset (see [MEMORY.md](../../MEMORY.md)). |
| **Entity id** (`ENTITY.*`) | Integer id into the entity store — **not** a VM heap slot | **`ENTITY.FREE`**, **`ENTITY.CLEARSCENE`** — Raylib models/meshes unloaded in engine order. |
| **Immediate draws** (`DRAW.*`) | No persistent GPU object per call | No handle to free for a single `DRAW.LINE`-style call. |
| **Globals + pen state** | moonBASIC has **no hidden 2D pen** | Pass **`r,g,b,a`** on each draw or keep locals. |

**Rule of thumb:** If a function **returns a handle**, something must **`FREE`** it unless the docs say “borrowed” or “singleton”.

---

## Section index

| Section | File |
|---------|------|
| Program / runtime | [runtime.md](runtime.md) |
| **Blitz flat API** (globals: `GRAPHICS`, `PLOT`, `CREATECUBE`, …) | [blitz-engine.md](blitz-engine.md) |
| Graphics (window, clear, present) | [graphics.md](graphics.md) |
| 2D drawing | [draw2d.md](draw2d.md) |
| 3D world (ambient, fog, wireframe) | [world3d.md](world3d.md) |
| Camera | [camera.md](camera.md) |
| Lights | [lights.md](lights.md) |
| Entities | [entities.md](entities.md) |
| Mesh / surface | [mesh-surface.md](mesh-surface.md) |
| Textures | [textures.md](textures.md) |
| Sprites (2D) | [sprites.md](sprites.md) |
| Sound | [sound.md](sound.md) |
| Input | [input.md](input.md) |
| Files | [files.md](files.md) |
| Math | [math.md](math.md) |
| MoonBASIC QOL (helpers) | [qol-moonbasic.md](qol-moonbasic.md) |
| Game — English / Blitz-style (PLAYER2D, OrbitCamera, …) | [game-english.md](game-english.md) |
| Physics 2D (Box2D) | [physics-2d.md](physics-2d.md) |
| Physics 3D (Jolt) | [physics-3d.md](physics-3d.md) |
| Networking (ENet) | [network-enet.md](network-enet.md) |
| Networking helpers | [network-helpers.md](network-helpers.md) |

---

## Related

- Legacy mappings: [Blitz index](../BLITZ_COMMAND_INDEX.md), [DBPro folder](../dbpro/README.md)  
- QOL shortcuts module: [QOL.md](../QOL.md), [GAMEHELPERS.md](../GAMEHELPERS.md)  
- Narrative references: [PHYSICS2D.md](../PHYSICS2D.md), [PHYSICS3D.md](../PHYSICS3D.md), [NETWORK.md](../NETWORK.md)  
- Registry source: [`compiler/builtinmanifest/commands.json`](../../../compiler/builtinmanifest/commands.json)  
