# Angles and rotation — radians, yaw, and quaternions

> Turn sprites, entities, and cameras without wrap bugs; convert deg/rad; use quats when Euler angles gimbal.

**Namespaces:** `MATH` · `VEC3` · `QUAT` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **2D aim:** [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md) · **3D yaw:** [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [Radians vs degrees](#radians-vs-degrees)
- [Choose the right tool](#choose-the-right-tool)
- [Wrap and shortest turn](#wrap-and-shortest-turn)
- [Smooth rotation](#smooth-rotation)
- [Quaternions (QUAT)](#quaternions-quat)
- [Full example — turret snap and smooth](#full-example--turret-snap-and-smooth)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **Default** | **Radians** for `SIN`, `COS`, `ATAN2`, `YAWFROMXZ`, `VEC2.ROTATE` |
| **Degrees** | `SIND`, `COSD`, `ANGLEDIFF`, `DEGPERSEC`, some sprite APIs |
| **3D facing** | Yaw on XZ + optional pitch for aim |
| **Advanced** | `QUAT.*` for compound 3D rotations |

---

## Radians vs degrees

| Task | Use |
|------|-----|
| `SIN` / `COS` / `ATAN2` | Radians |
| Human-readable turn speed “90°/s” | `MATH.DEGPERSEC(90, APP.DELTA())` |
| Sprite rotation in degrees | `MATH.RAD2DEG(angle)` before `SETROTATION` |
| Import from degrees | `MATH.DEG2RAD(deg)` |

Constants: `MATH.PI()`, `MATH.TAU()`.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| 2D aim on screen | `ATAN2(dy, dx)` | [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md) |
| 3D yaw on ground | `YAWFROMXZ(dx, dz)` | [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md) |
| Shortest turn (deg) | `MATH.ANGLEDIFF(from, to)` | `to - from` |
| Shortest turn (rad) | `MATH.ANGLEDIFFRAD` | Manual wrap |
| Smooth spin | `MATH.LERPANGLE` | `LERP` on angles |
| Keep angle in range | `WRAPANGLE` / `WRAPANGLE180` | Manual modulo |
| 3D arbitrary axis | `QUAT.*` | Euler x+y+z every frame |

---

## Wrap and shortest turn

```basic
; How far to rotate this frame (degrees)
delta = MATH.ANGLEDIFF(currentDeg, targetDeg)
currentDeg = currentDeg + MATH.CLAMP(delta, -maxStep, maxStep)

; Radians shortest difference
dr = MATH.ANGLEDIFFRAD(currentRad, targetRad)
```

**Why `ANGLEDIFF`:** `350°` to `10°` is `+20°`, not `-340°`.

---

## Smooth rotation

```basic
currentRad = MATH.LERPANGLE(currentRad, targetRad, 8 * APP.DELTA())
```

Combine with entity:

```basic
yaw = MATH.YAWFROMXZ(dx, dz)
ENTITY.SETYAW(hero, yaw)
```

---

## Quaternions (QUAT)

For full 3D orientation (not just yaw), use `QUAT` namespace — see [reference/VEC_QUAT.md](../../../reference/VEC_QUAT.md) and [API_CONSISTENCY.md](../../../API_CONSISTENCY.md#quat).

| Pattern | API |
|---------|-----|
| From Euler | `QUAT.FROMEULER` / similar |
| Rotate vector | `VEC3.ROTATEBYQUAT(v, q)` |
| Multiply orientations | `QUAT.MUL` |
| Free handle | `QUAT.FREE` |

**When to use quats:** Turret with pitch + yaw, camera free look, physics alignment to surface normal.

---

## Full example — turret snap and smooth

```basic
APP.OPEN(400, 300, "Angles")
APP.SETFPS(60)

px = 200
py = 150
angle = 0

WHILE NOT APP.SHOULDCLOSE()
    mx = INPUT.MOUSEX()
    my = INPUT.MOUSEY()
    target = MATH.ATAN2(my - py, mx - px)

    IF INPUT.KEYDOWN(KEY_LSHIFT) THEN
        angle = MATH.LERPANGLE(angle, target, 10 * APP.DELTA())
    ELSE
        angle = target
    ENDIF

    DRAW.LINE(px, py, px + MATH.COS(angle) * 60, py + MATH.SIN(angle) * 60, 255, 200, 80, 255)
    DRAW.CIRCLE(px, py, 8, 100, 180, 255, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Mix deg and rad in one formula | Convert at boundaries only |
| `LERP` on angles | Use `LERPANGLE` |
| Gimbal lock on 3D | Quats or clamp pitch |
| Sprite upside down | Add π/2 offset to aim |
| `ATAN(y/x)` quadrant bugs | `ATAN2(y, x)` |

---

## See also

- [VEC3-MATH.md](VEC3-MATH.md) — `ROTATEBYQUAT`
- [CAMERA-AND-INPUT.md](../CAMERA-AND-INPUT.md) — FPS look
- [INTERPOLATION-AND-EASING.md](INTERPOLATION-AND-EASING.md) — `LERPANGLE` context
