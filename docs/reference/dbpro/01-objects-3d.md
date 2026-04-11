# DBPro — Objects / 3D engine (core)

See [README.md](README.md) for legend. Deeper Blitz-style entity naming also appears in [../BLITZ3D.md](../BLITZ3D.md).

---

## Create / load / delete

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE OBJECT (obj, file)** | ≈ **`Model.Load()`** / **`Entity.Load()`** | DBPro reused integer slots; moon uses **handles** + **entity ids** depending on path. [MODEL.md](../MODEL.md), [ENTITY.md](../ENTITY.md). |
| **LOAD OBJECT (file, obj)** | ≈ **`Model.Load()`**, **`Entity.Load()`** | Order of args differs. |
| **DELETE OBJECT (obj)** | ≈ **`Entity.Free()`**, **`Model.Free()`** | What to call depends on whether you used **entity** or **model** handle. |
| **CLONE OBJECT** / **INSTANCE OBJECT** / **COPY OBJECT** | ≈ **`Entity.Copy()`**, **`Model.Clone()`** | No single “instance” keyword; see manifest. |
| **HIDE OBJECT** / **SHOW OBJECT** | ≈ **`Entity.Visible()`**, **`Model.Hide()`** / **`Model.Show()`** | |
| **LOCK OBJECT ON** / **OFF** | — | Use **physics freeze** / **custom flag** if needed; not one builtin. |

---

## Position / rotate / scale

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **POSITION OBJECT (obj, x, y, z)** | ✓ **`Entity.Position()`**, **`Model.SetPos()`** / transforms | |
| **ROTATE OBJECT** | ✓ **`Entity.SetRotation()`**, **`Model.SetRot()`** | Radians vs degrees: check each command. |
| **MOVE OBJECT (obj, distance)** | ≈ **`Entity.Move()`**, **`Model.Move()`** | Axis semantics differ from DBPro “forward”. |
| **TURN OBJECT LEFT/RIGHT/UP/DOWN** | ≈ **`Entity.Turn()`**, **`Model.Rotate()`** | Incremental rotation. |
| **SCALE OBJECT (obj, sx, sy, sz)** | ✓ **`Entity.Scale()`**, **`Model.SetScale()`** | |
| **POINT OBJECT (obj, x, y, z)** | ✓ **`Entity.LookAt()`** | |

---

## Getters

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **OBJECT POSITION X/Y/Z** | ✓ **`Entity.X()`**, **`Model.X()`** | |
| **OBJECT ANGLE X/Y/Z** | ≈ **`Entity.Pitch()`**, **`Model.GetRot()`** | |
| **OBJECT SIZE X/Y/Z** | ≈ **`Entity.Scale()`**, **`Model.GetScale()`** | |

---

## Appearance (color / alpha / FX)

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **COLOR OBJECT** | ✓ **`Entity.Color()`**, **`Model.SetColor()`** | |
| **SET OBJECT AMBIENT/DIFFUSE** | ≈ **`Material.SetColor()`**, **`Model.SetMetal()`**, **`Model.SetRough()`** | No full fixed-function material stack like DBPro. |
| **SET OBJECT ALPHA** | ✓ **`Entity.Alpha()`**, model alpha paths | |
| **SET OBJECT LIGHT** / **WIREFRAME** / **TRANSPARENCY** / **CULL** / **FILTER** / **FOG** / **SHADING** / **EFFECT** | ≈ **`Render.Clear()`**, **`Model.DrawWires()`**, **`Light.Make()`**, **`Shader.Load()`** | Feature split across modules. |

---

## Textures (object)

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **TEXTURE OBJECT** / **SET OBJECT TEXTURE*** | ≈ **`Entity.Texture()`**, **`Texture.Load()`**, **`Model` material** | Multi-stage UV pipeline differs; see [TEXTURE.md](../TEXTURE.md). |

---

## Collision

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **SET OBJECT COLLISION*** | ≈ **`Entity.Type()`**, **`Physics3D.Start()`** | Not a single DBPro-style collision setup. [COLLISION.md](../COLLISION.md), [PHYSICS3D.md](../PHYSICS3D.md). |
| **OBJECT COLLISION** / **OBJECT HIT** | ≈ **`Entity.Collided()`**, **`Entity.Pick()`** | |
| **OBJECT COLLISION X/Y/Z** | ✓ **`Entity.CollisionX()`** … | |
