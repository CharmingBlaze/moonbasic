# Blitz3D-style helpers

High-level commands inspired by Blitz3D naming. **3D camera** and **entity** natives require **CGO** (Raylib). **Input** aliases work wherever **`INPUT.*`** works.

---

## Camera (`Camera.*` → `CAMERA.*`)

| Command | Purpose |
|--------|---------|
| **`Camera.Turn(cam, dpitch#, dyaw#, droll#)`** | Incremental rotation (**radians**): yaw around world **+Y**, pitch around camera right, roll around view. Keeps eye–target distance. |
| **`Camera.Rotate(cam, pitch#, yaw#, roll#)`** | Absolute orientation (**radians**): builds forward from pitch/yaw, applies roll to **up**, keeps distance from eye to target. |
| **`Camera.Orbit(...)`** | Same arguments as **`Camera.SetOrbit`** — spherical orbit around a target (**alias**). |
| **`Camera.Zoom(cam, amount#)`** | Adds **amount** to vertical FOV (**degrees**), clamped **10–120**. |
| **`Camera.Follow(cam, tx#, ty#, tz#, yaw#, dist#, height#, smooth#)`** | Third-person follow: camera lerps behind **(tx,ty,tz)** on **XZ** at **yaw**, fixed world **height** for the eye, target lerps toward the subject. **`smooth`** is blended with frame time (~`smooth×8×dt` cap 1). Uses **`Time.Delta`** internally. |
| **`Camera.FollowEntity(cam, entity#, dist#, height#, smooth#)`** | Same as **`Follow`**, but target position and **yaw** come from an **entity** id (see below). |
| **`Camera.OrbitEntity(cam, entity#, yaw#, pitch#, dist#)`** | Orbits the camera around the entity’s **world** position using the same math as **`Camera.SetOrbit`** (see [CAMERA.md](CAMERA.md)). |

See also [CAMERA.md](CAMERA.md) for **`SetOrbit`**, **`OrbitAround`**, **`GetRay`**, etc.

---

## Entities (`Entity.*` → `ENTITY.*`)

Lightweight **integer ids** (not heap handles). **Platforms**: **`Entity.CreateBox`** / primitives below return ids; place with **`Entity.SetPosition`** or **`Entity.PositionEntity`**. **Players**: **`Entity.Create`** makes a small dynamic sphere character (default radius **0.5**, gravity **-28**).

Optional **`global`** arguments (where documented) use **`TRUE`/`FALSE`** or **`1`/`0`**: when **`TRUE`**, position/rotation/translation are interpreted or applied in **world** space vs parent-local space (see parenting).

### Creation

| Command | Purpose |
|--------|---------|
| **`Entity.Create()`** / **`Entity.CreateEntity()`** | Empty marker entity; add **`Radius`** / **`Box`** to make it a dynamic body. Returns **id**. |
| **`Entity.CreateBox` / `CreateCube(w#, h#, d#)`** | Static axis-aligned box (**full** dimensions), centred at origin until moved. |
| **`Entity.CreateSphere(radius#, segments)`** | Static **sphere** (segments ≥ 3); used for drawing and static collision. |
| **`Entity.CreateCylinder(radius#, height#, segments)`** | Static **cylinder** (simplified collision). |
| **`Entity.CreatePlane(size#)`** | Static **XZ** plane tile (extent **size**). |
| **`Entity.CreateMesh()`** | Procedural mesh placeholder (no file path); **`Entity.Copy`** is not supported until a reload path exists. |
| **`Entity.LoadMesh(path$)`** | **`rl.LoadModel`**: static model for drawing; animations are **not** loaded (use **`LoadAnimatedMesh`**). |
| **`Entity.LoadAnimatedMesh(path$)`** | Loads model plus **`LoadModelAnimations`** (Raylib 5.x); animation is advanced in **`Entity.Update`** when **`Entity.Animate`** sets a non-zero speed. |

### Transforms (MoonBASIC names)

| Command | Purpose |
|--------|---------|
| **`Entity.SetPosition(id, x#, y#, z# [, global])`** | Set position; optional **`global`** for parented entities. Aliases: **`PositionEntity`**. |
| **`Entity.RotateEntity(id, pitch#, yaw#, roll# [, global])`** | **Absolute** euler (**radians**). |
| **`Entity.TurnEntity(id, dpitch#, dyaw#, droll# [, global])`** | **Add** to euler (same as **`Entity.Rotate`**). |
| **`Entity.Move` / `MoveEntity(id, f#, r#, u#)`** | Move along **local** forward/right/up from pitch/yaw. |
| **`Entity.Translate` / `TranslateEntity(id, dx#, dy#, dz# [, global])`** | World-space delta (optional **`global`** flag matches Blitz overloads). |
| **`Entity.Scale(id, sx#, sy#, sz#)`** | Non-uniform scale. |

