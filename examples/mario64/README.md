# 3D Hop (`mario64`)

Small **third-person** demo: walk on a plane, jump, land on boxes, **orbit the camera** with **Q/E**. Same gameplay in three versions:

| File | What to notice |
|------|----------------|
| **`main.mb`** | Longer “first draft”: manual **sin/cos** basis vectors, manual WASD, manual landing **Y** after `BOXTOPLAND`, `IF`/`ELSE` HUD. |
| **`main_v2.mb`** | **Recommended teaching path:** parallel arrays for platforms, but **`Input.Axis`**, **`MOVEX`/`MOVEZ`**, **`BOXTOPLAND`** float return, **`IIF$`**, and **one line** for orbit yaw (`Input.Axis(KEY_Q, KEY_E) * DEGPERSEC(...)`). Heavily commented. |
| **`main_v3.mb`** | Same logic with **`TYPE` / `DIM AS`** — one `Platform` array instead of nine arrays. Uses **`Input.Orbit`** and **`MOVESTEPX`/`MOVESTEPZ`**. Landing still uses **`BOXTOPLAND`** in a loop ( **`LANDBOXES`** needs parallel **`DIM`** arrays). |

## Run

```bash
CGO_ENABLED=1 go run . examples/mario64/main_v2.mb
```

**Controls:** **mouse** moves the orbit camera (horizontal); **WASD** move; **Q/E** still nudge orbit; **Space** jump; **Esc** quit.

## Docs to read

- **Orbit camera** — `Camera.OrbitAround` in [CAMERA.md](../../docs/reference/CAMERA.md) (third-person on XZ + fixed eye height).
- **Walk + orbit input** — `Input.Axis` in [INPUT.md](../../docs/reference/INPUT.md); pair **Q/E** with **`DEGPERSEC`** for degrees-per-second yaw.
- **Movement** — `MOVEX` / `MOVEZ` in [MATH.md](../../docs/reference/MATH.md).
- **Landing** — `BOXTOPLAND` in [GAMEHELPERS.md](../../docs/reference/GAMEHELPERS.md).
