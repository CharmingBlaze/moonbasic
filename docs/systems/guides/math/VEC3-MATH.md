# 3D vector math — VEC3 directions, products, and transforms

> Positions and directions in 3D: add, dot, cross, normalize, reflect, and rotate with quaternions.

**Namespaces:** `VEC3` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **Ground scalars:** [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md) · **Reference:** [reference/VEC3.md](../../../reference/VEC3.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use VEC3](#when-to-use-vec3)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Key commands](#key-commands)
- [Patterns](#patterns)
- [Full example — direction to target](#full-example--direction-to-target)
- [Memory notes](#memory-notes)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | 3D handles; dot/cross; distance; lerp; reflect/project; mat4 transform |
| **You need first** | 3D entities or camera ([ENTITY-SYSTEM.md](../ENTITY-SYSTEM.md)) |
| **Typical games** | Aim vectors, knockback, surface normals, lighting direction |
| **Flat movement** | Often [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md) `HDIST` is enough |

---

## When to use VEC3

**Use when:**

- Knockback along `(target - self)` normalized.
- “Is enemy in front?” with `DOT(forward, toEnemy)`.
- Cross product for surface normal or strafe axis.
- Reflect velocity off wall normal.
- Transform point by camera matrix.

**Skip when:**

- Only XZ range and yaw — [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md).
- Character on slopes — `CHAR.*` integrates floor normal.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Ground aggro radius | `MATH.HDISTSQ` | Full `VEC3.DIST` |
| Direction A → B in 3D | `SUB` + `NORMALIZE` | Manual component divide |
| Facing cone | `VEC3.DOT` | Angle from components by hand |
| Strafe axis | `VEC3.CROSS(up, forward)` | Guess perpendicular |
| Bounce off wall | `VEC3.REFLECT` | Flip all velocity signs |
| Entity rotation | `QUAT` / entity yaw | Manual Euler every frame |

---

## Core workflow

1. **Create or scalar** — `VEC3.CREATE(x,y,z)` or pass 6 floats to `DIST`.
2. **Difference** — `to = VEC3.SUB(targetPos, selfPos)`.
3. **Normalize** — `dir = VEC3.NORMALIZE(to)` for unit step.
4. **Apply** — `pos = VEC3.ADD(pos, VEC3.MUL(dir, speed * dt))`.
5. **Free** temporary handles.

Scalar overloads: `VEC3.LENGTH(x,y,z)`, `VEC3.DIST(x1,y1,z1,x2,y2,z2)`, `VEC3.DISTSQ(...)`.

---

## Key commands

| Command | Why |
|---------|-----|
| `VEC3.CREATE(x,y,z)` | New handle |
| `VEC3.ADD` / `SUB` / `MUL` / `DIV` | Arithmetic |
| `VEC3.NEGATE` | Flip direction |
| `VEC3.LENGTH` | Magnitude |
| `VEC3.NORMALIZE` | Unit vector |
| `VEC3.DIST` / `DISTANCE` / `DISTSQ` | Separation |
| `VEC3.DOT(a,b)` | Facing, projection scale |
| `VEC3.CROSS(a,b)` | Perpendicular vector |
| `VEC3.LERP(a,b,t)` | Blend positions |
| `VEC3.PROJECT(a, onto)` | Slide along surface |
| `VEC3.REFLECT(dir, normal)` | Bounce |
| `VEC3.ROTATEBYQUAT(v, quat)` | Orientation math |
| `VEC3.TRANSFORMMAT4(v, mat)` | Camera / model space |
| `VEC3.ORTHONORMALIZE` | Build basis from vectors |
| `VEC3.X` / `Y` / `Z` | Read components |
| `VEC3.FREE` | Release handle |

Aliases: `VECADD`, `VECSUB`, `VECSCALE`, `VECLENGTH`, `VECNORMALIZE`, `VECDOT`, `VECCROSS`.

---

## Patterns

### Direction to target (3D)

```basic
dx = tx - px
dy = ty - py
dz = tz - pz
dist = VEC3.DIST(px, py, pz, tx, ty, tz)
IF dist > 0.01 THEN
    dir = VEC3.NORMALIZE(dx, dy, dz)
    px = px + VEC3.X(dir) * speed * APP.DELTA()
    py = py + VEC3.Y(dir) * speed * APP.DELTA()
    pz = pz + VEC3.Z(dir) * speed * APP.DELTA()
    VEC3.FREE(dir)
ENDIF
```

### In front of player (dot cone)

```basic
fwd = VEC3.CREATE(sin(yaw), 0, cos(yaw))
toEnemy = VEC3.SUB(enemyPos, selfPos)
toN = VEC3.NORMALIZE(toEnemy)
IF VEC3.DOT(fwd, toN) > 0.7 THEN spotted = 1
```

### Compare range cheaply

```basic
IF VEC3.DISTSQ(px,py,pz, ex,ey,ez) < r*r THEN ...
```

---

## Full example — direction to target

```basic
APP.OPEN(640, 480, "VEC3 aim")
APP.SETFPS(60)

px = 0
py = 1
pz = 0
tx = 5
ty = 2
tz = 4
speed = 2

WHILE NOT APP.SHOULDCLOSE()
    dx = tx - px
    dy = ty - py
    dz = tz - pz
    dist = VEC3.DIST(px, py, pz, tx, ty, tz)
    IF dist > 0.05 THEN
        d = VEC3.NORMALIZE(dx, dy, dz)
        step = speed * APP.DELTA()
        px = px + VEC3.X(d) * step
        py = py + VEC3.Y(d) * step
        pz = pz + VEC3.Z(d) * step
        VEC3.FREE(d)
    ENDIF
    DRAW.TEXT("dist " + dist, 10, 10, 16, 255, 255, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
```

---

## Memory notes

- Free handles from `CREATE`, `NORMALIZE`, `ADD`, etc. when chaining many ops.
- Hot path: scalar `DIST`/`DISTSQ`/`NORMALIZE(x,y,z)` overloads reduce heap churn.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Normalize zero vector | Check `LENGTH` or `DIST` first |
| `HDIST` for flying enemies | Use full `VEC3.DIST` when Y matters |
| Cross product order | `CROSS(a,b)` ≠ `CROSS(b,a)` — pick consistent basis |
| Confuse dot with cross | Dot → scalar facing; cross → perpendicular vector |
| `MAKE` deprecated | Use `CREATE` |

---

## See also

- [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md) — XZ scalars
- [ANGLES-AND-ROTATION.md](ANGLES-AND-ROTATION.md) — quaternions
- [CAMERA-AND-INPUT.md](../CAMERA-AND-INPUT.md) — camera rays
- [reference/VEC_QUAT.md](../../../reference/VEC_QUAT.md)
