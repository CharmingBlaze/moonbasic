# 2D game math — screen space, distances, and aim angles

> Scalar `MATH.*` helpers for platformers, top-down shooters, and HUD logic on the X/Y plane.

**Namespaces:** `MATH` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **Vectors:** [VEC2-MATH.md](VEC2-MATH.md) · **3D ground plane:** [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use scalar 2D math](#when-to-use-scalar-2d-math)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Distance on X/Y](#distance-on-xy)
- [Angles and facing (2D)](#angles-and-facing-2d)
- [Movement helpers](#movement-helpers)
- [Full example — top-down chase](#full-example--top-down-chase)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | `DIST2D`, `DISTSQ2D`, `ATAN2`, `ANGLETO`, `CIRCLEPOINT`, clamp/lerp on scalars |
| **You need first** | Loop + delta time ([GAME-LOOP-AND-RENDERING.md](../GAME-LOOP-AND-RENDERING.md)) |
| **Typical games** | 2D platformer, twin-stick, mouse aim, UI bars |
| **Upgrade when** | Many vector ops per frame → [VEC2-MATH.md](VEC2-MATH.md) |

**Why not only `VEC2`:** Four floats and `MATH.DIST2D` are enough for aggro radius and mouse aim. Vectors add normalize/rotate when direction math repeats.

---

## When to use scalar 2D math

**Use when:**

- Player and enemy positions are `px, py` floats.
- You compare “is mouse near button?” with distance threshold.
- You rotate a sprite with `ATAN2(dy, dx)`.
- You spawn loot on a ring with `CIRCLEPOINT`.

**Skip when:**

- Every frame you normalize, rotate, and push circles apart → [VEC2-MATH.md](VEC2-MATH.md).
- Gameplay is on XZ ground in 3D → [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md).

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Distance between two points | `MATH.DIST2D(x1,y1,x2,y2)` | `SQRT` by hand |
| “Inside radius?” without sqrt | `MATH.DISTSQ2D` vs `r*r` | `DIST2D` every check |
| Angle from A to B (screen) | `MATH.ATAN2(dy, dx)` | `ATAN(dy/dx)` alone |
| Angle between two headings (deg) | `MATH.ANGLEDIFF(a, b)` | Subtract angles raw |
| Smooth value toward target | `MATH.APPROACH` / `LERP` | [INTERPOLATION-AND-EASING.md](INTERPOLATION-AND-EASING.md) |
| Points on a loot ring | `MATH.CIRCLEPOINT` | Manual sin/cos loop |

---

## Core workflow

1. **Read positions** — sprite `x,y` or `INPUT.MOUSEX/MOUSEY`.
2. **Compare range** — `DISTSQ2D` for AI wake; `DIST2D` when you need exact meters.
3. **Aim** — `angle = ATAN2(ty - py, tx - px)` (radians).
4. **Move** — `px += cos(angle) * speed * APP.DELTA()` (or `VEC2` normalize).
5. **Clamp** — `MATH.CLAMP` for health, stay on screen.

Always multiply speed by **`APP.DELTA()`**.

---

## Distance on X/Y

```basic
dx = tx - px
dy = ty - py
distSq = MATH.DISTSQ2D(px, py, tx, ty)
IF distSq < 64 * 64 THEN
    ; within 64 px — no sqrt needed
ENDIF
dist = MATH.DIST2D(px, py, tx, ty)
```

| Command | Why |
|---------|-----|
| `MATH.DIST2D(x1,y1,x2,y2)` | Euclidean distance on screen / tile plane |
| `MATH.DISTSQ2D(...)` | Cheaper inside-radius tests |
| `MATH.HDIST` / `HDISTSQ` | **XZ only** — use in 3D top-down ([MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md)) |

---

## Angles and facing (2D)

**Radians** for `SIN`/`COS`/`ATAN2`. **Degrees** for `ANGLEDIFF` and `SIND`/`COSD`.

```basic
; Face mouse (screen Y often increases downward)
mx = INPUT.MOUSEX()
my = INPUT.MOUSEY()
aim = MATH.ATAN2(my - py, mx - px)
SPRITE.SETROTATION(hero, MATH.RAD2DEG(aim))
```

| Command | Why |
|---------|-----|
| `MATH.ATAN2(y, x)` | Signed angle from x-axis to point |
| `MATH.ANGLETO(x1,y1,x2,y2)` | Angle toward second point (check registry arity) |
| `MATH.ANGLEDIFF(a, b)` | Shortest turn in **degrees** (−180…180) |
| `MATH.ANGLEDIFFRAD(a, b)` | Same in radians |
| `MATH.LERPANGLE(a, b, t)` | Smooth rotation |
| `MATH.WRAPANGLE` / `WRAPANGLE180` | Keep angle in range |

Details: [ANGLES-AND-ROTATION.md](ANGLES-AND-ROTATION.md).

---

## Movement helpers

| Command | Why |
|---------|-----|
| `MATH.LERP(a, b, t)` | Smooth follow on one axis |
| `MATH.APPROACH(cur, target, step)` | Move at most `step` per frame |
| `MATH.CLAMP(v, lo, hi)` | Stay in arena |
| `MATH.PINGPONG(t, len)` | Oscillating platforms |
| `MATH.WRAP(x, lo, hi)` | Wrap screen torus |

For diagonal keyboard move at constant speed, normalize `(dx,dy)` — see [VEC2-MATH.md](VEC2-MATH.md).

---

## Full example — top-down chase

**Runnable:** [examples/guides/math/math_2d_chase.mb](../../../examples/guides/math/math_2d_chase.mb)

```basic
; Check: moonbasic --check examples/guides/math/math_2d_chase.mb
; Run:   moonrun examples/guides/math/math_2d_chase.mb

APP.OPEN(640, 480, "2D chase")
APP.SETFPS(60)

px = 100
py = 240
tx = 500
ty = 240
speed = 180

WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYDOWN(KEY_W) THEN ty = ty - speed * APP.DELTA()
    IF INPUT.KEYDOWN(KEY_S) THEN ty = ty + speed * APP.DELTA()
    IF INPUT.KEYDOWN(KEY_A) THEN tx = tx - speed * APP.DELTA()
    IF INPUT.KEYDOWN(KEY_D) THEN tx = tx + speed * APP.DELTA()

    dist = MATH.DIST2D(px, py, tx, ty)
    IF dist > 4 THEN
        t = MATH.CLAMP(speed * APP.DELTA() / dist, 0, 1)
        px = MATH.LERP(px, tx, t)
        py = MATH.LERP(py, ty, t)
    ENDIF

    DRAW.RECTANGLE(tx - 6, ty - 6, 12, 12, 255, 80, 80, 255)
    DRAW.RECTANGLE(px - 8, py - 8, 16, 16, 80, 200, 255, 255)
    DRAW.TEXT("dist " + dist, 10, 10, 16, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| `ATAN(dy/dx)` without quadrant | Use `ATAN2(dy, dx)` |
| Distance every AI without sqrt need | Compare `DISTSQ2D` to `radius*radius` |
| Forgot `APP.DELTA()` | Motion tied to FPS |
| Screen Y down vs world Y up | Flip dy in aim math for sprites |
| Used `HDIST` for screen X/Y | `HDIST` is XZ ground in 3D |

---

## See also

- [VEC2-MATH.md](VEC2-MATH.md) — normalize, rotate, pushout
- [SPRITES-TILEMAPS-2D.md](../SPRITES-TILEMAPS-2D.md) — tile coords
- [COLLISION-2D.md](../COLLISION-2D.md) — `VEC2` overlap tests
- [reference/GAME_MATH_HELPERS.md](../../../reference/GAME_MATH_HELPERS.md)
