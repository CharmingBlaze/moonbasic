# Modern Blitz-Style Command Surface (MoonBasic)

This document lists **friendly global names** (Blitz3D-style) mapped to the **register-based VM** implementation. Names are **case-insensitive** at compile time and normalize to dotted uppercase for `NAMESPACE.METHOD` forms.

For a **full program** that combines `Graphics3D`, `LoadMesh`, `EntityPBR`, `RENDER.Begin3D`, `DrawEntities`, and the frame contract (`Render.Clear` / `Render.Frame`), see [GETTING_STARTED.md](../GETTING_STARTED.md) (**Modern Blitz-style 3D**) and [EXAMPLES.md](../EXAMPLES.md).

## Conventions

- **Entity** — integer entity id from `ENTITY.CREATE*` / `Create*` helpers.
- **camera** — heap handle from **`CAMERA.CREATE`** (canonical) or deprecated `CAMERA.MAKE` / `CreateCamera()`.
- **Parent `0`** — no parent (world root). Optional parent arguments attach via `ENTITY.PARENT` / `EntityParent`.
- **Dotted forms** — every friendly command has an equivalent `ENTITY.*` / `CAMERA.*` / `FOG.*` / `TERRAIN.*` name where applicable.

## Entity lifecycle and hierarchy

| Command | Arguments | Implementation |
|--------|-----------|------------------|
| `CreatePivot` | `[parent]` | Empty hidden entity; optional parent. |
| `CreateCamera` | `[parent]` | **`CAMERA.CREATE`** (thin wrapper); parent attachment not applied (camera is heap object). |
| `CreateCube` | `()` / `(parent)` / `(w,h,d)` / `(parent,w,h,d)` | Box primitive. |
| `CreateSphere` | (see `ENTITY.CREATESPHERE`) | Radius + segments; parent overloads same pattern as cube where registered. |
| `CreateCylinder` | `(radius, height, segments)` | Cylinder primitive. |
| `CreateCone` | `()` / `(parent)` / `(r,h,seg)` / `(parent,r,h,seg)` | Cone primitive (`entKindCone`). |
| `CreatePlane` | `(size)` | Plane primitive. |
| `CreateMesh` | `[parent]` | Procedural mesh builder (`ENTITY.CREATEMESH`). |
| `CreateTerrain` | — | **Not an entity**: use `TERRAIN.MAKE(w, h [, cellSize])` (returns terrain handle). |
| `CreateMirror` | — | **Not implemented** (planar reflection deferred). |
| `CreateSprite` | (see `ENTITY.CREATESPRITE`) | Billboard sprite entity. |
| `LoadMesh` | `(path [, parent])` | `ENTITY.LOADMESH` |
| `LoadAnimMesh` / `LoadAnimMesh` | — | `ENTITY.LOADANIMATEDMESH` |
| `CopyEntity` | `(entity [, parent])` | `ENTITY.COPY` + optional `ENTITY.PARENT` |
| `FreeEntity` | `(entity)` | `ENTITY.FREE` |
| `EntityParent` | `(entity, parent [, global])` | `ENTITY.PARENT` |
| `GetParent` | `(entity)` | Returns parent entity or 0 |
| `FindChild` | `(entity, name)` | `ENTITY.FINDCHILD` |
| `GetChild` | `(entity, index)` | `ENTITY.GETCHILD` |
| `CountChildren` | `(entity)` | `ENTITY.COUNTCHILDREN` |

## Transformation and state

