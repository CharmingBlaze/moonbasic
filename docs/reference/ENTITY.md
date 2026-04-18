# Entity Commands

Lightweight integer-id entities for 3D scene objects: load, position, rotate, parent, animate, and draw.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create with `ENTITY.LOAD`, `ENTITY.CREATE`, `ENTITY.CREATECUBE`, etc.
2. Position/rotate with `ENTITY.SETPOS`, `ENTITY.SETROT`, `ENTITY.SETSCALE`.
3. Call `ENTITY.UPDATE(dt)` each frame to advance physics and animation.
4. Draw with `ENTITY.DRAWALL` (sorted scene) or `ENTITY.DRAW` (single entity).
5. Free with `ENTITY.FREE`.

For scene switching see [SCENE.md](SCENE.md). For level loading see [LEVEL.md](LEVEL.md).

---

### `ENTITY.LOAD(path)`
Loads a 3D model from a file path.

- **Arguments**:
    - `path`: (String) File path to a model (glTF, OBJ, etc.).
- **Returns**: (Integer) The new entity ID.
- **Example**:
    ```basic
    hero = ENTITY.LOAD("hero.glb")
    ```

---

### `ENTITY.CREATECUBE(w, h, d)` / `CREATESPHERE` / `CREATEPLANE`
Creates a primitive geometric entity.

- **Arguments**:
    - `w, h, d`: (Float) Dimensions.
- **Returns**: (Integer) The new entity ID.

---

### `ENTITY.SETPOS(id, x, y, z [, world])`
Sets the position of an entity.

- **Arguments**:
    - `id`: (Integer) Entity ID.
    - `x, y, z`: (Float) Coordinates.
    - `world`: (Boolean, Optional) `TRUE` for world-space, `FALSE` for relative to parent.
- **Returns**: (Integer) The entity ID (for chaining).

---

### `ENTITY.SETROT(id, pitch, yaw, roll)` / `TURN`
Sets the absolute Euler rotation (degrees) or adds a relative rotation.

- **Returns**: (Integer) The entity ID (for chaining).

---

### `ENTITY.SETSCALE(id, sx, sy, sz)`
Sets the non-uniform scale of an entity.

- **Returns**: (Integer) The entity ID (for chaining).

---

### `ENTITY.PARENT(child, parent)` / `UNPARENT`
Establishes a hierarchy between two entities.

- **Returns**: (Integer) The child entity ID (for chaining).

---

### `ENTITY.UPDATE(dt)`
Advances physics, animation, and lerps for all entities.

- **Returns**: (None)

---

### `ENTITY.DRAWALL()` / `DRAW(id)`
Renders the scene graph or a specific entity.

---

### `ENTITY.FREE(id)`
Frees the entity and its resources.

---

## Spatial Macros (`ENTITY.X`, `ENTITY.Y`, ...)

Shorthand for reading/writing coordinates directly from the entity store. These compile to fast bytecode.

- `ENTITY.X(id)` / `ENTITY.Y(id)` / `ENTITY.Z(id)` — Position.
- `ENTITY.P(id)` / `ENTITY.YAW(id)` / `ENTITY.R(id)` — Rotation (Pitch, Yaw, Roll).

---

## Examples

### Creating and Moving an Entity 
```basic
WINDOW.OPEN(1280, 720, "Entity Example")
cam = CAMERA.CREATE()

; Create a cube (w, h, d); SETPOS last arg = world-space when parented
cube = ENTITY.CREATECUBE(2, 2, 2)
ENTITY.SETPOS(cube, 0, 5, 0, TRUE)

WHILE NOT WINDOW.SHOULDCLOSE()
    ENTITY.TURN(cube, 0, 1.0, 0)

    RENDER.CLEAR(0, 0, 0)
    RENDER.Begin3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(cube)
WINDOW.CLOSE()
```

---

## Quick links

