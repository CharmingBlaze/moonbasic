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

### `Draw3D.Grid(slices, spacing#)`

Reference grid in the XZ plane. `slices`: number of divisions; `spacing#`: world units between lines.

### `Draw3D.Line(x1#, y1#, z1#, x2#, y2#, z2#, r, g, b, a)`

Line segment in world space; `r,g,b,a` are 0–255.

### `Draw3D.Point(x#, y#, z#, r, g, b, a)`

Single point (pixel) in 3D.

### `Draw3D.Sphere(x#, y#, z#, radius#, r, g, b, a)`

Solid sphere.

### `Draw3D.SphereWires(x#, y#, z#, radius#, rings, slices, r, g, b, a)`

Wireframe sphere; `rings` and `slices` are integer segment counts.

### `Draw3D.Cube(x#, y#, z#, w#, h#, d#, r, g, b, a)`

Axis-aligned solid cube centered at `(x,y,z)` with size `(w,h,d)`.

### `Draw3D.CubeWires(...)`

Same arguments as `Draw3D.Cube`; wireframe box.

### `Draw3D.Cylinder(x#, y#, z#, rTop#, rBot#, h#, slices, r, g, b, a)`

Solid cylinder; `rTop#` / `rBot#` are top and bottom radii.

### `Draw3D.CylinderWires(...)`

Same arity as `Draw3D.Cylinder`.

### `Draw3D.Capsule(sx#, sy#, sz#, ex#, ey#, ez#, radius#, slices, rings, r, g, b, a)`

Solid capsule between start `(sx,sy,sz)` and end `(ex,ey,ez)`.

### `Draw3D.CapsuleWires(...)`

Same arity as `Draw3D.Capsule`.

### `Draw3D.Plane(x#, y#, z#, width#, depth#, r, g, b, a)`

Horizontal plane through `(x,y,z)` with size `(width, depth)`.

### `Draw3D.BBox(minx#, miny#, minz#, maxx#, maxy#, maxz#, r, g, b, a)`

Axis-aligned bounding box wireframe.

### `Draw3D.Ray(rayArray, r, g, b, a)`

`rayArray` is a **6-element float array handle**: origin `(x,y,z)` then direction `(dx,dy,dz)`. Use with `Camera.GetRay` / `Camera.GetViewRay` or your own array.

### `Draw3D.Billboard(tex, x#, y#, z#, size#, r, g, b, a)`

Textured billboard facing the **active 3D camera**. **Must** be called inside `Camera.Begin` / `Camera.End` (the runtime needs the current camera).

### `Draw3D.BillboardRec(tex, srcX#, srcY#, srcW#, srcH#, x#, y#, z#, w#, h#, r, g, b, a)`

Billboard with a source rectangle on the texture and destination size `(w#, h#)` in world units. Same active-camera requirement as `Draw3D.Billboard`.

---

## See also

- [DRAW_WRAPPERS.md](DRAW_WRAPPERS.md) — object-style **`DRAWCUBE()`**, **`DRAWSPHERE()`**, … (short methods instead of long **`DRAW3D.*`** argument lists).
- [CAMERA.md](CAMERA.md) — 3D camera setup and picking rays.
- [DRAW2D.md](DRAW2D.md) — 2D drawing (screen space, `Camera2D`).
