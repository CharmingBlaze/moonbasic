# Collision Commands

Stateless geometry collision tests for 2D and 3D — boxes, circles, spheres, AABBs, lines, and distances. Plus handle-based **`BBOX`** / **`BSPHERE`** objects for persistent bounds queries.

No physics world required. For physics-driven collision see [PHYSICS2D.md](PHYSICS2D.md) / [PHYSICS3D.md](PHYSICS3D.md). For raycasts see [RAYCAST.md](RAYCAST.md).

## Core Workflow

Call any function with world-space coordinates. Each returns `TRUE` / `FALSE` (or a `float`). For reusable bounds, create a **`BBOX`** or **`BSPHERE`** handle and call `.check()` each frame.

---

## 2. 2D Tests

### `BOXCOLLIDE(x1, y1, w1, h1, x2, y2, w2, h2)`
Returns `TRUE` if two axis-aligned 2D rectangles overlap.

- **Arguments**:
    - `x1, y1, w1, h1`: (Float) Rect 1 position and size.
    - `x2, y2, w2, h2`: (Float) Rect 2 position and size.
- **Returns**: (Boolean)

---

### `CIRCLECOLLIDE(x1, y1, r1, x2, y2, r2)`
Returns `TRUE` if two circles overlap.

- **Arguments**:
    - `x1, y1, r1`: (Float) Circle 1 center and radius.
    - `x2, y2, r2`: (Float) Circle 2 center and radius.
- **Returns**: (Boolean)

---

### `CIRCLEBOXCOLLIDE(cx, cy, cr, bx, by, bw, bh)`
Returns `TRUE` if a circle and an axis-aligned rectangle overlap.

- **Arguments**:
    - `cx, cy, cr`: (Float) Circle center and radius.
    - `bx, by, bw, bh`: (Float) Box position and size.
- **Returns**: (Boolean)

---

### `POINTINBOX(px, py, bx, by, bw, bh)`
Returns `TRUE` if point `(px, py)` is inside the rectangle.

- **Returns**: (Boolean)

---

### `POINTINCIRCLE(px, py, cx, cy, cr)`
Returns `TRUE` if point `(px, py)` is inside the circle.

- **Returns**: (Boolean)

---

### `LINECOLLIDE(x1, y1, x2, y2, x3, y3, x4, y4)`
Returns `TRUE` if two line segments intersect.

- **Arguments**:
    - `x1, y1` to `x2, y2`: (Float) Line 1.
    - `x3, y3` to `x4, y4`: (Float) Line 2.
- **Returns**: (Boolean)

---

### `POINTONLINE(px, py, x1, y1, x2, y2)`
Returns `TRUE` if point `(px, py)` lies on the line segment.

- **Returns**: (Boolean)

---

## 3D Tests

### `SPHERECOLLIDE(x1, y1, z1, r1, x2, y2, z2, r2)`
Returns `TRUE` if two 3D spheres overlap.

- **Arguments**:
    - `x, y, z, r`: (Float) Center and radius for sphere 1 and 2.
- **Returns**: (Boolean)

---

### `AABBCOLLIDE(ax, ay, az, aw, ah, ad, bx, by, bz, bw, bh, bd)`
Returns `TRUE` if two 3D axis-aligned bounding boxes overlap.

- **Arguments**:
    - `x, y, z, w, h, d`: (Float) Pos and size for AABB A and B.
- **Returns**: (Boolean)

---

### `SPHEREBOXCOLLIDE(sx, sy, sz, sr, bx, by, bz, bw, bh, bd)`
Returns `TRUE` if a sphere and a 3D AABB overlap.

- **Returns**: (Boolean)

---

### `POINTINAABB(px, py, pz, bx, by, bz, bw, bh, bd)`
Returns `TRUE` if a 3D point is inside the AABB.

- **Returns**: (Boolean)

---

### `BOXTOPLAND(sx, sy, sz, sr, bx, by, bz, bw, bh, bd)`
Returns the landing-surface Y if a sphere lands on top of an AABB, or `0.0` if no landing.

- **Returns**: (Float) The Y coordinate of the surface.

---

## Distance

### `DISTANCE2D(x1, y1, x2, y2)` / `DISTANCE3D(x1, y1, z1, x2, y2, z2)`
Returns the Euclidean distance between two points.

- **Returns**: (Float)

---

### `DISTANCESQ2D(x1, y1, x2, y2)` / `DISTANCESQ3D(x1, y1, z1, x2, y2, z2)`
Returns the squared distance. No square root — faster for comparisons.

- **Returns**: (Float)

---

## Frustum

### `CHECK.INVIEW(entityIndex)`
Returns `TRUE` if the indexed entity is inside the current camera frustum.

- **Arguments**:
    - `entityIndex`: (Integer) The ID of the entity to check.
- **Returns**: (Boolean)

---

## BBox Handle (`BBOX.*`)

A persistent 3D axis-aligned bounding box handle. Create once, update bounds, test each frame.

### `BBOX.CREATE(minX, minY, minZ, maxX, maxY, maxZ)`
Creates a BBox handle with the given min/max corners.