- **3D skeletal clips & unified model API** — [ANIMATION_3D.md](ANIMATION_3D.md) (**`ENTITY.PLAY`** / **`PLAYNAME`**, **`ENTITY.LOADANIMATIONS`**, **`ENTITY.DRAW`**, **`GETBOUNDS`**, **`RAYHIT`**, …).
- **glTF level markers / layers** — [LEVEL.md](LEVEL.md) (**`LEVEL.LOAD`**, **`LEVEL.GETSPAWN`**, **`LEVEL.SHOWLAYER`** — distinct from **`SCENE.*`** game scenes).
- **Blitz-style names** (`PositionEntity`, `CreateSphere`, …) are mapped under **`ENTITY.POSITIONENTITY`**, **`ENTITY.CREATESPHERE`**, etc. — see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go), the **[Blitz command index](BLITZ_COMMAND_INDEX.md)**, and the concise **[Blitz essential API](BLITZ_ESSENTIAL_API.md)** (Position vs Move, Rotate vs Turn, Parent, Distance, …).
- **Dot-syntax handles** (`cube.Pos`, `sphere.Turn`) use **`ENTITYREF`** from **`CUBE()`** / **`SPHERE()`** — [BLITZ3D.md](BLITZ3D.md).
- **Scene save/load / clear** — [BLITZ2025.md](BLITZ2025.md), **`ENTITY.SAVESCENE`**, **`ENTITY.LOADSCENE`**, **`ENTITY.CLEARSCENE`**.

## Modern Blitz-style shorthands

- **`UPDATEPHYSICS`** / **`UpdatePhysics()`** — blitzengine bundle: one call per frame for **`ENTITY.UPDATE(TIME.DELTA)`** plus best-effort player / world / 2D / 3D steps (**`PHYSICS3D.UPDATE`** = **`STEP`**). Your draw pass stays **`RENDER.CLEAR`** → **`RENDER.Begin3D`** / **`ENTITY.DRAWALL`** / **`RENDER.END3D`** → **`RENDER.FRAME`**.
- **`DrawEntities()`** — same as **`ENTITY.DRAWALL`** (scene graph draw pass).
- **`CreatePivot()`** — empty transform node (invisible, for parenting).
- **`CreateCube(...)`** — `CreateCube()` / `CreateCube(w,h,d)` / `CreateCube(parent)` / `CreateCube(parent, w,h,d)`; see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go).
- **Jolt (Linux+CGO):** **`ENTITY.LINKPHYSBUFFER(entity, bufferIndex)`** ties an entity to a **`BODY3D`** matrix slot (from **`BODY3D.BUFFERINDEX`**). After **`PHYSICS3D.UPDATE`** / **`PHYSICS3D.STEP`**, translation from the shared buffer updates the entity pose. **`ENTITY.CLEARPHYSBUFFER(entity)`** removes the link.
- **Traffic cop (Jolt-linked entities):** **`ENTITY.ADDPHYSICS`** / **`ENTITY.PHYSICS`** marks the entity **physics-driven** (scripted gravity/velocity integration in **`ENTITY.UPDATE`** is skipped). **`ENTITY.SETPOS`** (canonical) / deprecated **`ENTITY.SETPOSITION`**, or dot **`Pos`**, also teleports the **Jolt** body so meshes do not rubber-band. **`ENTITY.MOVE`** sets **linear velocity** on the body; **`ENTITY.PUSH`** applies an **impulse**. Grounding for **`ENTITY.GROUNDED`** / **`IsGrounded`** uses a short downward ray after the physics sync. See **`examples/mario64/modern_blitz_hop.mb`**.

## Jolt collision groups, queries, and AI helpers (Linux + **`PHYSICS3D.START`**)

