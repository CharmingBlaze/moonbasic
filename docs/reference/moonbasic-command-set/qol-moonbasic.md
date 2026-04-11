# MoonBASIC-only quality of life

These are **real** engine helpers — not legacy Blitz/DBPro names. Prefer them for small games.

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **InputAxis(neg, pos)** | **`Input.Axis()`** | Returns -1, 0, or 1. |
| **InputOrbit(cam, target, dist, sens)** | **`Camera.SetOrbit()`** | |
| **MoveStepX(yaw, f, s, spd, dt)** | **`MoveStepX()`** | |
| **MoveStepZ(yaw, f, s, spd, dt)** | **`MoveStepZ()`** | |
| **LandBoxes(id, boxes)** | **`LandBoxes()`** | Snap to platform array. |
| **LandHeightmap(id, terrain)** | **`Terrain.SnapY()`** | |
| **CameraOrbit(cam, entity, dist)** | **`Camera.Orbit()`** | Auto-follow entity. |
| **LoadScene(path)** | **`Entity.LoadScene()`** | Clears then loads. |
| **SaveScene(path)** | **`Entity.SaveScene()`** | |
| **Entity.X(id)** | **`Entity.X()`** | |
| **Entity.Y(id)** | **`Entity.Y()`** | |
| **Entity.Z(id)** | **`Entity.Z()`** | |
| **Entity.Pitch(id)** | **`Entity.Pitch()`** | |
| **Entity.Yaw(id)** | **`Entity.Yaw()`** | |
| **Entity.Roll(id)** | **`Entity.Roll()`** | |
| **Entity.Distance(a, b)** | **`Entity.Distance()`** | |
| **UpdatePhysics()** | **`Entity.Update()`** | |
| **DrawEntities()** | **`DrawEntities()`** | Renders all visible. |

See [QOL.md](../QOL.md) for **`SCREENW`**, **`DT`**, **`ENDGAME`**, etc.
