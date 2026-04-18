# Prop Commands

Quick scene prop placement helpers. Spawn a mesh at a world position as a reusable prop handle, draw all props in one call, and free individually.

## Core Workflow

1. Load a mesh: `mesh = MESH.LOAD("assets/tree.glb")`.
2. `PROP.PLACE(mesh, x, y, z)` for each instance.
3. Each frame inside `RENDER.BEGIN3D` / `RENDER.END3D`: `PROP.DRAWALL()`.
4. `PROP.FREE(handle)` for individual instances on scene change.

---

## Commands

### `PROP.PLACE(mesh, x, y, z)` 

Places a prop instance of `mesh` at world position `(x, y, z)`. Returns a **prop handle**.

---

### `PROP.DRAWALL()` 

Draws all active prop instances in one pass. Call inside `RENDER.BEGIN3D` / `RENDER.END3D`.

---

### `PROP.FREE(propHandle)` 

Frees a single prop instance.

---

## Full Example

```basic
WINDOW.OPEN(960, 540, "Prop Demo")
WINDOW.SETFPS(60)

cam  = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -14)
CAMERA.SETTARGET(cam, 0, 0, 0)

tree = MESH.LOAD("assets/tree.glb")

FOR i = 1 TO 20
    PROP.PLACE(tree, RNDF(-12, 12), 0, RNDF(-12, 12))
NEXT i

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(80, 120, 60)
    RENDER.BEGIN3D(cam)
        PROP.DRAWALL()
        DRAW.GRID(24, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [ENTITY.md](ENTITY.md) — full entity system
- [SPAWNER.md](SPAWNER.md) — timed entity spawner
