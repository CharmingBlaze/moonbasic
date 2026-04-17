# Blitz3D-style helpers

High-level commands inspired by Blitz3D naming. **3D camera** and **entity** natives require **CGO** (Raylib). **Input** aliases work wherever **`INPUT.*`** works.

**Full Blitz name ↔ moonBASIC map** (2D, world, entity, camera, mesh, texture, audio, file, math): **[BLITZ_COMMAND_INDEX.md](BLITZ_COMMAND_INDEX.md)**. **Curated “essential” list** (entities, meshes, camera, 2D, **`CURVEVALUE`**, **`RAND`**): **[BLITZ_ESSENTIAL_API.md](BLITZ_ESSENTIAL_API.md)**.

---

## Core Workflow

This page is a compatibility alias reference. Blitz-style names map directly to moonBASIC registry commands:

1. Use **`Camera.Make()`** / **`CAMERA.CREATE`** → position → begin/end render pair.
2. Use **`Entity.Turn()`** / **`Entity.Move()`** each frame.
3. Drive the loop with **`RENDER.CLEAR`** → **`RENDER.BEGIN3D`** → **`ENTITY.DRAWALL`** → **`RENDER.END3D`** → **`RENDER.FRAME`**.

---

## Raylib render pipeline (replaces Blitz **Flip** / **RenderWorld** / **UpdateWorld**)

Classic Blitz3D drove the frame with **Flip**, **RenderWorld**, and **UpdateWorld**. moonBASIC targets **Raylib**: you control **when** the GPU draws and **which** camera is active. Those three Blitz entrypoints are **not** separate builtins — use the following instead:

| Blitz3D | moonBASIC / Raylib |
|---------|---------------------|
| **Flip** | **`RENDER.FRAME`** — presents the back buffer (call once at end of each frame after all drawing). |
| **RenderWorld** | After **`RENDER.CLEAR`**: **`RENDER.BEGIN3D(cam)`** … **`ENTITY.DRAWALL`** (and any **`DRAW3D.*`**) … **`RENDER.END3D()`** (same as **`CAMERA.BEGIN`/`CAMERA.END`**) — then 2D HUD with **`CAMERA2D.BEGIN()`** … **`CAMERA2D.END()`** if needed. |
| **UpdateWorld** | **`ENTITY.UPDATE(dt)`** with **`dt = TIME.DELTA()`** (or **`DT()`**) — steps entities, collision discovery, particles tied to the entity system, etc. Optional **`UPDATEPHYSICS`** bundles **`ENTITY.UPDATE`** + world/physics steps. |

**Note:** **`FlipMesh`** (mesh winding) is unrelated to **Flip** — see [MESH.md](MESH.md) / entity mesh helpers.

---

## Mental model: Blitz3D → moonBASIC

| Blitz3D idea | moonBASIC |
|--------------|-----------|
| **`Graphics3D width, height, depth`** | **`WINDOW.OPEN(w, h, title)`** — depth is handled by the 3D camera / Z-buffer, not a third dimension argument. |
| **`AmbientLight` / `CameraClsMode`** | **`RENDER.CLEAR(r,g,b)`** before **`RENDER.BEGIN3D(cam)`** (or **`CAMERA.BEGIN`**); sky colour is your clear. |
| **`CreateCamera` / orbit the view** | **`cam = CreateCamera()`** (→ **`CAMERA.CREATE`**) then **`CAMERA.SETORBIT`** or **`CAMERA.ORBIT`** (same math — see [CAMERA.md](CAMERA.md)). Third-person yaw/pitch/distance helpers: **`ORBITYAWDELTA`**, **`ORBITPITCHDELTA`**, **`ORBITDISTDELTA`** ([GAMEHELPERS.md](GAMEHELPERS.md)). |
| **`WireCube` / `Cube` (immediate)** | Short globals **`WIRECUBE`** / **`BOX`** (same as **`DRAW3D.CUBEWIRES`** / **`DRAW3D.CUBE`**) — see [DRAW3D.md](DRAW3D.md). Optional OOP-style **`DRAWCUBE()`** / **`DRAWSPHERE()`** wrappers: [DRAW_WRAPPERS.md](DRAW_WRAPPERS.md) (distinct from **`CUBE()`** entities). |
| **`MoveEntity` / `PositionEntity`** | **`Entity.MoveEntity`** / **`Entity.PositionEntity`** for **entity ids**; **dot-syntax** on **`CUBE`/`SPHERE`** handles below; for raw floats + **`LANDBOXES`** / **`LANDBOX`**, see **`examples/mario64/main_orbit_simple.mb`**. |
| **`KeyHit` / `KeyDown`** | Flat **`KEYHIT(key)`** / **`KEYDOWN(key)`** (also **`GAME.KEYHIT`** / **`GAME.KEYDOWN`**). |
| **`MouseXSpeed()` / `MouseYSpeed()`** | **`MOUSEXSPEED`** / **`MOUSEYSPEED`** or **`MDX`** / **`MDY`**. |
| **`EndGraphics`** | **`WINDOW.CLOSE()`** after **`ERASE ALL`** if you used VM handles — [MEMORY.md](../MEMORY.md). |

