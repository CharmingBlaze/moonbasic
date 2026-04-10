# Color commands

Color APIs use a color handle plus conversion/component helpers.

## Constructors

- `COLOR.RGB(r, g, b) -> handle`
- `COLOR.RGBA(r, g, b, a) -> handle`
- `COLOR.HEX("#RRGGBB" | "#RRGGBBAA") -> handle`
- `COLOR.HSV(h, s, v) -> handle`
- `COLOR.HSV(index, total) -> handle` — **evenly spaced hues** on the wheel for **`1 .. total`** (uses **`index`** as a 1-based slot). Saturation and value are fixed (high saturation, full value) so you can color many objects without hand-written sine/RGB clamp blocks.
- `COLOR.FROMHSV(h, s, v) -> handle` (alias of `COLOR.HSV`)

`COLOR.FROMHSV` / the 3-argument `COLOR.HSV` accept hue in either `0..360` (degrees) or `0..1` (normalized convenience form).

**Handle lifetime:** `COLOR.HSV` returns a **heap color handle**. If you allocate every frame, call **`COLOR.FREE`** when done, or reuse one handle. **`EntityColor(entity, colorHandle)`** accepts a color handle as the second argument (see **`ENTITY.COLOR`** overloads in [API_CONSISTENCY.md](../API_CONSISTENCY.md)).

## Operations

- `COLOR.LERP(aColor, bColor, t) -> handle`
- `COLOR.FADE(color, alpha01) -> handle`
- `COLOR.INVERT(color) -> handle`
- `COLOR.CONTRAST(color, amount) -> handle`
- `COLOR.BRIGHTNESS(color, amount) -> handle`
- `COLOR.CLAMP(r, g, b) -> handle` (returns clamped tuple-like 3-float array for destructuring)

## Components and conversion

- `COLOR.R(color) -> int`
- `COLOR.G(color) -> int`
- `COLOR.B(color) -> int`
- `COLOR.A(color) -> int`
- `COLOR.TOHSV(color) -> handle` — **tuple** **`(h, s, v)`** in one call (three floats; destructuring `h, s, v = COLOR.TOHSV(c)`), same Raylib `ColorToHSV` components as the scalars below.
- `COLOR.TOHSVX(color) -> float`
- `COLOR.TOHSVY(color) -> float`
- `COLOR.TOHSVZ(color) -> float`
- `COLOR.TOHEX(color) -> string`
- `COLOR.FREE(color)`

```basic
; Procedural palette (normalized hue)
c = COLOR.FROMHSV(FLOAT(i) / FLOAT(N_ENEMY), 0.8, 1.0)
h#, s#, v# = COLOR.TOHSV(c)   ; optional: read back HSV as a triple
EntityColor(enemy(i), COLOR.R(c), COLOR.G(c), COLOR.B(c), 255)
COLOR.FREE(c)

; Clamp precomputed RGB values
cr, cg, cb = COLOR.CLAMP(cr, cg, cb)
EntityColor(enemy(i), cr, cg, cb, 255)
```
