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

| Designed | Implementation | Arguments (typical) |
|----------|----------------|---------------------|
| **PositionEntity** | **`ENTITY.POSITIONENTITY`**, **`SETPOSITION`** | **`(entity#, x, y, z [, global])`** |
| **RotateEntity** | **`ENTITY.ROTATEENTITY`** | **`(entity#, pitch, yaw, roll)`** |
| **ScaleEntity** | **`ENTITY.SCALE`** | **`(entity#, sx, sy, sz)`** |
| **MoveEntity** | **`ENTITY.MOVE`**, **`MOVEENTITY`** | **`(entity#, forward, right, up)`** — **local** move along facing |
| **TranslateEntity** | **`ENTITY.TRANSLATE`**, **`ENTITY.TRANSLATEENTITY`** | **`(entity#, dx, dy, dz)`** — **world** delta |
| **TFormVector** | **`ENTITY.TFORMVECTOR`** | **`(x, y, z, srcEntity#, dstEntity#)`** → **3-float array handle** |
| **TurnEntity** | **`ENTITY.TURNENTITY`** | Delta angles |
| **PointEntity** | **`ENTITY.POINTENTITY`** | |
| **AlignToVector** | **`ENTITY.ALIGNTOVECTOR`** | |

## Rule-based collision types (with `ENTITY.UPDATE`)

| Designed | Implementation | Arguments / returns |
|----------|------------------|---------------------|
| **EntityType** | **`ENTITY.TYPE`** | **`(entity#, typeId#)`** — used as **`src`/`dst`** in **`COLLISIONS`** |
| **EntityHitsType** | (bool wrapper) | **`(entity#, type#)`** → **`TRUE`/`FALSE`** if any hit matches **`type#`** |
| **ENTITYCOLLIDED** | **`ENTITYCOLLIDED`** | **`(entity#, type#)`** → **other entity id** or **0** |

## Getters

**ENTITYX** … **ENTITYROLL**, **ENTITYSCALE** patterns — see **`ENTITY.ENTITYX`** … in registry.

## Visibility / state

**SHOW/HIDE**, **ALPHA**, **COLOR**, **FX**, **ORDER**, **PARENT** — **`ENTITY.SHOW`**, **`HIDE`**, **`ALPHA`**, **`COLOR`**, **`FX`**, **`ORDER`**, **`PARENT`**.