**Reference sketch** (immediate-mode 3D + orbit, no entities): [`examples/mario64/main_orbit_simple.mb`](../../examples/mario64/main_orbit_simple.mb) — **`DT`**, **`KEYHIT`**, **`BOX`**, **`WIRECUBE`**, **`FLAT`**, **`GRID3`**, **`CAMERA.SETORBIT`**, **`ERASE ALL`**.

---

## Dot-syntax entities (`CUBE` / `SPHERE` → `ENTITYREF`)

Blitz3D used commands like **`ScaleEntity`**, **`PositionEntity`**, **`EntityColor`**. moonBASIC keeps the same **runtime** (**`ENTITY.*`** with integer ids) but also exposes **short constructors** that return a **heap handle** (**`ENTITYREF`**) so you can use **modern dot-syntax**:

```moonbasic
cube = CUBE()           ; or CUBE(w, h, d)
cube.Pos(0, 1.2, 5)
cube.Scale(3, 0.35, 3)
cube.Col(255, 140, 96)

player = SPHERE(0.45)
player.Pos(px, py, pz)
player.Col(40, 95, 200)
```

Handle methods are normal calls: **`receiver.Method(args)`** (parentheses required), not Blitz’s space-separated **`PositionEntity id, x, y, z`** form.

| Dot method | Maps to | Notes |
|------------|---------|--------|
| **`Pos`** | **`Entity.SetPosition`** | Position in world space |
| **`Scale`** | **`Entity.Scale`** | Non-uniform scale |
| **`Rot`** | **`Entity.RotateEntity`** | **Absolute** euler (**radians**): pitch, yaw, roll |
| **`Turn`** | **`Entity.TurnEntity`** | **Incremental** euler (**radians**) |
| **`Move`** | **`Entity.Translate`** | World-space delta (**`Entity.Translate`**) |
| **`Col`** / **`Color`** | **`Entity.Color`** | RGB **0–255** |
| **`A`** | **`Entity.Alpha`** | |
| **`Free`** | **`Entity.Free`** | Unloads native resources; invalidates the handle |
| **`Hide`** / **`Show`** | **`Entity.Hide`** / **`Show`** | |

**Immediate-mode** drawing (**`BOX`**, **`WIRECUBE`**, **`DRAW3D.SPHERE`**, …) is unchanged: those are **not** entity constructors. **`SPHERE(radius)`** here is the **entity** primitive; for one-off draws, use **`BALL`** / **`DRAW3D.SPHERE`** (see [DRAW3D.md](DRAW3D.md)).

---

## Short camera handle (`CAM()` / `Cam()`)

**`cam = CAM()`** (alias of **`CAMERA.CREATE`**) returns the same **`Camera3D`** handle as **`CreateCamera()`** / **`Camera.Make()`** (deprecated). Prefer **dot methods** on the handle for short code:

| Dot (as a call) | Registry |
|-----|----------|
| **`cam.Pos(x,y,z)`** | **`Camera.SetPos`** |
| **`cam.FOV(deg)`** | **`Camera.SetFOV`** |
| **`cam.Look(...)`** / **`cam.LookAt(...)`** | **`Camera.LookAt`** (same as **`SetTarget`**) |
| **`cam.Orbit(entity, distance)`** | **`Camera.Orbit`** (3-arg **entity** orbit-follow — engine yaw/pitch/dist; see [CAMERA.md](CAMERA.md)) |
| **`cam.Orbit(tx,ty,tz,yaw,pitch,dist)`** / **`cam.SetOrbit(...)`** | **`Camera.SetOrbit`** / 7-arg **`Orbit`** |
| **`cam.Zoom(amount)`** | **`Camera.Zoom`** |

