# 3D Drawing Commands

Primitive 3D drawing via Raylib. **Call these between `RENDER.BEGIN3D(cam)` and `RENDER.END3D()`** (equivalent **`CAMERA.BEGIN(cam)`** … **`CAMERA.END()`**) so depth testing and the view/projection matrix are correct.

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern** where sections use **`###`** + **`---`**).

## Core Workflow

**`RENDER.BEGIN3D(cam)`** → **`DRAW3D.*`** (or short globals **`BOX`**, **`GRID3`**, …) → **`RENDER.END3D()`**. Same pass works with **`CAMERA.BEGIN`/`END`** if you prefer that bracketing.

moonBASIC registers each command as **`DRAW3D.*`** and also exposes the same behavior under **`Draw.*`** aliases (for example **`Draw.Grid`** → **`DRAW3D.GRID`**). Use either style; they are equivalent.

### Short global names (easier to type)

These builtins are **aliases** of the matching **`DRAW3D.*`** command (same argument lists). Handy in tight loops; long forms remain supported.

| Short | Same as — meaning |
|-------|---------------------|
| **`BOX`** | **`DRAW3D.CUBE`** — solid axis-aligned box |
| **`BOXW`** | **`DRAW3D.CUBEWIRES`** — wire box |
| **`WIRECUBE`** | Same as **`BOXW`** — Blitz3D **`WireCube`** spelling |
| **`BALL`** | **`DRAW3D.SPHERE`** — solid sphere |
| **`BALLW`** | **`DRAW3D.SPHEREWIRES`** — wire sphere |
| **`GRID3`** | **`DRAW3D.GRID`** — XZ reference grid |
| **`FLAT`** | **`DRAW3D.PLANE`** — horizontal plane patch |
| **`CAP`** | **`DRAW3D.CAPSULE`** — solid capsule |
| **`CAPW`** | **`DRAW3D.CAPSULEWIRES`** — wire capsule |

| **`DRAW3D` / `Draw` alias** | Notes |
|-------------------------|--------|
| **`DRAW3D.GRID`** / **`Draw.Grid`** | 2D name; same as **`DRAW3D.GRID`**. |
| **`DRAW3D.LINE`** / **`Draw.Line3D`** | 3D line segment. |
| **`DRAW3D.POINT`** / **`Draw.Point3D`** | 3D point. |
| **`DRAW3D.BBOX`** / **`Draw.BoundingBox`** | Wire-style name maps to **`DRAW3D.BBOX`**. |

---

## Primitives

### `DRAW3D.CUBE(x, y, z, w, h, d, r, g, b, a)`
Draws an axis-aligned solid box centered at `(x, y, z)`.

### `DRAW3D.SPHERE(x, y, z, radius, r, g, b, a)`
Draws a solid sphere.

---

### `DRAW3D.LINE(x1, y1, z1, x2, y2, z2, r, g, b, a)`
Draws a line segment in world space.

### `DRAW3D.GRID(slices, spacing)`
Draws a reference grid in the XZ plane. `slices`: number of divisions; `spacing`: world units between lines.

---

### `DRAW3D.BILLBOARD(tex, x, y, z, size, r, g, b, a)`
Draws a textured billboard facing the **active 3D camera**. **Must** be called inside **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`** (or **`CAMERA.BEGIN`/`CAMERA.END`**).

### `DRAW3D.RAY(handle, r, g, b, a)`
Draws a debug ray from a **6-element float array handle**: origin `(x, y, z)` then direction `(dx, dy, dz)`.

---

## Full Example

```basic
WINDOW.OPEN(800, 600, "3D primitives")
WINDOW.SETFPS(60)
cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 4, 10)
CAMERA.SETTARGET(cam, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(12, 14, 20)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(20, 1.0)
        DRAW3D.CUBE(0, 0.5, 0, 1, 1, 1, 100, 180, 255, 255)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [DRAW_WRAPPERS.md](DRAW_WRAPPERS.md) — object-style **`DRAWCUBE()`**, **`DRAWSPHERE()`**, … (short methods instead of long **`DRAW3D.*`** argument lists).
- [CAMERA.md](CAMERA.md) — 3D camera setup and picking rays.
- [DRAW2D.md](DRAW2D.md) — 2D drawing (screen space, **`CAMERA2D.*`**).