| Command | Arguments | Notes |
|--------|-----------|------|
| `PositionEntity` | `(entity, x, y, z [, global])` | Same as `ENTITY.POSITIONENTITY` / **`ENTITY.SETPOS`** (canonical) / deprecated `ENTITY.SETPOSITION`. |
| `MoveEntity` | `(entity, forward, right, up)` | Same as `MOVEENTITY` / `ENTITY.MOVE` — **local** axes from pitch/yaw. |
| `TranslateEntity` | `(entity, dx, dy, dz)` | Same as `ENTITY.TRANSLATE` / `ENTITY.TRANSLATEENTITY` — **world** delta. |
| `RotateEntity`, `TurnEntity`, `ScaleEntity` | (see `ENTITY.ROTATEENTITY`, `ENTITY.TURNENTITY`, `ENTITY.SCALE`) | Absolute vs delta rotation; non-uniform scale. |
| `TFormVector` | `(x, y, z, srcEntity, dstEntity)` → handle | Alias of `ENTITY.TFORMVECTOR`; **3-float array** handle in `dst` space. |
| `PointEntity` | `(entity, targetEntity [, roll])` | `ENTITY.POINTENTITY` (roll optional not fully wired). |
| `AlignToVector` | `(entity, vx, vy, vz, axis)` | `ENTITY.ALIGNTOVECTOR` |
| `EntityX` … `EntityRoll` | `(entity [, global])` | Same as `ENTITY.ENTITYX` … |
| `EntityVisible` | `(entity, visible)` | `ENTITY.VISIBLE` |
| `EntityDistance` | `(a, b)` | `ENTITY.DISTANCE` |
| `EntityInView` | `(entity, camera)` | `ENTITY.INVIEW` |
| `EntityName` / `NameEntity` | `(entity)` / `(entity, name)` | Get / set name (`ENTITY.SETNAME`) |

## Physics and collisions

| Command | Arguments | Notes |
|--------|-----------|------|
| `EntityType` | `(entity, type)` | Sets **collision type** id used by `COLLISIONS` / `ENTITYCOLLIDED` / `EntityHitsType`. |
| `EntityRadius`, `EntityBox` | `(entity, r)` / `(entity, w, h, d)` | `ENTITY.RADIUS` / `ENTITY.BOX` |
| `EntityMass`, `EntityFriction`, `EntityRestitution` | … | `ENTITY.SETMASS`, `SETFRICTION`, `SETBOUNCE` |
| `Collisions` | `(srcType, dstType, method, response)` | `COLLISIONS` — register **rule-based** pairs (e.g. sphere–box `method` **2**). |
| `ApplyEntityImpulse` / `ApplyEntityForce` | `(entity, fx, fy, fz)` | `ENTITY.ADDFORCE` (adds to velocity, mass-weighted). |
| `ApplyEntityTorque` | — | **Stub** — no torque integrator |
| `EntityHitsType` | `(entity, type)` → **bool** | **`TRUE`** if `entity` hit any body with **`EntityType == type`** after last `ENTITY.UPDATE`. |
| `ENTITYCOLLIDED` | `(entity, type)` → **int** | Returns **other entity id** or **0** (same hit test as `EntityHitsType`). |
| `EntityCollided` | `(entity, entity)` → **bool** | **Jolt** pairwise contact (Linux + `ENTITY.LINKPHYSBUFFER`); not the same as `EntityHitsType`. |
| `CountCollisions` | (0 args) | Jolt contact count (`PhysicsContactCount` alias) |
| `COUNTCOLLISIONS` | `(entity)` | Rule-based hit **count** |
| `GETCOLLISIONENTITY` / `CollisionEntity` | `(entity, index)` | Hit entity at index |
| `ENTITY.COLLISIONX` … / `CollisionX` … | Last hit **or** `(entity, index)` for indexed rule hits |
| `PhysicsCollisionNX` … `PhysicsContactCount` | (0 args) | **Jolt** last-contact globals |

## Picking

| Command | Notes |
|--------|------|
| `LinePick` | Ray vs static entity AABBs; sets **Picked\*** |
| `CameraPick` | Screen ray vs static AABBs |
| `PickedX` … `PickedNZ`, `PickedEntity`, `PickedDistance` | Last pick result |
| `PickedSurface` / `PickedTriangle` | Return `-1` until mesh triangle index is tracked |
| `EntityPick` | `ENTITY.PICK` — forward pick along entity yaw |

## Rendering and post