See [CAMERA.md](CAMERA.md) for full **`CAMERA.*`** reference.

---

## Landing helper alias

**`LANDBOX`** is an alias of **`LANDBOXES`** (same 12 arguments). See [GAMEHELPERS.md](GAMEHELPERS.md).

---

## Camera (`Camera.*` → `CAMERA.*`)

| Command | Purpose |
|--------|---------|
| **`Camera.Turn(cam, dpitch, dyaw, droll)`** | Incremental rotation (**radians**): yaw around world **+Y**, pitch around camera right, roll around view. Keeps eye–target distance. |
| **`Camera.Rotate(cam, pitch, yaw, roll)`** | Absolute orientation (**radians**): builds forward from pitch/yaw, applies roll to **up**, keeps distance from eye to target. |
| **`Camera.Orbit(cam, entity, dist)`** | **Entity** orbit-follow: internal yaw/pitch/distance + input (see [CAMERA.md](CAMERA.md)). |
| **`Camera.Orbit(cam, tx, ty, tz, yaw, pitch, dist)`** | Same arguments as **`Camera.SetOrbit`** — explicit spherical orbit (**7-arg** overload). |
| **`Camera.Zoom(cam, amount)`** | Adds **amount** to vertical FOV (**degrees**), clamped **10–120**. |
| **`Camera.Follow(cam, tx, ty, tz, yaw, dist, height, smooth)`** | Third-person follow: camera lerps behind **(tx,ty,tz)** on **XZ** at **yaw**, fixed world **height** for the eye, target lerps toward the subject. **`smooth`** is blended with frame time (~`smooth×8×dt` cap 1). Uses **`TIME.DELTA`** internally (registry **`CAMERA.FOLLOW`**). |
| **`Camera.FollowEntity(cam, entity, dist, height, smooth)`** | Same as **`Follow`**, but target position and **yaw** come from an **entity** id (see below). |
| **`Camera.OrbitEntity(cam, entity, yaw, pitch, dist)`** | Orbits the camera around the entity’s **world** position using the same math as **`Camera.SetOrbit`** (see [CAMERA.md](CAMERA.md)). |

See also [CAMERA.md](CAMERA.md) for **`SetOrbit`**, **`OrbitAround`**, **`GetRay`**, etc.

---

## Entities (`Entity.*` → `ENTITY.*`)

Lightweight **integer ids** (not heap handles). **Platforms**: **`Entity.CreateBox`** / primitives below return ids; place with **`Entity.SetPosition`** or **`Entity.PositionEntity`**. **Players**: **`Entity.Create`** makes a small dynamic sphere character (default radius **0.5**, gravity **-28**).

Optional **`global`** arguments (where documented) use **`TRUE`/`FALSE`** or **`1`/`0`**: when **`TRUE`**, position/rotation/translation are interpreted or applied in **world** space vs parent-local space (see parenting).

### Creation 

| Command | Purpose |
|--------|---------|
| **`Entity.Create()`** / **`Entity.CreateEntity()`** | Empty marker entity; add **`Radius`** / **`Box`** to make it a dynamic body. Returns **id**. |
| **`Entity.CreateBox` / `CreateCube(w, h, d)`** | Static axis-aligned box (**full** dimensions), centred at origin until moved. |
| **`Entity.CreateSphere(radius, segments)`** | Static **sphere** (segments ≥ 3); used for drawing and static collision. |
| **`Entity.CreateCylinder(radius, height, segments)`** | Static **cylinder** (simplified collision). |
| **`Entity.CreatePlane(size)`** | Static **XZ** plane tile (extent **size**). |
| **`Entity.CreateMesh()`** | Procedural mesh placeholder (no file path); **`Entity.Copy`** is not supported until a reload path exists. |
| **`Entity.LoadMesh(path)`** | **`rl.LoadModel`**: static model for drawing; animations are **not** loaded (use **`LoadAnimatedMesh`**). |
| **`Entity.LoadAnimatedMesh(path)`** | Loads model plus **`LoadModelAnimations`** (Raylib 5.x); animation is advanced in **`Entity.Update`** when **`Entity.Animate`** sets a non-zero speed. |

