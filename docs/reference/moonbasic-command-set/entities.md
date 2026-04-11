# Entities (core)

**Two styles:** integer **entity id** (`ENTITY.*`) and **heap** **`ENTITYREF`** from **`CUBE()`** / **`SPHERE()`** (dot methods). See [BLITZ3D.md](../BLITZ3D.md), [ENTITY.md](../ENTITY.md).

## Creation / lifetime

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **Entity.Load(path)** | **`Entity.Load()`** | Returns entity ID. |
| **Entity.CreateCube(size)** | **`Entity.CreateCube()`** | |
| **Entity.CreateSphere(radius)** | **`Entity.CreateSphere()`**, **`SPHERE()`** | |
| **Entity.CreatePlane(width, height)** | **`Entity.CreatePlane()`** | |
| **Entity.CreateMesh(path)** | **`Entity.CreateMesh()`**, **`Entity.LoadMesh()`** | **Raylib model/mesh** — must **`FREE`** — [MEMORY.md](../../MEMORY.md). |
| **Entity.Copy(id)** | **`Entity.Copy()`** | New id — both need lifecycle rules. |
| **Entity.Free(id)** | **`Entity.Free()`** | Unloads model/animations in safe order. |

## Transform

| Designed | Implementation | Arguments (typical) |
|----------|----------------|---------------------|
| **Entity.Position(id, x, y, z)** | **`Entity.Position()`**, **`Entity.SetPosition()`** | **`(entity, x, y, z [, global])`** |
| **Entity.Move(id, f, r, u)** | **`Entity.Move()`**, **`MoveEntity()`** | **`(entity, forward, right, up)`** — **local** move along facing |
| **Entity.Turn(id, p, y, r)** | **`Entity.Turn()`** | Delta angles |
| **Entity.Scale(id, x, y, z)** | **`Entity.Scale()`** | **`(entity, sx, sy, sz)`** |
| **Entity.Translate(id, dx, dy, dz)** | **`Entity.Translate()`**, **`TranslateEntity()`** | **`(entity, dx, dy, dz)`** — **world** delta |
| **Entity.TFormVector(x, y, z, srcEntity, dstEntity)** | **`Entity.TFormVector()`** | **`(x, y, z, srcEntity, dstEntity)`** → **3-float array handle** |
| **Entity.LookAt(id, x, y, z)** | **`Entity.LookAt()`** | |
| **Entity.AlignToVector(id, x, y, z)** | **`Entity.AlignToVector()`** | |

## Rule-based collision types (with `Entity.Update`)

| Designed | Implementation | Arguments / returns |
|----------|------------------|---------------------|
| **Entity.Type(id, typeId)** | **`Entity.Type()`** | **`(entity, typeId)`** — used as **`src`/`dst`** in **`COLLISIONS`** |
| **EntityHitsType(id, type)** | **`EntityHitsType()`** | **`(entity, type)`** → **`TRUE`/`FALSE`** if any hit matches **`type`** |
| **EntityCollided(id, type)** | **`EntityCollided()`** | **`(entity, type)`** → **other entity id** or **0** |

## Getters

**ENTITYX** … **ENTITYROLL**, **ENTITYSCALE** patterns — see **`ENTITY.ENTITYX`** … in registry.

## Visibility / state

**SHOW/HIDE**, **ALPHA**, **COLOR**, **FX**, **ORDER**, **PARENT** — **`Entity.Show()`**, **`Entity.Hide()`**, **`Entity.Alpha()`**, **`Entity.Color()`**, **`Entity.FX()`**, **`Entity.Order()`**, **`Entity.Parent()`**.
