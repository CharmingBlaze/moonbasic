# Game math helpers (2D / 3D)

High-frequency **distance**, **easing**, and **angle** helpers for gameplay code. They live in the **`MATH.*`** namespace (and **flat aliases** without the prefix, same as `LERP` / `MATH.LERP`). **No Raylib** — safe in stub builds; use alongside [LESS_MATH.md](LESS_MATH.md), [MATH.md](MATH.md), and collision helpers in [COLLISION.md](COLLISION.md) / `BOXCOLLIDE`, `DISTANCE2D`, …

**One-page map:** **`MATH.APPROACH`**, **`MATH.LERP`**, **`MATH.CURVE`**, **`MATH.NEWX`** / **`MATH.NEWZ`** (XZ heading + radians), **`MATH.ANGLEDIFF`** (degrees) — see the **“Short game-logic helpers”** table in [MATH.md](MATH.md).

---

## Horizontal distance (3D on XZ)

Use when **Y does not matter** (aggro radius, top-down distance, AI on a flat plane):

| Command | Returns |
|--------|---------|
| **`HDIST(x1, z1, x2, z2)`** | `sqrt((x2-x1)² + (z2-z1)²)` |
| **`MATH.HDIST`** | same |
| **`HDISTSQ(...)`** | squared distance — compare to `r*r` without `sqrt` |
| **`MATH.HDISTSQ`** | same |

---

## 2D distance under `MATH.*`

Same numbers as **`DISTANCE2D`** / **`DISTANCESQ2D`** in the game module; exposed here so everything “math-shaped” is discoverable next to **`MATH.DIST`** patterns:

| Command | Returns |
|--------|---------|
| **`DIST2D(x1, y1, x2, y2)`** | Euclidean distance in 2D |
| **`DISTSQ2D(...)`** | squared distance |

---

## Yaw from flat direction

| Command | Returns |
|--------|---------|
| **`YAWFROMXZ(dx, dz)`** | **radians** — `atan2(dx, dz)` consistent with **`MOVEX` / `MOVEZ`** / **`INPUT.MOVEDIR`** (forward uses `sin(yaw)`, `cos(yaw)` on X/Z). |

Use to face a target or stick from a world-space delta on the ground plane.

---

## Angles in radians

| Command | Returns |
|--------|---------|
| **`ANGLEDIFFRAD(a, b)`** | Shortest signed difference **b − a** in **radians** (−π..π). Use for “how far to rotate” without degree/`ANGLEDIFF` conversions. |

For **interpolation** toward an angle, use **`MATH.LERPANGLE`** (already in [MATH.md](MATH.md)).

---

## Easing: smootherstep

| Command | Returns |
|--------|---------|
| **`SMOOTHERSTEP(edge0, edge1, x)`** | Ken Perlin **smootherstep**: clamp **x** to [edge0, edge1], normalize to **t ∈ [0,1]**, then **6t⁵ − 15t⁴ + 10t³**. Smoother ends than **`SMOOTHSTEP`**. |
| **`MATH.SMOOTHERSTEP`** | same |

---

## See also

- [MATH.md](MATH.md) — `LERP`, `SMOOTHSTEP`, `LERPANGLE`, `DIST3D`, …
- [LESS_MATH.md](LESS_MATH.md) — `INPUT.MOVEDIR`, `VEC2.*`, terrain snap
- [EASY_LANGUAGE.md](../EASY_LANGUAGE.md) — design stance: helpers + full math