---

### Transforms (MoonBASIC names) 

| Command | Purpose |
|--------|---------|
| **`Entity.SetPosition(id, x, y, z [, global])`** | Set position; optional **`global`** for parented entities. Aliases: **`PositionEntity`**. |
| **`Entity.RotateEntity(id, pitch, yaw, roll [, global])`** | **Absolute** euler (**radians**). |
| **`Entity.TurnEntity(id, dpitch, dyaw, droll [, global])`** | **Add** to euler (same as **`Entity.Rotate`**). |
| **`Entity.Move` / `MoveEntity(id, f, r, u)`** | Move along **local** forward/right/up from pitch/yaw. |
| **`Entity.Translate` / `TranslateEntity(id, dx, dy, dz [, global])`** | World-space delta (optional **`global`** flag matches Blitz overloads). |
| **`Entity.Scale(id, sx, sy, sz)`** | Non-uniform scale. |

---

### Getters 

| Command | Purpose |
|--------|---------|
| **`Entity.EntityX` / `Y` / `Z(id [, global])`** | Position component (world by default; **`global`** controls local vs world when parented). |
| **`Entity.EntityPitch` / `Yaw` / `Roll(id [, global])`** | Orientation (**radians**). |
| **`Entity.GetPosition(id)`** | **`Vec3`** handle for world centre. |

---

### Hierarchy 

| Command | Purpose |
|--------|---------|
| **`Entity.Parent(child, parent [, global])`** | Parent **child** to **parent** (integer id **0** clears). |
| **`Entity.ParentClear(child)`** | Detach from parent. |

---

### Visuals (drawing / materials) 

Static primitives and loaded models participate in **`Entity.DrawAll`** (sorted by **`Entity.Order`** when set). **`Entity.Texture`** accepts a texture **handle** from **`Texture.Load`** (or **0** to clear).

| Command | Purpose |
|--------|---------|
| **`Entity.Color`**, **`Alpha`**, **`Shininess`**, **`Texture`**, **`FX`**, **`Blend`**, **`Order`** | Material-style fields for drawn entities. |

---

### Collision and hit data 

| Command | Purpose |
|--------|---------|
| **`Collisions` / `COLLISIONS(srcType, dstType, method, response)`** | Register **rule-based** pairs (e.g. sphere–box **`method`** **2**); resolved inside **`ENTITY.UPDATE`**. |
| **`EntityType` / `ENTITY.TYPE`** | Integer **collision type** id for **`COLLISIONS`** **`src`/`dst`** matching. |
| **`EntityHitsType(entity, type)`** → **bool** | **`TRUE`** if **`entity`** hit any other entity with **`EntityType == type`** after the last **`ENTITY.UPDATE`**. |
| **`ENTITYCOLLIDED(entity, type)`** → **int** | Same hit test; returns **other entity id** or **0**. |
| **`EntityCollided(a, b)`** → **bool** | **Jolt** pairwise contact (Linux + physics buffer link); **not** the same as **`EntityHitsType`**. |
| **`Entity.Radius`**, **`Box`**, **`Type`**, **`Collide`** | Blitz-style **type** masks and **`Collide`**: which other types this entity hits. |
| **`Entity.Collided`**, **`CollisionOther`** | Pairwise **dynamic**–**dynamic** overlap from last **`Update`**. |
| **`Entity.CollisionX/Y/Z`**, **`CollisionNX/Y/Z`** | Last resolved **contact** point and **normal** (static resolution, sphere/box). |
| **`Entity.Distance(a, b)`** | Distance between **world** positions. |
| **`Entity.TFormVector` / `TFormVector(x,y,z, src, dst)`** | Direction in **`src`** local space → **3-float array handle** in **`dst`** local space (see [ENTITY.md](ENTITY.md)). |

---

### Physics helpers 

| Command | Purpose |
|--------|---------|
| **`Entity.SetGravity` / `Gravity`**, **`Entity.Velocity`**, **`Entity.AddForce`**, **`Entity.Jump`** | Simple integrator in **`Entity.Update`**. |
| **`Entity.Slide(id [, on])`** | Slide along surfaces when resolving static hits. |
| **`Entity.Pick`**, **`Entity.PickMode`** | Forward ray pick from entity (simplified). |
| **`Entity.Floor(id)`** | Highest static **top** under entity **XZ** (same family as **`BOXTOPLAND`**). |
| **`Entity.MoveRelative`**, **`Entity.ApplyGravity`**, **`Entity.Grounded`** | Convenience wrappers (**`Grounded`** ↔ internal **`onGround`**). |
| **`Entity.SetMass`**, **`SetFriction`**, **`SetBounce`** | Used in collision response. |

