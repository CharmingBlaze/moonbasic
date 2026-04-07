# Entities (core)

**Two styles:** integer **entity id** (`ENTITY.*`) and **heap** **`ENTITYREF`** from **`CUBE()`** / **`SPHERE()`** (dot methods). See [BLITZ3D.md](../BLITZ3D.md), [ENTITY.md](../ENTITY.md).

## Creation / lifetime

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateCube (parent)** | **`ENTITY.CREATECUBE`**, **`CUBE()`** | **FREE** via **`ENTITY.FREE`** or **`ENTITY.CLEARSCENE`**. |
| **CreateSphere** | **`ENTITY.CREATESPHERE`**, **`SPHERE()`** | |
| **CreatePlane** | **`ENTITY.CREATEPLANE`** | |
| **CreateMesh** | **`ENTITY.CREATEMESH`**, **`ENTITY.LOADMESH`** | **Raylib model/mesh** — must **`FREE`** — [MEMORY.md](../../MEMORY.md). |
| **CreateSprite3D** | ≈ **`SPRITE`** / billboard / model — project-specific | |
| **CopyEntity** | **`ENTITY.COPY`** | New id — both need lifecycle rules. |
| **FreeEntity** | **`ENTITY.FREE`** | Unloads model/animations in safe order. |

## Transform

| Designed | Implementation |
|----------|----------------|
| **PositionEntity** | **`ENTITY.POSITIONENTITY`**, **`SETPOSITION`** |
| **RotateEntity** | **`ENTITY.ROTATEENTITY`** |
| **ScaleEntity** | **`ENTITY.SCALE`** |
| **MoveEntity** | **`ENTITY.MOVE`**, **`MOVEENTITY`** |
| **TurnEntity** | **`ENTITY.TURNENTITY`** |
| **PointEntity** | **`ENTITY.POINTENTITY`** |
| **AlignToVector** | **`ENTITY.ALIGNTOVECTOR`** |

## Getters

**ENTITYX** … **ENTITYROLL**, **ENTITYSCALE** patterns — see **`ENTITY.ENTITYX`** … in registry.

## Visibility / state

**SHOW/HIDE**, **ALPHA**, **COLOR**, **FX**, **ORDER**, **PARENT** — **`ENTITY.SHOW`**, **`HIDE`**, **`ALPHA`**, **`COLOR`**, **`FX`**, **`ORDER`**, **`PARENT`**.