| Command | Purpose |
|--------|---------|
| **`ENTITY.SETCOLLISIONGROUP(id, group)`** | Alias of **`ENTITY.COLLISIONLAYER`** — stores **0..31** for **`PICK.LAYERMASK`** / future simulation filtering. |
| **`ENTITY.CHECKCOLLISION(a, b)`** | Same as **`EntityCollided`** — **`true`** if the pair had a Jolt contact since the last **`PHYSICS3D.UPDATE`** / **`STEP`** (requires **`ENTITY.LINKPHYSBUFFER`** on both sides where applicable). |
| **`ENTITY.RAYCAST(ox, oy, oz, dx, dy, dz, maxDist)`** → **entity** | First hit entity along the ray segment (**Jolt** query path shared with **`PICK.*`** / **`PickCastEntityID`**). Returns **0** if none. |
| **`ENTITY.GETGROUNDNORMAL(id)`** → **vec3 handle** | With **`PLAYER.CREATE`**, uses **`CharacterVirtual.GetGroundNormal`**; otherwise a short downward Jolt ray. Fallback normal **`(0,1,0)`** if no hit. |
| **`ENTITY.APPLYIMPULSE(id, fx, fy, fz)`** | Same as **`ENTITY.ADDFORCE`** / **`ApplyEntityImpulse`** (host velocity integration). Not **`BodyInterface::AddImpulse`** until the Jolt C wrapper exposes it. |
| **`ENTITY.CANSEE(observer, target, fovDeg, maxDist)`** → **bool** | Vision cone + line-of-sight: forward from observer eye height (**~1.65**), aim at target at the same offset, **`PickCastEntityID`** along that segment must hit **target** (or no physics hit). |
| **`ENTITY.GETCLOSESTWITHTAG(id, radius, tag)`** → **entity** | Same tag rules as **`PLAYER.GETNEARBY`**, but only the **nearest** match (**0** if none). |
| **`ENTITY.PUSHOUTOFGEOMETRY(id)`** | Best-effort depenetration: nudges world **Y** up slightly; full recovery belongs in Jolt body / character settings when exposed. |
| **`ENTITY.HASTAG(id, pattern)`** → **bool** | **`path.Match`** on **Blender `tag`** or **entity name** only (stricter than **`ENTITY.ISTYPE`**, which also checks metadata **type** fields). **`EntityHasTag`** alias. |
| **`ENTITY.INFRUSTUM(id)`** → **bool** | Same frustum test as **`ENTITY.INVIEW`**, but uses the **active** **`CAMERA.BEGIN`** camera (no camera handle argument). Returns **false** outside a Begin/End 3D block. |
| **`ENTITY.LINEOFSIGHT(observer, target)`** → **bool** | Straight segment from observer eye (~**1.65** m) to target “eye” height — first **Jolt** hit must be **target** (or no physics hit). Does **not** skip trigger/sensor bodies until those use collision layers / masks. |
| **`ENTITY.GETOVERLAPCOUNT(zoneId, tag)`** → **int** | Counts entities matching **`tag`** (same glob rules as **`PLAYER.GETNEARBY`**) whose **pivot** lies inside **zoneId**’s world **AABB** (sphere prefilter + axis test). |
| **`ENTITY.ANIMATETOWARD(id, x, y, z, duration)`** | Linear **world** lerp of the entity root to **(x,y,z)** over **duration** seconds (driven inside **`ENTITY.UPDATE(dt)`**). |

Detailed normals along an arbitrary ray: **`PHYSICS3D.RAYCAST`** (returns a small result array including the surface normal).

## Movement, rule collisions, and space transforms

These globals mirror Blitz-style names; canonical forms are **`ENTITY.*`** / **`MOVEENTITY`** where noted.

### `ENTITY.GETPOS(entity)` → **handle** 

- **Arguments:** `entity` (int entity id).
- **Returns:** a 3-float tuple-like array handle `[x, y, z]` for destructuring.
- **Use case:** convenient one-call position reads in game loops.

```basic
px, py, pz = ENTITY.GETPOS(player)
```

`ENTITY.GETPOSITION(entity)` remains available and returns a vec3 handle for the handle-based vector API.

---

### `MoveEntity(entity, forward, right, up)` 

- **Arguments:** **`entity`** — entity id; **`forward`**, **`right`**, **`up`** — distances to move along that entity’s **local** axes (from its **pitch** and **yaw**; roll is not used for the basis).
- **Behavior:** Same as **`MOVEENTITY`** and **`ENTITY.MOVE`**. The engine builds a forward vector from yaw/pitch, derives right from the world up cross forward, then up from right cross forward, and adds **`forward·fwd + right·right + up·up`** to the entity’s **world** position (parent-aware).
- **Use for:** Walking relative to facing (e.g. set **`RotateEntity(player, 0, camYaw, 0)`** then **`MoveEntity(player, speed*dt, 0, 0)`** for forward).
- **Not for:** A fixed world offset — use **`TranslateEntity`** instead.

---

### `TranslateEntity(entity, dx, dy, dz)` 

- **Arguments:** **`entity`**; **`dx`**, **`dy`**, **`dz`** — delta in **world** space (applied to world position, then converted back to local if parented).
- **Behavior:** Same as **`ENTITY.TRANSLATE`** / **`ENTITY.TRANSLATEENTITY`**.
- **Use for:** Nudging lights, props, or anything that should move **`(1,0,0)`** in world axes regardless of rotation.

---

### `EntityHitsType(entity, type)` → **bool** 

- **Arguments:** **`entity`** — mover/query entity; **`type`** — integer **collision type** previously set with **`EntityType`** / **`ENTITY.TYPE`** on **other** entities (e.g. ground = **`2`**).
- **Returns:** **`TRUE`** if, **after the last `ENTITY.UPDATE(dt)`** (or **`UPDATEPHYSICS`**), **`entity`** has a rule-based hit whose other body’s **`EntityType`** equals **`type`**. Otherwise **`FALSE`**.
- **Relation to `ENTITYCOLLIDED`:** Same test as **`ENTITYCOLLIDED(entity, type) <> 0`**; **`ENTITYCOLLIDED`** returns the **other entity’s id** or **`0`** if you need the handle.
- **Prerequisites:** Register pairs with **`COLLISIONS(srcType, dstType, method, response)`** (e.g. sphere-vs-box **`method`** **`2`**) and run **`ENTITY.UPDATE`** each frame. **Not** the same as **`EntityCollided(a, b)`**, which is the **two-entity Jolt** contact query (Linux + linked buffers).

