# Ray Commands

Ray creation and intersection tests against boxes, spheres, planes, triangles, meshes, and models.

A **ray** is an origin point + direction vector. Create one, then query results with the `RAY.HIT*` family. Result accessors (`_HIT`, `_DISTANCE`, `_POINTX/Y/Z`, `_NORMALX/Y/Z`) are read after the corresponding test call.

## Core Workflow

1. `RAY.CREATE(ox, oy, oz, dx, dy, dz)` — create a ray (origin + direction).
2. Call a test: `RAY.HITBOX`, `RAY.HITSPHERE`, `RAY.HITPLANE`, `RAY.HITMESH`, `RAY.HITMODEL`, etc.
3. Read results: `RAY.HITBOX_HIT`, `RAY.HITBOX_DISTANCE`, `RAY.HITBOX_POINTX/Y/Z`, `RAY.HITBOX_NORMALX/Y/Z`.
4. `RAY.FREE(ray)` when done.

For camera-based screen picking see `PICK.*` commands — they wrap ray creation from a camera.

---

## Creation & Lifetime

### `RAY.CREATE(ox, oy, oz, dx, dy, dz)` 

Creates a ray with origin `(ox, oy, oz)` and direction `(dx, dy, dz)`. Direction does not need to be normalised. Returns a **ray handle**.

---

### `RAY.FREE(ray)` 

Frees the ray handle.

---

### `RAY.SETPOS(ray, x, y, z)` 

Sets the ray origin. Returns the ray handle for chaining.

---

### `RAY.SETDIR(ray, dx, dy, dz)` 

Sets the ray direction. Returns the ray handle for chaining.

---

### `RAY.GETPOS(ray)` / `RAY.POS(ray)` 

Returns the origin as a `VEC3` handle.

---

### `RAY.GETDIR(ray)` / `RAY.DIR(ray)` 

Returns the direction as a `VEC3` handle.

---

## Box Tests

### `RAY.HITBOX(ray, minX, minY, minZ, maxX, maxY, maxZ)` 

Tests the ray against an axis-aligned box. Returns `TRUE` on hit. Also writes internal result data read by the `RAY.HITBOX_*` accessors below.

---

### `RAY.HITBOX_HIT(ray, minX, minY, minZ, maxX, maxY, maxZ)` 

Returns `TRUE` / `FALSE` — same as `RAY.HITBOX` but explicit.

---

### `RAY.HITBOX_DISTANCE(ray, minX, minY, minZ, maxX, maxY, maxZ)` 

Returns the hit distance from the ray origin.

---

### `RAY.HITBOX_POINTX/Y/Z(ray, minX, minY, minZ, maxX, maxY, maxZ)` 

Returns the X, Y, or Z component of the hit point.

---

### `RAY.HITBOX_NORMALX/Y/Z(ray, minX, minY, minZ, maxX, maxY, maxZ)` 

Returns the X, Y, or Z component of the surface normal at the hit.

---

## Sphere Tests

### `RAY.HITSPHERE(ray, cx, cy, cz, radius)` 

Tests against a sphere at `(cx, cy, cz)` with `radius`. Returns `TRUE` on hit.

---

`RAY.HITSPHERE_HIT` / `RAY.HITSPHERE_DISTANCE` / `RAY.HITSPHERE_POINTX/Y/Z` / `RAY.HITSPHERE_NORMALX/Y/Z` — result accessors (same argument signature as `RAY.HITSPHERE`).

---

## Plane Tests

### `RAY.HITPLANE(ray, nx, ny, nz, d)` 

Tests against an infinite plane defined by normal `(nx, ny, nz)` and distance `d` from origin. Returns `TRUE` on hit.

---

`RAY.HITPLANE_HIT` / `RAY.HITPLANE_DISTANCE` / `RAY.HITPLANE_POINTX/Y/Z` / `RAY.HITPLANE_NORMALX/Y/Z` — result accessors.

---

## Triangle Tests

### `RAY.HITTRIANGLE(ray, ax, ay, az, bx, by, bz, cx, cy, cz)` 

Tests against a triangle defined by three vertices. Returns `TRUE` on hit.

---

`RAY.HITTRIANGLE_HIT` / `RAY.HITTRIANGLE_DISTANCE` / `RAY.HITTRIANGLE_POINTX/Y/Z` / `RAY.HITTRIANGLE_NORMALX/Y/Z` — result accessors.

---

## Mesh Tests

### `RAY.HITMESH(ray, meshHandle, transformHandle)` 

Tests against a CPU-side mesh with a transform matrix applied. Returns `TRUE` on hit. Slower than AABB tests; use for precise picking.

---

`RAY.HITMESH_HIT` / `RAY.HITMESH_DISTANCE` / `RAY.HITMESH_POINTX/Y/Z` / `RAY.HITMESH_NORMALX/Y/Z` — result accessors (same args).

---

## Model Tests

### `RAY.HITMODEL(ray, modelHandle)` 

