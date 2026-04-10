# World streaming (`WORLD.*`)

## Role in the stack

**World** does not own heightmap data. It drives **which part of the terrain** should be resident: you set **`WORLD.SETCENTER`** to the player or camera XZ, then **`WORLD.UPDATE`** once per frame so the bound terrain can load/unload **chunks** according to **`CHUNK.SETRANGE`**. Use **`WORLD.PRELOAD`** for a startup burst so the first view is filled in. See the narrative in [ARCHITECTURE.md](../../ARCHITECTURE.md) §11 (*Conceptual overview*).

The **world manager** ties into the active [`terrain`](TERRAIN.md) module: it updates **stream center**, runs **chunk load/unload** each frame, and exposes **preload** and **ready** queries. It does **not** implement separate `REGION.*` file I/O — that remains future work.

**CGO** required for real terrain streaming; stubs fail without it.

---

## `World.SetCenter(x#, z#)`

Sets the streaming focal point (usually the camera or player XZ). The bound terrain uses this to decide which chunks to load.

---

## `World.SetCenterEntity(entity#)`

Same as **`World.SetCenter`** but takes an **entity id** and uses that entity’s world **X** and **Z**. Prefer this when the player (or camera rig) is already positioned each frame — avoids duplicating **`EntityX` / `EntityZ`** calls. See [LESS_MATH.md](LESS_MATH.md).

---

## `World.Update(dt#)`

Advances streaming for one frame. Call **once per frame** after moving the center. **`dt`** is accepted for API symmetry; the current implementation uses the terrain’s tick path.

---

## `World.StreamEnable(enabled?)`

Enables or disables chunk streaming on the bound terrain.

---

## `World.Preload(terrain, radius)`

Loads chunks in a **Manhattan** or radius-based neighborhood around the current center (implementation: terrain `PreloadTerrain`). Use after **`World.SetCenter`** to avoid pop-in at start.

---

## `World.Status()` → string$

Human-readable status for debugging (implementation-defined).

---

## `World.IsReady(terrain)` → bool

Returns whether the given terrain handle has finished **initial** chunk work relevant to the current stream state (implementation-defined readiness).

---

## `World.SetVegetation(terrain#, billboard#, density#)`

Populates an internal **`SCATTER`** instance with random **XZ** samples over a fixed area and snaps **Y** to **`Terrain.GetHeight`** (same placement rule as **`Scatter.Apply`**). The **billboard** handle is **reserved** for future instanced mesh drawing; today **`Scatter.DrawAll`** uses simple debug spheres unless you extend the scatter renderer.

---

## See also

- [TERRAIN.md](TERRAIN.md) — heightfield and `Chunk.SetRange`
- [openworld_complete.mb](../../testdata/openworld_complete.mb) — integration sample