| Command | Notes |
|--------|------|
| `EntityColor`, `EntityAlpha`, … | `ENTITY.COLOR`, `ENTITY.ALPHA`, … |
| `EntityPBR`, `EntityNormalMap`, `EntityEmission` | Modern PBR path (see `entity_modern_fx_cgo.go`) |
| `EntityAutoFade` | **Stub** |
| `SetMSAA`, `SetSSAO`, `SetBloom`, `SetPostProcess` | Window / `EFFECT.*` / `POST.ADDSHADER` |
| `CameraRange`, `CameraZoom`, `CameraFOV` | `CAMERA.SETRANGE`, `CAMERA.ZOOM`, `CAMERA.SETFOV` |
| `CameraFogMode`, `CameraFogRange`, `CameraFogColor` | Forward to `FOG.*` (global fog; camera arg accepted for API parity) |
| `CameraProject` | Alias of `CAMERA.WORLDTOSCREEN` (returns 2-element array handle) |

## Animation

Friendly names `Animate`, `SetAnimTime`, `EntityAnimTime`, `ExtractAnimSeq`, `AnimLength`, `FindBone` map to existing `ENTITY.*` builtins.

## Surfaces and mesh helpers

`CreateSurface`, `AddVertex`, `AddTriangle`, `UpdateMesh`, `VertexX/Y/Z` match `ENTITY.*`.  
`UpdateNormals`, `FlipMesh`, `FitMesh`, `VertexNX/Y/Z`, `VertexU/V`, `CountVertices`, `CountTriangles` are **stubs or partial** until the procedural mesh builder exposes normals/UV topology.

## Audio

`Load3DSound`, `EmitSound`, `SoundVolume`, `SoundPitch`, `Listener` — see `runtime/audio/spatial_cgo.go`.

## Brush and materials

| Command | Arguments | Notes |
|--------|-----------|------|
| `CreateBrush` | `()` or `(r, g, b)` | Default RGB white. |
| `LoadBrush` | `(path [, flags, uScale, vScale])` | Loads texture + white brush; owns texture until `FreeBrush`. |
| `FreeBrush` | `(brush)` | Frees brush; unloads embedded texture from `LoadBrush`. |
| `BrushColor` | `(brush, r, g, b)` | 0–1 or 0–255 channel convention (same as entity color). |
| `BrushAlpha` | `(brush, alpha)` | `>1` treated as 0–255 → normalized. |
| `BrushBlend` | `(brush, mode)` | `0/1` opaque/alpha → alpha blend; `2` multiply; `3` additive (Raylib blend). |
| `BrushTexture` | `(brush, texture [, frame, uvIndex])` | `frame` reserved; `uvIndex` → internal UV slot. Replaces texture; frees prior embedded tex from `LoadBrush` if any. |
| `BrushFX` / `BrushShininess` | | Existing builtins. |
| `PaintEntity` | | Copies brush color/FX/shininess/alpha and texture ref to entity. |
| `GetEntityBrush` | `(entity)` | Brush handle or `0`. |
| `PaintSurface` | `(surface, brush)` | Mesh-builder surface handle (`CreateSurface`). |
| `GetSurfaceBrush` | `(surface)` | Brush handle or `0`. |

## Textures (friendly names)

| Command | Notes |
|--------|------|
| `LoadTexture` | Alias of `TEXTURE.LOAD` / `LOADTEXTURE`; optional second `flags` (default `1` = trilinear + repeat). |
| `CreateTexture` | `(w, h [, flags])` — blank RGBA. |
| `LoadAnimTexture` | `(path, flags, cellW, cellH, firstFrame, frameCount)` — **one cell** from a horizontal strip (`frameCount` reserved for future multi-frame). |
| `TextureWidth` / `TextureHeight` | Same as `TEXTUREWIDTH` / `TEXTUREHEIGHT`. |
| `TextureName` | Source path when loaded from disk; else `""`. |
| `FreeTexture` | `TEXTURE.FREE` / `FREETEXTURE`. |
| `ScaleTexture`, `PositionTexture`, `RotateTexture` | Store UV metadata on the texture object (for materials that read these fields). |
| `TextureCoords` | `(texture, coords)` — integer coord mode tag. |
| `SetCubeFace`, `SetCubeMode` | Metadata for cubemap-style assets (no cubemap GPU path yet). |

## Terrain (friendly names)

