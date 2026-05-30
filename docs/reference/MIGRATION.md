# BlitzBASIC → moonBASIC migration

Commands marked **`stub`** in `commands.json` are rejected at **compile time** with a clear message — they will not surprise you at runtime.

## Unavailable in this release (use alternatives)

| Command | Use instead |
|---------|-------------|
| `PHYSICS3D.DEBUGDRAW` | Omit debug overlay until Jolt wireframe draw lands |
| `PHYSICS.SPHERECAST` / `PHYSICS.BOXCAST` | `PHYSICS3D.RAYCAST` |
| `PHYSICS.ENABLE` / `PHYSICS.DISABLE` | Remove body or use ghost/trigger flags |
| `JOINT3D.FIXED` / `JOINT3D.SLIDER` / `JOINT3D.CONE` | `JOINT3D.HINGE` where applicable |
| `PHYSICS2D.ONCOLLISION` / `PHYSICS2D.PROCESSCOLLISIONS` | Implemented — callbacks run after each `PHYSICS2D.STEP` |
| `WORLD.SETREFLECTION` | Deferred — no reflection probes yet |
| `GAME.BURSTSPAWN` | Implemented — particle burst at (x,y,z) |
| `GAME.SPRITETILEBRIDGE` | Collision math + sprite bounds in script |
| `CREATESPRITE3D` | `ENTITY.CREATECUBE` + `ENTITY.TEXTURE` |
| `ENTITY.INSTANCE` | `MODEL.MAKEINSTANCED` |
| `LEVEL.OPTIMIZE` | `MODEL.MAKEINSTANCED` for batches |
| `LEVEL.APPLYPHYSICS` | `BODY3D.*` + `COMMIT` per mesh (see PHYSICS3D.md) |
| `MESH.GENERATELOD` / `MESH.GENERATELODCHAIN` | `MODEL.LOADLOD` (file-based LOD) |
| `MODEL.GETCHILD` | `ENTITY.GETCHILD` |
| `CREATEMIRROR` | Not available — planar reflections deferred |
| `FITMESH` / `FLIPMESH` / `UPDATENORMALS` | Import meshes with tools; use `ENTITY.LOADMESH` |
| `ADDSURFACE` | `MESH.MAKECUSTOM` / `ENTITY.LOADMESH` |

## Partial APIs

| Command | Notes |
|---------|--------|
| `PICK.RADIUS` | Ray pick works with radius `0`; non-zero sphere pick is not implemented yet |
| `ENTITY.ANIMATE` / cross-fade | **`ENTITY.CROSSFADE` / `ENTITY.TRANSITION`** — dual-pose bone blend over a duration (see ANIMATION_3D.md) |
| `PHYSICS2D.*` (other) | Box2D bodies/joints are implemented; collision callbacks run after each `PHYSICS2D.STEP` |

## Language gotchas

- **Namespace shadowing:** a variable named `time` shadows `TIME.*`. The compiler warns when a variable name matches a built-in namespace.
- **Game loop condition:** `NOT a OR b` parses as `(NOT a) OR b`. Use `NOT (a OR b)` — the compiler warns on the common mistake.

Regenerate stub annotations after manifest edits: `go run ./tools/annotate_stubs` (maintainers).
