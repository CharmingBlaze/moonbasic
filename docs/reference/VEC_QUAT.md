# Vec2, Vec3, and Quat — Heap Math Handles

These commands wrap Raylib’s `raymath` types as **handles** on the VM heap. They are useful for physics-facing math, camera helpers, and quaternion-based rotation without gimbal lock.

**Convention:** Angles are in **radians** unless noted. Free handles with `Vec2.Free`, `Vec3.Free`, or `Quat.Free` when you no longer need them (or rely on `Heap.FreeAll` at shutdown).

---

## Vec3

| Command | Notes |
|--------|--------|
| `Vec3.Make(x#, y#, z#)` | New vec3 handle. |
| `Vec3.Free(h)` | |
| `Vec3.X(h)` / `Y` / `Z` | Components. |
| `Vec3.Set(h, x#, y#, z#)` | Mutate in place. |
| `Vec3.Add` / `Sub` / `Mul` / `Div` / `Dot` / `Cross` / `Length` / `Normalize` / `Lerp` / `Distance` / `Reflect` / `Negate` / `Equals` | See compiler manifest for arities. |
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

---

## Vec2

| Command | Notes |
|--------|--------|
| `Vec2.TransformMat4(v, mat)` | Returns new vec2 (homogeneous transform). |

```basic
; After building a 2D camera matrix from Camera2D.GetMatrix
p = Vec2.Make(100, 200)
p2 = Vec2.TransformMat4(p, cam_mat)
Vec2.Free(p)
; use p2, then Vec2.Free(p2)
```

---

## Quat

| Command | Notes |
|--------|--------|
| `Quat.Identity()` | |
| `Quat.FromEuler(pitch#, yaw#, roll#)` | Radians (Raylib `QuaternionFromEuler`). |
| `Quat.FromAxisAngle(ax#, ay#, az#, angle#)` | Axis need not be normalized; Raylib normalizes internally. |
| `Quat.Multiply(a, b)` | Returns new quaternion. |
| `Quat.Slerp(a, b, t#)` | Spherical interpolation, `t` in 0..1. |
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
