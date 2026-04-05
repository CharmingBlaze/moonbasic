# Scatter and props (`SCATTER.*`, `PROP.*`)

**Scatter** places many decorative instances (implementation uses simple **draw** primitives for markers). **Prop** stores placed objects with optional model handles for future mesh draw. **CGO** required.

---

## `SCATTER.*`

| Command | Role |
|--------|------|
| `Scatter.Create(name$)` | New scatter set (`TagScatterSet`); **`name$`** is a label for debugging. |
| `Scatter.Free(scatter)` | Frees the set. |
| `Scatter.Apply(scatter, terrain, density#, scale#, seed)` | Populates instances over terrain (algorithm is implementation-defined). |
| `Scatter.DrawAll(scatter)` | Draws all instances. |

---

## `PROP.*`

| Command | Role |
|--------|------|
| `Prop.Place(model, x#, y#, z#)` | Records a prop at a point (model handle reserved for future use). |
| `Prop.Free(prop)` | Frees one prop handle. |
| `Prop.DrawAll()` | Draws all registered props. |

---

## Common mistake

Calling **`Scatter.Apply`** with extreme **density** — costs scale with instance count; start low.

---

## See also

- [TERRAIN.md](TERRAIN.md) — height sampling for placement
- [MODEL.md](MODEL.md) — when mesh props are wired
