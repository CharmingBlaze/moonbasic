# Collision Commands

Stateless geometry tests for 2D and 3D — boxes, circles, spheres, AABBs, lines, and distances. Persistent **`BBOX`** / **`BSPHERE`** handles are supported for repeated queries. **Preferred for new code:** the **`COLLISION.*`** helpers below take **`VEC2`** / **`VEC3`** handles so you pass **two to four arguments** instead of long scalar lists (see [STYLE_GUIDE.md](../../STYLE_GUIDE.md) on short parameter lists).

No physics world is required. Physics-driven collision is documented in [PHYSICS2D.md](PHYSICS2D.md) and [PHYSICS3D.md](PHYSICS3D.md). Ray queries: [RAYCAST.md](RAYCAST.md).

## Core Workflow

1. Build positions and sizes with **`VEC2.CREATE`** / **`VEC3.CREATE`** (or reuse handles you already have for movement).
2. Call **`COLLISION.BOXOVERLAP2D`**, **`COLLISION.SPHEREOVERLAP3D`**, and related **`COLLISION.*`** commands — they delegate to the same pure math as the legacy scalar builtins.
3. For one-off tests without vectors, the legacy **`BOXCOLLIDE`**, **`SPHERECOLLIDE`**, … functions remain available.
4. For bounds you update every frame, consider **`BBOX.*`** / **`BSPHERE.*`** or keep **`VEC2`** / **`VEC3`** fields and use **`COLLISION.*`**.

### Method chaining

Where **`BBOX`** / **`BSPHERE`** expose handle methods (for example **`bbox.check(other)`**), prefer those over repeating **`BBOX.CHECK`** when it reads more clearly in your loop.

---

## COLLISION namespace (VEC helpers)

These commands expect **heap handles** from **`VEC2.*`** / **`VEC3.*`**. They require a bound runtime heap (normal full-runtime games).

### `COLLISION.BOXOVERLAP2D(posA, sizeA, posB, sizeB)` 

Returns **`TRUE`** if two axis-aligned 2D rectangles overlap. Each rectangle is **position (min corner)** and **width/height** as separate **`VEC2`** values — **four arguments** instead of eight floats.

- **Arguments**:
  - `posA`, `sizeA`, `posB`, `sizeB` (**handle**): **`VEC2`** handles.
- **Returns**: (Boolean)

---

### `COLLISION.CIRCLEOVERLAP2D(center1, radius1, center2, radius2)` 

Returns **`TRUE`** if two circles overlap. **Four arguments:** two centers as **`VEC2`** and two radii as floats.

- **Arguments**:
  - `center1`, `center2` (**handle**): **`VEC2`** centers.
  - `radius1`, `radius2` (**float**): Radii.
- **Returns**: (Boolean)

---

### `COLLISION.POINTINBOX2D(point, boxPos, boxSize)` 

Returns **`TRUE`** if **`point`** lies inside the axis-aligned box with minimum corner **`boxPos`** and size **`boxSize`**.

- **Arguments**:
  - `point`, `boxPos`, `boxSize` (**handle**): **`VEC2`** handles.
- **Returns**: (Boolean)

---

### `COLLISION.CIRCLEBOX2D(circleCenter, radius, boxPos, boxSize)` 

Returns **`TRUE`** if a circle overlaps an axis-aligned rectangle.

- **Arguments**:
  - `circleCenter`, `boxPos`, `boxSize` (**handle**): **`VEC2`**.
  - `radius` (**float**): Circle radius.
- **Returns**: (Boolean)

---

### `COLLISION.LINESEGINTERSECT2D(a1, a2, b1, b2)` 

Returns **`TRUE`** if the two 2D segments intersect. Each endpoint is a **`VEC2`**.

- **Arguments**:
  - `a1`, `a2`, `b1`, `b2` (**handle**): Segment endpoints as **`VEC2`**.
- **Returns**: (Boolean)

---

