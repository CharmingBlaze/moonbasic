# BBox Commands

Axis-aligned bounding box (AABB) handles for fast overlap and sphere tests. Used for broad-phase culling and simple collision without physics.

## Core Workflow

1. `BBOX.CREATE(minX, minY, minZ, maxX, maxY, maxZ)` — create a box handle.
2. `BBOX.CHECK(a, b)` — test two boxes for overlap.
3. `BBOX.CHECKSPHERE(box, cx, cy, cz, radius)` — test box vs sphere.
4. `BBOX.FREE(box)` when done.

---

## Creation

### `BBOX.CREATE(minX, minY, minZ, maxX, maxY, maxZ)` 

Creates an AABB from min and max corners. Returns a **bbox handle**.

---

### `BBOX.FROMMODEL(modelHandle)` 

Creates a bbox that wraps a loaded model's geometry. Returns a **bbox handle**. Useful for automatic bounds extraction.

---

## Modification

### `BBOX.SETMIN(box, x, y, z)` 

Sets the min corner. Returns the bbox handle.

---

### `BBOX.SETMAX(box, x, y, z)` 

Sets the max corner. Returns the bbox handle.

---

## Inspection

### `BBOX.GETMIN(box)` / `BBOX.MIN(box)` 

Returns the min corner as a `VEC3` handle.

---

### `BBOX.GETMAX(box)` / `BBOX.MAX(box)` 

Returns the max corner as a `VEC3` handle.

---

## Tests

### `BBOX.CHECK(boxA, boxB)` 

Returns `TRUE` if the two bounding boxes overlap.

---

### `BBOX.CHECKSPHERE(box, cx, cy, cz, radius)` 

Returns `TRUE` if the bounding box overlaps a sphere at `(cx, cy, cz)` with `radius`.

---

## Lifetime

### `BBOX.FREE(box)` 

Frees the bounding box handle.

---

## Full Example

Simple AABB overlap check between two moving entities.

```basic
WINDOW.OPEN(800, 450, "BBox Demo")
WINDOW.SETFPS(60)

px = 100.0

boxA = BBOX.CREATE(0, 0, 0, 40, 40, 40)    ; static
boxB = BBOX.CREATE(0, 0, 0, 40, 40, 40)    ; moves with px

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 100 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 100 * dt

    ; update boxB position each frame
    BBOX.SETMIN(boxB, px,      0,  0)
    BBOX.SETMAX(boxB, px + 40, 40, 40)

    hit = BBOX.CHECK(boxA, boxB)

    RENDER.CLEAR(15, 15, 30)
    col = if(hit, 255, 80)
    DRAW.RECT(0,   0, 40, 40, 80, 80, 200, 255)
    DRAW.RECT(INT(px), 0, 40, 40, col, 80, 80, 255)
    IF hit THEN DRAW.TEXT("OVERLAP", 300, 200, 24, 255, 80, 80, 255)
    RENDER.FRAME()
WEND

BBOX.FREE(boxA)
BBOX.FREE(boxB)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `BBOX.MAKE(minX,minY,minZ, maxX,maxY,maxZ)` | Deprecated alias of `BBOX.CREATE`. |

---

## See also

- [BSPHERE.md](BSPHERE.md) — bounding sphere handles
- [COLLISION.md](COLLISION.md) — entity-level collision tests
- [SHAPE.md](SHAPE.md) — Jolt physics shapes
