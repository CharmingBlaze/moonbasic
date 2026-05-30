# Angle Commands

Angle utility helpers.

## Commands

### `ANGLE.DIFFERENCE(a, b)` 

Returns the shortest signed angle (in degrees) from `a` to `b`, wrapping correctly across the 0°/360° boundary. Result is in the range `-180` to `180`.

Alias of `MATH.ANGLEDIFF`.

```basic
diff = ANGLE.DIFFERENCE(350, 10)   ; returns 20 (not -340)
diff = ANGLE.DIFFERENCE(10, 350)   ; returns -20
```

---

## See also

- [MATH.md](MATH.md) — `MATH.ANGLEDIFF`, `MATH.WRAPANGLE`, `MATH.NORMALIZEANGLE`