### `COLLISION.POINTONSEG2D(point, segA, segB, threshold)` 

Returns **`TRUE`** if **`point`** is within **`threshold`** distance of the segment **`segA`–`segB`** (same idea as **`POINTONLINE`**).

- **Arguments**:
  - `point`, `segA`, `segB` (**handle**): **`VEC2`**.
  - `threshold` (**float**): Maximum distance from the segment.
- **Returns**: (Boolean)

---

### `COLLISION.SPHEREOVERLAP3D(center1, radius1, center2, radius2)` 

Returns **`TRUE`** if two spheres overlap. **Four arguments** instead of eight floats.

- **Arguments**:
  - `center1`, `center2` (**handle**): **`VEC3`** centers.
  - `radius1`, `radius2` (**float**): Radii.
- **Returns**: (Boolean)

---

### `COLLISION.AABBOVERLAP3D(minA, maxA, minB, maxB)` 

Returns **`TRUE`** if two axis-aligned 3D boxes overlap. Each box is described by **minimum and maximum corners** (same convention as **`AABBCOLLIDE`**): **four `VEC3` handles**.

- **Arguments**:
  - `minA`, `maxA`, `minB`, `maxB` (**handle**): **`VEC3`** corners.
- **Returns**: (Boolean)

---

### `COLLISION.SPHEREBOX3D(sphereCenter, radius, boxMin, boxSize)` 

Returns **`TRUE`** if a sphere intersects an axis-aligned box. The box is **minimum corner** plus **size** along each axis (**five** values via **four** handles + one float).

- **Arguments**:
  - `sphereCenter` (**handle**): **`VEC3`** sphere center.
  - `radius` (**float**): Sphere radius.
  - `boxMin` (**handle**): **`VEC3`** minimum corner of the box.
  - `boxSize` (**handle**): **`VEC3`** width, height, and depth.
- **Returns**: (Boolean)

---

### `COLLISION.POINTINAABB3D(point, boxMin, boxSize)` 

Returns **`TRUE`** if a 3D point lies inside the axis-aligned box given by **minimum corner** and **size**.

- **Arguments**:
  - `point`, `boxMin`, `boxSize` (**handle**): **`VEC3`**.
- **Returns**: (Boolean)

---

## Scalar 2D tests

### `BOXCOLLIDE(x1, y1, w1, h1, x2, y2, w2, h2)` 

Returns **`TRUE`** if two axis-aligned 2D rectangles overlap. Prefer **`COLLISION.BOXOVERLAP2D`** when you already use **`VEC2`** handles.

- **Arguments**:
  - `x1, y1, w1, h1`: (Float) First rectangle position and size.
  - `x2, y2, w2, h2`: (Float) Second rectangle position and size.
- **Returns**: (Boolean)

---

### `CIRCLECOLLIDE(x1, y1, r1, x2, y2, r2)` 

Returns **`TRUE`** if two circles overlap.

- **Arguments**:
  - `x1, y1, r1`: (Float) First circle center and radius.
  - `x2, y2, r2`: (Float) Second circle center and radius.
- **Returns**: (Boolean)

---

### `CIRCLEBOXCOLLIDE(cx, cy, cr, bx, by, bw, bh)` 

Returns **`TRUE`** if a circle and an axis-aligned rectangle overlap.

- **Returns**: (Boolean)

---

### `POINTINBOX(px, py, bx, by, bw, bh)` 

Returns **`TRUE`** if point **`(px, py)`** is inside the rectangle.

- **Returns**: (Boolean)

---

### `POINTINCIRCLE(px, py, cx, cy, cr)` 

Returns **`TRUE`** if point **`(px, py)`** is inside the circle.

- **Returns**: (Boolean)

---

### `LINECOLLIDE(x1, y1, x2, y2, x3, y3, x4, y4)` 

Returns **`TRUE`** if two line segments intersect.