---

### `ENTITY.TFORMVECTOR(x, y, z, srcEntity, dstEntity)` → **handle** 

- **Arguments:** Direction or vector components **`x`**, **`y`**, **`z`** in **`srcEntity`**’s **local** space; **`srcEntity`** and **`dstEntity`** are entity ids.
- **Returns:** **Heap handle** to a **3-element float array** (same convention as **`ENTITY.GETPOSITION`**): read components via array access or helpers your script style supports.
- **Behavior:** Alias of **`ENTITY.TFORMVECTOR`**. Transforms the vector by the **linear** part of the world matrix chain (direction only, no translation).
- **Use for:** Camera-relative directions, wind in ship space, etc. **Note:** There is no **`entity = 0`** “world” shortcut; use an axis-aligned **pivot entity** at the origin if you need world as a space.

---

## Scene hierarchy & world utilities (Blitz-style)

- **`Entity.Visible(entity, visible)`** / **`EntityVisible`** — sets the same flag as **`ENTITY.HIDE`** / **`ENTITY.SHOW`** (`visible` = false hides the entity).
- **`Entity.CountChildren(parent)`** — number of **direct** children (stable order = reparent / create order).
- **`Entity.GetChild(parent, index)`** — direct child entity at `index` (0-based).
- **`Entity.FindChild(rootEntity, name)`** — breadth-first search **under** `root` (not global; use **`ENTITY.FIND`** for global name lookup). Names come from **`ENTITY.SETNAME`**.
- **`Entity.TFormPoint(x, y, z, srcEntity, dstEntity)`** / **`Entity.TFormVector(...)`** — same semantics as **`TFormVector`** / **`ENTITY.TFORMVECTOR`** above; **`TFORMPOINT`** includes translation (full matrix); **`TFORMVECTOR`** is direction-only. Returns a **3-float numeric array handle** (same pattern as **`ENTITY.GETPOSITION`**).
- **`Entity.DeltaX`** / **`DeltaY`** / **`DeltaZ(entityA, entityB)`** — world-space axis delta **B − A** between origins.
- **`Entity.MatrixElement(entity, row, col)`** — one element of the **world** matrix; **row/col 0..3**, **column-major** (same as **`MAT4.GETELEMENT`** / Raylib `rl.Matrix`).
- **`Entity.InView(entity, camera)`** — conservative frustum test for the entity bounds vs the given **`CAMERA.CREATE`** handle (aspect from current framebuffer). **`Entity.SetCullMode`** force visible/hidden still applies first.

## 3D sprites (billboards)

- **`LoadSprite(path)`** / **`LoadSprite(path, parent)`** — **`ENTITY.LOADSPRITE`** / **`ENTITY.CREATESPRITE`** are aliases; optional **parent** parents the new sprite like **`Entity.Parent`** (child starts at local origin).
- **`SpriteMode`** / **`Entity.SpriteViewMode`** / **`SpriteViewMode`** — **`1`** = Y-axis billboard, **`2`** = full camera-facing billboard, **`3`** = static quad (see implementation in [`entity_cgo.go`](../../runtime/mbentity/entity_cgo.go)).

## Bulk free (`FreeEntities` / `Entity.FreeEntities`)

**`FreeEntities(arrayHandle)`** walks a **numeric** entity array (e.g. **`DIM badGuy AS HANDLE(n)`** / integer slots holding entity ids) and calls **`FreeEntity`** on each non-zero entry. Use at shutdown or level unload instead of hand-written **`FOR i = 1 TO n : FreeEntity(...) : NEXT`**.

## Terrain vs entity

Heightmap **terrain** is a separate **`TERRAIN.*`** heap object ([TERRAIN.md](TERRAIN.md)), not an entity. Use **`Terrain.GetHeight`** for height queries, **`Terrain.Place`** to position an entity on the surface in one call, or **`Terrain.SnapY`** to adjust **Y** only. **`Terrain.Raise`** / **`Terrain.Lower`** edit heights. **Jolt** heightfield shapes for terrain are **not** exposed in the current `jolt-go` binding; physics for terrain remains mesh/other shapes until extended.

