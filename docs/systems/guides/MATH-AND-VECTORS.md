# Math and vectors — hub

> Clamp health, roll dice, aim in 2D/3D, and move with vectors — pick the **focused guide** for your problem.

**Namespaces:** `MATH` · `VEC2` · `VEC3` · `QUAT` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../COMMAND_REGISTRY.md#data) · **Overview:** [09-DATA.md](../09-DATA.md)

---

## Math guide library

**Full index:** [math/README.md](math/README.md)

| Guide | When to read |
|-------|----------------|
| [math/MATH-2D-GAMEPLAY.md](math/MATH-2D-GAMEPLAY.md) | Screen X/Y, `DIST2D`, `ATAN2`, 2D aim |
| [math/MATH-3D-GAMEPLAY.md](math/MATH-3D-GAMEPLAY.md) | XZ ground, `HDIST`, `YAWFROMXZ`, `NEWX`/`NEWZ` |
| [math/VEC2-MATH.md](math/VEC2-MATH.md) | 2D normalize, rotate, pushout, `MOVE_TOWARD` |
| [math/VEC3-MATH.md](math/VEC3-MATH.md) | 3D dot, cross, reflect, distance |
| [math/INTERPOLATION-AND-EASING.md](math/INTERPOLATION-AND-EASING.md) | Lerp, smoothstep, approach, remap |
| [math/ANGLES-AND-ROTATION.md](math/ANGLES-AND-ROTATION.md) | Wrap angles, `LERPANGLE`, quaternions |
| [math/RANDOMNESS-AND-PROCEDURE.md](math/RANDOMNESS-AND-PROCEDURE.md) | `RAND`, `CHANCE`, seeds, loot tables |

---

## Quick picker

```
2D screen / tile distances?     → MATH-2D-GAMEPLAY
3D floor distance / yaw?        → MATH-3D-GAMEPLAY
2D direction handles?           → VEC2-MATH
3D direction / dot / cross?     → VEC3-MATH
Smooth bars / camera lag?       → INTERPOLATION-AND-EASING
Spin sprites / entity yaw?      → ANGLES-AND-ROTATION
Dice / loot / scatter spawn?    → RANDOMNESS-AND-PROCEDURE
```

**Rule:** Compare squared distance (`DISTSQ2D`, `HDISTSQ`, `VEC3.DISTSQ`) when you only need inside/outside radius.

Always multiply motion by **`APP.DELTA()`**.

---

## Minimal example

**Runnable:** [examples/guides/math/math_2d_chase.mb](../../../examples/guides/math/math_2d_chase.mb)

```basic
; Check: moonbasic --check examples/guides/math/math_2d_chase.mb
; Run:   moonrun examples/guides/math/math_2d_chase.mb
dist = MATH.DIST2D(px, py, tx, ty)
IF dist > 2 THEN
    t = MATH.CLAMP(speed * APP.DELTA() / dist, 0, 1)
    px = MATH.LERP(px, tx, t)
    py = MATH.LERP(py, ty, t)
ENDIF
```

---

## Reference (exhaustive)

- [reference/MATH.md](../../reference/MATH.md)
- [reference/VEC2.md](../../reference/VEC2.md)
- [reference/VEC3.md](../../reference/VEC3.md)
- [reference/GAME_MATH_HELPERS.md](../../reference/GAME_MATH_HELPERS.md)

---

## See also

- [CAMERA-AND-INPUT.md](CAMERA-AND-INPUT.md) — aim and movedir
- [COLLISION-2D.md](COLLISION-2D.md) — `VEC2` in overlap tests
- [CHARACTER-3D-WALKING.md](CHARACTER-3D-WALKING.md) — slope movement
