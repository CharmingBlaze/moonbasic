# World Commands

**`WORLD.*`** does **not** own heightmap data. It drives **which part of the terrain** stays resident: set **`WORLD.SETCENTER`** (player or camera **XZ**), call **`WORLD.UPDATE(dt)`** each frame so the bound terrain can load/unload chunks per **`CHUNK.SETRANGE`**, and use **`WORLD.PRELOAD`** for a startup burst. See [ARCHITECTURE.md](../../ARCHITECTURE.md) §11 (*Conceptual overview*).

The world manager ties into the active [terrain](TERRAIN.md) module. It does **not** implement separate **`REGION.*`** file I/O — that remains future work. **CGO** required for real terrain streaming; stubs fail without it.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Optionally **`WORLD.SETUP()`** or **`WORLD.SETUP(gravity)`** to initialize physics-world defaults (see [PHYSICS3D.md](PHYSICS3D.md) for Jolt). Each frame: **`WORLD.SETCENTER`** (or **`WORLD.SETCENTERENTITY`**), then **`WORLD.UPDATE(dt)`**. Tune streaming with **`WORLD.STREAMENABLE`**, warm caches with **`WORLD.PRELOAD`**, and poll **`WORLD.ISREADY`** / **`WORLD.STATUS`** while loading.

---

### `WORLD.SETUP()` / `WORLD.SETUP(gravity)` 

Initializes the physics world with default **Y** gravity (**-9.81**) or a custom **`gravity`** value. Expects **`PHYSICS3D.START()`** / Jolt when using full **3D** physics on desktop **Windows** and **Linux** with **CGO**.

---

### `WORLD.UPDATE(dt)` 

Updates world streaming and related state. **`dt`** is seconds (use **`TIME.DELTA()`** in the game loop).

---

### `WORLD.SETCENTER(x, z)` / `WORLD.SETCENTERENTITY(id)` 

Sets the streaming focal point for chunk paging.

---

### `WORLD.PRELOAD(terrain, radius)` 

Forces initial chunk loading around the current center.

---

### `WORLD.STREAMENABLE(toggle)` 

Enables or disables automatic chunk paging.

---

### `WORLD.ISREADY(terrain)` 

Returns **`TRUE`** when initial chunk work is complete enough to play.

---

### `WORLD.STATUS()` 

Returns a debug status string.

---

### `WORLD.SETVEGETATION(terrain, billboard, density)` 

Populates an internal **`SCATTER`** sample set over terrain (**`TERRAIN.GETHEIGHT`**) with the given **density**. The **billboard** handle is reserved for future instanced drawing; **`Scatter.DrawAll`** may use simple debug spheres until extended.

---

## Full Example

Integration sample: [`testdata/openworld_complete.mb`](../../testdata/openworld_complete.mb). Minimal sketch:

```basic
WORLD.SETUP()
terrain = TERRAIN.CREATE(64, 64, 8.0)
TERRAIN.FILLFLAT(terrain, 0.0)
CHUNK.SETRANGE(terrain, 80.0, 140.0)
WORLD.PRELOAD(terrain, 3)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    ; camX, camZ = focal point for streaming (e.g. from camera or player)
    WORLD.SETCENTER(camX, camZ)
    WORLD.UPDATE(dt)
    ; ... draw terrain, entities ...
WEND

TERRAIN.FREE(terrain)
```

---

## See also

- [TERRAIN.md](TERRAIN.md) — heightfield and **`CHUNK.SETRANGE`**
- [`testdata/openworld_complete.mb`](../../testdata/openworld_complete.mb)
