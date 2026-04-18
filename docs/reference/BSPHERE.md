# BSphere Commands

Bounding sphere handles for fast overlap tests against other spheres and boxes. Complement to [BBOX.md](BBOX.md) for broad-phase culling.

## Core Workflow

1. `BSPHERE.CREATE(x, y, z, radius)` — create a sphere handle.
2. `BSPHERE.CHECK(a, b)` — sphere vs sphere overlap.
3. `BSPHERE.CHECKBOX(sphere, bbox)` — sphere vs AABB overlap.
4. `BSPHERE.FREE(sphere)` when done.

---

## Creation

### `BSPHERE.CREATE(x, y, z, radius)` 

Creates a bounding sphere centred at `(x, y, z)` with `radius`. Returns a **bsphere handle**.

---

## Position & Radius

### `BSPHERE.SETPOS(sphere, x, y, z)` 

Sets the sphere centre. Returns the sphere handle.

---

### `BSPHERE.GETPOS(sphere)` / `BSPHERE.POS(sphere)` 

Returns the centre as a `VEC3` handle.

---

### `BSPHERE.SETRADIUS(sphere, radius)` 

Sets the sphere radius. Returns the sphere handle.

---

### `BSPHERE.GETRADIUS(sphere)` / `BSPHERE.RADIUS(sphere)` 

Returns the radius as a float.

---

## Tests

### `BSPHERE.CHECK(sphereA, sphereB)` 

Returns `TRUE` if the two bounding spheres overlap.

---

### `BSPHERE.CHECKBOX(sphere, bboxHandle)` 

Returns `TRUE` if the sphere overlaps the axis-aligned bounding box.

---

## Lifetime

### `BSPHERE.FREE(sphere)` 

Frees the bounding sphere handle.

---

## Full Example

Projectile collision using bounding spheres.

```basic
WINDOW.OPEN(800, 450, "BSphere Demo")
WINDOW.SETFPS(60)

; enemy sphere at centre
enemy = BSPHERE.CREATE(400, 225, 0, 40)

; bullet starts at left
bx = 50.0
bullet = BSPHERE.CREATE(bx, 225, 0, 8)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    bx = bx + 300 * dt
    BSPHERE.SETPOS(bullet, bx, 225, 0)

    hit = BSPHERE.CHECK(bullet, enemy)

    RENDER.CLEAR(15, 15, 30)
    DRAW.CIRCLE(400, 225, 40, 80, 80, 200, 255)
    DRAW.CIRCLE(INT(bx), 225, 8, 255, 200, 60, 255)
    IF hit THEN
        DRAW.TEXT("HIT!", 360, 100, 28, 255, 80, 80, 255)
        bx = 50.0
    END IF
    RENDER.FRAME()
WEND

BSPHERE.FREE(bullet)
BSPHERE.FREE(enemy)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `BSPHERE.MAKE(cx,cy,cz, r)` | Deprecated alias of `BSPHERE.CREATE`. |
| `BSPHERE.SETPOSITION(bs, x,y,z)` | Alias of `BSPHERE.SETCENTER`. |

---

## See also

- [BBOX.md](BBOX.md) — axis-aligned bounding box handles
- [COLLISION.md](COLLISION.md) — entity-level collision tests
- [SHAPE.md](SHAPE.md) — Jolt physics sphere shapes
