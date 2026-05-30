# Scatter and prop commands

**`SCATTER.*`** holds a named set of decorative points sampled on terrain (runtime draws simple markers). **`PROP.*`** records ad-hoc world positions and draws placeholder cubes. **CGO + Raylib** required for draw paths.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — prefer **`CREATE`/`FREE`** where registered.

**Note:** There is no **`SCATTER.UPDATE`** in the current runtime — this module does not animate instances.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**Scatter:** **`SCATTER.CREATE(name)`** → **`SCATTER.APPLY(scatter, terrain, density)`** after you have a valid **terrain** handle → **`SCATTER.DRAWALL`** inside your **3D** camera block each frame → **`SCATTER.FREE`** when done.

**Props:** **`PROP.PLACE`** → **`PROP.DRAWALL`** → **`PROP.FREE`**.

---

## `SCATTER.*`

### `SCATTER.CREATE(name)` 

Allocates a scatter object tagged with **`name`** (string). Returns a **heap handle**.

---

### `SCATTER.FREE(handle)` 

Frees the scatter handle.

---

### `SCATTER.APPLY(scatter, terrain, density)` 

Repopulates samples on the given **terrain** using **`density`** (internal heuristic — start **low** and profile).

---

### `SCATTER.DRAWALL(scatter)` 

Draws scatter markers. Call inside **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`** (or **`CAMERA.BEGIN`** / **`CAMERA.END`**).

---

## `PROP.*`

### `PROP.PLACE(model, x, y, z)` 

Records a prop at **(x, y, z)**. The first argument is **reserved** for a future mesh association; pass **`0`** or **`NIL`** where your build allows.

---

### `PROP.FREE(handle)` 

Frees a stored prop handle.

---

### `PROP.DRAWALL()` 

Draws all registered props (placeholder cubes in the reference implementation).

---

## Full Example

```basic
terrain = TERRAIN.LOAD("heightmap.png")
grass = SCATTER.CREATE("grass")
SCATTER.APPLY(grass, terrain, 0.08)
prop = PROP.PLACE(0, 5.0, 2.0, 5.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    WORLD.SETCENTER(0, 0)
    WORLD.UPDATE(dt)
    RENDER.CLEAR(30, 40, 50)
    RENDER.BEGIN3D(cam)
        TERRAIN.DRAW(terrain)
        SCATTER.DRAWALL(grass)
        PROP.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

PROP.FREE(prop)
SCATTER.FREE(grass)
TERRAIN.FREE(terrain)
```

**Common mistake:** **`SCATTER.APPLY`** with very high **density** — cost scales with instance count; start low and profile.

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `SCATTER.MAKE(...)` | Deprecated alias of `SCATTER.CREATE`. |

---

## See also

- [TERRAIN.md](TERRAIN.md) — height sampling / **`TERRAIN.LOAD`**
- [MODEL.md](MODEL.md) — when mesh-backed props are added
