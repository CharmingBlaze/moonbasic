# Navigation mesh and terrain (integration)

Open-world **heightfields** are **not** automatically converted into the grid nav system. Use existing **`NAV.*`** / **`PATH.*`** / **`NAVAGENT.*`** from [`runtime/mbnav`](../../runtime/mbnav) as documented in [NAV_AI.md](NAV_AI.md).

---

## Terrain workflow (grid nav)

1. **`Nav.Make()`** then **`Nav.SetGrid(...)`** to match your world scale.
2. Optionally **`Nav.AddTerrain(nav, modelHandle)`** — uses a **model’s axis-aligned bounding box** to mark walkable XZ cells and ground height (see [NAV_AI.md](NAV_AI.md)).
3. For pure heightfield worlds without a proxy model, build walkability from your own logic (e.g. sample **`Terrain.GetHeight`**, mark steep slopes blocked in a future helper) or use **`Nav.AddObstacle`** for blockers.
4. **`Nav.Build(nav)`** before **`Nav.FindPath`** or **`NavAgent.MoveTo`**.

---

## No duplicate APIs

Do **not** introduce a second **`NAV.FINDPATH`** or conflicting **`PATH.*`** names — extend **`mbnav`** or add a **new prefix** (e.g. `NAVTERRAIN.*`) if a dedicated bake pipeline is added later.

---

## See also

- [NAV_AI.md](NAV_AI.md) — full API tables
- [WORLD.md](WORLD.md) — streaming vs static nav rebuilds (rebuild when chunks change if you use baked meshes)
