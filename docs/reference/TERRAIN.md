# Terrain Commands

**`TERRAIN.*`** holds the **2D height grid** and builds **one mesh per loaded chunk** so the whole world is not a single giant mesh. **`CHUNK.*`** sets streaming distances and answers load-state queries. **[`WORLD.*`](WORLD.md)** drives the stream center and per-frame paging — terrain does not replace **`WORLD.UPDATE`**.

**Performance, loading mode, and chunk-size guidance:** [docs/TERRAIN.md](../TERRAIN.md). **Architecture:** [ARCHITECTURE.md](../../ARCHITECTURE.md) §11 (*Open-world runtime*). **CGO + Raylib** required for mesh build/draw; stubs error without it.

**Draw order:** Typically **sky → terrain → opaque props → water → weather** (see [ARCHITECTURE.md](../../ARCHITECTURE.md)).

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Create terrain with **`TERRAIN.CREATE(...)`** or **`TERRAIN.LOAD(path)`**, set origin/sample size (**`TERRAIN.SETPOS`**, **`TERRAIN.SETCHUNKSIZE`**), fill height (**`TERRAIN.FILLPERLIN`**, **`TERRAIN.FILLFLAT`**, …), configure **`CHUNK.SETRANGE(terrain, load, unload)`**, then each frame set the stream center via **`WORLD.SETCENTER`** / **`WORLD.UPDATE`** (see [WORLD.md](WORLD.md)). Draw with **`TERRAIN.DRAW`** inside **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`**.

---

### `TERRAIN.CREATE(...)` 

Creates procedural heightfield terrain. **`TERRAIN.MAKE`** is a deprecated alias. Several arities exist — see **`commands.json`**.

---

### `TERRAIN.LOAD(path)` 

Loads a heightmap image as terrain. Returns a **handle**.

---

### `TERRAIN.FREE(handle)` 

Frees terrain and chunk meshes.

---

### `TERRAIN.SETPOS(handle, x, y, z)` / `TERRAIN.SETCHUNKSIZE(handle, size)` 

World origin and chunk sample size (see reference for overloads).

---

### `TERRAIN.FILLPERLIN(handle, ...)` / `TERRAIN.FILLFLAT(handle, ...)` 

Procedural or flat height fill.

---

### `TERRAIN.GETHEIGHT(handle, x, z)` / `TERRAIN.GETSLOPE(handle, x, z)` 

Sample height and slope at **XZ**.

---

### `TERRAIN.RAISE(handle, ...)` / `TERRAIN.LOWER(handle, ...)` 

Brush sculpting helpers (see manifest).

---

### `TERRAIN.PLACE(handle, id, x, z)` 

Positions an entity on the terrain surface.

---

### `TERRAIN.SNAPY(handle, id)` 

Snaps an entity to the surface **Y**.

---

### `TERRAIN.DRAW(handle)` 

Renders terrain. Must be inside an active **3D** camera block.

---

### `CHUNK.SETRANGE(handle, load, unload)` 

Sets **load** / **unload** radii for chunk paging (world units).

---

### `CHUNK.GENERATE(handle, ix, iz)` / `CHUNK.COUNT(handle)` / `CHUNK.ISLOADED(handle, ix, iz)` 

Build or query chunk meshes (see manifest).

---

## Full Example

Sketch (stream center omitted; **`cam`** is your active 3D camera handle):

```basic
terrain = TERRAIN.CREATE(128, 128, 4.0)
TERRAIN.FILLPERLIN(terrain, 0.12, 42)
CHUNK.SETRANGE(terrain, 120.0, 200.0)
WORLD.SETCENTER(0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    WORLD.UPDATE(dt)
    RENDER.CLEAR(25, 35, 45)
    ; cam = your active CAMERA.* handle
    RENDER.BEGIN3D(cam)
        TERRAIN.DRAW(terrain)
    RENDER.END3D()
    RENDER.FRAME()
WEND

TERRAIN.FREE(terrain)
```

**Common mistake:** Confusing **chunk grid indices** with **world XZ** — streaming uses world position; chunk **i, j** are grid addresses.

---

## See also

- [WORLD.md](WORLD.md) — **`WORLD.SETCENTER`**, **`WORLD.UPDATE`**, **`WORLD.PRELOAD`**
- [WATER.md](WATER.md) — water plane vs terrain height
- [SCATTER.md](SCATTER.md) — foliage on terrain height
