# Entity commands

**Style:** New scripts should prefer **`ENTITY.SETPOS`**, **`ENTITY.DRAWALL`**, **`WINDOW.*`**, **`RENDER.*`** — see the [API Standardization Directive](../API_STANDARDIZATION_DIRECTIVE.md) and [STYLE_GUIDE.md](../../STYLE_GUIDE.md). Sections below still show **Blitz / Easy Mode** names (`Entity.Position`, `Window.Open`, …) where those facades exist for migration.

moonBASIC **entities** are lightweight **integer ids** (handles to rows in the host entity store—not heavy “OOP objects”). **`ENTITY.*`** is the canonical prefix; **PascalCase** **`Entity.*`** names in this doc are **facade aliases** where registered (see [`entity_blitz_facade_cgo.go`](../../runtime/mbentity/entity_blitz_facade_cgo.go)). **CGO** builds link the full implementation; see [`runtime/mbentity/`](../../runtime/mbentity/).

**Scene vs level:** **`SCENE.*`** is for **game scene** switching (**mbscene**). Loading glTF levels uses **`LEVEL.LOAD`** / **`LEVEL.FINDENTITY`** (or **`ENTITY.FIND`** by name)—not **`Scene.FindEntity`**.

### `Entity.Load(path)`
Loads a 3D model from a file path. Returns an **entity id**. 
- `path`: String path to model file (glTF, OBJ, etc.).

### `Entity.CreateCube(size)`
Creates a cube primitive entity.
- `size`: Uniform size of the cube.

---

### `Entity.Position(id, x, y, z)`
Sets the world position of an entity.
- `id`: Entity id.
- `x, y, z`: World coordinates.

### `Entity.Move(id, forward, right, up)`
Moves the entity along its local axes (from its current pitch and yaw).
- `forward, right, up`: Distances to move along local vectors.

### `Entity.SetRotation(id, pitch, yaw, roll)`
Sets the absolute rotation of an entity in degrees.

### `Entity.Turn(id, pitch, yaw, roll)`
Adds to the current rotation of an entity (delta rotation in degrees).

### `Entity.Scale(id, sx, sy, sz)`
Sets the non-uniform scale of an entity.

---

### `Entity.Parent(child, parent)`
Parents one entity to another. The child inherits the parent's transforms.

### `Entity.Unparent(id)`
Removes the parent from an entity while maintaining its world position.

---

### `Entity.Visible(id, toggle)`
Sets whether the entity is visible.
- `toggle`: Boolean (`TRUE` or `FALSE`).

### `Entity.Free(id)`
Frees the entity and its resources from memory.

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

### `MoveEntity(entity, forward, right, up)`

- **Arguments:** **`entity`** — entity id; **`forward`**, **`right`**, **`up`** — distances to move along that entity’s **local** axes (from its **pitch** and **yaw**; roll is not used for the basis).
- **Behavior:** Same as **`MOVEENTITY`** and **`ENTITY.MOVE`**. The engine builds a forward vector from yaw/pitch, derives right from the world up cross forward, then up from right cross forward, and adds **`forward·fwd + right·right + up·up`** to the entity’s **world** position (parent-aware).
- **Use for:** Walking relative to facing (e.g. set **`RotateEntity(player, 0, camYaw, 0)`** then **`MoveEntity(player, speed*dt, 0, 0)`** for forward).
- **Not for:** A fixed world offset — use **`TranslateEntity`** instead.

### `TranslateEntity(entity, dx, dy, dz)`

- **Arguments:** **`entity`**; **`dx`**, **`dy`**, **`dz`** — delta in **world** space (applied to world position, then converted back to local if parented).
- **Behavior:** Same as **`ENTITY.TRANSLATE`** / **`ENTITY.TRANSLATEENTITY`**.
- **Use for:** Nudging lights, props, or anything that should move **`(1,0,0)`** in world axes regardless of rotation.

### `EntityHitsType(entity, type)` → **bool**

- **Arguments:** **`entity`** — mover/query entity; **`type`** — integer **collision type** previously set with **`EntityType`** / **`ENTITY.TYPE`** on **other** entities (e.g. ground = **`2`**).
- **Returns:** **`TRUE`** if, **after the last `ENTITY.UPDATE(dt)`** (or **`UPDATEPHYSICS`**), **`entity`** has a rule-based hit whose other body’s **`EntityType`** equals **`type`**. Otherwise **`FALSE`**.
- **Relation to `ENTITYCOLLIDED`:** Same test as **`ENTITYCOLLIDED(entity, type) <> 0`**; **`ENTITYCOLLIDED`** returns the **other entity’s id** or **`0`** if you need the handle.
- **Prerequisites:** Register pairs with **`COLLISIONS(srcType, dstType, method, response)`** (e.g. sphere-vs-box **`method`** **`2`**) and run **`ENTITY.UPDATE`** each frame. **Not** the same as **`EntityCollided(a, b)`**, which is the **two-entity Jolt** contact query (Linux + linked buffers).

### `Entity.TFormVector(x, y, z, srcEntity, dstEntity)` → **handle**

- **Arguments:** Direction or vector components **`x`**, **`y`**, **`z`** in **`srcEntity`**’s **local** space; **`srcEntity`** and **`dstEntity`** are entity ids.
- **Returns:** **Heap handle** to a **3-element float array** (same convention as **`ENTITY.GETPOSITION`**): read components via array access or helpers your script style supports.
- **Behavior:** Alias of **`ENTITY.TFORMVECTOR`**. Transforms the vector by the **linear** part of the world matrix chain (direction only, no translation).
- **Use for:** Camera-relative directions, wind in ship space, etc. **Note:** There is no **`entity = 0`** “world” shortcut; use an axis-aligned **pivot entity** at the origin if you need world as a space.

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