## Performance / roadmap notes

- **Stencil mirrors** (reflection planes) are not implemented yet.
- **Heavy billboard counts:** many **`LOADSPRITE`** instances may benefit from future batching in **`ENTITY.DRAWALL`**; profile hot paths first.

## Procedural meshes (`ENTITY.CREATEMESH`)

- **`ENTITY.CREATEMESH`** / **`CreateMesh([parent])`** — allocates a **blank** procedural mesh (no default cube). The entity stays **hidden** until **`UpdateMesh`**. Optional **parent** works like **`CreateCube(parent)`**.
- **`CreateSurface` / `ENTITY.CREATESURFACE`** — returns the **surface handle** (same as the internal **`MeshBuilder`** heap object) used by **`AddVertex`** / **`AddTriangle`**.
- **`AddVertex`**, **`AddTriangle`**, **`UpdateMesh`**, **`VertexX` / `Y` / `Z`** — CPU-side builder + GPU upload via **`LoadModelFromMesh`** (smooth normals from triangle fans). **`ENTITY.FREE`** unloads the model and frees the builder.
- **`EmitSound(sound, entity)`** — plays a one-shot sound with **distance attenuation** and **stereo pan** vs the last **`Listener(cam)`** / **`AUDIO.LISTENERCAMERA`** (see [AUDIO.md](AUDIO.md) spatial notes).

## Skeletal animation, bone sockets, materials (Raylib)

Full **command matrix** (entities vs **`MODEL.*`**), time scaling, and limitations: **[ANIMATION_3D.md](ANIMATION_3D.md)**.

- **`ENTITY.LOADANIMATEDMESH` / `LoadAnimMesh`** — loads a model and **`LoadModelAnimations`**; first pose is applied with **`UpdateModelAnimation`** + **`UpdateModelAnimationBones`**.
- **`ENTITY.ANIMATE` / `Animate`** — **`ENTITY.ANIMATE(entity [, mode, speed])`**: mode **`0`–`1`** = loop, **`2`** = ping-pong, **`3`+** = clamp at end of clip. (Older scripts used mode **`1`** for clamp; use **`3`** now.) **100 ms skeletal cross-fade** between clips is not implemented (single active pose from Raylib).
- **`ENTITY.EXTRACTANIMSEQ` / `ExtractAnimSeq`** — **`(entity, startFrame, endFrame)`** inclusive clip range for the **current** animation; use **`ENTITY.SETANIMINDEX`** to pick which **`ModelAnimation`** is active.
- **`ENTITY.SETANIMINDEX`** — select animation clip index (resets time to 0). **`ENTITY.ANIMINDEX`** / **`ENTITY.ANIMCOUNT`** — read active index and loaded clip count.
- **`ENTITY.FINDBONE` / `FindBone`** — returns a **hidden** entity whose **world matrix** tracks a named bone on the host model each frame (parent props with **`ENTITY.LOADMESH(path, parent)`** or **`ENTITY.PARENT`**). If the host model is freed, sockets are invalidated.
- **`ENTITY.SETANIMTIME` / `SetAnimTime`**, **`ENTITY.ANIMTIME` / `EntityAnimTime`** — continuous animation time (not always an integer frame index).
- **Brushes:** **`CreateBrush`**, **`BrushTexture`**, **`BrushFX`**, **`BrushShininess`**, **`PaintEntity`** — heap **`Brush`** handle; **`PaintEntity`** copies color/texture/FX/shininess onto the entity. **`BrushFX`**: bit **`1`** = full-bright tint boost in **`ENTITY.DRAWALL`**, bit **`16`** = additive blending (Raylib **`BlendAdditive`**). Full PBR/shader swaps are future work.
- **`EntityShadow`** — stores **`shadowCast`** on the entity (**`0`** default, **`1`** / **`2`** reserved); hooking into the deferred shadow pass for **`ENTITY.DRAWALL`** models is not wired yet.

## Reference tables

- **[API_CONSISTENCY.md](../API_CONSISTENCY.md)** — search for **`ENTITY.`** for every overload and arity.
- **[ANIMATION_3D.md](ANIMATION_3D.md)** — skeletal clips, **`ENTITY.UPDATE`**, bone sockets.
- **[GAMEHELPERS.md](GAMEHELPERS.md)** — movement, landing, camera follow.

---

## Full Example

A spinning cube entity with camera orbit.

