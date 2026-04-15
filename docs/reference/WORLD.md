# World streaming (`WORLD.*`)

## Role in the stack

**World** does not own heightmap data. It drives **which part of the terrain** should be resident: you set **`WORLD.SETCENTER`** to the player or camera XZ, then **`WORLD.UPDATE`** once per frame so the bound terrain can load/unload **chunks** according to **`CHUNK.SETRANGE`**. Use **`WORLD.PRELOAD`** for a startup burst so the first view is filled in. See the narrative in [ARCHITECTURE.md](../../ARCHITECTURE.md) §11 (*Conceptual overview*).

The **world manager** ties into the active [`terrain`](TERRAIN.md) module: it updates **stream center**, runs **chunk load/unload** each frame, and exposes **preload** and **ready** queries. It does **not** implement separate `REGION.*` file I/O — that remains future work.

**CGO** required for real terrain streaming; stubs fail without it.

---

### `World.Setup([gravity#])`
**Easy Mode** physics entry point. Initializes global gravity and expects **`PHYSICS3D.START()`** / Jolt when using full **3D** physics (see [PHYSICS3D.md](PHYSICS3D.md)). On desktop **Windows and Linux** with **CGO + Jolt**, behavior matches the native world step. Default gravity is `-9.81`.

---

### `World.Update(dt)`
Updates world streaming and entity spatial SoA.

### `World.SetCenter(x, z)` / `World.SetCenterEntity(id)`
Sets the streaming focal point for chunks.

---

### `World.Preload(terrain, radius)`
Forces initial chunk loading around the center.

### `World.StreamEnable(toggle)`
Enables or disables automatic chunk paging.

---

### `World.IsReady(terrain)`
Returns `TRUE` if initial chunk work is complete.

### `World.Status()`
Returns a status string for debugging.

---

## `World.SetVegetation(terrain, billboard, density)`

Populates an internal **`SCATTER`** instance with random **XZ** samples over a fixed area and snaps **Y** to **`Terrain.GetHeight`** (same placement rule as **`Scatter.Apply`**). The **billboard** handle is **reserved** for future instanced mesh drawing; today **`Scatter.DrawAll`** uses simple debug spheres unless you extend the scatter renderer.

---

## See also

- [TERRAIN.md](TERRAIN.md) — heightfield and `Chunk.SetRange`
- [openworld_complete.mb](../../testdata/openworld_complete.mb) — integration sample
