# 3D Drawing Commands

Primitive 3D drawing via Raylib. **Call these between `Camera.Begin(cam)` and `Camera.End()`** so depth testing and the view/projection matrix are correct.

moonBASIC registers each command as **`Draw3D.*`** and also exposes the same behavior under **`Draw.*`** aliases (for example `Draw.Grid` → `DRAW3D.GRID`). Use either style; they are equivalent.

### Short global names (easier to type)

These builtins are **aliases** of the matching **`DRAW3D.*`** command (same argument lists). Handy in tight loops; long forms remain supported.

| Short | Same as — meaning |
|-------|---------------------|
| **`BOX`** | `DRAW3D.CUBE` — solid axis-aligned box |
| **`BOXW`** | `DRAW3D.CUBEWIRES` — wire box |
| **`WIRECUBE`** | Same as **`BOXW`** — Blitz3D **`WireCube`** spelling |
| **`BALL`** | `DRAW3D.SPHERE` — solid sphere |
| **`BALLW`** | `DRAW3D.SPHEREWIRES` — wire sphere |
| **`GRID3`** | `DRAW3D.GRID` — XZ reference grid |
| **`FLAT`** | `DRAW3D.PLANE` — horizontal plane patch |
| **`CAP`** | `DRAW3D.CAPSULE` — solid capsule |
| **`CAPW`** | `DRAW3D.CAPSULEWIRES` — wire capsule |

| `Draw3D` / `Draw` alias | Notes |
|-------------------------|--------|
| `Draw3D.Grid` / `Draw.Grid` | 2D name; same as `Draw3D.GRID`. |
| `Draw3D.Line` / `Draw.Line3D` | 3D line segment. |
| `Draw3D.Point` / `Draw.Point3D` | 3D point. |
| `Draw3D.BoundingBox` / `Draw.BoundingBox` | Wire-style name maps to `DRAW3D.BBOX`. |

---

## Primitives

### `Draw3D.Cube(x, y, z, w, h, d, r, g, b, a)`
Draws an axis-aligned solid box centered at `(x, y, z)`.

### `Draw3D.Sphere(x, y, z, radius, r, g, b, a)`
Draws a solid sphere.

---

### `Draw3D.Line(x1, y1, z1, x2, y2, z2, r, g, b, a)`
Draws a line segment in world space.

### `Draw3D.Grid(slices, spacing)`
Draws a reference grid in the XZ plane. `slices`: number of divisions; `spacing`: world units between lines.

---

### `Draw3D.Billboard(tex, x, y, z, size, r, g, b, a)`
Draws a textured billboard facing the **active 3D camera**. **Must** be called inside `Camera.Begin()` / `Camera.End()`.

### `Draw3D.Ray(handle, r, g, b, a)`
Draws a debug ray from a **6-element float array handle**: origin `(x, y, z)` then direction `(dx, dy, dz)`.

---

## See also

- [DRAW_WRAPPERS.md](DRAW_WRAPPERS.md) — object-style **`DRAWCUBE()`**, **`DRAWSPHERE()`**, … (short methods instead of long **`DRAW3D.*`** argument lists).
- [CAMERA.md](CAMERA.md) — 3D camera setup and picking rays.
- [DRAW2D.md](DRAW2D.md) — 2D drawing (screen space, `Camera2D`).
