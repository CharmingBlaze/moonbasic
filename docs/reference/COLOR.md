# Color Commands

Color handle creation, component access, conversion (RGB, HSV, hex), and manipulation.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a color with `COLOR.RGB`, `COLOR.RGBA`, `COLOR.HEX`, or `COLOR.HSV`.
2. Read components with `COLOR.R` / `COLOR.G` / `COLOR.B` / `COLOR.A`.
3. Transform with `COLOR.LERP`, `COLOR.FADE`, `COLOR.INVERT`, `COLOR.CONTRAST`, etc.
4. Free with `COLOR.FREE` when done.

---

### `COLOR.RGB(r, g, b)` / `RGBA`
Creates a color handle from component values (0–255).

- **Returns**: (Handle) The new color handle.
- **Example**:
    ```basic
    red = COLOR.RGB(255, 0, 0)
    ```

---

### `COLOR.HEX(hexString)`
Creates a color from a hexadecimal string.

- **Returns**: (Handle) The new color handle.

---

### `COLOR.HSV(h, s, v)`
Creates a color from hue (0–360), saturation (0.0–1.0), and value (0.0–1.0).

- **Returns**: (Handle) The new color handle.

---

### `COLOR.R(handle)` / `G` / `B` / `A`
Returns the specific component value (0–255).

- **Returns**: (Float)

---

### `COLOR.LERP(a, b, t)`
Returns a new color interpolated between two colors.

- **Returns**: (Handle) The new color handle.

---

### `COLOR.FADE(handle, alpha)`
Returns a new color with adjusted alpha transparency (0.0–1.0).

- **Returns**: (Handle) The new color handle.

---

### `COLOR.TOHEX(handle)`
Returns the color as a hex string.

- **Returns**: (String)

---

### `COLOR.FREE(handle)`
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

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `COLOR.FROMHSV(h, s, v)` | Create color from hue (0–360), saturation and value (0–1). |
| `COLOR.TOHSV(color)` | Returns `[h, s, v]` array. |
| `COLOR.TOHSVX(color)` / `TOHSVY` / `TOHSVZ` | Individual H / S / V component. |
| `COLOR.BRIGHTNESS(color, factor)` | Returns color brightened/darkened by `factor`. |
| `COLOR.CLAMP(color)` | Clamp all channels to 0–255 and return result. |

## See also

- [DRAW2D.md](DRAW2D.md) — `DRAW.RECTANGLE(x,y,w,h, r,g,b,a)`
- [IMAGE.md](IMAGE.md) — `IMAGE.COLORTINT`, `IMAGE.COLORREPLACE`
