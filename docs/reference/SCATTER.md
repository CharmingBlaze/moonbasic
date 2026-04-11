# Scatter and props (`SCATTER.*`, `PROP.*`)

**Scatter** places many decorative instances (implementation uses simple **draw** primitives for markers). **Prop** stores placed objects with optional model handles for future mesh draw. **CGO** required.

---

## `SCATTER.*`

### `Scatter.Make()`
Creates a new scatter system for distributing decorative objects over terrain. Returns a **handle**.

### `Scatter.Free(handle)`
Frees the scatter system and its associated resources.

### `Scatter.Apply(handle, terrain)`
Snaps all scatter items to the surface of the specified terrain.

### `Scatter.Update(handle, dt)`
Updates scatter item animations or internal logic based on elapsed time.

### `Scatter.Draw(handle)`
Renders all scatter items. This must be called within a **`Camera.Begin()`** / **`Camera.End()`** block.

### `Scatter.SetDensity(handle, density)`
Sets the item density for the scatter system.

### `Scatter.SetArea(handle, x, z, w, d)`
Sets the world-space bounding area for scattering.

---

## `PROP.*`

### `Prop.Place(model, x, y, z)`
Records a prop at a world position.

### `Prop.Free(handle)`
Frees a specific prop handle.

### `Prop.DrawAll()`
Draws all registered props.

---

## Common mistake

Calling **`Scatter.Apply`** with extreme **density** — costs scale with instance count; start low.

---

## See also

- [TERRAIN.md](TERRAIN.md) — height sampling for placement
- [MODEL.md](MODEL.md) — when mesh props are wired
