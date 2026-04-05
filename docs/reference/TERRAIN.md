# Terrain and chunks (`TERRAIN.*`, `CHUNK.*`)

## How this subsystem fits the open-world stack

This module holds the **2D height grid** and builds **one mesh per loaded chunk** so the whole world is not a single giant draw. **[`WORLD.*`](WORLD.md)** moves the **stream center** (usually camera/player XZ) and calls into terrain so chunks **inside** the load radius get meshes and chunks **outside** the unload radius can be released. **`CHUNK.*`** adjusts streaming distances and answers questions like how many chunks are loaded.

For the **architecture-level** story (terrain vs world, draw order, handles, navigation), read **§11** in [ARCHITECTURE.md](../../ARCHITECTURE.md) (*Open-world runtime* and *Conceptual overview*).

Heightfield terrain with **chunked** mesh generation (Raylib **`GenMeshHeightmap`**), streaming driven by [`WORLD.*`](WORLD.md) and the terrain module’s internal center. **CGO + Raylib** required for draw and mesh build; without CGO, stubs report an error.

**Draw order:** Typically **sky → terrain → opaque props → water → weather** (see [ARCHITECTURE.md](../../ARCHITECTURE.md)).

---

## `Terrain.Make(worldW, worldH [, cellSize#])` → handle

Creates a terrain object: `worldW` × `worldH` height samples in world units (integer dimensions ≥ 2). Optional **`cellSize#`** scales world spacing per sample (default **1**). Internal chunk size defaults to **64**; change with **`TERRAIN.SETCHUNKSIZE`**.

**Returns:** Terrain handle (`TagTerrain`).

**Common mistake:** Passing tiny dimensions (e.g. 2×2) — you get almost no usable surface; use hundreds of samples for open areas.

---

## `Terrain.Free(terrain)`

Frees the terrain and GPU meshes for loaded chunks.

---

## `Terrain.SetPos(terrain, x#, y#, z#)`

Offsets the terrain origin in world space for drawing.

---

## `Terrain.SetChunkSize(terrain, size)`

Sets the edge length in **height samples** per chunk (must match streaming expectations; affects mesh granularity).

---

## `Terrain.FillPerlin(terrain, scale#, amplitude#)`

Fills heights with layered value noise (implementation-defined seed). **`scale#`** controls feature size; **`amplitude#`** vertical range.

---

## `Terrain.FillFlat(terrain, height#)`

Sets every sample to **`height#`**.

---

## `Terrain.GetHeight(terrain, x#, z#)` → float

Bilinear height at world XZ (clamped to valid range).

---

## `Terrain.GetSlope(terrain, x#, z#)` → float

Approximate slope angle in degrees at XZ.

---

## `Terrain.Raise` / `Terrain.Lower(terrain, x#, z#, radius#, delta#)`

Brush edit: raise or lower height within **`radius#`** by **`delta#`** per call (used for sculpting).

---

## `Terrain.Draw(terrain)`

Draws all **loaded** chunk meshes for the current streaming state. Rebuilds meshes on the **main thread** when chunks load or heights change.

---

## `CHUNK.*` — chunk queries

| Command | Role |
|--------|------|
| `Chunk.Generate(terrain, cx, cz)` | Ensures the chunk at grid **(cx, cz)** is built (if in range). |
| `Chunk.Count(terrain)` | Number of chunks currently holding meshes. |
| `Chunk.SetRange(terrain, loadDist, unloadDist)` | World-unit distances from stream center for load/unload. |
| `Chunk.IsLoaded(terrain, cx, cz)` | Whether that chunk slot has a mesh. |

**Common mistake:** Confusing **chunk indices** with **world XZ** — use terrain/world docs: streaming uses world position; chunk indices are grid addresses.

---

## See also

- [WORLD.md](WORLD.md) — `World.SetCenter`, `World.Update`, preload
- [WATER.md](WATER.md) — water plane vs terrain height
- [SCATTER.md](SCATTER.md) — foliage on terrain height
