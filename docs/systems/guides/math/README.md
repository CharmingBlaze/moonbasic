# Math & vector guides

> Deep **how/why** guides for gameplay math — split by 2D vs 3D and scalars vs vectors.

**Runnable copies:** [examples/guides/README.md](../../../examples/guides/README.md) (`examples/guides/math/*.mb`)

**Quick hub:** [../GUIDES.md](../GUIDES.md) · **All overloads:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **Overview:** [09-DATA.md](../../09-DATA.md)

**Intro (short):** [../MATH-AND-VECTORS.md](../MATH-AND-VECTORS.md)

---

## Pick a guide

| Guide | Use when… |
|-------|-----------|
| [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md) | Screen positions, top-down/side distances, aim angles on X/Y |
| [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md) | Ground-plane distance, yaw, XZ movement, 3D facing |
| [VEC2-MATH.md](VEC2-MATH.md) | 2D vector handles: normalize, rotate, pushout, move_toward |
| [VEC3-MATH.md](VEC3-MATH.md) | 3D vectors: dot, cross, reflect, project, quat rotate |
| [INTERPOLATION-AND-EASING.md](INTERPOLATION-AND-EASING.md) | Lerp, smoothstep, approach, remap, ping-pong |
| [ANGLES-AND-ROTATION.md](ANGLES-AND-ROTATION.md) | Wrap angles, lerp rotation, deg/rad, quaternions |
| [RANDOMNESS-AND-PROCEDURE.md](RANDOMNESS-AND-PROCEDURE.md) | Dice, loot tables, seeds, chance checks |

---

## Scalar vs vector — rule of thumb

```
Only x,y floats, no direction algebra?  → MATH-2D or MATH-3D gameplay
Need unit direction, push-out, lerp point? → VEC2 or VEC3
Compare distance often?                 → DISTSQ / HDISTSQ (skip sqrt)
Physics collision shapes?               → COLLISION guides (use VEC2 handles)
```

**Angles:** Engine math uses **radians** unless a command says degrees (`MATH.SIND`, `MATH.ANGLEDIFF`).

---

## Reference (exhaustive)

- [reference/MATH.md](../../../reference/MATH.md)
- [reference/VEC2.md](../../../reference/VEC2.md)
- [reference/VEC3.md](../../../reference/VEC3.md)
- [reference/GAME_MATH_HELPERS.md](../../../reference/GAME_MATH_HELPERS.md)
- [reference/VEC_QUAT.md](../../../reference/VEC_QUAT.md)