- **Arguments**:
  - `x1, y1` to `x2, y2`: (Float) Line 1.
  - `x3, y3` to `x4, y4`: (Float) Line 2.
- **Returns**: (Boolean)

---

### `POINTONLINE(px, py, x1, y1, x2, y2, threshold)` 

Returns **`TRUE`** if point **`(px, py)`** lies near the segment within **`threshold`**.

- **Returns**: (Boolean)

---

## Scalar 3D tests

### `SPHERECOLLIDE(x1, y1, z1, r1, x2, y2, z2, r2)` 

Returns **`TRUE`** if two 3D spheres overlap.

- **Arguments**:
  - `x, y, z, r`: (Float) Center and radius for each sphere.
- **Returns**: (Boolean)

---

### `AABBCOLLIDE(ax, ay, az, aw, ah, ad, bx, by, bz, bw, bh, bd)` 

Returns **`TRUE`** if two axis-aligned 3D boxes overlap. Arguments are the **minimum corner** **`(x,y,z)`** and **maximum corner** **`(x,y,z)`** for box A, then the same for box B (twelve floats total: **not** position + size). Prefer **`COLLISION.AABBOVERLAP3D`** with **`VEC3`** min/max pairs when possible.

- **Returns**: (Boolean)

---

### `SPHEREBOXCOLLIDE(sx, sy, sz, sr, bx, by, bz, bw, bh, bd)` 

Returns **`TRUE`** if a sphere overlaps a 3D axis-aligned box. The box uses **minimum corner** **`(bx,by,bz)`** and **size** **`(bw,bh,bd)`**.

- **Returns**: (Boolean)

---

### `POINTINAABB(px, py, pz, bx, by, bz, bw, bh, bd)` 

Returns **`TRUE`** if a 3D point is inside the axis-aligned box (**min corner + size**).

- **Returns**: (Boolean)

---

### `BOXTOPLAND(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd)` 

Returns the landing-surface Y if a sphere (center **`px,py,pz`**, radius **`pr`**, vertical velocity **`pvy`**) lands on the **top** face of an axis-aligned box (**min corner** **`bx,by,bz`**, **size** **`bw,bh,bd`**), or **`0.0`** if there is no valid landing.

- **Returns**: (Float)

---

## Distance

### `DISTANCE2D(x1, y1, x2, y2)` / `DISTANCE3D(x1, y1, z1, x2, y2, z2)` 

Returns the Euclidean distance between two points.

- **Returns**: (Float)

---

### `DISTANCESQ2D(x1, y1, x2, y2)` / `DISTANCESQ3D(x1, y1, z1, x2, y2, z2)` 

Returns the squared distance (no square root).

- **Returns**: (Float)

---

## Frustum

### `CHECK.INVIEW(entityIndex)` 

Returns **`TRUE`** if the indexed entity is inside the current camera frustum.

- **Arguments**:
  - `entityIndex`: (Integer) Entity id.
- **Returns**: (Boolean)

---

## BBox handle (`BBOX.*`)

A persistent 3D axis-aligned bounding box handle. Create once, update corners, test each frame.

### `BBOX.CREATE(minX, minY, minZ, maxX, maxY, maxZ)` 

Creates a **`BBOX`** handle with the given min/max corners.

- **Arguments**:
  - `minX, minY, minZ`: (Float) Minimum corner.
  - `maxX, maxY, maxZ`: (Float) Maximum corner.
- **Returns**: (Handle)

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

Returns the corner as a **`[x, y, z]`** array handle.

- **Returns**: (Handle) Array handle.
- *Handle shortcut*: `bbox.getMin()`, `bbox.getMax()`

---

### `BBOX.CHECK(bbox, other)` 

Returns **`TRUE`** if this AABB overlaps another **`BBOX`** handle.

- **Arguments**:
  - `other`: (Handle) The other **`BBOX`**.
- **Returns**: (Boolean)
- *Handle shortcut*: `bbox.check(other)`

