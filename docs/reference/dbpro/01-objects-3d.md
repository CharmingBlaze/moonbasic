# DBPro — Objects / 3D engine (core)

See [README.md](README.md) for legend. Deeper Blitz-style entity naming also appears in [../BLITZ3D.md](../BLITZ3D.md).

---

## Create / load / delete

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE OBJECT (obj, file)** | ≈ **`MODEL.LOAD`** / **`ENTITY.LOADMESH`** / **`ENTITY.CREATE*`** | DBPro reused integer slots; moon uses **handles** + **entity ids** depending on path. [MODEL.md](../MODEL.md), [ENTITY.md](../ENTITY.md). |
| **LOAD OBJECT (file, obj)** | ≈ **`MODEL.LOAD`**, **`ENTITY.LOADMESH`** | Order of args differs. |
| **DELETE OBJECT (obj)** | ≈ **`ENTITY.FREE`**, **`MODEL` unload patterns** | What to call depends on whether you used **entity** or **model** handle. |
| **CLONE OBJECT** / **INSTANCE OBJECT** / **COPY OBJECT** | ≈ **`ENTITY.COPY`**, **`MODEL`** parenting / duplicate workflows | No single “instance” keyword; see manifest. |
| **HIDE OBJECT** / **SHOW OBJECT** | ≈ **`ENTITY.HIDE`** / **`SHOW`**, **`MODEL.HIDE`** / **`SHOW`** | |
| **LOCK OBJECT ON** / **OFF** | — | Use **physics freeze** / **custom flag** if needed; not one builtin. |

---

## Position / rotate / scale

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **POSITION OBJECT (obj, x, y, z)** | ✓ **`ENTITY.POSITIONENTITY`**, **`ENTITY.SETPOSITION`**, **`MODEL.SETPOS`** / transforms | |
| **ROTATE OBJECT** | ✓ **`ENTITY.ROTATEENTITY`**, **`MODEL.SETROT`**, **`MODEL.ROTATE`** | Radians vs degrees: check each command. |
| **MOVE OBJECT (obj, distance)** | ≈ **`ENTITY.MOVE`**, **`MODEL.MOVE`** | Axis semantics differ from DBPro “forward”. |
| **TURN OBJECT LEFT/RIGHT/UP/DOWN** | ≈ **`ENTITY.TURNENTITY`**, **`MODEL.ROTATE`** | Incremental rotation. |
| **SCALE OBJECT (obj, sx, sy, sz)** | ✓ **`ENTITY.SCALE`**, **`MODEL.SETSCALE`** | |
| **POINT OBJECT (obj, x, y, z)** | ✓ **`ENTITY.POINTENTITY`** | |

---

## Getters

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **OBJECT POSITION X/Y/Z** | ✓ **`ENTITY.ENTITYX`** … **`ENTITYZ`**, **`MODEL.X/Y/Z`** | |
| **OBJECT ANGLE X/Y/Z** | ≈ **`ENTITY.ENTITYPITCH/YAW/ROLL`**, **`MODEL.GETROT`** | |
| **OBJECT SIZE X/Y/Z** | ≈ **`ENTITY` scale getters**, **`MODEL.GETSCALE`** | |

---

## Appearance (color / alpha / FX)

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **COLOR OBJECT** | ✓ **`ENTITY.COLOR`**, **`MODEL.SETCOLOR`** | |
| **SET OBJECT AMBIENT/DIFFUSE/SPECULAR/EMISSIVE** | ≈ **`MATERIAL.*`**, **`MODEL.SETMETAL`**, **`SETROUGH`** | No full fixed-function material stack like DBPro. |
| **SET OBJECT ALPHA** | ✓ **`ENTITY.ALPHA`**, model alpha paths | |
| **SET OBJECT LIGHT** / **WIREFRAME** / **TRANSPARENCY** / **CULL** / **FILTER** / **FOG** / **SHADING** / **EFFECT** | ≈ **`RENDER.*`**, **`MODEL.DRAWWIRES`**, **`LIGHT.*`**, **`SHADER.*`** | Feature split across modules. |

---

## Textures (object)

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **TEXTURE OBJECT** / **SET OBJECT TEXTURE*** | ≈ **`ENTITY.TEXTURE`**, **`TEXTURE.*`**, **`MODEL` material** | Multi-stage UV pipeline differs; see [TEXTURE.md](../TEXTURE.md). |

---

## Collision

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **SET OBJECT COLLISION*** | ≈ **`ENTITY.PICKMODE`**, **`ENTITY.RADIUS`**, **`ENTITY.BOX`**, **`PHYSICS3D.*`** | Not a single DBPro-style collision setup. [COLLISION.md](../COLLISION.md), [PHYSICS3D.md](../PHYSICS3D.md). |
| **OBJECT COLLISION** / **OBJECT HIT** | ≈ **`ENTITY.COLLIDED`**, **`ENTITY.PICK`** | |
| **OBJECT COLLISION X/Y/Z** | ✓ **`ENTITY.COLLISIONX`** … | |
