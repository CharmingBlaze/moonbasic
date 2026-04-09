# 3D Hop (`mario64`)

Small **third-person** demo in a **Blitz3D-style** spirit: walk on a plane, jump, land on boxes, **orbit the camera**, immediate-mode primitives (**`BOX`**, **`WIRECUBE`**, **`FLAT`**, … — see [BLITZ3D.md](../../docs/reference/BLITZ3D.md)). Several variants:

| File | What to notice |
|------|----------------|
| **`main_orbit_simple.mb`** | **Easiest read:** commented “map” at the top, **`CONST`** palette + world bounds, one floor + one box — **`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`**, **`MOVESTEPX`/`Z`**, **`LANDBOXES`**, **`Camera.SetOrbit`**, **`ERASE ALL`**. |
| **`main.mb`** | **Default pick:** same hop as before, **implicit types** (no `#` / `$` / `?` suffixes), **Draw3D** + **Camera** only — no skybox or entity graph. |
| **`main_entities.mb`** | **Engine-style:** **CreateCube** / **CreateSphere**, **COLLISIONS**, **EntityGrounded** (coyote), **EntityMoveCameraRelative**, **Camera.OrbitEntity**, **CopyEntity** platforms, **ENTITY.UPDATE**, **DrawEntities**, child **hat** on **player**. |
| **`main_v2.mb`** | **Recommended teaching path:** parallel arrays for platforms, but **`Input.Axis`**, **`MOVEX`/`MOVEZ`**, **`BOXTOPLAND`** float return, **`IIF$`**, and **one line** for orbit yaw (`Input.Axis(KEY_Q, KEY_E) * DEGPERSEC(...)`). Heavily commented. |
| **`main_v3.mb`** | Same logic with **`TYPE` / `DIM AS`** — one `Platform` array instead of nine arrays. Uses **`Input.Orbit`** and **`MOVESTEPX`/`MOVESTEPZ`**. Landing still uses **`BOXTOPLAND`** in a loop ( **`LANDBOXES`** needs parallel **`DIM`** arrays). |

## Run

**`main.mb`** (Draw3D path, full runtime + Raylib):

```bash
go run -tags fullruntime . --run examples/mario64/main.mb
```

**Blitz-style variants** (often need **CGO** and the same **fullruntime** build if you use **`--run`** from the repo root):

```bash
go run -tags fullruntime . --run examples/mario64/main_orbit_simple.mb
```

**Controls (`main.mb` / orbit samples):** **Q/E** yaw, **right-drag** yaw/pitch, **wheel** zoom, **WASD** move, **Space** jump, **Esc** quit.

## Docs to read

- **[BLITZ3D.md](../../docs/reference/BLITZ3D.md)** — BlitzBasic3D → moonBASIC map (**`KEYHIT`**, **`WIRECUBE`**, **`Camera.Orbit`**, entities, …).  
- **`main_orbit_simple.mb`** — orbit deltas **`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`** in [GAMEHELPERS.md](../../docs/reference/GAMEHELPERS.md); **`Camera.SetOrbit`** in [CAMERA.md](../../docs/reference/CAMERA.md); teardown **`ERASE ALL`** in [MEMORY.md](../../docs/MEMORY.md).
- **Orbit camera** — `Camera.OrbitAround` in [CAMERA.md](../../docs/reference/CAMERA.md) (third-person on XZ + fixed eye height).
- **Walk + orbit input** — `Input.Axis` in [INPUT.md](../../docs/reference/INPUT.md); pair **Q/E** with **`DEGPERSEC`** for degrees-per-second yaw.
- **Movement** — `MOVEX` / `MOVEZ` in [MATH.md](../../docs/reference/MATH.md).
- **Landing** — `BOXTOPLAND` / `LANDBOXES` in [GAMEHELPERS.md](../../docs/reference/GAMEHELPERS.md).
