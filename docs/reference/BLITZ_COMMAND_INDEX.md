# Blitz-Style Command Index

A comprehensive mapping of classic Blitz3D and BlitzPlus command names to their modern MoonBASIC equivalents.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Namespace First:** While global aliases like `CreateCube` exist, we recommend using `ENTITY.CREATECUBE` or the short `CUBE()` constructor.
2. **Method Chaining:** All "Create" commands return handles. Configure them fluently: `CUBE().POS(0,5,0).COLOR(255,0,0)`.
3. **Logic & Render:** Use `ENTITY.UPDATE` to solve physics and `ENTITY.DRAWALL` for the main render pass.

---

## 1. Lifecycle & Hierarchy

| Command | Arguments | MoonBASIC Equivalent |
|---------|-----------|-----------------------|
| `CreatePivot` | `[parent]` | `ENTITY.CREATEPIVOT` |
| `CreateCamera`| `[parent]` | `CAMERA.CREATE` |
| `CreateCube` | `[parent]` | `ENTITY.CREATECUBE` / `CUBE()` |
| `CreateSphere`| `[parent]` | `ENTITY.CREATESPHERE` / `SPHERE()` |
| `LoadMesh` | `path` | `ENTITY.LOADMESH` |
| `FreeEntity` | `id` | `ENTITY.FREE` |
| `EntityParent`| `id, parent`| `ENTITY.PARENT` |

---

## 2. Transformation & State

| Command | Arguments | MoonBASIC Equivalent |
|---------|-----------|-----------------------|
| `PositionEntity`| `id, x, y, z`| `ENTITY.SETPOS` |
| `RotateEntity` | `id, p, y, r`| `ENTITY.SETROT` |
| `ScaleEntity` | `id, x, y, z`| `ENTITY.SETSCALE` |
| `MoveEntity` | `id, f, r, u`| `ENTITY.MOVE` |
| `TurnEntity` | `id, p, y, r`| `ENTITY.TURN` |
| `PointEntity` | `id, target` | `ENTITY.POINT` |

---

## 3. Physics & Collisions

| Command | Arguments | MoonBASIC Equivalent |
|---------|-----------|-----------------------|
| `EntityType` | `id, type` | `ENTITY.TYPE` |
| `EntityRadius`| `id, r` | `ENTITY.RADIUS` |
| `EntityBox` | `id, w, h, d`| `ENTITY.BOX` |
| `Collisions` | `src, dst...`| `COLLISIONS` |
| `EntityHit` | `id, type` | `ENTITY.COLLIDED` |

---

## 4. Picking & Rays

| Command | Arguments | MoonBASIC Equivalent |
|---------|-----------|-----------------------|
| `LinePick` | `x,y,z, dx,dy,dz`| `PICK.LINE` |
| `CameraPick`| `cam, sx, sy` | `PICK.CAMERA` |
| `PickedX` | None | `PICK.X` |
| `PickedEntity`| None | `PICK.ENTITY` |

---

## 5. Input & Utilities

| Command | Arguments | MoonBASIC Equivalent |
|---------|-----------|-----------------------|
| `KeyHit` | `key` | `INPUT.KEYHIT` |
| `KeyDown` | `key` | `INPUT.KEYDOWN` |
| `MouseX` | None | `INPUT.MOUSEX` |
| `MouseDown` | `button` | `INPUT.MOUSEDOWN` |
| `Rnd` | `min, max` | `MATH.RANDOM` |
| `MilliSecs` | None | `TIME.MILLIS` |

---

## 6. Advanced Transforms & Animation

| Command | MoonBASIC Equivalent | Notes |
|---------|-----------------------|-------|
| `TFormPoint` | `ENTITY.TFORMPOINT` | Returns a 3-float array handle. |
| `Animate` | `ENTITY.PLAY` | Skeletal animation playback. |
| `SetAnimTime`| `ENTITY.SETTIME` | Seek to specific animation time. |
| `FindBone` | `ENTITY.FINDBONE` | Returns a bone socket entity. |

---

## 7. Surfaces & Brushes

| Command | MoonBASIC Equivalent | Notes |
|---------|-----------------------|-------|
| `CreateMesh` | `ENTITY.CREATEMESH` | Procedural mesh container. |
| `CreateSurface`| `ENTITY.CREATESURFACE`| Builder surface handle. |
| `AddVertex` | `ENTITY.ADDVERTEX` | CPU-side vertex buffer. |
| `CreateBrush` | `BRUSH.CREATE` | Material property container. |
| `PaintEntity` | `ENTITY.PAINT` | Apply brush to entity. |

---

## 8. IO & System

| Command | MoonBASIC Equivalent | Notes |
|---------|-----------------------|-------|
| `WriteFile` | `FILE.OPENWRITE` | Returns a file handle. |
| `ReadFile` | `FILE.OPENREAD` | Returns a file handle. |
| `WriteLine` | `FILE.WRITELN` | String output to stream. |
| `ReadInt` | `FILE.READINT` | Binary LE integer read. |
| `CreateBank` | `MEM.CREATE` | Heap-allocated byte buffer. |
| `PeekByte` | `MEM.GETBYTE` | Direct memory access. |

---

## 9. String & Math Helpers

| Command | MoonBASIC Equivalent | Notes |
|---------|-----------------------|-------|
| `Sin` / `Cos` | `SIN` / `COS` | Standard degree-based trig. |
| `Left` / `Right`| `LEFT` / `RIGHT` | Unicode-safe substrings. |
| `Instr` | `INSTR` | 1-based character search. |
| `MilliSecs` | `TIME.MILLIS` | High-precision program timer. |

---

## See also

- [API_CONSISTENCY.md](../API_CONSISTENCY.md) — Every registered command name
- [BLITZ3D.md](BLITZ3D.md) — Narrative guide for Blitz developers
- [BLITZ_ESSENTIAL_API.md](BLITZ_ESSENTIAL_API.md) — Concise quick-reference
- [BLITZ2025.md](BLITZ2025.md) — Modern vision for MoonBASIC
