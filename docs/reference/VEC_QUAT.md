# Vec2, Vec3, and Quat

These commands support both classic heap-handle vector math and scalar convenience overloads for gameplay code.

**Convention:** Angles are in **radians** unless noted. Free heap handles with `Vec2.Free`, `Vec3.Free`, or `Quat.Free` when you no longer need them (or rely on `FREE.ALL` at shutdown).

---

## Vec3

| Command | Notes |
|--------|--------|
| `Vec3.Make(x, y, z)` | New vec3 handle. |
| `Vec3.Free(h)` | |
| `Vec3.X(h)` / `Y` / `Z` | Components. |
| `Vec3.Set(h, x, y, z)` | Mutate in place. |
| `Vec3.Add` / `Sub` / `Mul` / `Div` / `Dot` / `Cross` / `Length` / `Normalize` / `Lerp` / `Distance` / `Reflect` / `Negate` / `Equals` | Handle-based overloads. |
| `Vec3.TransformMat4(v, mat)` | Returns new vec3: `Vector3Transform`. |
| `Vec3.Angle(a, b)` | Returns float, radians between vectors. |
| `Vec3.Project(v, onto)` | Returns new vec3 projection of `v` onto `onto`. |
| `Vec3.OrthoNormalize(v1, v2)` | **Mutates both** handles (Gram–Schmidt). |
| `Vec3.RotateByQuat(v, q)` | Returns new vec3. |

```basic
v = Vec3.Make(0, 1, 0)
u = Vec3.Make(1, 0, 0)
PRINT "angle rad=" + STR$(Vec3.Angle(v, u))
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
dist# = VEC3.LENGTH(dx, dy, dz)
d# = VEC3.DIST(px#, py#, pz#, ex#, ey#, ez#)
IF VEC3.DISTSQ(px#, py#, pz#, ex#, ey#, ez#) < 4.0 THEN
    ; within 2 units
ENDIF
```

---

## Vec2

| Command | Notes |
|--------|--------|
| `Vec2.Make(x, y)` / `Vec2.Free(h)` / `Vec2.X(h)` / `Vec2.Y(h)` / `Vec2.Set(h, x, y)` | Handle-based vec2 API. |
| `Vec2.Add` / `Sub` / `Mul` / `Length` / `Normalize` / `Lerp` / `Distance` / `Angle` / `Rotate` | Handle-based overloads. |
| `Vec2.TransformMat4(v, mat)` | Returns new vec2 (homogeneous transform). |

Scalar convenience overloads:

- `VEC2.LENGTH(x, y) -> float`
- `VEC2.NORMALIZE(x, y) -> handle` (tuple-like 2-float array)
- `VEC2.MOVE_TOWARD(fromX, fromY, toX, toY, maxDist) -> handle`

```basic
f#, s# = VEC2.NORMALIZE(f#, s#)
ex#, ez# = VEC2.MOVE_TOWARD(ex#, ez#, px#, pz#, chaseSpeed# * dt#)
dist# = VEC2.LENGTH(ex# - px#, ez# - pz#)
```

---

## Quat

| Command | Notes |
|--------|--------|
| `Quat.Identity()` | |
| `Quat.FromEuler(pitch, yaw, roll)` | Radians (Raylib `QuaternionFromEuler`). |
| `Quat.FromAxisAngle(ax, ay, az, angle)` | Axis need not be normalized; Raylib normalizes internally. |
| `Quat.Multiply(a, b)` | Returns new quaternion. |
| `Quat.Slerp(a, b, t)` | Spherical interpolation, `t` in 0..1. |
| `Quat.ToMat4(q)` | Returns new transform matrix handle (same type as `Transform.*`). |
| `Quat.ToEuler(q)` | Returns new **Vec3** handle: **X=roll, Y=pitch, Z=yaw** (radians), per Raylib `QuaternionToEuler`. |
| `Quat.FromVec3ToVec3(from, to)` | Shortest rotation from direction `from` to `to` (new quaternion). |
| `Quat.FromMat4(mat)` | Rotation part of matrix → quaternion. |
| `Quat.Transform(q, mat)` | `QuaternionTransform` — applies 4×4 matrix to quaternion. |
| `Quat.Normalize` / `Quat.Invert` | Return new quaternions. |
| `Quat.Free(q)` | |

```basic
q = Quat.FromEuler(0, PI() / 4, 0)
e = Quat.ToEuler(q)
PRINT "roll=" + STR$(Vec3.X(e))
Vec3.Free(e)
m = Quat.ToMat4(q)
Mesh.Draw(cube, mat, m)
Transform.Free(m)
Quat.Free(q)
```

---

## See also

- [TRANSFORM.md](TRANSFORM.md) — transform matrices; legacy spellings in [MAT4.md](MAT4.md).
