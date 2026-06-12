# Guide examples — runnable copies

Each script matches the **Full example** in `docs/systems/guides/` (and `docs/systems/guides/math/`).

| Script | Guide |
|--------|--------|
| `game_loop.mb` | [GAME-LOOP-AND-RENDERING.md](../../docs/systems/guides/GAME-LOOP-AND-RENDERING.md) |
| `lighting.mb` | [LIGHTING.md](../../docs/systems/guides/LIGHTING.md) |
| `particles.mb` | [PARTICLES.md](../../docs/systems/guides/PARTICLES.md) |
| `meshes_materials.mb` | [MESHES-MODELS-MATERIALS.md](../../docs/systems/guides/MESHES-MODELS-MATERIALS.md) |
| `files_json.mb` | [FILES-AND-JSON.md](../../docs/systems/guides/FILES-AND-JSON.md) |
| `save_progress.mb` | [SAVE-AND-PROGRESS.md](../../docs/systems/guides/SAVE-AND-PROGRESS.md) |
| `math/math_2d_chase.mb` | [MATH-2D-GAMEPLAY.md](../../docs/systems/guides/math/MATH-2D-GAMEPLAY.md) |
| `math/math_3d_xz.mb` | [MATH-3D-GAMEPLAY.md](../../docs/systems/guides/math/MATH-3D-GAMEPLAY.md) |
| `math/vec2_arena.mb` | [VEC2-MATH.md](../../docs/systems/guides/math/VEC2-MATH.md) |
| `math/vec3_aim.mb` | [VEC3-MATH.md](../../docs/systems/guides/math/VEC3-MATH.md) |
| `math/interpolation_ui.mb` | [INTERPOLATION-AND-EASING.md](../../docs/systems/guides/math/INTERPOLATION-AND-EASING.md) |
| `math/angles_turret.mb` | [ANGLES-AND-ROTATION.md](../../docs/systems/guides/math/ANGLES-AND-ROTATION.md) |
| `math/random_loot.mb` | [RANDOMNESS-AND-PROCEDURE.md](../../docs/systems/guides/math/RANDOMNESS-AND-PROCEDURE.md) |

**Check all guide examples:**

```bash
moonbasic --check examples/guides/game_loop.mb
```

From repo root, check every file:

```bash
# PowerShell
Get-ChildItem -Recurse examples/guides -Filter *.mb | ForEach-Object { go run . --check $_.FullName }
```

**Run (full runtime):**

```bash
moonrun examples/guides/game_loop.mb
moonrun examples/guides/math/math_2d_chase.mb
```

Uses current checklist APIs: `APP.*`, `RENDER.*`, `moonbasic --check`, `moonrun`.
