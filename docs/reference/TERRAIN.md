# Terrain and chunks (`TERRAIN.*`, `CHUNK.*`)

**Performance, loading mode, and chunk-size guidance:** [docs/TERRAIN.md](../TERRAIN.md).

## How this subsystem fits the open-world stack

This module holds the **2D height grid** and builds **one mesh per loaded chunk** so the whole world is not a single giant draw. **[`WORLD.*`](WORLD.md)** moves the **stream center** (usually camera/player XZ) and calls into terrain so chunks **inside** the load radius get meshes and chunks **outside** the unload radius can be released. **`CHUNK.*`** adjusts streaming distances and answers questions like how many chunks are loaded.

For the **architecture-level** story (terrain vs world, draw order, handles, navigation), read **§11** in [ARCHITECTURE.md](../../ARCHITECTURE.md) (*Open-world runtime* and *Conceptual overview*).

Heightfield terrain with **chunked** mesh generation (Raylib **`GenMeshHeightmap`**), streaming driven by [`WORLD.*`](WORLD.md) and the terrain module’s internal center. **CGO + Raylib** required for draw and mesh build; without CGO, stubs report an error.

**Draw order:** Typically **sky → terrain → opaque props → water → weather** (see [ARCHITECTURE.md](../../ARCHITECTURE.md)).

---

### `Terrain.Load(path)`
Loads heightmap image as terrain. Returns a **handle**.

### `Terrain.Free(handle)`
Frees terrain and all chunk meshes.

---

### `Terrain.GetHeight(handle, x, z)`
Returns the height at world coordinates.

### `Terrain.Place(handle, id, x, z)`
Positions an entity on the terrain surface.

### `Terrain.SnapY(handle, id)`
Snaps an entity to the surface Y.

---

### `Chunk.SetRange(handle, load, unload)`
Sets paging distances for world chunks.

### `Terrain.Draw(handle)`
Renders the terrain. Must be inside `Camera.Begin()`.

**Common mistake:** Confusing **chunk indices** with **world XZ** — use terrain/world docs: streaming uses world position; chunk indices are grid addresses.

---

## See also

- [WORLD.md](WORLD.md) — `World.SetCenter`, `World.Update`, preload
- [WATER.md](WATER.md) — water plane vs terrain height
- [SCATTER.md](SCATTER.md) — foliage on terrain height
