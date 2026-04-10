# Terrain chase

Procedural **heightfield** (`Terrain.Make` + `FillPerlin`), **orbital camera** (hold **RMB**), **WASD** movement (camera-relative via `INPUT.MOVEDIR`) with **Y** from `Terrain.SnapY`, and **eight spheres** that slowly seek your **XZ** position.

No textures or downloads — only built-in drawing.

## Run

From the repository root (full runtime + CGO):

```bash
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/terrain_chase/main.mb
```

## Docs

- [LESS_MATH.md](../../docs/reference/LESS_MATH.md) — patterns this sample uses (`MATH.CIRCLEPOINT`, `VEC2.PUSHOUT`, `INPUT.MOVEDIR`, `Terrain.Place`, …)
- [TERRAIN.md](../../docs/reference/TERRAIN.md) — `Terrain.Place`, streaming
- [WORLD.md](../../docs/reference/WORLD.md)
- [ENTITY.md](../../docs/reference/ENTITY.md) — `Entity.GetPos`, `FREEENTITIES`, `Camera.OrbitEntity`
- [COLOR.md](../../docs/reference/COLOR.md) — `COLOR.HSV(index, total)`
- [TIME.md](../../docs/reference/TIME.md) — `Time.Delta(min, max)`
- [DRAW2D.md](../../docs/reference/DRAW2D.md) — `DEBUG.PRINT`
- [INPUT.md](../../docs/reference/INPUT.md) — `Input.MouseWheel`

## Notes

- Terrain is **not** a Jolt collider here; props use **`Terrain.Place`** or **`Terrain.SnapY`** so Y follows the heightfield (see [ENTITY.md](../../docs/reference/ENTITY.md) *Terrain vs entity*).
- You can swap in textures or models later (`LoadMesh`, `EntityPBR`, etc.); this sample stays dependency-free for CI and offline use.
- Convenience built-ins: `CLAMP`, `VEC2.MOVE_TOWARD`, `INPUT.MOUSEDELTA`, tuple destructuring (`a, b = expr`).
