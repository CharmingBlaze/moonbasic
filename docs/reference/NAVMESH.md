# Navigation mesh and terrain (integration)

Open-world **heightfields** are **not** automatically converted into the grid nav system. Use **`NAV.*`**, **`PATH.*`**, and **`NAVAGENT.*`** from [`runtime/mbnav`](../../runtime/mbnav) as documented in [NAV_AI.md](NAV_AI.md).

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**Create** a nav object with **`NAV.MAKE()`**, define the **grid** with **`NAV.SETGRID(...)`** to match world scale, optionally attach **`NAV.ADDTERRAIN(nav, modelHandle)`** (AABB walkability / height), add **`NAV.ADDOBSTACLE`** as needed, then **`NAV.BUILD(nav)`** before **`NAV.FINDPATH`** or **`NAVAGENT`** moves. For pure heightfields without a proxy model, sample **`TERRAIN.GETHEIGHT`** and feed walkability into your own marking logic or obstacles — see [NAV_AI.md](NAV_AI.md).

---

### `NAV.MAKE()`

Creates a navigation handle. See [NAV_AI.md](NAV_AI.md) for return type and pairing with **`NAV.FREE`**.

---

### `NAV.SETGRID(...)`

Defines grid resolution and extents — match your world scale (see [NAV_AI.md](NAV_AI.md)).

---

### `NAV.ADDTERRAIN(nav, modelHandle)`

Uses a **model** axis-aligned bounding box to mark walkable **XZ** cells and ground height.

---

### `NAV.ADDOBSTACLE(...)`

Registers blockers (see manifest / [NAV_AI.md](NAV_AI.md)).

---

### `NAV.BUILD(nav)`

Bakes the nav mesh **after** grid + obstacles are configured.

---

### `NAV.FINDPATH(...)`

Path query — full signature in [NAV_AI.md](NAV_AI.md).

---

## No duplicate APIs

Do **not** introduce a second **`NAV.FINDPATH`** or conflicting **`PATH.*`** names — extend **`mbnav`** or add a **new prefix** (for example **`NAVTERRAIN.*`**) if a dedicated bake pipeline is added later.

---

## Full Example

See runnable patterns and full API tables in [NAV_AI.md](NAV_AI.md). Sketch:

```basic
nav = NAV.MAKE()
NAV.SETGRID(nav, ...)
; NAV.ADDTERRAIN(nav, groundModel)
NAV.BUILD(nav)
; ... NAV.FINDPATH or NAVAGENT ... 
```

Rebuild nav when **streaming** changes walkable geometry if you bake from chunk meshes — see [WORLD.md](WORLD.md).

---

## See also

- [NAV_AI.md](NAV_AI.md) — full **`NAV.*`** / **`PATH.*`** / **`NAVAGENT.*`** tables
- [WORLD.md](WORLD.md) — streaming vs static nav rebuilds