| Command | Notes |
|--------|------|
| `TerrainHeight` | Alias of `TERRAIN.GETHEIGHT`. |
| `LoadTerrain` | `(path [, parent])` — greyscale image → heightfield; `parent` reserved. |
| `TerrainDetail` / `TerrainShading` | No-op placeholders (API parity). |
| `ModifyTerrain` | `(terrain, x, z, height [, realtime])` — sets nearest cell height; `realtime` reserved. |
| `TerrainX` / `TerrainZ` | World position → grid fractional index along X / Z (Y argument ignored). |
| `TerrainSize` | `(terrain)` → **2-element float array** `[cellsX, cellsY]`. |

Canonical chunk/streaming commands remain `TERRAIN.*` / `CHUNK.*` (see [TERRAIN.md](TERRAIN.md)).

## Input (flat globals)

| Command | Notes |
|--------|------|
| `MouseX`, `MouseY`, `MouseZ` | Aliases of `INPUT.MOUSE*` / wheel (`MOUSEZ`). |
| `MouseXSpeed`, `MouseYSpeed` | Delta (`INPUT.MOUSEXSPEED` / `YSPEED`). |
| `MouseDown`, `MouseHit` | Same as `INPUT.MOUSEDOWN` / `INPUT.MOUSEHIT`. |
| `FlushMouse`, `FlushKeys` | No-ops (Raylib is polled; documented for Blitz parity). |
| `WaitMouse`, `WaitKey` | Blocking poll until a mouse button / key press (returns button id / key code). |
| `MoveMouse` | `(x, y)` screen position (`INPUT.SETMOUSEPOS`). |
| `HidePointer` / `ShowPointer` | `CURSOR.HIDE` / `CURSOR.SHOW`. |
| `GetKey` | First key pressed this poll, or `0`. |
| `KeyHit`, `KeyDown`, `KEYHIT`, `KEYDOWN` | Existing. |

## Display and video memory

| Command | Notes |
|--------|------|
| `WindowWidth`, `WindowHeight` | Same as `WINDOW.WIDTH` / `HEIGHT`. |
| `ScreenWidth`, `ScreenHeight` | Same metrics (screen/window client size). |
| `GraphicsWidth`, `GraphicsHeight` | `RENDER.WIDTH` / `HEIGHT`. |
| `GraphicsDepth` | Returns `32` (placeholder; Raylib does not expose drawable depth bits per target). |
| `AvailVidMem`, `TotalVidMem` | `-1` = unknown (no portable VRAM query in core runtime). |
| `GpuName` | Alias of `SYSTEM.GPUNAME`. |

## File I/O (flat)

| Command | Notes |
|--------|------|
| `WriteFile` / `ReadFile` | `(path)` → open for write / read; returns **file handle** (same as `FILE.OPENWRITE` / `OPENREAD`). |
| `CloseFile` | `FILE.CLOSE`. |
| `WriteLine` / `ReadLine` | `FILE.WRITELN` / `READLINE`. |
| `WriteInt` / `ReadInt` | Little-endian 32-bit integers on the file stream. |
| `WriteFloat` / `ReadFloat` | Little-endian 32-bit float. |
| `EOF` | Already `EOF` → `FILE.EOF`. |

## System

| Command | Notes |
|--------|------|
| `MilliSecs` | `time` module — milliseconds since init |
| `Graphics3D` | `(w,h)` or `(w,h,depth,mode)` — resize; depth reserved; mode bit0 = high-DPI hint |
| `AppTitle` | `WINDOW.SETTITLE` |
| `UpdatePhysics` | `UPDATEPHYSICS` — `ENTITY.UPDATE(Time.Delta)` + best-effort `WORLD.UPDATE`, `PHYSICS2D.STEP`, `PHYSICS3D.STEP` (blitzengine) |
| `RENDER.BEGIN3D` / `RENDER.END3D` | Delegate to `CAMERA.BEGIN` / `CAMERA.END` |

---

## Math and trigonometry (Float64)

All math uses **float64** end-to-end. Friendly names live in the `math` namespace alongside `MATH.*` / flat `SIN` forms.