- **Arguments**:
    - `minX, minY, minZ`: (Float) Minimum corner.
    - `maxX, maxY, maxZ`: (Float) Maximum corner.
- **Returns**: (Handle) The BBox handle.

---

### `BBOX.SETMIN(bbox, x, y, z)` / `BBOX.SETMAX(bbox, x, y, z)`
Sets the corners of the bounding box.

- **Arguments**:
    - `bbox`: (Handle) The box to modify.
    - `x, y, z`: (Float) New corner coordinates.
- **Returns**: (Handle) The bbox handle (for chaining).
- *Handle shortcut*: `bbox.setMin(x, y, z)`, `bbox.setMax(x, y, z)`

---

### `BBOX.GETMIN(bbox)` / `BBOX.GETMAX(bbox)`
Returns the corner as a `[x, y, z]` array handle.

- **Returns**: (Handle) Array handle.
- *Handle shortcut*: `bbox.getMin()`, `bbox.getMax()`

---

### `BBOX.CHECK(bbox, other)`
Returns `TRUE` if this AABB overlaps another AABB handle.

- **Arguments**:
    - `other`: (Handle) The other BBox to test against.
- **Returns**: (Boolean)
- *Handle shortcut*: `bbox.check(other)`

---

### `BBOX.CHECKSPHERE(bbox, sx, sy, sz, r)`
Returns `TRUE` if this AABB overlaps a sphere.

- **Arguments**:
    - `sx, sy, sz`: (Float) Sphere center.
    - `r`: (Float) Sphere radius.
- **Returns**: (Boolean)
- *Handle shortcut*: `bbox.checkSphere(sx, sy, sz, r)`

---

### `BBOX.FREE(bbox)`
Freese the BBox handle.

- *Handle shortcut*: `bbox.free()`

---

## BSphere Handle (`BSPHERE.*`)

A persistent 3D bounding sphere handle.

### `BSPHERE.CREATE(x, y, z, radius)`
Creates a BSphere handle centred at `(x, y, z)` with `radius`.

- **Arguments**:
    - `x, y, z`: (Float) Center.
    - `radius`: (Float) Radius.
- **Returns**: (Handle) The BSphere handle.

---

### `BSPHERE.SETPOS(bsphere, x, y, z)`
Moves the sphere centre.

- **Returns**: (Handle) The bsphere handle (for chaining).
- *Handle shortcut*: `bsphere.setPos(x, y, z)`

---

### `BSPHERE.GETPOS(bsphere)`
Returns the sphere centre as a `[x, y, z]` array handle.

- *Handle shortcut*: `bsphere.getPos()`

---

### `BSPHERE.SETRADIUS(bsphere, r)`
Sets the sphere radius.

- **Returns**: (Handle) The bsphere handle (for chaining).
- *Handle shortcut*: `bsphere.setRadius(r)`

---

### `BSPHERE.GETRADIUS(bsphere)`
Returns the sphere radius as a float.

- *Handle shortcut*: `bsphere.getRadius()`

---

### `BSPHERE.CHECK(bsphere, other)`
Returns `TRUE` if this sphere overlaps another BSphere handle.

- **Arguments**:
    - `other`: (Handle) The other BSphere.
- **Returns**: (Boolean)
- *Handle shortcut*: `bsphere.check(other)`

---

### `BSPHERE.CHECKBOX(bsphere, bbox)`
Returns `TRUE` if this sphere overlaps a BBox handle.

- **Arguments**:
    - `bbox`: (Handle) The BBox.
- **Returns**: (Boolean)
- *Handle shortcut*: `bsphere.checkBox(bbox)`

---

### `BSPHERE.FREE(bsphere)`
Frees the BSphere handle.

- *Handle shortcut*: `bsphere.free()`

---

## Full Example

```basic
; Stateless 2D test
px = 100.0 : py = 150.0 : pr = 16.0
ex = 120.0 : ey = 140.0 : ew = 32.0 : eh = 32.0

IF CIRCLEBOXCOLLIDE(px, py, pr, ex, ey, ew, eh)
    PRINT "Circle hit box!"
END IF
PRINT "Distance: " + STR(DISTANCE2D(px, py, ex, ey))

; Handle-based 3D bounds
playerBox = BBOX.CREATE(-0.5, 0, -0.5, 0.5, 2, 0.5)
enemySphere = BSPHERE.CREATE(3, 1, 0, 1.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    ; update bounds to match entity positions each frame
    playerBox.setMin(px - 0.5, 0,   pz - 0.5)
    playerBox.setMax(px + 0.5, 2.0, pz + 0.5)

    IF enemySphere.checkBox(playerBox)
        PRINT "Enemy hit player!"
    END IF
    RENDER.FRAME()
WEND

BBOX.FREE(playerBox)
BSPHERE.FREE(enemySphere)
```

---

## See also

- [PHYSICS2D.md](PHYSICS2D.md) — Box2D physics collision
- [PHYSICS3D.md](PHYSICS3D.md) — Jolt 3D physics collision
- [RAYCAST.md](RAYCAST.md) — ray vs world queries
- [SPRITE.md](SPRITE.md) — sprite bounding-box collision