Tests against a loaded model (entity mesh) without a separate transform. Returns `TRUE` on hit. Alias: `RAY.INTERSECTSMODEL`.

---

`RAY.HITMODEL_HIT` / `RAY.HITMODEL_DISTANCE` / `RAY.HITMODEL_POINTX/Y/Z` / `RAY.HITMODEL_NORMALX/Y/Z` — result accessors (same args).

`RAY.INTERSECTSMODEL_*` — identical aliases for all result accessors.

---

## Full Example

Screen-to-world ray picking — click to find what the mouse points at.

```basic
WINDOW.OPEN(960, 540, "Ray Picking")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 6, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW3D.CUBE(0, 0, 0, 2, 2, 2, 80, 120, 200, 255)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()

    IF MOUSE.PRESSED(0)
        ; build a ray from camera through mouse position
        ray = CAMERA.GETRAY(cam, MOUSE.X(), MOUSE.Y())
        ox = VEC3.X(RAY.GETPOS(ray))
        oy = VEC3.Y(RAY.GETPOS(ray))
        oz = VEC3.Z(RAY.GETPOS(ray))
        dx = VEC3.X(RAY.GETDIR(ray))
        dy = VEC3.Y(RAY.GETDIR(ray))
        dz = VEC3.Z(RAY.GETDIR(ray))

        hit = RAY.HITBOX(ray, -1, -1, -1, 1, 1, 1)
        IF hit THEN
            d = RAY.HITBOX_DISTANCE(ray, -1, -1, -1, 1, 1, 1)
            PRINT "Hit cube at dist " + STR(d)
        END IF
        RAY.FREE(ray)
    END IF

    RENDER.FRAME()
WEND

CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

### Creation aliases

| Command | Description |
|--------|-------------|
| `RAY.MAKE(ox,oy,oz, dx,dy,dz)` | Deprecated alias of `RAY.CREATE`. |
| `RAY.SETPOSITION(ray, ox,oy,oz)` | Set ray origin in place. |

### Hit result Y/Z accessors

All intersection commands return a result struct. The per-axis accessors listed below complement the `X` forms already documented.

| Command | Description |
|--------|-------------|
| `RAY.HITBOX_POINTY(ray, min, max)` / `HITBOX_POINTZ` | Y/Z of hit point on AABB. |
| `RAY.HITBOX_NORMALY(ray, min, max)` / `HITBOX_NORMALZ` | Y/Z of hit normal on AABB. |
| `RAY.HITSPHERE_POINTY(ray, cx,cy,cz, r)` / `HITSPHERE_POINTZ` | Y/Z of hit point on sphere. |
| `RAY.HITSPHERE_NORMALY(...)` / `HITSPHERE_NORMALZ` | Y/Z of hit normal on sphere. |
| `RAY.HITPLANE_POINTY(ray, nx,ny,nz, d)` / `HITPLANE_POINTZ` | Y/Z of hit point on plane. |
| `RAY.HITPLANE_NORMALY(...)` / `HITPLANE_NORMALZ` | Y/Z of hit normal on plane. |
| `RAY.HITTRIANGLE_POINTY(...)` / `HITTRIANGLE_POINTZ` | Y/Z of hit point on triangle. |
| `RAY.HITTRIANGLE_NORMALY(...)` / `HITTRIANGLE_NORMALZ` | Y/Z of hit normal on triangle. |
| `RAY.HITMESH_POINTY(...)` / `HITMESH_POINTZ` | Y/Z of hit point on mesh. |
| `RAY.HITMESH_NORMALY(...)` / `HITMESH_NORMALZ` | Y/Z of hit normal on mesh. |
| `RAY.HITMODEL_POINTY(...)` / `HITMODEL_POINTZ` | Y/Z of hit point on model. |
| `RAY.HITMODEL_NORMALY(...)` / `HITMODEL_NORMALZ` | Y/Z of hit normal on model. |
| `RAY.INTERSECTSMODEL_HIT(ray, mdl)` | Returns `TRUE` if ray hits model. |
| `RAY.INTERSECTSMODEL_DISTANCE(ray, mdl)` | Returns hit distance. |
| `RAY.INTERSECTSMODEL_POINTX(ray, mdl)` / `RAY.INTERSECTSMODEL_POINTY(ray, mdl)` / `RAY.INTERSECTSMODEL_POINTZ(ray, mdl)` | Hit point per axis. |
| `RAY.INTERSECTSMODEL_NORMALX(ray, mdl)` / `RAY.INTERSECTSMODEL_NORMALY(ray, mdl)` / `RAY.INTERSECTSMODEL_NORMALZ(ray, mdl)` | Hit normal per axis. |

---

## See also

- [PICK.md](PICK.md) — higher-level `PICK.FROMCAMERA`, entity picking
- [CAMERA.md](CAMERA.md) — `CAMERA.GETRAY` for screen-to-world rays
- [COLLISION.md](COLLISION.md) — overlap tests and distance queries
- [RAY2D.md](RAY2D.md) — 2D ray/circle/rect intersection tests
