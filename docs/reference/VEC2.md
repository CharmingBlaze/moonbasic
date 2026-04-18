# Vec2 Commands

2D vector math with heap-handle and scalar-convenience overloads. All angles in **radians**.

## Core Workflow

1. Create with `VEC2.CREATE(x, y)` — returns a heap handle.
2. Perform arithmetic (`VEC2.ADD`, `VEC2.SUB`, `VEC2.MUL`, `VEC2.NORMALIZE`, etc.).
3. Read components with `VEC2.X(v)` / `VEC2.Y(v)`.
4. Free with `VEC2.FREE(v)` when no longer needed.

For scalar-only operations (no heap allocation) use the overloads that take raw floats instead of handles.

---

## Creation

### `VEC2.CREATE(x, y)` 

Creates a new 2D vector handle. Returns a **heap handle**.

---

### `VEC2.FREE(v)` 

Frees the vector handle.

---

## Components

### `VEC2.X(v)` 

Returns the X component as a float.

---

### `VEC2.Y(v)` 

Returns the Y component as a float.

---

### `VEC2.SET(v, x, y)` 

Updates the components of an existing vector handle in place.

---

## Arithmetic

### `VEC2.ADD(a, b)` 

Returns a new handle equal to `a + b`.

---

### `VEC2.SUB(a, b)` 

Returns a new handle equal to `a - b`.

---

### `VEC2.MUL(v, scalar)` 

Returns a new handle equal to `v * scalar`.

---

## Length & Distance

### `VEC2.LENGTH(v)` / `VEC2.LENGTH(x, y)` 

Returns the vector magnitude. Accepts a handle or raw floats.

---

### `VEC2.DIST(a, b)` / `VEC2.DIST(x1, y1, x2, y2)` 

Returns the distance between two points. Accepts handles or raw floats.

---

### `VEC2.DISTSQ(x1, y1, x2, y2)` 

Returns the squared distance (avoids sqrt — useful for comparisons).

---

### `VEC2.DISTANCE(a, b)` 

Alias of `VEC2.DIST` (handle overload).

---

## Normalise & Direction

### `VEC2.NORMALIZE(v)` / `VEC2.NORMALIZE(x, y)` 

Returns a new unit vector in the same direction. Accepts a handle or raw floats.

---

### `VEC2.ANGLE(a, b)` 

Returns the signed angle in **radians** between vectors `a` and `b`.

---

### `VEC2.ROTATE(v, radians)` 

Returns a new vector rotated by `radians`.

---

## Interpolation & Movement

### `VEC2.LERP(a, b, t)` 

Returns a new handle linearly interpolated between `a` and `b` by `t` (0–1).

---

### `VEC2.MOVE_TOWARD(x, y, tx, ty, step)` 

Returns a new handle moved from `(x, y)` toward `(tx, ty)` by at most `step` units.

---

### `VEC2.PUSHOUT(x, y, cx, cy, radius)` 

Returns a new handle pushed out of a circle at `(cx, cy)` with `radius`. Useful for circle overlap resolution.

---

## Transform

### `VEC2.TRANSFORMMAT4(v, mat4Handle)` 

Returns a new handle with the vector transformed by a 4×4 matrix.

---

## Full Example

Player movement clamped to a circular arena using Vec2.

```basic
WINDOW.OPEN(800, 600, "Vec2 Demo")
WINDOW.SETFPS(60)

px = 400.0
py = 300.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    dx = 0.0
    dy = 0.0
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN dx =  1.0
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN dx = -1.0
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN dy =  1.0
    IF INPUT.KEYDOWN(KEY_UP)    THEN dy = -1.0

    ; normalise diagonal movement
    IF dx <> 0 AND dy <> 0 THEN
        vel = VEC2.CREATE(dx, dy)
        nrm = VEC2.NORMALIZE(vel)
        dx  = VEC2.X(nrm) * 200 * dt
        dy  = VEC2.Y(nrm) * 200 * dt
        VEC2.FREE(nrm)
        VEC2.FREE(vel)
    ELSE
        dx = dx * 200 * dt
        dy = dy * 200 * dt
    END IF

    px = px + dx
    py = py + dy

    ; clamp inside arena circle
    pushed = VEC2.PUSHOUT(px, py, 400, 300, 180)
    px = VEC2.X(pushed)
    py = VEC2.Y(pushed)
    VEC2.FREE(pushed)

    RENDER.CLEAR(20, 20, 40)
    DRAW.CIRCLEOUTLINE(400, 300, 180, 60, 60, 100, 255)
    DRAW.CIRCLE(INT(px), INT(py), 12, 80, 160, 255, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `VEC2.MAKE(x, y)` | Deprecated alias of `VEC2.CREATE`. |

---

## See also

- [VEC3.md](VEC3.md) — 3D vector math
- [VEC_QUAT.md](VEC_QUAT.md) — quaternion and combined examples
- [MATH.md](MATH.md) — scalar trig, lerp, clamp