```basic
WINDOW.OPEN(960, 540, "Entity Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

cube = ENTITY.CREATECUBE(2.0)
ENTITY.SETPOS(cube, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    ENTITY.UPDATE(dt)
    ENTITY.TURN(cube, 0, 30.0 * dt, 0)

    RENDER.CLEAR(20, 30, 50)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(cube)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

Commands not covered in detail above. All accept an entity id or handle as the first argument unless noted.

### Creation aliases

| Command | Description |
|--------|-------------|
| `ENTITY.MAKE(path)` | Alias of `ENTITY.LOAD`. |
| `ENTITY.CREATEENTITY()` / `ENTITY.MAKEENTITY()` | Create an empty entity with no mesh. |
| `ENTITY.CREATEBOX(w,h,d)` / `ENTITY.MAKEBOX(w,h,d)` | Box mesh entity. |
| `ENTITY.MAKECUBE(size)` | Cube entity (uniform box). |
| `ENTITY.CREATECONE(r,h)` / `ENTITY.MAKECONE(r,h)` | Cone mesh entity. |
| `ENTITY.CREATECYLINDER(r,h)` / `ENTITY.MAKECYLINDER(r,h)` | Cylinder mesh entity. |
| `ENTITY.MAKESPHERE(r)` | Sphere mesh entity. |
| `ENTITY.MAKEMESH(meshHandle)` | Entity from an existing mesh handle. |
| `ENTITY.MAKESURFACE()` | Empty surface entity for dynamic geometry. |
| `ENTITY.MAKEPLANE(w,d)` | Flat plane entity. |
| `ENTITY.MAKESPRITE(texHandle)` | Billboard sprite entity. |
| `ENTITY.COPY(id)` | Clone an entity (same mesh, new id). |

---

### Position / pose shortcuts

| Command | Description |
|--------|-------------|
| `ENTITY.ENTITYX(id)` / `ENTITYY` / `ENTITYZ` | World X/Y/Z (same as `ENTITY.X/Y/Z`). |
| `ENTITY.ENTITYPITCH(id)` / `ENTITYYAW` / `ENTITYROLL` | Euler rotation accessors. |
| `ENTITY.GETROT(id)` | Returns `[pitch, yaw, roll]` array. |
| `ENTITY.GETSCALE(id)` | Returns `[sx, sy, sz]` array. |
| `ENTITY.GETXZ(id)` | Returns `[x, z]` ground-plane position array. |
| `ENTITY.FLOOR(id)` | Snaps entity Y to the terrain height at its XZ. |
| `ENTITY.SNAPTO(id, targetId)` | Teleports entity to target's world position. |
| `ENTITY.CLAMPTOTERRAIN(id)` | Keeps entity Y at or above terrain surface. |
| `ENTITY.SETROTATION(id, p, y, r)` | Alias of `ENTITY.SETROT`. |
| `ENTITY.TURNENTITY(id, dp, dy, dr)` | Incremental rotation each frame. |

---

### Movement helpers

| Command | Description |
|--------|-------------|
| `ENTITY.MOVETOWARD(id, tx, ty, tz, speed)` | Move entity toward target position at speed. |
| `ENTITY.MOVERELATIVE(id, dx, dz)` | Move relative to entity's own yaw. |
| `ENTITY.MOVECAMERARELATIVE(id, cam, f, s, speed)` | Camera-relative WASD movement. |
| `ENTITY.MOVEWITHCAMERA(id, cam, f, s, speed)` | Alias of `MOVECAMERARELATIVE`. |
| `ENTITY.LOOKAT(id, tx, ty, tz)` | Face toward world position. |
| `ENTITY.POINTAT(id, targetId)` | Face toward another entity. |
| `ENTITY.POINTENTITY(id, targetId)` | Alias of `POINTAT`. |
| `ENTITY.TURNTOWARD(id, targetId, speed)` | Smoothly rotate toward target entity at speed. |
| `ENTITY.ALIGNTOVECTOR(id, nx, ny, nz)` | Align entity up-axis to a normal vector. |
| `ENTITY.NAVTO(id, x, z)` | Set navigation destination (uses navmesh). |
| `ENTITY.WANDER(id, radius, speed)` | Random wandering within radius. |
| `ENTITY.FLEE(id, fromId, speed)` | Move away from target entity. |
| `ENTITY.PATROL(id, ...)` | Follow a waypoint patrol route. |
| `ENTITY.MAGNETTO(id, targetId, force)` | Pull toward target with a force. |
| `ENTITY.SLIDE(id, dx, dz)` | Slide entity along XZ with wall collision response. |
| `ENTITY.JUMP(id, impulse)` | Apply upward velocity impulse. |
| `ENTITY.CUTJUMP(id)` | Cancel current jump arc (for variable jump height). |
| `ENTITY.WASGROUNDED(id)` | Returns `TRUE` if entity was grounded last frame. |
| `ENTITY.ISWALLSLIDING(id)` | Returns `TRUE` if entity is sliding against a wall. |

---

### Physics & material

| Command | Description |
|--------|-------------|
| `ENTITY.APPLYGRAVITY(id, dt)` | Apply gravity step to entity velocity. |
| `ENTITY.APPLYTORQUE(id, tx, ty, tz)` | Apply torque to Jolt body. |
| `ENTITY.PHYSICSMOTION(id, vx, vy, vz)` | Set kinematic body motion target. |
| `ENTITY.SETGRAVITY(id, g)` | Per-entity gravity override. |
| `ENTITY.SETGRAVITYSCALE(id, scale)` | Scale world gravity for this entity. |
| `ENTITY.SETMASS(id, mass)` | Set Jolt body mass. |
| `ENTITY.SETBOUNCE(id, r)` / `ENTITY.SETBOUNCINESS(id, r)` | Set restitution. |
| `ENTITY.SETFRICTION(id, f)` | Set surface friction. |
| `ENTITY.SETBUOYANCY(id, b)` | Set buoyancy factor for water interaction. |
| `ENTITY.GETBUOYANCY(id)` | Get current buoyancy setting. |
| `ENTITY.ISSUBMERGED(id)` | Returns `TRUE` if entity is in a water volume. |
| `ENTITY.SETSTATIC(id, bool)` | Mark entity as static for collision baking. |
| `ENTITY.SETTRIGGER(id, bool)` | Mark entity as a trigger (sensor) body. |
| `ENTITY.PICKMODE(id, mode)` | Set Jolt pick/query mode. |
| `ENTITY.SETWEIGHT(id, w)` | Gameplay weight (affects `PLAYER.PUSH`). |

---

### Collision queries

| Command | Description |
|--------|-------------|
| `ENTITY.COLLISIONX(id)` / `COLLISIONY` / `COLLISIONZ` | Last collision hit point. |
| `ENTITY.COLLISIONNX(id)` / `COLLISIONNY` / `COLLISIONNZ` | Last collision normal. |
| `ENTITY.COLLISIONOTHER(id)` | Entity id of last collision partner. |
| `ENTITY.CHECKRADIUS(id, r)` | Returns closest entity id within radius `r`. |
| `ENTITY.DISTANCETO(id, otherId)` | World distance to another entity. |
| `ENTITY.GETDISTANCE(id, otherId)` | Alias of `DISTANCETO`. |
| `ENTITY.WITHINRADIUS(id, otherId, r)` | Returns `TRUE` if distance ≤ r. |
| `ENTITY.ENTITIESINRADIUS(id, r)` | Array of entity ids within radius. |
| `ENTITY.ENTITIESINBOX(cx,cy,cz, hw,hh,hd)` | Array of entity ids inside AABB. |
| `ENTITY.ENTITIESINGROUP(group)` | Array of ids in a named group. |
| `ENTITY.FINDBYPROPERTY(key, value)` | Find entities matching a metadata property. |

---

### Animation extensions

| Command | Description |
|--------|-------------|
| `ENTITY.SETANIMATION(id, name)` | Play animation by name. |
| `ENTITY.SETANIMFRAME(id, frame)` | Seek to a specific frame. |
| `ENTITY.SETANIMLOOP(id, bool)` | Enable or disable loop. |
| `ENTITY.SETANIMSPEED(id, speed)` | Set playback speed multiplier. |
| `ENTITY.STOPANIM(id)` | Stop the current animation. |
| `ENTITY.CURRENTANIM(id)` / `ENTITY.CURRENTANIM$(id)` | Returns current animation name. |
| `ENTITY.ANIMNAME(id, index)` / `ENTITY.ANIMNAME$(id, index)` | Returns animation name by index. |
| `ENTITY.ANIMLENGTH(id, name)` | Returns animation duration in seconds. |
| `ENTITY.ISPLAYING(id)` | Returns `TRUE` if an animation is playing. |
| `ENTITY.CROSSFADE(id, name, duration)` | Cross-fade to another animation. |
| `ENTITY.TRANSITION(id, name, blend)` | Transition to animation with blend factor. |
| `ENTITY.GETBONEPOS(id, boneName)` | Returns `[x,y,z]` world position of a bone. |
| `ENTITY.GETBONEROT(id, boneName)` | Returns `[p,y,r]` world rotation of a bone. |

---

### Visual effects & rendering

| Command | Description |
|--------|-------------|
| `ENTITY.ALPHA(id, a)` | Set entity alpha (0.0–1.0). |
| `ENTITY.GETALPHA(id)` | Get current alpha. |
| `ENTITY.GETCOLOR(id)` | Returns `[r,g,b,a]` tint array. |
| `ENTITY.RGB(id, r, g, b)` | Set RGB tint (full alpha). |
| `ENTITY.COLORPULSE(id, r, g, b, speed)` | Animate color pulsing at speed. |
| `ENTITY.OUTLINE(id, r, g, b, thickness)` | Enable/configure outline effect. |
| `ENTITY.GHOSTMODE(id, bool)` | Semi-transparent ghost rendering. |
| `ENTITY.SQUASH(id, sy, speed)` | Squash-and-stretch on Y axis. |
| `ENTITY.ADDTRAIL(id, length, r, g, b, a)` | Attach a motion trail effect. |
| `ENTITY.ADDWOBBLE(id, amp, freq)` | Attach a wobble deformation. |
| `ENTITY.SCROLLMATERIAL(id, ux, uy)` | Scroll UV texture coordinates per frame. |
| `ENTITY.SETDETAILTEXTURE(id, texHandle)` | Assign a detail/overlay texture. |
| `ENTITY.SETTEXTUREMAP(id, slot, texHandle)` | Set texture in a specific material slot. |
| `ENTITY.SETTEXTURESCROLL(id, ux, uy)` | Set persistent UV scroll speed. |
| `ENTITY.SETTEXTUREFLIP(id, flipX, flipY)` | Flip texture UVs. |
| `ENTITY.SETSHADER(id, shaderHandle)` | Assign a custom shader. |
| `ENTITY.SETSPRITEFRAME(id, frame)` | Set sprite sheet frame index. |
| `ENTITY.EMITPARTICLES(id, count)` | Burst-emit from entity's attached emitter. |
| `ENTITY.EXPLODE(id, force, radius)` | Destroy and scatter physics fragments. |
| `ENTITY.INSTANCEGRID(id, cols, rows, spacing)` | Stamp a grid of instances. |

---

### State & messaging

| Command | Description |
|--------|-------------|
| `ENTITY.GETSTATE(id)` | Returns current AI/gameplay state string. |
| `ENTITY.SETHEALTH(id, hp)` | Set entity health value. |
| `ENTITY.ISALIVE(id)` | Returns `TRUE` if health > 0. |
| `ENTITY.DAMAGE(id, amount)` | Apply damage and reduce health. |
| `ENTITY.SETTAG(id, tag)` | Set entity tag string. |
| `ENTITY.GETMETADATA(id, key)` | Read a glTF `extras` metadata property. |
| `ENTITY.SENDMESSAGE(id, msg)` | Post a string message to entity's queue. |
| `ENTITY.POLLMESSAGE(id)` | Pop and return the next queued message. |
| `ENTITY.ONHIT(id, callback)` | Register callback for hit events. |
| `ENTITY.ONDEATHDROP(id, template)` | Spawn `template` entity when entity dies. |
| `ENTITY.ATTACH(id, parentId)` | Parent entity to another (alias of `ENTITY.SETPARENT`). |
| `ENTITY.PARENTCLEAR(id)` | Clear parent, return to world space. |

---

### Groups

| Command | Description |
|--------|-------------|
| `ENTITY.GROUPCREATE(name)` | Create a named entity group. |
| `ENTITY.GROUPADD(id, group)` | Add entity to a named group. |
| `ENTITY.GROUPREMOVE(id, group)` | Remove entity from a named group. |

---

### Visibility & settings

| Command | Description |
|--------|-------------|
| `ENTITY.SETVISIBLE(id, bool)` | Show or hide entity. |

### Surface vertex accessors

| Command | Description |
|--------|-------------|
| `ENTITY.VERTEXY(id, index)` | Returns the Y component of vertex `index` on a surface entity. |
| `ENTITY.VERTEXZ(id, index)` | Returns the Z component of vertex `index` on a surface entity. |

---

## See also

- [ANIMATION_3D.md](ANIMATION_3D.md) — skeletal clips, bone sockets
- [SCENE.md](SCENE.md) — scene switching and entity groups
- [PHYSICS3D.md](PHYSICS3D.md) — physics body bridge
- [GAMEHELPERS.md](GAMEHELPERS.md) — movement helpers
