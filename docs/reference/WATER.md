# Water (`WATER.*`)

A **horizontal water plane** with simple wave motion, color gradients, and queries for camera/gameplay (**depth**, **underwater**). **CGO + Raylib** required.

**Draw order:** After opaque terrain and props, before transparent weather/particles when possible.

---

## `Water.Make(width#, length#)` → handle

Creates a subdivided plane mesh for rendering. Size is in world units.

---

## `Water.Free(water)`

Frees the water object.

---

## `Water.SetPos(water, x#, y#, z#)`

Places the water plane ( **`y`** is the nominal surface **`BedY`** reference for depth queries).

---

## `Water.Draw(water)` / `Water.Update(dt#)`

**Draw** renders the plane with animated normals/colors. **`Water.Update`** advances wave phase for **all** active water instances (module-level tick).

---

## `Water.SetWaveHeight(water, amp#)`

Sets vertical wave amplitude for rendering and **`GetWaveY`**.

---

## `Water.GetWaveY(water, x#, z#)` → float

Surface Y including wave offset at XZ (approximate).

---

## `Water.GetDepth(water, x#, z#)` → float

Returns a **column depth** metric from the animated surface to the bed plane at **XZ** (see runtime — not a ray through arbitrary **Y**). For point-in-water tests use **`Water.IsUnder`**.

---

## `Water.IsUnder(water, x#, y#, z#)` → bool

**True** if the point is below the animated surface.

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