### Getters

| Command | Purpose |
|--------|---------|
| **`Entity.EntityX` / `Y` / `Z(id [, global])`** | Position component (world by default; **`global`** controls local vs world when parented). |
| **`Entity.EntityPitch` / `Yaw` / `Roll(id [, global])`** | Orientation (**radians**). |
| **`Entity.GetPosition(id)`** | **`Vec3`** handle for world centre. |

### Hierarchy

| Command | Purpose |
|--------|---------|
| **`Entity.Parent(child, parent [, global])`** | Parent **child** to **parent** (integer id **0** clears). |
| **`Entity.ParentClear(child)`** | Detach from parent. |

### Visuals (drawing / materials)

Static primitives and loaded models participate in **`Entity.DrawAll`** (sorted by **`Entity.Order`** when set). **`Entity.Texture`** accepts a texture **handle** from **`Texture.Load`** (or **0** to clear).

| Command | Purpose |
|--------|---------|
| **`Entity.Color`**, **`Alpha`**, **`Shininess`**, **`Texture`**, **`FX`**, **`Blend`**, **`Order`** | Material-style fields for drawn entities. |

### Collision and hit data

| Command | Purpose |
|--------|---------|
| **`Entity.Radius`**, **`Box`**, **`Type`**, **`Collide`** | Blitz-style **type** masks and **`Collide`**: which other types this entity hits. |
| **`Entity.Collided`**, **`CollisionOther`** | Pairwise **dynamic**–**dynamic** overlap from last **`Update`**. |
| **`Entity.CollisionX/Y/Z`**, **`CollisionNX/Y/Z`** | Last resolved **contact** point and **normal** (static resolution, sphere/box). |
| **`Entity.Distance(a, b)`** | Distance between **world** positions. |

### Physics helpers

| Command | Purpose |
|--------|---------|
| **`Entity.SetGravity` / `Gravity`**, **`Entity.Velocity`**, **`Entity.AddForce`**, **`Entity.Jump`** | Simple integrator in **`Entity.Update`**. |
| **`Entity.Slide(id [, on])`** | Slide along surfaces when resolving static hits. |
| **`Entity.Pick`**, **`Entity.PickMode`** | Forward ray pick from entity (simplified). |
| **`Entity.Floor(id)`** | Highest static **top** under entity **XZ** (same family as **`BOXTOPLAND`**). |
| **`Entity.MoveRelative`**, **`Entity.ApplyGravity`**, **`Entity.Grounded`** | Convenience wrappers (**`Grounded`** ↔ internal **`onGround`**). |
| **`Entity.SetMass`**, **`SetFriction`**, **`SetBounce`** | Used in collision response. |

### Pointing

| **`Entity.PointEntity(id, targetId)`** | Aim **+Z** toward another entity’s **world** position (yaw; pitch flattened for stability). |
| **`Entity.AlignToVector(id, vx#, vy#, vz#, axis)`** | Align local **+Z** to a **world** direction (**axis** reserved / simplified). |

### Animation

| Command | Purpose |
|--------|---------|
| **`Entity.Animate(id [, mode, speed#])`** | **`speed`** **0** = frozen (use **`SetAnimTime`** only); non-zero = advance in **`Update`**. |
| **`Entity.SetAnimTime`**, **`AnimTime`**, **`AnimLength`** | Frame index / clip length (Raylib **`ModelAnimation.FrameCount`** when loaded). |

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
| **`JoyX` / `JoyY` / `JoyButton`** | **`Input.JoyX`**, **`Input.JoyY`**, **`Input.JoyButton`**, or **`GAME.***` — default **gamepad 0**, left stick **X**/**Y**; optional args **`(gamepad#)`** or **`(gamepad#, axis#)`** for **JoyX/JoyY**; **`JoyButton(button#)`** or **`(gamepad#, button#)`**. |

Axis and button indices follow **Raylib** (`GamepadAxis*`, `GamepadButton*`).

---

## Related

- [BLITZ2025.md](BLITZ2025.md) — wider Blitz-style name map (scene files, groups, **`PHYSICS.*`** aliases, file/JSON helpers)  
- [CAMERA.md](CAMERA.md) — full camera API  
- [INPUT.md](INPUT.md) — keyboard, mouse, gamepad actions  
- [COLLISION.md](COLLISION.md) — **`BOXTOPLAND`**, overlap tests  