---

### `BBOX.CHECKSPHERE(bbox, sx, sy, sz, r)` 

Returns **`TRUE`** if this AABB overlaps a sphere.

- **Arguments**:
  - `sx, sy, sz`: (Float) Sphere center.
  - `r`: (Float) Sphere radius.
- **Returns**: (Boolean)
- *Handle shortcut*: `bbox.checkSphere(sx, sy, sz, r)`

---

### `BBOX.FREE(bbox)` 

Frees the **`BBOX`** handle.

- *Handle shortcut*: `bbox.free()`

---

## BSphere handle (`BSPHERE.*`)

A persistent 3D bounding sphere handle.

### `BSPHERE.CREATE(x, y, z, radius)` 

Creates a **`BSPHERE`** centered at **`(x, y, z)`** with **`radius`**.

- **Arguments**:
  - `x, y, z`: (Float) Center.
  - `radius`: (Float) Radius.
- **Returns**: (Handle)

---

### `BSPHERE.SETPOS(bsphere, x, y, z)` 

Moves the sphere centre.

- **Returns**: (Handle) The bsphere handle (for chaining).
- *Handle shortcut*: `bsphere.setPos(x, y, z)`

---

### `BSPHERE.GETPOS(bsphere)` 

Returns the sphere centre as a **`[x, y, z]`** array handle.

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

Returns **`TRUE`** if this sphere overlaps another **`BSPHERE`** handle.

- **Arguments**:
  - `other`: (Handle) The other **`BSPHERE`**.
- **Returns**: (Boolean)
- *Handle shortcut*: `bsphere.check(other)`

---

### `BSPHERE.CHECKBOX(bsphere, bbox)` 

Returns **`TRUE`** if this sphere overlaps a **`BBOX`** handle.

- **Arguments**:
  - `bbox`: (Handle) The **`BBOX`**.
- **Returns**: (Boolean)
- *Handle shortcut*: `bsphere.checkBox(bbox)`

---

### `BSPHERE.FREE(bsphere)` 

Frees the **`BSPHERE`** handle.

- *Handle shortcut*: `bsphere.free()`

---

## Full Example

This sample uses **`COLLISION.BOXOVERLAP2D`** with **`VEC2`** handles, then shows a scalar **`CIRCLEBOXCOLLIDE`** call for comparison.

```basic
; VEC2-based overlap (four handles)
pa = VEC2.CREATE(0.0, 0.0)
sa = VEC2.CREATE(32.0, 32.0)
pb = VEC2.CREATE(16.0, 16.0)
sb = VEC2.CREATE(32.0, 32.0)

IF COLLISION.BOXOVERLAP2D(pa, sa, pb, sb)
    PRINT "VEC2 boxes overlap"
ENDIF

VEC2.FREE(pa)
VEC2.FREE(sa)
VEC2.FREE(pb)
VEC2.FREE(sb)

; Same idea with scalars (eight floats)
IF CIRCLEBOXCOLLIDE(100.0, 150.0, 16.0, 120.0, 140.0, 32.0, 32.0)
    PRINT "Circle hits box"
ENDIF
PRINT "Distance: " + STR(DISTANCE2D(100.0, 150.0, 120.0, 140.0))

; Handle-based 3D bounds (optional) — px/pz stand in for entity position
px = 0.0
pz = 0.0
playerBox = BBOX.CREATE(-0.5, 0, -0.5, 0.5, 2, 0.5)
enemySphere = BSPHERE.CREATE(3, 1, 0, 1.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    playerBox.setMin(px - 0.5, 0,   pz - 0.5)
    playerBox.setMax(px + 0.5, 2.0, pz + 0.5)

    IF enemySphere.checkBox(playerBox)
        PRINT "Enemy hit player!"
    ENDIF
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
- [STYLE_GUIDE.md](../../STYLE_GUIDE.md) — **`Namespace.Method`**, chaining, and short parameter lists
