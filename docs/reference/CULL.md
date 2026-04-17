# Cull Commands

Visibility culling, distance queries, and occlusion helpers for large scenes.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Set the maximum draw distance with `CULL.SETMAXDISTANCE`.
2. Optionally enable occlusion culling with `CULL.OCCLUSIONENABLE` and register occluders with `CULL.OCCLUDERADD`.
3. Each frame, test objects with `CULL.SPHEREVISIBLE`, `CULL.AABBVISIBLE`, or `CULL.POINTVISIBLE` before drawing.
4. Use `CULL.DISTANCE` / `CULL.INRANGE` for LOD or AI decisions.
5. Monitor culling efficiency with `CULL.STATS*` counters.

---

### `CULL.SPHEREVISIBLE(x, y, z, radius)` 

Returns `TRUE` if a sphere at the given position is within the view frustum and draw distance.

---

### `CULL.AABBVISIBLE(minX, minY, minZ, maxX, maxY, maxZ)` 

Returns `TRUE` if the axis-aligned bounding box intersects the view frustum.

---

### `CULL.POINTVISIBLE(x, y, z)` 

Returns `TRUE` if the point is inside the view frustum.

---

### `CULL.DISTANCE(x1, y1, z1, x2, y2, z2)` 

Returns the Euclidean distance between two 3D points.

---

### `CULL.DISTANCESQ(x1, y1, z1, x2, y2, z2)` 

Returns the squared distance between two 3D points (faster than `CULL.DISTANCE` — no square root).

---

### `CULL.INRANGE(x1, y1, z1, x2, y2, z2, maxDist)` 

Returns `TRUE` if the distance between the two points is less than `maxDist`.

---

### `CULL.SETMAXDISTANCE(dist)` 

Sets the maximum draw distance. Objects beyond this are culled regardless of frustum.

---

### `CULL.GETMAXDISTANCE()` 

Returns the current maximum draw distance.

---

### `CULL.SETBACKFACECULLING(enabled)` 

Enables or disables back-face culling globally.

---

### `CULL.BEHINDHORIZON(x, y, z)` 

Returns `TRUE` if the point is below the camera's horizon (useful for terrain).

---

### `CULL.ISOCCLUDED(x, y, z)` 

Returns `TRUE` if the point is hidden behind a registered occluder.

---

### `CULL.OCCLUDERADD(minX, minY, minZ, maxX, maxY, maxZ)` 

Adds an axis-aligned box as an occluder for the current frame.

---

### `CULL.OCCLUDERCLEAR()` 

Removes all registered occluders. Call at the start of each frame before re-adding.

---

### `CULL.OCCLUSIONENABLE(enabled)` 

Enables or disables the occlusion culling system.

---

### `CULL.BATCHSPHERE(entityHandle, radius)` 

Adds an entity to the batched sphere-cull list. Returns a cull handle for tracking.

---

### `CULL.STATSTOTAL()` 

Returns the total number of objects tested this frame.

---

### `CULL.STATSVISIBLE()` 

Returns the number of objects that passed all cull tests.

---

### `CULL.STATSCULLED()` 

Returns the total number of objects culled (frustum + distance + horizon + occlusion).

---

### `CULL.STATSFRUSTUMCULLED()` 

Returns the number of objects culled by frustum.

---

### `CULL.STATSDISTANCECULLED()` 

Returns the number of objects culled by distance.

---

### `CULL.STATSHORIZONCULLED()` 

Returns the number of objects culled by horizon.

---

### `CULL.STATSOCCLUSIONCULLED()` 

Returns the number of objects culled by occlusion.

---

### `CULL.STATSRESET()` 

Resets all cull stat counters to zero.

---

## Full Example

This example sets up distance and occlusion culling and only draws visible entities.

```basic
; Configure culling
CULL.SETMAXDISTANCE(500.0)
CULL.OCCLUSIONENABLE(TRUE)

WHILE NOT WINDOW.SHOULDCLOSE()
    CULL.OCCLUDERCLEAR()

    ; Register a building as an occluder
    CULL.OCCLUDERADD(10, 0, 10, 15, 20, 15)

    ; Only draw entities that pass visibility
    FOR i = 0 TO entityCount - 1
        ex = ENTITY.GETX(entities(i))
        ey = ENTITY.GETY(entities(i))
        ez = ENTITY.GETZ(entities(i))
        IF CULL.SPHEREVISIBLE(ex, ey, ez, 2.0)
            ENTITY.DRAW(entities(i))
        END IF
    NEXT

    ; Print stats
    PRINT "Visible: " + STR(CULL.STATSVISIBLE()) + " / " + STR(CULL.STATSTOTAL())
    CULL.STATSRESET()

    RENDER.BEGINFRAME()
    RENDER.ENDFRAME()
WEND
```
