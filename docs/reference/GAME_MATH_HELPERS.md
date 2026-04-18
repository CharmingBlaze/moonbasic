# Game math helpers (2D / 3D)

High-frequency **distance**, **easing**, and **angle** helpers for gameplay code. They live in the **`MATH.*`** namespace (and **flat aliases** without the prefix, same as `LERP` / `MATH.LERP`). **No Raylib** — safe in stub builds; use alongside [LESS_MATH.md](LESS_MATH.md), [MATH.md](MATH.md), and collision helpers in [COLLISION.md](COLLISION.md) / `BOXCOLLIDE`, `DISTANCE2D`, …

**One-page map:** **`MATH.APPROACH`**, **`MATH.LERP`**, **`MATH.CURVE`**, **`MATH.NEWX`** / **`MATH.NEWZ`** (XZ heading + radians), **`MATH.ANGLEDIFF`** (degrees) — see the **"Short game-logic helpers"** table in [MATH.md](MATH.md).

## Core Workflow

- Use `HDIST` / `DIST2D` to check aggro radii without a physics query.
- Use `LERP` / `APPROACH` to smooth movement and color transitions.
- Use `YAWFROMXZ` to face an enemy: `ENTITY.SETYAW(e, YAWFROMXZ(dx, dz))`.
- Use `SMOOTHSTEP` / `SMOOTHERSTEP` for eased UI animations.

---

## Horizontal distance (3D on XZ)

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

## Full Example

Smooth color pulse and approach-to-target using game math helpers.

```basic
WINDOW.OPEN(800, 450, "GameMath Demo")
WINDOW.SETFPS(60)

current = 0.0
target  = 1.0
t       = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    t  = t + dt

    ; APPROACH: move current toward target by max 2*dt per frame
    current = APPROACH(current, target, 2.0 * dt)
    IF ABS(current - target) < 0.01 THEN target = 1.0 - target

    ; SMOOTHSTEP for eased color
    brightness = INT(SMOOTHSTEP(0.0, 1.0, current) * 255)

    RENDER.CLEAR(0, 0, 0)
    DRAW.RECTANGLE(200, 150, 400, 150, brightness, INT(brightness * 0.5), 255 - brightness, 255)
    DRAW.TEXT("APPROACH: " + STR(current), 10, 10, 18, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [MATH.md](MATH.md) — `LERP`, `SMOOTHSTEP`, `LERPANGLE`, `DIST3D`, …
- [LESS_MATH.md](LESS_MATH.md) — `INPUT.MOVEDIR`, `VEC2.*`, terrain snap
- [EASY_LANGUAGE.md](../EASY_LANGUAGE.md) — design stance: helpers + full math
