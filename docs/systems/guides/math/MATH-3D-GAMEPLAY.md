# 3D game math — ground plane, yaw, and horizontal distance

> Scalar `MATH.*` for 3D games where gameplay mostly lives on the **XZ floor** (yaw, horizontal range, step forward).

**Namespaces:** `MATH` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **3D vectors:** [VEC3-MATH.md](VEC3-MATH.md) · **2D screen:** [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use scalar 3D gameplay math](#when-to-use-scalar-3d-gameplay-math)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Horizontal distance (XZ)](#horizontal-distance-xz)
- [Yaw and facing](#yaw-and-facing)
- [Step on the ground plane](#step-on-the-ground-plane)
- [Full example — face and chase on XZ](#full-example--face-and-chase-on-xz)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | `HDIST`, `HDISTSQ`, `YAWFROMXZ`, `NEWX`/`NEWZ`, `ANGLETO`, full 3D dist via scalars |
| **You need first** | Entities on XZ ([ENTITY-SYSTEM.md](../ENTITY-SYSTEM.md)) |
| **Typical games** | Third-person action, RTS-style ground units, open world |
| **Upgrade when** | Full 3D flight, slopes, dot/cross → [VEC3-MATH.md](VEC3-MATH.md) |

**Why XZ helpers exist:** Most chase, aggro, and WASD movement ignore height for distance checks. `HDIST` skips Y so flying enemies do not skew range on the map.

---

## When to use scalar 3D gameplay math

**Use when:**

- Enemy aggro uses ground distance, not vertical offset.
- Character turns to face target yaw on the floor.
- WASD moves along camera yaw (`INPUT.MOVEDIR` pattern).
- Spawn points on a circle in the world (`CIRCLEPOINT` on XZ).

**Skip when:**

- Vertical gameplay (jump arcs only) — still use Y separately.
- True 3D vectors every frame — [VEC3-MATH.md](VEC3-MATH.md).
- Physics handles motion — [COLLISION-3D.md](../COLLISION-3D.md).

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Range on map (ignore height) | `MATH.HDIST` / `HDISTSQ` | Full 3D distance for flat arenas |
| Face target on ground | `MATH.YAWFROMXZ(dx, dz)` | `ATAN2` with wrong axis order |
| Step forward from yaw | `MATH.NEWX` / `NEWZ` | Manual sin/cos each time |
| Full 3D gap (flying) | `VEC3.DIST` or 6-float dist | `HDIST` when Y matters |
| Slope walking | `CHAR.*` | Raw yaw only |

---

## Core workflow

1. **Positions** — `ENTITY` x,y,z or `px, py, pz`.
2. **Horizontal delta** — `dx = tx - px`, `dz = tz - pz`.
3. **Range** — `HDISTSQ(px, pz, tx, tz)` vs `radius²`.
4. **Facing** — `yaw = MATH.YAWFROMXZ(dx, dz)` → `ENTITY.SETYAW` or camera.
5. **Move** — `px = MATH.NEWX(px, yaw, speed * APP.DELTA())` (see registry arity).
6. **Height** — `py` from terrain `GETHEIGHT` or character controller.

---

## Horizontal distance (XZ)

```basic
; Aggro — ignore vertical separation
rangeSq = MATH.HDISTSQ(px, pz, ex, ez)
IF rangeSq < 15 * 15 THEN wakeAI = 1

; Exact meters on ground
groundDist = MATH.HDIST(px, pz, ex, ez)
```

| Command | Why |
|---------|-----|
| `MATH.HDIST(x1,z1,x2,z2)` | √(Δx² + Δz²) |
| `MATH.HDISTSQ(...)` | Compare to r² without sqrt |
| `VEC3.DIST` / `DISTSQ` | Full 3D when Y matters |

---

## Yaw and facing

**Yaw** is rotation around **Y**, in **radians**, consistent with `MOVEX`/`MOVEZ` and `INPUT.MOVEDIR`.

```basic
dx = tx - px
dz = tz - pz
yaw = MATH.YAWFROMXZ(dx, dz)
ENTITY.SETYAW(hero, yaw)
```

| Command | Why |
|---------|-----|
| `MATH.YAWFROMXZ(dx, dz)` | Aim on XZ from delta |
| `MATH.ANGLETO(x1,z1,x2,z2)` | Yaw from two world points |
| `MATH.ANGLEDIFFRAD(a, b)` | How far to rotate (radians) |
| `MATH.LERPANGLE(a, b, t)` | Smooth turn |
| `MATH.DEGPERSEC(deg, dt)` | Turn speed in deg/s → this frame |

Camera orbit: [CAMERA-AND-INPUT.md](../CAMERA-AND-INPUT.md). Character slopes: [CHARACTER-3D-WALKING.md](../CHARACTER-3D-WALKING.md).

---

## Step on the ground plane

```basic
yaw = MATH.YAWFROMXZ(dx, dz)
step = 5 * APP.DELTA()
px = MATH.NEWX(px, yaw, step)
pz = MATH.NEWZ(pz, yaw, step)
```

**Why `NEWX`/`NEWZ`:** Same convention as movement helpers — fewer sign mistakes than raw `SIN`/`COS`.

Spawn ring on ground:

```basic
pt = MATH.CIRCLEPOINT(cx, cz, radius, i, count)
; returns handle or components — see COMMAND_REGISTRY
```

---

## Full example — face and chase on XZ

```basic
APP.OPEN(800, 600, "3D ground math")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 8, -12)
CAMERA.LOOKAT(cam, 0, 0, 0)

hero = ENTITY.CREATECUBE(1, 1, 1)
hero.pos(0, 0, 0)
target = ENTITY.CREATECUBE(0.5, 0.5, 0.5)
target.pos(6, 0, 4)

WHILE NOT APP.SHOULDCLOSE()
    hx = hero.x()
    hz = hero.z()
    tx = target.x()
    tz = target.z()

    dx = tx - hx
    dz = tz - hz
    dist = MATH.HDIST(hx, hz, tx, tz)

    IF dist > 0.5 THEN
        yaw = MATH.YAWFROMXZ(dx, dz)
        step = 3 * APP.DELTA()
        hero.pos(MATH.NEWX(hx, yaw, step), hero.y(), MATH.NEWZ(hz, yaw, step))
        ENTITY.SETYAW(hero, yaw)
    ENDIF

    RENDER.CLEAR(25, 28, 35)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

APP.CLOSE()
```

Verify `NEWX`/`NEWZ` arity in [COMMAND_REGISTRY.md](../../COMMAND_REGISTRY.md) — signatures may include current position + yaw + distance.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| `HDIST` with x,y screen coords | Use `DIST2D` for screen ([MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md)) |
| Yaw in degrees with `SIN` | Convert with `DEG2RAD` or use `SIND` |
| Full 3D dist for flat AI | Use `HDIST` on XZ |
| Move without delta | Multiply step by `APP.DELTA()` |
| Ignore terrain Y | Snap Y with [TERRAIN-OPEN-WORLD.md](../TERRAIN-OPEN-WORLD.md) |

---

## See also

- [VEC3-MATH.md](VEC3-MATH.md) — dot, cross, 3D direction
- [TERRAIN-OPEN-WORLD.md](../TERRAIN-OPEN-WORLD.md) — height at x,z
- [reference/LESS_MATH.md](../../../reference/LESS_MATH.md) — camera-relative shortcuts
- [reference/GAME_MATH_HELPERS.md](../../../reference/GAME_MATH_HELPERS.md)