---

### Pointing 

| **`Entity.PointEntity(id, targetId)`** | Aim **+Z** toward another entity’s **world** position (yaw; pitch flattened for stability). |
| **`Entity.AlignToVector(id, vx, vy, vz, axis)`** | Align local **+Z** to a **world** direction (**axis** reserved / simplified). |

---

### Animation 

| Command | Purpose |
|--------|---------|
| **`Entity.Animate(id [, mode, speed])`** | **`speed`** **0** = frozen (use **`SetAnimTime`** only); non-zero = advance in **`Update`**. |
| **`Entity.SetAnimTime`**, **`AnimTime`**, **`AnimLength`** | Frame index / clip length (Raylib **`ModelAnimation.FrameCount`** when loaded). |

---

### Management 

| **`Entity.Hide`**, **`Show`**, **`Free`**, **`Copy`**, **`SetName`**, **`Find`** | Visibility, unload (**`UnloadModel`** + animations), duplicate (**reload** from path for models), registry by **case-insensitive** name. |

**Notes:** Physics is **simple** (slide/separate, not a full physics engine). For **Jolt**/**CharacterController**, keep using [PHYSICS3D.md](PHYSICS3D.md) / [CHARCONTROLLER.md](CHARCONTROLLER.md).

---

## Input aliases

| Blitz-style | Equivalent |
|-------------|------------|
| **`KeyHit(key)`** | **`Input.KeyHit`** / **`KeyHit`** / **`GAME.KEYHIT`** — same as **`KeyPressed`**. |
| **`KeyDown(key)`** | **`Input.KeyDown`** / **`KeyDown`** (unchanged). |
| **`MouseXSpeed` / `MouseYSpeed`** | **`Input.MouseXSpeed`**, **`Input.MouseYSpeed`**, or top-level **`MouseXSpeed`**, **`MouseYSpeed`** — same as mouse delta **X** / **Y**. |
| **`JoyX` / `JoyY` / `JoyButton`** | **`Input.JoyX`**, **`Input.JoyY`**, **`Input.JoyButton`**, or **`GAME.***` — default **gamepad 0**, left stick **X**/**Y**; optional args **`(gamepad)`** or **`(gamepad, axis)`** for **JoyX/JoyY**; **`JoyButton(button)`** or **`(gamepad, button)`**. |

Axis and button indices follow **Raylib** (`GamepadAxis*`, `GamepadButton*`).

---

## Full Example

Blitz3D-style spin demo: load a mesh, orbit a camera, drive with keyboard input.

```basic
WINDOW.OPEN(960, 540, "Blitz3D Style")
WINDOW.SETFPS(60)

cam = Camera.Make()
Camera.SetPos(cam, 0, 4, -10)
Camera.SetTarget(cam, 0, 0, 0)

cube = Entity.CreateCube(2.0)
Entity.SetPos(cube, 0, 0, 0)

yaw = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = DT()
    IF KeyDown(KEY_LEFT)  THEN yaw = yaw - 90 * dt
    IF KeyDown(KEY_RIGHT) THEN yaw = yaw + 90 * dt
    Entity.SetRotation(cube, 0, yaw, 0)
    Entity.Update(dt)

    RENDER.CLEAR(15, 20, 35)
    RENDER.BEGIN3D(cam)
        Entity.DrawAll()
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

Entity.Free(cube)
Camera.Free(cam)
WINDOW.CLOSE()
```

---

## See also

- [BLITZ2025.md](BLITZ2025.md) — wider Blitz-style name map (scene files, groups, **`PHYSICS.*`** aliases, file/JSON helpers)
- [BLITZ_ESSENTIAL_API.md](BLITZ_ESSENTIAL_API.md) — curated essential list
- [CAMERA.md](CAMERA.md) — full camera API
- [INPUT.md](INPUT.md) — keyboard, mouse, gamepad actions
- [COLLISION.md](COLLISION.md) — **`BOXTOPLAND`**, overlap tests
