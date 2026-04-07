# Entity commands

moonBASIC **entities** are lightweight **integer ids** (physics/visual game objects) with **`ENTITY.*`** builtins. **CGO** builds link the full implementation; see the registry in [`runtime/mbentity/`](../../runtime/mbentity/).

## Quick links

- **Blitz-style names** (`PositionEntity`, `CreateSphere`, …) are mapped under **`ENTITY.POSITIONENTITY`**, **`ENTITY.CREATESPHERE`**, etc. — see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go) and the **[Blitz command index](BLITZ_COMMAND_INDEX.md)**.
- **Dot-syntax handles** (`cube.Pos`, `sphere.Turn`) use **`ENTITYREF`** from **`CUBE()`** / **`SPHERE()`** — [BLITZ3D.md](BLITZ3D.md).
- **Scene save/load / clear** — [BLITZ2025.md](BLITZ2025.md), **`ENTITY.SAVESCENE`**, **`ENTITY.LOADSCENE`**, **`ENTITY.CLEARSCENE`**.

## Reference tables

- **[API_CONSISTENCY.md](../API_CONSISTENCY.md)** — search for **`ENTITY.`** for every overload and arity.
- **[GAMEHELPERS.md](GAMEHELPERS.md)** — movement, landing, camera follow.
