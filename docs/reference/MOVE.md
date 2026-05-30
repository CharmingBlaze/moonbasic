# Move Commands

Scalar movement math aliases. Thin wrappers over `MATH.*` for convenience.

## Commands

### `MOVE.TOWARD(current, target, maxDelta)` 

Moves `current` toward `target` by at most `maxDelta`. Returns the new value. Alias of `MATH.APPROACH`.

```basic
speed = MOVE.TOWARD(speed, maxSpeed, accel * dt)
```

---

### `MOVE.LERP(a, b, t)` 

Returns `a + (b - a) * t`. Alias of `MATH.LERP`.

```basic
px = MOVE.LERP(px, targetX, 5 * dt)
```

---

## See also

- [MATH.md](MATH.md) — `MATH.APPROACH`, `MATH.LERP`, `MATH.SMOOTHSTEP`
- [VEC2.md](VEC2.md) — `VEC2.LERP`, `VEC2.MOVETOWARD`
- [VEC3.md](VEC3.md) — `VEC3.LERP`, `VEC3.MOVETOWARD`
