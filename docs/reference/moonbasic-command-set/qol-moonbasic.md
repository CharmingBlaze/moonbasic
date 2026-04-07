# MoonBASIC-only quality of life

These are **real** engine helpers — not legacy Blitz/DBPro names. Prefer them for small games.

## Input

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **InputAxis (neg_key, pos_key)** | **`INPUT.AXIS`** / **`Input.Axis`** in [INPUT.md](../INPUT.md) | No heap. |
| **InputOrbit (cam, target, dist, sens)** | **`ORBITYAWDELTA`**, **`ORBITPITCHDELTA`**, **`ORBITDISTDELTA`** + **`Camera.SetOrbit`** | Floats only — [GAMEHELPERS.md](../GAMEHELPERS.md). |

## Movement

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **MoveStepX / Z / Y** | **`MOVESTEPX`**, **`MOVESTEPZ`** (math), **`MOVER.MOVESTEPX`**, **`MOVER.MOVESTEPZ`** (facade) | **MOVER** uses a **heap handle** — **`MOVER.FREE`**. No **`MoveStepY`** as a single global — use **`JUMP`** / physics / custom. |

## Collision / landing

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **LandBoxes (entity, boxes)** | **`LANDBOXES`**, **`LANDBOX`**, **`MOVER.LAND`** | **Platform array** is often a **heap handle** — own **`FREE`** for that array. |
| **LandHeightmap** | **`BOXTOPLAND`**, terrain helpers in [GAMEHELPERS.md](../GAMEHELPERS.md) | |

## Camera

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CameraOrbit** | **`CAMERA.SETORBIT`**, **`ORBIT`**, **`ORBITENTITY`** | |
| **CameraFollow** | **`CAMERA.FOLLOW`**, **`FOLLOWENTITY`** | |

## Scene

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **LoadScene (file$)** | **`ENTITY.LOADSCENE`**, **`SCENE.LOADSCENE`** | Clears entities then loads — **native** resources tracked per [MEMORY.md](../../MEMORY.md). |
| **SaveScene (file$)** | **`ENTITY.SAVESCENE`**, **`SCENE.SAVESCENE`** | |
| **SceneEntities ()** | **`ENTITY.ENTITIESINRADIUS`**, **`ENTITIESINBOX`**, **groups** | Returns **array handles** — **`FREE`** when done. |

See [QOL.md](../QOL.md) for **`SCREENW`**, **`DT`**, **`ENDGAME`**, etc.
