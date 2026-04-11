# Vec2, Vec3, and Quat

These commands support both classic heap-handle vector math and scalar convenience overloads for gameplay code.

**Convention:** Angles are in **radians** unless noted. Free heap handles with `Vec2.Free`, `Vec3.Free`, or `Quat.Free` when you no longer need them (or rely on `FREE.ALL` at shutdown).

---

### `Vec3.Make(x, y, z)`
Creates a new 3D vector handle.

### `Vec3.X(handle)` / `Y` / `Z`
Returns the individual red, green, blue, or alpha component (0-255).

### `Vec3.Set(handle, x, y, z)`
Updates the components of an existing vector handle in place.

---

### `Vec3.Add(a, b)` / `Sub` / `Mul` / `Div`
Arithmetic on vector handles. Returns a new handle.

### `Vec3.Dot(a, b)` / `Cross(a, b)`
Vector products. `Dot` returns a float; `Cross` returns a new handle.

### `Vec3.Length(handle)` / `Normalize(handle)`
Returns the length of a vector or a new normalized vector handle.

### `Vec3.Distance(a, b)`
Returns the distance between two 3D points.

---

### `Vec3.Free(handle)`
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

### `Vec2.Make(x, y)`
Creates a new 2D vector handle.

### `Vec2.X(handle)` / `Y`
Returns the individual red, green, blue, or alpha component (0-255).

### `Vec2.Set(handle, x, y)`
Updates the components of an existing vector handle in place.

---

### `Vec2.Add(a, b)` / `Sub` / `Mul`
Arithmetic on vector handles. Returns a new handle.

### `Vec2.Length(handle)` / `Normalize(handle)`
Returns the length of a vector or a new normalized vector handle.

### `Vec2.Free(handle)`
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

### `Quat.Identity()`
Returns an identity quaternion handle.

### `Quat.Make(x, y, z, w)`
Creates a new quaternion handle with explicit components.

### `Quat.FromEuler(p, y, r)`
Creates a quaternion from Euler angles in **radians**.

---

### `Quat.ToEuler(q)`
Returns a 3-float array handle `[p, y, r]` from a quaternion.

### `Quat.Lerp(a, b, t)` / `Quat.Slerp(a, b, t)`
Linear and spherical interpolation between quaternions `a` and `b` by factor `t` (0–1). Returns a new handle.

### `Quat.Free(handle)`
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

## See also

- [TRANSFORM.md](TRANSFORM.md) — transform matrices; legacy spellings in [MAT4.md](MAT4.md).
