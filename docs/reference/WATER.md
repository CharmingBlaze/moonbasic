# Water (`WATER.*`)

A **horizontal water plane** with simple wave motion, color gradients, and queries for camera/gameplay (**depth**, **underwater**). **CGO + Raylib** required.

**Draw order:** After opaque terrain and props, before transparent weather/particles when possible.

---

### `Water.Make(width, length)`
Creates a subdivided water plane. Returns a **handle**.

### `Water.Free(handle)`
Frees the water resource.

---

### `Water.SetPos(handle, x, y, z)`
Sets the water surface world position.

### `Water.Draw(handle)`
Renders the water surface. Must be inside `Camera.Begin()`.

---

### `Water.GetWaveY(handle, x, z)`
Returns surface Y including wave offset at XZ.

### `Water.GetDepth(handle, x, z)`
Returns depth from surface to bed at XZ.

### `Water.IsUnder(handle, x, y, z)`
Returns `TRUE` if the point is below the animated surface.

---

## `Water.SetShallowColor` / `Water.SetDeepColor(water, r, g, b, a)`

Tint multipliers or colors for shallow vs deep regions (see runtime for exact blending).

---

## Common mistake

Using **`Terrain.GetHeight`** for water level — water has its **own** Y from **`SetPos`**; compare **`GetWaveY`** or **`IsUnder`** for consistency.

---

## See also

- [TERRAIN.md](TERRAIN.md)
- [SKY.md](SKY.md) — horizon tint
