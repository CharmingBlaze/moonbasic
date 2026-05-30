# Vec3 Commands

3D vector math with heap-handle and scalar-convenience overloads. All angles in **radians**.

## Core Workflow

1. Create with `VEC3.CREATE(x, y, z)` — returns a heap handle.
2. Perform arithmetic (`VEC3.ADD`, `VEC3.DOT`, `VEC3.CROSS`, `VEC3.NORMALIZE`, etc.).
3. Read components with `VEC3.X(v)` / `VEC3.Y(v)` / `VEC3.Z(v)`.
4. Free with `VEC3.FREE(v)`.

Scalar overloads (passing raw floats instead of handles) avoid heap allocation for simple one-off calculations.

---

## Creation

### `VEC3.CREATE(x, y, z)` 

Creates a new 3D vector heap handle. Alias: `VEC3.VEC3(x, y, z)`.

---

### `VEC3.FREE(v)` 

Frees the vector handle.

---

## Components

### `VEC3.X(v)` / `VEC3.Y(v)` / `VEC3.Z(v)` 

Returns the X, Y, or Z component as a float scalar.

---

### `VEC3.SET(v, x, y, z)` 

Updates the components of an existing vector handle in place.

---

## Arithmetic

### `VEC3.ADD(a, b)` 

Returns a new handle equal to `a + b`. Alias: `VEC3.VECADD`.

---

### `VEC3.SUB(a, b)` 

Returns a new handle equal to `a - b`. Alias: `VEC3.VECSUB`.

---

### `VEC3.MUL(v, scalar)` 

Returns a new handle scaled by `scalar`. Alias: `VEC3.VECSCALE`.

---

### `VEC3.DIV(v, scalar)` 

Returns a new handle divided by `scalar`.

---

### `VEC3.NEGATE(v)` 

Returns a new handle with all components negated (`-v`).

---

## Length & Distance

### `VEC3.LENGTH(v)` / `VEC3.LENGTH(x, y, z)` 

Returns the vector magnitude. Accepts a handle or three raw floats. Alias: `VEC3.VECLENGTH`.

---

### `VEC3.DIST(a, b)` / `VEC3.DIST(x1, y1, z1, x2, y2, z2)` 

Returns the distance between two points.

---

### `VEC3.DISTSQ(x1, y1, z1, x2, y2, z2)` 

Returns the squared distance (no sqrt — faster for comparisons).

---

### `VEC3.DISTANCE(a, b)` 

Alias of `VEC3.DIST` (handle overload).

---

## Normalise & Direction

### `VEC3.NORMALIZE(v)` / `VEC3.NORMALIZE(x, y, z)` 

Returns a new unit vector. Accepts handle or raw floats. Alias: `VEC3.VECNORMALIZE`.

---

### `VEC3.DOT(a, b)` 

Returns the dot product as a float. Alias: `VEC3.VECDOT`.

---

### `VEC3.CROSS(a, b)` 

Returns a new handle perpendicular to both `a` and `b`. Alias: `VEC3.VECCROSS`.

---

### `VEC3.REFLECT(v, normal)` 

Returns the reflection of `v` about `normal` (both handles). Useful for bounce directions.

---

### `VEC3.ANGLE(a, b)` 

Returns the angle in **radians** between `a` and `b`.

---

### `VEC3.PROJECT(v, onto)` 

Returns the projection of `v` onto `onto` as a new handle.

---

### `VEC3.ORTHONORMALIZE(a, b)` 

Orthonormalises vectors `a` and `b` in place (modifies handles, Gram-Schmidt).

---

## Interpolation

### `VEC3.LERP(a, b, t)` 

Returns a new handle linearly interpolated between `a` and `b` by `t` (0–1).

---

## Quaternion Rotation

### `VEC3.ROTATEBYQUAT(v, q)` 

Returns a new handle with `v` rotated by quaternion `q`.

---

## Transform

### `VEC3.TRANSFORMMAT4(v, mat4Handle)` 

Returns a new handle with the vector transformed by a 4×4 matrix.

---

### `VEC3.EQUALS(a, b)` 

Returns `TRUE` if `a` and `b` have identical components.

---

## Full Example

Calculating a look direction from two world positions and reflecting off a surface.

```basic
WINDOW.OPEN(960, 540, "Vec3 Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 3, -8)
CAMERA.SETTARGET(cam, 0, 0, 0)

origin = VEC3.CREATE(0, 5, 0)
target = VEC3.CREATE(3, 0, 0)
up     = VEC3.CREATE(0, 1, 0)

; direction from origin to target
dir = VEC3.SUB(target, origin)
dir = VEC3.NORMALIZE(dir)

; reflect off a flat floor (normal = up)
reflected = VEC3.REFLECT(dir, up)

PRINT "Reflected Y: " + STR(VEC3.Y(reflected))

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(20, 20, 35)
    RENDER.BEGIN3D(cam)
        ; draw direction line
        ox = VEC3.X(origin)
        oy = VEC3.Y(origin)
        oz = VEC3.Z(origin)
        rx = ox + VEC3.X(dir) * 4
        ry = oy + VEC3.Y(dir) * 4
        rz = oz + VEC3.Z(dir) * 4
        DRAW3D.LINE(ox, oy, oz, rx, ry, rz, 80, 200, 255, 255)
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

VEC3.FREE(origin)
VEC3.FREE(target)
VEC3.FREE(up)
VEC3.FREE(dir)
VEC3.FREE(reflected)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `VEC3.MAKE(x, y, z)` | Deprecated alias of `VEC3.CREATE`. |

---

## See also

- [VEC2.md](VEC2.md) — 2D vector math
- [VEC_QUAT.md](VEC_QUAT.md) — quaternions and combined examples
- [TRANSFORM.md](TRANSFORM.md) — matrix transforms using Vec3
- [MATH.md](MATH.md) — scalar helpers
