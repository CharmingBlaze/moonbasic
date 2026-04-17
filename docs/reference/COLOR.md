# Color Commands

Color handle creation, component access, conversion (RGB, HSV, hex), and manipulation.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a color with `COLOR.RGB`, `COLOR.RGBA`, `COLOR.HEX`, or `COLOR.HSV`.
2. Read components with `COLOR.R` / `COLOR.G` / `COLOR.B` / `COLOR.A`.
3. Transform with `COLOR.LERP`, `COLOR.FADE`, `COLOR.INVERT`, `COLOR.CONTRAST`, etc.
4. Free with `COLOR.FREE` when done.

---

### `COLOR.RGB(r, g, b)` 

Creates a color handle from red, green, and blue values (0–255).

---

### `COLOR.RGBA(r, g, b, a)` 

Creates a color handle from red, green, blue, and alpha values (0–255).

---

### `COLOR.HEX(hexString)` 

Creates a color from a hexadecimal string (e.g. `"#FF0000"` or `"FF0000"`).

---

### `COLOR.HSV(h, s, v)` 

Creates a color from hue (0–360), saturation (0.0–1.0), and value (0.0–1.0).

---

### `COLOR.FROMHSV(h, s, v)` 

Alias for `COLOR.HSV`.

---

### `COLOR.CLAMP(r, g, b)` 

Clamps RGB float values to valid range. Returns a color handle.

---

### `COLOR.R(colorHandle)` 

Returns the red component (0–255).

---

### `COLOR.G(colorHandle)` 

Returns the green component (0–255).

---

### `COLOR.B(colorHandle)` 

Returns the blue component (0–255).

---

### `COLOR.A(colorHandle)` 

Returns the alpha component (0–255).

---

### `COLOR.LERP(a, b, t)` 

Returns a new color interpolated between colors `a` and `b` by factor `t` (0.0–1.0).

---

### `COLOR.FADE(colorHandle, alpha)` 

Returns a new color with the specified alpha transparency (0.0–1.0).

---

### `COLOR.TOHSVX(colorHandle)` 

Returns the hue component (H) of the color as a float.

---

### `COLOR.TOHSVY(colorHandle)` 

Returns the saturation component (S) of the color as a float.

---

### `COLOR.TOHSVZ(colorHandle)` 

Returns the value component (V) of the color as a float.

---

### `COLOR.TOHSV(colorHandle)` 

Returns the HSV representation as a handle.

---

### `COLOR.TOHEX(colorHandle)` 

Returns the color as a hex string (e.g. `"FF0000FF"`).

---

### `COLOR.INVERT(colorHandle)` 

Returns a new color with inverted RGB channels.

---

### `COLOR.CONTRAST(colorHandle, factor)` 

Returns a new color with contrast adjusted by `factor`.

---

### `COLOR.BRIGHTNESS(colorHandle, factor)` 

Returns a new color with brightness adjusted by `factor`.

---

### `COLOR.FREE(colorHandle)` 

Releases the color handle from memory.

---

## Full Example

This example creates a procedural palette and applies colors to entities.

```basic
FOR i = 0 TO 9
    hue = FLOAT(i) / 10.0 * 360.0
    c = COLOR.HSV(hue, 0.8, 1.0)
    ENTITY.SETCOLOR(enemies(i), COLOR.R(c), COLOR.G(c), COLOR.B(c), 255)

    ; Also get the hex representation
    hex = COLOR.TOHEX(c)
    PRINT "Enemy " + STR(i) + ": " + hex

    COLOR.FREE(c)
NEXT
```