| Command | Arguments | Notes |
|--------|-----------|--------|
| `Sin` / `Cos` / `Tan` | `(angle)` | **Degrees** (same as `SIN` / `COS` / `TAN`). |
| `SINRAD` / `COSRAD` / `TANRAD` | `(radians)` | Radian trig (physics / legacy). |
| `ASin` / `ACos` / `ATan` | `(x)` | Inverse trig; **radians** (same as `ASIN` / `ACOS` / `ATAN`). |
| `ATan2` | `(y, x)` | `math.Atan2` (radians). |
| `Sqrt` / `Abs` | `(x)` | Same as `SQRT` / `ABS`. |
| `Floor` / `Ceil` | `(x)` | Same as `FLOOR` / `CEIL`. |
| `ROUND` / `MATH.ROUND` | `(x [, decimals])` | Half-away-from-zero; optional second arg = decimal places (see `runtime/mathmod/basic.go`). |
| `Exp` / `Log` / `Log10` | `(x)` | Same as `EXP` / `LOG` / `LOG10`. |
| `Rnd` | `(min, max)` | Random **float** in `[min, max]` (order-corrected if reversed). |
| `Rand` | `(min, max)` | Random **int**, **inclusive** range. |
| `SeedRnd` | `(seed)` | Reseeds the runtime PRNG (same pool as `RNDSEED`). |
| `MILLISECS` | `()` | **Float64** milliseconds since mathmod init (sub-ms precision). |

## String manipulation (UTF-8 runes)

| Command | Arguments | Notes |
|--------|-----------|--------|
| `LEN` | `(s)` | Rune count (not byte length). |
| `LEFT` / `RIGHT` | `(s, n)` | Substrings by rune count. |
| `MID` | `(s, start [, count])` | **`start` is 1-based** rune index; `count` defaults to end of string. |
| `UPPER` / `LOWER` | `(s)` | Unicode case mapping. |
| `INSTR` / `Instr` | `(hay, needle [, start])` | 1-based start index; returns **1-based** position or `0`. |
| `REPLACE` | `(s, find, repl)` | Global literal replacement. |
| `TRIM` / `LTRIM` / `RTRIM` | `(s)` | Whitespace trim. |
| `CHR` / `ASC` | | Codepoint ↔ first character. |
| `STRING` | `(char, n)` **or** `(n, char)` | Repeats first **rune** of `char`; legacy order still accepted. |

## Banks and memory buffers

Banks are `MEM.*` heap objects; Blitz names are aliases.

| Command | Arguments | Notes |
|--------|-----------|--------|
| `CreateBank` | `(size)` | **`MEM.CREATE`** (canonical) / deprecated **`MEM.MAKE`** |
| `FreeBank` | `(bank)` | `MEM.FREE` |
| `BankSize` | `(bank)` | `MEM.SIZE` |
| `ResizeBank` | `(bank, newSize)` | `MEM.RESIZE` |
| `CopyBank` | `(src, srcOff, dest, destOff, count)` | Forwards to `MEM.COPY` |
| `PeekByte` / `PokeByte` | `(bank, off [, val])` | 8-bit |
| `PeekShort` / `PokeShort` | | 16-bit LE |
| `PeekInt` / `PokeInt` | | 32-bit LE |
| `PeekFloat` / `PokeFloat` | | **64-bit LE IEEE float** (`MEM.GETDOUBLE` / `SETDOUBLE`; **8-byte aligned** offset). |
| `ReadBank` | `(bank, file, offset, count)` | Read bytes from file stream into bank |
| `WriteBank` | `(bank, file, offset, count)` | Write bytes from bank to file stream |

## Time and system info

| Command | Arguments | Notes |
|--------|-----------|--------|
| `CurrentTime` | `()` | `"HH:MM:SS"` (local wall clock); same as `TIME`. |
| `CurrentDate` | `()` | `"DD Mon YYYY"` (e.g. `02 Jan 2006` layout). |
| `MilliSecs` | `()` | **Float64** ms (Raylib time when CGO; see `runtime/time/millis_*.go`). |
| `Delay` / `DELAY` | `(ms)` | Sleep; **numeric** ms (int or float). |
| `SystemProperty` | `(key)` | Keys: `os`, `os_version`, `arch`, `cpu_cores` (int), `compiler` (Go runtime string). Unknown keys → `""`. |

---

For the full machine-checked arity table, see `compiler/builtinmanifest/commands.json` (keys are normalized uppercase).
