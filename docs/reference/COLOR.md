# Color commands

Color APIs use a color handle plus conversion/component helpers.

### `Color.RGB(r, g, b)`
Creates a color handle from red, green, and blue values (0-255).

### `Color.RGBA(r, g, b, a)`
Creates a color handle from red, green, blue, and alpha values (0-255).

### `Color.Hex(hexString)`
Creates a color from a hexadecimal string (e.g., `"#FF0000"` or `"FF0000"`).

---

### `Color.R(handle)` / `Color.G(handle)` / `Color.B(handle)` / `Color.A(handle)`
Returns the individual red, green, blue, or alpha component (0-255).

### `Color.Set(handle, r, g, b, a)`
Updates the components of an existing color handle in place.

---

### `Color.Lerp(a, b, t)`
Returns a new color handle interpolated between colors `a` and `b` by factor `t`.

### `Color.Fade(handle, alpha)`
Returns a new color handle with the specified alpha transparency (0.0-1.0).

---

### `Color.Free(handle)`
Releases the color handle from the heap and frees its memory.

```basic
; Procedural palette (normalized hue)
c = COLOR.FROMHSV(FLOAT(i) / FLOAT(N_ENEMY), 0.8, 1.0)
h, s, v = COLOR.TOHSV(c)   ; optional: read back HSV as a triple
EntityColor(enemy(i), COLOR.R(c), COLOR.G(c), COLOR.B(c), 255)
COLOR.FREE(c)

; Clamp precomputed RGB values
cr, cg, cb = COLOR.CLAMP(cr, cg, cb)
EntityColor(enemy(i), cr, cg, cb, 255)
```
