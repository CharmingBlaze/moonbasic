# Vec2, Vec3, and Quat Commands

Vector and quaternion math with heap-handle and scalar-convenience overloads.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create vectors with `VEC3.MAKE` / `VEC2.MAKE`, quaternions with `QUAT.FROMEULER`.
2. Perform arithmetic (`VEC3.ADD`, `VEC3.DOT`, etc.) or use scalar overloads (`VEC3.LENGTH(x,y,z)`).
3. Free handles with `VEC3.FREE` / `VEC2.FREE` / `QUAT.FREE`.

Angles are in **radians**. See also [TRANSFORM.md](TRANSFORM.md).

---

### `VEC3.MAKE(x, y, z)` 
Creates a new 3D vector handle.

---

### `VEC3.X(handle)` / `Y` / `Z` 
Returns the individual red, green, blue, or alpha component (0-255).

---

### `VEC3.SET(handle, x, y, z)` 
Updates the components of an existing vector handle in place.

---

### `VEC3.ADD(a, b)` / `SUB` / `MUL` / `DIV` 
Arithmetic on vector handles. Returns a new handle.

---

### `VEC3.DOT(a, b)` / `VEC3.CROSS(a, b)` 
Vector products. `Dot` returns a float; `Cross` returns a new handle.

---

### `VEC3.LENGTH(handle)` / `VEC3.NORMALIZE(handle)` 
Returns the length of a vector or a new normalized vector handle.

---

### `VEC3.DISTANCE(a, b)` 
Returns the distance between two 3D points.

---

### `VEC3.FREE(handle)` 
Releases the vector handle from the heap and frees its memory.

```basic
v = Vec3.Make(0, 1, 0)
u = Vec3.Make(1, 0, 0)
PRINT "angle rad=" + STR(Vec3.Angle(v, u))
Vec3.Free(v)
Vec3.Free(u)
```

Scalar convenience overloads (no vec3 handle required):

- `VEC3.LENGTH(x, y, z) -> float`
- `VEC3.NORMALIZE(x, y, z) -> handle` (tuple-like 3-float array for destructuring)
- `VEC3.DIST(x1, y1, z1, x2, y2, z2) -> float` — distance between two points; overload with **two vec3 handles** (same as `VEC3.Distance`).
- `VEC3.DISTSQ(x1, y1, z1, x2, y2, z2) -> float` — squared distance (cheap radius checks without `SQRT`).

```basic
dx, dy, dz = VEC3.NORMALIZE(dx, dy, dz)
dist = VEC3.LENGTH(dx, dy, dz)
d = VEC3.DIST(px, py, pz, ex, ey, ez)
IF VEC3.DISTSQ(px, py, pz, ex, ey, ez) < 4.0 THEN
    ; within 2 units
ENDIF
```

---

### `VEC2.MAKE(x, y)` 
Creates a new 2D vector handle.

---

### `VEC2.X(handle)` / `Y` 
Returns the individual red, green, blue, or alpha component (0-255).

---

### `VEC2.SET(handle, x, y)` 
Updates the components of an existing vector handle in place.

---

### `VEC2.ADD(a, b)` / `SUB` / `MUL` 
Arithmetic on vector handles. Returns a new handle.

---

### `VEC2.LENGTH(handle)` / `VEC2.NORMALIZE(handle)` 
Returns the length of a vector or a new normalized vector handle.

---

### `VEC2.FREE(handle)` 
Releases the vector handle from the heap and frees its memory.

Scalar convenience overloads:

- `VEC2.LENGTH(x, y) -> float`
- `VEC2.NORMALIZE(x, y) -> handle` (tuple-like 2-float array)
- `VEC2.MOVE_TOWARD(fromX, fromY, toX, toY, maxDist) -> handle`

```basic
f, s = VEC2.NORMALIZE(f, s)
ex, ez = VEC2.MOVE_TOWARD(ex, ez, px, pz, chaseSpeed * dt)
dist = VEC2.LENGTH(ex - px, ez - pz)
```

---

## Quat

### `QUAT.IDENTITY()` 
Returns an identity quaternion handle.

---

### `QUAT.FROMAXISANGLE(axis, angle)` / `QUAT.FROMVEC3TOVEC3(from, to)` 
Creates a quaternion from an axis+angle or from two direction vectors.

---

### `QUAT.FROMEULER(p, y, r)` 
Creates a quaternion from Euler angles in **radians**.

---

### `QUAT.TOEULER(q)` 
Returns a 3-float array handle `[p, y, r]` from a quaternion.

---

### `QUAT.SLERP(a, b, t)` 
Spherical interpolation between quaternions `a` and `b` by factor `t` (0–1). Returns a new handle.

---

### `QUAT.FREE(handle)` 
Releases the quaternion handle from the heap and frees its memory.

```basic
q = Quat.FromEuler(0, PI() / 4, 0)
e = Quat.ToEuler(q)
PRINT "roll=" + STR(Vec3.X(e))
Vec3.Free(e)
m = Quat.ToMat4(q)
Mesh.Draw(cube, mat, m)
Transform.Free(m)
Quat.Free(q)
```

---

## Full Example

Rotating a mesh using a quaternion built from euler angles each frame.

```basic
WINDOW.OPEN(960, 540, "Vec/Quat Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 3, -8)
CAMERA.SETTARGET(cam, 0, 0, 0)

mesh = MESH.CREATECUBE(2, 2, 2)
mat  = MATERIAL.CREATEDEFAULT()
yaw  = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt  = TIME.DELTA()
    yaw = yaw + 1.2 * dt          ; radians/sec

    q = QUAT.FROMEULER(0, yaw, 0)
    m = QUAT.TOMAT4(q)
    QUAT.FREE(q)

    RENDER.CLEAR(20, 20, 40)
    RENDER.BEGIN3D(cam)
        MESH.DRAW(mesh, mat, m)
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    TRANSFORM.FREE(m)
    RENDER.FRAME()
WEND

MESH.UNLOAD(mesh)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [TRANSFORM.md](TRANSFORM.md) — transform matrices
- [MAT4.md](MAT4.md) — legacy matrix spellings
- [MATH.md](MATH.md) — scalar trig and lerp helpers
